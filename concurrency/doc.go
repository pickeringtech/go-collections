// Package concurrency provides helpers for running work functions concurrently
// while bounding the number of workers that may be in flight at once.
//
// BlockingWorkLimiter.Run takes the full slice of work up front and blocks until
// it is all done. BackgroundWorkLimiter instead streams work in over time and
// must be driven through its lifecycle in order:
//
//	wl := NewBackgroundWorkLimiter(4)
//	wl.Start()        // open the limiter and start the background worker
//	wl.Add(work)      // hand off work; repeat as many times as needed
//	wl.Stop()         // close the limiter to new work
//	wl.Wait()         // block until all added work has finished
//	errs := wl.Errors()
//
// Calling Add or Stop before Start, or Add after Stop, is a programming error
// and panics rather than blocking forever or panicking obscurely on a nil or
// closed channel. A redundant Start or Stop is a no-op.
package concurrency
