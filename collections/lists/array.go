package lists

import "github.com/pickeringtech/go-collections/slices"

type Array[T any] struct {
	elements []T
}

func NewArray[T any](elements ...T) *Array[T] {
	return &Array[T]{
		elements: elements,
	}
}

// Interface guards
var _ List[int] = &Array[int]{}
var _ MutableList[int] = &Array[int]{}

func (a *Array[T]) AllMatch(fn func(T) bool) bool {
	return slices.AllMatch(a.elements, fn)
}

func (a *Array[T]) AnyMatch(fn func(T) bool) bool {
	return slices.AnyMatch(a.elements, fn)
}

func (a *Array[T]) Dequeue() (T, bool, []T) {
	return slices.PopFront(a.elements)
}

func (a *Array[T]) DequeueInPlace() (T, bool) {
	res, ok, newSli := slices.PopFront(a.elements)
	a.elements = newSli
	return res, ok
}

func (a *Array[T]) Enqueue(element T) []T {
	return slices.Push(a.elements, element)
}

func (a *Array[T]) EnqueueInPlace(element T) {
	a.elements = slices.Push(a.elements, element)
}

func (a *Array[T]) Filter(fn func(T) bool) []T {
	return slices.Filter(a.elements, fn)
}

func (a *Array[T]) FilterInPlace(fn func(T) bool) {
	a.elements = slices.Filter(a.elements, fn)
}

func (a *Array[T]) Find(fn func(T) bool) (T, bool) {
	return slices.Find(a.elements, fn)
}

func (a *Array[T]) FindIndex(fn func(T) bool) int {
	return slices.FindIndex(a.elements, fn)
}

func (a *Array[T]) ForEach(fn EachFunc[T]) {
	for _, element := range a.elements {
		fn(element)
	}
}

func (a *Array[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	for idx, element := range a.elements {
		fn(idx, element)
	}
}

func (a *Array[T]) Get(index int, defaultValue T) T {
	return slices.Get(a.elements, index, defaultValue)
}

func (a *Array[T]) GetAsSlice() []T {
	return a.elements
}

func (a *Array[T]) Insert(index int, element ...T) []T {
	return slices.Insert(a.elements, index, element...)
}

func (a *Array[T]) InsertInPlace(index int, element ...T) {
	a.elements = slices.Insert(a.elements, index, element...)
}

func (a *Array[T]) Length() int {
	return slices.Length(a.elements)
}

func (a *Array[T]) PeekEnd() (T, bool) {
	return slices.PeekEnd(a.elements)
}

func (a *Array[T]) PeekFront() (T, bool) {
	return slices.PeekFront(a.elements)
}

func (a *Array[T]) Pop() (T, bool, []T) {
	return slices.Pop(a.elements)
}

func (a *Array[T]) PopInPlace() (T, bool) {
	res, ok, newSli := slices.Pop(a.elements)
	a.elements = newSli
	return res, ok
}

func (a *Array[T]) Push(element T) []T {
	return slices.Push(a.elements, element)
}

func (a *Array[T]) PushInPlace(element T) {
	a.elements = slices.Push(a.elements, element)
}

func (a *Array[T]) Sort(fn func(T, T) bool) []T {
	return slices.Sort(a.elements, fn)
}

func (a *Array[T]) SortInPlace(fn func(T, T) bool) {
	slices.SortInPlace(a.elements, fn)
}
