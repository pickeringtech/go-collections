package concurrency_test

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/pickeringtech/go-collections/concurrency"
)

// ExampleMap shows the order-preserving concurrent map: work may finish in any
// order, but output[i] always corresponds to input[i], so the printed slice is
// deterministic.
func ExampleMap() {
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

// ExampleForEach shows side-effecting parallel iteration. Because the work runs
// concurrently, shared state is updated atomically; the total is
// order-independent and therefore deterministic.
func ExampleForEach() {
	input := []int{1, 2, 3, 4, 5}

	var sum int64
	err := concurrency.ForEach(context.Background(), input,
		func(_ context.Context, n int) error {
			atomic.AddInt64(&sum, int64(n))
			return nil
		},
	)

	fmt.Printf("sum=%d err=%v\n", atomic.LoadInt64(&sum), err)
	// Output: sum=15 err=<nil>
}

// ExampleBatch shows chunked concurrent processing: the input is split into
// consecutive batches of at most size elements, and each batch is processed
// concurrently. Here every batch is summed into a shared atomic total.
func ExampleBatch() {
	input := []int{1, 2, 3, 4, 5, 6, 7}

	var total int64
	err := concurrency.Batch(context.Background(), input, 3,
		func(_ context.Context, batch []int) error {
			var sub int64
			for _, n := range batch {
				sub += int64(n)
			}
			atomic.AddInt64(&total, sub)
			return nil
		},
	)

	fmt.Printf("total=%d err=%v\n", atomic.LoadInt64(&total), err)
	// Output: total=28 err=<nil>
}

// ExampleMap_collectErrors shows the CollectErrors policy: every item runs to
// completion and all failures are reported together, so errors.Is can match any
// of them.
func ExampleMap_collectErrors() {
	input := []int{1, 2, 3, 4}

	_, err := concurrency.Map(context.Background(), input,
		func(_ context.Context, n int) (int, error) {
			if n%2 == 0 {
				return 0, fmt.Errorf("cannot process even number %d", n)
			}
			return n, nil
		},
		concurrency.WithErrorPolicy(concurrency.CollectErrors),
	)

	fmt.Println(err != nil)
	// Output: true
}
