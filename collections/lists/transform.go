package lists

// Map, FlatMap and Reduce are free functions rather than methods because Go
// still does not allow method type parameters (golang/go#49085). A general
// T -> U transform needs a fresh type parameter on the operation itself, which
// only a free function can express. Filter remains a method because it is a
// same-element-type transform (T -> List[T]) and so needs no extra type
// parameter.
//
// Like Filter and the other immutable list operations, these functions return
// the List interface (backed by the default Array implementation) so their
// results chain straight into other collection helpers. Empty or nil input
// yields an initialised, non-nil empty List, matching slices.Map.

// Map applies fn to every element of l, in order, returning a new List holding
// the results. The output element type U may differ from the input type T.
func Map[T, U any](l List[T], fn func(T) U) List[U] {
	out := []U{}
	l.ForEach(func(element T) {
		out = append(out, fn(element))
	})
	return NewArray(out...)
}

// FlatMap applies fn to every element of l, in order, concatenating the Lists it
// returns into a single new List. It is the natural choice when each input
// element expands into zero or more output elements.
func FlatMap[T, U any](l List[T], fn func(T) List[U]) List[U] {
	out := []U{}
	l.ForEach(func(element T) {
		out = append(out, fn(element).AsSlice()...)
	})
	return NewArray(out...)
}

// Reduce folds l into a single accumulated value of type A, starting from init
// and applying fn to each element in order. For an empty or nil List it returns
// init unchanged.
func Reduce[T, A any](l List[T], init A, fn func(A, T) A) A {
	accumulator := init
	l.ForEach(func(element T) {
		accumulator = fn(accumulator, element)
	})
	return accumulator
}
