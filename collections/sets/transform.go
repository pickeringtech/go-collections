package sets

import "github.com/pickeringtech/go-collections/constraints"

// Map, MapSorted and Reduce are free functions rather than methods because Go
// still does not allow method type parameters (golang/go#49085). A general
// T -> U transform needs a fresh type parameter on the operation itself, which
// only a free function can express. Filter remains a method because it keeps the
// same element type.
//
// Unlike Filter, which a Set implements directly, Map returns the Set interface
// (backed by the default Hash implementation) so its result chains straight
// into other collection helpers. An empty input Set yields an initialised,
// non-nil empty Set. The input itself must be a non-nil Set: these helpers call
// methods on it and do not guard against a nil Set value.

// Map applies fn to every element of s, returning a new Set holding the results.
// The output element type U may differ from the input type T. Because a Set
// deduplicates, distinct inputs that fn maps to the same value collapse into a
// single element, so the result may be smaller than the input.
//
// The result is always a Hash-backed Set, regardless of s's backing. This is a
// deliberate downgrade: because Map may change the element type and U is only
// constrained to comparable, the output element need not be ordered, so a sorted
// (TreeSet) input cannot in general be carried through to a sorted output. When
// the output element type is constraints.Ordered and you want to preserve sorted
// iteration, use MapSorted.
func Map[T comparable, U comparable](s Set[T], fn func(T) U) Set[U] {
	out := NewHash[U]()
	s.ForEach(func(element T) {
		out.AddInPlace(fn(element))
	})
	return out
}

// MapSorted is the ordered-preserving counterpart to Map: it applies fn to every
// element of s and returns a TreeSet-backed SortedSet whose elements are kept in
// ascending order. It requires the output element type U to be
// constraints.Ordered, which is what lets the result remain sorted. As with Map,
// distinct inputs that fn maps to the same value collapse into a single element.
//
// Use MapSorted when you want a sorted result; use Map when the output element
// type is merely comparable or sorted iteration is not needed.
func MapSorted[T comparable, U constraints.Ordered](s Set[T], fn func(T) U) SortedSet[U] {
	out := NewTreeSet[U]()
	s.ForEach(func(element T) {
		out.AddInPlace(fn(element))
	})
	return out
}

// Reduce folds s into a single accumulated value of type A, starting from init
// and applying fn to each element. Iteration order over a Set is unspecified, so
// fn should be order-independent. For an empty Set it returns init unchanged.
func Reduce[T comparable, A any](s Set[T], init A, fn func(A, T) A) A {
	accumulator := init
	s.ForEach(func(element T) {
		accumulator = fn(accumulator, element)
	})
	return accumulator
}
