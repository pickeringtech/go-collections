package lru_test

import (
	"testing"
	"time"

	"github.com/pickeringtech/go-collections/collections/lru"
)

// assertNoReentrantDeadlock runs op in a goroutine and fails if it does not
// complete within the timeout. A timeout means ForEach held the lock across the
// user callback, so re-entering the cache from inside the callback deadlocked
// (see issue #75).
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

// TestConcurrentLRUForEachIsReentrant verifies that a ForEach callback may
// safely call back into the same concurrent cache (including a write method,
// which on the RWMutex variant would otherwise be a read->write upgrade
// deadlock).
func TestConcurrentLRUForEachIsReentrant(t *testing.T) {
	t.Run("ConcurrentLRU", func(t *testing.T) {
		c := lru.NewConcurrentLRU[string, int](10)
		c.PutInPlace("a", 1)
		c.PutInPlace("b", 2)
		assertNoReentrantDeadlock(t, "ConcurrentLRU", func() {
			c.ForEach(func(k string, v int) {
				c.PutInPlace("reentry", v)
			})
		})
	})
	t.Run("ConcurrentLRURW", func(t *testing.T) {
		c := lru.NewConcurrentLRURW[string, int](10)
		c.PutInPlace("a", 1)
		c.PutInPlace("b", 2)
		assertNoReentrantDeadlock(t, "ConcurrentLRURW", func() {
			c.ForEach(func(k string, v int) {
				c.PutInPlace("reentry", v)
			})
		})
	})
}
