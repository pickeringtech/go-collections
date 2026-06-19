package concurrency_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/pickeringtech/go-collections/concurrency"
)

// errBoom is a sentinel returned by failing work in the tests below.
var errBoom = errors.New("boom")

// doubleOK is the trivial happy-path MapFunc used across cases.
func doubleOK(_ context.Context, n int) (int, error) { return n * 2, nil }

func TestMap(t *testing.T) {
	type args struct {
		input []int
		fn    concurrency.MapFunc[int, int]
		opts  []concurrency.Option
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "doubles every element preserving order",
			args: args{input: []int{1, 2, 3, 4, 5}, fn: doubleOK},
			want: []int{2, 4, 6, 8, 10},
		},
		{
			name: "single worker still preserves order",
			args: args{input: []int{1, 2, 3}, fn: doubleOK, opts: []concurrency.Option{concurrency.WithConcurrency(1)}},
			want: []int{2, 4, 6},
		},
		{
			name: "nil input results in empty output",
			args: args{input: nil, fn: doubleOK},
			want: []int{},
		},
		{
			name: "empty input results in empty output",
			args: args{input: []int{}, fn: doubleOK},
			want: []int{},
		},
		{
			name: "stop on error reports first error by index",
			args: args{
				input: []int{1, 2, 3},
				fn: func(_ context.Context, n int) (int, error) {
					if n == 2 {
						return 0, errBoom
					}
					return n * 2, nil
				},
			},
			// output[1] stays zero (failed); other positions may or may not be
			// filled depending on scheduling, so only the error is asserted.
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := concurrency.Map(context.Background(), tt.args.input, tt.args.fn, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Map() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
			if got == nil {
				t.Errorf("Map() returned nil slice, want non-nil")
			}
		})
	}
}

// TestMapPreservesOrderUnderConcurrency hammers Map with reversed, deliberately
// staggered work to prove output ordering follows input, not completion order.
func TestMapPreservesOrderUnderConcurrency(t *testing.T) {
	const n = 200
	input := make([]int, n)
	for i := range input {
		input[i] = i
	}
	got, err := concurrency.Map(context.Background(), input, func(_ context.Context, i int) (int, error) {
		// Earlier indices yield more often, so they tend to finish last; a
		// completion-ordered implementation would scramble the output.
		for j := 0; j < (n - i); j++ {
			runtime.Gosched()
		}
		return i * i, nil
	}, concurrency.WithConcurrency(8))
	if err != nil {
		t.Fatalf("Map() unexpected error: %v", err)
	}
	for i := range input {
		if got[i] != i*i {
			t.Fatalf("Map() out of order at %d: got %d, want %d", i, got[i], i*i)
		}
	}
}

func TestMapErrorPolicies(t *testing.T) {
	failEven := func(_ context.Context, n int) (int, error) {
		if n%2 == 0 {
			return 0, fmt.Errorf("even %d: %w", n, errBoom)
		}
		return n, nil
	}
	input := []int{1, 2, 3, 4, 5} // evens 2 and 4 fail

	t.Run("collect all joins every error", func(t *testing.T) {
		_, err := concurrency.Map(context.Background(), input, failEven, concurrency.WithErrorPolicy(concurrency.CollectErrors))
		if err == nil {
			t.Fatal("Map() error = nil, want joined errors")
		}
		if !errors.Is(err, errBoom) {
			t.Errorf("Map() error %v does not wrap errBoom", err)
		}
		// errors.Join exposes its parts via Unwrap() []error.
		joined, ok := err.(interface{ Unwrap() []error })
		if !ok {
			t.Fatalf("Map() error %T is not a joined error", err)
		}
		if got := len(joined.Unwrap()); got != 2 {
			t.Errorf("Map() collected %d errors, want 2", got)
		}
	})

	t.Run("continue on error reports no error and fills successes", func(t *testing.T) {
		got, err := concurrency.Map(context.Background(), input, failEven, concurrency.WithErrorPolicy(concurrency.ContinueOnError))
		if err != nil {
			t.Fatalf("Map() error = %v, want nil", err)
		}
		want := []int{1, 0, 3, 0, 5} // failed evens stay zero
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Map() = %v, want %v", got, want)
		}
	})
}

func TestForEach(t *testing.T) {
	t.Run("runs every element exactly once", func(t *testing.T) {
		var sum int64
		input := []int{1, 2, 3, 4, 5}
		err := concurrency.ForEach(context.Background(), input, func(_ context.Context, n int) error {
			atomic.AddInt64(&sum, int64(n))
			return nil
		})
		if err != nil {
			t.Fatalf("ForEach() error = %v, want nil", err)
		}
		if sum != 15 {
			t.Errorf("ForEach() sum = %d, want 15", sum)
		}
	})

	t.Run("nil input never calls fn", func(t *testing.T) {
		called := false
		err := concurrency.ForEach(context.Background(), nil, func(context.Context, int) error {
			called = true
			return nil
		})
		if err != nil {
			t.Fatalf("ForEach() error = %v, want nil", err)
		}
		if called {
			t.Error("ForEach() called fn on empty input")
		}
	})

	t.Run("stop on error surfaces the error", func(t *testing.T) {
		err := concurrency.ForEach(context.Background(), []int{1, 2, 3}, func(_ context.Context, n int) error {
			if n == 2 {
				return errBoom
			}
			return nil
		})
		if !errors.Is(err, errBoom) {
			t.Errorf("ForEach() error = %v, want errBoom", err)
		}
	})
}

func TestBatch(t *testing.T) {
	t.Run("processes consecutive chunks covering every element", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		var mu sync.Mutex
		var seen []int
		sizes := map[int]int{}
		err := concurrency.Batch(context.Background(), input, 2, func(_ context.Context, batch []int) error {
			mu.Lock()
			defer mu.Unlock()
			sizes[len(batch)]++
			seen = append(seen, batch...)
			return nil
		})
		if err != nil {
			t.Fatalf("Batch() error = %v, want nil", err)
		}
		sort.Ints(seen)
		if !reflect.DeepEqual(seen, input) {
			t.Errorf("Batch() covered %v, want %v", seen, input)
		}
		// Chunk([1 2 3 4 5], 2) -> two batches of 2, one of 1.
		if sizes[2] != 2 || sizes[1] != 1 {
			t.Errorf("Batch() chunk sizes = %v, want two of size 2 and one of size 1", sizes)
		}
	})

	t.Run("non-positive size processes nothing", func(t *testing.T) {
		called := false
		err := concurrency.Batch(context.Background(), []int{1, 2, 3}, 0, func(context.Context, []int) error {
			called = true
			return nil
		})
		if err != nil {
			t.Fatalf("Batch() error = %v, want nil", err)
		}
		if called {
			t.Error("Batch() called fn with non-positive size")
		}
	})

	t.Run("collect errors aggregates failing batches", func(t *testing.T) {
		err := concurrency.Batch(context.Background(), []int{1, 2, 3, 4}, 1, func(_ context.Context, batch []int) error {
			if batch[0]%2 == 0 {
				return errBoom
			}
			return nil
		}, concurrency.WithErrorPolicy(concurrency.CollectErrors))
		if !errors.Is(err, errBoom) {
			t.Errorf("Batch() error = %v, want errBoom", err)
		}
	})
}

// TestContextCancellation proves a cancelled context wins over the policy and is
// surfaced as the context error, and that no further work starts after cancel.
func TestContextCancellation(t *testing.T) {
	t.Run("pre-cancelled context skips all work", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		var ran int64
		err := concurrency.ForEach(ctx, []int{1, 2, 3}, func(context.Context, int) error {
			atomic.AddInt64(&ran, 1)
			return nil
		}, concurrency.WithConcurrency(1))
		if !errors.Is(err, context.Canceled) {
			t.Errorf("ForEach() error = %v, want context.Canceled", err)
		}
		if ran != 0 {
			t.Errorf("ForEach() ran %d items after cancel, want 0", ran)
		}
	})

	t.Run("cancellation wins over continue policy", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := concurrency.Map(ctx, []int{1, 2, 3}, doubleOK, concurrency.WithErrorPolicy(concurrency.ContinueOnError))
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Map() error = %v, want context.Canceled even under ContinueOnError", err)
		}
	})

	t.Run("mid-run cancellation stops further dispatch", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		var ran int64
		// One worker, so items run strictly one at a time; the first cancels.
		err := concurrency.ForEach(ctx, []int{1, 2, 3, 4, 5}, func(_ context.Context, n int) error {
			atomic.AddInt64(&ran, 1)
			cancel()
			return nil
		}, concurrency.WithConcurrency(1))
		if !errors.Is(err, context.Canceled) {
			t.Errorf("ForEach() error = %v, want context.Canceled", err)
		}
		if got := atomic.LoadInt64(&ran); got != 1 {
			t.Errorf("ForEach() ran %d items, want exactly 1 before cancel took hold", got)
		}
	})
}

// TestConcurrencyIsBounded confirms the configured degree caps in-flight work.
func TestConcurrencyIsBounded(t *testing.T) {
	const limit = 4
	var inFlight, peak int64
	input := make([]int, 64)
	err := concurrency.ForEach(context.Background(), input, func(_ context.Context, _ int) error {
		cur := atomic.AddInt64(&inFlight, 1)
		for {
			p := atomic.LoadInt64(&peak)
			if cur <= p || atomic.CompareAndSwapInt64(&peak, p, cur) {
				break
			}
		}
		runtime.Gosched()
		atomic.AddInt64(&inFlight, -1)
		return nil
	}, concurrency.WithConcurrency(limit))
	if err != nil {
		t.Fatalf("ForEach() error = %v, want nil", err)
	}
	if peak > limit {
		t.Errorf("ForEach() peak concurrency %d exceeded limit %d", peak, limit)
	}
}
