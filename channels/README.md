# Channels - Pipeline Processing Made Simple

The `channels` package provides pipeline patterns for Go channels, letting you build concurrent data processing without managing channels by hand. Each stage transform consumes a channel and returns a new one, so stages compose by nesting. A `context.Context` governs the lifetime of every stage; cancelling it tears the whole pipeline down and reclaims its goroutines.

The snippets below mirror the compiled package documentation in [`doc.go`](doc.go), whose Quick Start runs as `Example_quickStart` in the test suite and is guaranteed to track the real API.

## Quick Start

```go
import (
    "context"

    "github.com/pickeringtech/go-collections/channels"
)

ctx := context.Background()

// Feed numbers into a channel.
input := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

// Build a pipeline: square every number, then keep the even squares.
pipeline := channels.NewPipeline[int, int](ctx, input, func(ctx context.Context, in <-chan int) <-chan int {
    squares := channels.Map(ctx, in, func(n int) int { return n * n })
    return channels.Filter(ctx, squares, func(n int) bool { return n%2 == 0 })
})

// CollectAsSlice drains the pipeline once the input channel is closed.
results := pipeline.CollectAsSlice()
fmt.Printf("Results: %v\n", results) // [4 16 36 64 100]
```

## Why Use Channel Pipelines?

Native Go channel processing requires complex coordination:

```go
// Manual approach - complex and error-prone
input := make(chan int)
squares := make(chan int)
evens := make(chan int)

// Stage 1: Square numbers
go func() {
    defer close(squares)
    for n := range input {
        squares <- n * n
    }
}()

// Stage 2: Filter evens
go func() {
    defer close(evens)
    for n := range squares {
        if n%2 == 0 {
            evens <- n
        }
    }
}()

// Collect results
var results []int
for n := range evens {
    results = append(results, n)
}
```

The standalone `Map`, `Filter`, and `Reduce` helpers each own one stage's goroutine and channel lifecycle, so the same computation reads as a straight data flow. Each takes a context, so cancelling it reclaims the stage's goroutine:

```go
squares := channels.Map(ctx, input, func(n int) int { return n * n })
evens := channels.Filter(ctx, squares, func(n int) bool { return n%2 == 0 })
results := channels.CollectAsSlice(evens)
```

## Core Concepts

Stage transforms (each consumes a channel and returns a new channel):

- `Map` - transform each element, producing a channel of the result type
- `Filter` - forward only the elements matching a predicate
- `Reduce` - fold the stream into a single running value

Sources and sinks:

- `FromSlice` / `FromMap` - turn a slice or map into a channel
- `CollectAsSlice` / `CollectNAsSlice` - drain a channel into a slice
- `CollectAsMap` / `BuildMapFromEntries` - drain a channel into a map

Pipeline:

- `NewPipeline` pins the input and output types and wires the stages via a `PipelineCreationFunc`; `Pipeline.CollectAsSlice` drains the result.

## Core Operations

### Map - Transform Each Element

```go
// Transform numbers to strings
numbers := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5})
strings := channels.Map(ctx, numbers, func(n int) string { return fmt.Sprintf("num_%d", n) })
result := channels.CollectAsSlice(strings)
// Result: ["num_1", "num_2", "num_3", "num_4", "num_5"]

// Extract fields from structs
users := channels.FromSlice(ctx, []User{
    {Name: "Alice", Age: 25},
    {Name: "Bob", Age: 30},
})
names := channels.Map(ctx, users, func(u User) string { return u.Name })
// channels.CollectAsSlice(names): ["Alice", "Bob"]
```

### Filter - Keep Matching Elements

```go
// Filter even numbers
numbers := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
evens := channels.Filter(ctx, numbers, func(n int) bool { return n%2 == 0 })
result := channels.CollectAsSlice(evens)
// Result: [2, 4, 6, 8, 10]

// Filter active users
users := channels.FromSlice(ctx, []User{
    {Name: "Alice", Active: true},
    {Name: "Bob", Active: false},
    {Name: "Charlie", Active: true},
})
active := channels.Filter(ctx, users, func(u User) bool { return u.Active })
// channels.CollectAsSlice(active): [Alice, Charlie]
```

### Reduce - Fold a Stream

`Reduce` folds a channel down to a single running value, emitted on its own channel so it still composes with the other stages:

```go
input := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5})
totals := channels.Reduce(ctx, input, func(acc, n int) int { return acc + n })
total := channels.CollectAsSlice(totals) // [15]
```

## Composing Stages

Because every stage transform takes a channel and returns a channel, stages compose by nesting - the output of one becomes the input of the next:

```go
input := channels.FromSlice(ctx, []string{"one", "two", "three", "four", "five"})

lengths := channels.Map(ctx, input, func(s string) int { return len(s) })
longish := channels.Filter(ctx, lengths, func(n int) bool { return n >= 4 })

results := channels.CollectAsSlice(longish) // [5 4 4]
```

`NewPipeline` captures that wiring behind a single value with fixed input and output types, which is handy when a pipeline is passed around or returned. It threads its context into the supplied function so every stage shares one cancellation signal:

```go
pipeline := channels.NewPipeline[string, int](ctx, input, func(ctx context.Context, in <-chan string) <-chan int {
    lengths := channels.Map(ctx, in, func(s string) int { return len(s) })
    return channels.Filter(ctx, lengths, func(n int) bool { return n >= 4 })
})
results := pipeline.CollectAsSlice()
```

## Error Handling

The stage helpers do not expose a separate error channel; the idiomatic approach is to carry the error alongside each result and partition downstream:

```go
type Result struct {
    Value int
    Err   error
}

parsed := channels.Map(ctx, input, func(s string) Result {
    n, err := strconv.Atoi(s)
    return Result{Value: n, Err: err}
})

ok := channels.Filter(ctx, parsed, func(r Result) bool { return r.Err == nil })
values := channels.CollectAsSlice(ok)
```

## Performance Considerations

Channel pipelines add per-element goroutine and channel overhead in exchange for streaming, backpressure, and composable stages. Use them for data- and stream-processing workflows; prefer the `slices` package (or a hand-written loop) for in-memory data on a performance-critical hot path.

## Integration with Other Packages

Channels interoperate with the `slices` and `collections` packages - drain a pipeline into a slice, then build a collection from it:

```go
input := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5})
evens := channels.Filter(ctx, input, func(n int) bool { return n%2 == 0 })
results := channels.CollectAsSlice(evens)

resultSet := collections.NewSet(results...)
resultDict := collections.NewDict(
    slices.Map(results, func(n int) collections.Pair[int, int] {
        return collections.Pair[int, int]{Key: n, Value: n * n}
    })...,
)
```

## Quick Reference

```go
// Sources
channels.FromSlice(ctx, []T)                  // Channel from a slice
channels.FromMap(ctx, map[K]V)                // Channel of maps.Entry from a map

// Stage transforms (channel in, channel out)
channels.Map(ctx, in, func(I) O)              // Transform each element
channels.Filter(ctx, in, func(T) bool)        // Keep matching elements
channels.Reduce(ctx, in, func(O, I) O)        // Fold into a running value

// Pipeline
channels.NewPipeline[I, O](ctx, in, fn)       // Wire stages, fixed I/O types
pipeline.CollectAsSlice()                     // Drain the pipeline into a slice

// Sinks
channels.CollectAsSlice(in)                   // Drain a channel into a slice
channels.CollectNAsSlice(in, n)               // Drain the first n elements
channels.CollectAsMap(in, fn)                 // Drain into a map
channels.BuildMapFromEntries(entries)         // Build a map from maps.Entry values
```

Start with `Map` and `Filter`, reach for `Reduce` when folding a stream, and wrap a multi-stage flow in `NewPipeline` when you need to pass it around by value.
