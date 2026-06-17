package collections

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/heaps"
	"github.com/pickeringtech/go-collections/collections/lru"
)

func TestHeapConstructors(t *testing.T) {
	// less sorts ascending, so the comparator-based constructors behave as min-heaps.
	less := heaps.Min[int]

	heapFns := map[string]func(...int) heaps.Heap[int]{
		"NewHeap":                func(v ...int) heaps.Heap[int] { return NewHeap(less, v...) },
		"NewMinHeap":             func(v ...int) heaps.Heap[int] { return NewMinHeap(v...) },
		"NewConcurrentHeap":      func(v ...int) heaps.Heap[int] { return NewConcurrentHeap(less, v...) },
		"NewConcurrentMinHeap":   func(v ...int) heaps.Heap[int] { return NewConcurrentMinHeap(v...) },
		"NewConcurrentRWHeap":    func(v ...int) heaps.Heap[int] { return NewConcurrentRWHeap(less, v...) },
		"NewConcurrentRWMinHeap": func(v ...int) heaps.Heap[int] { return NewConcurrentRWMinHeap(v...) },
	}
	for name, make := range heapFns {
		t.Run(name, func(t *testing.T) {
			h := make(5, 1, 3, 2, 4)
			if h.Length() != 5 {
				t.Errorf("Length() = %d, want 5", h.Length())
			}
			top, ok := h.Peek()
			if !ok || top != 1 {
				t.Errorf("Peek() = (%d, %t), want (1, true)", top, ok)
			}
			if got := h.AsSortedSlice(); !equalInts(got, []int{1, 2, 3, 4, 5}) {
				t.Errorf("AsSortedSlice() = %v, want [1 2 3 4 5]", got)
			}
		})
	}
}

func TestMaxHeapConstructors(t *testing.T) {
	maxHeapFns := map[string]func(...int) heaps.Heap[int]{
		"NewMaxHeap":             func(v ...int) heaps.Heap[int] { return NewMaxHeap(v...) },
		"NewConcurrentMaxHeap":   func(v ...int) heaps.Heap[int] { return NewConcurrentMaxHeap(v...) },
		"NewConcurrentRWMaxHeap": func(v ...int) heaps.Heap[int] { return NewConcurrentRWMaxHeap(v...) },
	}
	for name, make := range maxHeapFns {
		t.Run(name, func(t *testing.T) {
			h := make(5, 1, 3, 2, 4)
			top, ok := h.Peek()
			if !ok || top != 5 {
				t.Errorf("Peek() = (%d, %t), want (5, true)", top, ok)
			}
			if got := h.AsSortedSlice(); !equalInts(got, []int{5, 4, 3, 2, 1}) {
				t.Errorf("AsSortedSlice() = %v, want [5 4 3 2 1]", got)
			}
		})
	}
}

func TestLRUConstructors(t *testing.T) {
	lruFns := map[string]func(int, ...lru.Option[string, int]) lru.MutableCache[string, int]{
		"NewLRU":             NewLRU[string, int],
		"NewConcurrentLRU":   NewConcurrentLRU[string, int],
		"NewConcurrentRWLRU": NewConcurrentRWLRU[string, int],
	}
	for name, make := range lruFns {
		t.Run(name, func(t *testing.T) {
			c := make(2)
			c.PutInPlace("a", 1)
			c.PutInPlace("b", 2)
			// Get is the recency-marking read that lives on MutableCache, so the
			// facade must return MutableCache for the cache to be useful.
			if v, ok := c.Get("a"); !ok || v != 1 {
				t.Errorf("Get(a) = (%d, %t), want (1, true)", v, ok)
			}
			// "a" was just used, so adding "c" evicts the least-recently-used "b".
			c.PutInPlace("c", 3)
			if c.Length() != 2 {
				t.Errorf("Length() = %d, want 2", c.Length())
			}
			if c.Contains("b") {
				t.Errorf("Contains(b) = true, want false (should have been evicted)")
			}
			if !c.Contains("a") || !c.Contains("c") {
				t.Errorf("expected a and c to remain present")
			}
		})
	}
}

func TestNewLRUWithOptions(t *testing.T) {
	var evicted []string
	c := NewLRU[string, int](1,
		lru.WithOnEvict(func(k string, v int) { evicted = append(evicted, k) }),
		lru.WithEntries(lru.Pair[string, int]{Key: "seed", Value: 1}),
	)
	if v, ok := c.Peek("seed"); !ok || v != 1 {
		t.Errorf("Peek(seed) = (%d, %t), want (1, true)", v, ok)
	}
	c.PutInPlace("next", 2) // evicts "seed" (capacity 1)
	if len(evicted) != 1 || evicted[0] != "seed" {
		t.Errorf("evicted = %v, want [seed]", evicted)
	}
}
