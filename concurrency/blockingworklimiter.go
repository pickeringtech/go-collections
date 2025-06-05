package concurrency

import "sync"

// BlockingWorkLimiter runs work functions asynchronously with a limit on the number of concurrent workers, blocking
// until all work assigned is done.
type BlockingWorkLimiter struct {
	max int
}

func NewBlockingWorkLimiter(max int) *BlockingWorkLimiter {
	return &BlockingWorkLimiter{
		max: max,
	}
}

// Run runs the work functions asynchronously with a limit on the number of concurrent workers.
// If the limit is reached, the function will wait until a worker is available before assigning the next work.
// The function will block until all work is done, and then return a list of the errors encountered.
func (wl *BlockingWorkLimiter) Run(workToDo []WorkFunc) []error {
	var errors []error

	workLimiter := make(chan struct{}, wl.max)
	var wg sync.WaitGroup

	wg.Add(len(workToDo))

	for _, work := range workToDo {
		work := work
		go func() {
			// Release worker lock when work is done
			defer func() { <-workLimiter }()
			defer wg.Done()

			// Run the work
			err := work()

			// Collect errors
			if err != nil {
				errors = append(errors, err)
			}
		}()
	}

	wg.Wait()

	return errors
}
