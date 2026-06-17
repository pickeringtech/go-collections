package dicts

import (
	"cmp"
	"reflect"

	"github.com/pickeringtech/go-collections/constraints"
)

// node represents a single node in the binary search tree.
//
// Height is the length of the longest path from this node down to a leaf,
// counting nodes (a leaf has Height 1). It is maintained by the AVL rebalancing
// in insertNode/removeNode and read by the balancing helpers.
type node[K constraints.Ordered, V any] struct {
	Key    K
	Value  V
	Height int
	Left   *node[K, V]
	Right  *node[K, V]
}

// Tree is a self-balancing binary search tree (AVL) implementation of a
// dictionary. It maintains keys in sorted order and guarantees O(log n)
// worst-case performance for Get/Contains/Put/Remove and the ordered
// navigation operations (Floor/Ceiling/Min/Max/Range), regardless of insertion
// order — including the degenerate sorted-insert case that turns a plain BST
// into a linked list.
// Keys must implement constraints.Ordered (integers, floats, strings).
type Tree[K constraints.Ordered, V any] struct {
	root *node[K, V]
	size int
}

// NewTree creates a new Tree dictionary with the given key-value pairs.
func NewTree[K constraints.Ordered, V any](entries ...Pair[K, V]) *Tree[K, V] {
	t := &Tree[K, V]{}
	for _, entry := range entries {
		t.PutInPlace(entry.Key, entry.Value)
	}
	return t
}

// Interface guards to ensure Tree implements the required interfaces
var _ Dict[string, int] = &Tree[string, int]{}
var _ MutableDict[string, int] = &Tree[string, int]{}

// Get retrieves the value associated with the given key.
// If the key is not found, returns the default value and false.
func (t *Tree[K, V]) Get(key K, defaultValue V) (V, bool) {
	node := t.findNode(key)
	if node != nil {
		return node.Value, true
	}
	return defaultValue, false
}

// Contains checks if the given key exists in the dictionary.
func (t *Tree[K, V]) Contains(key K) bool {
	return t.findNode(key) != nil
}

// Length returns the number of key-value pairs in the dictionary.
func (t *Tree[K, V]) Length() int {
	return t.size
}

// IsEmpty returns true if the dictionary contains no key-value pairs.
func (t *Tree[K, V]) IsEmpty() bool {
	return t.size == 0
}

// findNode searches for a node with the given key.
func (t *Tree[K, V]) findNode(key K) *node[K, V] {
	current := t.root
	for current != nil {
		switch cmp.Compare(key, current.Key) {
		case -1:
			current = current.Left
		case 1:
			current = current.Right
		case 0:
			return current
		}
	}
	return nil
}

// PutInPlace adds or updates the given key-value pair in the dictionary.
func (t *Tree[K, V]) PutInPlace(key K, value V) {
	t.root = t.insertNode(t.root, key, value)
}

// UpdateInPlace reads the value at key, applies fn to it, and stores the result
// back under key, returning the new value. fn receives the current value (the
// zero value if the key is absent) and whether the key existed.
func (t *Tree[K, V]) UpdateInPlace(key K, fn func(old V, existed bool) V) V {
	var old V
	existed := false
	if node := t.findNode(key); node != nil {
		old, existed = node.Value, true
	}
	newValue := fn(old, existed)
	t.root = t.insertNode(t.root, key, newValue)
	return newValue
}

// insertNode inserts key/value into the subtree rooted at n, rebalancing on the
// way back up so the AVL height invariant is preserved. It returns the new root
// of the subtree.
func (t *Tree[K, V]) insertNode(n *node[K, V], key K, value V) *node[K, V] {
	if n == nil {
		t.size++
		return &node[K, V]{Key: key, Value: value, Height: 1}
	}
	switch cmp.Compare(key, n.Key) {
	case -1:
		n.Left = t.insertNode(n.Left, key, value)
	case 1:
		n.Right = t.insertNode(n.Right, key, value)
	case 0:
		// Key already exists, update value. Structure is unchanged.
		n.Value = value
		return n
	}
	return t.rebalance(n)
}

// height returns the AVL height of n, treating a nil subtree as height 0.
func (t *Tree[K, V]) height(n *node[K, V]) int {
	if n == nil {
		return 0
	}
	return n.Height
}

// updateHeight recomputes n.Height from its children. n must be non-nil.
func (t *Tree[K, V]) updateHeight(n *node[K, V]) {
	left, right := t.height(n.Left), t.height(n.Right)
	if left > right {
		n.Height = left + 1
	} else {
		n.Height = right + 1
	}
}

// balanceFactor returns left height minus right height for n (0 for nil). A
// magnitude greater than 1 means the subtree violates the AVL invariant.
func (t *Tree[K, V]) balanceFactor(n *node[K, V]) int {
	if n == nil {
		return 0
	}
	return t.height(n.Left) - t.height(n.Right)
}

// rotateRight performs a right rotation around y and returns the new subtree
// root. y must have a non-nil left child.
func (t *Tree[K, V]) rotateRight(y *node[K, V]) *node[K, V] {
	x := y.Left
	y.Left = x.Right
	x.Right = y
	t.updateHeight(y)
	t.updateHeight(x)
	return x
}

// rotateLeft performs a left rotation around x and returns the new subtree
// root. x must have a non-nil right child.
func (t *Tree[K, V]) rotateLeft(x *node[K, V]) *node[K, V] {
	y := x.Right
	x.Right = y.Left
	y.Left = x
	t.updateHeight(x)
	t.updateHeight(y)
	return y
}

// rebalance updates n's height and, if n violates the AVL balance invariant,
// performs the appropriate single or double rotation. It returns the new root
// of the (now balanced) subtree.
func (t *Tree[K, V]) rebalance(n *node[K, V]) *node[K, V] {
	t.updateHeight(n)
	balance := t.balanceFactor(n)

	// Left-heavy.
	if balance > 1 {
		// Left-Right case: convert to Left-Left first.
		if t.balanceFactor(n.Left) < 0 {
			n.Left = t.rotateLeft(n.Left)
		}
		return t.rotateRight(n)
	}

	// Right-heavy.
	if balance < -1 {
		// Right-Left case: convert to Right-Right first.
		if t.balanceFactor(n.Right) > 0 {
			n.Right = t.rotateRight(n.Right)
		}
		return t.rotateLeft(n)
	}

	return n
}

// Put creates a new dictionary with the given key-value pair added or updated.
// Returns the new dictionary without modifying the original.
func (t *Tree[K, V]) Put(key K, value V) Dict[K, V] {
	return t.put(key, value)
}

// put is the concrete-typed implementation of Put, returning *Tree so callers
// that need the concrete type (e.g. the concurrent wrappers) avoid a type assertion.
func (t *Tree[K, V]) put(key K, value V) *Tree[K, V] {
	newTree := t.copy()
	newTree.PutInPlace(key, value)
	return newTree
}

// PutMany creates a new dictionary with all given key-value pairs added or updated.
// Returns the new dictionary without modifying the original.
func (t *Tree[K, V]) PutMany(pairs ...Pair[K, V]) Dict[K, V] {
	return t.putMany(pairs...)
}

// putMany is the concrete-typed implementation of PutMany, returning *Tree so callers
// that need the concrete type (e.g. the concurrent wrappers) avoid a type assertion.
func (t *Tree[K, V]) putMany(pairs ...Pair[K, V]) *Tree[K, V] {
	newTree := t.copy()
	newTree.PutManyInPlace(pairs...)
	return newTree
}

// PutManyInPlace adds or updates all given key-value pairs in the dictionary.
func (t *Tree[K, V]) PutManyInPlace(pairs ...Pair[K, V]) {
	for _, pair := range pairs {
		t.PutInPlace(pair.Key, pair.Value)
	}
}

// copy creates a deep copy of the tree.
func (t *Tree[K, V]) copy() *Tree[K, V] {
	newTree := &Tree[K, V]{size: t.size}
	newTree.root = t.copyNode(t.root)
	return newTree
}

// copyNode recursively copies a node and its subtrees.
func (t *Tree[K, V]) copyNode(n *node[K, V]) *node[K, V] {
	if n == nil {
		return nil
	}
	return &node[K, V]{
		Key:    n.Key,
		Value:  n.Value,
		Height: n.Height,
		Left:   t.copyNode(n.Left),
		Right:  t.copyNode(n.Right),
	}
}

// ForEach executes the given function for each key-value pair in sorted order.
func (t *Tree[K, V]) ForEach(fn func(key K, value V)) {
	t.inOrderTraversal(t.root, fn)
}

// ForEachKey executes the given function for each key in sorted order.
func (t *Tree[K, V]) ForEachKey(fn func(key K)) {
	t.inOrderTraversal(t.root, func(key K, _ V) {
		fn(key)
	})
}

// ForEachValue executes the given function for each value in key-sorted order.
func (t *Tree[K, V]) ForEachValue(fn func(value V)) {
	t.inOrderTraversal(t.root, func(_ K, value V) {
		fn(value)
	})
}

// inOrderTraversal performs an in-order traversal of the tree.
func (t *Tree[K, V]) inOrderTraversal(n *node[K, V], fn func(key K, value V)) {
	if n != nil {
		t.inOrderTraversal(n.Left, fn)
		fn(n.Key, n.Value)
		t.inOrderTraversal(n.Right, fn)
	}
}

// Filter returns a new dictionary containing only the key-value pairs
// that satisfy the given predicate function.
func (t *Tree[K, V]) Filter(fn func(key K, value V) bool) Dict[K, V] {
	result := NewTree[K, V]()
	t.ForEach(func(key K, value V) {
		if fn(key, value) {
			result.PutInPlace(key, value)
		}
	})
	return result
}

// FilterInPlace removes all key-value pairs that do not satisfy
// the given predicate function, modifying the dictionary in place.
func (t *Tree[K, V]) FilterInPlace(fn func(key K, value V) bool) {
	var toRemove []K
	t.ForEach(func(key K, value V) {
		if !fn(key, value) {
			toRemove = append(toRemove, key)
		}
	})
	for _, key := range toRemove {
		t.RemoveInPlace(key)
	}
}

// AllMatch returns true if every key-value pair satisfies the given predicate.
// It is vacuously true for an empty dictionary. It short-circuits on the first
// pair that fails the predicate.
func (t *Tree[K, V]) AllMatch(fn func(key K, value V) bool) bool {
	return !t.anyInOrder(t.root, func(key K, value V) bool {
		return !fn(key, value)
	})
}

// AnyMatch returns true if at least one key-value pair satisfies the given
// predicate. It is false for an empty dictionary. It short-circuits on the
// first matching pair.
func (t *Tree[K, V]) AnyMatch(fn func(key K, value V) bool) bool {
	return t.anyInOrder(t.root, fn)
}

// anyInOrder reports whether any key-value pair in the subtree rooted at n
// satisfies fn, traversing in order and returning as soon as a match is found.
func (t *Tree[K, V]) anyInOrder(n *node[K, V], fn func(key K, value V) bool) bool {
	if n == nil {
		return false
	}
	if t.anyInOrder(n.Left, fn) {
		return true
	}
	if fn(n.Key, n.Value) {
		return true
	}
	return t.anyInOrder(n.Right, fn)
}

// NoneMatch returns true if no key-value pair satisfies the given predicate.
// It is vacuously true for an empty dictionary.
func (t *Tree[K, V]) NoneMatch(fn func(key K, value V) bool) bool {
	return !t.AnyMatch(fn)
}

// Find returns the first key-value pair that satisfies the given predicate.
// Returns the key, value, and true if found; zero values and false otherwise.
func (t *Tree[K, V]) Find(fn func(key K, value V) bool) (K, V, bool) {
	var foundKey K
	var foundValue V
	found := false

	t.ForEach(func(key K, value V) {
		if !found && fn(key, value) {
			foundKey = key
			foundValue = value
			found = true
		}
	})

	if found {
		return foundKey, foundValue, true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// FindKey returns the first key that satisfies the given predicate.
// Returns the key and true if found; zero value and false otherwise.
func (t *Tree[K, V]) FindKey(fn func(key K) bool) (K, bool) {
	var foundKey K
	found := false

	t.ForEach(func(key K, _ V) {
		if !found && fn(key) {
			foundKey = key
			found = true
		}
	})

	if found {
		return foundKey, true
	}
	var zeroK K
	return zeroK, false
}

// FindValue returns the first value that satisfies the given predicate.
// Returns the value and true if found; zero value and false otherwise.
func (t *Tree[K, V]) FindValue(fn func(value V) bool) (V, bool) {
	var foundValue V
	found := false

	t.ForEach(func(_ K, value V) {
		if !found && fn(value) {
			foundValue = value
			found = true
		}
	})

	if found {
		return foundValue, true
	}
	var zeroV V
	return zeroV, false
}

// ContainsValue checks if the given value exists in the dictionary.
//
// Values are compared with reflect.DeepEqual, matching the equality semantics
// used by list removal. This supports non-comparable value types (slices, maps,
// funcs) without panicking.
func (t *Tree[K, V]) ContainsValue(value V) bool {
	return t.AnyMatch(func(_ K, v V) bool {
		return reflect.DeepEqual(v, value)
	})
}

// Keys returns a slice containing all keys in sorted order.
func (t *Tree[K, V]) Keys() []K {
	keys := make([]K, 0, t.size)
	t.ForEach(func(key K, _ V) {
		keys = append(keys, key)
	})
	return keys
}

// Values returns a slice containing all values in key-sorted order.
func (t *Tree[K, V]) Values() []V {
	values := make([]V, 0, t.size)
	t.ForEach(func(_ K, value V) {
		values = append(values, value)
	})
	return values
}

// Items returns a slice containing all key-value pairs as Pair structs in sorted order.
func (t *Tree[K, V]) Items() []Pair[K, V] {
	items := make([]Pair[K, V], 0, t.size)
	t.ForEach(func(key K, value V) {
		items = append(items, Pair[K, V]{Key: key, Value: value})
	})
	return items
}

// AsMap returns the dictionary as a native Go map.
func (t *Tree[K, V]) AsMap() map[K]V {
	result := make(map[K]V, t.size)
	t.ForEach(func(key K, value V) {
		result[key] = value
	})
	return result
}

// Remove creates a new dictionary with the given key removed.
// Returns the new dictionary without modifying the original.
func (t *Tree[K, V]) Remove(key K) Dict[K, V] {
	return t.remove(key)
}

// remove is the concrete-typed implementation of Remove, returning *Tree so callers
// that need the concrete type (e.g. the concurrent wrappers) avoid a type assertion.
func (t *Tree[K, V]) remove(key K) *Tree[K, V] {
	newTree := t.copy()
	newTree.RemoveInPlace(key)
	return newTree
}

// RemoveMany creates a new dictionary with all given keys removed.
// Returns the new dictionary without modifying the original.
func (t *Tree[K, V]) RemoveMany(keys ...K) Dict[K, V] {
	return t.removeMany(keys...)
}

// removeMany is the concrete-typed implementation of RemoveMany, returning *Tree so callers
// that need the concrete type (e.g. the concurrent wrappers) avoid a type assertion.
func (t *Tree[K, V]) removeMany(keys ...K) *Tree[K, V] {
	newTree := t.copy()
	newTree.RemoveManyInPlace(keys...)
	return newTree
}

// RemoveInPlace removes the given key from the dictionary.
// Returns the removed value and true if the key existed; zero value and false otherwise.
func (t *Tree[K, V]) RemoveInPlace(key K) (V, bool) {
	var removedValue V
	var found bool
	t.root, removedValue, found = t.removeNode(t.root, key)
	if found {
		t.size--
	}
	return removedValue, found
}

// RemoveManyInPlace removes all given keys from the dictionary.
func (t *Tree[K, V]) RemoveManyInPlace(keys ...K) {
	for _, key := range keys {
		t.RemoveInPlace(key)
	}
}

// Clear removes all key-value pairs from the dictionary.
func (t *Tree[K, V]) Clear() {
	t.root = nil
	t.size = 0
}

// removeNode removes a node with the given key from the subtree rooted at n.
// Returns the new root of the subtree, the removed value, and whether the key was found.
func (t *Tree[K, V]) removeNode(n *node[K, V], key K) (*node[K, V], V, bool) {
	var zeroV V
	if n == nil {
		return nil, zeroV, false
	}

	// cmp.Compare returns -1, 0, or +1, so these three branches are exhaustive.
	comparison := cmp.Compare(key, n.Key)
	if comparison < 0 {
		var removedValue V
		var found bool
		n.Left, removedValue, found = t.removeNode(n.Left, key)
		return t.rebalance(n), removedValue, found
	}
	if comparison > 0 {
		var removedValue V
		var found bool
		n.Right, removedValue, found = t.removeNode(n.Right, key)
		return t.rebalance(n), removedValue, found
	}

	// Found the node to remove
	removedValue := n.Value

	// Case 1 & 2: Node has at most one child (right child, possibly nil).
	if n.Left == nil {
		return n.Right, removedValue, true
	}

	// Case 3: Node has only left child
	if n.Right == nil {
		return n.Left, removedValue, true
	}

	// Case 4: Node has both children
	// Find the inorder successor (smallest node in right subtree), copy it into
	// this node, then delete it from the right subtree and rebalance up.
	successor := t.findMin(n.Right)
	n.Key = successor.Key
	n.Value = successor.Value
	n.Right, _, _ = t.removeNode(n.Right, successor.Key)
	return t.rebalance(n), removedValue, true
}

// findMin finds the node with the minimum key in the subtree rooted at n.
func (t *Tree[K, V]) findMin(n *node[K, V]) *node[K, V] {
	for n.Left != nil {
		n = n.Left
	}
	return n
}
