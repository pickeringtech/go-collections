// Package concurrency provides helpers for running work functions concurrently
// while bounding the number of workers that may be in flight at once.
//
// BlockingWorkLimiter.Run takes the full slice of work up front and blocks until
// it is all done. BackgroundWorkLimiter instead streams work in over time and
// must be driven through its lifecycle in order:
//
//	wl := NewBackgroundWorkLimiter(4)
//	wl.Start()        // open the limiter and start the background worker
//	wl.Add(work)      // hand off work; repeat as many times as needed
//	wl.Stop()         // close the limiter to new work
//	wl.Wait()         // block until all added work has finished
//	errs := wl.Errors()
//
// Calling Add or Stop before Start, or Add after Stop, is a programming error
// and panics rather than blocking forever or panicking obscurely on a nil or
// closed channel. A redundant Start or Stop is a no-op.
//
// # Data-parallel transforms
//
// Map, ForEach and Batch are slice-shaped transforms layered on the same
// bounded primitive, so they never spawn unbounded goroutines:
//
//	squares, err := concurrency.Map(ctx, nums,
//		func(ctx context.Context, n int) (int, error) { return n * n, nil },
//		concurrency.WithConcurrency(8))
//	// squares[i] == nums[i]*nums[i] - output is always order-preserving
//
// Map returns an order-preserving result slice; ForEach is its side-effecting
// counterpart; Batch chunks the input (via slices.Chunk) and processes each
// chunk concurrently. All three thread a context.Context to every work function
// and accept Options:
//
//   - WithConcurrency sets the degree of parallelism (default
//     runtime.GOMAXPROCS(0)).
//   - WithErrorPolicy selects how a failing work function is handled:
//     StopOnError (default) cancels the rest and reports the first error;
//     CollectErrors runs everything and joins all errors; ContinueOnError runs
//     everything and reports none.
//
// A cancelled context always wins over the error policy: it is reported as the
// context error rather than as a work error, because a cancelled run did not
// complete.
package concurrency
