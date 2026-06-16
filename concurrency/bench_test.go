package concurrency_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/concurrency"
)

// This file backfills the Benchmark leg of the Example+Test+Benchmark trio for
// the public constructors of the concurrency package (issue #52). Each
// benchmark sub-benchmarks across a ladder of work-item counts via b.Run,
// driving a fixed worker limit so the timed variable is the amount of work, not
// the pool size.
//
// The ladder is capped at 10_000 because every work item is dispatched on its
// own goroutine: the scheduler cost the limiters exist to bound is real
// per-item work the framework cannot amortise, so larger cells would dominate
// CI wall-clock for no extra signal. This mirrors the deliberate caps elsewhere
// in the suite.
const benchWorkerLimit = 4

var workLadder = []struct {
	name string
	n    int
}{
	{"3 items", 3},
	{"10 items", 10},
	{"100 items", 100},
	{"1_000 items", 1_000},
	{"10_000 items", 10_000},
}

// noopWork returns n trivial work functions, so the benchmark measures the
// limiter's dispatch and synchronisation overhead rather than the work itself.
func noopWork(n int) []concurrency.WorkFunc {
	work := make([]concurrency.WorkFunc, n)
	for i := range work {
		work[i] = func() error { return nil }
	}
	return work
}

func BenchmarkNewBlockingWorkLimiter(b *testing.B) {
	for _, bm := range workLadder {
		work := noopWork(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				limiter := concurrency.NewBlockingWorkLimiter(benchWorkerLimit)
				_ = limiter.Run(work)
			}
		})
	}
}

func BenchmarkNewBackgroundWorkLimiter(b *testing.B) {
	for _, bm := range workLadder {
		work := noopWork(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				limiter := concurrency.NewBackgroundWorkLimiter(benchWorkerLimit)
				limiter.Start()
				for _, w := range work {
					limiter.Add(w)
				}
				limiter.Stop()
				limiter.Wait()
			}
		})
	}
}
