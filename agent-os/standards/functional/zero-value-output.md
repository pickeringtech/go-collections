# Zero-Value Output Contract

For `nil` or empty input, return an **initialized-but-empty, non-nil** result — never `nil`.

```go
// slices
func Filter[T any](input []T, fn FilterFunc[T]) []T {
	output := []T{}            // non-nil, even when nothing matches
	for _, e := range input {
		if fn(e) { output = append(output, e) }
	}
	return output
}

// maps
func Map[...](input map[K]V, fn MapFunc[...]) map[OK]OV {
	results := map[OK]OV{}     // non-nil empty map
	// ...
	return results
}
```

- **slices** → `[]T{}`
- **maps** → `map[K]V{}`
- **channels** → an open channel that is `close`d immediately when the (empty) input drains

Callers can always range/len the result without a nil check. Tests assert the empty non-nil form (see [[coverage-requirements]]).

> Gap: existing `slices` functions return `nil` on empty/nil input (they declare `var output []T` and never append). Migrate them to initialize `[]T{}`, and update their tests from `want: nil` to `want: []T{}`.
