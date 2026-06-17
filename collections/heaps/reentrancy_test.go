package heaps_test

import (
	"testing"
	"time"

	"github.com/pickeringtech/go-collections/collections/heaps"
)

// assertNoReentrantDeadlock runs op in a goroutine and fails if it does not
// complete within the timeout. A timeout means ForEach held the lock across the
// user callback, so re-entering the heap from inside the callback deadlocked
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

// TestConcurrentHeapForEachIsReentrant verifies that a ForEach callback may
// safely call back into the same concurrent heap (including a write method,
// which on the RWMutex variant would otherwise be a read->write upgrade
// deadlock).
func TestConcurrentHeapForEachIsReentrant(t *testing.T) {
	t.Run("ConcurrentBinary", func(t *testing.T) {
		h := heaps.NewConcurrentMin[int](3, 1, 2)
		assertNoReentrantDeadlock(t, "ConcurrentBinary", func() {
			h.ForEach(func(v int) {
				h.PushInPlace(v + 100)
			})
		})
	})
	t.Run("ConcurrentRWBinary", func(t *testing.T) {
		h := heaps.NewConcurrentRWMin[int](3, 1, 2)
		assertNoReentrantDeadlock(t, "ConcurrentRWBinary", func() {
			h.ForEach(func(v int) {
				h.PushInPlace(v + 100)
			})
		})
	})
}
