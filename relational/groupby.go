package relational

import "iter"

// KeyFunc derives the grouping key K from a value V. It is the single point
// that decides which group a value belongs to, named (rather than inlined as a
// bare func(V) K) so call sites read as data-engineering verbs — "group by
// department", "count by status" — and so the same key extractor can be shared
// across GroupBy, GroupBySeq and CountBy.
type KeyFunc[K comparable, V any] func(V) K

// GroupBy partitions input into buckets keyed by keyFn(value), returning a map
// from key to the slice of values that produced it. This is the relational
// "GROUP BY" primitive: the foundation the Aggregate pipeline reduces over, so
// it deliberately keeps every value rather than summarising — summarising is
// Aggregate's job, kept separate so any aggregator (sum, mean, your own) can be
// applied to the same grouping.
//
// Within each group the values appear in first-seen order: the order they occur
// in input. That makes results deterministic and lets order-sensitive
// aggregators (first, last, moving windows) behave predictably, which a map's
// own iteration order could not guarantee.
//
// The input slice is never mutated. Empty or nil input yields a non-nil empty
// map, so callers can range or len the result without a nil check.
func GroupBy[K comparable, V any](input []V, keyFn KeyFunc[K, V]) map[K][]V {
	groups := map[K][]V{}
	for _, value := range input {
		key := keyFn(value)
		groups[key] = append(groups[key], value)
	}
	return groups
}

// GroupBySeq is GroupBy over an iter.Seq pull sequence rather than a
// materialised slice, so a stream (a database cursor, a generator, a channel
// adapter) can be grouped without first collecting it into a slice. The
// grouping semantics are identical to GroupBy: first-seen order within each
// group, a non-nil empty map for an empty sequence.
//
// The sequence is consumed exactly once. A nil seq is treated as empty and
// yields a non-nil empty map.
func GroupBySeq[K comparable, V any](seq iter.Seq[V], keyFn KeyFunc[K, V]) map[K][]V {
	groups := map[K][]V{}
	if seq == nil {
		return groups
	}
	for value := range seq {
		key := keyFn(value)
		groups[key] = append(groups[key], value)
	}
	return groups
}

// CountBy returns how many values fall into each group keyed by keyFn(value).
// It is the common "GROUP BY … COUNT(*)" shortcut: when you only need the size
// of each bucket, CountBy avoids retaining every value the way GroupBy does, so
// it is cheaper in memory on large inputs while giving the same key set.
//
// The input slice is never mutated. Empty or nil input yields a non-nil empty
// map.
func CountBy[K comparable, V any](input []V, keyFn KeyFunc[K, V]) map[K]int {
	counts := map[K]int{}
	for _, value := range input {
		counts[keyFn(value)]++
	}
	return counts
}
