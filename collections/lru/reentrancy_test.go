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

// TestConcurrentLRUOnEvictIsReentrant verifies that an onEvict callback fired by
// a capacity-driven eviction may safely call back into the same concurrent
// cache. Before issue #155 the callback ran while the cache lock was held, so
// re-entering the cache from it deadlocked.
func TestConcurrentLRUOnEvictIsReentrant(t *testing.T) {
	t.Run("ConcurrentLRU/PutInPlace", func(t *testing.T) {
		var c *lru.ConcurrentLRU[string, int]
		var once bool
		c = lru.NewConcurrentLRU[string, int](1, lru.WithOnEvict(func(k string, v int) {
			// Re-enter the cache: a read and a write that both need the lock. The
			// once guard stops the write from cascading, since a capacity-1 cache
			// evicts on every insert.
			c.Length()
			if !once {
				once = true
				c.PutInPlace("evicted-"+k, v)
			}
		}))
		c.PutInPlace("a", 1)
		assertNoReentrantDeadlock(t, "ConcurrentLRU/PutInPlace", func() {
			c.PutInPlace("b", 2) // evicts "a", fires onEvict
		})
	})
	t.Run("ConcurrentLRU/Put", func(t *testing.T) {
		var c *lru.ConcurrentLRU[string, int]
		c = lru.NewConcurrentLRU[string, int](1, lru.WithOnEvict(func(k string, v int) {
			c.Length()
		}))
		c.PutInPlace("a", 1)
		assertNoReentrantDeadlock(t, "ConcurrentLRU/Put", func() {
			c.Put("b", 2) // immutable put on a full cache evicts, fires onEvict
		})
	})
	t.Run("ConcurrentLRURW/PutInPlace", func(t *testing.T) {
		var c *lru.ConcurrentLRURW[string, int]
		var once bool
		c = lru.NewConcurrentLRURW[string, int](1, lru.WithOnEvict(func(k string, v int) {
			c.Length() // read lock
			if !once {
				once = true
				c.PutInPlace("evicted-"+k, v) // write lock
			}
		}))
		c.PutInPlace("a", 1)
		assertNoReentrantDeadlock(t, "ConcurrentLRURW/PutInPlace", func() {
			c.PutInPlace("b", 2)
		})
	})
	t.Run("ConcurrentLRURW/Put", func(t *testing.T) {
		var c *lru.ConcurrentLRURW[string, int]
		c = lru.NewConcurrentLRURW[string, int](1, lru.WithOnEvict(func(k string, v int) {
			c.Length() // read lock taken from within a read-locked Put
		}))
		c.PutInPlace("a", 1)
		assertNoReentrantDeadlock(t, "ConcurrentLRURW/Put", func() {
			c.Put("b", 2)
		})
	})
}
