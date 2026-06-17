package sets

import (
	"iter"

	"github.com/pickeringtech/go-collections/collections/dicts"
	"github.com/pickeringtech/go-collections/constraints"
)

// present is the zero-size value stored for every element in the backing tree.
var present = struct{}{}

// TreeSet is a sorted set implementation backed by the binary search tree from
// the dicts package. It keeps elements in sorted order and, on top of the usual
// mathematical set operations, exposes ordered navigation and iteration
// (Min/Max, Floor/Ceiling, Range and ascending/descending iterators).
//
// Elements must implement constraints.Ordered (integers, floats, strings). Use
// TreeSet when you need sorted iteration or range queries; prefer Hash for
// unordered membership testing at O(1).
//
// Backed by dicts.Tree, TreeSet inherits its float key contract: all NaN
// elements collapse to a single element that sorts as the minimum (ahead of
// -Inf), and -0.0 and +0.0 are the same element. See dicts.Tree for details.
type TreeSet[T constraints.Ordered] struct {
	tree *dicts.Tree[T, struct{}]
}

// NewTreeSet creates a new TreeSet with the given elements.
// Duplicate elements are automatically removed.
func NewTreeSet[T constraints.Ordered](elements ...T) *TreeSet[T] {
	s := &TreeSet[T]{tree: dicts.NewTree[T, struct{}]()}
	for _, element := range elements {
		s.tree.PutInPlace(element, present)
	}
	return s
}

// Interface guards to ensure TreeSet implements the required interfaces.
var _ Set[string] = &TreeSet[string]{}
var _ MutableSet[string] = &TreeSet[string]{}
var _ SortedSet[string] = &TreeSet[string]{}
var _ MutableSortedSet[string] = &TreeSet[string]{}

// Contains checks if the given element exists in the set.
func (s *TreeSet[T]) Contains(element T) bool {
	return s.tree.Contains(element)
}

// Length returns the number of elements in the set.
func (s *TreeSet[T]) Length() int {
	return s.tree.Length()
}

// IsEmpty returns true if the set contains no elements.
func (s *TreeSet[T]) IsEmpty() bool {
	return s.tree.IsEmpty()
}

// ForEach executes the given function for each element in ascending order.
func (s *TreeSet[T]) ForEach(fn func(element T)) {
	s.tree.ForEachKey(fn)
}

// Filter returns a new set containing only the elements that satisfy the given
// predicate function.
func (s *TreeSet[T]) Filter(fn func(element T) bool) Set[T] {
	result := NewTreeSet[T]()
	s.tree.ForEachKey(func(element T) {
		if fn(element) {
			result.tree.PutInPlace(element, present)
		}
	})
	return result
}

// FilterInPlace removes all elements that do not satisfy the given predicate
// function, modifying the set in place.
func (s *TreeSet[T]) FilterInPlace(fn func(element T) bool) {
	s.tree.FilterInPlace(func(element T, _ struct{}) bool {
		return fn(element)
	})
}

// Find returns the first element (in ascending order) that satisfies the given predicate.
// Returns the element and true if found; zero value and false otherwise.
func (s *TreeSet[T]) Find(fn func(element T) bool) (T, bool) {
	return s.tree.FindKey(fn)
}

// AllMatch returns true if all elements satisfy the given predicate.
// It is vacuously true for an empty set.
func (s *TreeSet[T]) AllMatch(fn func(element T) bool) bool {
	return s.tree.AllMatch(func(element T, _ struct{}) bool {
		return fn(element)
	})
}

// AnyMatch returns true if any element satisfies the given predicate.
// It is false for an empty set.
func (s *TreeSet[T]) AnyMatch(fn func(element T) bool) bool {
	return s.tree.AnyMatch(func(element T, _ struct{}) bool {
		return fn(element)
	})
}

// NoneMatch returns true if no element satisfies the given predicate.
// It is vacuously true for an empty set.
func (s *TreeSet[T]) NoneMatch(fn func(element T) bool) bool {
	return !s.AnyMatch(fn)
}

// AsSlice returns the set as a slice in ascending order.
func (s *TreeSet[T]) AsSlice() []T {
	return s.tree.Keys()
}

// AsMap returns the set as a native Go map with struct{} values.
func (s *TreeSet[T]) AsMap() map[T]struct{} {
	result := make(map[T]struct{}, s.tree.Length())
	s.tree.ForEachKey(func(element T) {
		result[element] = present
	})
	return result
}

// Add creates a new set with the given element added.
// Returns the new set without modifying the original.
func (s *TreeSet[T]) Add(element T) Set[T] {
	result := s.clone()
	result.tree.PutInPlace(element, present)
	return result
}

// AddMany creates a new set with all given elements added.
// Returns the new set without modifying the original.
func (s *TreeSet[T]) AddMany(elements ...T) Set[T] {
	result := s.clone()
	for _, element := range elements {
		result.tree.PutInPlace(element, present)
	}
	return result
}

// Union creates a new set containing all elements from this set and the other set.
func (s *TreeSet[T]) Union(other Set[T]) Set[T] {
	result := s.clone()
	other.ForEach(func(element T) {
		result.tree.PutInPlace(element, present)
	})
	return result
}

// AddInPlace adds the given element to the set.
func (s *TreeSet[T]) AddInPlace(element T) {
	s.tree.PutInPlace(element, present)
}

// AddManyInPlace adds all given elements to the set.
func (s *TreeSet[T]) AddManyInPlace(elements ...T) {
	for _, element := range elements {
		s.tree.PutInPlace(element, present)
	}
}

// UnionInPlace adds all elements from the other set to this set.
func (s *TreeSet[T]) UnionInPlace(other Set[T]) {
	other.ForEach(func(element T) {
		s.tree.PutInPlace(element, present)
	})
}

// Remove creates a new set with the given element removed.
// Returns the new set without modifying the original.
func (s *TreeSet[T]) Remove(element T) Set[T] {
	result := s.clone()
	result.tree.RemoveInPlace(element)
	return result
}

// RemoveMany creates a new set with all given elements removed.
// Returns the new set without modifying the original.
func (s *TreeSet[T]) RemoveMany(elements ...T) Set[T] {
	result := s.clone()
	result.tree.RemoveManyInPlace(elements...)
	return result
}

// Difference creates a new set containing elements in this set but not in the other set.
func (s *TreeSet[T]) Difference(other Set[T]) Set[T] {
	result := NewTreeSet[T]()
	s.tree.ForEachKey(func(element T) {
		if !other.Contains(element) {
			result.tree.PutInPlace(element, present)
		}
	})
	return result
}

// RemoveInPlace removes the given element from the set.
// Returns true if the element was present and removed; false otherwise.
func (s *TreeSet[T]) RemoveInPlace(element T) bool {
	_, found := s.tree.RemoveInPlace(element)
	return found
}

// RemoveManyInPlace removes all given elements from the set.
func (s *TreeSet[T]) RemoveManyInPlace(elements ...T) {
	s.tree.RemoveManyInPlace(elements...)
}

// DifferenceInPlace removes all elements that are present in the other set.
//
// The elements to remove are snapshotted before any mutation, so passing the
// receiver itself (s.DifferenceInPlace(s), which empties the set) does not
// mutate the backing tree while it is being traversed.
func (s *TreeSet[T]) DifferenceInPlace(other Set[T]) {
	toRemove := other.AsSlice()
	for _, element := range toRemove {
		s.tree.RemoveInPlace(element)
	}
}

// Clear removes all elements from the set.
func (s *TreeSet[T]) Clear() {
	s.tree.Clear()
}

// Intersection creates a new set containing elements present in both sets.
func (s *TreeSet[T]) Intersection(other Set[T]) Set[T] {
	result := NewTreeSet[T]()
	s.tree.ForEachKey(func(element T) {
		if other.Contains(element) {
			result.tree.PutInPlace(element, present)
		}
	})
	return result
}

// IsSubsetOf returns true if this set is a subset of the other set.
func (s *TreeSet[T]) IsSubsetOf(other Set[T]) bool {
	return s.AllMatch(other.Contains)
}

// IsSupersetOf returns true if this set is a superset of the other set.
func (s *TreeSet[T]) IsSupersetOf(other Set[T]) bool {
	return other.IsSubsetOf(s)
}

// IsDisjoint returns true if this set has no elements in common with the other set.
func (s *TreeSet[T]) IsDisjoint(other Set[T]) bool {
	return s.NoneMatch(other.Contains)
}

// Equals returns true if both sets contain exactly the same elements.
func (s *TreeSet[T]) Equals(other Set[T]) bool {
	if s.Length() != other.Length() {
		return false
	}
	return s.IsSubsetOf(other)
}

// IntersectionInPlace keeps only elements that are present in both sets.
func (s *TreeSet[T]) IntersectionInPlace(other Set[T]) {
	s.tree.FilterInPlace(func(element T, _ struct{}) bool {
		return other.Contains(element)
	})
}

// Min returns the smallest element.
// Returns the element and true if the set is non-empty; zero value and false otherwise.
func (s *TreeSet[T]) Min() (T, bool) {
	element, _, ok := s.tree.Min()
	return element, ok
}

// Max returns the largest element.
// Returns the element and true if the set is non-empty; zero value and false otherwise.
func (s *TreeSet[T]) Max() (T, bool) {
	element, _, ok := s.tree.Max()
	return element, ok
}

// Floor returns the largest element less than or equal to the given element.
// Returns the element and true if such an element exists; zero value and false otherwise.
func (s *TreeSet[T]) Floor(element T) (T, bool) {
	found, _, ok := s.tree.Floor(element)
	return found, ok
}

// Ceiling returns the smallest element greater than or equal to the given element.
// Returns the element and true if such an element exists; zero value and false otherwise.
func (s *TreeSet[T]) Ceiling(element T) (T, bool) {
	found, _, ok := s.tree.Ceiling(element)
	return found, ok
}

// Range returns all elements within the inclusive range [lo, hi], in ascending
// order. Returns a non-nil (possibly empty) slice.
func (s *TreeSet[T]) Range(lo, hi T) []T {
	pairs := s.tree.Range(lo, hi)
	result := make([]T, 0, len(pairs))
	for _, pair := range pairs {
		result = append(result, pair.Key)
	}
	return result
}

// All returns an iterator over all elements in ascending order.
func (s *TreeSet[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for element := range s.tree.All() {
			if !yield(element) {
				return
			}
		}
	}
}

// Backward returns an iterator over all elements in descending order.
func (s *TreeSet[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		for element := range s.tree.Backward() {
			if !yield(element) {
				return
			}
		}
	}
}

// RangeAll returns an iterator over the elements within the inclusive range
// [lo, hi], in ascending order.
func (s *TreeSet[T]) RangeAll(lo, hi T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for element := range s.tree.RangeAll(lo, hi) {
			if !yield(element) {
				return
			}
		}
	}
}

// clone returns a deep copy of the set, sharing no state with the original.
func (s *TreeSet[T]) clone() *TreeSet[T] {
	result := NewTreeSet[T]()
	s.tree.ForEachKey(func(element T) {
		result.tree.PutInPlace(element, present)
	})
	return result
}
