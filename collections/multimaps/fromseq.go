package multimaps

import "iter"

// ListMultimapFromSeq2 builds a list-backed ListMultimap from the key/value
// pairs produced by seq, preserving insertion order and keeping duplicate
// values. It is the inbound counterpart to the All iterator.
func ListMultimapFromSeq2[K comparable, V any](seq iter.Seq2[K, V]) ListMultimap[K, V] {
	m := make(ListMultimap[K, V])
	for key, value := range seq {
		m.PutInPlace(key, value)
	}
	return m
}

// SetMultimapFromSeq2 builds a set-backed SetMultimap from the key/value pairs
// produced by seq, collapsing duplicate values bound to the same key. It is the
// inbound counterpart to the All iterator.
func SetMultimapFromSeq2[K comparable, V comparable](seq iter.Seq2[K, V]) SetMultimap[K, V] {
	m := make(SetMultimap[K, V])
	for key, value := range seq {
		m.PutInPlace(key, value)
	}
	return m
}
