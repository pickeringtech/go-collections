package sets_test

import (
	"testing"
	"time"

	"github.com/pickeringtech/go-collections/collections/sets"
)

// concurrentSetFactory builds a concurrent MutableSet[int] from the given
// elements. Each of the four concurrent set implementations is exercised so the
// self-operation regression tests cover every locking strategy.
type concurrentSetFactory struct {
	name string
	make func(elements ...int) sets.MutableSet[int]
}

var concurrentSetFactories = []concurrentSetFactory{
	{"ConcurrentHash", func(e ...int) sets.MutableSet[int] { return sets.NewConcurrentHash(e...) }},
	{"ConcurrentHashRW", func(e ...int) sets.MutableSet[int] { return sets.NewConcurrentHashRW(e...) }},
	{"ConcurrentTreeSet", func(e ...int) sets.MutableSet[int] { return sets.NewConcurrentTreeSet(e...) }},
	{"ConcurrentTreeSetRW", func(e ...int) sets.MutableSet[int] { return sets.NewConcurrentTreeSetRW(e...) }},
}

// runWithDeadlockGuard runs fn in a goroutine and fails the test if it does not
// finish promptly. A self-aliased algebra method that re-acquires the receiver's
// lock would hang forever, so the timeout turns that deadlock into a test
// failure rather than a stuck `go test` run (see issue #81).
func runWithDeadlockGuard(t *testing.T, fn func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("operation did not complete within 5s — likely self-deadlock")
	}
}

// TestConcurrentSet_SelfAlgebra_NoDeadlock verifies that every concurrent set
// algebra method is safe — and correct — when called with the receiver as the
// argument. Run under -race to also surface data races. Regression for issue #81.
func TestConcurrentSet_SelfAlgebra_NoDeadlock(t *testing.T) {
	for _, f := range concurrentSetFactories {
		f := f
		t.Run(f.name, func(t *testing.T) {
			t.Run("Union", func(t *testing.T) {
				s := f.make(1, 2, 3)
				var got sets.Set[int]
				runWithDeadlockGuard(t, func() { got = s.Union(s) })
				if !intSlicesEqual(got.AsSlice(), []int{1, 2, 3}) {
					t.Errorf("Union(self) = %v, want [1 2 3]", got.AsSlice())
				}
			})

			t.Run("UnionInPlace", func(t *testing.T) {
				s := f.make(1, 2, 3)
				runWithDeadlockGuard(t, func() { s.UnionInPlace(s) })
				if !intSlicesEqual(s.AsSlice(), []int{1, 2, 3}) {
					t.Errorf("after UnionInPlace(self) = %v, want [1 2 3]", s.AsSlice())
				}
			})

			t.Run("Difference", func(t *testing.T) {
				s := f.make(1, 2, 3)
				var got sets.Set[int]
				runWithDeadlockGuard(t, func() { got = s.Difference(s) })
				if got.Length() != 0 {
					t.Errorf("Difference(self) = %v, want empty", got.AsSlice())
				}
			})

			t.Run("DifferenceInPlace", func(t *testing.T) {
				s := f.make(1, 2, 3)
				runWithDeadlockGuard(t, func() { s.DifferenceInPlace(s) })
				if s.Length() != 0 {
					t.Errorf("after DifferenceInPlace(self) = %v, want empty", s.AsSlice())
				}
			})

			t.Run("Intersection", func(t *testing.T) {
				s := f.make(1, 2, 3)
				var got sets.Set[int]
				runWithDeadlockGuard(t, func() { got = s.Intersection(s) })
				if !intSlicesEqual(got.AsSlice(), []int{1, 2, 3}) {
					t.Errorf("Intersection(self) = %v, want [1 2 3]", got.AsSlice())
				}
			})

			t.Run("IntersectionInPlace", func(t *testing.T) {
				s := f.make(1, 2, 3)
				runWithDeadlockGuard(t, func() { s.IntersectionInPlace(s) })
				if !intSlicesEqual(s.AsSlice(), []int{1, 2, 3}) {
					t.Errorf("after IntersectionInPlace(self) = %v, want [1 2 3]", s.AsSlice())
				}
			})

			t.Run("IsSubsetOf", func(t *testing.T) {
				s := f.make(1, 2, 3)
				var got bool
				runWithDeadlockGuard(t, func() { got = s.IsSubsetOf(s) })
				if !got {
					t.Error("IsSubsetOf(self) = false, want true")
				}
			})

			t.Run("IsSupersetOf", func(t *testing.T) {
				s := f.make(1, 2, 3)
				var got bool
				runWithDeadlockGuard(t, func() { got = s.IsSupersetOf(s) })
				if !got {
					t.Error("IsSupersetOf(self) = false, want true")
				}
			})

			t.Run("IsDisjoint_NonEmpty", func(t *testing.T) {
				s := f.make(1, 2, 3)
				var got bool
				runWithDeadlockGuard(t, func() { got = s.IsDisjoint(s) })
				if got {
					t.Error("IsDisjoint(self) = true for a non-empty set, want false")
				}
			})

			t.Run("IsDisjoint_Empty", func(t *testing.T) {
				s := f.make()
				var got bool
				runWithDeadlockGuard(t, func() { got = s.IsDisjoint(s) })
				if !got {
					t.Error("IsDisjoint(self) = false for an empty set, want true")
				}
			})

			t.Run("Equals", func(t *testing.T) {
				s := f.make(1, 2, 3)
				var got bool
				runWithDeadlockGuard(t, func() { got = s.Equals(s) })
				if !got {
					t.Error("Equals(self) = false, want true")
				}
			})
		})
	}
}
