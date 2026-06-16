package dicts

import (
	"cmp"
	"iter"
)

// Interface guards to ensure Tree implements the ordered interfaces.
var _ SortedDict[string, int] = &Tree[string, int]{}
var _ MutableSortedDict[string, int] = &Tree[string, int]{}

// Min returns the entry with the smallest key.
// Returns the key, value, and true if the tree is non-empty; zero values and false otherwise.
func (t *Tree[K, V]) Min() (K, V, bool) {
	if t.root == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	n := t.findMin(t.root)
	return n.Key, n.Value, true
}

// Max returns the entry with the largest key.
// Returns the key, value, and true if the tree is non-empty; zero values and false otherwise.
func (t *Tree[K, V]) Max() (K, V, bool) {
	if t.root == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	n := t.root
	for n.Right != nil {
		n = n.Right
	}
	return n.Key, n.Value, true
}

// Floor returns the entry with the largest key less than or equal to the given key.
// Returns the key, value, and true if such an entry exists; zero values and false otherwise.
func (t *Tree[K, V]) Floor(key K) (K, V, bool) {
	var best *node[K, V]
	current := t.root
	for current != nil {
		comparison := cmp.Compare(key, current.Key)
		if comparison == 0 {
			return current.Key, current.Value, true
		}
		if comparison < 0 {
			current = current.Left
		} else {
			best = current
			current = current.Right
		}
	}
	if best == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	return best.Key, best.Value, true
}

// Ceiling returns the entry with the smallest key greater than or equal to the given key.
// Returns the key, value, and true if such an entry exists; zero values and false otherwise.
func (t *Tree[K, V]) Ceiling(key K) (K, V, bool) {
	var best *node[K, V]
	current := t.root
	for current != nil {
		comparison := cmp.Compare(key, current.Key)
		if comparison == 0 {
			return current.Key, current.Value, true
		}
		if comparison > 0 {
			current = current.Right
		} else {
			best = current
			current = current.Left
		}
	}
	if best == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	return best.Key, best.Value, true
}

// Range returns all entries whose key is within the inclusive range [lo, hi],
// in ascending key order. Returns a non-nil (possibly empty) slice.
func (t *Tree[K, V]) Range(lo, hi K) []Pair[K, V] {
	result := make([]Pair[K, V], 0)
	t.rangeInto(t.root, lo, hi, &result)
	return result
}

// rangeInto appends, in ascending order, the entries of the subtree rooted at n
// whose key falls within the inclusive range [lo, hi]. Subtrees that cannot
// contain in-range keys are pruned.
func (t *Tree[K, V]) rangeInto(n *node[K, V], lo, hi K, out *[]Pair[K, V]) {
	if n == nil {
		return
	}
	// Only descend left when keys smaller than n.Key could be in range.
	if cmp.Compare(lo, n.Key) < 0 {
		t.rangeInto(n.Left, lo, hi, out)
	}
	if cmp.Compare(lo, n.Key) <= 0 && cmp.Compare(n.Key, hi) <= 0 {
		*out = append(*out, Pair[K, V]{Key: n.Key, Value: n.Value})
	}
	// Only descend right when keys larger than n.Key could be in range.
	if cmp.Compare(hi, n.Key) > 0 {
		t.rangeInto(n.Right, lo, hi, out)
	}
}

// All returns an iterator over all entries in ascending key order.
func (t *Tree[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		t.ascend(t.root, yield)
	}
}

// ascend performs an in-order traversal of the subtree rooted at n, yielding each
// entry. It returns false as soon as yield returns false so iteration stops early.
func (t *Tree[K, V]) ascend(n *node[K, V], yield func(K, V) bool) bool {
	if n == nil {
		return true
	}
	if !t.ascend(n.Left, yield) {
		return false
	}
	if !yield(n.Key, n.Value) {
		return false
	}
	return t.ascend(n.Right, yield)
}

// Backward returns an iterator over all entries in descending key order.
func (t *Tree[K, V]) Backward() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		t.descend(t.root, yield)
	}
}

// descend performs a reverse in-order traversal of the subtree rooted at n,
// yielding each entry. It returns false as soon as yield returns false.
func (t *Tree[K, V]) descend(n *node[K, V], yield func(K, V) bool) bool {
	if n == nil {
		return true
	}
	if !t.descend(n.Right, yield) {
		return false
	}
	if !yield(n.Key, n.Value) {
		return false
	}
	return t.descend(n.Left, yield)
}

// RangeAll returns an iterator over the entries whose key is within the inclusive
// range [lo, hi], in ascending key order.
func (t *Tree[K, V]) RangeAll(lo, hi K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		t.rangeAscend(t.root, lo, hi, yield)
	}
}

// rangeAscend yields, in ascending order, the in-range entries of the subtree
// rooted at n, pruning subtrees that cannot contain in-range keys. It returns
// false as soon as yield returns false so iteration stops early.
func (t *Tree[K, V]) rangeAscend(n *node[K, V], lo, hi K, yield func(K, V) bool) bool {
	if n == nil {
		return true
	}
	if cmp.Compare(lo, n.Key) < 0 {
		if !t.rangeAscend(n.Left, lo, hi, yield) {
			return false
		}
	}
	if cmp.Compare(lo, n.Key) <= 0 && cmp.Compare(n.Key, hi) <= 0 {
		if !yield(n.Key, n.Value) {
			return false
		}
	}
	if cmp.Compare(hi, n.Key) > 0 {
		return t.rangeAscend(n.Right, lo, hi, yield)
	}
	return true
}
