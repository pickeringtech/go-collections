package channels_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/channels"
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

		got := channels.CollectAsSlice(channels.FromSlice(input))

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

		got := channels.CollectNAsSlice(channels.FromSlice(input), n)

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

		filtered := channels.Filter(channels.FromSlice(input), keep)
		mapped := channels.Map(filtered, transform)
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
