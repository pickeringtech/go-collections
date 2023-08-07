package lists

type Filterable[T any] interface {
	Filter(fun func(T) bool) []T
	FilterInPlace(fun func(T) bool)
}

type Indexable[T any] interface {
	Get(index int, defaultValue T) T
	Length() int
}

type Mappable[I, O any] interface {
	Map(fun func(I) O) []O
}

type Searchable[T any] interface {
	AllMatch(fun func(T) bool) bool
	AnyMatch(fun func(T) bool) bool
	Find(fun func(T) bool) (T, bool)
	FindIndex(fun func(T) bool) int
}

type Sortable[T any] interface {
	Sort(fun func(T, T) bool) []T
	SortInPlace(fun func(T, T) bool)
}

type Reducible[I, O any] interface {
	Reduce(fun func(O, I) O, initial O) O
}
