package dicts_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// The tree orders keys exclusively through cmp.Compare, which defines a total
// order over float64 even in the presence of NaN: every NaN compares equal to
// every other NaN, and a NaN sorts as less than every non-NaN value (including
// -Inf). These tests pin that contract down so a future switch to a raw `<`
// comparison — which would make NaN compare unordered and silently corrupt the
// BST — fails loudly.
//
// Note this makes the tree *better behaved* than a native Go map, where a NaN
// key can be stored but never read back and accumulates a fresh entry on every
// insert. That is exactly why these tests assert against hand-computed
// expectations rather than a map oracle.

// TestTree_NaNKey_RoundTrips verifies that a NaN key can be stored, found, read
// back and removed — the basic insert/lookup/delete consistency a raw `<`
// comparison would break.
func TestTree_NaNKey_RoundTrips(t *testing.T) {
	nan := math.NaN()
	tree := dicts.NewTree[float64, string]()

	tree.PutInPlace(nan, "nan-value")

	if !tree.Contains(nan) {
		t.Fatal("Contains(NaN) = false, want true")
	}
	if got, ok := tree.Get(nan, "missing"); !ok || got != "nan-value" {
		t.Fatalf("Get(NaN) = (%q, %v), want (%q, true)", got, ok, "nan-value")
	}
	if tree.Length() != 1 {
		t.Fatalf("Length() = %d, want 1", tree.Length())
	}

	removed, ok := tree.RemoveInPlace(nan)
	if !ok || removed != "nan-value" {
		t.Fatalf("RemoveInPlace(NaN) = (%q, %v), want (%q, true)", removed, ok, "nan-value")
	}
	if tree.Contains(nan) {
		t.Fatal("Contains(NaN) = true after removal, want false")
	}
	if tree.Length() != 0 {
		t.Fatalf("Length() = %d after removal, want 0", tree.Length())
	}
}

// TestTree_NaNKeys_CollapseToOne verifies that all NaN bit patterns are treated
// as a single key: re-putting NaN overwrites rather than duplicating, so the
// tree never accumulates phantom entries the way a native map does.
func TestTree_NaNKeys_CollapseToOne(t *testing.T) {
	// math.NaN() and a NaN from a different operation can have different bit
	// patterns, yet cmp.Compare must treat them as the same key.
	nan1 := math.NaN()
	nan2 := math.Inf(1) - math.Inf(1) // 0 * Inf style NaN, distinct bit pattern

	tree := dicts.NewTree[float64, int]()
	tree.PutInPlace(nan1, 1)
	tree.PutInPlace(nan2, 2)

	if tree.Length() != 1 {
		t.Fatalf("Length() = %d after two NaN puts, want 1", tree.Length())
	}
	if got, ok := tree.Get(nan1, -1); !ok || got != 2 {
		t.Fatalf("Get(NaN) = (%d, %v), want (2, true) — second put should overwrite", got, ok)
	}
}

// TestTree_NaNKey_SortsAsMinimum verifies that NaN occupies the minimum position
// in the total order: it is the Min, it leads ordered iteration, and Floor/
// Ceiling resolve it correctly against the other special float values.
func TestTree_NaNKey_SortsAsMinimum(t *testing.T) {
	nan := math.NaN()
	negInf := math.Inf(-1)
	posInf := math.Inf(1)

	tree := dicts.NewTree[float64, string]()
	tree.PutInPlace(posInf, "+inf")
	tree.PutInPlace(1.5, "one-point-five")
	tree.PutInPlace(negInf, "-inf")
	tree.PutInPlace(nan, "nan")

	// Min must be NaN — it sorts below even -Inf.
	minKey, _, ok := tree.Min()
	if !ok || !math.IsNaN(minKey) {
		t.Fatalf("Min() key = %v (ok=%v), want NaN", minKey, ok)
	}

	// Max must be +Inf.
	if maxKey, _, ok := tree.Max(); !ok || maxKey != posInf {
		t.Fatalf("Max() key = %v (ok=%v), want +Inf", maxKey, ok)
	}

	// Ordered iteration must place NaN first, then ascending non-NaN values.
	var order []float64
	tree.ForEachKey(func(k float64) { order = append(order, k) })
	want := []float64{nan, negInf, 1.5, posInf}
	if len(order) != len(want) {
		t.Fatalf("iteration produced %d keys, want %d: %v", len(order), len(want), order)
	}
	if !math.IsNaN(order[0]) {
		t.Fatalf("first key in order = %v, want NaN", order[0])
	}
	for i := 1; i < len(want); i++ {
		if order[i] != want[i] {
			t.Fatalf("order[%d] = %v, want %v", i, order[i], want[i])
		}
	}

	// Floor(NaN) and Ceiling(NaN) must both resolve to the NaN entry itself.
	if k, _, ok := tree.Floor(nan); !ok || !math.IsNaN(k) {
		t.Fatalf("Floor(NaN) key = %v (ok=%v), want NaN", k, ok)
	}
	if k, _, ok := tree.Ceiling(nan); !ok || !math.IsNaN(k) {
		t.Fatalf("Ceiling(NaN) key = %v (ok=%v), want NaN", k, ok)
	}

	// Ceiling(-Inf) must skip NaN and return -Inf, since NaN < -Inf.
	if k, _, ok := tree.Ceiling(negInf); !ok || k != negInf {
		t.Fatalf("Ceiling(-Inf) key = %v (ok=%v), want -Inf", k, ok)
	}
}

// TestTree_NegativeZeroKey verifies the other float edge case cmp.Compare
// folds together: -0.0 and +0.0 are the same key, so the second put overwrites.
func TestTree_NegativeZeroKey(t *testing.T) {
	negZero := math.Copysign(0, -1)
	posZero := 0.0

	tree := dicts.NewTree[float64, string]()
	tree.PutInPlace(negZero, "neg-zero")
	tree.PutInPlace(posZero, "pos-zero")

	if tree.Length() != 1 {
		t.Fatalf("Length() = %d, want 1 (-0.0 and +0.0 are the same key)", tree.Length())
	}
	if got, ok := tree.Get(negZero, "missing"); !ok || got != "pos-zero" {
		t.Fatalf("Get(-0.0) = (%q, %v), want (%q, true)", got, ok, "pos-zero")
	}
}

// TestTree_NaNKey_RemoveWithNonNaNNeighbours verifies that removing the NaN key
// from a populated tree leaves the remaining ordered structure intact — the
// removal path also routes through cmp.Compare, so a regression there would
// strand or duplicate neighbouring keys.
func TestTree_NaNKey_RemoveWithNonNaNNeighbours(t *testing.T) {
	nan := math.NaN()
	tree := dicts.NewTree[float64, int]()
	for i, v := range []float64{3, 1, 4, 1.5, 9, 2.6} {
		tree.PutInPlace(v, i)
	}
	tree.PutInPlace(nan, 99)

	if _, ok := tree.RemoveInPlace(nan); !ok {
		t.Fatal("RemoveInPlace(NaN) ok = false, want true")
	}
	if tree.Contains(nan) {
		t.Fatal("Contains(NaN) = true after removal, want false")
	}

	// The remaining keys must still iterate in ascending order with no NaN.
	keys := tree.Keys()
	wantKeys := []float64{1, 1.5, 2.6, 3, 4, 9}
	if len(keys) != len(wantKeys) {
		t.Fatalf("Keys() = %v, want %v", keys, wantKeys)
	}
	for i, k := range keys {
		if k != wantKeys[i] {
			t.Fatalf("Keys()[%d] = %v, want %v", i, k, wantKeys[i])
		}
	}
}
