# Benchmark Scaling Ladder

Every public function gets a `Benchmark` that runs sub-benchmarks across a fixed size ladder:

**3, 10, 100, 1_000, 10_000, 100_000, 1_000_000**

```go
func BenchmarkFilter(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
		fn   func(int) bool
	}{
		{name: "3 elements", sli: []int{1, 2, 3}, fn: func(e int) bool { return e >= 2 }},
		{name: "10 elements", sli: slices.Generate(10, slices.NumericIdentityGenerator[int]), fn: func(e int) bool { return e >= 5 }},
		// ... 100, 1_000, 10_000, 100_000, 1_000_000
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Filter(bm.sli, bm.fn)
			}
		})
	}
}
```

- Generate large inputs with `slices.Generate(n, slices.NumericIdentityGenerator[int])`.
- Use `b.Run(bm.name, ...)` per size; size names use `_` digit separators (`1_000_000`).
- Assign the result to `_` to prevent the call being optimized away.
