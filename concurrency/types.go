package concurrency

// WorkFunc is a unit of work to be run by a work limiter; it returns an error if the work fails.
type WorkFunc func() error
