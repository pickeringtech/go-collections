package sets

import (
	"iter"
	"sync"

	"github.com/pickeringtech/go-collections/constraints"
)

// ConcurrentTreeSet is a thread-safe sorted set backed by a TreeSet with a mutex
// for synchronization. Elements are kept in sorted order and all operations are
// protected by a single mutex. Use it when reads and writes are balanced; prefer
// ConcurrentTreeSetRW for read-heavy workloads.
type ConcurrentTreeSet[T constraints.Ordered] struct {
	set  *TreeSet[T]
	lock *sync.Mutex
}

// NewConcurrentTreeSet creates a new ConcurrentTreeSet with the given elements.
// Duplicate elements are automatically removed.
func NewConcurrentTreeSet[T constraints.Ordered](elements ...T) *ConcurrentTreeSet[T] {
	return &ConcurrentTreeSet[T]{
		set:  NewTreeSet[T](elements...),
		lock: &sync.Mutex{},
	}
}

// Interface guards to ensure ConcurrentTreeSet implements the required interfaces.
var _ Set[string] = &ConcurrentTreeSet[string]{}
var _ MutableSet[string] = &ConcurrentTreeSet[string]{}
var _ SortedSet[string] = &ConcurrentTreeSet[string]{}
var _ MutableSortedSet[string] = &ConcurrentTreeSet[string]{}

// wrapConcurrentTreeSet builds a new ConcurrentTreeSet, with its own lock, around the given set.
func wrapConcurrentTreeSet[T constraints.Ordered](set *TreeSet[T]) *ConcurrentTreeSet[T] {
	return &ConcurrentTreeSet[T]{set: set, lock: &sync.Mutex{}}
}

// Contains checks if the given element exists in the set.
func (ch *ConcurrentTreeSet[T]) Contains(element T) bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.Contains(element)
}

// Length returns the number of elements in the set.
func (ch *ConcurrentTreeSet[T]) Length() int {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.Length()
}

// IsEmpty returns true if the set contains no elements.
func (ch *ConcurrentTreeSet[T]) IsEmpty() bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.IsEmpty()
}

// ForEach executes the given function for each element in ascending order. fn is
// invoked after the lock is released, against a point-in-time snapshot taken under
// the lock, so fn may safely call back into the collection.
func (ch *ConcurrentTreeSet[T]) ForEach(fn func(element T)) {
	ch.lock.Lock()
	elements := ch.set.AsSlice()
	ch.lock.Unlock()

	for _, element := range elements {
		fn(element)
	}
}

// Filter returns a new set containing only the elements that satisfy the given
// predicate. The result is a new thread-safe ConcurrentTreeSet. The predicate is
// evaluated after the lock is released, against a point-in-time snapshot taken
// under the lock, so it may safely call back into the collection.
func (ch *ConcurrentTreeSet[T]) Filter(fn func(element T) bool) Set[T] {
	ch.lock.Lock()
	elements := ch.set.AsSlice()
	ch.lock.Unlock()

	result := NewTreeSet[T]()
	for _, element := range elements {
		if fn(element) {
			result.AddInPlace(element)
		}
	}
	return wrapConcurrentTreeSet(result)
}

// FilterInPlace removes all elements that do not satisfy the given predicate,
// modifying the set in place. The predicate is evaluated after the lock is
// released, against a point-in-time snapshot taken under the lock, so it may
// safely call back into the collection. Modifications made concurrently with
// evaluation are not reflected in the retained set.
func (ch *ConcurrentTreeSet[T]) FilterInPlace(fn func(element T) bool) {
	ch.lock.Lock()
	elements := ch.set.AsSlice()
	ch.lock.Unlock()

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
func (ch *ConcurrentTreeSet[T]) Find(fn func(element T) bool) (T, bool) {
	ch.lock.Lock()
	elements := ch.set.AsSlice()
	ch.lock.Unlock()

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
func (ch *ConcurrentTreeSet[T]) AllMatch(fn func(element T) bool) bool {
	ch.lock.Lock()
	elements := ch.set.AsSlice()
	ch.lock.Unlock()

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
func (ch *ConcurrentTreeSet[T]) AnyMatch(fn func(element T) bool) bool {
	ch.lock.Lock()
	elements := ch.set.AsSlice()
	ch.lock.Unlock()

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
func (ch *ConcurrentTreeSet[T]) NoneMatch(fn func(element T) bool) bool {
	ch.lock.Lock()
	elements := ch.set.AsSlice()
	ch.lock.Unlock()

	for _, element := range elements {
		if fn(element) {
			return false
		}
	}
	return true
}

// AsSlice returns the set as a slice in ascending order.
func (ch *ConcurrentTreeSet[T]) AsSlice() []T {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.AsSlice()
}

// AsMap returns the set as a native Go map with struct{} values.
func (ch *ConcurrentTreeSet[T]) AsMap() map[T]struct{} {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.AsMap()
}

// Add creates a new set with the given element added.
// Returns a new thread-safe ConcurrentTreeSet without modifying the original.
func (ch *ConcurrentTreeSet[T]) Add(element T) Set[T] {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return wrapConcurrentTreeSet(ch.set.Add(element).(*TreeSet[T]))
}

// AddMany creates a new set with all given elements added.
// Returns a new thread-safe ConcurrentTreeSet without modifying the original.
func (ch *ConcurrentTreeSet[T]) AddMany(elements ...T) Set[T] {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return wrapConcurrentTreeSet(ch.set.AddMany(elements...).(*TreeSet[T]))
}

// Union creates a new set containing all elements from this set and the other set.
// Returns a new thread-safe ConcurrentTreeSet without modifying the original.
func (ch *ConcurrentTreeSet[T]) Union(other Set[T]) Set[T] {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return wrapConcurrentTreeSet(ch.set.Union(other).(*TreeSet[T]))
}

// AddInPlace adds the given element to the set.
func (ch *ConcurrentTreeSet[T]) AddInPlace(element T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.set.AddInPlace(element)
}

// AddManyInPlace adds all given elements to the set.
func (ch *ConcurrentTreeSet[T]) AddManyInPlace(elements ...T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.set.AddManyInPlace(elements...)
}

// UnionInPlace adds all elements from the other set to this set.
func (ch *ConcurrentTreeSet[T]) UnionInPlace(other Set[T]) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.set.UnionInPlace(other)
}

// Remove creates a new set with the given element removed.
// Returns a new thread-safe ConcurrentTreeSet without modifying the original.
func (ch *ConcurrentTreeSet[T]) Remove(element T) Set[T] {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return wrapConcurrentTreeSet(ch.set.Remove(element).(*TreeSet[T]))
}

// RemoveMany creates a new set with all given elements removed.
// Returns a new thread-safe ConcurrentTreeSet without modifying the original.
func (ch *ConcurrentTreeSet[T]) RemoveMany(elements ...T) Set[T] {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return wrapConcurrentTreeSet(ch.set.RemoveMany(elements...).(*TreeSet[T]))
}

// Difference creates a new set containing elements in this set but not in the other set.
// Returns a new thread-safe ConcurrentTreeSet without modifying the original.
func (ch *ConcurrentTreeSet[T]) Difference(other Set[T]) Set[T] {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return wrapConcurrentTreeSet(ch.set.Difference(other).(*TreeSet[T]))
}

// RemoveInPlace removes the given element from the set.
// Returns true if the element was present and removed; false otherwise.
func (ch *ConcurrentTreeSet[T]) RemoveInPlace(element T) bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.RemoveInPlace(element)
}

// RemoveManyInPlace removes all given elements from the set.
func (ch *ConcurrentTreeSet[T]) RemoveManyInPlace(elements ...T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.set.RemoveManyInPlace(elements...)
}

// DifferenceInPlace removes all elements that are present in the other set.
func (ch *ConcurrentTreeSet[T]) DifferenceInPlace(other Set[T]) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.set.DifferenceInPlace(other)
}

// Clear removes all elements from the set.
func (ch *ConcurrentTreeSet[T]) Clear() {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.set.Clear()
}

// Intersection creates a new set containing elements present in both sets.
// Returns a new thread-safe ConcurrentTreeSet without modifying the original.
func (ch *ConcurrentTreeSet[T]) Intersection(other Set[T]) Set[T] {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return wrapConcurrentTreeSet(ch.set.Intersection(other).(*TreeSet[T]))
}

// IsSubsetOf returns true if this set is a subset of the other set.
func (ch *ConcurrentTreeSet[T]) IsSubsetOf(other Set[T]) bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.IsSubsetOf(other)
}

// IsSupersetOf returns true if this set is a superset of the other set.
func (ch *ConcurrentTreeSet[T]) IsSupersetOf(other Set[T]) bool {
	return other.IsSubsetOf(ch)
}

// IsDisjoint returns true if this set has no elements in common with the other set.
func (ch *ConcurrentTreeSet[T]) IsDisjoint(other Set[T]) bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.IsDisjoint(other)
}

// Equals returns true if both sets contain exactly the same elements.
func (ch *ConcurrentTreeSet[T]) Equals(other Set[T]) bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.Equals(other)
}

// IntersectionInPlace keeps only elements that are present in both sets.
func (ch *ConcurrentTreeSet[T]) IntersectionInPlace(other Set[T]) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.set.IntersectionInPlace(other)
}

// Min returns the smallest element.
func (ch *ConcurrentTreeSet[T]) Min() (T, bool) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.Min()
}

// Max returns the largest element.
func (ch *ConcurrentTreeSet[T]) Max() (T, bool) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.Max()
}

// Floor returns the largest element less than or equal to the given element.
func (ch *ConcurrentTreeSet[T]) Floor(element T) (T, bool) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.Floor(element)
}

// Ceiling returns the smallest element greater than or equal to the given element.
func (ch *ConcurrentTreeSet[T]) Ceiling(element T) (T, bool) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.Ceiling(element)
}

// Range returns all elements within the inclusive range [lo, hi], in ascending
// order. Returns a non-nil (possibly empty) slice.
func (ch *ConcurrentTreeSet[T]) Range(lo, hi T) []T {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	return ch.set.Range(lo, hi)
}

// All returns an iterator over all elements in ascending order. The elements are
// snapshotted under the lock, so iteration is safe against concurrent mutation
// and never holds the lock while calling the consumer.
func (ch *ConcurrentTreeSet[T]) All() iter.Seq[T] {
	ch.lock.Lock()
	elements := ch.set.AsSlice()
	ch.lock.Unlock()
	return seqFromSlice(elements)
}

// Backward returns an iterator over all elements in descending order, snapshotted under the lock.
func (ch *ConcurrentTreeSet[T]) Backward() iter.Seq[T] {
	ch.lock.Lock()
	elements := ch.set.AsSlice()
	ch.lock.Unlock()
	return seqFromSliceReverse(elements)
}

// RangeAll returns an iterator over the elements within the inclusive range
// [lo, hi], in ascending order, snapshotted under the lock.
func (ch *ConcurrentTreeSet[T]) RangeAll(lo, hi T) iter.Seq[T] {
	ch.lock.Lock()
	elements := ch.set.Range(lo, hi)
	ch.lock.Unlock()
	return seqFromSlice(elements)
}

// seqFromSlice returns an iterator that yields the given elements in order.
func seqFromSlice[T constraints.Ordered](elements []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, element := range elements {
			if !yield(element) {
				return
			}
		}
	}
}

// seqFromSliceReverse returns an iterator that yields the given elements in reverse order.
func seqFromSliceReverse[T constraints.Ordered](elements []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(elements) - 1; i >= 0; i-- {
			if !yield(elements[i]) {
				return
			}
		}
	}
}
