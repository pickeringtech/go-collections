package heaps_test

import (
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/heaps"
)

// concurrentFactories names each thread-safe constructor under test so the
// shared behavioural suite runs against every concurrent variant.
func concurrentFactories() []struct {
	name string
	rw   bool // true for the RWMutex-backed variants
	make func(values ...int) heaps.MutableHeap[int]
} {
	return []struct {
		name string
		rw   bool
		make func(values ...int) heaps.MutableHeap[int]
	}{
		{name: "ConcurrentBinary", rw: false, make: func(values ...int) heaps.MutableHeap[int] {
			return heaps.NewConcurrentMin(values...)
		}},
		{name: "ConcurrentBinaryMax", rw: false, make: func(values ...int) heaps.MutableHeap[int] {
			return heaps.NewConcurrentMax(values...)
		}},
		{name: "ConcurrentRWBinary", rw: true, make: func(values ...int) heaps.MutableHeap[int] {
			return heaps.NewConcurrentRWMin(values...)
		}},
		{name: "ConcurrentRWBinaryMax", rw: true, make: func(values ...int) heaps.MutableHeap[int] {
			return heaps.NewConcurrentRWMax(values...)
		}},
	}
}

func TestConcurrent_BehavesLikeBinary(t *testing.T) {
	values := []int{5, 1, 4, 1, 5, 9, 2, 6}
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			h := f.make(values...)

			// AsSortedSlice and Peek agree on the priority head (compute the
			// relatively expensive sorted drain once).
			sorted := h.AsSortedSlice()
			if len(sorted) != len(values) {
				t.Fatalf("AsSortedSlice length = %d, want %d", len(sorted), len(values))
			}
			head, ok := h.Peek()
			if !ok {
				t.Fatal("Peek() ok = false, want true")
			}
			if head != sorted[0] {
				t.Errorf("Peek() = %d, want %d (drain head)", head, sorted[0])
			}

			// Immutable ops return a thread-safe heap of the same family and
			// leave the receiver untouched. Keyed off f.rw so every variant —
			// including the Max constructors — verifies the return type.
			before := h.Length()
			pushed := h.Push(7)
			if f.rw {
				if _, ok := pushed.(*heaps.ConcurrentRWBinary[int]); !ok {
					t.Errorf("RW Push returned %T, want *ConcurrentRWBinary", pushed)
				}
			} else {
				if _, ok := pushed.(*heaps.ConcurrentBinary[int]); !ok {
					t.Errorf("mutex Push returned %T, want *ConcurrentBinary", pushed)
				}
			}
			if h.Length() != before {
				t.Errorf("Push mutated receiver: Length() = %d, want %d", h.Length(), before)
			}
			if pushed.Length() != before+1 {
				t.Errorf("Push result Length() = %d, want %d", pushed.Length(), before+1)
			}

			multi := h.PushMany(0, 8)
			if multi.Length() != before+2 {
				t.Errorf("PushMany result Length() = %d, want %d", multi.Length(), before+2)
			}

			_, popOK, rest := h.Pop()
			if !popOK {
				t.Error("Pop() ok = false, want true")
			}
			if rest.Length() != before-1 {
				t.Errorf("Pop rest Length() = %d, want %d", rest.Length(), before-1)
			}
			if h.Length() != before {
				t.Errorf("Pop mutated receiver: Length() = %d, want %d", h.Length(), before)
			}

			// In-place mutation.
			h.PushInPlace(3)
			h.PushManyInPlace(2, 2)
			if h.Length() != before+3 {
				t.Errorf("after in-place pushes Length() = %d, want %d", h.Length(), before+3)
			}
			if _, ok := h.PopInPlace(); !ok {
				t.Error("PopInPlace() ok = false, want true")
			}

			// IsEmpty / ForEach / All / AsSlice / Drain.
			fresh := f.make(3, 1, 2)
			if fresh.IsEmpty() {
				t.Error("IsEmpty() = true, want false")
			}
			var each []int
			fresh.ForEach(func(e int) { each = append(each, e) })
			if len(each) != 3 {
				t.Errorf("ForEach visited %d, want 3", len(each))
			}
			var all []int
			for v := range fresh.All() {
				all = append(all, v)
			}
			if len(all) != 3 {
				t.Errorf("All() visited %d, want 3", len(all))
			}
			if got := fresh.AsSlice(); len(got) != 3 {
				t.Errorf("AsSlice() length = %d, want 3", len(got))
			}
			var drained []int
			for v := range fresh.Drain() {
				drained = append(drained, v)
			}
			if len(drained) != 3 {
				t.Errorf("Drain() visited %d, want 3", len(drained))
			}
		})
	}
}

func TestConcurrent_IteratorEarlyBreak(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			h := f.make(5, 1, 3, 2, 4)
			for range h.All() {
				break
			}
			for range h.Drain() {
				break
			}
			if h.Length() != 5 {
				t.Errorf("early break mutated receiver: Length() = %d, want 5", h.Length())
			}
		})
	}
}

func TestConcurrent_EmptyState(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			h := f.make()
			if !h.IsEmpty() {
				t.Error("IsEmpty() = false, want true")
			}
			if h.Length() != 0 {
				t.Errorf("Length() = %d, want 0", h.Length())
			}
			if v, ok := h.Peek(); ok || v != 0 {
				t.Errorf("Peek() = (%d, %v), want (0, false)", v, ok)
			}
			if v, ok := h.PopInPlace(); ok || v != 0 {
				t.Errorf("PopInPlace() = (%d, %v), want (0, false)", v, ok)
			}
		})
	}
}

// TestConcurrent_RaceSafety hammers a shared heap from many goroutines; run with
// -race to surface data races.
func TestConcurrent_RaceSafety(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			h := f.make()
			const workers = 50

			var wg sync.WaitGroup
			wg.Add(workers)
			for w := 0; w < workers; w++ {
				go func(base int) {
					defer wg.Done()
					h.PushInPlace(base)
					h.PushManyInPlace(base+1, base+2)
					_, _ = h.Peek()
					_ = h.Length()
					_ = h.AsSortedSlice()
					for range h.All() {
						break
					}
					_, _ = h.PopInPlace()
				}(w * 3)
			}
			wg.Wait()

			// 50 workers each net +2 elements (3 pushed, 1 popped).
			if got := h.Length(); got != workers*2 {
				t.Errorf("after race Length() = %d, want %d", got, workers*2)
			}
			// The remaining elements still drain in priority order — ascending
			// for the min variants, descending for the max variants.
			drained := h.AsSortedSlice()
			if !isMonotonic(drained) {
				t.Errorf("AsSortedSlice() = %v is not in priority order", drained)
			}
		})
	}
}

func TestConcurrent_NewConcurrentComparator(t *testing.T) {
	// The general comparator-driven constructors accept an arbitrary LessFunc.
	mx := heaps.NewConcurrent(heaps.Max[int], 1, 9, 4)
	if v, _ := mx.Peek(); v != 9 {
		t.Errorf("NewConcurrent(Max) Peek() = %d, want 9", v)
	}
	rw := heaps.NewConcurrentRW(heaps.Min[int], 7, 2, 5)
	if v, _ := rw.Peek(); v != 2 {
		t.Errorf("NewConcurrentRW(Min) Peek() = %d, want 2", v)
	}
}

// isMonotonic reports whether the slice is sorted ascending or descending,
// covering both the min and max priority drains.
func isMonotonic(values []int) bool {
	return sort.IntsAreSorted(values) || sort.SliceIsSorted(values, func(i, j int) bool {
		return values[i] > values[j]
	})
}

// assertSorted is a tiny guard used by the concurrent suite.
func assertSorted(t *testing.T, got, want []int) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestConcurrent_DrainMatchesOracle(t *testing.T) {
	values := []int{4, 2, 7, 1}
	want := []int{1, 2, 4, 7}
	assertSorted(t, heaps.NewConcurrentMin(values...).AsSortedSlice(), want)
	assertSorted(t, heaps.NewConcurrentRWMin(values...).AsSortedSlice(), want)
}
