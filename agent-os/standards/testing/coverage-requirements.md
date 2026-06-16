# Test Coverage Requirements

Every function ships with tests before it's considered done.

## Public (exported) functions & types — the trio

1. **Example** — runnable godoc example with `// Output:` (or `// Unordered output:` for maps/sets)
2. **Test** — table-driven (see [[table-driven-tests]])
3. **Benchmark** — scaling ladder (see [[benchmark-scaling]])

```go
func ExampleFilter() { ... }   // doc + compile-checked usage
func TestFilter(t *testing.T)  // correctness
func BenchmarkFilter(b *testing.B) // performance
```

## Private (unexported) functions — Test only

- Table-driven `Test` required.
- No `Example` (not in public docs), no `Benchmark` (benchmark via the public caller).

## Always cover these edge cases

- `nil` input → empty non-nil output (see [[zero-value-output]])
- empty input → empty non-nil output

```go
{name: "nil input results in empty output", args: args{input: nil}, want: []string{}},
{name: "empty input results in empty output", args: args{input: []string{}}, want: []string{}},
```

## Concurrent types

Add a runnable `Example` that spawns goroutines with `sync.WaitGroup` to demonstrate thread safety. Run the suite with `-race`.

## Collection types & transformation functions — add a fuzz target

Beyond the trio, collection types and the slice/map/channel transformation
functions ship with a Go native fuzz target (`FuzzXxx`, `go test -fuzz`). Real
verification — not just hand-picked table cases — is core to the product
mission. Fuzzing surfaces edge-case panics, invariant violations, and inputs
the author never thought to test.

Prefer **invariant / differential** assertions over exact outputs:

- **Round-trip invariants** — `Push` then `Pop` returns the same element;
  `FromSlice` then `CollectAsSlice` reproduces the slice; `Reverse(Reverse(x))`
  equals `x`.
- **Differential oracles** — compare behaviour against the native Go
  equivalent (a `dicts.Hash` vs a built-in `map`, a set's membership vs a
  `map[T]struct{}`).
- **No-panic guarantees** — feed arbitrary operation sequences; assert the
  structure never panics and its invariants hold (length consistent, no
  duplicate set members, list ordering preserved).
- **Concurrent safety** — fuzz operation sequences run under `-race` against
  the `Concurrent*` variants.

Conventions:

- `FuzzXxx` functions live in `fuzz_test.go` alongside the existing
  Example/Test/Benchmark trio.
- Seed each corpus with `f.Add(...)` using the same edge cases the table tests
  cover (empty, nil, single element, duplicates).
- Keep targets deterministic and fast so CI can run them with a short
  `-fuzztime`.
- When an operation has a domain the implementation does not support (e.g. a
  negative page size), return early from the fuzz body rather than asserting on
  undefined behaviour.

## General

- Tests live in an external `_test` package (e.g. `package slices_test`) — black-box, public API only.
- Use the standard library only: `reflect.DeepEqual` + `t.Errorf`. No testify or assertion libraries.
