package lists

type EachFunc[T any] func(element T)

type IndexedEachFunc[T any] func(idx int, element T)
