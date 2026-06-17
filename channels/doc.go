// Package channels provides pipeline patterns and utilities for Go channels,
// enabling concurrent data processing, stream processing, and producer-consumer
// patterns without hand-rolled channel coordination.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/channels"
//
//	// Feed numbers into a channel.
//	input := channels.FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
//
//	// Build a pipeline: square every number, then keep the even squares.
//	// A Pipeline pins down the input and output types; the supplied function
//	// wires the intermediate stages together with the standalone Map and Filter
//	// helpers.
//	pipeline := channels.NewPipeline[int, int](input, func(in <-chan int) <-chan int {
//		squares := channels.Map(in, func(n int) int { return n * n })
//		return channels.Filter(squares, func(n int) bool { return n%2 == 0 })
//	})
//
//	// CollectAsSlice drains the pipeline once the input channel is closed.
//	results := pipeline.CollectAsSlice()
//	// results: [4 16 36 64 100]
//
// This Quick Start is compiled and run as Example_quickStart in the package's
// test suite, so it is guaranteed to track the real API.
//
// # Why Use Channel Pipelines?
//
// Native Go channel processing requires complex coordination:
//
//	// Manual approach - complex and error-prone
//	input := make(chan int)
//	squares := make(chan int)
//	evens := make(chan int)
//
//	// Stage 1: Square numbers
//	go func() {
//		defer close(squares)
//		for n := range input {
//			squares <- n * n
//		}
//	}()
//
//	// Stage 2: Filter evens
//	go func() {
//		defer close(evens)
//		for n := range squares {
//			if n%2 == 0 {
//				evens <- n
//			}
//		}
//	}()
//
//	// Collect results
//	var results []int
//	for n := range evens {
//		results = append(results, n)
//	}
//
// The standalone Map, Filter, and Reduce helpers each own one stage's goroutine
// and channel lifecycle, so the same computation reads as a straight data flow:
//
//	squares := channels.Map(input, func(n int) int { return n * n })
//	evens := channels.Filter(squares, func(n int) bool { return n%2 == 0 })
//	results := channels.CollectAsSlice(evens)
//
// # Core Concepts
//
// Stage transforms (each consumes a channel and returns a new channel):
//   - Map: transform each element, producing a channel of the result type
//   - Filter: forward only the elements matching a predicate
//   - Reduce: fold the stream into a single running value
//
// Sources and sinks:
//   - FromSlice / FromMap: turn a slice or map into a channel
//   - CollectAsSlice / CollectNAsSlice: drain a channel into a slice
//   - CollectAsMap / BuildMapFromEntries: drain a channel into a map
//
// Pipeline:
//   - NewPipeline pins the input and output types and wires the stages via a
//     PipelineCreationFunc; Pipeline.CollectAsSlice drains the result.
//
// # Composing Stages
//
// Because every stage transform takes a channel and returns a channel, stages
// compose by nesting - the output of one becomes the input of the next:
//
//	input := channels.FromSlice([]string{"one", "two", "three", "four", "five"})
//
//	lengths := channels.Map(input, func(s string) int { return len(s) })
//	longish := channels.Filter(lengths, func(n int) bool { return n >= 4 })
//
//	results := channels.CollectAsSlice(longish) // [5 4 4]
//
// NewPipeline captures that wiring behind a single value with fixed input and
// output types, which is handy when a pipeline is passed around or returned:
//
//	pipeline := channels.NewPipeline[string, int](input, func(in <-chan string) <-chan int {
//		lengths := channels.Map(in, func(s string) int { return len(s) })
//		return channels.Filter(lengths, func(n int) bool { return n >= 4 })
//	})
//	results := pipeline.CollectAsSlice()
//
// # Reducing a Stream
//
// Reduce folds a channel down to a single running value, emitted on its own
// channel so it still composes with the other stages:
//
//	input := channels.FromSlice([]int{1, 2, 3, 4, 5})
//	totals := channels.Reduce(input, func(acc, n int) int { return acc + n })
//	total := channels.CollectAsSlice(totals) // [15]
//
// # Error Handling
//
// The stage helpers do not expose a separate error channel; the idiomatic
// approach is to carry the error alongside each result and partition downstream:
//
//	type Result struct {
//		Value int
//		Err   error
//	}
//
//	parsed := channels.Map(input, func(s string) Result {
//		n, err := strconv.Atoi(s)
//		return Result{Value: n, Err: err}
//	})
//
//	ok := channels.Filter(parsed, func(r Result) bool { return r.Err == nil })
//	values := channels.CollectAsSlice(ok)
//
// # Performance Considerations
//
// Channel pipelines add per-element goroutine and channel overhead in exchange
// for streaming, backpressure, and composable stages. Use them for data- and
// stream-processing workflows; prefer the slices package (or a hand-written
// loop) for in-memory data on a performance-critical hot path.
//
// # Integration with Other Packages
//
// Channels interoperate with the slices and collections packages - drain a
// pipeline into a slice, then build a collection from it:
//
//	input := channels.FromSlice([]int{1, 2, 3, 4, 5})
//	evens := channels.Filter(input, func(n int) bool { return n%2 == 0 })
//	results := channels.CollectAsSlice(evens)
//
//	resultSet := collections.NewSet(results...)
//	resultDict := collections.NewDict(
//		slices.Map(results, func(n int) collections.Pair[int, int] {
//			return collections.Pair[int, int]{Key: n, Value: n * n}
//		})...,
//	)
//
// Start with Map and Filter, reach for Reduce when folding a stream, and wrap a
// multi-stage flow in NewPipeline when you need to pass it around by value.
package channels
