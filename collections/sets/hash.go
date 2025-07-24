package sets

// Hash is a fast set implementation using Go's built-in map for O(1) operations.
// It automatically handles deduplication and provides rich mathematical set operations
// like union, intersection, and difference.
//
// Hash is perfect for:
//   - Membership testing (Contains)
//   - Removing duplicates from data
//   - Mathematical set operations
//   - Fast lookups and insertions
//
// Example usage:
//
//	// Create permission sets
//	adminPerms := sets.NewHash("read", "write", "delete", "admin")
//	userPerms := sets.NewHash("read", "write")
//
//	// Mathematical operations
//	common := adminPerms.Intersection(userPerms)  // {"read", "write"}
//	extra := adminPerms.Difference(userPerms)     // {"delete", "admin"}
//
//	// Membership testing
//	canDelete := adminPerms.Contains("delete")    // true
type Hash[T comparable] map[T]struct{}

// NewHash creates a new Hash set with the given elements.
// Duplicate elements are automatically removed.
//
// Example:
//
//	// Empty set
//	empty := sets.NewHash[string]()
//
//	// Set with initial elements (duplicates removed)
//	colors := sets.NewHash("red", "green", "blue", "red")  // {"red", "green", "blue"}
//
//	// Set from slice (deduplication)
//	items := []string{"apple", "banana", "apple", "cherry"}
//	unique := sets.NewHash(items...)  // {"apple", "banana", "cherry"}
func NewHash[T comparable](values ...T) Hash[T] {
	m := make(Hash[T])
	for _, value := range values {
		m[value] = struct{}{}
	}
	return m
}

// Interface guards to ensure Hash implements the required interfaces
var _ Set[string] = Hash[string]{}
var _ MutableSet[string] = Hash[string]{}

// Contains checks if the given element exists in the set.
// This is the primary operation for membership testing.
//
// Example:
//
//	permissions := sets.NewHash("read", "write", "execute")
//
//	if permissions.Contains("write") {
//		fmt.Println("User can write")  // This will print
//	}
//
//	if !permissions.Contains("delete") {
//		fmt.Println("User cannot delete")  // This will print
//	}
func (h Hash[T]) Contains(element T) bool {
	_, exists := h[element]
	return exists
}

// Length returns the number of elements in the set.
func (h Hash[T]) Length() int {
	return len(h)
}

// IsEmpty returns true if the set contains no elements.
func (h Hash[T]) IsEmpty() bool {
	return len(h) == 0
}

// ForEach executes the given function for each element.
func (h Hash[T]) ForEach(fn func(element T)) {
	for element := range h {
		fn(element)
	}
}

// Filter returns a new set containing only the elements
// that satisfy the given predicate function.
func (h Hash[T]) Filter(fn func(element T) bool) Set[T] {
	result := make(Hash[T])
	for element := range h {
		if fn(element) {
			result[element] = struct{}{}
		}
	}
	return result
}

// FilterInPlace removes all elements that do not satisfy
// the given predicate function, modifying the set in place.
func (h Hash[T]) FilterInPlace(fn func(element T) bool) {
	for element := range h {
		if !fn(element) {
			delete(h, element)
		}
	}
}

// Find returns the first element that satisfies the given predicate.
// Returns the element and true if found; zero value and false otherwise.
func (h Hash[T]) Find(fn func(element T) bool) (T, bool) {
	for element := range h {
		if fn(element) {
			return element, true
		}
	}
	var zero T
	return zero, false
}

// AllMatch returns true if all elements satisfy the given predicate.
func (h Hash[T]) AllMatch(fn func(element T) bool) bool {
	for element := range h {
		if !fn(element) {
			return false
		}
	}
	return true
}

// AnyMatch returns true if any element satisfies the given predicate.
func (h Hash[T]) AnyMatch(fn func(element T) bool) bool {
	for element := range h {
		if fn(element) {
			return true
		}
	}
	return false
}

// AsSlice returns the set as a slice.
func (h Hash[T]) AsSlice() []T {
	result := make([]T, 0, len(h))
	for element := range h {
		result = append(result, element)
	}
	return result
}

// AsMap returns the set as a native Go map with struct{} values.
func (h Hash[T]) AsMap() map[T]struct{} {
	result := make(map[T]struct{}, len(h))
	for element := range h {
		result[element] = struct{}{}
	}
	return result
}

// Add creates a new set with the given element added.
// Returns the new set without modifying the original.
func (h Hash[T]) Add(element T) Set[T] {
	result := make(Hash[T], len(h)+1)
	for e := range h {
		result[e] = struct{}{}
	}
	result[element] = struct{}{}
	return result
}

// AddMany creates a new set with all given elements added.
// Returns the new set without modifying the original.
func (h Hash[T]) AddMany(elements ...T) Set[T] {
	result := make(Hash[T], len(h)+len(elements))
	for e := range h {
		result[e] = struct{}{}
	}
	for _, element := range elements {
		result[element] = struct{}{}
	}
	return result
}

// Union creates a new set containing all elements from this set and the other set.
func (h Hash[T]) Union(other Set[T]) Set[T] {
	result := make(Hash[T], h.Length()+other.Length())
	for e := range h {
		result[e] = struct{}{}
	}
	other.ForEach(func(element T) {
		result[element] = struct{}{}
	})
	return result
}

// AddInPlace adds the given element to the set.
func (h Hash[T]) AddInPlace(element T) {
	h[element] = struct{}{}
}

// AddManyInPlace adds all given elements to the set.
func (h Hash[T]) AddManyInPlace(elements ...T) {
	for _, element := range elements {
		h[element] = struct{}{}
	}
}

// UnionInPlace adds all elements from the other set to this set.
func (h Hash[T]) UnionInPlace(other Set[T]) {
	other.ForEach(func(element T) {
		h[element] = struct{}{}
	})
}

// Remove creates a new set with the given element removed.
// Returns the new set without modifying the original.
func (h Hash[T]) Remove(element T) Set[T] {
	result := make(Hash[T], len(h))
	for e := range h {
		if e != element {
			result[e] = struct{}{}
		}
	}
	return result
}

// RemoveMany creates a new set with all given elements removed.
// Returns the new set without modifying the original.
func (h Hash[T]) RemoveMany(elements ...T) Set[T] {
	// Create a set of elements to remove for efficient lookup
	toRemove := make(map[T]struct{}, len(elements))
	for _, element := range elements {
		toRemove[element] = struct{}{}
	}

	result := make(Hash[T], len(h))
	for e := range h {
		if _, shouldRemove := toRemove[e]; !shouldRemove {
			result[e] = struct{}{}
		}
	}
	return result
}

// Difference creates a new set containing elements in this set but not in the other set.
func (h Hash[T]) Difference(other Set[T]) Set[T] {
	result := make(Hash[T])
	for element := range h {
		if !other.Contains(element) {
			result[element] = struct{}{}
		}
	}
	return result
}

// RemoveInPlace removes the given element from the set.
// Returns true if the element was present and removed; false otherwise.
func (h Hash[T]) RemoveInPlace(element T) bool {
	if _, exists := h[element]; exists {
		delete(h, element)
		return true
	}
	return false
}

// RemoveManyInPlace removes all given elements from the set.
func (h Hash[T]) RemoveManyInPlace(elements ...T) {
	for _, element := range elements {
		delete(h, element)
	}
}

// DifferenceInPlace removes all elements that are present in the other set.
func (h Hash[T]) DifferenceInPlace(other Set[T]) {
	other.ForEach(func(element T) {
		delete(h, element)
	})
}

// Clear removes all elements from the set.
func (h Hash[T]) Clear() {
	for element := range h {
		delete(h, element)
	}
}

// Intersection creates a new set containing elements present in both sets.
func (h Hash[T]) Intersection(other Set[T]) Set[T] {
	result := make(Hash[T])
	for element := range h {
		if other.Contains(element) {
			result[element] = struct{}{}
		}
	}
	return result
}

// IsSubsetOf returns true if this set is a subset of the other set.
func (h Hash[T]) IsSubsetOf(other Set[T]) bool {
	for element := range h {
		if !other.Contains(element) {
			return false
		}
	}
	return true
}

// IsSupersetOf returns true if this set is a superset of the other set.
func (h Hash[T]) IsSupersetOf(other Set[T]) bool {
	return other.IsSubsetOf(h)
}

// IsDisjoint returns true if this set has no elements in common with the other set.
func (h Hash[T]) IsDisjoint(other Set[T]) bool {
	for element := range h {
		if other.Contains(element) {
			return false
		}
	}
	return true
}

// Equals returns true if both sets contain exactly the same elements.
func (h Hash[T]) Equals(other Set[T]) bool {
	if h.Length() != other.Length() {
		return false
	}
	for element := range h {
		if !other.Contains(element) {
			return false
		}
	}
	return true
}

// IntersectionInPlace keeps only elements that are present in both sets.
func (h Hash[T]) IntersectionInPlace(other Set[T]) {
	for element := range h {
		if !other.Contains(element) {
			delete(h, element)
		}
	}
}
