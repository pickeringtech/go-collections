package lists

import (
	"github.com/pickeringtech/go-collections/slices"
	"sync"
)

type ConcurrentRWArray[T any] struct {
	elements []T
	lock     *sync.RWMutex
}

func (a *ConcurrentRWArray[T]) FilterInPlace(fn func(T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Filter(a.elements, fn)
}

func (a *ConcurrentRWArray[T]) InsertInPlace(index int, element ...T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Insert(a.elements, index, element...)
}

func (a *ConcurrentRWArray[T]) PushInPlace(element T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Push(a.elements, element)
}

func (a *ConcurrentRWArray[T]) PopInPlace() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	res, ok, newSli := slices.Pop(a.elements)
	a.elements = newSli
	return res, ok
}

func (a *ConcurrentRWArray[T]) EnqueueInPlace(element T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Push(a.elements, element)
}

func (a *ConcurrentRWArray[T]) DequeueInPlace() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	res, ok, newSli := slices.PopFront(a.elements)
	a.elements = newSli
	return res, ok
}

func NewConcurrentRWArray[T any](elements ...T) *ConcurrentRWArray[T] {
	return &ConcurrentRWArray[T]{
		elements: elements,
		lock:     &sync.RWMutex{},
	}
}

// Interface guards
var _ Filterable[int] = &ConcurrentRWArray[int]{}
var _ Indexable[int] = &ConcurrentRWArray[int]{}
var _ Iterable[int] = &ConcurrentRWArray[int]{}
var _ Searchable[int] = &ConcurrentRWArray[int]{}
var _ Sortable[int] = &ConcurrentRWArray[int]{}
var _ List[int] = &ConcurrentRWArray[int]{}
var _ MutableList[int] = &ConcurrentRWArray[int]{}

func (a *ConcurrentRWArray[T]) AllMatch(fun func(T) bool) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.AllMatch(a.elements, fun)
}

func (a *ConcurrentRWArray[T]) AnyMatch(fun func(T) bool) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.AnyMatch(a.elements, fun)
}

func (a *ConcurrentRWArray[T]) Dequeue() (T, bool, []T) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.PopFront(a.elements)
}

func (a *ConcurrentRWArray[T]) Enqueue(element T) []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Push(a.elements, element)
}

func (a *ConcurrentRWArray[T]) Filter(fun func(T) bool) []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Filter(a.elements, fun)
}

func (a *ConcurrentRWArray[T]) Find(fun func(T) bool) (T, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Find(a.elements, fun)
}

func (a *ConcurrentRWArray[T]) FindIndex(fun func(T) bool) int {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.FindIndex(a.elements, fun)
}

func (a *ConcurrentRWArray[T]) ForEach(fun EachFunc[T]) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	for _, element := range a.elements {
		fun(element)
	}
}

func (a *ConcurrentRWArray[T]) ForEachWithIndex(fun IndexedEachFunc[T]) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	for idx, element := range a.elements {
		fun(idx, element)
	}
}

func (a *ConcurrentRWArray[T]) Get(index int, defaultValue T) T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Get(a.elements, index, defaultValue)
}

func (a *ConcurrentRWArray[T]) GetAsSlice() []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Copy(a.elements)
}

func (a *ConcurrentRWArray[T]) Insert(index int, element ...T) []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Insert(a.elements, index, element...)
}

func (a *ConcurrentRWArray[T]) Length() int {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Length(a.elements)
}

func (a *ConcurrentRWArray[T]) PeekEnd() (T, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.PeekEnd(a.elements)
}

func (a *ConcurrentRWArray[T]) PeekFront() (T, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.PeekFront(a.elements)
}

func (a *ConcurrentRWArray[T]) Pop() (T, bool, []T) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Pop(a.elements)
}

func (a *ConcurrentRWArray[T]) Push(element T) []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Push(a.elements, element)
}

func (a *ConcurrentRWArray[T]) Sort(lessThan func(T, T) bool) []T {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return slices.Sort(a.elements, lessThan)
}

func (a *ConcurrentRWArray[T]) SortInPlace(lessThan func(T, T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	slices.SortInPlace(a.elements, lessThan)
}
