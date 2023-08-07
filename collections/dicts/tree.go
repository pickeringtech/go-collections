package dicts

type node[K comparable, V any] struct {
	Key   K
	Value V
	Left  *node[K, V]
	Right *node[K, V]
}

type Tree[K comparable, V any] struct {
	Root *node[K, V]
}

func NewTree[K comparable, V any](entries ...Pair[K, V]) Tree[K, V] {
	t := Tree[K, V]{}
	for _, entry := range entries {
		t.Put(entry.Key, entry.Value)
	}
	return t
}

func (t Tree[K, V]) Put(key K, value V) {
	if t.Root == nil {
		t.Root = &node[K, V]{Key: key, Value: value}
		return
	}
	//t.Root = put(t.Root, key, value)
}
