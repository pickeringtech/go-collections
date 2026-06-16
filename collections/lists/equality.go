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
