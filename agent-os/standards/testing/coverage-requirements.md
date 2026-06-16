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

- Tests live in an external `_test` package (e.g. `package slices_test`) — black-box, public API only.
- Use the standard library only: `reflect.DeepEqual` + `t.Errorf`. No testify or assertion libraries.
