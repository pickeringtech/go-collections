// Command worker-pipeline processes a stream of integers through a bounded
// worker pool and then folds the results back together.
//
// It demonstrates channels + concurrency working together:
//   - channels.FromSlice turns the inputs into a stream to fan out from;
//   - concurrency.BlockingWorkLimiter caps how many transforms run at once
//     (the fan-out), with results collected under a mutex (the fan-in);
//   - a channels reduce pipeline sums the fanned-in results.
//
// Worker scheduling is non-deterministic, so the squares are sorted before
// printing and the sum is order-independent — keeping the output reproducible.
package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/pickeringtech/go-collections/channels"
	"github.com/pickeringtech/go-collections/concurrency"
	"github.com/pickeringtech/go-collections/slices"
)

func main() {
	n := flag.Int("n", 12, "process the integers 1..n")
	workers := flag.Int("workers", 4, "maximum number of concurrent workers")
	flag.Parse()

	// slices: generate the inputs 1..n.
	inputs := slices.Generate(*n, func(i int) int { return i + 1 })

	// channels: present the inputs as a stream to fan out from.
	stream := channels.FromSlice(inputs)

	// concurrency: a bounded pool squares each item; no more than `workers`
	// run concurrently. The mutex-guarded append is the fan-in.
	limiter := concurrency.NewBlockingWorkLimiter(*workers)

	var mu sync.Mutex
	var squares []int
	var work []concurrency.WorkFunc
	for v := range stream {
		work = append(work, func() error {
			squared := v * v
			mu.Lock()
			squares = append(squares, squared)
			mu.Unlock()
			return nil
		})
	}
	limiter.Run(work)

	// channels: reduce the fanned-in results to a single sum.
	sum := channels.NewPipeline[int, int](channels.FromSlice(squares), func(in <-chan int) <-chan int {
		return channels.Reduce(in, func(acc, x int) int { return acc + x })
	}).CollectAsSlice()

	fmt.Printf("inputs:  %v\n", inputs)
	fmt.Printf("squares: %v\n", slices.SortOrderedAsc(squares))
	fmt.Printf("sum:     %d\n", sum[0])
}
