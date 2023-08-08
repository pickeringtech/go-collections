package collections

import "github.com/pickeringtech/go-collections/collections/lists"

type ListBuilder[T any] struct {
	isConcurrent bool
	isRW         bool
	values       []T
}

func NewListBuilder[T any]() ListBuilder[T] {
	return ListBuilder[T]{}
}

func (b ListBuilder[T]) Concurrent() ListBuilder[T] {
	b.isConcurrent = true
	return b
}

func (b ListBuilder[T]) RW() ListBuilder[T] {
	b.isRW = true
	return b
}

func (b ListBuilder[T]) Add(values ...T) ListBuilder[T] {
	b.values = append(b.values, values...)
	return b
}

func (b ListBuilder[T]) Build() lists.List[T] {
	//if b.isRW {
	//	return NewConcurrentRWList[T](b.values...)
	//}
	if b.isConcurrent {
		return NewConcurrentList[T](b.values...)
	}
	return NewList[T](b.values...)
}
