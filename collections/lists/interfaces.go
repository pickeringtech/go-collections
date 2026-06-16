package lists

// Filterable is implemented by collections that can be filtered into a new
// slice without modifying the receiver.
type Filterable[T any] interface {
	Filter(fn func(T) bool) []T
}

// MutableFilterable is implemented by collections that can be filtered in place,
// modifying the receiver.
type MutableFilterable[T any] interface {
	FilterInPlace(fn func(T) bool)
}

// Indexable is implemented by collections that support positional access and
// reporting their length.
type Indexable[T any] interface {
	Get(index int, defaultValue T) (T, bool)
	Length() int
}

// Convertible is implemented by collections that can be converted into a slice.
// Whether the returned slice aliases the collection's backing storage or is an
// independent copy is implementation-defined.
type Convertible[T any] interface {
	AsSlice() []T
}

// Insertable is implemented by collections that can produce a new slice with
// elements inserted at a given index, without modifying the receiver.
type Insertable[T any] interface {
	Insert(index int, element ...T) []T
}

// MutableInsertable is implemented by collections that can insert elements at a
// given index in place, modifying the receiver.
type MutableInsertable[T any] interface {
	InsertInPlace(index int, element ...T)
}

// Iterable is implemented by collections that can be iterated over, with or
// without element indices.
type Iterable[T any] interface {
	ForEach(fn EachFunc[T])
	ForEachWithIndex(fn IndexedEachFunc[T])
}

// List is the read-oriented collection interface, combining filtering,
// indexing, insertion, iteration, searching, sorting, stack and queue
// behaviours along with conversion to a slice.
type List[T any] interface {
	Filterable[T]
	Indexable[T]
	Insertable[T]
	Iterable[T]
	Searchable[T]
	Sortable[T]
	Stack[T]
	Queue[T]
	Convertible[T]
}

// MutableList extends List with the in-place mutation operations for filtering,
// insertion, sorting, stack and queue behaviours.
type MutableList[T any] interface {
	List[T]
	MutableFilterable[T]
	MutableInsertable[T]
	MutableSortable[T]
	MutableStack[T]
	MutableQueue[T]
}

// Searchable is implemented by collections that can be searched with predicates.
type Searchable[T any] interface {
	AllMatch(fn func(T) bool) bool
	AnyMatch(fn func(T) bool) bool
	Find(fn func(T) bool) (T, bool)
	FindIndex(fn func(T) bool) int
}

// Sortable is implemented by collections that can be sorted, either into a new
// slice or in place, using a less-than comparison.
type Sortable[T any] interface {
	Sort(fn func(T, T) bool) []T
	SortInPlace(fn func(T, T) bool)
}

// MutableSortable is implemented by collections that can be sorted in place.
type MutableSortable[T any] interface {
	SortInPlace(fn func(T, T) bool)
}

// Stack is implemented by collections supporting LIFO operations that return a
// new slice rather than mutating the receiver.
type Stack[T any] interface {
	Push(element T) []T
	Pop() (T, bool, []T)
	PeekEnd() (T, bool)
}

// MutableStack is implemented by collections supporting in-place LIFO
// operations that modify the receiver.
type MutableStack[T any] interface {
	PushInPlace(element T)
	PopInPlace() (T, bool)
}

// Queue is implemented by collections supporting FIFO operations that return a
// new slice rather than mutating the receiver.
type Queue[T any] interface {
	Enqueue(element T) []T
	Dequeue() (T, bool, []T)
	PeekFront() (T, bool)
}

// MutableQueue is implemented by collections supporting in-place FIFO
// operations that modify the receiver.
type MutableQueue[T any] interface {
	EnqueueInPlace(element T)
	DequeueInPlace() (T, bool)
}
