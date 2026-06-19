package concurrency

import (
	"context"
	"errors"
	"runtime"

	"github.com/pickeringtech/go-collections/slices"
)

// MapFunc transforms one input element into one output value, returning an
// error if the work fails. It receives the (possibly cancelled) context so a
// long-running transform can abort early when a sibling fails or the caller
// cancels.
type MapFunc[I, O any] func(context.Context, I) (O, error)

// EachFunc performs a side effect for one input element, returning an error if
// the work fails. Like MapFunc it receives the context so it can observe
// cancellation.
type EachFunc[T any] func(context.Context, T) error

// BatchFunc processes one chunk of input elements, returning an error if the
// work fails. Batches are the consecutive groups produced by slices.Chunk, so
// the slice it receives is a view into the input and must not be retained or
// mutated beyond the call.
type BatchFunc[T any] func(context.Context, []T) error

// ErrorPolicy selects how a parallel transform reacts to a work function that
// returns an error. Context cancellation is handled separately and always
// surfaces regardless of policy - see the package's run helper.
type ErrorPolicy int

const (
	// StopOnError (the default) cancels the remaining work the moment any item
	// fails and reports the first error in input order. Items already in flight
	// run to completion; items not yet started are skipped.
	StopOnError ErrorPolicy = iota
	// CollectErrors runs every item to completion and reports all the errors
	// joined together (via errors.Join), so errors.Is/As can match any of them.
	CollectErrors
	// ContinueOnError runs every item to completion and reports no work error at
	// all - a best-effort mode where the caller only cares about the successful
	// results (or the side effects that did happen).
	ContinueOnError
)

// config holds the resolved options for a parallel transform.
type config struct {
	concurrency int
	policy      ErrorPolicy
}

// Option customises a parallel transform. Options are applied in order, so a
// later option of the same kind wins.
type Option func(*config)

// WithConcurrency sets the maximum number of work functions that may run at
// once. The degree is enforced by the same work-limiter the rest of the package
// uses, which clamps a value below 1 to 1. When unset, the degree defaults to
// runtime.GOMAXPROCS(0).
func WithConcurrency(n int) Option {
	return func(c *config) {
		c.concurrency = n
	}
}

// WithErrorPolicy selects how the transform reacts to a failing work function.
// When unset, the policy defaults to StopOnError.
func WithErrorPolicy(p ErrorPolicy) Option {
	return func(c *config) {
		c.policy = p
	}
}

// newConfig resolves the defaults and then applies the caller's options.
func newConfig(opts []Option) config {
	cfg := config{
		concurrency: runtime.GOMAXPROCS(0),
		policy:      StopOnError,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	// Normalise an out-of-range policy back to the default so run and resolve
	// cannot disagree about it: run only cancels for StopOnError, while resolve
	// treats any unrecognised value as StopOnError, and that split would be a
	// surprising contract for an exported API.
	if cfg.policy < StopOnError || cfg.policy > ContinueOnError {
		cfg.policy = StopOnError
	}
	return cfg
}

// Map applies fn to every element of input concurrently, with parallelism
// bounded by the configured concurrency (default runtime.GOMAXPROCS(0)), and
// returns the results. The output is always order-preserving: output[i] holds
// fn's result for input[i], so the result is deterministic regardless of the
// order work finished in.
//
// The returned slice always has len(input) elements. Positions whose work
// failed or was skipped (because an earlier failure or the context cancelled
// the run) hold the zero value of O - inspect the error, or use a result type
// that distinguishes "unset", before trusting them.
//
// If input is empty or nil the output is an initialised, non-nil empty slice
// and fn is never called. The error follows the configured ErrorPolicy; a
// cancelled context always wins and is returned in preference to any work error.
func Map[I, O any](ctx context.Context, input []I, fn MapFunc[I, O], opts ...Option) ([]O, error) {
	cfg := newConfig(opts)
	output := make([]O, len(input))
	errs := run(ctx, len(input), cfg, func(c context.Context, i int) error {
		o, err := fn(c, input[i])
		if err != nil {
			return err
		}
		output[i] = o
		return nil
	})
	return output, resolve(ctx, errs, cfg.policy)
}

// ForEach applies fn to every element of input concurrently for its side
// effects, with parallelism bounded by the configured concurrency (default
// runtime.GOMAXPROCS(0)). It is the side-effecting counterpart of Map.
//
// fn runs on many goroutines at once, so any shared state it touches must be
// synchronised by the caller. If input is empty or nil, fn is never called and
// the result is nil. The error follows the configured ErrorPolicy; a cancelled
// context always wins and is returned in preference to any work error.
func ForEach[T any](ctx context.Context, input []T, fn EachFunc[T], opts ...Option) error {
	cfg := newConfig(opts)
	errs := run(ctx, len(input), cfg, func(c context.Context, i int) error {
		return fn(c, input[i])
	})
	return resolve(ctx, errs, cfg.policy)
}

// Batch splits input into consecutive chunks of at most size elements (reusing
// slices.Chunk, so ordering and the partial final chunk match it exactly) and
// applies fn to each chunk concurrently, with parallelism bounded by the
// configured concurrency (default runtime.GOMAXPROCS(0)).
//
// Batching amortises per-call overhead for work that is cheaper in bulk - a
// bulk insert, a vectorised feature computation - while still bounding how many
// batches run at once. It is the batched counterpart of ForEach; for a batched
// transform, compose slices.Chunk with Map instead.
//
// If size <= 0, or input is empty or nil, there are no chunks, fn is never
// called and the result is nil. The error follows the configured ErrorPolicy; a
// cancelled context always wins and is returned in preference to any work error.
func Batch[T any](ctx context.Context, input []T, size int, fn BatchFunc[T], opts ...Option) error {
	cfg := newConfig(opts)
	batches := slices.Chunk(input, size)
	errs := run(ctx, len(batches), cfg, func(c context.Context, i int) error {
		return fn(c, batches[i])
	})
	return resolve(ctx, errs, cfg.policy)
}

// run executes n bounded, context-aware tasks and returns their errors indexed
// by task position, so callers can resolve a deterministic first error. task(i)
// performs item i; a task whose index is skipped (because the context is already
// done) is never invoked.
//
// The bound is provided by BlockingWorkLimiter - the package's existing
// semaphore-plus-WaitGroup primitive - rather than a hand-rolled goroutine loop,
// so the no-unbounded-goroutines guarantee is shared with the rest of the
// package. The limiter's own []error return is ignored: errors are recorded by
// index here instead, which is what makes StopOnError's "first error" and
// CollectErrors' ordering deterministic.
//
// StopOnError cancels an internal child context on the first failure; because
// that child is distinct from the caller's ctx, the cancellation does not make
// ctx.Err() report a context error for what is really a work failure. The
// limiter still dispatches the remaining items, but each returns immediately
// once the child context is done, so no further user work runs.
func run(ctx context.Context, n int, cfg config, task func(context.Context, int) error) []error {
	errs := make([]error, n)
	if n == 0 {
		return errs
	}

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	work := make([]WorkFunc, n)
	for i := 0; i < n; i++ {
		i := i
		work[i] = func() error {
			// Skip once the context is done - either the caller cancelled, or
			// StopOnError cancelled us after an earlier failure.
			if cctx.Err() != nil {
				return nil
			}
			err := task(cctx, i)
			if err != nil {
				errs[i] = err
				if cfg.policy == StopOnError {
					cancel()
				}
			}
			// Always report success to the limiter: its []error return is
			// unused (errors are tracked by index above), so returning the real
			// error would only push it onto the limiter's slow lock+append path.
			return nil
		}
	}

	_ = NewBlockingWorkLimiter(cfg.concurrency).Run(work)
	return errs
}

// resolve turns the per-index errors into a single error according to policy.
// A cancelled caller context takes precedence over any work error and over the
// policy, because a cancelled context means the run did not complete - a
// different fact from "some item failed".
func resolve(ctx context.Context, errs []error, policy ErrorPolicy) error {
	ctxErr := ctx.Err()
	if ctxErr != nil {
		return ctxErr
	}
	switch policy {
	case ContinueOnError:
		return nil
	case CollectErrors:
		// errors.Join skips nils and returns nil when every entry is nil.
		return errors.Join(errs...)
	default: // StopOnError
		for _, err := range errs {
			if err != nil {
				return err
			}
		}
		return nil
	}
}
