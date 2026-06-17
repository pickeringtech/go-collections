package dicts

import (
	"math"
	"math/rand"
	"testing"

	"github.com/pickeringtech/go-collections/constraints"
)

// checkAVL recursively verifies the AVL invariants for the subtree rooted at n:
//   - the balance factor (left height − right height) is within [-1, 1]; and
//   - the cached node.Height equals the height computed by walking the subtree.
//
// It returns the actual height so a parent can validate its own cached height.
// This is what guards against the regression in #97: a plain BST built from
// sorted inserts degenerates into a linked list (height == n), making every
// lookup O(n) and TreeSet.Intersection O(n²). A balanced tree keeps height at
// O(log n).
func checkAVL[K constraints.Ordered, V any](t *testing.T, n *node[K, V]) int {
	t.Helper()
	if n == nil {
		return 0
	}
	leftHeight := checkAVL(t, n.Left)
	rightHeight := checkAVL(t, n.Right)

	balance := leftHeight - rightHeight
	if balance < -1 || balance > 1 {
		t.Fatalf("AVL balance invariant violated at key %v: balance factor %d", n.Key, balance)
	}

	actual := leftHeight
	if rightHeight > actual {
		actual = rightHeight
	}
	actual++
	if n.Height != actual {
		t.Fatalf("cached height wrong at key %v: cached %d, actual %d", n.Key, n.Height, actual)
	}
	return actual
}

// maxAVLHeight returns the tightest known upper bound on the height of an AVL
// tree holding n elements: 1.44 * log2(n + 2) − 0.328, rounded up.
func maxAVLHeight(n int) int {
	if n == 0 {
		return 0
	}
	return int(math.Ceil(1.4405*math.Log2(float64(n)+2) - 0.3277))
}

func TestTree_StaysBalancedOnSortedInsert(t *testing.T) {
	const n = 10_000

	t.Run("ascending", func(t *testing.T) {
		tree := NewTree[int, int]()
		for i := 0; i < n; i++ {
			tree.PutInPlace(i, i)
		}
		height := checkAVL(t, tree.root)
		if max := maxAVLHeight(n); height > max {
			t.Fatalf("ascending inserts left the tree too tall: height %d, expected <= %d", height, max)
		}
	})

	t.Run("descending", func(t *testing.T) {
		tree := NewTree[int, int]()
		for i := n - 1; i >= 0; i-- {
			tree.PutInPlace(i, i)
		}
		height := checkAVL(t, tree.root)
		if max := maxAVLHeight(n); height > max {
			t.Fatalf("descending inserts left the tree too tall: height %d, expected <= %d", height, max)
		}
	})
}

func TestTree_StaysBalancedThroughRandomInsertAndRemove(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	tree := NewTree[int, int]()
	present := make(map[int]bool)

	// Mixed insert/remove churn; verify the invariant after every mutation.
	for step := 0; step < 5_000; step++ {
		key := rng.Intn(500)
		if rng.Intn(2) == 0 {
			tree.PutInPlace(key, key)
			present[key] = true
		} else {
			tree.RemoveInPlace(key)
			delete(present, key)
		}
		checkAVL(t, tree.root)
	}

	// Cached size must match reality, and every surviving key must be findable.
	if tree.size != len(present) {
		t.Fatalf("size mismatch: tree.size %d, expected %d", tree.size, len(present))
	}
	for key := range present {
		if _, ok := tree.Get(key, -1); !ok {
			t.Fatalf("expected key %d to be present after churn", key)
		}
	}

	if max := maxAVLHeight(tree.size); tree.size > 0 {
		if height := checkAVL(t, tree.root); height > max {
			t.Fatalf("tree too tall after churn: height %d, expected <= %d", height, max)
		}
	}
}

func TestTree_RemoveKeepsKeysSorted(t *testing.T) {
	tree := NewTree[int, int]()
	for i := 0; i < 256; i++ {
		tree.PutInPlace(i, i)
	}
	// Remove every third key, then confirm the remaining keys are still sorted
	// and the tree is still balanced.
	for i := 0; i < 256; i += 3 {
		tree.RemoveInPlace(i)
	}
	checkAVL(t, tree.root)

	keys := tree.Keys()
	for i := 1; i < len(keys); i++ {
		if keys[i-1] >= keys[i] {
			t.Fatalf("keys not sorted at index %d: %d >= %d", i, keys[i-1], keys[i])
		}
	}
}
