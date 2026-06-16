package dicts_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// TestConcurrentHash_ImmutableOpsReturnConcurrentType asserts that every
// immutable operation on ConcurrentHash returns a new ConcurrentHash (the same
// thread-safe type) rather than a plain, non-thread-safe Hash.
func TestConcurrentHash_ImmutableOpsReturnConcurrentType(t *testing.T) {
	base := func() *dicts.ConcurrentHash[string, int] {
		return dicts.NewConcurrentHash(
			dicts.Pair[string, int]{Key: "a", Value: 1},
			dicts.Pair[string, int]{Key: "b", Value: 2},
		)
	}

	tests := map[string]func() dicts.Dict[string, int]{
		"Filter":     func() dicts.Dict[string, int] { return base().Filter(func(k string, v int) bool { return v > 1 }) },
		"Put":        func() dicts.Dict[string, int] { return base().Put("c", 3) },
		"PutMany":    func() dicts.Dict[string, int] { return base().PutMany(dicts.Pair[string, int]{Key: "c", Value: 3}) },
		"Remove":     func() dicts.Dict[string, int] { return base().Remove("a") },
		"RemoveMany": func() dicts.Dict[string, int] { return base().RemoveMany("a", "b") },
	}

	for name, op := range tests {
		t.Run(name, func(t *testing.T) {
			got := op()
			if _, ok := got.(*dicts.ConcurrentHash[string, int]); !ok {
				t.Errorf("%s returned %T, want *dicts.ConcurrentHash[string, int]", name, got)
			}
		})
	}
}

// TestConcurrentHashRW_ImmutableOpsReturnConcurrentType asserts the same for the
// read-write mutex variant.
func TestConcurrentHashRW_ImmutableOpsReturnConcurrentType(t *testing.T) {
	base := func() *dicts.ConcurrentHashRW[string, int] {
		return dicts.NewConcurrentHashRW(
			dicts.Pair[string, int]{Key: "a", Value: 1},
			dicts.Pair[string, int]{Key: "b", Value: 2},
		)
	}

	tests := map[string]func() dicts.Dict[string, int]{
		"Filter":     func() dicts.Dict[string, int] { return base().Filter(func(k string, v int) bool { return v > 1 }) },
		"Put":        func() dicts.Dict[string, int] { return base().Put("c", 3) },
		"PutMany":    func() dicts.Dict[string, int] { return base().PutMany(dicts.Pair[string, int]{Key: "c", Value: 3}) },
		"Remove":     func() dicts.Dict[string, int] { return base().Remove("a") },
		"RemoveMany": func() dicts.Dict[string, int] { return base().RemoveMany("a", "b") },
	}

	for name, op := range tests {
		t.Run(name, func(t *testing.T) {
			got := op()
			if _, ok := got.(*dicts.ConcurrentHashRW[string, int]); !ok {
				t.Errorf("%s returned %T, want *dicts.ConcurrentHashRW[string, int]", name, got)
			}
		})
	}
}

// TestConcurrentHash_ImmutableOpsIndependentOfReceiver verifies that the
// dictionary returned by an immutable op does not share state with the receiver.
func TestConcurrentHash_ImmutableOpsIndependentOfReceiver(t *testing.T) {
	original := dicts.NewConcurrentHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)

	put := original.Put("c", 3)
	removed := original.Remove("a")

	// Mutate the receiver in place; results must be unaffected.
	original.PutInPlace("z", 99)
	original.RemoveInPlace("b")

	if _, ok := put.Get("z", -1); ok {
		t.Error("Put result was affected by later in-place mutation of the receiver")
	}
	if v, ok := put.Get("a", -1); !ok || v != 1 {
		t.Error("Put result lost key \"a\" after receiver mutation")
	}
	if _, ok := removed.Get("z", -1); ok {
		t.Error("Remove result was affected by later in-place mutation of the receiver")
	}
	if removed.Length() != 1 {
		t.Errorf("Remove result length = %d, want 1", removed.Length())
	}
}
