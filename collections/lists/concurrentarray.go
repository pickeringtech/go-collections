package lists

import (
	"github.com/pickeringtech/go-collectionutil/slices"
	"sync"
)

type ConcurrentArray[T any] struct {
	elements []T
	lock     *sync.Mutex
}

func NewConcurrentArray[T any](elements ...T) ConcurrentArray[T] {
	return ConcurrentArray[T]{
		elements: elements,
		lock:     &sync.Mutex{},
	}
}

// Interface guards
var _ Filterable[int] = &ConcurrentArray[int]{}
var _ Indexable[int] = &ConcurrentArray[int]{}
var _ Searchable[int] = &ConcurrentArray[int]{}
var _ Sortable[int] = &ConcurrentArray[int]{}

func (a ConcurrentArray[T]) AllMatch(fun func(T) bool) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.AllMatch(a.elements, fun)
}

func (a ConcurrentArray[T]) AnyMatch(fun func(T) bool) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.AnyMatch(a.elements, fun)
}

func (a ConcurrentArray[T]) Filter(fun func(T) bool) []T {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Filter(a.elements, fun)
}

func (a ConcurrentArray[T]) FilterInPlace(fun func(T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	slices.FilterInPlace(a.elements, fun)
}

func (a ConcurrentArray[T]) Find(fun func(T) bool) (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Find(a.elements, fun)
}

func (a ConcurrentArray[T]) FindIndex(fun func(T) bool) int {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.FindIndex(a.elements, fun)
}

func (a ConcurrentArray[T]) ForEach(fun EachFunc[T]) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for _, element := range a.elements {
		fun(element)
	}
}

func (a ConcurrentArray[T]) ForEachWithIndex(fun EachFuncWithIndex[T]) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for idx, element := range a.elements {
		fun(idx, element)
	}
}

func (a ConcurrentArray[T]) Get(index int, defaultValue T) T {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Get(a.elements, index, defaultValue)
}

func (a ConcurrentArray[T]) Length() int {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Length(a.elements)
}

func (a ConcurrentArray[T]) Sort(lessThan func(T, T) bool) []T {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Sort(a.elements, lessThan)
}

func (a ConcurrentArray[T]) SortInPlace(lessThan func(T, T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	slices.SortInPlace(a.elements, lessThan)
}
