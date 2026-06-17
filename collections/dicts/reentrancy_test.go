package dicts_test

import (
	"testing"
	"time"

	"github.com/pickeringtech/go-collections/collections/dicts"
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

func dictReentrancyFactories() []struct {
	name string
	make func() dicts.MutableDict[string, int]
} {
	return []struct {
		name string
		make func() dicts.MutableDict[string, int]
	}{
		{"ConcurrentHash", func() dicts.MutableDict[string, int] {
			return dicts.NewConcurrentHash[string, int](dicts.Pair[string, int]{Key: "a", Value: 1}, dicts.Pair[string, int]{Key: "b", Value: 2})
		}},
		{"ConcurrentHashRW", func() dicts.MutableDict[string, int] {
			return dicts.NewConcurrentHashRW[string, int](dicts.Pair[string, int]{Key: "a", Value: 1}, dicts.Pair[string, int]{Key: "b", Value: 2})
		}},
		{"ConcurrentTree", func() dicts.MutableDict[string, int] {
			return dicts.NewConcurrentTree[string, int](dicts.Pair[string, int]{Key: "a", Value: 1}, dicts.Pair[string, int]{Key: "b", Value: 2})
		}},
		{"ConcurrentTreeRW", func() dicts.MutableDict[string, int] {
			return dicts.NewConcurrentTreeRW[string, int](dicts.Pair[string, int]{Key: "a", Value: 1}, dicts.Pair[string, int]{Key: "b", Value: 2})
		}},
	}
}

// TestConcurrentDictCallbacksAreReentrant verifies that a callback may safely
// call back into the same concurrent dictionary (including a write method,
// which on the RWMutex variants would otherwise be a read->write upgrade
// deadlock).
func TestConcurrentDictCallbacksAreReentrant(t *testing.T) {
	for _, f := range dictReentrancyFactories() {
		t.Run(f.name+"/ForEach", func(t *testing.T) {
			d := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				d.ForEach(func(k string, v int) {
					d.PutInPlace("reentry", v)
				})
			})
		})
		t.Run(f.name+"/Filter", func(t *testing.T) {
			d := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				_ = d.Filter(func(k string, v int) bool {
					d.PutInPlace("reentry", v)
					return true
				})
			})
		})
		t.Run(f.name+"/FilterInPlace", func(t *testing.T) {
			d := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				d.FilterInPlace(func(k string, v int) bool {
					d.PutInPlace("reentry", v)
					return true
				})
			})
		})
		t.Run(f.name+"/AllMatch", func(t *testing.T) {
			d := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				_ = d.AllMatch(func(k string, v int) bool {
					d.PutInPlace("reentry", v)
					return true
				})
			})
		})
		t.Run(f.name+"/Find", func(t *testing.T) {
			d := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				_, _, _ = d.Find(func(k string, v int) bool {
					d.PutInPlace("reentry", v)
					return false
				})
			})
		})
	}
}
