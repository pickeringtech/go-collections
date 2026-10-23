package dicts

import "iter"

// pairsSeq returns an iterator over the pairs that snapshot returns. snapshot is
// called each time a range over the iterator begins — never when the iterator
// is created — so a stored iterator always reflects the dictionary as it is
// when it is ranged, and the loop body never runs while snapshot holds a lock.
func pairsSeq[K comparable, V any](snapshot func() []Pair[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, pair := range snapshot() {
			if !yield(pair.Key, pair.Value) {
				return
			}
		}
	}
}

// reversePairsSeq is pairsSeq yielding the snapshot in reverse order.
func reversePairsSeq[K comparable, V any](snapshot func() []Pair[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		pairs := snapshot()
		for i := len(pairs) - 1; i >= 0; i-- {
			if !yield(pairs[i].Key, pairs[i].Value) {
				return
			}
		}
	}
}

// itemsSeq is pairsSeq for single-valued iterators (keys or values).
func itemsSeq[T any](snapshot func() []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range snapshot() {
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

// KeysSeq returns an iterator over the keys in ascending order, snapshotted when
// a range over it begins (see Tree.All). The ascending All and Backward
// iterators live in tree_ordered.go alongside the other SortedDict operations.
func (t *Tree[K, V]) KeysSeq() iter.Seq[K] { return itemsSeq(t.Keys) }

// ValuesSeq returns an iterator over the values in ascending key order,
// snapshotted when a range over it begins (see Tree.All).
func (t *Tree[K, V]) ValuesSeq() iter.Seq[V] { return itemsSeq(t.Values) }

// All returns an iterator over key/value pairs. Each range over it snapshots the
// entries under the lock when it begins — not when All is called — and yields
// them after releasing it, so it is safe for concurrent use and the loop body
// may call back into the dictionary. Iteration order is unspecified.
func (ch *ConcurrentHash[K, V]) All() iter.Seq2[K, V] { return pairsSeq(ch.Items) }

// KeysSeq returns an iterator over the keys, snapshotted under the lock when a
// range over it begins (see All). Iteration order is unspecified.
func (ch *ConcurrentHash[K, V]) KeysSeq() iter.Seq[K] { return itemsSeq(ch.Keys) }

// ValuesSeq returns an iterator over the values, snapshotted under the lock when
// a range over it begins (see All). Iteration order is unspecified.
func (ch *ConcurrentHash[K, V]) ValuesSeq() iter.Seq[V] { return itemsSeq(ch.Values) }

// All returns an iterator over key/value pairs. Each range over it snapshots the
// entries under the read lock when it begins — not when All is called — and
// yields them after releasing it, so it is safe for concurrent use and the loop
// body may call back into the dictionary. Iteration order is unspecified.
func (ch *ConcurrentHashRW[K, V]) All() iter.Seq2[K, V] { return pairsSeq(ch.Items) }

// KeysSeq returns an iterator over the keys, snapshotted under the read lock
// when a range over it begins (see All). Iteration order is unspecified.
func (ch *ConcurrentHashRW[K, V]) KeysSeq() iter.Seq[K] { return itemsSeq(ch.Keys) }

// ValuesSeq returns an iterator over the values, snapshotted under the read lock
// when a range over it begins (see All). Iteration order is unspecified.
func (ch *ConcurrentHashRW[K, V]) ValuesSeq() iter.Seq[V] { return itemsSeq(ch.Values) }

// KeysSeq returns an iterator over the keys in ascending order, snapshotted
// under the lock when a range over it begins (see ConcurrentTree.All). The
// ascending All and descending Backward iterators live alongside the
// SortedDict operations.
func (ch *ConcurrentTree[K, V]) KeysSeq() iter.Seq[K] { return itemsSeq(ch.Keys) }

// ValuesSeq returns an iterator over the values in ascending key order,
// snapshotted under the lock when a range over it begins (see
// ConcurrentTree.All).
func (ch *ConcurrentTree[K, V]) ValuesSeq() iter.Seq[V] { return itemsSeq(ch.Values) }

// KeysSeq returns an iterator over the keys in ascending order, snapshotted
// under the read lock when a range over it begins (see ConcurrentTreeRW.All).
func (ch *ConcurrentTreeRW[K, V]) KeysSeq() iter.Seq[K] { return itemsSeq(ch.Keys) }

// ValuesSeq returns an iterator over the values in ascending key order,
// snapshotted under the read lock when a range over it begins (see
// ConcurrentTreeRW.All).
func (ch *ConcurrentTreeRW[K, V]) ValuesSeq() iter.Seq[V] { return itemsSeq(ch.Values) }
