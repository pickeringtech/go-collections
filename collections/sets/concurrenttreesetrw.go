package sets

import (
	"iter"
	"sync"

	"github.com/pickeringtech/go-collections/constraints"
)

// ConcurrentTreeSetRW is a thread-safe sorted set backed by a TreeSet with a
// read-write mutex for synchronization. Read operations use read locks so many
// readers can proceed concurrently, while writes are exclusive. Elements are
// kept in sorted order. Prefer it over ConcurrentTreeSet for read-heavy workloads.
type ConcurrentTreeSetRW[T constraints.Ordered] struct {
	set  *TreeSet[T]
	lock *sync.RWMutex
}

// NewConcurrentTreeSetRW creates a new ConcurrentTreeSetRW with the given elements.
// Duplicate elements are automatically removed.
func NewConcurrentTreeSetRW[T constraints.Ordered](elements ...T) *ConcurrentTreeSetRW[T] {
	return &ConcurrentTreeSetRW[T]{
		set:  NewTreeSet[T](elements...),
		lock: &sync.RWMutex{},
	}
}

// Interface guards to ensure ConcurrentTreeSetRW implements the required interfaces.
var _ Set[string] = &ConcurrentTreeSetRW[string]{}
var _ MutableSet[string] = &ConcurrentTreeSetRW[string]{}
var _ SortedSet[string] = &ConcurrentTreeSetRW[string]{}
var _ MutableSortedSet[string] = &ConcurrentTreeSetRW[string]{}

// wrapConcurrentTreeSetRW builds a new ConcurrentTreeSetRW, with its own lock, around the given set.
func wrapConcurrentTreeSetRW[T constraints.Ordered](set *TreeSet[T]) *ConcurrentTreeSetRW[T] {
	return &ConcurrentTreeSetRW[T]{set: set, lock: &sync.RWMutex{}}
}

// Contains checks if the given element exists in the set.
func (ch *ConcurrentTreeSetRW[T]) Contains(element T) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.set.Contains(element)
}

// Length returns the number of elements in the set.
func (ch *ConcurrentTreeSetRW[T]) Length() int {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.set.Length()
}

// IsEmpty returns true if the set contains no elements.
func (ch *ConcurrentTreeSetRW[T]) IsEmpty() bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.set.IsEmpty()
}

// ForEach executes the given function for each element in ascending order. fn is
// invoked after the lock is released, against a point-in-time snapshot taken under
// the lock, so fn may safely call back into the collection.
func (ch *ConcurrentTreeSetRW[T]) ForEach(fn func(element T)) {
	ch.lock.RLock()
	elements := ch.set.AsSlice()
	ch.lock.RUnlock()

	for _, element := range elements {
		fn(element)
	}
}

// Filter returns a new set containing only the elements that satisfy the given
// predicate. The result is a new thread-safe ConcurrentTreeSetRW. The predicate is
// evaluated after the lock is released, against a point-in-time snapshot taken
// under the lock, so it may safely call back into the collection.
func (ch *ConcurrentTreeSetRW[T]) Filter(fn func(element T) bool) Set[T] {
	ch.lock.RLock()
	elements := ch.set.AsSlice()
	ch.lock.RUnlock()

	result := NewTreeSet[T]()
	for _, element := range elements {
		if fn(element) {
			result.AddInPlace(element)
		}
	}
	return wrapConcurrentTreeSetRW(result)
}

// FilterInPlace removes all elements that do not satisfy the given predicate,
// modifying the set in place. The predicate is evaluated after the lock is
// released, against a point-in-time snapshot taken under the lock, so it may
// safely call back into the collection. Modifications made concurrently with
// evaluation are not reflected in the retained set.
func (ch *ConcurrentTreeSetRW[T]) FilterInPlace(fn func(element T) bool) {
	ch.lock.RLock()
	elements := ch.set.AsSlice()
	ch.lock.RUnlock()

	var toRemove []T
	for _, element := range elements {
		if !fn(element) {
			toRemove = append(toRemove, element)
		}
	}

	ch.lock.Lock()
	ch.set.RemoveManyInPlace(toRemove...)
	ch.lock.Unlock()
}

// Find returns the first element (in ascending order) that satisfies the predicate.
// The predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (ch *ConcurrentTreeSetRW[T]) Find(fn func(element T) bool) (T, bool) {
	ch.lock.RLock()
	elements := ch.set.AsSlice()
	ch.lock.RUnlock()

	for _, element := range elements {
		if fn(element) {
			return element, true
		}
	}
	var zero T
	return zero, false
}

// AllMatch returns true if all elements satisfy the given predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (ch *ConcurrentTreeSetRW[T]) AllMatch(fn func(element T) bool) bool {
	ch.lock.RLock()
	elements := ch.set.AsSlice()
	ch.lock.RUnlock()

	for _, element := range elements {
		if !fn(element) {
			return false
		}
	}
	return true
}

// AnyMatch returns true if any element satisfies the given predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (ch *ConcurrentTreeSetRW[T]) AnyMatch(fn func(element T) bool) bool {
	ch.lock.RLock()
	elements := ch.set.AsSlice()
	ch.lock.RUnlock()

	for _, element := range elements {
		if fn(element) {
			return true
		}
	}
	return false
}

// NoneMatch returns true if no element satisfies the given predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (ch *ConcurrentTreeSetRW[T]) NoneMatch(fn func(element T) bool) bool {
	ch.lock.RLock()
	elements := ch.set.AsSlice()
	ch.lock.RUnlock()

	for _, element := range elements {
		if fn(element) {
			return false
		}
	}
	return true
}

// AsSlice returns the set as a slice in ascending order.
func (ch *ConcurrentTreeSetRW[T]) AsSlice() []T {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.set.AsSlice()
}

// AsMap returns the set as a native Go map with struct{} values.
func (ch *ConcurrentTreeSetRW[T]) AsMap() map[T]struct{} {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.set.AsMap()
}

// Add creates a new set with the given element added.
// Returns a new thread-safe ConcurrentTreeSetRW without modifying the original.
func (ch *ConcurrentTreeSetRW[T]) Add(element T) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return wrapConcurrentTreeSetRW(ch.set.Add(element).(*TreeSet[T]))
}

// AddMany creates a new set with all given elements added.
// Returns a new thread-safe ConcurrentTreeSetRW without modifying the original.
func (ch *ConcurrentTreeSetRW[T]) AddMany(elements ...T) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return wrapConcurrentTreeSetRW(ch.set.AddMany(elements...).(*TreeSet[T]))
}

// Union creates a new set containing all elements from this set and the other set.
// Returns a new thread-safe ConcurrentTreeSetRW without modifying the original.
func (ch *ConcurrentTreeSetRW[T]) Union(other Set[T]) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	// When other is the receiver, operate on the inner (non-locking) set so the
	// delegated call doesn't re-acquire ch.lock (recursive read locks can
	// deadlock with a queued writer).
	if other == Set[T](ch) {
		other = ch.set
	}
	return wrapConcurrentTreeSetRW(ch.set.Union(other).(*TreeSet[T]))
}

// AddInPlace adds the given element to the set.
func (ch *ConcurrentTreeSetRW[T]) AddInPlace(element T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.set.AddInPlace(element)
}

// AddManyInPlace adds all given elements to the set.
func (ch *ConcurrentTreeSetRW[T]) AddManyInPlace(elements ...T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.set.AddManyInPlace(elements...)
}

// UnionInPlace adds all elements from the other set to this set.
func (ch *ConcurrentTreeSetRW[T]) UnionInPlace(other Set[T]) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	if other == Set[T](ch) {
		other = ch.set
	}
	ch.set.UnionInPlace(other)
}

// Remove creates a new set with the given element removed.
// Returns a new thread-safe ConcurrentTreeSetRW without modifying the original.
func (ch *ConcurrentTreeSetRW[T]) Remove(element T) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return wrapConcurrentTreeSetRW(ch.set.Remove(element).(*TreeSet[T]))
}

// RemoveMany creates a new set with all given elements removed.
// Returns a new thread-safe ConcurrentTreeSetRW without modifying the original.
func (ch *ConcurrentTreeSetRW[T]) RemoveMany(elements ...T) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return wrapConcurrentTreeSetRW(ch.set.RemoveMany(elements...).(*TreeSet[T]))
}

// Difference creates a new set containing elements in this set but not in the other set.
// Returns a new thread-safe ConcurrentTreeSetRW without modifying the original.
func (ch *ConcurrentTreeSetRW[T]) Difference(other Set[T]) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	if other == Set[T](ch) {
		other = ch.set
	}
	return wrapConcurrentTreeSetRW(ch.set.Difference(other).(*TreeSet[T]))
}

// RemoveInPlace removes the given element from the set.
// Returns true if the element was present and removed; false otherwise.
func (ch *ConcurrentTreeSetRW[T]) RemoveInPlace(element T) bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.RemoveInPlace(element)
}

// RemoveManyInPlace removes all given elements from the set.
func (ch *ConcurrentTreeSetRW[T]) RemoveManyInPlace(elements ...T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.set.RemoveManyInPlace(elements...)
}

// DifferenceInPlace removes all elements that are present in the other set.
func (ch *ConcurrentTreeSetRW[T]) DifferenceInPlace(other Set[T]) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	if other == Set[T](ch) {
		other = ch.set
	}
	ch.set.DifferenceInPlace(other)
}

// Clear removes all elements from the set.
func (ch *ConcurrentTreeSetRW[T]) Clear() {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.set.Clear()
}

// Intersection creates a new set containing elements present in both sets.
// Returns a new thread-safe ConcurrentTreeSetRW without modifying the original.
func (ch *ConcurrentTreeSetRW[T]) Intersection(other Set[T]) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	if other == Set[T](ch) {
		other = ch.set
	}
	return wrapConcurrentTreeSetRW(ch.set.Intersection(other).(*TreeSet[T]))
}

// IsSubsetOf returns true if this set is a subset of the other set.
func (ch *ConcurrentTreeSetRW[T]) IsSubsetOf(other Set[T]) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	if other == Set[T](ch) {
		other = ch.set
	}
	return ch.set.IsSubsetOf(other)
}

// IsSupersetOf returns true if this set is a superset of the other set.
func (ch *ConcurrentTreeSetRW[T]) IsSupersetOf(other Set[T]) bool {
	return other.IsSubsetOf(ch)
}

// IsDisjoint returns true if this set has no elements in common with the other set.
func (ch *ConcurrentTreeSetRW[T]) IsDisjoint(other Set[T]) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	if other == Set[T](ch) {
		other = ch.set
	}
	return ch.set.IsDisjoint(other)
}

// Equals returns true if both sets contain exactly the same elements.
func (ch *ConcurrentTreeSetRW[T]) Equals(other Set[T]) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	if other == Set[T](ch) {
		other = ch.set
	}
	return ch.set.Equals(other)
}

// IntersectionInPlace keeps only elements that are present in both sets.
func (ch *ConcurrentTreeSetRW[T]) IntersectionInPlace(other Set[T]) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	if other == Set[T](ch) {
		other = ch.set
	}
	ch.set.IntersectionInPlace(other)
}

// Min returns the smallest element.
func (ch *ConcurrentTreeSetRW[T]) Min() (T, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.set.Min()
}

// Max returns the largest element.
func (ch *ConcurrentTreeSetRW[T]) Max() (T, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.set.Max()
}

// Floor returns the largest element less than or equal to the given element.
func (ch *ConcurrentTreeSetRW[T]) Floor(element T) (T, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.set.Floor(element)
}

// Ceiling returns the smallest element greater than or equal to the given element.
func (ch *ConcurrentTreeSetRW[T]) Ceiling(element T) (T, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.set.Ceiling(element)
}

// Range returns all elements within the inclusive range [lo, hi], in ascending
// order. Returns a non-nil (possibly empty) slice.
func (ch *ConcurrentTreeSetRW[T]) Range(lo, hi T) []T {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return ch.set.Range(lo, hi)
}

// All returns an iterator over all elements in ascending order. The elements are
// snapshotted under the read lock, so iteration is safe against concurrent
// mutation and never holds the lock while calling the consumer.
func (ch *ConcurrentTreeSetRW[T]) All() iter.Seq[T] {
	ch.lock.RLock()
	elements := ch.set.AsSlice()
	ch.lock.RUnlock()
	return seqFromSlice(elements)
}

// Backward returns an iterator over all elements in descending order, snapshotted under the read lock.
func (ch *ConcurrentTreeSetRW[T]) Backward() iter.Seq[T] {
	ch.lock.RLock()
	elements := ch.set.AsSlice()
	ch.lock.RUnlock()
	return seqFromSliceReverse(elements)
}

// RangeAll returns an iterator over the elements within the inclusive range
// [lo, hi], in ascending order, snapshotted under the read lock.
func (ch *ConcurrentTreeSetRW[T]) RangeAll(lo, hi T) iter.Seq[T] {
	ch.lock.RLock()
	elements := ch.set.Range(lo, hi)
	ch.lock.RUnlock()
	return seqFromSlice(elements)
}
