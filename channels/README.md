# Channels - Pipeline Processing Made Simple

The `channels` package provides pipeline patterns and utilities for Go channels, enabling concurrent data processing, stream processing, and producer-consumer patterns without hand-rolled channel coordination.

Every stage transform takes a `context.Context` so cancelling it tears the whole pipeline down and reclaims its goroutines deterministically.

> The snippets below mirror the package's compiled, runnable examples (`Example_quickStart`, `ExampleMap`, `ExampleFilter`, `ExampleReduce`, and friends in the `*_test.go` files), so they track the real API. For the authoritative reference, see the package godoc in [`doc.go`](doc.go).

## Quick Start

```go
import (
    "context"

    "github.com/pickeringtech/go-collections/channels"
)

// A context governs the lifetime of every stage; cancelling it tears the
// whole pipeline down and reclaims its goroutines.
ctx := context.Background()

// Feed numbers into a channel.
input := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

// Build a pipeline: square every number, then keep the even squares.
// A Pipeline pins down the input and output types; the supplied function
// wires the intermediate stages together with the standalone Map and Filter
// helpers, threading the context through each one.
pipeline := channels.NewPipeline[int, int](ctx, input, func(ctx context.Context, in <-chan int) <-chan int {
    squares := channels.Map(ctx, in, func(n int) int { return n * n })
    return channels.Filter(ctx, squares, func(n int) bool { return n%2 == 0 })
})

// CollectAsSlice drains the pipeline once the input channel is closed.
results := pipeline.CollectAsSlice()
fmt.Printf("Results: %v\n", results) // [4 16 36 64 100]
```

## Why Use Channel Pipelines?

**Native Go channels require complex coordination:**
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

**Stage helpers are clean and safe:**
```go
// The standalone Map, Filter, and Reduce helpers each own one stage's
// goroutine and channel lifecycle, so the same computation reads as a straight
// data flow. Each takes a context, so cancelling it reclaims the stage.
squares := channels.Map(ctx, input, func(n int) int { return n * n })
evens := channels.Filter(ctx, squares, func(n int) bool { return n%2 == 0 })
results := channels.CollectAsSlice(evens)
```

## Core Concepts

Stage transforms each consume a channel and return a new channel:

- **`Map`** — transform each element, producing a channel of the result type.
- **`Filter`** — forward only the elements matching a predicate.
- **`Reduce`** — fold the stream into a single running value, emitted on its own channel.

Sources and sinks:

- **`FromSlice` / `FromMap`** — turn a slice or map into a channel.
- **`CollectAsSlice` / `CollectNAsSlice`** — drain a channel into a slice.
- **`CollectAsMap` / `BuildMapFromEntries`** — drain a channel into a map.

Pipeline:

- **`NewPipeline`** pins the input and output types and wires the stages via a `PipelineCreationFunc`; **`Pipeline.CollectAsSlice`** drains the result.

## Stage Transforms

### Map - Transform Each Element
```go
ctx := context.Background()
input := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5})
doubled := channels.Map(ctx, input, func(n int) int { return n * 2 })
results := channels.CollectAsSlice(doubled)
// results: [2 4 6 8 10]
```

`Map` can change the element type. The type parameters are usually inferred, but you can spell them out when it aids clarity:
```go
input := channels.FromSlice(ctx, []int{1, 2, 3})
labels := channels.Map[int, string](ctx, input, func(n int) string {
    return fmt.Sprintf("num_%d", n)
})
results := channels.CollectAsSlice(labels)
// results: [num_1 num_2 num_3]
```

### Filter - Keep Matching Elements
```go
ctx := context.Background()
input := channels.FromSlice(ctx, []string{"hello", "everyone", "world", "goodness", "gracious"})
output := channels.Filter(ctx, input, func(element string) bool {
    return len(element) > 5
})
results := channels.CollectAsSlice(output)
// results: [everyone goodness gracious]
```

### Reduce - Fold a Stream
```go
ctx := context.Background()
input := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5})
totals := channels.Reduce(ctx, input, func(acc, n int) int { return acc + n })
total := channels.CollectAsSlice(totals)
// total: [15]
```

`Reduce` emits its running value on a channel so it still composes with the other stages.

## Composing Stages

Because every stage transform takes a channel and returns a channel, stages compose by nesting — the output of one becomes the input of the next:

```go
ctx := context.Background()
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

A pipeline can change types between its ends — for example, total the input and then stringify it:

```go
ctx := context.Background()
input := channels.FromSlice(ctx, []int{1, 2, 5, 4, 3})

pipeline := channels.NewPipeline[int, string](ctx, input, func(ctx context.Context, in <-chan int) <-chan string {
    totals := channels.Reduce(ctx, in, func(acc, n int) int { return acc + n })
    return channels.Map[int, string](ctx, totals, func(n int) string {
        return strconv.Itoa(n)
    })
})

results := pipeline.CollectAsSlice() // ["15"]
```

## Sources and Sinks

### FromSlice / FromMap - Build a Channel
```go
ctx := context.Background()

// From a slice
nums := channels.FromSlice(ctx, []int{1, 2, 5, 4, 3})

// From a map (each entry becomes a maps.Entry on the channel)
entries := channels.FromMap(ctx, map[int]string{1: "one", 2: "two"})
```

### Collecting Results
```go
ctx := context.Background()

// Drain everything into a slice (blocks until the channel closes).
input := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5})
doubled := channels.Map(ctx, input, func(n int) int { return n * 2 })
all := channels.CollectAsSlice(doubled) // [2 4 6 8 10]

// Or take only the first N elements.
more := channels.Map(ctx, channels.FromSlice(ctx, []int{1, 2, 3, 4, 5}), func(n int) int { return n * 2 })
firstThree := channels.CollectNAsSlice(more, 3) // [2 4 6]
```

`CollectAsMap` drains a channel into a map by deriving a key/value entry from each element:
```go
ctx := context.Background()
input := channels.FromSlice(ctx, []string{"hello", "and", "world"})
lengths := channels.CollectAsMap(input, func(s string) maps.Entry[string, int] {
    return maps.Entry[string, int]{Key: s, Value: len(s)}
})
// lengths: map[and:3 hello:5 world:5]
```

## Error Handling

The stage helpers do not expose a separate error channel. The idiomatic approach is to carry the error alongside each result and partition downstream:

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

## Cancellation

Every stage transform takes a context. Cancelling it stops each stage's goroutine and reclaims it, so a pipeline never leaks goroutines when its consumer goes away:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

input := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5})
doubled := channels.Map(ctx, input, func(n int) int { return n * 2 })

// Take what you need, then cancel to tear the rest of the pipeline down.
firstTwo := channels.CollectNAsSlice(doubled, 2)
cancel()
```

## Performance Considerations

Channel pipelines add per-element goroutine and channel overhead in exchange for streaming, backpressure, and composable stages. Use them for data- and stream-processing workflows; prefer the [`slices`](../slices) package (or a hand-written loop) for in-memory data on a performance-critical hot path.

## Integration with Other Packages

Channels interoperate with the `slices` and `collections` packages — drain a pipeline into a slice, then build a collection from it:

```go
ctx := context.Background()
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
channels.FromSlice(ctx, []T{...})            // slice  -> channel
channels.FromMap(ctx, map[K]V{...})          // map    -> channel of maps.Entry[K, V]

// Stage transforms (each: channel -> channel)
channels.Map(ctx, in, func(I) O)             // transform each element
channels.Filter(ctx, in, func(T) bool)       // keep matching elements
channels.Reduce(ctx, in, func(O, I) O)       // fold to a running value

// Pipeline (fixed input/output types)
pipeline := channels.NewPipeline[I, O](ctx, in, func(ctx, in) <-chan O { ... })
pipeline.CollectAsSlice()                    // drain the pipeline into a slice

// Sinks
channels.CollectAsSlice(in)                  // channel -> []T
channels.CollectNAsSlice(in, n)              // channel -> first n elements
channels.CollectAsMap(in, func(I) maps.Entry[OK, OV]) // channel -> map[OK]OV
channels.BuildMapFromEntries(entries)        // []maps.Entry[K, V] -> map[K]V
```

Start with `Map` and `Filter`, reach for `Reduce` when folding a stream, and wrap a multi-stage flow in `NewPipeline` when you need to pass it around by value.
