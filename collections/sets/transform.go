package sets

// Map and Reduce are free functions rather than methods because Go still does
// not allow method type parameters (golang/go#49085). A general T -> U
// transform needs a fresh type parameter on the operation itself, which only a
// free function can express. Filter remains a method because it keeps the same
// element type.
//
// Unlike Filter, which a Set implements directly, Map returns the Set interface
// (backed by the default Hash implementation) so its result chains straight
// into other collection helpers. Empty or nil input yields an initialised,
// non-nil empty Set.

// Map applies fn to every element of s, returning a new Set holding the results.
// The output element type U may differ from the input type T. Because a Set
// deduplicates, distinct inputs that fn maps to the same value collapse into a
// single element, so the result may be smaller than the input.
func Map[T comparable, U comparable](s Set[T], fn func(T) U) Set[U] {
	out := NewHash[U]()
	s.ForEach(func(element T) {
		out.AddInPlace(fn(element))
	})
	return out
}

// Reduce folds s into a single accumulated value of type A, starting from init
// and applying fn to each element. Iteration order over a Set is unspecified, so
// fn should be order-independent. For an empty or nil Set it returns init
// unchanged.
func Reduce[T comparable, A any](s Set[T], init A, fn func(A, T) A) A {
	accumulator := init
	s.ForEach(func(element T) {
		accumulator = fn(accumulator, element)
	})
	return accumulator
}
