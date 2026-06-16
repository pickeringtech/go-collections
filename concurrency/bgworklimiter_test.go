package concurrency

import (
	"errors"
	"sync/atomic"
	"testing"
)

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
