package concurrency_test

import (
	"context"
	"testing"

	"github.com/pickeringtech/go-collections/concurrency"
)

// These benchmarks reuse the workLadder and benchWorkerLimit defined in
// bench_test.go: a fixed worker limit so the timed variable is the amount of
// work, not the pool size. The work itself is trivial so the measurement is the
// transform's dispatch, ordering, and synchronisation overhead. The ladder is
// capped at 10_000 for the same reason as the limiter benchmarks: each item is
// dispatched on its own goroutine.

// intInput builds a deterministic slice of n ints for the benchmarks to map
// over, so the benchmark allocation is the transform's, not the input's.
func intInput(n int) []int {
	in := make([]int, n)
	for i := range in {
		in[i] = i
	}
	return in
}

func BenchmarkMap(b *testing.B) {
	for _, bm := range workLadder {
		input := intInput(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = concurrency.Map(context.Background(), input,
					func(_ context.Context, n int) (int, error) { return n * 2, nil },
					concurrency.WithConcurrency(benchWorkerLimit))
			}
		})
	}
}

func BenchmarkForEach(b *testing.B) {
	for _, bm := range workLadder {
		input := intInput(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = concurrency.ForEach(context.Background(), input,
					func(context.Context, int) error { return nil },
					concurrency.WithConcurrency(benchWorkerLimit))
			}
		})
	}
}

func BenchmarkBatch(b *testing.B) {
	for _, bm := range workLadder {
		input := intInput(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = concurrency.Batch(context.Background(), input, 16,
					func(context.Context, []int) error { return nil },
					concurrency.WithConcurrency(benchWorkerLimit))
			}
		})
	}
}
