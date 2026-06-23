package channels_test

import (
	"context"
	"testing"

	"github.com/pickeringtech/go-collections/channels"
)

// The windowing benchmarks reuse the package ladder (capped at 10_000): every
// element crosses an unbuffered channel, a goroutine handoff the framework
// cannot amortise, so larger cells dominate CI wall-clock for no extra signal.
// This matches the cap documented in bench_test.go.

func BenchmarkTumblingWindow(b *testing.B) {
	for _, bm := range ladder {
		sli := intSlice(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for range channels.TumblingWindow(benchCtx, channels.FromSlice(benchCtx, sli), 8) {
				}
			}
		})
	}
}

func BenchmarkSlidingWindow(b *testing.B) {
	for _, bm := range ladder {
		sli := intSlice(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for range channels.SlidingWindow(benchCtx, channels.FromSlice(benchCtx, sli), 8, 2) {
				}
			}
		})
	}
}

func BenchmarkSessionWindow(b *testing.B) {
	// A gap every 8th element, so sessions are bounded and emission is regular.
	gap := func(prev, next int) bool { return next%8 != 0 }
	for _, bm := range ladder {
		sli := intSlice(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for range channels.SessionWindow(benchCtx, channels.FromSlice(benchCtx, sli), gap) {
				}
			}
		})
	}
}

func BenchmarkWindowedReduce(b *testing.B) {
	windower := func(ctx context.Context, in <-chan int) <-chan []int {
		return channels.TumblingWindow(ctx, in, 8)
	}
	sum := func(w []int) int {
		total := 0
		for _, v := range w {
			total += v
		}
		return total
	}
	for _, bm := range ladder {
		sli := intSlice(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for range channels.WindowedReduce(benchCtx, channels.FromSlice(benchCtx, sli), windower, sum) {
				}
			}
		})
	}
}
