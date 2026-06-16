package lists

import "reflect"

// indexOfDeepEqual returns the index of the first element in elements that is
// deeply equal to target (per reflect.DeepEqual), or -1 if none match.
//
// Lists are parameterized [T any] and so cannot use the == operator. Value-based
// removal therefore relies on reflect.DeepEqual, matching the equality semantics
// dicts uses for ContainsValue. Element types that are comparable and want native
// == semantics should use a ComparableList.
func indexOfDeepEqual[T any](elements []T, target T) int {
	for i := range elements {
		if reflect.DeepEqual(elements[i], target) {
			return i
		}
	}
	return -1
}

// deleteOwned removes the element at index from a slice the caller already owns,
// mutating it in place to avoid a defensive copy. The slice is returned
// unchanged when index is out of bounds. Callers must pass a slice they own
// (for example one from slices.Copy or GetAsSlice), never a shared backing
// array, since the returned slice must be independent of any receiver.
func deleteOwned[T any](owned []T, index int) []T {
	if index < 0 || index >= len(owned) {
		return owned
	}
	return append(owned[:index], owned[index+1:]...)
}
