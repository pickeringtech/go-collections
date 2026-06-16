package dicts

import "iter"

// seqPairs returns an iterator over the key/value pairs.
func seqPairs[K comparable, V any](pairs []Pair[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, pair := range pairs {
			if !yield(pair.Key, pair.Value) {
				return
			}
		}
	}
}

// seqOf returns an iterator over the given items.
func seqOf[T any](items []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	}
}

// FromSeq2 builds a Hash dictionary from the key/value pairs produced by seq.
// When seq yields the same key more than once, the last value wins. It is the
// inbound counterpart to the All iterator.
func FromSeq2[K comparable, V any](seq iter.Seq2[K, V]) Hash[K, V] {
	h := make(Hash[K, V])
	for key, value := range seq {
		h[key] = value
	}
	return h
}

// All returns an iterator over key/value pairs. Iteration order is unspecified.
func (h Hash[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key, value := range h {
			if !yield(key, value) {
				return
			}
		}
	}
}

// KeysSeq returns an iterator over the keys. Iteration order is unspecified.
func (h Hash[K, V]) KeysSeq() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range h {
			if !yield(key) {
				return
			}
		}
	}
}

// ValuesSeq returns an iterator over the values. Iteration order is unspecified.
func (h Hash[K, V]) ValuesSeq() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, value := range h {
			if !yield(value) {
				return
			}
		}
	}
}

// All returns an iterator over key/value pairs in ascending key order.
func (t *Tree[K, V]) All() iter.Seq2[K, V] { return seqPairs(t.Items()) }

// KeysSeq returns an iterator over the keys in ascending order.
func (t *Tree[K, V]) KeysSeq() iter.Seq[K] { return seqOf(t.Keys()) }

// ValuesSeq returns an iterator over the values in ascending key order.
func (t *Tree[K, V]) ValuesSeq() iter.Seq[V] { return seqOf(t.Values()) }

// All returns an iterator over key/value pairs, over a snapshot taken under the
// lock. It is safe for concurrent use. Iteration order is unspecified.
func (ch *ConcurrentHash[K, V]) All() iter.Seq2[K, V] { return seqPairs(ch.Items()) }

// KeysSeq returns an iterator over the keys, over a snapshot taken under the
// lock. It is safe for concurrent use. Iteration order is unspecified.
func (ch *ConcurrentHash[K, V]) KeysSeq() iter.Seq[K] { return seqOf(ch.Keys()) }

// ValuesSeq returns an iterator over the values, over a snapshot taken under the
// lock. It is safe for concurrent use. Iteration order is unspecified.
func (ch *ConcurrentHash[K, V]) ValuesSeq() iter.Seq[V] { return seqOf(ch.Values()) }

// All returns an iterator over key/value pairs, over a snapshot taken under the
// lock. It is safe for concurrent use. Iteration order is unspecified.
func (ch *ConcurrentHashRW[K, V]) All() iter.Seq2[K, V] { return seqPairs(ch.Items()) }

// KeysSeq returns an iterator over the keys, over a snapshot taken under the
// lock. It is safe for concurrent use. Iteration order is unspecified.
func (ch *ConcurrentHashRW[K, V]) KeysSeq() iter.Seq[K] { return seqOf(ch.Keys()) }

// ValuesSeq returns an iterator over the values, over a snapshot taken under the
// lock. It is safe for concurrent use. Iteration order is unspecified.
func (ch *ConcurrentHashRW[K, V]) ValuesSeq() iter.Seq[V] { return seqOf(ch.Values()) }
