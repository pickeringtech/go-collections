package sets

import "sync"

// ConcurrentHashRW is a thread-safe set implementation using Go's built-in map
// with a read-write mutex for synchronization. Read operations use read locks for better
// performance when there are many concurrent readers.
type ConcurrentHashRW[T comparable] struct {
	data map[T]struct{}
	lock sync.RWMutex
}

// NewConcurrentHashRW creates a new ConcurrentHashRW set with the given elements.
func NewConcurrentHashRW[T comparable](values ...T) *ConcurrentHashRW[T] {
	s := &ConcurrentHashRW[T]{
		data: make(map[T]struct{}),
	}
	for _, value := range values {
		s.data[value] = struct{}{}
	}
	return s
}

// Interface guards to ensure ConcurrentHashRW implements the required interfaces
var _ Set[string] = &ConcurrentHashRW[string]{}
var _ MutableSet[string] = &ConcurrentHashRW[string]{}

// Contains checks if the given element exists in the set.
func (ch *ConcurrentHashRW[T]) Contains(element T) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	_, exists := ch.data[element]
	return exists
}

// Length returns the number of elements in the set.
func (ch *ConcurrentHashRW[T]) Length() int {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	return len(ch.data)
}

// IsEmpty returns true if the set contains no elements.
func (ch *ConcurrentHashRW[T]) IsEmpty() bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	return len(ch.data) == 0
}

// ForEach executes the given function for each element. fn is invoked after the
// lock is released, against a point-in-time snapshot taken under the lock, so fn
// may safely call back into the collection.
func (ch *ConcurrentHashRW[T]) ForEach(fn func(element T)) {
	ch.lock.RLock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
	ch.lock.RUnlock()

	for _, element := range elements {
		fn(element)
	}
}

// Filter returns a new set containing only the elements
// that satisfy the given predicate function. The returned set is a new
// thread-safe ConcurrentHashRW, independent of the receiver. The predicate is
// evaluated after the lock is released, against a point-in-time snapshot taken
// under the lock, so it may safely call back into the collection.
func (ch *ConcurrentHashRW[T]) Filter(fn func(element T) bool) Set[T] {
	ch.lock.RLock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
	ch.lock.RUnlock()

	result := NewConcurrentHashRW[T]()
	for _, element := range elements {
		if fn(element) {
			result.data[element] = struct{}{}
		}
	}
	return result
}

// FilterInPlace removes all elements that do not satisfy the given predicate
// function, modifying the set in place. The predicate is evaluated after the
// lock is released, against a point-in-time snapshot taken under the lock, so
// it may safely call back into the collection.
//
// Only elements the predicate rejected are removed, and only if still present
// at apply time, so elements added concurrently in the evaluation window are
// preserved.
func (ch *ConcurrentHashRW[T]) FilterInPlace(fn func(element T) bool) {
	ch.lock.RLock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
	ch.lock.RUnlock()

	var toRemove []T
	for _, element := range elements {
		if !fn(element) {
			toRemove = append(toRemove, element)
		}
	}

	ch.lock.Lock()
	for _, element := range toRemove {
		delete(ch.data, element)
	}
	ch.lock.Unlock()
}

// Find returns the first element that satisfies the given predicate.
// Returns the element and true if found; zero value and false otherwise. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (ch *ConcurrentHashRW[T]) Find(fn func(element T) bool) (T, bool) {
	ch.lock.RLock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
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
func (ch *ConcurrentHashRW[T]) AllMatch(fn func(element T) bool) bool {
	ch.lock.RLock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
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
func (ch *ConcurrentHashRW[T]) AnyMatch(fn func(element T) bool) bool {
	ch.lock.RLock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
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
func (ch *ConcurrentHashRW[T]) NoneMatch(fn func(element T) bool) bool {
	ch.lock.RLock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
	ch.lock.RUnlock()

	for _, element := range elements {
		if fn(element) {
			return false
		}
	}
	return true
}

// AsSlice returns the set as a slice.
func (ch *ConcurrentHashRW[T]) AsSlice() []T {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make([]T, 0, len(ch.data))
	for element := range ch.data {
		result = append(result, element)
	}
	return result
}

// AsMap returns the set as a native Go map with struct{} values.
func (ch *ConcurrentHashRW[T]) AsMap() map[T]struct{} {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(map[T]struct{}, len(ch.data))
	for element := range ch.data {
		result[element] = struct{}{}
	}
	return result
}

// Add creates a new set with the given element added.
// Returns a new thread-safe ConcurrentHashRW without modifying the original.
func (ch *ConcurrentHashRW[T]) Add(element T) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := NewConcurrentHashRW[T]()
	for e := range ch.data {
		result.data[e] = struct{}{}
	}
	result.data[element] = struct{}{}
	return result
}

// AddMany creates a new set with all given elements added.
// Returns a new thread-safe ConcurrentHashRW without modifying the original.
func (ch *ConcurrentHashRW[T]) AddMany(elements ...T) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := NewConcurrentHashRW[T]()
	for e := range ch.data {
		result.data[e] = struct{}{}
	}
	for _, element := range elements {
		result.data[element] = struct{}{}
	}
	return result
}

// Union creates a new set containing all elements from this set and the other set.
// Returns a new thread-safe ConcurrentHashRW without modifying the original.
func (ch *ConcurrentHashRW[T]) Union(other Set[T]) Set[T] {
	// Snapshot other before locking ch. Holding ch.lock while calling a locking
	// method on other inverts lock order against the opposite-operand call and
	// deadlocks (a.Union(b) racing b.Union(a)). Self-union needs no snapshot.
	var otherElements []T
	if other != Set[T](ch) {
		otherElements = other.AsSlice()
	}

	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := NewConcurrentHashRW[T]()
	for e := range ch.data {
		result.data[e] = struct{}{}
	}
	for _, element := range otherElements {
		result.data[element] = struct{}{}
	}
	return result
}

// AddInPlace adds the given element to the set.
func (ch *ConcurrentHashRW[T]) AddInPlace(element T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	ch.data[element] = struct{}{}
}

// AddManyInPlace adds all given elements to the set.
func (ch *ConcurrentHashRW[T]) AddManyInPlace(elements ...T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for _, element := range elements {
		ch.data[element] = struct{}{}
	}
}

// UnionInPlace adds all elements from the other set to this set.
func (ch *ConcurrentHashRW[T]) UnionInPlace(other Set[T]) {
	// Unioning a set with itself is a no-op; otherwise snapshot other before
	// locking ch so the snapshot's own lock acquisition can't invert lock order
	// and deadlock against b.UnionInPlace(a).
	if other == Set[T](ch) {
		return
	}
	otherElements := other.AsSlice()

	ch.lock.Lock()
	defer ch.lock.Unlock()
	for _, element := range otherElements {
		ch.data[element] = struct{}{}
	}
}

// Remove creates a new set with the given element removed.
// Returns a new thread-safe ConcurrentHashRW without modifying the original.
func (ch *ConcurrentHashRW[T]) Remove(element T) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := NewConcurrentHashRW[T]()
	for e := range ch.data {
		if e != element {
			result.data[e] = struct{}{}
		}
	}
	return result
}

// RemoveMany creates a new set with all given elements removed.
// Returns a new thread-safe ConcurrentHashRW without modifying the original.
func (ch *ConcurrentHashRW[T]) RemoveMany(elements ...T) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	// Create a set of elements to remove for efficient lookup
	toRemove := make(map[T]struct{}, len(elements))
	for _, element := range elements {
		toRemove[element] = struct{}{}
	}

	result := NewConcurrentHashRW[T]()
	for e := range ch.data {
		if _, shouldRemove := toRemove[e]; !shouldRemove {
			result.data[e] = struct{}{}
		}
	}
	return result
}

// Difference creates a new set containing elements in this set but not in the other set.
// Returns a new thread-safe ConcurrentHashRW without modifying the original.
func (ch *ConcurrentHashRW[T]) Difference(other Set[T]) Set[T] {
	result := NewConcurrentHashRW[T]()
	// Self-difference is empty; otherwise snapshot other before locking ch so the
	// snapshot can't invert lock order and deadlock against b.Difference(a).
	if other == Set[T](ch) {
		return result
	}
	otherElements := other.AsMap()

	ch.lock.RLock()
	defer ch.lock.RUnlock()
	for element := range ch.data {
		_, inOther := otherElements[element]
		if !inOther {
			result.data[element] = struct{}{}
		}
	}
	return result
}

// RemoveInPlace removes the given element from the set.
// Returns true if the element was present and removed; false otherwise.
func (ch *ConcurrentHashRW[T]) RemoveInPlace(element T) bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	if _, exists := ch.data[element]; exists {
		delete(ch.data, element)
		return true
	}
	return false
}

// RemoveManyInPlace removes all given elements from the set.
func (ch *ConcurrentHashRW[T]) RemoveManyInPlace(elements ...T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for _, element := range elements {
		delete(ch.data, element)
	}
}

// DifferenceInPlace removes all elements that are present in the other set.
func (ch *ConcurrentHashRW[T]) DifferenceInPlace(other Set[T]) {
	// Removing a set's own elements from itself empties it; handle that directly.
	// Otherwise snapshot other before locking ch so the snapshot can't invert
	// lock order and deadlock against b.DifferenceInPlace(a).
	if other == Set[T](ch) {
		ch.lock.Lock()
		clear(ch.data)
		ch.lock.Unlock()
		return
	}
	otherElements := other.AsSlice()

	ch.lock.Lock()
	defer ch.lock.Unlock()
	for _, element := range otherElements {
		delete(ch.data, element)
	}
}

// Clear removes all elements from the set.
func (ch *ConcurrentHashRW[T]) Clear() {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for element := range ch.data {
		delete(ch.data, element)
	}
}

// Intersection creates a new set containing elements present in both sets.
// Returns a new thread-safe ConcurrentHashRW without modifying the original.
func (ch *ConcurrentHashRW[T]) Intersection(other Set[T]) Set[T] {
	result := NewConcurrentHashRW[T]()
	// Self-intersection is just a copy; otherwise snapshot other before locking ch
	// so the snapshot can't invert lock order and deadlock against b.Intersection(a).
	if other == Set[T](ch) {
		ch.lock.RLock()
		defer ch.lock.RUnlock()
		for element := range ch.data {
			result.data[element] = struct{}{}
		}
		return result
	}
	otherElements := other.AsMap()

	ch.lock.RLock()
	defer ch.lock.RUnlock()
	for element := range ch.data {
		_, inOther := otherElements[element]
		if inOther {
			result.data[element] = struct{}{}
		}
	}
	return result
}

// IsSubsetOf returns true if this set is a subset of the other set.
func (ch *ConcurrentHashRW[T]) IsSubsetOf(other Set[T]) bool {
	// A set is always a subset of itself; otherwise snapshot other before locking
	// ch so the snapshot can't invert lock order and deadlock against
	// b.IsSubsetOf(a) (or the AB/BA superset race, since IsSupersetOf delegates here).
	if other == Set[T](ch) {
		return true
	}
	otherElements := other.AsMap()

	ch.lock.RLock()
	defer ch.lock.RUnlock()
	for element := range ch.data {
		_, inOther := otherElements[element]
		if !inOther {
			return false
		}
	}
	return true
}

// IsSupersetOf returns true if this set is a superset of the other set.
func (ch *ConcurrentHashRW[T]) IsSupersetOf(other Set[T]) bool {
	return other.IsSubsetOf(ch)
}

// IsDisjoint returns true if this set has no elements in common with the other set.
func (ch *ConcurrentHashRW[T]) IsDisjoint(other Set[T]) bool {
	// A set is disjoint with itself only when empty; otherwise snapshot other
	// before locking ch so the snapshot can't invert lock order and deadlock
	// against b.IsDisjoint(a).
	if other == Set[T](ch) {
		ch.lock.RLock()
		defer ch.lock.RUnlock()
		return len(ch.data) == 0
	}
	otherElements := other.AsMap()

	ch.lock.RLock()
	defer ch.lock.RUnlock()
	for element := range ch.data {
		_, inOther := otherElements[element]
		if inOther {
			return false
		}
	}
	return true
}

// Equals returns true if both sets contain exactly the same elements.
func (ch *ConcurrentHashRW[T]) Equals(other Set[T]) bool {
	// A set always equals itself; otherwise snapshot other before locking ch so
	// the snapshot can't invert lock order and deadlock against b.Equals(a).
	if other == Set[T](ch) {
		return true
	}
	otherElements := other.AsMap()

	ch.lock.RLock()
	defer ch.lock.RUnlock()
	if len(ch.data) != len(otherElements) {
		return false
	}
	for element := range ch.data {
		_, inOther := otherElements[element]
		if !inOther {
			return false
		}
	}
	return true
}

// IntersectionInPlace keeps only elements that are present in both sets.
func (ch *ConcurrentHashRW[T]) IntersectionInPlace(other Set[T]) {
	// Intersecting a set with itself is a no-op; otherwise snapshot other before
	// locking ch so the snapshot can't invert lock order and deadlock against
	// b.IntersectionInPlace(a).
	if other == Set[T](ch) {
		return
	}
	otherElements := other.AsMap()

	ch.lock.Lock()
	defer ch.lock.Unlock()
	for element := range ch.data {
		_, inOther := otherElements[element]
		if !inOther {
			delete(ch.data, element)
		}
	}
}
