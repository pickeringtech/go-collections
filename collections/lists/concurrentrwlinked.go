package lists

import "sync"

// ConcurrentRWLinked is a thread-safe linked list implementation
// using a read-write mutex for synchronization. Read operations use read locks for better
// performance when there are many concurrent readers.
type ConcurrentRWLinked[T any] struct {
	data *Linked[T]
	lock *sync.RWMutex
}

// NewConcurrentRWLinked creates a new ConcurrentRWLinked with the given values.
func NewConcurrentRWLinked[T any](values ...T) *ConcurrentRWLinked[T] {
	return &ConcurrentRWLinked[T]{
		data: NewLinked(values...),
		lock: &sync.RWMutex{},
	}
}

// NewConcurrentRWLinkedCircular creates a new circular ConcurrentRWLinked with the given values.
func NewConcurrentRWLinkedCircular[T any](values ...T) *ConcurrentRWLinked[T] {
	return &ConcurrentRWLinked[T]{
		data: NewLinkedCircular(values...),
		lock: &sync.RWMutex{},
	}
}

// Interface guards to ensure ConcurrentRWLinked implements the required interfaces
var _ List[int] = &ConcurrentRWLinked[int]{}
var _ MutableList[int] = &ConcurrentRWLinked[int]{}

// AllMatch returns true if all elements satisfy the given predicate.
func (cl *ConcurrentRWLinked[T]) AllMatch(fn func(T) bool) bool {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.AllMatch(fn)
}

// AnyMatch returns true if any element satisfies the given predicate.
func (cl *ConcurrentRWLinked[T]) AnyMatch(fn func(T) bool) bool {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.AnyMatch(fn)
}

// Find returns the first element that satisfies the given predicate.
func (cl *ConcurrentRWLinked[T]) Find(fn func(T) bool) (T, bool) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Find(fn)
}

// FindIndex returns the index of the first element that satisfies the given predicate.
func (cl *ConcurrentRWLinked[T]) FindIndex(fn func(T) bool) int {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.FindIndex(fn)
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func (cl *ConcurrentRWLinked[T]) Filter(fn func(T) bool) []T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Filter(fn)
}

// FilterInPlace removes elements that don't satisfy the predicate.
func (cl *ConcurrentRWLinked[T]) FilterInPlace(fn func(T) bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.FilterInPlace(fn)
}

// Get returns the element at the given index, or defaultValue if out of bounds.
func (cl *ConcurrentRWLinked[T]) Get(index int, defaultValue T) T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Get(index, defaultValue)
}

// Length returns the number of elements in the list.
func (cl *ConcurrentRWLinked[T]) Length() int {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Length()
}

// ForEach executes the given function for each element.
func (cl *ConcurrentRWLinked[T]) ForEach(fn EachFunc[T]) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	cl.data.ForEach(fn)
}

// ForEachWithIndex executes the given function for each element with its index.
func (cl *ConcurrentRWLinked[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	cl.data.ForEachWithIndex(fn)
}

// GetAsSlice returns the list as a slice.
func (cl *ConcurrentRWLinked[T]) GetAsSlice() []T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.GetAsSlice()
}

// Insert creates a new slice with elements inserted at the given index.
func (cl *ConcurrentRWLinked[T]) Insert(index int, elements ...T) []T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Insert(index, elements...)
}

// InsertInPlace inserts elements at the given index.
func (cl *ConcurrentRWLinked[T]) InsertInPlace(index int, elements ...T) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.InsertInPlace(index, elements...)
}

// Sort returns a new sorted slice.
func (cl *ConcurrentRWLinked[T]) Sort(lessThan func(T, T) bool) []T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Sort(lessThan)
}

// SortInPlace sorts the list in place.
func (cl *ConcurrentRWLinked[T]) SortInPlace(lessThan func(T, T) bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.SortInPlace(lessThan)
}

// Push adds an element to the end and returns a new slice.
func (cl *ConcurrentRWLinked[T]) Push(element T) []T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Push(element)
}

// PushInPlace adds an element to the end.
func (cl *ConcurrentRWLinked[T]) PushInPlace(element T) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.PushInPlace(element)
}

// Pop removes and returns the last element.
func (cl *ConcurrentRWLinked[T]) Pop() (T, bool, []T) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Pop()
}

// PopInPlace removes and returns the last element.
func (cl *ConcurrentRWLinked[T]) PopInPlace() (T, bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.PopInPlace()
}

// PeekEnd returns the last element without removing it.
func (cl *ConcurrentRWLinked[T]) PeekEnd() (T, bool) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.PeekEnd()
}

// Enqueue adds an element to the end and returns a new slice.
func (cl *ConcurrentRWLinked[T]) Enqueue(element T) []T {
	return cl.Push(element)
}

// EnqueueInPlace adds an element to the end.
func (cl *ConcurrentRWLinked[T]) EnqueueInPlace(element T) {
	cl.PushInPlace(element)
}

// Dequeue removes and returns the first element.
func (cl *ConcurrentRWLinked[T]) Dequeue() (T, bool, []T) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Dequeue()
}

// DequeueInPlace removes and returns the first element.
func (cl *ConcurrentRWLinked[T]) DequeueInPlace() (T, bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.DequeueInPlace()
}

// PeekFront returns the first element without removing it.
func (cl *ConcurrentRWLinked[T]) PeekFront() (T, bool) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.PeekFront()
}
