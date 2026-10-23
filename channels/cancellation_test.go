package channels_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/pickeringtech/go-collections/channels"
)

// TestStagesDoNotReadPastCancellation pins that cancellation takes priority over
// a ready input in every stage. A bare select over {<-ctx.Done(), <-input} picks
// a ready case at random, so a stage whose context was already cancelled used to
// consume (and silently drop) a buffered element - running fn on it - about half
// the time. Each case starts a stage on an already-cancelled context with one
// buffered element and repeats it enough times that a coin-flip select would
// lose at least once.
func TestStagesDoNotReadPastCancellation(t *testing.T) {
	const iterations = 200

	type stage func(ctx context.Context, input <-chan int, calls *atomic.Int64) (closed func(*testing.T))
	tests := []struct {
		name  string
		start stage
	}{
		{
			name: "Map",
			start: func(ctx context.Context, input <-chan int, calls *atomic.Int64) func(*testing.T) {
				out := channels.Map(ctx, input, func(v int) int {
					calls.Add(1)
					return v
				})
				return func(t *testing.T) { assertChannelCloses(t, out, "Map") }
			},
		},
		{
			name: "Filter",
			start: func(ctx context.Context, input <-chan int, calls *atomic.Int64) func(*testing.T) {
				out := channels.Filter(ctx, input, func(int) bool {
					calls.Add(1)
					return true
				})
				return func(t *testing.T) { assertChannelCloses(t, out, "Filter") }
			},
		},
		{
			name: "Reduce",
			start: func(ctx context.Context, input <-chan int, calls *atomic.Int64) func(*testing.T) {
				out := channels.Reduce(ctx, input, func(acc, v int) int {
					calls.Add(1)
					return acc + v
				})
				return func(t *testing.T) { assertChannelCloses(t, out, "Reduce") }
			},
		},
		{
			name: "TumblingWindow",
			start: func(ctx context.Context, input <-chan int, _ *atomic.Int64) func(*testing.T) {
				out := channels.TumblingWindow(ctx, input, 1)
				return func(t *testing.T) { assertChannelCloses(t, out, "TumblingWindow") }
			},
		},
		{
			name: "SlidingWindow",
			start: func(ctx context.Context, input <-chan int, _ *atomic.Int64) func(*testing.T) {
				out := channels.SlidingWindow(ctx, input, 1, 1)
				return func(t *testing.T) { assertChannelCloses(t, out, "SlidingWindow") }
			},
		},
		{
			name: "SessionWindow",
			start: func(ctx context.Context, input <-chan int, calls *atomic.Int64) func(*testing.T) {
				out := channels.SessionWindow(ctx, input, func(int, int) bool {
					calls.Add(1)
					return true
				})
				return func(t *testing.T) { assertChannelCloses(t, out, "SessionWindow") }
			},
		},
		{
			name: "WindowedReduce",
			start: func(ctx context.Context, _ <-chan int, calls *atomic.Int64) func(*testing.T) {
				// WindowedReduce reads windows from its windower, so hand it one
				// already-buffered window to tempt the read.
				windows := make(chan []int, 1)
				windows <- []int{1}
				windower := func(context.Context, <-chan int) <-chan []int { return windows }
				out := channels.WindowedReduce(ctx, nil, windower, func(w []int) int {
					calls.Add(1)
					return len(w)
				})
				return func(t *testing.T) {
					assertChannelCloses(t, out, "WindowedReduce")
					if len(windows) != 1 {
						t.Fatal("WindowedReduce consumed a buffered window after cancellation")
					}
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				input := make(chan int, 1)
				input <- 1
				var calls atomic.Int64

				tt.start(ctx, input, &calls)(t)

				if len(input) != 1 {
					t.Fatalf("iteration %d: %s consumed the buffered element after cancellation", i, tt.name)
				}
				got := calls.Load()
				if got != 0 {
					t.Fatalf("iteration %d: %s called its function %d times after cancellation, want 0", i, tt.name, got)
				}
			}
		})
	}
}
