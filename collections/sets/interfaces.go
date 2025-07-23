package sets

// Indexable provides basic element access operations for sets.
type Indexable[T comparable] interface {
	// Contains checks if the given element exists in the set.
	Contains(element T) bool

	// Length returns the number of elements in the set.
	Length() int

	// IsEmpty returns true if the set contains no elements.
	IsEmpty() bool
}

// Iterable provides iteration capabilities for sets.
type Iterable[T comparable] interface {
	// ForEach executes the given function for each element.
	ForEach(fn func(element T))
}

// Filterable provides filtering capabilities for sets.
type Filterable[T comparable] interface {
	// Filter returns a new set containing only the elements
	// that satisfy the given predicate function.
	Filter(fn func(element T) bool) Set[T]
}

// MutableFilterable provides in-place filtering capabilities.
type MutableFilterable[T comparable] interface {
	// FilterInPlace removes all elements that do not satisfy
	// the given predicate function, modifying the set in place.
	FilterInPlace(fn func(element T) bool)
}

// Searchable provides search capabilities for sets.
type Searchable[T comparable] interface {
	// Find returns the first element that satisfies the given predicate.
	// Returns the element and true if found; zero value and false otherwise.
	Find(fn func(element T) bool) (T, bool)

	// AllMatch returns true if all elements satisfy the given predicate.
	AllMatch(fn func(element T) bool) bool

	// AnyMatch returns true if any element satisfies the given predicate.
	AnyMatch(fn func(element T) bool) bool
}

// Convertible provides conversion capabilities for sets.
type Convertible[T comparable] interface {
	// AsSlice returns the set as a slice.
	AsSlice() []T

	// AsMap returns the set as a native Go map with struct{} values.
	AsMap() map[T]struct{}
}

// Insertable provides insertion capabilities for sets.
type Insertable[T comparable] interface {
	// Add creates a new set with the given element added.
	// Returns the new set without modifying the original.
	Add(element T) Set[T]

	// AddMany creates a new set with all given elements added.
	// Returns the new set without modifying the original.
	AddMany(elements ...T) Set[T]

	// Union creates a new set containing all elements from this set and the other set.
	Union(other Set[T]) Set[T]
}

// MutableInsertable provides in-place insertion capabilities.
type MutableInsertable[T comparable] interface {
	// AddInPlace adds the given element to the set.
	AddInPlace(element T)

	// AddManyInPlace adds all given elements to the set.
	AddManyInPlace(elements ...T)

	// UnionInPlace adds all elements from the other set to this set.
	UnionInPlace(other Set[T])
}

// Removable provides removal capabilities for sets.
type Removable[T comparable] interface {
	// Remove creates a new set with the given element removed.
	// Returns the new set without modifying the original.
	Remove(element T) Set[T]

	// RemoveMany creates a new set with all given elements removed.
	// Returns the new set without modifying the original.
	RemoveMany(elements ...T) Set[T]

	// Difference creates a new set containing elements in this set but not in the other set.
	Difference(other Set[T]) Set[T]
}

// MutableRemovable provides in-place removal capabilities.
type MutableRemovable[T comparable] interface {
	// RemoveInPlace removes the given element from the set.
	// Returns true if the element was present and removed; false otherwise.
	RemoveInPlace(element T) bool

	// RemoveManyInPlace removes all given elements from the set.
	RemoveManyInPlace(elements ...T)

	// DifferenceInPlace removes all elements that are present in the other set.
	DifferenceInPlace(other Set[T])

	// Clear removes all elements from the set.
	Clear()
}

// SetOperations provides mathematical set operations.
type SetOperations[T comparable] interface {
	// Intersection creates a new set containing elements present in both sets.
	Intersection(other Set[T]) Set[T]

	// IsSubsetOf returns true if this set is a subset of the other set.
	IsSubsetOf(other Set[T]) bool

	// IsSupersetOf returns true if this set is a superset of the other set.
	IsSupersetOf(other Set[T]) bool

	// IsDisjoint returns true if this set has no elements in common with the other set.
	IsDisjoint(other Set[T]) bool

	// Equals returns true if both sets contain exactly the same elements.
	Equals(other Set[T]) bool
}

// MutableSetOperations provides in-place mathematical set operations.
type MutableSetOperations[T comparable] interface {
	// IntersectionInPlace keeps only elements that are present in both sets.
	IntersectionInPlace(other Set[T])
}

// Set represents an immutable set interface that provides
// comprehensive element operations without modifying the original set.
type Set[T comparable] interface {
	Indexable[T]
	Iterable[T]
	Filterable[T]
	Searchable[T]
	Convertible[T]
	Insertable[T]
	Removable[T]
	SetOperations[T]
}

// MutableSet represents a mutable set interface that provides
// comprehensive element operations with the ability to modify the set in place.
type MutableSet[T comparable] interface {
	Set[T]
	MutableFilterable[T]
	MutableInsertable[T]
	MutableRemovable[T]
	MutableSetOperations[T]
}
