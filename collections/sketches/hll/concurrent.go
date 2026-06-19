package hll

import (
	"sync"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
)

// ConcurrentSketch is a goroutine-safe HyperLogLog. It wraps a plain Sketch with
// a sync.RWMutex — all the logic lives on Sketch, and every method here takes
// the appropriate lock and delegates. Count and the getters take a read lock;
// Add, Merge and Clear take the write lock.
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
var _ Cardinality[string] = (*ConcurrentSketch[string])(nil)

// NewConcurrent creates a goroutine-safe HyperLogLog with the same options and
// validation as New.
func NewConcurrent[T comparable](opts ...Option[T]) (*ConcurrentSketch[T], error) {
	s, err := New(opts...)
	if err != nil {
		return nil, err
	}
	return &ConcurrentSketch[T]{sketch: s}, nil
}

// Add records an element. Safe for concurrent use.
func (c *ConcurrentSketch[T]) Add(value T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.sketch.Add(value)
}

// Count returns the estimated number of distinct elements added. Safe for
// concurrent use.
func (c *ConcurrentSketch[T]) Count() uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.sketch.Count()
}

// Merge folds other into the receiver under the write lock. other is read
// without locking, so it must not be mutated concurrently elsewhere.
func (c *ConcurrentSketch[T]) Merge(other *Sketch[T]) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.sketch.Merge(other)
}

// Clear resets every register. Safe for concurrent use.
func (c *ConcurrentSketch[T]) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.sketch.Clear()
}

// Precision returns the register-index width (p).
func (c *ConcurrentSketch[T]) Precision() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.sketch.Precision()
}

// RegisterCount returns the number of registers (m = 2^p).
func (c *ConcurrentSketch[T]) RegisterCount() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.sketch.RegisterCount()
}

// StandardError returns the expected relative error of Count.
func (c *ConcurrentSketch[T]) StandardError() float64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.sketch.StandardError()
}

// Snapshot returns a plain, unlocked copy of the underlying sketch, taken under
// the read lock — the safe way to obtain a *Sketch to pass to another sketch's
// Merge.
func (c *ConcurrentSketch[T]) Snapshot() *Sketch[T] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	cp := &Sketch[T]{
		registers: make([]uint8, len(c.sketch.registers)),
		precision: c.sketch.precision,
		seed:      c.sketch.seed,
		hasher:    c.sketch.hasher,
	}
	copy(cp.registers, c.sketch.registers)
	return cp
}
