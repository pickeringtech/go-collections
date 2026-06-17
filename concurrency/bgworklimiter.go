package concurrency

import "sync"

// bgwState is the lifecycle phase of a BackgroundWorkLimiter. The limiter moves
// strictly forward through these phases, guarding against out-of-order calls
// that would otherwise block forever, panic, or leak goroutines.
type bgwState int

const (
	// bgwNew is the state immediately after construction, before Start.
	bgwNew bgwState = iota
	// bgwStarted is the state after Start: Add and Stop are now legal.
	bgwStarted
	// bgwStopped is the state after Stop: no further work may be added.
	bgwStopped
)

// BackgroundWorkLimiter runs work functions asynchronously, allowing you to specify when the work begins, and finishes.
//
// The required call order is Start, then any number of Add, then Stop, then
// Wait. Calling Add or Stop before Start panics, and Add after Stop panics,
// because those are programming mistakes rather than recoverable conditions. A
// second Start is a no-op (it does not leak the running goroutine), as is a
// second Stop.
type BackgroundWorkLimiter struct {
	max         int
	workToDo    chan WorkFunc
	workLimiter chan struct{}
	waitGroup   *sync.WaitGroup

	// stateLock guards state and the lifecycle transitions of the channels
	// above. Add holds it across the hand-off so a concurrent Stop cannot close
	// workToDo out from under an in-flight send.
	stateLock sync.Mutex
	state     bgwState

	errLock sync.Mutex
	errs    []error
}

// NewBackgroundWorkLimiter creates a BackgroundWorkLimiter that runs at most limit work functions concurrently.
// A limit below 1 is clamped to 1: an unbuffered semaphore (limit == 0) would deadlock on the first send,
// and a negative size would panic at channel creation.
func NewBackgroundWorkLimiter(limit int) *BackgroundWorkLimiter {
	if limit < 1 {
		limit = 1
	}
	return &BackgroundWorkLimiter{
		max:       limit,
		waitGroup: &sync.WaitGroup{},
	}
}

// Start opens the limiter for work and begins processing items added via Add in
// the background. Calling Start a second time is a no-op: the already-running
// goroutine and its channels are left untouched rather than orphaned. Start
// after Stop is likewise a no-op - a stopped limiter cannot be restarted, so
// construct a fresh BackgroundWorkLimiter for a new batch of work.
func (wl *BackgroundWorkLimiter) Start() {
	wl.stateLock.Lock()
	defer wl.stateLock.Unlock()
	if wl.state != bgwNew {
		return
	}
	wl.state = bgwStarted
	wl.workToDo = make(chan WorkFunc)
	wl.workLimiter = make(chan struct{}, wl.max)
	go wl.run()
}

// run performs each of the functions in the workToDo channel, as and when they become available. Work is arranged with
// the waitGroup to allow the User to await the completion of all work.
func (wl *BackgroundWorkLimiter) run() {
	for work := range wl.workToDo {
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
// in process being completed. Calling Stop before Start panics; a second Stop is
// a no-op.
func (wl *BackgroundWorkLimiter) Stop() {
	wl.stateLock.Lock()
	defer wl.stateLock.Unlock()
	switch wl.state {
	case bgwNew:
		panic("concurrency: BackgroundWorkLimiter.Stop called before Start")
	case bgwStopped:
		return
	}
	wl.state = bgwStopped
	close(wl.workToDo)
}

// Wait awaits the completion of every item in the workToDo being completed. This includes work which is still in
// process at the point Stop was called. In order to await all work being completed after a Stop was called, call Wait.
func (wl *BackgroundWorkLimiter) Wait() {
	wl.waitGroup.Wait()
}

// Add adds an item of work to be completed in the background. It panics if the
// limiter is not currently running - that is, if Start has not been called or
// Stop already has been called - because sending on the unallocated or closed
// channel would otherwise block forever or panic obscurely.
//
// The wait group is incremented here, on the producer side, before the work is
// handed off, so that a subsequent Stop followed by Wait cannot race ahead of
// run registering it. The state lock is held across the hand-off so a concurrent
// Stop cannot close workToDo mid-send.
func (wl *BackgroundWorkLimiter) Add(work WorkFunc) {
	wl.stateLock.Lock()
	defer wl.stateLock.Unlock()
	if wl.state != bgwStarted {
		panic("concurrency: BackgroundWorkLimiter.Add called while not running (call Start first, and do not Add after Stop)")
	}
	wl.waitGroup.Add(1)
	wl.workToDo <- work
}

// Errors returns a copy of the errors collected from completed work. Call it
// after Wait so that every item of work has had a chance to finish.
func (wl *BackgroundWorkLimiter) Errors() []error {
	wl.errLock.Lock()
	defer wl.errLock.Unlock()
	return append([]error(nil), wl.errs...)
}
