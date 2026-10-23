package lists

import "sort"

// sortStable sorts elements in place by the less-than function lessThan,
// keeping elements that compare equal in their original relative order. Every
// list implementation sorts through it, so the same input and comparator give
// the same order whatever the list's backing.
func sortStable[T any](elements []T, lessThan func(T, T) bool) {
	sort.SliceStable(elements, func(i, j int) bool {
		return lessThan(elements[i], elements[j])
	})
}
