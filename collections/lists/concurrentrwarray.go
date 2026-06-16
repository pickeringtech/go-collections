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

// Dequeue returns the first element, whether one was present, and a new slice
// (independent of the receiver's backing array) with that element removed,
// without modifying the receiver. It takes a read lock and is safe for
// concurrent use.
func (a *ConcurrentRWArray[T]) Dequeue() (T, bool, []T) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	// Operate on a copy so the returned slice is independent of the receiver's
	// backing array (PopFront returns a sub-slice of its input).
	return slices.PopFront(slices.Copy(a.elements))
}

// Enqueue returns a new slice (independent of the receiver's backing array)
// with element appended to the end, without modifying the receiver. It takes a
// read lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Enqueue(element T) []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	// Operate on a copy so the returned slice is independent of the receiver's
	// backing array (Push may append into shared capacity).
	return slices.Push(slices.Copy(a.elements), element)
}

// Filter returns a new slice containing only the elements for which fun returns
// true, without modifying the receiver. It takes a read lock and is safe for
// concurrent use.
func (a *ConcurrentRWArray[T]) Filter(fun func(T) bool) []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Filter(a.elements, fun)
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

// Get returns the element at index, or defaultValue if the index is out of
// bounds. It takes a read lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Get(index int, defaultValue T) T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Get(a.elements, index, defaultValue)
}

// GetAsSlice returns a copy of the elements as a new slice, independent of the
// receiver's backing array. It takes a read lock and is safe for concurrent
// use.
func (a *ConcurrentRWArray[T]) GetAsSlice() []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Copy(a.elements)
}

// Insert returns a new slice (independent of the receiver) with the given
// elements inserted at index, without modifying the receiver. It takes a read
// lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Insert(index int, element ...T) []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	// Operate on a copy so the returned slice is independent of the receiver and
	// the receiver's backing array is never mutated by the insert.
	return slices.Insert(slices.Copy(a.elements), index, element...)
}

// Length returns the number of elements in the list. It takes a read lock and is
// safe for concurrent use.
func (a *ConcurrentRWArray[T]) Length() int {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Length(a.elements)
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

// Pop returns the last element, whether one was present, and a new slice
// (independent of the receiver's backing array) with that element removed,
// without modifying the receiver. It takes a read lock and is safe for
// concurrent use.
func (a *ConcurrentRWArray[T]) Pop() (T, bool, []T) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	// Operate on a copy so the returned slice is independent of the receiver's
	// backing array (Pop returns a sub-slice of its input).
	return slices.Pop(slices.Copy(a.elements))
}

// Push returns a new slice (independent of the receiver's backing array) with
// element appended to the end, without modifying the receiver. It takes a read
// lock and is safe for concurrent use.
func (a *ConcurrentRWArray[T]) Push(element T) []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	// Operate on a copy so the returned slice is independent of the receiver's
	// backing array (Push may append into shared capacity).
	return slices.Push(slices.Copy(a.elements), element)
}

// Sort returns a new slice sorted according to the less-than function lessThan,
// without modifying the receiver. It takes a read lock and is safe for
// concurrent use.
func (a *ConcurrentRWArray[T]) Sort(lessThan func(T, T) bool) []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Sort(a.elements, lessThan)
}

// SortInPlace sorts the receiver's elements according to the less-than function
// lessThan, under an exclusive lock. It is safe for concurrent use.
func (a *ConcurrentRWArray[T]) SortInPlace(lessThan func(T, T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	slices.SortInPlace(a.elements, lessThan)
}
