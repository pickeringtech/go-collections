package lists

import (
	"github.com/pickeringtech/go-collections/slices"
	"sync"
)

// ConcurrentRWArray is a slice-backed implementation of MutableList that is
// safe for concurrent use. It uses a sync.RWMutex so that read-only operations
// can proceed concurrently while mutating operations take an exclusive lock.
type ConcurrentRWArray[T any] struct {
	elements []T
	lock     *sync.RWMutex
}

// FilterInPlace retains only the elements for which fn returns true, modifying
// the receiver under an exclusive lock. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) FilterInPlace(fn func(T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Filter(a.elements, fn)
}

// InsertInPlace inserts the given elements at index, modifying the receiver
// under an exclusive lock. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) InsertInPlace(index int, element ...T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Insert(a.elements, index, element...)
}

// PushInPlace appends element to the end of the receiver under an exclusive
// lock. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) PushInPlace(element T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Push(a.elements, element)
}

// PopInPlace removes and returns the last element, reporting whether one was
// present, modifying the receiver under an exclusive lock. It is safe for
// concurrent use.
func (a *ConcurrentRWArray[T]) PopInPlace() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	res, ok, newSli := slices.Pop(a.elements)
	a.elements = newSli
	return res, ok
}

// EnqueueInPlace appends element to the end of the receiver under an exclusive
// lock. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) EnqueueInPlace(element T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Push(a.elements, element)
}

// DequeueInPlace removes and returns the first element, reporting whether one
// was present, modifying the receiver under an exclusive lock. It is safe for
// concurrent use.
func (a *ConcurrentRWArray[T]) DequeueInPlace() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	res, ok, newSli := slices.PopFront(a.elements)
	a.elements = newSli
	return res, ok
}

// NewConcurrentRWArray creates a new ConcurrentRWArray seeded with the given
// elements, preserving their order.
func NewConcurrentRWArray[T any](elements ...T) *ConcurrentRWArray[T] {
	return &ConcurrentRWArray[T]{
		elements: elements,
		lock:     &sync.RWMutex{},
	}
}

// Interface guards
var _ Filterable[int] = &ConcurrentRWArray[int]{}
var _ Indexable[int] = &ConcurrentRWArray[int]{}
var _ Iterable[int] = &ConcurrentRWArray[int]{}
var _ Searchable[int] = &ConcurrentRWArray[int]{}
var _ Sortable[int] = &ConcurrentRWArray[int]{}
var _ List[int] = &ConcurrentRWArray[int]{}
var _ MutableList[int] = &ConcurrentRWArray[int]{}

// AllMatch returns true if every element satisfies the predicate fun (vacuously
// true for an empty list). It takes a read lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) AllMatch(fun func(T) bool) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.AllMatch(a.elements, fun)
}

// AnyMatch returns true if at least one element satisfies the predicate fun. It
// takes a read lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) AnyMatch(fun func(T) bool) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.AnyMatch(a.elements, fun)
}

// NoneMatch returns true if no element satisfies the predicate fun (vacuously
// true for an empty list). It takes a read lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) NoneMatch(fun func(T) bool) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return !slices.AnyMatch(a.elements, fun)
}

// Dequeue returns the first element, whether one was present, and a new List
// (independent of the receiver's backing array) with that element removed,
// without modifying the receiver. It takes a read lock and is safe for
// concurrent use.
func (a *ConcurrentRWArray[T]) Dequeue() (T, bool, List[T]) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	// Operate on a copy so the returned List is independent of the receiver's
	// backing array (PopFront returns a sub-slice of its input).
	res, ok, rest := slices.PopFront(slices.Copy(a.elements))
	return res, ok, NewArray(rest...)
}

// Enqueue returns a new List (independent of the receiver's backing array)
// with element appended to the end, without modifying the receiver. It takes a
// read lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Enqueue(element T) List[T] {
	a.lock.RLock()
	defer a.lock.RUnlock()

	// Operate on a copy so the returned List is independent of the receiver's
	// backing array (Push may append into shared capacity).
	return NewArray(slices.Push(slices.Copy(a.elements), element)...)
}

// Filter returns a new List containing only the elements for which fun returns
// true, without modifying the receiver. It takes a read lock and is safe for
// concurrent use.
func (a *ConcurrentRWArray[T]) Filter(fun func(T) bool) List[T] {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return NewArray(slices.Filter(a.elements, fun)...)
}

// Find returns the first element for which fun returns true and whether such an
// element was found. It takes a read lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Find(fun func(T) bool) (T, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Find(a.elements, fun)
}

// FindIndex returns the index of the first element for which fun returns true,
// or -1 if none match. It takes a read lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) FindIndex(fun func(T) bool) int {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.FindIndex(a.elements, fun)
}

// ForEach calls fun once for each element in order while holding a read lock. It
// is safe for concurrent use.
func (a *ConcurrentRWArray[T]) ForEach(fun EachFunc[T]) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	for _, element := range a.elements {
		fun(element)
	}
}

// ForEachWithIndex calls fun once for each element in order, passing the
// element's index and value, while holding a read lock. It is safe for
// concurrent use.
func (a *ConcurrentRWArray[T]) ForEachWithIndex(fun IndexedEachFunc[T]) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	for idx, element := range a.elements {
		fun(idx, element)
	}
}

// Get returns the element at index and true, or defaultValue and false if the
// index is out of bounds. It takes a read lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Get(index int, defaultValue T) (T, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if index < 0 || index >= len(a.elements) {
		return defaultValue, false
	}
	return a.elements[index], true
}

// AsSlice returns a copy of the elements as a new slice, independent of the
// receiver's backing array. It takes a read lock and is safe for concurrent
// use.
func (a *ConcurrentRWArray[T]) AsSlice() []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Copy(a.elements)
}

// Insert returns a new List (independent of the receiver) with the given
// elements inserted at index, without modifying the receiver. It takes a read
// lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Insert(index int, element ...T) List[T] {
	a.lock.RLock()
	defer a.lock.RUnlock()

	// Operate on a copy so the returned List is independent of the receiver and
	// the receiver's backing array is never mutated by the insert.
	return NewArray(slices.Insert(slices.Copy(a.elements), index, element...)...)
}

// Length returns the number of elements in the list. It takes a read lock and is
// safe for concurrent use.
func (a *ConcurrentRWArray[T]) Length() int {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Length(a.elements)
}

// IsEmpty returns true if the list contains no elements. It takes a read lock
// and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) IsEmpty() bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Length(a.elements) == 0
}

// RemoveAt returns a new List (independent of the receiver's backing array)
// with the element at index removed, without modifying the receiver. If index
// is out of bounds the elements are returned unchanged. It takes a read lock and
// is safe for concurrent use.
func (a *ConcurrentRWArray[T]) RemoveAt(index int) List[T] {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return NewArray(deleteOwned(slices.Copy(a.elements), index)...)
}

// Remove returns a new List (independent of the receiver's backing array) with
// the first element deeply equal to element removed, without modifying the
// receiver. If no element matches, the elements are returned unchanged. It takes
// a read lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Remove(element T) List[T] {
	a.lock.RLock()
	defer a.lock.RUnlock()

	elements := slices.Copy(a.elements)
	return NewArray(deleteOwned(elements, indexOfDeepEqual(elements, element))...)
}

// RemoveAtInPlace removes the element at index, returning it and whether the
// index was in bounds, modifying the receiver under an exclusive lock. It is
// safe for concurrent use.
func (a *ConcurrentRWArray[T]) RemoveAtInPlace(index int) (T, bool) {
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
// whether an element was removed, modifying the receiver under an exclusive
// lock. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) RemoveInPlace(element T) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	index := indexOfDeepEqual(a.elements, element)
	if index < 0 {
		return false
	}
	a.elements = slices.Delete(a.elements, index)
	return true
}

// Clear removes all elements from the list under an exclusive lock. It is safe
// for concurrent use.
func (a *ConcurrentRWArray[T]) Clear() {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = nil
}

// PeekEnd returns the last element without removing it, and whether one was
// present. It takes a read lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) PeekEnd() (T, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.PeekEnd(a.elements)
}

// PeekFront returns the first element without removing it, and whether one was
// present. It takes a read lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) PeekFront() (T, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.PeekFront(a.elements)
}

// Pop returns the last element, whether one was present, and a new List
// (independent of the receiver's backing array) with that element removed,
// without modifying the receiver. It takes a read lock and is safe for
// concurrent use.
func (a *ConcurrentRWArray[T]) Pop() (T, bool, List[T]) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	// Operate on a copy so the returned List is independent of the receiver's
	// backing array (Pop returns a sub-slice of its input).
	res, ok, rest := slices.Pop(slices.Copy(a.elements))
	return res, ok, NewArray(rest...)
}

// Push returns a new List (independent of the receiver's backing array) with
// element appended to the end, without modifying the receiver. It takes a read
// lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Push(element T) List[T] {
	a.lock.RLock()
	defer a.lock.RUnlock()

	// Operate on a copy so the returned List is independent of the receiver's
	// backing array (Push may append into shared capacity).
	return NewArray(slices.Push(slices.Copy(a.elements), element)...)
}

// Sort returns a new List sorted according to the less-than function lessThan,
// without modifying the receiver. It takes a read lock and is safe for
// concurrent use.
func (a *ConcurrentRWArray[T]) Sort(lessThan func(T, T) bool) List[T] {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return NewArray(slices.Sort(a.elements, lessThan)...)
}

// SortInPlace sorts the receiver's elements according to the less-than function
// lessThan, under an exclusive lock. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) SortInPlace(lessThan func(T, T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	slices.SortInPlace(a.elements, lessThan)
}
