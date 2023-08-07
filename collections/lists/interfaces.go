package lists

type Filterable[T any] interface {
	Filter(fn func(T) bool) []T
}

type Indexable[T any] interface {
	Get(index int, defaultValue T) T
	Length() int
}

type Insertable[T any] interface {
	Insert(index int, element T) []T
	InsertAll(index int, elements []T) []T
	InsertInPlace(index int, element T)
	InsertAllInPlace(index int, elements []T)
}

type Iterable[T any] interface {
	ForEach(fn EachFunc[T])
	ForEachWithIndex(fn IndexedEachFunc[T])
}

type Mappable[I, O any] interface {
	Map(fn func(I) O) []O
}

type Searchable[T any] interface {
	AllMatch(fn func(T) bool) bool
	AnyMatch(fn func(T) bool) bool
	Find(fn func(T) bool) (T, bool)
	FindIndex(fn func(T) bool) int
}

type Sortable[T any] interface {
	Sort(fn func(T, T) bool) []T
	SortInPlace(fn func(T, T) bool)
}

type Reducible[I, O any] interface {
	Reduce(fn func(O, I) O, initial O) O
}
