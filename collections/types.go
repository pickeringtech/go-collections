package collections

type EachFunc[T any] func(element T)

type EachFuncWithIndex[T any] func(idx int, element T)
