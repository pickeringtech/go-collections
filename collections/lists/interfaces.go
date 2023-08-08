package lists

type Filterable[T any] interface {
	Filter(fn func(T) bool) []T
}

type MutableFilterable[T any] interface {
	FilterInPlace(fn func(T) bool)
}

type Indexable[T any] interface {
	Get(index int, defaultValue T) T
	Length() int
}

type Insertable[T any] interface {
	Insert(index int, element ...T) []T
}

type MutableInsertable[T any] interface {
	InsertInPlace(index int, element ...T)
}

type Iterable[T any] interface {
	ForEach(fn EachFunc[T])
	ForEachWithIndex(fn IndexedEachFunc[T])
}

type List[T any] interface {
	Filterable[T]
	Indexable[T]
	Insertable[T]
	Iterable[T]
	Searchable[T]
	Sortable[T]
	Stack[T]
	Queue[T]
	GetAsSlice() []T
}

type MutableList[T any] interface {
	List[T]
	MutableFilterable[T]
	MutableInsertable[T]
	MutableSortable[T]
	MutableStack[T]
	MutableQueue[T]
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

type MutableSortable[T any] interface {
	SortInPlace(fn func(T, T) bool)
}

type Stack[T any] interface {
	Push(element T) []T
	Pop() (T, bool, []T)
	PeekEnd() (T, bool)
}

type MutableStack[T any] interface {
	PushInPlace(element T)
	PopInPlace() (T, bool)
}

type Queue[T any] interface {
	Enqueue(element T) []T
	Dequeue() (T, bool, []T)
	PeekFront() (T, bool)
}

type MutableQueue[T any] interface {
	EnqueueInPlace(element T)
	DequeueInPlace() (T, bool)
}
