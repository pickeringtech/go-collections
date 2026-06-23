package channels_test

import (
	"context"
	"testing"

	"github.com/pickeringtech/go-collections/channels"
	"github.com/pickeringtech/go-collections/slices"
)

// FuzzFromSliceCollect asserts the round-trip invariant: streaming a slice
// through a channel and collecting it back reproduces the original slice
// exactly, preserving both element count and order.
func FuzzFromSliceCollect(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{1})
	f.Add([]byte{1, 2, 3, 4, 5})

	f.Fuzz(func(t *testing.T, data []byte) {
		input := make([]byte, len(data))
		copy(input, data)

		got := channels.CollectAsSlice(channels.FromSlice(context.Background(), input))

		if len(got) != len(input) {
			t.Fatalf("round-trip length = %d, want %d", len(got), len(input))
		}
		for i := range input {
			if got[i] != input[i] {
				t.Fatalf("round-trip[%d] = %d, want %d", i, got[i], input[i])
			}
		}
	})
}

// FuzzCollectNAsSlice asserts that CollectNAsSlice returns exactly the first
// min(n, len) elements in order.
func FuzzCollectNAsSlice(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4}, 2)
	f.Add([]byte{1, 2, 3}, 10)
	f.Add([]byte{}, 5)
	f.Add([]byte{9}, 0)

	f.Fuzz(func(t *testing.T, data []byte, n int) {
		// Negative counts are outside the meaningful domain.
		if n < 0 {
			return
		}
		input := make([]byte, len(data))
		copy(input, data)

		got := channels.CollectNAsSlice(channels.FromSlice(context.Background(), input), n)

		want := n
		if want > len(input) {
			want = len(input)
		}
		if len(got) != want {
			t.Fatalf("CollectN length = %d, want %d", len(got), want)
		}
		for i := 0; i < want; i++ {
			if got[i] != input[i] {
				t.Fatalf("CollectN[%d] = %d, want %d", i, got[i], input[i])
			}
		}
	})
}

// FuzzPipelineFilterMap asserts that a Filter→Map pipeline preserves the order
// of the surviving elements and produces exactly one output per element that
// passes the filter — matching a hand-rolled reference computation.
func FuzzPipelineFilterMap(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{1, 2, 3, 4, 5, 6})
	f.Add([]byte{2, 4, 6})

	f.Fuzz(func(t *testing.T, data []byte) {
		input := make([]byte, len(data))
		copy(input, data)

		keep := func(b byte) bool { return b%2 == 0 }
		transform := func(b byte) int { return int(b) * 10 }

		ctx := context.Background()
		filtered := channels.Filter(ctx, channels.FromSlice(ctx, input), keep)
		mapped := channels.Map(ctx, filtered, transform)
		got := channels.CollectAsSlice(mapped)

		var want []int
		for _, b := range input {
			if keep(b) {
				want = append(want, transform(b))
			}
		}
		if len(got) != len(want) {
			t.Fatalf("pipeline length = %d, want %d", len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("pipeline[%d] = %d, want %d", i, got[i], want[i])
			}
		}
	})
}

// FuzzTumblingWindow is a differential fuzz test: streaming an arbitrary byte
// slice through TumblingWindow must reproduce slices.Chunk with its trailing
// partial chunk dropped — same window count, same widths (all == size), same
// element order — and never panic.
func FuzzTumblingWindow(f *testing.F) {
	f.Add([]byte(nil), 2)
	f.Add([]byte{}, 2)
	f.Add([]byte{1}, 1)
	f.Add([]byte{1, 2, 3, 4, 5, 6}, 3)
	f.Add([]byte{1, 2, 3, 4, 5, 6, 7}, 3)
	f.Add([]byte{1, 2, 3}, 0)

	f.Fuzz(func(t *testing.T, data []byte, size int) {
		// Keep size in a sane band so fuzzing explores boundaries without
		// allocating absurd windows; non-positive stays meaningful.
		if size > len(data)+8 {
			size = len(data) + 8
		}
		if size < -2 {
			size = -2
		}

		ctx := context.Background()
		out := channels.TumblingWindow(ctx, channels.FromSlice(ctx, data), size)

		var want [][]byte
		for _, chunk := range slices.Chunk(data, size) {
			if len(chunk) == size {
				want = append(want, chunk)
			}
		}

		idx := 0
		for w := range out {
			if idx >= len(want) {
				t.Fatalf("TumblingWindow emitted more windows than the chunk-oracle (size=%d, data=%v)", size, data)
			}
			if len(w) != size {
				t.Fatalf("window %d width = %d, want %d", idx, len(w), size)
			}
			for j := range w {
				if w[j] != want[idx][j] {
					t.Fatalf("window %d element %d = %d, want %d", idx, j, w[j], want[idx][j])
				}
			}
			idx++
		}
		if idx != len(want) {
			t.Fatalf("TumblingWindow emitted %d windows, chunk-oracle has %d (size=%d)", idx, len(want), size)
		}
	})
}

// FuzzSlidingWindow is a differential fuzz test for the step==1 case: streaming
// an arbitrary byte slice through SlidingWindow with step 1 must reproduce
// slices.Window exactly — same window count, same widths (all == size), same
// element order — and never panic.
func FuzzSlidingWindow(f *testing.F) {
	f.Add([]byte(nil), 2)
	f.Add([]byte{}, 2)
	f.Add([]byte{1}, 1)
	f.Add([]byte{1, 2, 3, 4, 5}, 3)
	f.Add([]byte{1, 2, 3, 4, 5, 6}, 2)
	f.Add([]byte{1, 2, 3}, 0)

	f.Fuzz(func(t *testing.T, data []byte, size int) {
		if size > len(data)+8 {
			size = len(data) + 8
		}
		if size < -2 {
			size = -2
		}

		ctx := context.Background()
		out := channels.SlidingWindow(ctx, channels.FromSlice(ctx, data), size, 1)

		want := slices.Window(data, size) // step-1 sliding window is exactly this

		idx := 0
		for w := range out {
			if idx >= len(want) {
				t.Fatalf("SlidingWindow emitted more windows than the window-oracle (size=%d, data=%v)", size, data)
			}
			if len(w) != size {
				t.Fatalf("window %d width = %d, want %d", idx, len(w), size)
			}
			for j := range w {
				if w[j] != want[idx][j] {
					t.Fatalf("window %d element %d = %d, want %d", idx, j, w[j], want[idx][j])
				}
			}
			idx++
		}
		if idx != len(want) {
			t.Fatalf("SlidingWindow emitted %d windows, window-oracle has %d (size=%d)", idx, len(want), size)
		}
	})
}
