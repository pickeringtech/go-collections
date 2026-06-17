package channels_test

import (
	"context"
	"testing"

	"github.com/pickeringtech/go-collections/channels"
	"github.com/pickeringtech/go-collections/maps"
	"github.com/pickeringtech/go-collections/slices"
)

// benchCtx is the never-cancelled context every channel benchmark threads through the functions under test. The
// cancellation machinery only ever evaluates the not-yet-done branch of each select, so its cost is captured in the
// measurement without altering the produce-and-drain behaviour being benchmarked.
var benchCtx = context.Background()

// This file backfills the Benchmark leg of the Example+Test+Benchmark trio for
// every public function in the channels package (issue #52). It follows the
// scaling-ladder pattern from agent-os/standards/testing/benchmark-scaling.md:
// each benchmark sub-benchmarks across a fixed size set via b.Run, assigning the
// result to _ (or fully draining the output channel) so nothing is optimised
// away.
//
// Unlike the slices/maps suites, the channel ladder is capped at 10_000. Every
// element here crosses an unbuffered channel — a goroutine handoff that costs
// hundreds of nanoseconds and that the framework cannot amortise — so the 100k
// and 1M cells would dominate CI wall-clock for no extra signal. This mirrors
// the deliberate caps the collections mutate benchmarks already use.
var ladder = []struct {
	name string
	n    int
}{
	{"3 elements", 3},
	{"10 elements", 10},
	{"100 elements", 100},
	{"1_000 elements", 1_000},
	{"10_000 elements", 10_000},
}

// intSlice returns the values 0..n-1, the source data every channel benchmark
// streams through the function under test.
func intSlice(n int) []int {
	return slices.Generate(n, slices.NumericIdentityGenerator[int])
}

// intEntries returns n key-value entries, used by the entry-oriented benchmarks
// (BuildMapFromEntries) that take a slice rather than a channel.
func intEntries(n int) []maps.Entry[int, int] {
	entries := make([]maps.Entry[int, int], n)
	for i := 0; i < n; i++ {
		entries[i] = maps.Entry[int, int]{Key: i, Value: i}
	}
	return entries
}

func BenchmarkFromSlice(b *testing.B) {
	for _, bm := range ladder {
		sli := intSlice(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// FromSlice is lazy, so drain the channel to measure the full
				// produce-and-close cost.
				for range channels.FromSlice(benchCtx, sli) {
				}
			}
		})
	}
}

func BenchmarkFromMap(b *testing.B) {
	for _, bm := range ladder {
		m := make(map[int]int, bm.n)
		for _, k := range intSlice(bm.n) {
			m[k] = k
		}
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for range channels.FromMap(benchCtx, m) {
				}
			}
		})
	}
}

func BenchmarkCollectAsSlice(b *testing.B) {
	for _, bm := range ladder {
		sli := intSlice(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = channels.CollectAsSlice(channels.FromSlice(benchCtx, sli))
			}
		})
	}
}

func BenchmarkCollectNAsSlice(b *testing.B) {
	for _, bm := range ladder {
		sli := intSlice(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = channels.CollectNAsSlice(channels.FromSlice(benchCtx, sli), bm.n)
			}
		})
	}
}

func BenchmarkCollectAsMap(b *testing.B) {
	for _, bm := range ladder {
		sli := intSlice(bm.n)
		fn := func(in int) maps.Entry[int, int] {
			return maps.Entry[int, int]{Key: in, Value: in}
		}
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = channels.CollectAsMap(channels.FromSlice(benchCtx, sli), fn)
			}
		})
	}
}

func BenchmarkBuildMapFromEntries(b *testing.B) {
	for _, bm := range ladder {
		entries := intEntries(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = channels.BuildMapFromEntries(entries)
			}
		})
	}
}

func BenchmarkMap(b *testing.B) {
	for _, bm := range ladder {
		sli := intSlice(bm.n)
		fn := func(in int) int { return in * 2 }
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for range channels.Map(benchCtx, channels.FromSlice(benchCtx, sli), fn) {
				}
			}
		})
	}
}

func BenchmarkFilter(b *testing.B) {
	for _, bm := range ladder {
		sli := intSlice(bm.n)
		fn := func(in int) bool { return in%2 == 0 }
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for range channels.Filter(benchCtx, channels.FromSlice(benchCtx, sli), fn) {
				}
			}
		})
	}
}

func BenchmarkReduce(b *testing.B) {
	for _, bm := range ladder {
		sli := intSlice(bm.n)
		fn := func(acc, el int) int { return acc + el }
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Reduce emits a single accumulated value then closes its output
				// channel. Drain with range (rather than a single receive) so the
				// Reduce goroutine has fully exited before the next iteration,
				// avoiding cross-iteration overlap in the measurement.
				for range channels.Reduce(benchCtx, channels.FromSlice(benchCtx, sli), fn) {
				}
			}
		})
	}
}

func BenchmarkNewPipeline(b *testing.B) {
	for _, bm := range ladder {
		sli := intSlice(bm.n)
		stage := func(ctx context.Context, input <-chan int) <-chan int {
			return channels.Filter(ctx, input, func(in int) bool { return in%2 == 0 })
		}
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = channels.NewPipeline(benchCtx, channels.FromSlice(benchCtx, sli), stage).CollectAsSlice()
			}
		})
	}
}
