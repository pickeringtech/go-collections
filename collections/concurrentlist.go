package collections

import (
	"github.com/pickeringtech/go-collectionutil/slices"
	"sync"
)

type ConcurrentList[T any] struct {
	mu  sync.Mutex
	sli []T
}

func NewConcurrentList[T any]() *ConcurrentList[T] {
	return &ConcurrentList[T]{}
}

func (c *ConcurrentList[T]) Add(elements ...T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sli = append(c.sli, elements...)
}

func (c *ConcurrentList[T]) Delete(idx int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sli = slices.Delete(c.sli, idx)
}

func (c *ConcurrentList[T]) Get(idx int) (T, bool) {
	var defaultVal T
	if idx < 0 {
		return defaultVal, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	sliLen := len(c.sli)
	if idx >= sliLen {
		return defaultVal, false
	}
	return c.sli[idx], true
}

func (c *ConcurrentList[T]) Insert(idx int, elements ...T) {
}

func (c *ConcurrentList[T]) Pop() T {
	c.mu.Lock()
	defer c.mu.Unlock()
	lastElement, sli := slices.Pop(c.sli)
	c.sli = sli
	return lastElement
}

func (c *ConcurrentList[T]) PopFront() T {
	c.mu.Lock()
	defer c.mu.Unlock()
	firstElement, sli := slices.PopFront(c.sli)
	c.sli = sli
	return firstElement
}

func (c *ConcurrentList[T]) Push(elements ...T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sli = slices.Push(c.sli, elements...)
}

func (c *ConcurrentList[T]) PushFront(elements ...T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sli = slices.PushFront(c.sli, elements...)
}

func (c *ConcurrentList[T]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.sli)
}

func (c *ConcurrentList[T]) Cap() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return cap(c.sli)
}

func (c *ConcurrentList[T]) ToSlice() []T {
	c.mu.Lock()
	defer c.mu.Unlock()
	return slices.Copy(c.sli)
}

type StringifierFunc[T any] func(element T) string

func (c *ConcurrentList[T]) Stringify(fn StringifierFunc[T]) []string {
	return slices.Map[T, string](c.sli, func(t T) string {
		return fn(t)
	})
}
