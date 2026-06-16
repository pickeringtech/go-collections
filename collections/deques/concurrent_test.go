package deques_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/deques"
)

// TestConcurrentReturnType asserts the immutable return contract: operating on a
// concurrent deque yields a new instance of the same concurrent type, never the
// plain RingBuffer.
func TestConcurrentReturnType(t *testing.T) {
	t.Run("Mutex", func(t *testing.T) {
		d := deques.NewConcurrentRingBuffer[int](1, 2, 3)

		_, isType := d.PushBack(4).(*deques.ConcurrentRingBuffer[int])
		if !isType {
			t.Errorf("PushBack returned %T, want *ConcurrentRingBuffer", d.PushBack(4))
		}
		_, isType = d.PushFront(0).(*deques.ConcurrentRingBuffer[int])
		if !isType {
			t.Errorf("PushFront returned %T, want *ConcurrentRingBuffer", d.PushFront(0))
		}
		_, _, rest := d.PopFront()
		_, isType = rest.(*deques.ConcurrentRingBuffer[int])
		if !isType {
			t.Errorf("PopFront returned %T, want *ConcurrentRingBuffer", rest)
		}
		_, _, rest = d.PopBack()
		_, isType = rest.(*deques.ConcurrentRingBuffer[int])
		if !isType {
			t.Errorf("PopBack returned %T, want *ConcurrentRingBuffer", rest)
		}
	})

	t.Run("RWMutex", func(t *testing.T) {
		d := deques.NewConcurrentRWRingBuffer[int](1, 2, 3)

		_, isType := d.PushBack(4).(*deques.ConcurrentRWRingBuffer[int])
		if !isType {
			t.Errorf("PushBack returned %T, want *ConcurrentRWRingBuffer", d.PushBack(4))
		}
		_, isType = d.PushFront(0).(*deques.ConcurrentRWRingBuffer[int])
		if !isType {
			t.Errorf("PushFront returned %T, want *ConcurrentRWRingBuffer", d.PushFront(0))
		}
		_, _, rest := d.PopFront()
		_, isType = rest.(*deques.ConcurrentRWRingBuffer[int])
		if !isType {
			t.Errorf("PopFront returned %T, want *ConcurrentRWRingBuffer", rest)
		}
		_, _, rest = d.PopBack()
		_, isType = rest.(*deques.ConcurrentRWRingBuffer[int])
		if !isType {
			t.Errorf("PopBack returned %T, want *ConcurrentRWRingBuffer", rest)
		}
	})
}

// TestConcurrentIndependence asserts a deque produced by an immutable op has its
// own lock and backing store, so mutating it does not affect the original.
func TestConcurrentIndependence(t *testing.T) {
	original := deques.NewConcurrentRingBuffer[int](1, 2, 3)
	derived := original.PushBack(4).(*deques.ConcurrentRingBuffer[int])

	derived.PushBackInPlace(5)
	derived.PopFrontInPlace()

	if original.Length() != 3 {
		t.Errorf("original mutated: Length = %d, want 3", original.Length())
	}
}

// TestConcurrentAccess hammers a concurrent deque from many goroutines and
// asserts the length stays within bounds. Run with -race to catch data races.
func TestConcurrentAccess(t *testing.T) {
	deqs := []deques.MutableDeque[int]{
		deques.NewConcurrentRingBuffer[int](),
		deques.NewConcurrentRWRingBuffer[int](),
		deques.NewBoundedConcurrentRingBuffer[int](50, deques.OverwriteOldest),
		deques.NewBoundedConcurrentRWRingBuffer[int](50, deques.RejectWhenFull),
	}
	for _, d := range deqs {
		const workers = 8
		const perWorker = 100
		var wg sync.WaitGroup
		for w := 0; w < workers; w++ {
			wg.Add(1)
			go func(base int) {
				defer wg.Done()
				for i := 0; i < perWorker; i++ {
					d.PushBackInPlace(base + i)
					d.PushFrontInPlace(base + i)
					d.PopFrontInPlace()
					d.PopBackInPlace()
					_ = d.Length()
					_ = d.AsSlice()
					for range d.Values() {
					}
				}
			}(w * perWorker)
		}
		wg.Wait()

		if d.Length() < 0 {
			t.Errorf("%T: negative length %d", d, d.Length())
		}
		capacity := d.Capacity()
		if capacity != deques.Unbounded && d.Length() > capacity {
			t.Errorf("%T: length %d exceeds capacity %d", d, d.Length(), capacity)
		}
	}
}

// ExampleConcurrentRingBuffer demonstrates safe concurrent use: many goroutines
// push to a shared deque, and the final length is deterministic.
func ExampleConcurrentRingBuffer() {
	d := deques.NewConcurrentRingBuffer[int]()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			d.PushBackInPlace(v)
		}(i)
	}
	wg.Wait()

	fmt.Println(d.Length())
	// Output: 100
}
