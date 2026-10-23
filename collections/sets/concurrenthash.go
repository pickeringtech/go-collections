package sets

import (
	"sync"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
)

// ConcurrentHash is a thread-safe set implementation using Go's built-in map
// with a mutex for synchronization. All operations are protected by a single mutex.
//
// ConcurrentHash must not be copied after first use; copying after construction
// produces an independent lock over shared backing data, which breaks the
// thread-safety contract. go vet reports any such copy.
type ConcurrentHash[T comparable] struct {
	_    nocopy.NoCopy
	data map[T]struct{}
	lock sync.Mutex
}

// NewConcurrentHash creates a new ConcurrentHash set with the given elements.
func NewConcurrentHash[T comparable](values ...T) *ConcurrentHash[T] {
	s := &ConcurrentHash[T]{
		data: make(map[T]struct{}),
	}
	for _, value := range values {
		s.data[value] = struct{}{}
	}
	return s
}

// Interface guards to ensure ConcurrentHash implements the required interfaces
var _ Set[string] = &ConcurrentHash[string]{}
var _ MutableSet[string] = &ConcurrentHash[string]{}

// Contains checks if the given element exists in the set.
func (ch *ConcurrentHash[T]) Contains(element T) bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	_, exists := ch.data[element]
	return exists
}

// Length returns the number of elements in the set.
func (ch *ConcurrentHash[T]) Length() int {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	return len(ch.data)
}

// IsEmpty returns true if the set contains no elements.
func (ch *ConcurrentHash[T]) IsEmpty() bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	return len(ch.data) == 0
}

// ForEach executes the given function for each element. fn is invoked after the
// lock is released, against a point-in-time snapshot taken under the lock, so fn
// may safely call back into the collection.
func (ch *ConcurrentHash[T]) ForEach(fn func(element T)) {
	ch.lock.Lock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
	ch.lock.Unlock()

	for _, element := range elements {
		fn(element)
	}
}

// Filter returns a new set containing only the elements
// that satisfy the given predicate function. The returned set is a new
// thread-safe ConcurrentHash, independent of the receiver. The predicate is
// evaluated after the lock is released, against a point-in-time snapshot taken
// under the lock, so it may safely call back into the collection.
func (ch *ConcurrentHash[T]) Filter(fn func(element T) bool) Set[T] {
	ch.lock.Lock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
	ch.lock.Unlock()

	result := NewConcurrentHash[T]()
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
func (ch *ConcurrentHash[T]) FilterInPlace(fn func(element T) bool) {
	ch.lock.Lock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
	ch.lock.Unlock()

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
func (ch *ConcurrentHash[T]) Find(fn func(element T) bool) (T, bool) {
	ch.lock.Lock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
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
func (ch *ConcurrentHash[T]) AllMatch(fn func(element T) bool) bool {
	ch.lock.Lock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
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
func (ch *ConcurrentHash[T]) AnyMatch(fn func(element T) bool) bool {
	ch.lock.Lock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
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
func (ch *ConcurrentHash[T]) NoneMatch(fn func(element T) bool) bool {
	ch.lock.Lock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
	ch.lock.Unlock()

	for _, element := range elements {
		if fn(element) {
			return false
		}
	}
	return true
}

// AsSlice returns the set as a slice.
func (ch *ConcurrentHash[T]) AsSlice() []T {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	result := make([]T, 0, len(ch.data))
	for element := range ch.data {
		result = append(result, element)
	}
	return result
}

// AsMap returns the set as a native Go map with struct{} values.
func (ch *ConcurrentHash[T]) AsMap() map[T]struct{} {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	result := make(map[T]struct{}, len(ch.data))
	for element := range ch.data {
		result[element] = struct{}{}
	}
	return result
}

// Add creates a new set with the given element added.
// Returns a new thread-safe ConcurrentHash without modifying the original.
func (ch *ConcurrentHash[T]) Add(element T) Set[T] {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	result := NewConcurrentHash[T]()
	for e := range ch.data {
		result.data[e] = struct{}{}
	}
	result.data[element] = struct{}{}
	return result
}

// AddMany creates a new set with all given elements added.
// Returns a new thread-safe ConcurrentHash without modifying the original.
func (ch *ConcurrentHash[T]) AddMany(elements ...T) Set[T] {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	result := NewConcurrentHash[T]()
	for e := range ch.data {
		result.data[e] = struct{}{}
	}
	for _, element := range elements {
		result.data[element] = struct{}{}
	}
	return result
}

// Union creates a new set containing all elements from this set and the other set.
// Returns a new thread-safe ConcurrentHash without modifying the original.
func (ch *ConcurrentHash[T]) Union(other Set[T]) Set[T] {
	// Snapshot other before locking ch. Holding ch.lock while calling a locking
	// method on other inverts lock order against the opposite-operand call and
	// deadlocks (a.Union(b) racing b.Union(a)). Self-union needs no snapshot.
	var otherElements []T
	if other != Set[T](ch) {
		otherElements = other.AsSlice()
	}

	ch.lock.Lock()
	defer ch.lock.Unlock()

	result := NewConcurrentHash[T]()
	for e := range ch.data {
		result.data[e] = struct{}{}
	}
	for _, element := range otherElements {
		result.data[element] = struct{}{}
	}
	return result
}

// AddInPlace adds the given element to the set.
func (ch *ConcurrentHash[T]) AddInPlace(element T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	ch.data[element] = struct{}{}
}

// AddManyInPlace adds all given elements to the set.
func (ch *ConcurrentHash[T]) AddManyInPlace(elements ...T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for _, element := range elements {
		ch.data[element] = struct{}{}
	}
}

// UnionInPlace adds all elements from the other set to this set.
func (ch *ConcurrentHash[T]) UnionInPlace(other Set[T]) {
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
// Returns a new thread-safe ConcurrentHash without modifying the original.
func (ch *ConcurrentHash[T]) Remove(element T) Set[T] {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	result := NewConcurrentHash[T]()
	for e := range ch.data {
		if e != element {
			result.data[e] = struct{}{}
		}
	}
	return result
}

// RemoveMany creates a new set with all given elements removed.
// Returns a new thread-safe ConcurrentHash without modifying the original.
func (ch *ConcurrentHash[T]) RemoveMany(elements ...T) Set[T] {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	// Create a set of elements to remove for efficient lookup
	toRemove := make(map[T]struct{}, len(elements))
	for _, element := range elements {
		toRemove[element] = struct{}{}
	}

	result := NewConcurrentHash[T]()
	for e := range ch.data {
		if _, shouldRemove := toRemove[e]; !shouldRemove {
			result.data[e] = struct{}{}
		}
	}
	return result
}

// Difference creates a new set containing elements in this set but not in the other set.
// Returns a new thread-safe ConcurrentHash without modifying the original.
//
// Membership in other is decided by other.Contains, exactly as the
// non-concurrent sets do, so an operand whose equality is not == (such as a
// TreeSet holding NaN) is honoured. The receiver's elements are snapshotted
// under its lock, which is released before other is consulted, so the two
// sets' locks are never held together: a.Difference(b) racing b.Difference(a) cannot
// deadlock, and the result is not an atomic view of both sets.
func (ch *ConcurrentHash[T]) Difference(other Set[T]) Set[T] {
	return ch.Filter(func(element T) bool { return !other.Contains(element) })
}

// RemoveInPlace removes the given element from the set.
// Returns true if the element was present and removed; false otherwise.
func (ch *ConcurrentHash[T]) RemoveInPlace(element T) bool {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	if _, exists := ch.data[element]; exists {
		delete(ch.data, element)
		return true
	}
	return false
}

// RemoveManyInPlace removes all given elements from the set.
func (ch *ConcurrentHash[T]) RemoveManyInPlace(elements ...T) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for _, element := range elements {
		delete(ch.data, element)
	}
}

// DifferenceInPlace removes all elements that are present in the other set.
func (ch *ConcurrentHash[T]) DifferenceInPlace(other Set[T]) {
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
func (ch *ConcurrentHash[T]) Clear() {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	for element := range ch.data {
		delete(ch.data, element)
	}
}

// Intersection creates a new set containing elements present in both sets.
// Returns a new thread-safe ConcurrentHash without modifying the original.
//
// Membership in other is decided by other.Contains, exactly as the
// non-concurrent sets do, so an operand whose equality is not == (such as a
// TreeSet holding NaN) is honoured. The receiver's elements are snapshotted
// under its lock, which is released before other is consulted, so the two
// sets' locks are never held together: a.Intersection(b) racing b.Intersection(a) cannot
// deadlock, and the result is not an atomic view of both sets.
func (ch *ConcurrentHash[T]) Intersection(other Set[T]) Set[T] {
	return ch.Filter(other.Contains)
}

// IsSubsetOf returns true if this set is a subset of the other set.
//
// Membership in other is decided by other.Contains, exactly as the
// non-concurrent sets do, so an operand whose equality is not == (such as a
// TreeSet holding NaN) is honoured. The receiver's elements are snapshotted
// under its lock, which is released before other is consulted, so the two
// sets' locks are never held together: a.IsSubsetOf(b) racing b.IsSubsetOf(a) cannot
// deadlock, and the result is not an atomic view of both sets.
func (ch *ConcurrentHash[T]) IsSubsetOf(other Set[T]) bool {
	return ch.AllMatch(other.Contains)
}

// IsSupersetOf returns true if this set is a superset of the other set.
func (ch *ConcurrentHash[T]) IsSupersetOf(other Set[T]) bool {
	return other.IsSubsetOf(ch)
}

// IsDisjoint returns true if this set has no elements in common with the other set.
//
// Membership in other is decided by other.Contains, exactly as the
// non-concurrent sets do, so an operand whose equality is not == (such as a
// TreeSet holding NaN) is honoured. The receiver's elements are snapshotted
// under its lock, which is released before other is consulted, so the two
// sets' locks are never held together: a.IsDisjoint(b) racing b.IsDisjoint(a) cannot
// deadlock, and the result is not an atomic view of both sets.
func (ch *ConcurrentHash[T]) IsDisjoint(other Set[T]) bool {
	return ch.NoneMatch(other.Contains)
}

// Equals returns true if both sets contain exactly the same elements.
//
// Membership in other is decided by other.Contains, exactly as the
// non-concurrent sets do, so an operand whose equality is not == (such as a
// TreeSet holding NaN) is honoured. The receiver's elements are snapshotted
// under its lock, which is released before other is consulted, so the two
// sets' locks are never held together: a.Equals(b) racing b.Equals(a) cannot
// deadlock, and the result is not an atomic view of both sets.
func (ch *ConcurrentHash[T]) Equals(other Set[T]) bool {
	return sameMembers(ch.AsSlice(), other)
}

// IntersectionInPlace keeps only elements that are present in both sets. It is
// FilterInPlace with other.Contains as the predicate, so other is consulted
// outside the receiver's lock.
//
// Membership in other is decided by other.Contains, exactly as the
// non-concurrent sets do, so an operand whose equality is not == (such as a
// TreeSet holding NaN) is honoured. The receiver's elements are snapshotted
// under its lock, which is released before other is consulted, so the two
// sets' locks are never held together: a.IntersectionInPlace(b) racing b.IntersectionInPlace(a) cannot
// deadlock, and the result is not an atomic view of both sets.
func (ch *ConcurrentHash[T]) IntersectionInPlace(other Set[T]) {
	ch.FilterInPlace(other.Contains)
}

// sameMembers reports whether other holds exactly the snapshotted elements: it
// has the same number of elements and other.Contains reports every one of them.
// It is Equals for the concurrent sets, which snapshot their own elements and
// then consult other outside their lock.
func sameMembers[T comparable](elements []T, other Set[T]) bool {
	if len(elements) != other.Length() {
		return false
	}
	for _, element := range elements {
		if !other.Contains(element) {
			return false
		}
	}
	return true
}
