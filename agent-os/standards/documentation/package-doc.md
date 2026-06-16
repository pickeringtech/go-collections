# Package Documentation (doc.go)

Every package ships a `doc.go` holding the package-level godoc comment. Tab-indent all code blocks (godoc convention).

## Required sections

1. **Summary paragraph** — opens `// Package <name> provides ...`; one short paragraph on purpose.
2. **`# Quick Start`** — import line, then one idiomatic example annotated with `// Result:` comments.
3. **`# Why Use ...?`** or **`# When to Use X vs Y`** — justify the package with a native-vs-this comparison (see below).

```go
// Package slices provides functional programming utilities for Go slices.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/slices"
//
//	result := slices.Filter(numbers, func(n int) bool { return n%2 == 0 })
//	// Result: [2 4 6]
```

## Native-vs-this comparison

Show the verbose native Go, then the concise package equivalent:

```go
// # Why Use Slices Package?
//
// Native approach — verbose and error-prone:
//	var evens []int
//	for _, n := range numbers { if n%2 == 0 { evens = append(evens, n) } }
//
// With this package:
//	evens := slices.Filter(numbers, func(n int) bool { return n%2 == 0 })
```

## Optional sections — use the standard vocabulary

Add when relevant to the package; reuse these exact headings (don't invent synonyms):

`# Common Patterns` · `# Performance` · `# Thread Safety` · `# Immutable vs Mutable Operations` · `# Available Implementations` · `# Integration with Other Packages`
