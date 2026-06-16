package heaps_test

import (
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/heaps"
)

// isMinHeap verifies the binary-heap invariant over a heap-array snapshot: every
// parent is <= its children.
func isMinHeap(data []byte) bool {
	for i := range data {
		left := 2*i + 1
		right := 2*i + 2
		if left < len(data) && data[i] > data[left] {
			return false
		}
		if right < len(data) && data[i] > data[right] {
			return false
		}
	}
	return true
}

// FuzzBinaryOracle is a differential fuzz test: for arbitrary input bytes it
// checks the heap against a native sort oracle and asserts the structural heap
// invariant, exercising heapify, sift-up and sift-down on real data.
func FuzzBinaryOracle(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{1})
	f.Add([]byte{3, 1, 2})
	f.Add([]byte{5, 5, 5, 5})
	f.Add([]byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0})

	f.Fuzz(func(t *testing.T, data []byte) {
		values := make([]uint8, len(data))
		copy(values, data)

		// Native oracle: ascending sort for the min-heap.
		want := make([]uint8, len(values))
		copy(want, values)
		sort.Slice(want, func(i, j int) bool { return want[i] < want[j] })

		minLess := func(a, b uint8) bool { return a < b }
		built := heaps.New(minLess, values...)

		// Length and Peek agree with the oracle.
		if built.Length() != len(want) {
			t.Fatalf("Length() = %d, want %d", built.Length(), len(want))
		}
		if len(want) > 0 {
			top, ok := built.Peek()
			if !ok || top != want[0] {
				t.Fatalf("Peek() = (%d, %v), want (%d, true)", top, ok, want[0])
			}
		}

		// The heap-array snapshot honours the heap invariant.
		if !isMinHeap(built.AsSlice()) {
			t.Fatalf("AsSlice() = %v violates the min-heap invariant", built.AsSlice())
		}

		// Heapify-built and incrementally-pushed heaps both drain to the oracle.
		gotHeapify := drainUint8(t, heaps.New(minLess, values...))
		assertUint8Equal(t, "heapify drain", gotHeapify, want)

		incremental := heaps.New(minLess)
		for _, v := range values {
			incremental.PushInPlace(v)
		}
		assertUint8Equal(t, "incremental drain", drainUint8(t, incremental), want)

		// AsSortedSlice is non-mutating and matches the oracle too.
		assertUint8Equal(t, "AsSortedSlice", built.AsSortedSlice(), want)
		if built.Length() != len(want) {
			t.Fatalf("AsSortedSlice mutated the heap: Length() = %d, want %d", built.Length(), len(want))
		}

		// The max-heap drains to the reverse of the ascending oracle.
		maxWant := make([]uint8, len(want))
		for i, v := range want {
			maxWant[len(want)-1-i] = v
		}
		maxHeap := heaps.New(func(a, b uint8) bool { return a > b }, values...)
		assertUint8Equal(t, "max drain", maxHeap.AsSortedSlice(), maxWant)
	})
}

func drainUint8(t *testing.T, h *heaps.Binary[uint8]) []uint8 {
	t.Helper()
	out := make([]uint8, 0, h.Length())
	for {
		v, ok := h.PopInPlace()
		if !ok {
			return out
		}
		out = append(out, v)
	}
}

func assertUint8Equal(t *testing.T, label string, got, want []uint8) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: length = %d, want %d", label, len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s: at %d got %d, want %d (full: %v vs %v)", label, i, got[i], want[i], got, want)
		}
	}
}
