package concurrency

import (
	"errors"
	"strings"
	"sync/atomic"
	"testing"
)

// assertPanicContains is deferred in the misuse-path tests: it recovers the
// panic and fails unless it is a string containing want. Asserting the message
// (not merely that some panic fired) keeps the test from passing on an
// unrelated panic - e.g. a nil-pointer dereference - that would mean the
// lifecycle guard had regressed.
func assertPanicContains(t *testing.T, want string) {
	t.Helper()
	r := recover()
	if r == nil {
		t.Fatalf("expected panic containing %q, got none", want)
	}
	msg, ok := r.(string)
	if !ok {
		t.Fatalf("expected string panic containing %q, got %T: %v", want, r, r)
	}
	if !strings.Contains(msg, want) {
		t.Fatalf("panic = %q, want it to contain %q", msg, want)
	}
}

func TestBackgroundWorkLimiter_RunsAllWork(t *testing.T) {
	var counter int64

	limiter := NewBackgroundWorkLimiter(4)
	limiter.Start()
	for i := 0; i < 100; i++ {
		limiter.Add(func() error {
			atomic.AddInt64(&counter, 1)
			return nil
		})
	}
	limiter.Stop()
	limiter.Wait()

	if got := atomic.LoadInt64(&counter); got != 100 {
		t.Errorf("ran %d work items, want 100", got)
	}
	if errs := limiter.Errors(); len(errs) != 0 {
		t.Errorf("Errors() = %d, want 0", len(errs))
	}
}

func TestBackgroundWorkLimiter_CollectsErrors(t *testing.T) {
	failErr := errors.New("boom")

	limiter := NewBackgroundWorkLimiter(2)
	limiter.Start()
	limiter.Add(func() error { return nil })
	limiter.Add(func() error { return failErr })
	limiter.Add(func() error { return failErr })
	limiter.Stop()
	limiter.Wait()

	errs := limiter.Errors()
	if len(errs) != 2 {
		t.Fatalf("Errors() = %d, want 2", len(errs))
	}
	for _, err := range errs {
		if !errors.Is(err, failErr) {
			t.Errorf("unexpected error %v, want %v", err, failErr)
		}
	}
}

func TestBackgroundWorkLimiter_InvalidLimitIsClamped(t *testing.T) {
	// limit == 0 once produced an unbuffered semaphore that deadlocked on the first
	// send; limit < 0 panicked at channel creation. Both must now run to completion.
	for _, limit := range []int{0, -1, -5} {
		limit := limit
		t.Run("", func(t *testing.T) {
			var counter int64

			limiter := NewBackgroundWorkLimiter(limit)
			limiter.Start()
			for i := 0; i < 10; i++ {
				limiter.Add(func() error {
					atomic.AddInt64(&counter, 1)
					return nil
				})
			}
			limiter.Stop()
			limiter.Wait()

			if got := atomic.LoadInt64(&counter); got != 10 {
				t.Errorf("ran %d work items, want 10", got)
			}
			if errs := limiter.Errors(); len(errs) != 0 {
				t.Errorf("Errors() = %d, want 0", len(errs))
			}
		})
	}
}

func TestBackgroundWorkLimiter_AddBeforeStartPanics(t *testing.T) {
	// Before Start the workToDo channel is nil; sending on it would block
	// forever. The guard must turn that into an immediate panic instead.
	limiter := NewBackgroundWorkLimiter(2)
	defer assertPanicContains(t, "Add called while not running")
	limiter.Add(func() error { return nil })
}

func TestBackgroundWorkLimiter_StopBeforeStartPanics(t *testing.T) {
	// Before Start the workToDo channel is nil; close(nil) panics obscurely.
	// The guard must produce a clear panic instead.
	limiter := NewBackgroundWorkLimiter(2)
	defer assertPanicContains(t, "Stop called before Start")
	limiter.Stop()
}

func TestBackgroundWorkLimiter_AddAfterStopPanics(t *testing.T) {
	// Once Stop has closed workToDo, a further Add would panic with "send on
	// closed channel". The guard must reject it cleanly instead.
	limiter := NewBackgroundWorkLimiter(2)
	limiter.Start()
	limiter.Stop()
	limiter.Wait()
	defer assertPanicContains(t, "Add called while not running")
	limiter.Add(func() error { return nil })
}

func TestBackgroundWorkLimiter_DoubleStartIsNoOp(t *testing.T) {
	// A second Start must not orphan the first goroutine or replace the
	// channels; the limiter must still run all work to completion exactly once.
	var counter int64

	limiter := NewBackgroundWorkLimiter(2)
	limiter.Start()
	limiter.Start()
	for i := 0; i < 10; i++ {
		limiter.Add(func() error {
			atomic.AddInt64(&counter, 1)
			return nil
		})
	}
	limiter.Stop()
	limiter.Wait()

	if got := atomic.LoadInt64(&counter); got != 10 {
		t.Errorf("ran %d work items, want 10", got)
	}
}

func TestBackgroundWorkLimiter_DoubleStopIsNoOp(t *testing.T) {
	// A second Stop must not panic by closing an already-closed channel.
	limiter := NewBackgroundWorkLimiter(2)
	limiter.Start()
	limiter.Add(func() error { return nil })
	limiter.Stop()
	limiter.Stop()
	limiter.Wait()
}

func TestBackgroundWorkLimiter_NeverExceedsLimit(t *testing.T) {
	const limit = 3
	var inFlight int64
	var maxObserved int64

	limiter := NewBackgroundWorkLimiter(limit)
	limiter.Start()
	for i := 0; i < 50; i++ {
		limiter.Add(func() error {
			current := atomic.AddInt64(&inFlight, 1)
			for {
				observed := atomic.LoadInt64(&maxObserved)
				if current <= observed || atomic.CompareAndSwapInt64(&maxObserved, observed, current) {
					break
				}
			}
			for j := 0; j < 1000; j++ {
				_ = j
			}
			atomic.AddInt64(&inFlight, -1)
			return nil
		})
	}
	limiter.Stop()
	limiter.Wait()

	if got := atomic.LoadInt64(&maxObserved); got > limit {
		t.Errorf("observed %d concurrent workers, limit is %d", got, limit)
	}
}
