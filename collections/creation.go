package collections

import (
	"github.com/pickeringtech/go-collections/collections/dicts"
	"github.com/pickeringtech/go-collections/collections/lists"
	"github.com/pickeringtech/go-collections/collections/sets"
)

func NewList[T any](values ...T) lists.List[T] {
	return lists.NewArray(values...)
}

func NewConcurrentList[T any](values ...T) lists.List[T] {
	return lists.NewConcurrentArray[T](values...)
}

//func NewConcurrentRWList[T any](values ...T) lists.List[T] {
//	return lists.NewConcurrentArrayRW[T](values...)
//}

func NewQueue[T any](values ...T) lists.Queue[T] {
	return lists.NewArray(values...)
}

//func NewConcurrentRWQueue[T any](values ...T) lists.List[T] {
//	return lists.NewConcurrentArrayRW[T](values...)
//}

func NewConcurrentQueue[T any](values ...T) lists.Queue[T] {
	return lists.NewConcurrentArray[T](values...)
}

func NewStack[T any](values ...T) lists.Stack[T] {
	return lists.NewArray(values...)
}

func NewConcurrentStack[T any](values ...T) lists.Stack[T] {
	return lists.NewConcurrentArray[T](values...)
}

//func NewConcurrentRWStack[T any](values ...T) lists.List[T] {
//	return lists.NewConcurrentArrayRW[T](values...)
//}

func NewDict[K comparable, V any](entries ...dicts.Pair[K, V]) dicts.Dict[K, V] {
	return dicts.NewHash[K, V](entries...)
}

func NewConcurrentDict[K comparable, V any](entries ...dicts.Pair[K, V]) dicts.Dict[K, V] {
	return dicts.NewConcurrentHash[K, V](entries...)
}

func NewConcurrentRWDict[K comparable, V any](entries ...dicts.Pair[K, V]) dicts.Dict[K, V] {
	return dicts.NewConcurrentHashRW[K, V](entries...)
}

func NewSet[T comparable](elements ...T) sets.Set[T] {
	return sets.NewHash[T](elements...)
}

func NewConcurrentSet[T comparable](elements ...T) sets.Set[T] {
	return sets.NewConcurrentHash[T](elements...)
}

func NewConcurrentRWSet[T comparable](elements ...T) sets.Set[T] {
	return sets.NewConcurrentHashRW[T](elements...)
}

func NewLinkedList[T any](elements ...T) lists.List[T] {
	return lists.NewLinked[T](elements...)
}

func NewConcurrentLinkedList[T any](elements ...T) lists.List[T] {
	return lists.NewConcurrentLinked[T](elements...)
}

func NewConcurrentRWLinkedList[T any](elements ...T) lists.List[T] {
	return lists.NewConcurrentRWLinked[T](elements...)
}

func NewDoublyLinkedList[T any](elements ...T) lists.List[T] {
	return lists.NewDoublyLinked[T](elements...)
}

func NewConcurrentDoublyLinkedList[T any](elements ...T) lists.List[T] {
	return lists.NewConcurrentDoublyLinked[T](elements...)
}

func NewConcurrentRWDoublyLinkedList[T any](elements ...T) lists.List[T] {
	return lists.NewConcurrentRWDoublyLinked[T](elements...)
}


