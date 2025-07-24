package lists

import "sync"

// ConcurrentRWDoublyLinked is a thread-safe doubly linked list implementation
// using a read-write mutex for synchronization. Read operations use read locks for better
// performance when there are many concurrent readers.
type ConcurrentRWDoublyLinked[T any] struct {
	data *DoublyLinked[T]
	lock *sync.RWMutex
}

// NewConcurrentRWDoublyLinked creates a new ConcurrentRWDoublyLinked with the given values.
func NewConcurrentRWDoublyLinked[T any](values ...T) *ConcurrentRWDoublyLinked[T] {
	return &ConcurrentRWDoublyLinked[T]{
		data: NewDoublyLinked(values...),
		lock: &sync.RWMutex{},
	}
}

// NewConcurrentRWDoublyLinkedCircular creates a new circular ConcurrentRWDoublyLinked with the given values.
func NewConcurrentRWDoublyLinkedCircular[T any](values ...T) *ConcurrentRWDoublyLinked[T] {
	return &ConcurrentRWDoublyLinked[T]{
		data: NewDoublyLinkedCircular(values...),
		lock: &sync.RWMutex{},
	}
}

// Interface guards to ensure ConcurrentRWDoublyLinked implements the required interfaces
var _ List[int] = &ConcurrentRWDoublyLinked[int]{}
var _ MutableList[int] = &ConcurrentRWDoublyLinked[int]{}

// AllMatch returns true if all elements satisfy the given predicate.
func (cl *ConcurrentRWDoublyLinked[T]) AllMatch(fn func(T) bool) bool {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.AllMatch(fn)
}

// AnyMatch returns true if any element satisfies the given predicate.
func (cl *ConcurrentRWDoublyLinked[T]) AnyMatch(fn func(T) bool) bool {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.AnyMatch(fn)
}

// Find returns the first element that satisfies the given predicate.
func (cl *ConcurrentRWDoublyLinked[T]) Find(fn func(T) bool) (T, bool) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Find(fn)
}

// FindIndex returns the index of the first element that satisfies the given predicate.
func (cl *ConcurrentRWDoublyLinked[T]) FindIndex(fn func(T) bool) int {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.FindIndex(fn)
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func (cl *ConcurrentRWDoublyLinked[T]) Filter(fn func(T) bool) []T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Filter(fn)
}

// FilterInPlace removes elements that don't satisfy the predicate.
func (cl *ConcurrentRWDoublyLinked[T]) FilterInPlace(fn func(T) bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.FilterInPlace(fn)
}

// Get returns the element at the given index, or defaultValue if out of bounds.
func (cl *ConcurrentRWDoublyLinked[T]) Get(index int, defaultValue T) T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Get(index, defaultValue)
}

// Length returns the number of elements in the list.
func (cl *ConcurrentRWDoublyLinked[T]) Length() int {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Length()
}

// ForEach executes the given function for each element.
func (cl *ConcurrentRWDoublyLinked[T]) ForEach(fn EachFunc[T]) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	cl.data.ForEach(fn)
}

// ForEachWithIndex executes the given function for each element with its index.
func (cl *ConcurrentRWDoublyLinked[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	cl.data.ForEachWithIndex(fn)
}

// GetAsSlice returns the list as a slice.
func (cl *ConcurrentRWDoublyLinked[T]) GetAsSlice() []T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.GetAsSlice()
}

// Insert creates a new slice with elements inserted at the given index.
func (cl *ConcurrentRWDoublyLinked[T]) Insert(index int, elements ...T) []T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Insert(index, elements...)
}

// InsertInPlace inserts elements at the given index.
func (cl *ConcurrentRWDoublyLinked[T]) InsertInPlace(index int, elements ...T) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.InsertInPlace(index, elements...)
}

// Sort returns a new sorted slice.
func (cl *ConcurrentRWDoublyLinked[T]) Sort(lessThan func(T, T) bool) []T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Sort(lessThan)
}

// SortInPlace sorts the list in place.
func (cl *ConcurrentRWDoublyLinked[T]) SortInPlace(lessThan func(T, T) bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.SortInPlace(lessThan)
}

// Push adds an element to the end and returns a new slice.
func (cl *ConcurrentRWDoublyLinked[T]) Push(element T) []T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Push(element)
}

// PushInPlace adds an element to the end.
func (cl *ConcurrentRWDoublyLinked[T]) PushInPlace(element T) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.PushInPlace(element)
}

// Pop removes and returns the last element.
func (cl *ConcurrentRWDoublyLinked[T]) Pop() (T, bool, []T) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Pop()
}

// PopInPlace removes and returns the last element.
func (cl *ConcurrentRWDoublyLinked[T]) PopInPlace() (T, bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.PopInPlace()
}

// PeekEnd returns the last element without removing it.
func (cl *ConcurrentRWDoublyLinked[T]) PeekEnd() (T, bool) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.PeekEnd()
}

// Enqueue adds an element to the end and returns a new slice.
func (cl *ConcurrentRWDoublyLinked[T]) Enqueue(element T) []T {
	return cl.Push(element)
}

// EnqueueInPlace adds an element to the end.
func (cl *ConcurrentRWDoublyLinked[T]) EnqueueInPlace(element T) {
	cl.PushInPlace(element)
}

// Dequeue removes and returns the first element.
func (cl *ConcurrentRWDoublyLinked[T]) Dequeue() (T, bool, []T) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Dequeue()
}

// DequeueInPlace removes and returns the first element.
func (cl *ConcurrentRWDoublyLinked[T]) DequeueInPlace() (T, bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.DequeueInPlace()
}

// PeekFront returns the first element without removing it.
func (cl *ConcurrentRWDoublyLinked[T]) PeekFront() (T, bool) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.PeekFront()
}
