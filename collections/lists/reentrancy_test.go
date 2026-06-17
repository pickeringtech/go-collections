package lists_test

import (
	"testing"
	"time"

	"github.com/pickeringtech/go-collections/collections/lists"
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

func listReentrancyFactories() []struct {
	name string
	make func() lists.MutableList[int]
} {
	return []struct {
		name string
		make func() lists.MutableList[int]
	}{
		{"ConcurrentArray", func() lists.MutableList[int] { return lists.NewConcurrentArray[int](1, 2, 3) }},
		{"ConcurrentRWArray", func() lists.MutableList[int] { return lists.NewConcurrentRWArray[int](1, 2, 3) }},
		{"ConcurrentLinked", func() lists.MutableList[int] { return lists.NewConcurrentLinked[int](1, 2, 3) }},
		{"ConcurrentRWLinked", func() lists.MutableList[int] { return lists.NewConcurrentRWLinked[int](1, 2, 3) }},
		{"ConcurrentDoublyLinked", func() lists.MutableList[int] { return lists.NewConcurrentDoublyLinked[int](1, 2, 3) }},
		{"ConcurrentRWDoublyLinked", func() lists.MutableList[int] { return lists.NewConcurrentRWDoublyLinked[int](1, 2, 3) }},
	}
}

// TestConcurrentListCallbacksAreReentrant verifies that a callback may safely
// call back into the same concurrent list (including a write method, which on
// the RWMutex variants would otherwise be a read->write upgrade deadlock).
func TestConcurrentListCallbacksAreReentrant(t *testing.T) {
	for _, f := range listReentrancyFactories() {
		t.Run(f.name+"/ForEach", func(t *testing.T) {
			l := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				l.ForEach(func(v int) {
					l.PushInPlace(v)
				})
			})
		})
		t.Run(f.name+"/Filter", func(t *testing.T) {
			l := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				_ = l.Filter(func(v int) bool {
					l.PushInPlace(v)
					return true
				})
			})
		})
		t.Run(f.name+"/FilterInPlace", func(t *testing.T) {
			l := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				l.FilterInPlace(func(v int) bool {
					l.PushInPlace(v)
					return true
				})
			})
		})
		t.Run(f.name+"/AllMatch", func(t *testing.T) {
			l := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				_ = l.AllMatch(func(v int) bool {
					l.PushInPlace(v)
					return true
				})
			})
		})
		t.Run(f.name+"/Find", func(t *testing.T) {
			l := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				_, _ = l.Find(func(v int) bool {
					l.PushInPlace(v)
					return false
				})
			})
		})
	}
}
