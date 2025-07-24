package lists

import "sync"

// ConcurrentDoublyLinked is a thread-safe doubly linked list implementation
// using a mutex for synchronization. All operations are protected by a single mutex.
type ConcurrentDoublyLinked[T any] struct {
	data *DoublyLinked[T]
	lock *sync.Mutex
}

// NewConcurrentDoublyLinked creates a new ConcurrentDoublyLinked with the given values.
func NewConcurrentDoublyLinked[T any](values ...T) *ConcurrentDoublyLinked[T] {
	return &ConcurrentDoublyLinked[T]{
		data: NewDoublyLinked(values...),
		lock: &sync.Mutex{},
	}
}

// NewConcurrentDoublyLinkedCircular creates a new circular ConcurrentDoublyLinked with the given values.
func NewConcurrentDoublyLinkedCircular[T any](values ...T) *ConcurrentDoublyLinked[T] {
	return &ConcurrentDoublyLinked[T]{
		data: NewDoublyLinkedCircular(values...),
		lock: &sync.Mutex{},
	}
}

// Interface guards to ensure ConcurrentDoublyLinked implements the required interfaces
var _ List[int] = &ConcurrentDoublyLinked[int]{}
var _ MutableList[int] = &ConcurrentDoublyLinked[int]{}

// AllMatch returns true if all elements satisfy the given predicate.
func (cl *ConcurrentDoublyLinked[T]) AllMatch(fn func(T) bool) bool {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.AllMatch(fn)
}

// AnyMatch returns true if any element satisfies the given predicate.
func (cl *ConcurrentDoublyLinked[T]) AnyMatch(fn func(T) bool) bool {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.AnyMatch(fn)
}

// Find returns the first element that satisfies the given predicate.
func (cl *ConcurrentDoublyLinked[T]) Find(fn func(T) bool) (T, bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.Find(fn)
}

// FindIndex returns the index of the first element that satisfies the given predicate.
func (cl *ConcurrentDoublyLinked[T]) FindIndex(fn func(T) bool) int {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.FindIndex(fn)
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func (cl *ConcurrentDoublyLinked[T]) Filter(fn func(T) bool) []T {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.Filter(fn)
}

// FilterInPlace removes elements that don't satisfy the predicate.
func (cl *ConcurrentDoublyLinked[T]) FilterInPlace(fn func(T) bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.FilterInPlace(fn)
}

// Get returns the element at the given index, or defaultValue if out of bounds.
func (cl *ConcurrentDoublyLinked[T]) Get(index int, defaultValue T) T {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.Get(index, defaultValue)
}

// Length returns the number of elements in the list.
func (cl *ConcurrentDoublyLinked[T]) Length() int {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.Length()
}

// ForEach executes the given function for each element.
func (cl *ConcurrentDoublyLinked[T]) ForEach(fn EachFunc[T]) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.ForEach(fn)
}

// ForEachWithIndex executes the given function for each element with its index.
func (cl *ConcurrentDoublyLinked[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.ForEachWithIndex(fn)
}

// GetAsSlice returns the list as a slice.
func (cl *ConcurrentDoublyLinked[T]) GetAsSlice() []T {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.GetAsSlice()
}

// Insert creates a new slice with elements inserted at the given index.
func (cl *ConcurrentDoublyLinked[T]) Insert(index int, elements ...T) []T {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.Insert(index, elements...)
}

// InsertInPlace inserts elements at the given index.
func (cl *ConcurrentDoublyLinked[T]) InsertInPlace(index int, elements ...T) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.InsertInPlace(index, elements...)
}

// Sort returns a new sorted slice.
func (cl *ConcurrentDoublyLinked[T]) Sort(lessThan func(T, T) bool) []T {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.Sort(lessThan)
}

// SortInPlace sorts the list in place.
func (cl *ConcurrentDoublyLinked[T]) SortInPlace(lessThan func(T, T) bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.SortInPlace(lessThan)
}

// Push adds an element to the end and returns a new slice.
func (cl *ConcurrentDoublyLinked[T]) Push(element T) []T {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.Push(element)
}

// PushInPlace adds an element to the end.
func (cl *ConcurrentDoublyLinked[T]) PushInPlace(element T) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.PushInPlace(element)
}

// Pop removes and returns the last element.
func (cl *ConcurrentDoublyLinked[T]) Pop() (T, bool, []T) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.Pop()
}

// PopInPlace removes and returns the last element.
func (cl *ConcurrentDoublyLinked[T]) PopInPlace() (T, bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.PopInPlace()
}

// PeekEnd returns the last element without removing it.
func (cl *ConcurrentDoublyLinked[T]) PeekEnd() (T, bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.PeekEnd()
}

// Enqueue adds an element to the end and returns a new slice.
func (cl *ConcurrentDoublyLinked[T]) Enqueue(element T) []T {
	return cl.Push(element)
}

// EnqueueInPlace adds an element to the end.
func (cl *ConcurrentDoublyLinked[T]) EnqueueInPlace(element T) {
	cl.PushInPlace(element)
}

// Dequeue removes and returns the first element.
func (cl *ConcurrentDoublyLinked[T]) Dequeue() (T, bool, []T) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.Dequeue()
}

// DequeueInPlace removes and returns the first element.
func (cl *ConcurrentDoublyLinked[T]) DequeueInPlace() (T, bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.DequeueInPlace()
}

// PeekFront returns the first element without removing it.
func (cl *ConcurrentDoublyLinked[T]) PeekFront() (T, bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.PeekFront()
}
