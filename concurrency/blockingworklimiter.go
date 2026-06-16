package concurrency

import "sync"

// BlockingWorkLimiter runs work functions asynchronously with a limit on the number of concurrent workers, blocking
// until all work assigned is done.
type BlockingWorkLimiter struct {
	max int
}

// NewBlockingWorkLimiter creates a BlockingWorkLimiter that runs at most limit work functions concurrently.
func NewBlockingWorkLimiter(limit int) *BlockingWorkLimiter {
	return &BlockingWorkLimiter{
		max: limit,
	}
}

// Run runs the work functions asynchronously with a limit on the number of concurrent workers.
// If the limit is reached, the function will wait until a worker is available before assigning the next work.
// The function will block until all work is done, and then return a list of the errors encountered.
func (wl *BlockingWorkLimiter) Run(workToDo []WorkFunc) []error {
	var errors []error
	var errLock sync.Mutex

	workLimiter := make(chan struct{}, wl.max)
	var wg sync.WaitGroup

	wg.Add(len(workToDo))

	for _, work := range workToDo {
		work := work
		// Acquire a worker slot before launching; this blocks once max workers
		// are in flight, which is what bounds the concurrency.
		workLimiter <- struct{}{}
		go func() {
			defer wg.Done()
			// Release worker slot when work is done
			defer func() { <-workLimiter }()

			// Run the work
			err := work()

			// Collect errors
			if err != nil {
				errLock.Lock()
				errors = append(errors, err)
				errLock.Unlock()
			}
		}()
	}

	wg.Wait()

	return errors
}
