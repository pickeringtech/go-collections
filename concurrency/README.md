# Concurrency

The `concurrency` package bounds the number of work functions running at once.
Use a work limiter for individual tasks, or a parallel helper for slices.
See the [package API docs](https://pkg.go.dev/github.com/pickeringtech/go-collections/concurrency)
for the full reference.

## Choose an API

| API | Use when | Completion and results |
| --- | --- | --- |
| `NewBlockingWorkLimiter(limit)` | All tasks are ready as a `[]WorkFunc`. | `Run(work)` blocks until all tasks finish and returns their errors. |
| `NewBackgroundWorkLimiter(limit)` | Tasks arrive over time. | `Start` → `Add` → `Stop` → `Wait` → `Errors`; `Add` can block while handing off work. |
| `Map` | Transform each slice element. | Returns a result slice in input order and an error. |
| `ForEach` | Perform a side effect for each slice element. | Waits for running callbacks and returns an error. |
| `Batch` | Process consecutive chunks of a slice. | Runs chunk callbacks concurrently and returns an error. |

Callbacks run concurrently; synchronise access to shared state. `Batch` passes
views into the input slice, so copy a batch before modifying its elements if
the input must stay unchanged.

## Bounded Map

This runnable snippet mirrors `ExampleMap` in
[`parallel_example_test.go`](parallel_example_test.go):

```go
package main

import (
	"context"
	"fmt"

	"github.com/pickeringtech/go-collections/concurrency"
)

func main() {
	input := []int{1, 2, 3, 4, 5}

	squares, err := concurrency.Map(context.Background(), input,
		func(_ context.Context, n int) (int, error) {
			return n * n, nil
		},
		concurrency.WithConcurrency(3),
	)

	fmt.Println(squares, err)
	// Output: [1 4 9 16 25] <nil>
}
```

At most three callbacks run at once. Results retain input order regardless of
completion order. The result always has `len(input)` elements; failed or skipped
positions hold the output type's zero value, so check the error before using it.

`Map`, `ForEach` and `Batch` accept `WithConcurrency` (default:
`runtime.GOMAXPROCS(0)`, values below 1 are clamped to 1) and `WithErrorPolicy`:

- `StopOnError` (default) cancels remaining work on failure, waits for callbacks
  already running, and reports the first error in input order.
- `CollectErrors` runs all items and joins their errors with `errors.Join`.
- `ContinueOnError` runs all items and suppresses work errors.

Caller context cancellation takes precedence over every error policy and returns
`ctx.Err()`. Callbacks receive a context and should observe cancellation to stop
promptly; running callbacks must return before the helper returns. The work
limiters use `WorkFunc` (`func() error`) and collect errors without these options
or built-in context cancellation.

## Background limiter lifecycle

Call **Start → Add → Stop → Wait → Errors** in that order. `Stop` closes
submissions; `Wait` waits for all accepted work to finish. This snippet mirrors
`ExampleNewBackgroundWorkLimiter` in [`examples_test.go`](examples_test.go):

```go
package main

import (
	"fmt"
	"sync/atomic"

	"github.com/pickeringtech/go-collections/concurrency"
)

func main() {
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
```

`Add` or `Stop` before `Start`, and `Add` after `Stop`, panic. Repeated `Start`
or `Stop` calls are no-ops; a stopped limiter cannot be restarted.

For a complete cross-package flow, see the existing
[`worker-pipeline` example](../examples/cmd/worker-pipeline) and its
[run instructions](../examples/README.md#running-an-app).
