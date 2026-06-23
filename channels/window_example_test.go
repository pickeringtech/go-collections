package channels_test

import (
	"context"
	"fmt"

	"github.com/pickeringtech/go-collections/channels"
	"github.com/pickeringtech/go-collections/stats"
)

func ExampleTumblingWindow() {
	ctx := context.Background()
	input := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5, 6, 7})
	// Fixed-width windows of 3; the trailing partial window [7] is dropped.
	windows := channels.TumblingWindow(ctx, input, 3)

	for w := range windows {
		fmt.Println(w)
	}
	// Output:
	// [1 2 3]
	// [4 5 6]
}

func ExampleSlidingWindow() {
	ctx := context.Background()
	input := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5})
	// Width 3, step 1: overlapping windows, full windows only.
	windows := channels.SlidingWindow(ctx, input, 3, 1)

	for w := range windows {
		fmt.Println(w)
	}
	// Output:
	// [1 2 3]
	// [2 3 4]
	// [3 4 5]
}

func ExampleSessionWindow() {
	ctx := context.Background()
	input := channels.FromSlice(ctx, []int{1, 2, 10, 11, 30})
	// A gap opens when consecutive values differ by more than 5; the open
	// session is flushed when the input closes.
	gap := func(prev, next int) bool { return next-prev <= 5 }
	windows := channels.SessionWindow(ctx, input, gap)

	for w := range windows {
		fmt.Println(w)
	}
	// Output:
	// [1 2]
	// [10 11]
	// [30]
}

func ExampleWindowedReduce() {
	ctx := context.Background()
	input := channels.FromSlice(ctx, []float64{1, 2, 3, 4, 5, 6})
	// Cut the stream into tumbling windows of 3, then reduce each to its mean by
	// composing a stats aggregate directly.
	windower := func(ctx context.Context, in <-chan float64) <-chan []float64 {
		return channels.TumblingWindow(ctx, in, 3)
	}
	means := channels.WindowedReduce(ctx, input, windower, func(w []float64) float64 {
		m, _ := stats.Mean(w)
		return m
	})

	for m := range means {
		fmt.Println(m)
	}
	// Output:
	// 2
	// 5
}
