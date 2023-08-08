package collections

import "github.com/pickeringtech/go-collections/collections/lists"

type StackBuilder[T any] struct {
	isConcurrent bool
	isRW         bool
	values       []T
}

func NewStackBuilder[T any]() StackBuilder[T] {
	return StackBuilder[T]{}
}

func (b StackBuilder[T]) Concurrent() StackBuilder[T] {
	b.isConcurrent = true
	return b
}

func (b StackBuilder[T]) RW() StackBuilder[T] {
	b.isRW = true
	return b
}

func (b StackBuilder[T]) Add(values ...T) StackBuilder[T] {
	b.values = append(b.values, values...)
	return b
}

func (b StackBuilder[T]) Build() lists.Stack[T] {
	//if b.isRW {
	//	return NewConcurrentRWList[T](b.values...)
	//}
	if b.isConcurrent {
		return NewConcurrentList[T](b.values...)
	}
	return NewList[T](b.values...)
}
