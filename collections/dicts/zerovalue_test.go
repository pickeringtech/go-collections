package dicts_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// The concurrent dictionaries embed their mutex by value, so the lock on a
// zero-value instance is safe to take even though the backing data is nil. These
// tests pin that contract: a read on a bare &ConcurrentX{} must return empty
// rather than panicking on a nil lock. (Writes still require the constructor, as
// documented on each type.)

func TestConcurrentHash_ZeroValueLockSafeReads(t *testing.T) {
	var ch dicts.ConcurrentHash[string, int]

	if got := ch.Length(); got != 0 {
		t.Errorf("Length() = %d, want 0", got)
	}
	if ch.Contains("missing") {
		t.Error("Contains() = true, want false")
	}
	if _, ok := ch.Get("missing", -1); ok {
		t.Error("Get() ok = true, want false")
	}
	if !ch.IsEmpty() {
		t.Error("IsEmpty() = false, want true")
	}
}

func TestConcurrentHashRW_ZeroValueLockSafeReads(t *testing.T) {
	var ch dicts.ConcurrentHashRW[string, int]

	if got := ch.Length(); got != 0 {
		t.Errorf("Length() = %d, want 0", got)
	}
	if ch.Contains("missing") {
		t.Error("Contains() = true, want false")
	}
	if _, ok := ch.Get("missing", -1); ok {
		t.Error("Get() ok = true, want false")
	}
	if !ch.IsEmpty() {
		t.Error("IsEmpty() = false, want true")
	}
}

// Tree documents a usable zero value: a bare &Tree{} is a valid, empty tree that
// grows on the first write.
func TestTree_ZeroValueReadyForUse(t *testing.T) {
	var tr dicts.Tree[string, int]

	tr.PutInPlace("a", 1)
	tr.PutInPlace("b", 2)

	if got := tr.Length(); got != 2 {
		t.Fatalf("Length() = %d, want 2", got)
	}
	if v, ok := tr.Get("a", -1); !ok || v != 1 {
		t.Errorf("Get(a) = (%d, %t), want (1, true)", v, ok)
	}
}
