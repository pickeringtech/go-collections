package lists

import (
	"sync"

	"github.com/pickeringtech/go-collections/slices"
)

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

// AllMatch returns true if all elements satisfy the given predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (cl *ConcurrentRWDoublyLinked[T]) AllMatch(fn func(T) bool) bool {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	return slices.AllMatch(snapshot, fn)
}

// AnyMatch returns true if any element satisfies the given predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (cl *ConcurrentRWDoublyLinked[T]) AnyMatch(fn func(T) bool) bool {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	return slices.AnyMatch(snapshot, fn)
}

// NoneMatch returns true if no element satisfies the given predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (cl *ConcurrentRWDoublyLinked[T]) NoneMatch(fn func(T) bool) bool {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	return !slices.AnyMatch(snapshot, fn)
}

// Find returns the first element that satisfies the given predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (cl *ConcurrentRWDoublyLinked[T]) Find(fn func(T) bool) (T, bool) {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	return slices.Find(snapshot, fn)
}

// FindIndex returns the index of the first element that satisfies the given
// predicate. The predicate is evaluated after the lock is released, against a
// point-in-time snapshot taken under the lock, so it may safely call back into
// the collection.
func (cl *ConcurrentRWDoublyLinked[T]) FindIndex(fn func(T) bool) int {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	return slices.FindIndex(snapshot, fn)
}

// Filter returns a new List containing only elements that satisfy the
// predicate. The predicate is evaluated after the lock is released, against a
// point-in-time snapshot taken under the lock, so it may safely call back into
// the collection.
func (cl *ConcurrentRWDoublyLinked[T]) Filter(fn func(T) bool) List[T] {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	return NewArray(slices.Filter(snapshot, fn)...)
}

// FilterInPlace removes elements that don't satisfy the predicate. The predicate
// is evaluated after the lock is released, against a point-in-time snapshot
// taken under the lock, so it may safely call back into the collection.
// Modifications made concurrently with evaluation are not reflected in the
// retained set.
func (cl *ConcurrentRWDoublyLinked[T]) FilterInPlace(fn func(T) bool) {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()

	retained := slices.Filter(snapshot, fn)

	cl.lock.Lock()
	cl.data.Clear()
	for _, element := range retained {
		cl.data.PushInPlace(element)
	}
	cl.lock.Unlock()
}

// Get returns the element at the given index and true, or defaultValue and
// false if the index is out of bounds.
func (cl *ConcurrentRWDoublyLinked[T]) Get(index int, defaultValue T) (T, bool) {
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

// IsEmpty returns true if the list contains no elements.
func (cl *ConcurrentRWDoublyLinked[T]) IsEmpty() bool {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.IsEmpty()
}

// RemoveAt returns a new List with the element at index removed, without
// modifying the receiver.
func (cl *ConcurrentRWDoublyLinked[T]) RemoveAt(index int) List[T] {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.RemoveAt(index)
}

// Remove returns a new List with the first element deeply equal to element
// removed, without modifying the receiver.
func (cl *ConcurrentRWDoublyLinked[T]) Remove(element T) List[T] {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Remove(element)
}

// RemoveAtInPlace removes the element at index, returning it and whether the
// index was in bounds.
func (cl *ConcurrentRWDoublyLinked[T]) RemoveAtInPlace(index int) (T, bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.RemoveAtInPlace(index)
}

// RemoveInPlace removes the first element deeply equal to element, reporting
// whether an element was removed.
func (cl *ConcurrentRWDoublyLinked[T]) RemoveInPlace(element T) bool {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.RemoveInPlace(element)
}

// Clear removes all elements from the list.
func (cl *ConcurrentRWDoublyLinked[T]) Clear() {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.Clear()
}

// ForEach executes the given function for each element. fn is invoked after the
// lock is released, against a point-in-time snapshot taken under the lock, so fn
// may safely call back into the collection.
func (cl *ConcurrentRWDoublyLinked[T]) ForEach(fn EachFunc[T]) {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	for _, element := range snapshot {
		fn(element)
	}
}

// ForEachWithIndex executes the given function for each element with its index.
// fn is invoked after the lock is released, against a point-in-time snapshot
// taken under the lock, so fn may safely call back into the collection.
func (cl *ConcurrentRWDoublyLinked[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	for idx, element := range snapshot {
		fn(idx, element)
	}
}

// AsSlice returns the list as a slice.
func (cl *ConcurrentRWDoublyLinked[T]) AsSlice() []T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.AsSlice()
}

// Insert creates a new List with elements inserted at the given index.
func (cl *ConcurrentRWDoublyLinked[T]) Insert(index int, elements ...T) List[T] {
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

// Sort returns a new sorted List.
func (cl *ConcurrentRWDoublyLinked[T]) Sort(lessThan func(T, T) bool) List[T] {
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

// Push adds an element to the end and returns a new List.
func (cl *ConcurrentRWDoublyLinked[T]) Push(element T) List[T] {
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
func (cl *ConcurrentRWDoublyLinked[T]) Pop() (T, bool, List[T]) {
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

// Enqueue adds an element to the end and returns a new List.
func (cl *ConcurrentRWDoublyLinked[T]) Enqueue(element T) List[T] {
	return cl.Push(element)
}

// EnqueueInPlace adds an element to the end.
func (cl *ConcurrentRWDoublyLinked[T]) EnqueueInPlace(element T) {
	cl.PushInPlace(element)
}

// Dequeue removes and returns the first element.
func (cl *ConcurrentRWDoublyLinked[T]) Dequeue() (T, bool, List[T]) {
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
