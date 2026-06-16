# Worker-Pool Limiter

Bound concurrent work with a buffered-channel semaphore plus a `WaitGroup`.

```go
func (wl *BlockingWorkLimiter) Run(workToDo []WorkFunc) []error {
	var errors []error
	workLimiter := make(chan struct{}, wl.max) // semaphore sized to max workers
	var wg sync.WaitGroup
	wg.Add(len(workToDo))

	for _, work := range workToDo {
		work := work                 // capture loop var
		workLimiter <- struct{}{}    // acquire a slot (blocks at max)
		go func() {
			defer func() { <-workLimiter }() // release slot
			defer wg.Done()
			if err := work(); err != nil {
				errors = append(errors, err)
			}
		}()
	}
	wg.Wait()
	return errors
}
```

- Semaphore = `chan struct{}` buffered to `max`; send to acquire, receive (in `defer`) to release.
- Always `work := work` to capture the loop variable before the goroutine.
- `defer` both the slot release and `wg.Done()`.
- Work units are `WorkFunc` (`func() error`); collect and return all errors.
- Blocking variant waits for all work in `Run`; background variant feeds work via a channel and exposes `Stop`/`Wait`.
