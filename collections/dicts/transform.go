package dicts

import "github.com/pickeringtech/go-collections/constraints"

// Map, MapSorted and Reduce are free functions rather than methods because Go
// still does not allow method type parameters (golang/go#49085). A general
// transform to a different key/value type needs fresh type parameters on the
// operation itself, which only a free function can express. Filter remains a
// method because it keeps the same key/value types.
//
// Unlike Filter, which a Dict implements directly, Map returns the Dict
// interface (backed by the default Hash implementation) so its result chains
// straight into other collection helpers. An empty input Dict yields an
// initialised, non-nil empty Dict. The input itself must be a non-nil Dict:
// these helpers call methods on it and do not guard against a nil Dict value.
//
// FlatMap is intentionally omitted: flattening a dictionary of dictionaries has
// no unambiguous key-merging semantics, so it is left out until a concrete
// use-case appears.

// Map applies fn to every key-value pair of d, returning a new Dict built from
// the (key, value) pairs it produces. The output key and value types (OK, OV)
// may differ from the input types (K, V). If fn maps two input keys to the same
// output key, the later pair (in iteration order) wins; iteration order over a
// Dict is unspecified.
//
// The result is always a Hash-backed Dict, regardless of d's backing. This is a
// deliberate downgrade: because Map may change the key type and OK is only
// constrained to comparable, the output key need not be ordered, so a sorted
// (Tree) input cannot in general be carried through to a sorted output. When the
// output key type is constraints.Ordered and you want to preserve sorted
// iteration, use MapSorted.
func Map[K comparable, V any, OK comparable, OV any](d Dict[K, V], fn func(K, V) (OK, OV)) Dict[OK, OV] {
	out := NewHash[OK, OV]()
	d.ForEach(func(key K, value V) {
		ok, ov := fn(key, value)
		out.PutInPlace(ok, ov)
	})
	return out
}

// MapSorted is the ordered-preserving counterpart to Map: it applies fn to every
// key-value pair of d and returns a Tree-backed SortedDict whose keys are kept in
// ascending order. It requires the output key type OK to be constraints.Ordered,
// which is what lets the result remain sorted. If fn maps two input keys to the
// same output key, the later pair (in iteration order over d) wins.
//
// Use MapSorted when you want a sorted result; use Map when the output key type
// is merely comparable or sorted iteration is not needed.
func MapSorted[K comparable, V any, OK constraints.Ordered, OV any](d Dict[K, V], fn func(K, V) (OK, OV)) SortedDict[OK, OV] {
	out := NewTree[OK, OV]()
	d.ForEach(func(key K, value V) {
		ok, ov := fn(key, value)
		out.PutInPlace(ok, ov)
	})
	return out
}

// Reduce folds d into a single accumulated value of type A, starting from init
// and applying fn to each key-value pair. Iteration order over a Dict is
// unspecified, so fn should be order-independent. For an empty Dict it returns
// init unchanged.
func Reduce[K comparable, V any, A any](d Dict[K, V], init A, fn func(A, K, V) A) A {
	accumulator := init
	d.ForEach(func(key K, value V) {
		accumulator = fn(accumulator, key, value)
	})
	return accumulator
}
