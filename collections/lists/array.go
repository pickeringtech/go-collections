package lists

import "github.com/pickeringtech/go-collections/slices"

// Array is a slice-backed implementation of MutableList. It delegates each
// operation to the helpers in the slices package and is not safe for concurrent
// use.
type Array[T any] struct {
	elements []T
}

// NewArray creates a new Array seeded with the given elements, preserving their
// order.
func NewArray[T any](elements ...T) *Array[T] {
	return &Array[T]{
		elements: elements,
	}
}

// Interface guards
var _ List[int] = &Array[int]{}
var _ MutableList[int] = &Array[int]{}

// AllMatch returns true if every element satisfies the predicate fn (vacuously
// true for an empty list).
func (a *Array[T]) AllMatch(fn func(T) bool) bool {
	return slices.AllMatch(a.elements, fn)
}

// AnyMatch returns true if at least one element satisfies the predicate fn.
func (a *Array[T]) AnyMatch(fn func(T) bool) bool {
	return slices.AnyMatch(a.elements, fn)
}

// NoneMatch returns true if no element satisfies the predicate fn (vacuously
// true for an empty list).
func (a *Array[T]) NoneMatch(fn func(T) bool) bool {
	return !slices.AnyMatch(a.elements, fn)
}

// Dequeue returns the first element, whether one was present, and a new List
// with that element removed, without modifying the receiver.
func (a *Array[T]) Dequeue() (T, bool, List[T]) {
	res, ok, rest := slices.PopFront(a.elements)
	return res, ok, NewArray(rest...)
}

// DequeueInPlace removes and returns the first element, reporting whether one
// was present, modifying the receiver.
func (a *Array[T]) DequeueInPlace() (T, bool) {
	res, ok, newSli := slices.PopFront(a.elements)
	a.elements = newSli
	return res, ok
}

// Enqueue returns a new List with element appended to the end, without
// modifying the receiver.
func (a *Array[T]) Enqueue(element T) List[T] {
	return NewArray(slices.Push(a.elements, element)...)
}

// EnqueueInPlace appends element to the end of the receiver.
func (a *Array[T]) EnqueueInPlace(element T) {
	a.elements = slices.Push(a.elements, element)
}

// Filter returns a new List containing only the elements for which fn returns
// true, without modifying the receiver.
func (a *Array[T]) Filter(fn func(T) bool) List[T] {
	return NewArray(slices.Filter(a.elements, fn)...)
}

// FilterInPlace retains only the elements for which fn returns true, modifying
// the receiver.
func (a *Array[T]) FilterInPlace(fn func(T) bool) {
	a.elements = slices.Filter(a.elements, fn)
}

// Find returns the first element for which fn returns true and whether such an
// element was found.
func (a *Array[T]) Find(fn func(T) bool) (T, bool) {
	return slices.Find(a.elements, fn)
}

// FindIndex returns the index of the first element for which fn returns true, or
// -1 if none match.
func (a *Array[T]) FindIndex(fn func(T) bool) int {
	return slices.FindIndex(a.elements, fn)
}

// ForEach calls fn once for each element in order.
func (a *Array[T]) ForEach(fn EachFunc[T]) {
	for _, element := range a.elements {
		fn(element)
	}
}

// ForEachWithIndex calls fn once for each element in order, passing the
// element's index and value.
func (a *Array[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	for idx, element := range a.elements {
		fn(idx, element)
	}
}

// Get returns the element at index and true, or defaultValue and false if the
// index is out of bounds.
func (a *Array[T]) Get(index int, defaultValue T) (T, bool) {
	if index < 0 || index >= len(a.elements) {
		return defaultValue, false
	}
	return a.elements[index], true
}

// AsSlice returns the underlying backing slice.
func (a *Array[T]) AsSlice() []T {
	return a.elements
}

// Insert returns a new List with the given elements inserted at index, without
// modifying the receiver. The index may range over 0 <= index <= Length(): an
// index equal to the length appends. An out-of-range index returns the elements
// unchanged.
func (a *Array[T]) Insert(index int, element ...T) List[T] {
	return NewArray(slices.Insert(a.elements, index, element...)...)
}

// InsertInPlace inserts the given elements at index, modifying the receiver. The
// index may range over 0 <= index <= Length(): an index equal to the length
// appends. An out-of-range index leaves the receiver untouched.
func (a *Array[T]) InsertInPlace(index int, element ...T) {
	a.elements = slices.Insert(a.elements, index, element...)
}

// Length returns the number of elements in the list.
func (a *Array[T]) Length() int {
	return slices.Length(a.elements)
}

// IsEmpty returns true if the list contains no elements.
func (a *Array[T]) IsEmpty() bool {
	return a.Length() == 0
}

// RemoveAt returns a new, independent List with the element at index removed,
// without modifying the receiver. If index is out of bounds the elements are
// returned unchanged.
func (a *Array[T]) RemoveAt(index int) List[T] {
	return NewArray(deleteOwned(slices.Copy(a.elements), index)...)
}

// Remove returns a new, independent List with the first element deeply equal to
// element removed, without modifying the receiver. If no element matches, the
// elements are returned unchanged.
func (a *Array[T]) Remove(element T) List[T] {
	elements := slices.Copy(a.elements)
	return NewArray(deleteOwned(elements, indexOfDeepEqual(elements, element))...)
}

// RemoveAtInPlace removes the element at index, returning it and whether the
// index was in bounds, modifying the receiver.
func (a *Array[T]) RemoveAtInPlace(index int) (T, bool) {
	if index < 0 || index >= len(a.elements) {
		var zero T
		return zero, false
	}
	removed := a.elements[index]
	a.elements = slices.Delete(a.elements, index)
	return removed, true
}

// RemoveInPlace removes the first element deeply equal to element, reporting
// whether an element was removed, modifying the receiver.
func (a *Array[T]) RemoveInPlace(element T) bool {
	index := indexOfDeepEqual(a.elements, element)
	if index < 0 {
		return false
	}
	a.elements = slices.Delete(a.elements, index)
	return true
}

// Clear removes all elements from the list.
func (a *Array[T]) Clear() {
	a.elements = nil
}

// PeekEnd returns the last element without removing it, and whether one was
// present.
func (a *Array[T]) PeekEnd() (T, bool) {
	return slices.PeekEnd(a.elements)
}

// PeekFront returns the first element without removing it, and whether one was
// present.
func (a *Array[T]) PeekFront() (T, bool) {
	return slices.PeekFront(a.elements)
}

// Pop returns the last element, whether one was present, and a new List with
// that element removed, without modifying the receiver.
func (a *Array[T]) Pop() (T, bool, List[T]) {
	res, ok, rest := slices.Pop(a.elements)
	return res, ok, NewArray(rest...)
}

// PopInPlace removes and returns the last element, reporting whether one was
// present, modifying the receiver.
func (a *Array[T]) PopInPlace() (T, bool) {
	res, ok, newSli := slices.Pop(a.elements)
	a.elements = newSli
	return res, ok
}

// Push returns a new List with element appended to the end, without modifying
// the receiver.
func (a *Array[T]) Push(element T) List[T] {
	return NewArray(slices.Push(a.elements, element)...)
}

// PushInPlace appends element to the end of the receiver.
func (a *Array[T]) PushInPlace(element T) {
	a.elements = slices.Push(a.elements, element)
}

// Sort returns a new List sorted according to the less-than function fn,
// without modifying the receiver.
func (a *Array[T]) Sort(fn func(T, T) bool) List[T] {
	return NewArray(slices.Sort(a.elements, fn)...)
}

// SortInPlace sorts the receiver's elements according to the less-than function
// fn.
func (a *Array[T]) SortInPlace(fn func(T, T) bool) {
	slices.SortInPlace(a.elements, fn)
}
