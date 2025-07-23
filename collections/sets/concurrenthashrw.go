package sets

import "sync"

// ConcurrentHashRW is a thread-safe set implementation using Go's built-in map
// with a read-write mutex for synchronization. Read operations use read locks for better
// performance when there are many concurrent readers.
type ConcurrentHashRW[T comparable] struct {
	data map[T]struct{}
	lock *sync.RWMutex
}

// NewConcurrentHashRW creates a new ConcurrentHashRW set with the given elements.
func NewConcurrentHashRW[T comparable](values ...T) *ConcurrentHashRW[T] {
	s := &ConcurrentHashRW[T]{
		data: make(map[T]struct{}),
		lock: &sync.RWMutex{},
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

// ForEach executes the given function for each element.
func (ch *ConcurrentHashRW[T]) ForEach(fn func(element T)) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for element := range ch.data {
		fn(element)
	}
}

// Filter returns a new set containing only the elements
// that satisfy the given predicate function.
func (ch *ConcurrentHashRW[T]) Filter(fn func(element T) bool) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(Hash[T])
	for element := range ch.data {
		if fn(element) {
			result[element] = struct{}{}
		}
	}
	return result
}

// FilterInPlace removes all elements that do not satisfy
// the given predicate function, modifying the set in place.
func (ch *ConcurrentHashRW[T]) FilterInPlace(fn func(element T) bool) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for element := range ch.data {
		if !fn(element) {
			delete(ch.data, element)
		}
	}
}

// Find returns the first element that satisfies the given predicate.
// Returns the element and true if found; zero value and false otherwise.
func (ch *ConcurrentHashRW[T]) Find(fn func(element T) bool) (T, bool) {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for element := range ch.data {
		if fn(element) {
			return element, true
		}
	}
	var zero T
	return zero, false
}

// AllMatch returns true if all elements satisfy the given predicate.
func (ch *ConcurrentHashRW[T]) AllMatch(fn func(element T) bool) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for element := range ch.data {
		if !fn(element) {
			return false
		}
	}
	return true
}

// AnyMatch returns true if any element satisfies the given predicate.
func (ch *ConcurrentHashRW[T]) AnyMatch(fn func(element T) bool) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for element := range ch.data {
		if fn(element) {
			return true
		}
	}
	return false
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
// Returns the new set without modifying the original.
func (ch *ConcurrentHashRW[T]) Add(element T) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(Hash[T], len(ch.data)+1)
	for e := range ch.data {
		result[e] = struct{}{}
	}
	result[element] = struct{}{}
	return result
}

// AddMany creates a new set with all given elements added.
// Returns the new set without modifying the original.
func (ch *ConcurrentHashRW[T]) AddMany(elements ...T) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(Hash[T], len(ch.data)+len(elements))
	for e := range ch.data {
		result[e] = struct{}{}
	}
	for _, element := range elements {
		result[element] = struct{}{}
	}
	return result
}

// Union creates a new set containing all elements from this set and the other set.
func (ch *ConcurrentHashRW[T]) Union(other Set[T]) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(Hash[T], ch.Length()+other.Length())
	for e := range ch.data {
		result[e] = struct{}{}
	}
	other.ForEach(func(element T) {
		result[element] = struct{}{}
	})
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
	ch.lock.Lock()
	defer ch.lock.Unlock()

	other.ForEach(func(element T) {
		ch.data[element] = struct{}{}
	})
}

// Remove creates a new set with the given element removed.
// Returns the new set without modifying the original.
func (ch *ConcurrentHashRW[T]) Remove(element T) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(Hash[T], len(ch.data))
	for e := range ch.data {
		if e != element {
			result[e] = struct{}{}
		}
	}
	return result
}

// RemoveMany creates a new set with all given elements removed.
// Returns the new set without modifying the original.
func (ch *ConcurrentHashRW[T]) RemoveMany(elements ...T) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	// Create a set of elements to remove for efficient lookup
	toRemove := make(map[T]struct{}, len(elements))
	for _, element := range elements {
		toRemove[element] = struct{}{}
	}

	result := make(Hash[T], len(ch.data))
	for e := range ch.data {
		if _, shouldRemove := toRemove[e]; !shouldRemove {
			result[e] = struct{}{}
		}
	}
	return result
}

// Difference creates a new set containing elements in this set but not in the other set.
func (ch *ConcurrentHashRW[T]) Difference(other Set[T]) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(Hash[T])
	for element := range ch.data {
		if !other.Contains(element) {
			result[element] = struct{}{}
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
	ch.lock.Lock()
	defer ch.lock.Unlock()

	other.ForEach(func(element T) {
		delete(ch.data, element)
	})
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
func (ch *ConcurrentHashRW[T]) Intersection(other Set[T]) Set[T] {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	result := make(Hash[T])
	for element := range ch.data {
		if other.Contains(element) {
			result[element] = struct{}{}
		}
	}
	return result
}

// IsSubsetOf returns true if this set is a subset of the other set.
func (ch *ConcurrentHashRW[T]) IsSubsetOf(other Set[T]) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for element := range ch.data {
		if !other.Contains(element) {
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
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	for element := range ch.data {
		if other.Contains(element) {
			return false
		}
	}
	return true
}

// Equals returns true if both sets contain exactly the same elements.
func (ch *ConcurrentHashRW[T]) Equals(other Set[T]) bool {
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	if len(ch.data) != other.Length() {
		return false
	}
	for element := range ch.data {
		if !other.Contains(element) {
			return false
		}
	}
	return true
}

// IntersectionInPlace keeps only elements that are present in both sets.
func (ch *ConcurrentHashRW[T]) IntersectionInPlace(other Set[T]) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for element := range ch.data {
		if !other.Contains(element) {
			delete(ch.data, element)
		}
	}
}
