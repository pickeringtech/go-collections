package lru_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/lru"
)

// TestConcurrentLRU_ImmutableOpsReturnConcurrentType asserts that the immutable
// operations on ConcurrentLRU return a new ConcurrentLRU (the same thread-safe
// type) rather than a plain, non-thread-safe LRU — thread-safe in, thread-safe out.
func TestConcurrentLRU_ImmutableOpsReturnConcurrentType(t *testing.T) {
	base := func() *lru.ConcurrentLRU[string, int] {
		return lru.NewConcurrentLRU(2, lru.WithEntries(
			lru.Pair[string, int]{Key: "a", Value: 1},
			lru.Pair[string, int]{Key: "b", Value: 2},
		))
	}

	tests := map[string]func() lru.Cache[string, int]{
		"Put":    func() lru.Cache[string, int] { return base().Put("c", 3) },
		"Remove": func() lru.Cache[string, int] { return base().Remove("a") },
	}

	for name, op := range tests {
		t.Run(name, func(t *testing.T) {
			got := op()
			_, ok := got.(*lru.ConcurrentLRU[string, int])
			if !ok {
				t.Errorf("%s returned %T, want *lru.ConcurrentLRU[string, int]", name, got)
			}
		})
	}
}

// TestConcurrentLRURW_ImmutableOpsReturnConcurrentType asserts the same for the
// read-write mutex variant.
func TestConcurrentLRURW_ImmutableOpsReturnConcurrentType(t *testing.T) {
	base := func() *lru.ConcurrentLRURW[string, int] {
		return lru.NewConcurrentLRURW(2, lru.WithEntries(
			lru.Pair[string, int]{Key: "a", Value: 1},
			lru.Pair[string, int]{Key: "b", Value: 2},
		))
	}

	tests := map[string]func() lru.Cache[string, int]{
		"Put":    func() lru.Cache[string, int] { return base().Put("c", 3) },
		"Remove": func() lru.Cache[string, int] { return base().Remove("a") },
	}

	for name, op := range tests {
		t.Run(name, func(t *testing.T) {
			got := op()
			_, ok := got.(*lru.ConcurrentLRURW[string, int])
			if !ok {
				t.Errorf("%s returned %T, want *lru.ConcurrentLRURW[string, int]", name, got)
			}
		})
	}
}

// TestConcurrentLRU_ImmutableOpsIndependentOfReceiver verifies the cache
// returned by an immutable op shares no state with the receiver.
func TestConcurrentLRU_ImmutableOpsIndependentOfReceiver(t *testing.T) {
	original := lru.NewConcurrentLRU(3, lru.WithEntries(
		lru.Pair[string, int]{Key: "a", Value: 1},
		lru.Pair[string, int]{Key: "b", Value: 2},
	))

	put := original.Put("c", 3)
	removed := original.Remove("a")

	// Mutate the receiver in place; results must be unaffected.
	original.PutInPlace("z", 99)
	original.RemoveInPlace("b")

	if put.Contains("z") {
		t.Error("Put result was affected by later in-place mutation of the receiver")
	}
	if value, ok := put.Peek("a"); !ok || value != 1 {
		t.Error("Put result lost key \"a\" after receiver mutation")
	}
	if removed.Contains("z") {
		t.Error("Remove result was affected by later in-place mutation of the receiver")
	}
	if removed.Length() != 1 {
		t.Errorf("Remove result length = %d, want 1", removed.Length())
	}
}
