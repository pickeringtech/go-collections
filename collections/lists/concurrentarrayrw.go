package lists

import (
	"github.com/pickeringtech/go-collectionutil/slices"
	"sync"
)

type ConcurrentArrayRW[T any] struct {
	elements []T
	lock     *sync.RWMutex
}

func NewConcurrentArrayRW[T any](elements ...T) ConcurrentArrayRW[T] {
	return ConcurrentArrayRW[T]{
		elements: elements,
		lock:     &sync.RWMutex{},
	}
}

// Interface guards
var _ Filterable[int] = &ConcurrentArrayRW[int]{}
var _ Indexable[int] = &ConcurrentArrayRW[int]{}
var _ Searchable[int] = &ConcurrentArrayRW[int]{}
var _ Sortable[int] = &ConcurrentArrayRW[int]{}

func (a ConcurrentArrayRW[T]) AllMatch(fun func(T) bool) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.AllMatch(a.elements, fun)
}

func (a ConcurrentArrayRW[T]) AnyMatch(fun func(T) bool) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.AnyMatch(a.elements, fun)
}

func (a ConcurrentArrayRW[T]) Filter(fun func(T) bool) []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Filter(a.elements, fun)
}

func (a ConcurrentArrayRW[T]) FilterInPlace(fun func(T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	slices.FilterInPlace(a.elements, fun)
}

func (a ConcurrentArrayRW[T]) Find(fun func(T) bool) (T, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Find(a.elements, fun)
}

func (a ConcurrentArrayRW[T]) FindIndex(fun func(T) bool) int {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.FindIndex(a.elements, fun)
}

func (a ConcurrentArrayRW[T]) ForEach(fun EachFunc[T]) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	for _, element := range a.elements {
		fun(element)
	}
}

func (a ConcurrentArrayRW[T]) ForEachWithIndex(fun EachFuncWithIndex[T]) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	for idx, element := range a.elements {
		fun(idx, element)
	}
}

func (a ConcurrentArrayRW[T]) Get(index int, defaultValue T) T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Get(a.elements, index, defaultValue)
}

func (a ConcurrentArrayRW[T]) Length() int {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Length(a.elements)
}

func (a ConcurrentArrayRW[T]) Sort(lessThan func(T, T) bool) []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Sort(a.elements, lessThan)
}

func (a ConcurrentArrayRW[T]) SortInPlace(lessThan func(T, T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	slices.SortInPlace(a.elements, lessThan)
}
