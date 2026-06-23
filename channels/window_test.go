package channels_test

import (
	"context"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/pickeringtech/go-collections/channels"
	"github.com/pickeringtech/go-collections/slices"
)

// collectWindows drains a window channel into a slice of windows.
func collectWindows[T any](ch <-chan []T) [][]T {
	var out [][]T
	for w := range ch {
		out = append(out, w)
	}
	return out
}

func TestTumblingWindow(t *testing.T) {
	type args struct {
		input []int
		size  int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{name: "nil input yields no windows", args: args{input: nil, size: 3}, want: nil},
		{name: "empty input yields no windows", args: args{input: []int{}, size: 3}, want: nil},
		{name: "size zero yields no windows", args: args{input: []int{1, 2, 3}, size: 0}, want: nil},
		{name: "negative size yields no windows", args: args{input: []int{1, 2, 3}, size: -1}, want: nil},
		{name: "exact multiple", args: args{input: []int{1, 2, 3, 4}, size: 2}, want: [][]int{{1, 2}, {3, 4}}},
		{name: "trailing partial dropped", args: args{input: []int{1, 2, 3, 4, 5, 6, 7}, size: 3}, want: [][]int{{1, 2, 3}, {4, 5, 6}}},
		{name: "size larger than input drops everything", args: args{input: []int{1, 2}, size: 5}, want: nil},
		{name: "size one", args: args{input: []int{1, 2, 3}, size: 1}, want: [][]int{{1}, {2}, {3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			out := channels.TumblingWindow(ctx, channels.FromSlice(ctx, tt.args.input), tt.args.size)
			got := collectWindows(out)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TumblingWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTumblingWindow_MatchesChunkOracle cross-checks against slices.Chunk with
// its trailing partial chunk dropped — the full-windows-only equivalent.
func TestTumblingWindow_MatchesChunkOracle(t *testing.T) {
	inputs := [][]int{
		{1, 2, 3, 4, 5, 6},
		{1, 2, 3, 4, 5, 6, 7},
		{1},
		{},
	}
	sizes := []int{1, 2, 3, 4}
	for _, in := range inputs {
		for _, size := range sizes {
			ctx := context.Background()
			out := channels.TumblingWindow(ctx, channels.FromSlice(ctx, in), size)
			got := collectWindows(out)

			var want [][]int
			for _, chunk := range slices.Chunk(in, size) {
				if len(chunk) == size {
					want = append(want, chunk)
				}
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("TumblingWindow(%v, %d) = %v, chunk-oracle = %v", in, size, got, want)
			}
		}
	}
}

func TestSlidingWindow(t *testing.T) {
	type args struct {
		input []int
		size  int
		step  int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{name: "nil input", args: args{input: nil, size: 2, step: 1}, want: nil},
		{name: "size zero", args: args{input: []int{1, 2, 3}, size: 0, step: 1}, want: nil},
		{name: "step zero", args: args{input: []int{1, 2, 3}, size: 2, step: 0}, want: nil},
		{name: "negative step", args: args{input: []int{1, 2, 3}, size: 2, step: -1}, want: nil},
		{name: "overlap step 1", args: args{input: []int{1, 2, 3, 4, 5}, size: 3, step: 1}, want: [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}}},
		{name: "step equals size is tumbling", args: args{input: []int{1, 2, 3, 4, 5, 6}, size: 2, step: 2}, want: [][]int{{1, 2}, {3, 4}, {5, 6}}},
		{name: "step greater than size skips", args: args{input: []int{1, 2, 3, 4, 5, 6, 7, 8}, size: 2, step: 3}, want: [][]int{{1, 2}, {4, 5}, {7, 8}}},
		{name: "partial trailing dropped", args: args{input: []int{1, 2, 3, 4}, size: 3, step: 2}, want: [][]int{{1, 2, 3}}},
		{name: "size larger than input", args: args{input: []int{1, 2}, size: 5, step: 1}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			out := channels.SlidingWindow(ctx, channels.FromSlice(ctx, tt.args.input), tt.args.size, tt.args.step)
			got := collectWindows(out)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlidingWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSlidingWindow_MatchesWindowOracle cross-checks the step==1 case against
// slices.Window, which is exactly the advance-by-one sliding window.
func TestSlidingWindow_MatchesWindowOracle(t *testing.T) {
	inputs := [][]int{
		{1, 2, 3, 4, 5, 6},
		{1, 2, 3},
		{1},
		{},
	}
	sizes := []int{1, 2, 3}
	for _, in := range inputs {
		for _, size := range sizes {
			ctx := context.Background()
			out := channels.SlidingWindow(ctx, channels.FromSlice(ctx, in), size, 1)
			got := collectWindows(out)

			want := [][]int{}
			want = append(want, slices.Window(in, size)...)
			if len(got) == 0 && len(want) == 0 {
				continue
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("SlidingWindow(%v, %d, 1) = %v, window-oracle = %v", in, size, got, want)
			}
		}
	}
}

func TestSessionWindow(t *testing.T) {
	type args struct {
		input   []int
		maxStep int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{name: "nil input", args: args{input: nil, maxStep: 5}, want: nil},
		{name: "single element flushes on close", args: args{input: []int{1}, maxStep: 5}, want: [][]int{{1}}},
		{name: "one continuous session", args: args{input: []int{1, 2, 3}, maxStep: 5}, want: [][]int{{1, 2, 3}}},
		{name: "three sessions flushed", args: args{input: []int{1, 2, 10, 11, 30}, maxStep: 5}, want: [][]int{{1, 2}, {10, 11}, {30}}},
		{name: "every element a new session", args: args{input: []int{1, 100, 200}, maxStep: 5}, want: [][]int{{1}, {100}, {200}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			gap := func(prev, next int) bool { return next-prev <= tt.args.maxStep }
			out := channels.SessionWindow(ctx, channels.FromSlice(ctx, tt.args.input), gap)
			got := collectWindows(out)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindowedReduce(t *testing.T) {
	type args struct {
		input []int
		size  int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{name: "sum of tumbling windows", args: args{input: []int{1, 2, 3, 4, 5, 6}, size: 3}, want: []int{6, 15}},
		{name: "partial dropped before reduce", args: args{input: []int{1, 2, 3, 4, 5, 6, 7}, size: 3}, want: []int{6, 15}},
		{name: "no full window yields nothing", args: args{input: []int{1, 2}, size: 5}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			windower := func(ctx context.Context, in <-chan int) <-chan []int {
				return channels.TumblingWindow(ctx, in, tt.args.size)
			}
			sum := func(w []int) int {
				total := 0
				for _, v := range w {
					total += v
				}
				return total
			}
			out := channels.WindowedReduce(ctx, channels.FromSlice(ctx, tt.args.input), windower, sum)
			got := channels.CollectAsSlice(out)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WindowedReduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestWindowEmitsDefensiveCopies asserts that mutating an emitted window does not
// corrupt a later window — windows must be independent copies, never views into
// a reused buffer.
func TestWindowEmitsDefensiveCopies(t *testing.T) {
	ctx := context.Background()
	out := channels.TumblingWindow(ctx, channels.FromSlice(ctx, []int{1, 2, 3, 4}), 2)
	windows := collectWindows(out)
	if len(windows) != 2 {
		t.Fatalf("got %d windows, want 2", len(windows))
	}
	windows[0][0] = 999
	if windows[1][0] != 3 {
		t.Errorf("mutating window 0 changed window 1: %v", windows[1])
	}
}

// TestTumblingWindowCancellation asserts cancelling the context tears the
// goroutine down: the output channel closes even though the input never closes.
func TestTumblingWindowCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int) // never written, never closed
	out := channels.TumblingWindow(ctx, input, 2)

	cancel()

	assertChannelCloses(t, out, "TumblingWindow")
}

// TestSlidingWindowCancellation mirrors the tumbling case for the sliding verb.
func TestSlidingWindowCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int)
	out := channels.SlidingWindow(ctx, input, 2, 1)

	cancel()

	assertChannelCloses(t, out, "SlidingWindow")
}

// TestSessionWindowCancellation asserts the session verb tears down on cancel
// without flushing the open session.
func TestSessionWindowCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int, 1)
	input <- 1 // an element buffered into the open session

	out := channels.SessionWindow(ctx, input, func(prev, next int) bool { return true })

	cancel()

	assertChannelCloses(t, out, "SessionWindow")
}

// TestWindowedReduceCancellation asserts the composed reduce stage tears down on
// cancel.
func TestWindowedReduceCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int)
	windower := func(ctx context.Context, in <-chan int) <-chan []int {
		return channels.TumblingWindow(ctx, in, 2)
	}
	out := channels.WindowedReduce(ctx, input, windower, func(w []int) int { return len(w) })

	cancel()

	assertChannelCloses(t, out, "WindowedReduce")
}

// TestTumblingWindowCancellationWhileSending covers the other teardown path: a
// full window is ready and the goroutine is blocked trying to deliver it
// downstream (the output is unbuffered and unread). Cancelling unblocks that
// send, so the goroutine abandons the window and exits rather than leaking.
func TestTumblingWindowCancellationWhileSending(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int, 2)
	input <- 1
	input <- 2 // completes a window of size 2, so the goroutine reaches the send

	out := channels.TumblingWindow(ctx, input, 2)

	for i := 0; i < 1000; i++ {
		runtime.Gosched()
	}
	cancel()

	assertChannelCloses(t, out, "TumblingWindow")
}

// TestSlidingWindowCancellationWhileSending covers the sliding verb's blocked
// send: a full window is ready and the goroutine parks delivering it downstream
// until cancel unblocks the send and tears it down.
func TestSlidingWindowCancellationWhileSending(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int, 2)
	input <- 1
	input <- 2 // completes a window of size 2, so the goroutine reaches the send

	out := channels.SlidingWindow(ctx, input, 2, 1)

	for i := 0; i < 1000; i++ {
		runtime.Gosched()
	}
	cancel()

	assertChannelCloses(t, out, "SlidingWindow")
}

// TestSessionWindowCancellationWhileSending covers the session verb's blocked
// send: two same-session elements then a third that opens a new session, which
// emits the buffered session — the goroutine parks on that send until cancel.
func TestSessionWindowCancellationWhileSending(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int, 3)
	input <- 1
	input <- 2
	input <- 100 // gap from 2, so the [1 2] session is emitted on the blocked send

	gap := func(prev, next int) bool { return next-prev <= 5 }
	out := channels.SessionWindow(ctx, input, gap)

	for i := 0; i < 1000; i++ {
		runtime.Gosched()
	}
	cancel()

	assertChannelCloses(t, out, "SessionWindow")
}

// TestWindowedReduceCancellationWhileSending covers WindowedReduce's blocked
// send: a window is reduced and the goroutine parks delivering the result.
func TestWindowedReduceCancellationWhileSending(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int, 2)
	input <- 1
	input <- 2

	windower := func(ctx context.Context, in <-chan int) <-chan []int {
		return channels.TumblingWindow(ctx, in, 2)
	}
	out := channels.WindowedReduce(ctx, input, windower, func(w []int) int { return len(w) })

	for i := 0; i < 1000; i++ {
		runtime.Gosched()
	}
	cancel()

	assertChannelCloses(t, out, "WindowedReduce")
}

// assertChannelCloses fails unless ch closes (drains to a final !ok receive)
// within a generous deadline.
func assertChannelCloses[T any](t *testing.T, ch <-chan T, name string) {
	t.Helper()
	deadline := time.After(time.Second)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-deadline:
			t.Fatalf("%s() output channel did not close after cancellation", name)
		}
	}
}
