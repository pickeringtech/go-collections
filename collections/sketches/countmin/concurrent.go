package countmin

import (
	"sync"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
)

// ConcurrentSketch is a goroutine-safe Count-Min sketch. It wraps a plain Sketch
// with a sync.RWMutex — all the logic lives on Sketch, and every method here
// takes the appropriate lock and delegates. Estimate and the dimension getters
// take a read lock; the mutating operations take the write lock.
//
// ConcurrentSketch must not be copied after first use; go vet's copylocks pass
// reports any such copy because of the embedded nocopy guard.
//
// The zero value is not usable; construct one with NewConcurrent.
type ConcurrentSketch[T comparable] struct {
	_      nocopy.NoCopy
	sketch *Sketch[T]
	lock   sync.RWMutex
}

// Interface guard.
var _ Frequency[string] = (*ConcurrentSketch[string])(nil)

// NewConcurrent creates a goroutine-safe Count-Min sketch with the same error
// bounds as New, returning the same errors on invalid input.
func NewConcurrent[T comparable](epsilon, delta float64, opts ...Option[T]) (*ConcurrentSketch[T], error) {
	s, err := New(epsilon, delta, opts...)
	if err != nil {
		return nil, err
	}
	return &ConcurrentSketch[T]{sketch: s}, nil
}

// Add records one occurrence of value. Safe for concurrent use.
func (c *ConcurrentSketch[T]) Add(value T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.sketch.Add(value)
}

// AddCount records n occurrences of value. Safe for concurrent use.
func (c *ConcurrentSketch[T]) AddCount(value T, n uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.sketch.AddCount(value, n)
}

// Estimate returns the approximate count of value. Safe for concurrent use.
func (c *ConcurrentSketch[T]) Estimate(value T) uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.sketch.Estimate(value)
}

// Merge folds other into the receiver under the write lock. other is read
// without locking, so it must not be mutated concurrently elsewhere.
func (c *ConcurrentSketch[T]) Merge(other *Sketch[T]) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.sketch.Merge(other)
}

// Clear resets every counter. Safe for concurrent use.
func (c *ConcurrentSketch[T]) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.sketch.Clear()
}

// Total returns the sum of all counts added. Safe for concurrent use.
func (c *ConcurrentSketch[T]) Total() uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.sketch.Total()
}

// Width returns the number of columns per row (w).
func (c *ConcurrentSketch[T]) Width() uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.sketch.Width()
}

// Depth returns the number of rows (d).
func (c *ConcurrentSketch[T]) Depth() uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.sketch.Depth()
}

// Snapshot returns a plain, unlocked copy of the underlying sketch, taken under
// the read lock — the safe way to obtain a *Sketch to pass to another sketch's
// Merge.
func (c *ConcurrentSketch[T]) Snapshot() *Sketch[T] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	cp := &Sketch[T]{
		counts: make([]uint64, len(c.sketch.counts)),
		w:      c.sketch.w,
		d:      c.sketch.d,
		seed:   c.sketch.seed,
		hasher: c.sketch.hasher,
		total:  c.sketch.total,
	}
	copy(cp.counts, c.sketch.counts)
	return cp
}
