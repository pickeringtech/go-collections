package lists

// ComparableList extends MutableList with value-equality queries that are only
// possible when the element type is comparable.
//
// The base list interfaces are deliberately parameterized [T any], so they
// cannot offer == based operations such as membership testing — the value-based
// Remove on those interfaces falls back to reflect.DeepEqual. ComparableList
// fills the gap for the common case where T is comparable, exposing Contains
// backed by the native == operator.
type ComparableList[T comparable] interface {
	MutableList[T]

	// Contains reports whether element is present in the list, compared with the
	// == operator.
	Contains(element T) bool

	// IndexOf returns the index of the first element equal (==) to element, or -1
	// if no element matches.
	IndexOf(element T) int
}

// Comparable adapts any MutableList[T] into a ComparableList[T] by adding the
// == based queries. It embeds the wrapped list, so every list operation is
// available directly and operates on the same underlying data.
type Comparable[T comparable] struct {
	MutableList[T]
}

// Interface guard
var _ ComparableList[int] = &Comparable[int]{MutableList: NewArray[int]()}

// NewComparable creates a ComparableList backed by an Array seeded with the
// given elements, preserving their order.
func NewComparable[T comparable](elements ...T) *Comparable[T] {
	return &Comparable[T]{MutableList: NewArray(elements...)}
}

// NewComparableFrom wraps an existing MutableList so it gains the == based
// queries. The wrapper shares the underlying list, so mutations through either
// are visible to both. This allows any list implementation — including the
// concurrent ones — to be used as a ComparableList.
func NewComparableFrom[T comparable](list MutableList[T]) *Comparable[T] {
	return &Comparable[T]{MutableList: list}
}

// Contains reports whether element is present in the list, compared with the ==
// operator.
func (c *Comparable[T]) Contains(element T) bool {
	return c.IndexOf(element) >= 0
}

// IndexOf returns the index of the first element equal (==) to element, or -1 if
// no element matches.
func (c *Comparable[T]) IndexOf(element T) int {
	return c.FindIndex(func(value T) bool { return value == element })
}
