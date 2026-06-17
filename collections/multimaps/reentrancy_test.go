package multimaps_test

import (
	"testing"
	"time"

	"github.com/pickeringtech/go-collections/collections/multimaps"
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

func multimapReentrancyFactories() []struct {
	name string
	make func() multimaps.MutableMultimap[string, int]
} {
	seed := []multimaps.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}}
	return []struct {
		name string
		make func() multimaps.MutableMultimap[string, int]
	}{
		{"ConcurrentListMultimap", func() multimaps.MutableMultimap[string, int] {
			return multimaps.NewConcurrentListMultimap[string, int](seed...)
		}},
		{"ConcurrentRWListMultimap", func() multimaps.MutableMultimap[string, int] {
			return multimaps.NewConcurrentRWListMultimap[string, int](seed...)
		}},
		{"ConcurrentSetMultimap", func() multimaps.MutableMultimap[string, int] {
			return multimaps.NewConcurrentSetMultimap[string, int](seed...)
		}},
		{"ConcurrentRWSetMultimap", func() multimaps.MutableMultimap[string, int] {
			return multimaps.NewConcurrentRWSetMultimap[string, int](seed...)
		}},
	}
}

// TestConcurrentMultimapCallbacksAreReentrant verifies that a callback may
// safely call back into the same concurrent multimap (including a write method,
// which on the RWMutex variants would otherwise be a read->write upgrade
// deadlock).
func TestConcurrentMultimapCallbacksAreReentrant(t *testing.T) {
	for _, f := range multimapReentrancyFactories() {
		t.Run(f.name+"/ForEach", func(t *testing.T) {
			m := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				m.ForEach(func(k string, v int) {
					m.PutInPlace("reentry", v)
				})
			})
		})
		t.Run(f.name+"/ForEachKey", func(t *testing.T) {
			m := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				m.ForEachKey(func(k string, values []int) {
					m.PutInPlace("reentry", len(values))
				})
			})
		})
		t.Run(f.name+"/Filter", func(t *testing.T) {
			m := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				_ = m.Filter(func(k string, v int) bool {
					m.PutInPlace("reentry", v)
					return true
				})
			})
		})
		t.Run(f.name+"/FilterInPlace", func(t *testing.T) {
			m := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				m.FilterInPlace(func(k string, v int) bool {
					m.PutInPlace("reentry", v)
					return true
				})
			})
		})
		t.Run(f.name+"/AllMatch", func(t *testing.T) {
			m := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				_ = m.AllMatch(func(k string, v int) bool {
					m.PutInPlace("reentry", v)
					return true
				})
			})
		})
		t.Run(f.name+"/Find", func(t *testing.T) {
			m := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				_, _, _ = m.Find(func(k string, v int) bool {
					m.PutInPlace("reentry", v)
					return false
				})
			})
		})
		t.Run(f.name+"/All", func(t *testing.T) {
			m := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				for k, v := range m.All() {
					_ = k
					m.PutInPlace("reentry", v)
				}
			})
		})
		t.Run(f.name+"/KeysSeq", func(t *testing.T) {
			m := f.make()
			assertNoReentrantDeadlock(t, f.name, func() {
				for k := range m.KeysSeq() {
					m.PutInPlace("reentry", len(k))
				}
			})
		})
	}
}
