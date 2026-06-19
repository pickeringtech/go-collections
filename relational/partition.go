package relational

// Partition splits input into the values that satisfy predicate and those that
// do not, in a single pass, preserving input order within each side. It is the
// relational complement to a filter: a filter discards the rejected half, but
// data-engineering work often needs both halves — the rows to process and the
// rows to quarantine, the matched and the unmatched. Returning both at once
// avoids walking the slice twice with a predicate and its negation.
//
// The input slice is never mutated. Empty or nil input yields two non-nil empty
// slices, so both results are always safe to range or len.
func Partition[V any](input []V, predicate func(V) bool) (matched []V, unmatched []V) {
	matched = []V{}
	unmatched = []V{}
	for _, value := range input {
		if predicate(value) {
			matched = append(matched, value)
			continue
		}
		unmatched = append(unmatched, value)
	}
	return matched, unmatched
}
