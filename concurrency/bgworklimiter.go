package concurrency

import "sync"

// BackgroundWorkLimiter runs work functions asynchronously, allowing you to specify when the work begins, and finishes.
type BackgroundWorkLimiter struct {
	max         int
	workToDo    chan WorkFunc
	workLimiter chan struct{}
	waitGroup   *sync.WaitGroup

	errLock sync.Mutex
	errs    []error
}

// NewBackgroundWorkLimiter creates a BackgroundWorkLimiter that runs at most limit work functions concurrently.
func NewBackgroundWorkLimiter(limit int) *BackgroundWorkLimiter {
	return &BackgroundWorkLimiter{
		max:       limit,
		waitGroup: &sync.WaitGroup{},
	}
}

// Start opens the limiter for work and begins processing items added via Add in the background.
func (wl *BackgroundWorkLimiter) Start() {
	wl.workToDo = make(chan WorkFunc)
	wl.workLimiter = make(chan struct{}, wl.max)
	go wl.run()
}

// run performs each of the functions in the workToDo channel, as and when they become available. Work is arranged with
// the waitGroup to allow the User to await the completion of all work.
func (wl *BackgroundWorkLimiter) run() {
	for work := range wl.workToDo {
		wl.waitGroup.Add(1)
		work := work
		wl.workLimiter <- struct{}{}
		go func() {
			defer wl.waitGroup.Done()
			// Release worker lock when work is done
			defer func() { <-wl.workLimiter }()

			// Run the work
			err := work()

			// Collect errors
			if err != nil {
				wl.errLock.Lock()
				wl.errs = append(wl.errs, err)
				wl.errLock.Unlock()
			}
		}()
	}
}

// Stop shuts down the workToDo channel, preventing any new work from being added - but does not stop existing work
// in process being completed.
func (wl *BackgroundWorkLimiter) Stop() {
	close(wl.workToDo)
}

// Wait awaits the completion of every item in the workToDo being completed. This includes work which is still in
// process at the point Stop was called. In order to await all work being completed after a Stop was called, call Wait.
func (wl *BackgroundWorkLimiter) Wait() {
	wl.waitGroup.Wait()
}

// Add adds an item of work to be completed in the background.
func (wl *BackgroundWorkLimiter) Add(work WorkFunc) {
	wl.workToDo <- work
}

// Errors returns a copy of the errors collected from completed work. Call it
// after Wait so that every item of work has had a chance to finish.
func (wl *BackgroundWorkLimiter) Errors() []error {
	wl.errLock.Lock()
	defer wl.errLock.Unlock()
	return append([]error(nil), wl.errs...)
}
