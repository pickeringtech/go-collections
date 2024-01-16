package concurrency

type BackgroundWorkLimiter struct {
	max         int
	workToDo    chan WorkFunc
	workLimiter chan struct{}
}

func NewBackgroundWorkLimiter(max int) *BackgroundWorkLimiter {
	return &BackgroundWorkLimiter{
		max: max,
	}
}

func (wl *BackgroundWorkLimiter) Start() {
	wl.workToDo = make(chan WorkFunc)
	wl.workLimiter = make(chan struct{}, wl.max)
	go wl.run()
}

func (wl *BackgroundWorkLimiter) run() []error {
	var errors []error

	for work := range wl.workToDo {
		work := work
		wl.workLimiter <- struct{}{}
		go func() {
			// Release worker lock when work is done
			defer func() { <-wl.workLimiter }()

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

func (wl *BackgroundWorkLimiter) Stop() {
	close(wl.workToDo)
}

func (wl *BackgroundWorkLimiter) Add(work WorkFunc) {
	wl.workToDo <- work
}
