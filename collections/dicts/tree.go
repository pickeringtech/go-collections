package dicts

import (
	"cmp"
	"github.com/pickeringtech/go-collections/constraints"
)

// node represents a single node in the binary search tree.
type node[K constraints.Ordered, V any] struct {
	Key   K
	Value V
	Left  *node[K, V]
	Right *node[K, V]
}

// Tree is a binary search tree implementation of a dictionary.
// It maintains keys in sorted order and provides O(log n) average case performance.
// Note: This is a simple BST without self-balancing, so worst case is O(n).
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
	if t.root == nil {
		t.root = &node[K, V]{Key: key, Value: value}
		t.size++
		return
	}

	current := t.root
	for {
		switch cmp.Compare(key, current.Key) {
		case -1:
			if current.Left == nil {
				current.Left = &node[K, V]{Key: key, Value: value}
				t.size++
				return
			}
			current = current.Left
		case 1:
			if current.Right == nil {
				current.Right = &node[K, V]{Key: key, Value: value}
				t.size++
				return
			}
			current = current.Right
		case 0:
			// Key already exists, update value
			current.Value = value
			return
		}
	}
}

// Put creates a new dictionary with the given key-value pair added or updated.
// Returns the new dictionary without modifying the original.
func (t *Tree[K, V]) Put(key K, value V) Dict[K, V] {
	newTree := t.copy()
	newTree.PutInPlace(key, value)
	return newTree
}

// PutMany creates a new dictionary with all given key-value pairs added or updated.
// Returns the new dictionary without modifying the original.
func (t *Tree[K, V]) PutMany(pairs ...Pair[K, V]) Dict[K, V] {
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
		Key:   n.Key,
		Value: n.Value,
		Left:  t.copyNode(n.Left),
		Right: t.copyNode(n.Right),
	}
}

// ForEach executes the given function for each key-value pair in sorted order.
func (t *Tree[K, V]) ForEach(fn func(key K, value V)) {
	t.inOrderTraversal(t.root, fn)
}

// ForEachKey executes the given function for each key in sorted order.
func (t *Tree[K, V]) ForEachKey(fn func(key K)) {
	t.inOrderTraversal(t.root, func(key K, value V) {
		fn(key)
	})
}

// ForEachValue executes the given function for each value in key-sorted order.
func (t *Tree[K, V]) ForEachValue(fn func(value V)) {
	t.inOrderTraversal(t.root, func(key K, value V) {
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

	t.ForEach(func(key K, value V) {
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

	t.ForEach(func(key K, value V) {
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
func (t *Tree[K, V]) ContainsValue(value V) bool {
	found := false
	t.ForEach(func(key K, v V) {
		if !found && any(v) == any(value) {
			found = true
		}
	})
	return found
}

// Keys returns a slice containing all keys in sorted order.
func (t *Tree[K, V]) Keys() []K {
	keys := make([]K, 0, t.size)
	t.ForEach(func(key K, value V) {
		keys = append(keys, key)
	})
	return keys
}

// Values returns a slice containing all values in key-sorted order.
func (t *Tree[K, V]) Values() []V {
	values := make([]V, 0, t.size)
	t.ForEach(func(key K, value V) {
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
	newTree := t.copy()
	newTree.RemoveInPlace(key)
	return newTree
}

// RemoveMany creates a new dictionary with all given keys removed.
// Returns the new dictionary without modifying the original.
func (t *Tree[K, V]) RemoveMany(keys ...K) Dict[K, V] {
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

	switch cmp.Compare(key, n.Key) {
	case -1:
		var removedValue V
		var found bool
		n.Left, removedValue, found = t.removeNode(n.Left, key)
		return n, removedValue, found
	case 1:
		var removedValue V
		var found bool
		n.Right, removedValue, found = t.removeNode(n.Right, key)
		return n, removedValue, found
	case 0:
		// Found the node to remove
		removedValue := n.Value

		// Case 1: Node has no children
		if n.Left == nil && n.Right == nil {
			return nil, removedValue, true
		}

		// Case 2: Node has only right child
		if n.Left == nil {
			return n.Right, removedValue, true
		}

		// Case 3: Node has only left child
		if n.Right == nil {
			return n.Left, removedValue, true
		}

		// Case 4: Node has both children
		// Find the inorder successor (smallest node in right subtree)
		successor := t.findMin(n.Right)
		n.Key = successor.Key
		n.Value = successor.Value
		n.Right, _, _ = t.removeNode(n.Right, successor.Key)
		return n, removedValue, true
	}

	return n, zeroV, false
}

// findMin finds the node with the minimum key in the subtree rooted at n.
func (t *Tree[K, V]) findMin(n *node[K, V]) *node[K, V] {
	for n.Left != nil {
		n = n.Left
	}
	return n
}
