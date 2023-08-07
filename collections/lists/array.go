package lists

import "github.com/pickeringtech/go-collectionutil/slices"

type Array[T any] []T

func NewList[T any](elements ...T) Array[T] {
	return elements
}

// Interface guards
var _ Filterable[int] = &Array[int]{}
var _ Indexable[int] = &Array[int]{}
var _ Searchable[int] = &Array[int]{}
var _ Sortable[int] = &Array[int]{}

func (a Array[T]) AllMatch(fun func(T) bool) bool {
	return slices.AllMatch(a, fun)
}

func (a Array[T]) AnyMatch(fun func(T) bool) bool {
	return slices.AnyMatch(a, fun)
}

func (a Array[T]) Filter(fun func(T) bool) []T {
	return slices.Filter(a, fun)
}

func (a Array[T]) FilterInPlace(fun func(T) bool) {
	slices.FilterInPlace(a, fun)
}

func (a Array[T]) Find(fun func(T) bool) (T, bool) {
	return slices.Find(a, fun)
}

func (a Array[T]) FindIndex(fun func(T) bool) int {
	return slices.FindIndex(a, fun)
}

func (a Array[T]) ForEach(fun EachFunc[T]) {
	for _, element := range a {
		fun(element)
	}
}

func (a Array[T]) ForEachWithIndex(fun EachFuncWithIndex[T]) {
	for idx, element := range a {
		fun(idx, element)
	}
}

func (a Array[T]) Get(index int, defaultValue T) T {
	return slices.Get(a, index, defaultValue)
}

func (a Array[T]) Length() int {
	return slices.Length(a)
}

func (a Array[T]) Sort(fun func(T, T) bool) []T {
	return slices.Sort(a, fun)
}

func (a Array[T]) SortInPlace(fun func(T, T) bool) {
	slices.SortInPlace(a, fun)
}
