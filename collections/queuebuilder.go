package collections

import "github.com/pickeringtech/go-collections/collections/lists"

type QueueBuilder[T any] struct {
	isConcurrent bool
	isRW         bool
	values       []T
}

func NewQueueBuilder[T any]() QueueBuilder[T] {
	return QueueBuilder[T]{}
}

func (b QueueBuilder[T]) Concurrent() QueueBuilder[T] {
	b.isConcurrent = true
	return b
}

func (b QueueBuilder[T]) RW() QueueBuilder[T] {
	b.isRW = true
	return b
}

func (b QueueBuilder[T]) Add(values ...T) QueueBuilder[T] {
	b.values = append(b.values, values...)
	return b
}

func (b QueueBuilder[T]) Build() lists.Queue[T] {
	//if b.isRW {
	//	return NewConcurrentRWList[T](b.values...)
	//}
	if b.isConcurrent {
		return NewConcurrentList[T](b.values...)
	}
	return NewList[T](b.values...)
}
