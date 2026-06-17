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

// heapFuzzFactories returns the heap constructors to fuzz, each producing a
// fresh MutableHeap[uint8] from the given comparator and values. The concurrent
// variants share the sequential heap's semantics, so they must satisfy the same
// oracle.
func heapFuzzFactories() []struct {
	name string
	make func(heaps.LessFunc[uint8], ...uint8) heaps.MutableHeap[uint8]
} {
	return []struct {
		name string
		make func(heaps.LessFunc[uint8], ...uint8) heaps.MutableHeap[uint8]
	}{
		{"Binary", func(l heaps.LessFunc[uint8], v ...uint8) heaps.MutableHeap[uint8] { return heaps.New(l, v...) }},
		{"ConcurrentBinary", func(l heaps.LessFunc[uint8], v ...uint8) heaps.MutableHeap[uint8] {
			return heaps.NewConcurrent(l, v...)
		}},
		{"ConcurrentRWBinary", func(l heaps.LessFunc[uint8], v ...uint8) heaps.MutableHeap[uint8] {
			return heaps.NewConcurrentRW(l, v...)
		}},
	}
}

// FuzzBinaryOracle is a differential fuzz test: for every heap backend and for
// arbitrary input bytes it checks the heap against a native sort oracle and
// asserts the structural heap invariant, exercising heapify, sift-up and
// sift-down on real data.
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
		maxLess := func(a, b uint8) bool { return a > b }

		// The max-heap drains to the reverse of the ascending oracle.
		maxWant := make([]uint8, len(want))
		for i, v := range want {
			maxWant[len(want)-1-i] = v
		}

		for _, backend := range heapFuzzFactories() {
			built := backend.make(minLess, values...)

			// Length and Peek agree with the oracle.
			if built.Length() != len(want) {
				t.Fatalf("%s: Length() = %d, want %d", backend.name, built.Length(), len(want))
			}
			if len(want) > 0 {
				top, ok := built.Peek()
				if !ok || top != want[0] {
					t.Fatalf("%s: Peek() = (%d, %v), want (%d, true)", backend.name, top, ok, want[0])
				}
			}

			// The heap-array snapshot honours the heap invariant.
			snapshot := built.AsSlice()
			if !isMinHeap(snapshot) {
				t.Fatalf("%s: AsSlice() = %v violates the min-heap invariant", backend.name, snapshot)
			}

			// Heapify-built and incrementally-pushed heaps both drain to the oracle.
			gotHeapify := drainUint8(t, backend.make(minLess, values...))
			assertUint8Equal(t, backend.name+" heapify drain", gotHeapify, want)

			incremental := backend.make(minLess)
			for _, v := range values {
				incremental.PushInPlace(v)
			}
			assertUint8Equal(t, backend.name+" incremental drain", drainUint8(t, incremental), want)

			// AsSortedSlice is non-mutating and matches the oracle too.
			assertUint8Equal(t, backend.name+" AsSortedSlice", built.AsSortedSlice(), want)
			if built.Length() != len(want) {
				t.Fatalf("%s: AsSortedSlice mutated the heap: Length() = %d, want %d", backend.name, built.Length(), len(want))
			}

			// FromSeq(All) round-trips back to a heap that drains to the oracle.
			assertUint8Equal(t, backend.name+" FromSeq", heaps.FromSeq(heaps.Min[uint8], built.All()).AsSortedSlice(), want)

			maxHeap := backend.make(maxLess, values...)
			assertUint8Equal(t, backend.name+" max drain", maxHeap.AsSortedSlice(), maxWant)
		}
	})
}

func drainUint8(t *testing.T, h heaps.MutableHeap[uint8]) []uint8 {
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
