package lists

import "iter"

// Return-type policy for immutable list operations
//
// Every non-mutating transform below (Filter, Insert, RemoveAt, Remove, Sort,
// Push, Pop, Enqueue, Dequeue) returns a List[T] rather than a plain []T, so
// results chain straight into further list operations — mirroring the rich
// return types of dicts.Filter (Dict[K,V]) and sets.Filter (Set[T]), and the
// lists.Map free function (List[U]). The concrete value is an Array-backed
// List[T] regardless of the receiver's implementation; callers that need a raw
// slice call AsSlice on the result (zero-overhead for Array, whose AsSlice
// returns its backing slice). AsSlice on the Convertible interface remains the
// single, explicit slice escape hatch.

// Filterable is implemented by collections that can be filtered into a new
// List without modifying the receiver.
type Filterable[T any] interface {
	Filter(fn func(T) bool) List[T]
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
	// IsEmpty returns true if the list contains no elements.
	IsEmpty() bool
}

// Convertible is implemented by collections that can be converted into a slice.
// AsSlice returns an independent copy of the elements: mutating the returned
// slice never affects the collection's backing storage, and the collection's
// later mutations never affect a previously returned slice.
type Convertible[T any] interface {
	AsSlice() []T
}

// Insertable is implemented by collections that can produce a new List with
// elements inserted at a given index, without modifying the receiver.
type Insertable[T any] interface {
	Insert(index int, element ...T) List[T]
}

// MutableInsertable is implemented by collections that can insert elements at a
// given index in place, modifying the receiver.
type MutableInsertable[T any] interface {
	InsertInPlace(index int, element ...T)
}

// Iterable is implemented by collections that can be iterated over, with or
// without element indices, via callbacks or range-over-func iterators.
type Iterable[T any] interface {
	ForEach(fn EachFunc[T])
	ForEachWithIndex(fn IndexedEachFunc[T])
	// All returns an iterator over index/value pairs, front to back, suitable
	// for use with range-over-func.
	All() iter.Seq2[int, T]
	// Values returns an iterator over values, front to back, suitable for use
	// with range-over-func.
	Values() iter.Seq[T]
	// Backward returns an iterator over index/value pairs, back to front;
	// indices still count from the front (the last element has index
	// Length()-1).
	Backward() iter.Seq2[int, T]
}

// Removable is implemented by collections that can produce a new List with an
// element removed, without modifying the receiver.
//
// Index-based removal is unambiguous and needs no constraint on T. Value-based
// removal requires equality; because lists are parameterized [T any] (and so
// cannot use ==), Remove compares with reflect.DeepEqual, consistent with
// dicts.ContainsValue. Callers that want native == semantics for a comparable
// element type should use a ComparableList.
type Removable[T any] interface {
	// RemoveAt returns a new List with the element at index removed. If index
	// is out of bounds, the returned List holds the list's elements unchanged.
	RemoveAt(index int) List[T]

	// Remove returns a new List with the first element deeply equal (per
	// reflect.DeepEqual) to element removed. If no element matches, the returned
	// List holds the list's elements unchanged.
	Remove(element T) List[T]
}

// MutableRemovable is implemented by collections that can remove elements in
// place, modifying the receiver. See Removable for the equality semantics of
// value-based removal.
type MutableRemovable[T any] interface {
	// RemoveAtInPlace removes the element at index, returning it and true if the
	// index was in bounds; the zero value and false otherwise.
	RemoveAtInPlace(index int) (T, bool)

	// RemoveInPlace removes the first element deeply equal (per
	// reflect.DeepEqual) to element, reporting whether an element was removed.
	RemoveInPlace(element T) bool

	// Clear removes all elements from the list.
	Clear()
}

// List is the read-oriented collection interface, combining filtering,
// indexing, insertion, iteration, removal, searching, sorting, stack and queue
// behaviours along with conversion to a slice.
type List[T any] interface {
	Filterable[T]
	Indexable[T]
	Insertable[T]
	Iterable[T]
	Removable[T]
	Searchable[T]
	Sortable[T]
	Stack[T]
	Queue[T]
	Convertible[T]
}

// MutableList extends List with the in-place mutation operations for filtering,
// insertion, removal, sorting, stack and queue behaviours.
type MutableList[T any] interface {
	List[T]
	MutableFilterable[T]
	MutableInsertable[T]
	MutableRemovable[T]
	MutableSortable[T]
	MutableStack[T]
	MutableQueue[T]
}

// Searchable is implemented by collections that can be searched with predicates.
//
// AllMatch, AnyMatch, NoneMatch and Find form the search core shared across the
// lists, dicts and sets families. FindIndex is a deliberate list-specific
// extension that reflects the positional shape of a list.
type Searchable[T any] interface {
	AllMatch(fn func(T) bool) bool
	AnyMatch(fn func(T) bool) bool
	NoneMatch(fn func(T) bool) bool
	Find(fn func(T) bool) (T, bool)
	FindIndex(fn func(T) bool) int
}

// Sortable is implemented by collections that can be sorted into a new List,
// without modifying the receiver, using a less-than comparison.
type Sortable[T any] interface {
	Sort(fn func(T, T) bool) List[T]
}

// MutableSortable is implemented by collections that can be sorted in place.
type MutableSortable[T any] interface {
	SortInPlace(fn func(T, T) bool)
}

// Stack is implemented by collections supporting LIFO operations that return a
// new List rather than mutating the receiver.
type Stack[T any] interface {
	Push(element T) List[T]
	Pop() (T, bool, List[T])
	PeekEnd() (T, bool)
}

// MutableStack is implemented by collections supporting in-place LIFO
// operations that modify the receiver.
type MutableStack[T any] interface {
	PushInPlace(element T)
	PopInPlace() (T, bool)
}

// Queue is implemented by collections supporting FIFO operations that return a
// new List rather than mutating the receiver.
type Queue[T any] interface {
	Enqueue(element T) List[T]
	Dequeue() (T, bool, List[T])
	PeekFront() (T, bool)
}

// MutableQueue is implemented by collections supporting in-place FIFO
// operations that modify the receiver.
type MutableQueue[T any] interface {
	EnqueueInPlace(element T)
	DequeueInPlace() (T, bool)
}
