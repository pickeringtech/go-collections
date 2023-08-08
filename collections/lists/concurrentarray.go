package lists

import (
	"github.com/pickeringtech/go-collections/slices"
	"sync"
)

type ConcurrentArray[T any] struct {
	elements []T
	lock     *sync.Mutex
}

func NewConcurrentArray[T any](elements ...T) *ConcurrentArray[T] {
	return &ConcurrentArray[T]{
		elements: elements,
		lock:     &sync.Mutex{},
	}
}

// Interface guards
var _ List[int] = &ConcurrentArray[int]{}
var _ MutableList[int] = &ConcurrentArray[int]{}

func (a *ConcurrentArray[T]) AllMatch(fun func(T) bool) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.AllMatch(a.elements, fun)
}

func (a *ConcurrentArray[T]) AnyMatch(fun func(T) bool) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.AnyMatch(a.elements, fun)
}

func (a *ConcurrentArray[T]) Dequeue() (T, bool, []T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.PopFront(a.elements)
}

func (a *ConcurrentArray[T]) DequeueInPlace() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	res, ok, newSli := slices.PopFront(a.elements)
	a.elements = newSli
	return res, ok
}

func (a *ConcurrentArray[T]) Enqueue(element T) []T {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Push(a.elements, element)
}

func (a *ConcurrentArray[T]) EnqueueInPlace(element T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Push(a.elements, element)
}
func (a *ConcurrentArray[T]) Filter(fun func(T) bool) []T {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Filter(a.elements, fun)
}

func (a *ConcurrentArray[T]) FilterInPlace(fn func(T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Filter(a.elements, fn)
}

func (a *ConcurrentArray[T]) Find(fun func(T) bool) (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Find(a.elements, fun)
}

func (a *ConcurrentArray[T]) FindIndex(fun func(T) bool) int {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.FindIndex(a.elements, fun)
}

func (a *ConcurrentArray[T]) ForEach(fun EachFunc[T]) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for _, element := range a.elements {
		fun(element)
	}
}

func (a *ConcurrentArray[T]) ForEachWithIndex(fun IndexedEachFunc[T]) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for idx, element := range a.elements {
		fun(idx, element)
	}
}

func (a *ConcurrentArray[T]) Get(index int, defaultValue T) T {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Get(a.elements, index, defaultValue)
}

func (a *ConcurrentArray[T]) GetAsSlice() []T {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Copy(a.elements)
}

func (a *ConcurrentArray[T]) Insert(index int, element ...T) []T {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Insert(a.elements, index, element...)
}

func (a *ConcurrentArray[T]) InsertInPlace(index int, element ...T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Insert(a.elements, index, element...)
}

func (a *ConcurrentArray[T]) Length() int {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Length(a.elements)
}

func (a *ConcurrentArray[T]) PeekEnd() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.PeekEnd(a.elements)
}

func (a *ConcurrentArray[T]) PeekFront() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.PeekFront(a.elements)
}

func (a *ConcurrentArray[T]) Pop() (T, bool, []T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Pop(a.elements)
}

func (a *ConcurrentArray[T]) PopInPlace() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	res, ok, newSli := slices.Pop(a.elements)
	a.elements = newSli
	return res, ok
}

func (a *ConcurrentArray[T]) Push(element T) []T {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Push(a.elements, element)
}

func (a *ConcurrentArray[T]) PushInPlace(element T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.elements = slices.Push(a.elements, element)
}

func (a *ConcurrentArray[T]) Sort(lessThan func(T, T) bool) []T {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Sort(a.elements, lessThan)
}

func (a *ConcurrentArray[T]) SortInPlace(lessThan func(T, T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	slices.SortInPlace(a.elements, lessThan)
}
