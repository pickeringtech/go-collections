package concurrency_test

import (
	"fmt"
	"sync/atomic"

	"github.com/pickeringtech/go-collections/concurrency"
)

// ExampleNewBackgroundWorkLimiter shows the asynchronous limiter: Start opens it
// for work, Add hands items off to run in the background bounded by the limit,
// and Stop followed by Wait drains everything still in flight. The work runs
// concurrently, so the counter is updated atomically; the printed total is
// order-independent and therefore deterministic.
func ExampleNewBackgroundWorkLimiter() {
	limiter := concurrency.NewBackgroundWorkLimiter(2)
	limiter.Start()

	var counter int64
	for i := 0; i < 3; i++ {
		limiter.Add(func() error {
			atomic.AddInt64(&counter, 1)
			return nil
		})
	}

	limiter.Stop()
	limiter.Wait()

	fmt.Printf("completed %d items with %d errors\n", atomic.LoadInt64(&counter), len(limiter.Errors()))
	// Output: completed 3 items with 0 errors
}

// ExampleNewBlockingWorkLimiter shows the synchronous limiter: Run launches all
// the work bounded by the limit and blocks until every item is done, returning
// the errors collected along the way.
func ExampleNewBlockingWorkLimiter() {
	limiter := concurrency.NewBlockingWorkLimiter(2)

	var counter int64
	work := []concurrency.WorkFunc{
		func() error { atomic.AddInt64(&counter, 1); return nil },
		func() error { atomic.AddInt64(&counter, 1); return nil },
		func() error { atomic.AddInt64(&counter, 1); return nil },
	}

	errs := limiter.Run(work)

	fmt.Printf("completed %d items with %d errors\n", atomic.LoadInt64(&counter), len(errs))
	// Output: completed 3 items with 0 errors
}
