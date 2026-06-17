package sets_test

import (
	"testing"
	"time"

	"github.com/pickeringtech/go-collections/collections/sets"
)

// assertNoReentrantDeadlock runs op in a goroutine and fails if it does not
// complete within the timeout. A timeout means a callback-taking method held
// the lock across the user callback, so re-entering the collection from inside
// the callback deadlocked (see issue #75).
func assertNoReentrantDeadlock(t *testing.T, name string, op func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		op()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s: did not complete within 2s — the lock is held across the callback (reentrancy deadlock)", name)
	}
}

func setReentrancyFactories() []struct {
	name string
	make func() sets.MutableSet[int]
} {
	return []struct {
		name string
		make func() sets.MutableSet[int]
	}{
		{"ConcurrentHash", func() sets.MutableSet[int] { return sets.NewConcurrentHash[int](1, 2, 3) }},
		{"ConcurrentHashRW", func() sets.MutableSet[int] { return sets.NewConcurrentHashRW[int](1, 2, 3) }},
		{"ConcurrentTreeSet", func() sets.MutableSet[int] { return sets.NewConcurrentTreeSet[int](1, 2, 3) }},
		{"ConcurrentTreeSetRW", func() sets.MutableSet[int] { return sets.NewConcurrentTreeSetRW[int](1, 2, 3) }},
	}
}

// TestConcurrentSetCallbacksAreReentrant verifies that a callback may safely
// call back into the same concurrent set (including a write method, which on
// the RWMutex variants would otherwise be a read->write upgrade deadlock).
func TestConcurrentSetCallbacksAreReentrant(t *testing.T) {
	for _, f := range setReentrancyFactories() {
		t.Run(f.name+"/ForEach", func(t *testing.T) {
			s := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				s.ForEach(func(v int) {
					s.AddInPlace(v + 100)
				})
			})
		})
		t.Run(f.name+"/Filter", func(t *testing.T) {
			s := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				_ = s.Filter(func(v int) bool {
					s.AddInPlace(v + 100)
					return true
				})
			})
		})
		t.Run(f.name+"/FilterInPlace", func(t *testing.T) {
			s := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				s.FilterInPlace(func(v int) bool {
					s.AddInPlace(v + 100)
					return true
				})
			})
		})
		t.Run(f.name+"/AllMatch", func(t *testing.T) {
			s := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				_ = s.AllMatch(func(v int) bool {
					s.AddInPlace(v + 100)
					return true
				})
			})
		})
		t.Run(f.name+"/Find", func(t *testing.T) {
			s := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				_, _ = s.Find(func(v int) bool {
					s.AddInPlace(v + 100)
					return false
				})
			})
		})
	}
}
