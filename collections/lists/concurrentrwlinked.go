package lists

import (
	"sync"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
	"github.com/pickeringtech/go-collections/slices"
)

// ConcurrentRWLinked is a thread-safe linked list implementation
// using a read-write mutex for synchronization. Read operations use read locks for better
// performance when there are many concurrent readers.
//
// ConcurrentRWLinked must not be copied after first use; copying after construction
// produces an independent lock over shared backing data, which breaks the
// thread-safety contract. go vet reports any such copy.
type ConcurrentRWLinked[T any] struct {
	_    nocopy.NoCopy
	data *Linked[T]
	lock sync.RWMutex
}

// NewConcurrentRWLinked creates a new ConcurrentRWLinked with the given values.
func NewConcurrentRWLinked[T any](values ...T) *ConcurrentRWLinked[T] {
	return &ConcurrentRWLinked[T]{
		data: NewLinked(values...),
	}
}

// NewConcurrentRWLinkedCircular creates a new circular ConcurrentRWLinked with the given values.
func NewConcurrentRWLinkedCircular[T any](values ...T) *ConcurrentRWLinked[T] {
	return &ConcurrentRWLinked[T]{
		data: NewLinkedCircular(values...),
	}
}

// Interface guards to ensure ConcurrentRWLinked implements the required interfaces
var _ List[int] = &ConcurrentRWLinked[int]{}
var _ MutableList[int] = &ConcurrentRWLinked[int]{}

// AllMatch returns true if all elements satisfy the given predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (cl *ConcurrentRWLinked[T]) AllMatch(fn func(T) bool) bool {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	return slices.AllMatch(snapshot, fn)
}

// AnyMatch returns true if any element satisfies the given predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (cl *ConcurrentRWLinked[T]) AnyMatch(fn func(T) bool) bool {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	return slices.AnyMatch(snapshot, fn)
}

// NoneMatch returns true if no element satisfies the given predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (cl *ConcurrentRWLinked[T]) NoneMatch(fn func(T) bool) bool {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	return !slices.AnyMatch(snapshot, fn)
}

// Find returns the first element that satisfies the given predicate. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the collection.
func (cl *ConcurrentRWLinked[T]) Find(fn func(T) bool) (T, bool) {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	return slices.Find(snapshot, fn)
}

// FindIndex returns the index of the first element that satisfies the given
// predicate. The predicate is evaluated after the lock is released, against a
// point-in-time snapshot taken under the lock, so it may safely call back into
// the collection.
func (cl *ConcurrentRWLinked[T]) FindIndex(fn func(T) bool) int {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	return slices.FindIndex(snapshot, fn)
}

// Filter returns a new List containing only elements that satisfy the
// predicate. The predicate is evaluated after the lock is released, against a
// point-in-time snapshot taken under the lock, so it may safely call back into
// the collection.
func (cl *ConcurrentRWLinked[T]) Filter(fn func(T) bool) List[T] {
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
func (cl *ConcurrentRWLinked[T]) FilterInPlace(fn func(T) bool) {
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
func (cl *ConcurrentRWLinked[T]) Get(index int, defaultValue T) (T, bool) {
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

// IsEmpty returns true if the list contains no elements.
func (cl *ConcurrentRWLinked[T]) IsEmpty() bool {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.IsEmpty()
}

// RemoveAt returns a new List with the element at index removed, without
// modifying the receiver.
func (cl *ConcurrentRWLinked[T]) RemoveAt(index int) List[T] {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.RemoveAt(index)
}

// Remove returns a new List with the first element deeply equal to element
// removed, without modifying the receiver.
func (cl *ConcurrentRWLinked[T]) Remove(element T) List[T] {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.Remove(element)
}

// RemoveAtInPlace removes the element at index, returning it and whether the
// index was in bounds.
func (cl *ConcurrentRWLinked[T]) RemoveAtInPlace(index int) (T, bool) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.RemoveAtInPlace(index)
}

// RemoveInPlace removes the first element deeply equal to element, reporting
// whether an element was removed.
func (cl *ConcurrentRWLinked[T]) RemoveInPlace(element T) bool {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.data.RemoveInPlace(element)
}

// Clear removes all elements from the list.
func (cl *ConcurrentRWLinked[T]) Clear() {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	cl.data.Clear()
}

// ForEach executes the given function for each element. fn is invoked after the
// lock is released, against a point-in-time snapshot taken under the lock, so fn
// may safely call back into the collection.
func (cl *ConcurrentRWLinked[T]) ForEach(fn EachFunc[T]) {
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
func (cl *ConcurrentRWLinked[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	cl.lock.RLock()
	snapshot := cl.data.AsSlice()
	cl.lock.RUnlock()
	for idx, element := range snapshot {
		fn(idx, element)
	}
}

// AsSlice returns the list as a slice.
func (cl *ConcurrentRWLinked[T]) AsSlice() []T {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	return cl.data.AsSlice()
}

// Insert creates a new List with elements inserted at the given index.
func (cl *ConcurrentRWLinked[T]) Insert(index int, elements ...T) List[T] {
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

// Sort returns a new sorted List.
func (cl *ConcurrentRWLinked[T]) Sort(lessThan func(T, T) bool) List[T] {
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

// Push adds an element to the end and returns a new List.
func (cl *ConcurrentRWLinked[T]) Push(element T) List[T] {
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
func (cl *ConcurrentRWLinked[T]) Pop() (T, bool, List[T]) {
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

// Enqueue adds an element to the end and returns a new List.
func (cl *ConcurrentRWLinked[T]) Enqueue(element T) List[T] {
	return cl.Push(element)
}

// EnqueueInPlace adds an element to the end.
func (cl *ConcurrentRWLinked[T]) EnqueueInPlace(element T) {
	cl.PushInPlace(element)
}

// Dequeue removes and returns the first element.
func (cl *ConcurrentRWLinked[T]) Dequeue() (T, bool, List[T]) {
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
