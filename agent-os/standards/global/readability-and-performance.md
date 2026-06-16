# Readability AND Performance

Production code must be **easy to read** and **highly performant**. Anyone should be able to open a function and immediately understand what it does.

## Default: write the clear version

The exported function is the readable version. Optimize for comprehension first.

```go
// Filter returns a new slice of elements for which fn returns true.
func Filter[T any](input []T, fn FilterFunc[T]) []T {
	var output []T
	for _, element := range input {
		if fn(element) {
			output = append(output, element)
		}
	}
	return output
}
```

## When perf conflicts with readability: add a `Fast` variant

Do **not** sacrifice the default's clarity. Instead expose a sibling with the `Fast` suffix.

```go
// FilterFast is an optimized Filter that preallocates capacity to avoid
// repeated growth. Trades readability for speed on large inputs.
func FilterFast[T any](input []T, fn FilterFunc[T]) []T { ... }
```

Rules for a `Fast` variant:

- **Earn it with a benchmark.** Only add one when the scaling ladder ([[benchmark-scaling]]) shows the readable version is a real bottleneck. Don't pre-emptively double the API.
- **Identical results.** `Fast` must return exactly what the default returns — enforce with a shared equivalence test.
- **Same coverage.** It gets the same test trio as any public function ([[coverage-requirements]]).
- **Comment the trade.** State what readability/safety is given up and why.

The readable default is canonical; `Fast` is an opt-in customization for hot paths.
