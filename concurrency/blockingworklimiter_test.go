package concurrency

import (
	"errors"
	"sync/atomic"
	"testing"
)

func TestBlockingWorkLimiter_RunsAllWork(t *testing.T) {
	var counter int64
	work := make([]WorkFunc, 100)
	for i := range work {
		work[i] = func() error {
			atomic.AddInt64(&counter, 1)
			return nil
		}
	}

	limiter := NewBlockingWorkLimiter(4)
	errs := limiter.Run(work)

	if len(errs) != 0 {
		t.Errorf("Run() returned %d errors, want 0", len(errs))
	}
	if got := atomic.LoadInt64(&counter); got != 100 {
		t.Errorf("ran %d work items, want 100", got)
	}
}

func TestBlockingWorkLimiter_CollectsErrors(t *testing.T) {
	failErr := errors.New("boom")
	work := []WorkFunc{
		func() error { return nil },
		func() error { return failErr },
		func() error { return failErr },
		func() error { return nil },
	}

	limiter := NewBlockingWorkLimiter(2)
	errs := limiter.Run(work)

	if len(errs) != 2 {
		t.Fatalf("Run() returned %d errors, want 2", len(errs))
	}
	for _, err := range errs {
		if !errors.Is(err, failErr) {
			t.Errorf("unexpected error %v, want %v", err, failErr)
		}
	}
}

func TestBlockingWorkLimiter_NeverExceedsLimit(t *testing.T) {
	const limit = 3
	var inFlight int64
	var maxObserved int64

	work := make([]WorkFunc, 50)
	for i := range work {
		work[i] = func() error {
			current := atomic.AddInt64(&inFlight, 1)
			for {
				observed := atomic.LoadInt64(&maxObserved)
				if current <= observed || atomic.CompareAndSwapInt64(&maxObserved, observed, current) {
					break
				}
			}
			// Brief busy work so multiple items genuinely overlap.
			for j := 0; j < 1000; j++ {
				_ = j
			}
			atomic.AddInt64(&inFlight, -1)
			return nil
		}
	}

	limiter := NewBlockingWorkLimiter(limit)
	limiter.Run(work)

	if got := atomic.LoadInt64(&maxObserved); got > limit {
		t.Errorf("observed %d concurrent workers, limit is %d", got, limit)
	}
}

func TestBlockingWorkLimiter_InvalidLimitIsClamped(t *testing.T) {
	// limit == 0 once produced an unbuffered semaphore that deadlocked on the first
	// send; limit < 0 panicked at channel creation. Both must now run to completion.
	for _, limit := range []int{0, -1, -5} {
		limit := limit
		t.Run("", func(t *testing.T) {
			var counter int64
			work := make([]WorkFunc, 10)
			for i := range work {
				work[i] = func() error {
					atomic.AddInt64(&counter, 1)
					return nil
				}
			}

			limiter := NewBlockingWorkLimiter(limit)
			errs := limiter.Run(work)

			if len(errs) != 0 {
				t.Errorf("Run() returned %d errors, want 0", len(errs))
			}
			if got := atomic.LoadInt64(&counter); got != 10 {
				t.Errorf("ran %d work items, want 10", got)
			}
		})
	}
}

func TestBlockingWorkLimiter_EmptyWork(t *testing.T) {
	limiter := NewBlockingWorkLimiter(2)
	errs := limiter.Run(nil)
	if len(errs) != 0 {
		t.Errorf("Run(nil) returned %d errors, want 0", len(errs))
	}
}
