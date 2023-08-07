package collections

import (
	"github.com/pickeringtech/go-collectionutil/collections/dicts"
	"github.com/pickeringtech/go-collectionutil/collections/lists"
)

func NewList[T any](values ...T) lists.Array[T] {
	return lists.NewArray(values...)
}

func NewConcurrentArray[T any](values ...T) lists.ConcurrentArray[T] {
	return lists.NewConcurrentArray[T](values...)
}

func NewMap[K comparable, V any]() dicts.Hash[K, V] {
	return dicts.NewHash[K, V]()
}

func NewSet[T comparable]() dicts.Hash[T, struct{}] {
	return dicts.NewHash[T, struct{}]()
}
