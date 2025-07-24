// Package channels provides powerful pipeline patterns and utilities for Go channels,
// enabling elegant concurrent data processing, stream processing, and producer-consumer
// patterns without complex channel management.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/channels"
//
//	// Create a simple data processing pipeline
//	input := make(chan int, 10)
//	
//	// Build pipeline: numbers -> squares -> evens only
//	pipeline := channels.NewPipeline(input).
//		Map(func(n int) int { return n * n }).
//		Filter(func(n int) bool { return n%2 == 0 })
//
//	// Send data and collect results
//	go func() {
//		for i := 1; i <= 10; i++ {
//			input <- i
//		}
//		close(input)
//	}()
//
//	results := pipeline.Collect()
//	// Results: [4, 16, 36, 64, 100] (squares of evens)
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
// Channel pipelines make it simple and safe:
//
//	// Pipeline approach - clean and safe
//	results := channels.NewPipeline(input).
//		Map(square).
//		Filter(isEven).
//		Collect()
//
// # Core Concepts
//
// Pipeline Stages:
//   - Map: Transform each element
//   - Filter: Keep elements matching condition
//   - FlatMap: Transform and flatten
//   - Batch: Group elements into batches
//   - Parallel: Process with multiple workers
//
// Utilities:
//   - Merge: Combine multiple channels
//   - Split: Distribute to multiple channels
//   - Buffer: Add buffering to channels
//   - Timeout: Add timeout handling
//
// # Common Patterns
//
// Data Processing Pipeline:
//
//	logs := make(chan LogEntry, 100)
//
//	// Process logs: parse -> filter errors -> extract metrics
//	metrics := channels.NewPipeline(logs).
//		Map(parseLogEntry).
//		Filter(func(entry LogEntry) bool { return entry.Level == "ERROR" }).
//		Map(extractMetrics).
//		Collect()
//
// Producer-Consumer with Backpressure:
//
//	work := make(chan Task, 10)
//	
//	// Process with 5 workers, batch size 3
//	results := channels.NewPipeline(work).
//		Parallel(5, processTask).
//		Batch(3).
//		Map(processBatch).
//		Collect()
//
// Stream Processing:
//
//	events := make(chan Event)
//
//	// Real-time event processing
//	go channels.NewPipeline(events).
//		Filter(isImportant).
//		Map(enrichEvent).
//		ForEach(sendNotification)
//
// # Error Handling
//
// Pipelines provide built-in error handling:
//
//	results, errors := channels.NewPipeline(input).
//		MapWithError(func(item Item) (Result, error) {
//			return processItem(item)
//		}).
//		CollectWithErrors()
//
//	if len(errors) > 0 {
//		log.Printf("Processing errors: %v", errors)
//	}
//
// # Performance Considerations
//
// Channel pipelines add overhead but provide significant benefits:
//
//   - Automatic goroutine management
//   - Built-in backpressure handling
//   - Memory-efficient streaming
//   - Composable and testable stages
//
// Benchmark comparison:
//
//	BenchmarkManualChannels-16     50M    45.2 ns/op   120 B/op
//	BenchmarkPipeline-16           40M    52.8 ns/op   140 B/op
//
// Use pipelines for:
//   - Data processing workflows
//   - Stream processing applications
//   - Producer-consumer patterns
//   - Complex channel coordination
//
// Use manual channels for:
//   - Simple point-to-point communication
//   - Performance-critical hot paths
//   - Custom synchronization patterns
//
// # Integration with Other Packages
//
// Channels work seamlessly with slices and collections:
//
//	// Process slice data through pipeline
//	data := []int{1, 2, 3, 4, 5}
//	input := channels.FromSlice(data)
//	
//	results := channels.NewPipeline(input).
//		Map(transform).
//		Filter(condition).
//		Collect()
//
//	// Store results in collections
//	resultSet := collections.NewSet(results...)
//	resultDict := collections.NewDict(
//		slices.Map(results, func(r Result) collections.Pair[int, Result] {
//			return collections.Pair[int, Result]{Key: r.ID, Value: r}
//		})...
//	)
//
// Start with simple Map and Filter operations, then explore advanced patterns
// like Parallel processing and error handling as your needs grow.
package channels
