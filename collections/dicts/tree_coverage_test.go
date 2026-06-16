package dicts_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// TestTree_AnyMatch_LeftSubtree forces anyInOrder to find its match in the left
// subtree, exercising the early-return branch of the in-order recursion.
func TestTree_AnyMatch_LeftSubtree(t *testing.T) {
	// Inserting in this order makes "bravo" the root with "alpha" on the left.
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "bravo", Value: 2},
		dicts.Pair[string, int]{Key: "alpha", Value: 1},
		dicts.Pair[string, int]{Key: "charlie", Value: 3},
	)

	if !tree.AnyMatch(func(key string, _ int) bool { return key == "alpha" }) {
		t.Fatalf("AnyMatch should find a key located in the left subtree")
	}
}

// TestTree_Find_NoMatch exercises the not-found return path of Find.
func TestTree_Find_NoMatch(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[string, int]{Key: "alpha", Value: 1},
		dicts.Pair[string, int]{Key: "bravo", Value: 2},
	)

	key, value, ok := tree.Find(func(_ string, value int) bool { return value > 100 })
	if ok {
		t.Fatalf("Find should report no match, got (%q, %d, %v)", key, value, ok)
	}
	if key != "" || value != 0 {
		t.Fatalf("Find should return zero values on no match, got (%q, %d)", key, value)
	}
}

// TestTree_RemoveInPlace_SuccessorHasLeftChild removes a node with two children
// whose in-order successor is not the immediate right child, forcing findMin to
// walk left at least once.
func TestTree_RemoveInPlace_SuccessorHasLeftChild(t *testing.T) {
	// Tree shape (insertion order chosen deliberately):
	//
	//        50
	//       /  \
	//     30    70
	//          /
	//        60
	//          \
	//           65
	//
	// Removing 50 picks the successor 60 from the right subtree; findMin walks
	// from 70 down its left child 60.
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 50, Value: "fifty"},
		dicts.Pair[int, string]{Key: 30, Value: "thirty"},
		dicts.Pair[int, string]{Key: 70, Value: "seventy"},
		dicts.Pair[int, string]{Key: 60, Value: "sixty"},
		dicts.Pair[int, string]{Key: 65, Value: "sixty-five"},
	)

	value, ok := tree.RemoveInPlace(50)
	if !ok || value != "fifty" {
		t.Fatalf("RemoveInPlace(50) = (%q, %v), want (\"fifty\", true)", value, ok)
	}

	if tree.Contains(50) {
		t.Fatalf("tree should no longer contain the removed key 50")
	}

	// The successor (60) must now be the value reachable at the old root's key
	// neighbourhood, and every remaining key must still be present and ordered.
	wantKeys := []int{30, 60, 65, 70}
	gotKeys := tree.Keys()
	if len(gotKeys) != len(wantKeys) {
		t.Fatalf("remaining keys = %v, want %v", gotKeys, wantKeys)
	}
	for i, k := range wantKeys {
		if gotKeys[i] != k {
			t.Fatalf("remaining keys = %v, want %v", gotKeys, wantKeys)
		}
	}
}
