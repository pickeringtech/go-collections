package lists

import (
	"github.com/pickeringtech/go-collections/slices"
	"sync"
)

// ConcurrentArray is a slice-backed implementation of MutableList that is safe
// for concurrent use. Every operation is guarded by a sync.Mutex.
type ConcurrentArray[T any] struct {
	elements []T
	lock     *sync.Mutex
}

// NewConcurrentArray creates a new ConcurrentArray seeded with the given
// elements, preserving their order.
func NewConcurrentArray[T any](elements ...T) *ConcurrentArray[T] {
	return &ConcurrentArray[T]{
		elements: elements,
		lock:     &sync.Mutex{},
	}
}

// Interface guards
var _ List[int] = &ConcurrentArray[int]{}
var _ MutableList[int] = &ConcurrentArray[int]{}

// AllMatch returns true if every element satisfies the predicate fun (vacuously
// true for an empty list). The predicate is evaluated after the lock is
// released, against a point-in-time snapshot taken under the lock, so it may
// safely call back into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) AllMatch(fun func(T) bool) bool {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	return slices.AllMatch(snapshot, fun)
}

// AnyMatch returns true if at least one element satisfies the predicate fun. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the
// collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) AnyMatch(fun func(T) bool) bool {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	return slices.AnyMatch(snapshot, fun)
}

// NoneMatch returns true if no element satisfies the predicate fun (vacuously
// true for an empty list). The predicate is evaluated after the lock is
// released, against a point-in-time snapshot taken under the lock, so it may
// safely call back into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) NoneMatch(fun func(T) bool) bool {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	return !slices.AnyMatch(snapshot, fun)
}

// Dequeue returns the first element, whether one was present, and a new List
// (independent of the receiver's backing array) with that element removed,
// without modifying the receiver. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Dequeue() (T, bool, List[T]) {
	a.lock.Lock()
	defer a.lock.Unlock()

	// Operate on a copy so the returned List is independent of the receiver's
	// backing array (PopFront returns a sub-slice of its input).
	res, ok, rest := slices.PopFront(slices.Copy(a.elements))
	return res, ok, NewArray(rest...)
}

// DequeueInPlace removes and returns the first element, reporting whether one
// was present, modifying the receiver. It is safe for concurrent use.
func (a *ConcurrentArray[T]) DequeueInPlace() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	res, ok, newSli := slices.PopFront(a.elements)
	a.elements = newSli
	return res, ok
}

// Enqueue returns a new List (independent of the receiver's backing array)
// with element appended to the end, without modifying the receiver. It is safe
// for concurrent use.
func (a *ConcurrentArray[T]) Enqueue(element T) List[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	// PushCopy copies and appends in a single allocation, so the new List is
	// independent of the receiver (plain Push could append into shared capacity).
	return NewArray(slices.PushCopy(a.elements, element)...)
}

// EnqueueInPlace appends element to the end of the receiver. It is safe for
// concurrent use.
func (a *ConcurrentArray[T]) EnqueueInPlace(element T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Push(a.elements, element)
}

// Filter returns a new List containing only the elements for which fun returns
// true, without modifying the receiver. The predicate is evaluated after the
// lock is released, against a point-in-time snapshot taken under the lock, so it
// may safely call back into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Filter(fun func(T) bool) List[T] {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	return NewArray(slices.Filter(snapshot, fun)...)
}

// FilterInPlace retains only the elements for which fn returns true, modifying
// the receiver. The predicate is evaluated after the lock is released, against
// a point-in-time snapshot taken under the lock, so it may safely call back into
// the collection. Modifications made concurrently with evaluation are not
// reflected in the retained set. It is safe for concurrent use.
func (a *ConcurrentArray[T]) FilterInPlace(fn func(T) bool) {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	retained := slices.Filter(snapshot, fn)

	a.lock.Lock()
	a.elements = retained
	a.lock.Unlock()
}

// Find returns the first element for which fun returns true and whether such an
// element was found. The predicate is evaluated after the lock is released,
// against a point-in-time snapshot taken under the lock, so it may safely call
// back into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Find(fun func(T) bool) (T, bool) {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	return slices.Find(snapshot, fun)
}

// FindIndex returns the index of the first element for which fun returns true,
// or -1 if none match. The predicate is evaluated after the lock is released,
// against a point-in-time snapshot taken under the lock, so it may safely call
// back into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) FindIndex(fun func(T) bool) int {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	return slices.FindIndex(snapshot, fun)
}

// ForEach calls fun once for each element in order. fun is invoked after the
// lock is released, against a point-in-time snapshot taken under the lock, so
// fun may safely call back into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) ForEach(fun EachFunc[T]) {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	for _, element := range snapshot {
		fun(element)
	}
}

// ForEachWithIndex calls fun once for each element in order, passing the
// element's index and value. fun is invoked after the lock is released, against
// a point-in-time snapshot taken under the lock, so fun may safely call back
// into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) ForEachWithIndex(fun IndexedEachFunc[T]) {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	for idx, element := range snapshot {
		fun(idx, element)
	}
}

// Get returns the element at index and true, or defaultValue and false if the
// index is out of bounds. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Get(index int, defaultValue T) (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if index < 0 || index >= len(a.elements) {
		return defaultValue, false
	}
	return a.elements[index], true
}

// AsSlice returns a copy of the elements as a new slice, independent of the
// receiver's backing array. It is safe for concurrent use.
func (a *ConcurrentArray[T]) AsSlice() []T {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Copy(a.elements)
}

// Insert returns a new List (independent of the receiver) with the given
// elements inserted at index, without modifying the receiver. It is safe for
// concurrent use.
func (a *ConcurrentArray[T]) Insert(index int, element ...T) List[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	// Operate on a copy so the returned List is independent of the receiver and
	// the receiver's backing array is never mutated by the insert.
	return NewArray(slices.Insert(slices.Copy(a.elements), index, element...)...)
}

// InsertInPlace inserts the given elements at index, modifying the receiver. It
// is safe for concurrent use.
func (a *ConcurrentArray[T]) InsertInPlace(index int, element ...T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Insert(a.elements, index, element...)
}

// Length returns the number of elements in the list. It is safe for concurrent
// use.
func (a *ConcurrentArray[T]) Length() int {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Length(a.elements)
}

// IsEmpty returns true if the list contains no elements. It is safe for
// concurrent use.
func (a *ConcurrentArray[T]) IsEmpty() bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Length(a.elements) == 0
}

// RemoveAt returns a new List (independent of the receiver's backing array)
// with the element at index removed, without modifying the receiver. If index
// is out of bounds the elements are returned unchanged. It is safe for
// concurrent use.
func (a *ConcurrentArray[T]) RemoveAt(index int) List[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	return NewArray(deleteOwned(slices.Copy(a.elements), index)...)
}

// Remove returns a new List (independent of the receiver's backing array) with
// the first element deeply equal to element removed, without modifying the
// receiver. If no element matches, the elements are returned unchanged. It is
// safe for concurrent use.
func (a *ConcurrentArray[T]) Remove(element T) List[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	elements := slices.Copy(a.elements)
	return NewArray(deleteOwned(elements, indexOfDeepEqual(elements, element))...)
}

// RemoveAtInPlace removes the element at index, returning it and whether the
// index was in bounds, modifying the receiver. It is safe for concurrent use.
func (a *ConcurrentArray[T]) RemoveAtInPlace(index int) (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if index < 0 || index >= len(a.elements) {
		var zero T
		return zero, false
	}
	removed := a.elements[index]
	a.elements = slices.Delete(a.elements, index)
	return removed, true
}

// RemoveInPlace removes the first element deeply equal to element, reporting
// whether an element was removed, modifying the receiver. It is safe for
// concurrent use.
func (a *ConcurrentArray[T]) RemoveInPlace(element T) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	index := indexOfDeepEqual(a.elements, element)
	if index < 0 {
		return false
	}
	a.elements = slices.Delete(a.elements, index)
	return true
}

// Clear removes all elements from the list. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Clear() {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = nil
}

// PeekEnd returns the last element without removing it, and whether one was
// present. It is safe for concurrent use.
func (a *ConcurrentArray[T]) PeekEnd() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.PeekEnd(a.elements)
}

// PeekFront returns the first element without removing it, and whether one was
// present. It is safe for concurrent use.
func (a *ConcurrentArray[T]) PeekFront() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.PeekFront(a.elements)
}

// Pop returns the last element, whether one was present, and a new List
// (independent of the receiver's backing array) with that element removed,
// without modifying the receiver. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Pop() (T, bool, List[T]) {
	a.lock.Lock()
	defer a.lock.Unlock()

	// Operate on a copy so the returned List is independent of the receiver's
	// backing array (Pop returns a sub-slice of its input).
	res, ok, rest := slices.Pop(slices.Copy(a.elements))
	return res, ok, NewArray(rest...)
}

// PopInPlace removes and returns the last element, reporting whether one was
// present, modifying the receiver. It is safe for concurrent use.
func (a *ConcurrentArray[T]) PopInPlace() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	res, ok, newSli := slices.Pop(a.elements)
	a.elements = newSli
	return res, ok
}

// Push returns a new List (independent of the receiver's backing array) with
// element appended to the end, without modifying the receiver. It is safe for
// concurrent use.
func (a *ConcurrentArray[T]) Push(element T) List[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	// PushCopy copies and appends in a single allocation, so the new List is
	// independent of the receiver (plain Push could append into shared capacity).
	return NewArray(slices.PushCopy(a.elements, element)...)
}

// PushInPlace appends element to the end of the receiver. It is safe for
// concurrent use.
func (a *ConcurrentArray[T]) PushInPlace(element T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Push(a.elements, element)
}

// Sort returns a new List sorted according to the less-than function lessThan,
// without modifying the receiver. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Sort(lessThan func(T, T) bool) List[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	return NewArray(slices.Sort(a.elements, lessThan)...)
}

// SortInPlace sorts the receiver's elements according to the less-than function
// lessThan. It is safe for concurrent use.
func (a *ConcurrentArray[T]) SortInPlace(lessThan func(T, T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	slices.SortInPlace(a.elements, lessThan)
}
