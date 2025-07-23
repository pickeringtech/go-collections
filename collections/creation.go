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

func NewSet[T comparable]() sets.Set[T] {
	return dicts.NewHash[T, struct{}]()
}

func NewConcurrentSet[T comparable]() sets.Set[T] {
	return dicts.NewHash[T, struct{}]()
}

func NewConcurrentRWSet[T comparable]() sets.Set[T] {
	return dicts.NewHash[T, struct{}]()
}
