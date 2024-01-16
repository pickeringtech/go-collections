package concurrency

type ForegroundWorkLimiter struct {
	max int
}

func NewForegroundWorkLimiter(max int) *ForegroundWorkLimiter {
	return &ForegroundWorkLimiter{
		max: max,
	}
}

func (wl *ForegroundWorkLimiter) Run(workToDo []WorkFunc) []error {
	var errors []error

	workLimiter := make(chan struct{}, wl.max)

	for _, work := range workToDo {
		work := work
		workLimiter <- struct{}{}
		go func() {
			// Release worker lock when work is done
			defer func() { <-workLimiter }()

			// Run the work
			err := work()

			// Collect errors
			if err != nil {
				errors = append(errors, err)
			}
		}()
	}

	return errors
}
