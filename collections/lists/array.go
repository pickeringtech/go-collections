package lists

import "github.com/pickeringtech/go-collectionutil/slices"

type Array[T any] []T

func NewArray[T any](elements ...T) Array[T] {
	return elements
}

// Interface guards
var _ Filterable[int] = &Array[int]{}
var _ Indexable[int] = &Array[int]{}
var _ Insertable[int] = &Array[int]{}
var _ Iterable[int] = &Array[int]{}
var _ Searchable[int] = &Array[int]{}
var _ Sortable[int] = &Array[int]{}

func (a Array[T]) AllMatch(fn func(T) bool) bool {
	return slices.AllMatch(a, fn)
}

func (a Array[T]) AnyMatch(fn func(T) bool) bool {
	return slices.AnyMatch(a, fn)
}

func (a Array[T]) Filter(fn func(T) bool) []T {
	return slices.Filter(a, fn)
}

func (a Array[T]) Find(fn func(T) bool) (T, bool) {
	return slices.Find(a, fn)
}

func (a Array[T]) FindIndex(fn func(T) bool) int {
	return slices.FindIndex(a, fn)
}

func (a Array[T]) ForEach(fn EachFunc[T]) {
	for _, element := range a {
		fn(element)
	}
}

func (a Array[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	for idx, element := range a {
		fn(idx, element)
	}
}

func (a Array[T]) Get(index int, defaultValue T) T {
	return slices.Get(a, index, defaultValue)
}

func (a Array[T]) Insert(index int, element T) []T {
	//TODO implement me
	panic("implement me")
}

func (a Array[T]) InsertAll(index int, elements []T) []T {
	//TODO implement me
	panic("implement me")
}

func (a Array[T]) InsertInPlace(index int, element T) {
	//TODO implement me
	panic("implement me")
}

func (a Array[T]) InsertAllInPlace(index int, elements []T) {
	//TODO implement me
	panic("implement me")
}

func (a Array[T]) Length() int {
	return slices.Length(a)
}

func (a Array[T]) Sort(fn func(T, T) bool) []T {
	return slices.Sort(a, fn)
}

func (a Array[T]) SortInPlace(fn func(T, T) bool) {
	slices.SortInPlace(a, fn)
}
