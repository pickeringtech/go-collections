package sets_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
)

// TestConcurrentHash_ImmutableOpsReturnConcurrentType asserts that every
// immutable operation on ConcurrentHash returns a new ConcurrentHash (the same
// thread-safe type) rather than a plain, non-thread-safe Hash.
func TestConcurrentHash_ImmutableOpsReturnConcurrentType(t *testing.T) {
	base := func() *sets.ConcurrentHash[int] { return sets.NewConcurrentHash(1, 2, 3) }
	other := sets.NewConcurrentHash(3, 4)

	tests := map[string]func() sets.Set[int]{
		"Filter":       func() sets.Set[int] { return base().Filter(func(e int) bool { return e > 1 }) },
		"Add":          func() sets.Set[int] { return base().Add(9) },
		"AddMany":      func() sets.Set[int] { return base().AddMany(9, 10) },
		"Union":        func() sets.Set[int] { return base().Union(other) },
		"Remove":       func() sets.Set[int] { return base().Remove(1) },
		"RemoveMany":   func() sets.Set[int] { return base().RemoveMany(1, 2) },
		"Difference":   func() sets.Set[int] { return base().Difference(other) },
		"Intersection": func() sets.Set[int] { return base().Intersection(other) },
	}

	for name, op := range tests {
		t.Run(name, func(t *testing.T) {
			got := op()
			if _, ok := got.(*sets.ConcurrentHash[int]); !ok {
				t.Errorf("%s returned %T, want *sets.ConcurrentHash[int]", name, got)
			}
		})
	}
}

// TestConcurrentHashRW_ImmutableOpsReturnConcurrentType asserts the same for the
// read-write mutex variant.
func TestConcurrentHashRW_ImmutableOpsReturnConcurrentType(t *testing.T) {
	base := func() *sets.ConcurrentHashRW[int] { return sets.NewConcurrentHashRW(1, 2, 3) }
	other := sets.NewConcurrentHashRW(3, 4)

	tests := map[string]func() sets.Set[int]{
		"Filter":       func() sets.Set[int] { return base().Filter(func(e int) bool { return e > 1 }) },
		"Add":          func() sets.Set[int] { return base().Add(9) },
		"AddMany":      func() sets.Set[int] { return base().AddMany(9, 10) },
		"Union":        func() sets.Set[int] { return base().Union(other) },
		"Remove":       func() sets.Set[int] { return base().Remove(1) },
		"RemoveMany":   func() sets.Set[int] { return base().RemoveMany(1, 2) },
		"Difference":   func() sets.Set[int] { return base().Difference(other) },
		"Intersection": func() sets.Set[int] { return base().Intersection(other) },
	}

	for name, op := range tests {
		t.Run(name, func(t *testing.T) {
			got := op()
			if _, ok := got.(*sets.ConcurrentHashRW[int]); !ok {
				t.Errorf("%s returned %T, want *sets.ConcurrentHashRW[int]", name, got)
			}
		})
	}
}

// TestConcurrentHash_ImmutableOpsIndependentOfReceiver verifies that the set
// returned by an immutable op does not share state with the receiver: mutating
// the receiver in place afterwards must not change the returned set.
func TestConcurrentHash_ImmutableOpsIndependentOfReceiver(t *testing.T) {
	original := sets.NewConcurrentHash(1, 2, 3)

	added := original.Add(4)
	filtered := original.Filter(func(e int) bool { return e > 1 })

	// Mutate the receiver in place; results must be unaffected.
	original.AddInPlace(100)
	original.RemoveInPlace(1)

	if added.Contains(100) {
		t.Error("Add result was affected by later in-place mutation of the receiver")
	}
	if !added.Contains(1) {
		t.Error("Add result lost element 1 after receiver mutation")
	}
	if filtered.Contains(100) {
		t.Error("Filter result was affected by later in-place mutation of the receiver")
	}
	if filtered.Length() != 2 {
		t.Errorf("Filter result length = %d, want 2", filtered.Length())
	}
}

// TestConcurrentHashRW_ImmutableOpsIndependentOfReceiver verifies independence
// for the read-write mutex variant.
func TestConcurrentHashRW_ImmutableOpsIndependentOfReceiver(t *testing.T) {
	original := sets.NewConcurrentHashRW(1, 2, 3)

	removed := original.Remove(2)

	original.AddInPlace(100)

	if removed.Contains(100) {
		t.Error("Remove result was affected by later in-place mutation of the receiver")
	}
	if removed.Length() != 2 {
		t.Errorf("Remove result length = %d, want 2", removed.Length())
	}
}
