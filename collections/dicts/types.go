package dicts

type EachFunc[T any] func(element T)

type EachFuncWithIndex[T any] func(idx int, element T)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}
