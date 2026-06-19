package bloom

import (
	"sync"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
)

// ConcurrentFilter is a goroutine-safe Bloom filter. It wraps a plain Filter
// with a sync.RWMutex — all the filter logic lives on Filter, and every method
// here simply takes the appropriate lock and delegates, so there is a single
// implementation to reason about. Membership queries take a read lock (they
// dominate the typical workload); mutations take the write lock.
//
// ConcurrentFilter must not be copied after first use; copying duplicates the
// lock over shared backing data and breaks the safety contract. go vet's
// copylocks pass reports any such copy because of the embedded nocopy guard.
//
// The zero value is not usable; construct one with NewConcurrent.
type ConcurrentFilter[T comparable] struct {
	_      nocopy.NoCopy
	filter *Filter[T]
	lock   sync.RWMutex
}

// Interface guard.
var _ Membership[string] = (*ConcurrentFilter[string])(nil)

// NewConcurrent creates a goroutine-safe Bloom filter sized for expectedItems
// insertions at the target falsePositiveRate. It validates its arguments like
// New and returns the same errors.
func NewConcurrent[T comparable](expectedItems int, falsePositiveRate float64, opts ...Option[T]) (*ConcurrentFilter[T], error) {
	f, err := New(expectedItems, falsePositiveRate, opts...)
	if err != nil {
		return nil, err
	}
	return &ConcurrentFilter[T]{filter: f}, nil
}

// Add inserts an element. Safe for concurrent use.
func (c *ConcurrentFilter[T]) Add(value T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.filter.Add(value)
}

// Contains reports whether value may be present. Safe for concurrent use.
func (c *ConcurrentFilter[T]) Contains(value T) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.filter.Contains(value)
}

// Merge folds other into the receiver under the write lock. other is read
// without locking, so it must not be mutated concurrently by another goroutine.
func (c *ConcurrentFilter[T]) Merge(other *Filter[T]) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.filter.Merge(other)
}

// Clear empties the filter. Safe for concurrent use.
func (c *ConcurrentFilter[T]) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.filter.Clear()
}

// ApproxCount estimates the number of distinct elements added. Safe for
// concurrent use.
func (c *ConcurrentFilter[T]) ApproxCount() uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.filter.ApproxCount()
}

// EstimatedFalsePositiveRate returns the current false-positive probability.
// Safe for concurrent use.
func (c *ConcurrentFilter[T]) EstimatedFalsePositiveRate() float64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.filter.EstimatedFalsePositiveRate()
}

// BitSize returns the underlying bit-array size (m).
func (c *ConcurrentFilter[T]) BitSize() uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.filter.BitSize()
}

// HashCount returns the number of hash functions per element (k).
func (c *ConcurrentFilter[T]) HashCount() uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.filter.HashCount()
}

// Snapshot returns a plain, unlocked copy of the underlying filter. It is the
// safe way to obtain a *Filter to pass to another filter's Merge: the copy is
// taken under the read lock and is independent of the receiver.
func (c *ConcurrentFilter[T]) Snapshot() *Filter[T] {
	c.lock.RLock()
	defer c.lock.RUnlock()
	cp := &Filter[T]{
		bitsArr: make([]uint64, len(c.filter.bitsArr)),
		m:       c.filter.m,
		k:       c.filter.k,
		seed:    c.filter.seed,
		hasher:  c.filter.hasher,
	}
	copy(cp.bitsArr, c.filter.bitsArr)
	return cp
}
