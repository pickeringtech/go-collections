package maps_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/maps"
	"github.com/pickeringtech/go-collections/slices"
)

// This file backfills the Benchmark leg of the Example+Test+Benchmark trio for
// every public function in the maps package (issue #52), using the scaling
// ladder from agent-os/standards/testing/benchmark-scaling.md: each benchmark
// sub-benchmarks across a fixed size set via b.Run, with the result assigned to
// _ so the call is not optimised away.

// ladder is the shared element-count matrix for read- and allocate-style
// benchmarks that build their input once per size and reuse it. The name field
// gives each sub-benchmark the underscore-separated digits the standard asks
// for (e.g. "1_000_000 elements").
var ladder = []struct {
	name string
	n    int
}{
	{"3 elements", 3},
	{"10 elements", 10},
	{"100 elements", 100},
	{"1_000 elements", 1_000},
	{"10_000 elements", 10_000},
	{"100_000 elements", 100_000},
	{"1_000_000 elements", 1_000_000},
}

// intMap returns a map of n entries keyed 0..n-1 with the value mirroring the
// key, so benchmarks have deterministic, comparably-sized input at every size.
func intMap(n int) map[int]int {
	m := make(map[int]int, n)
	for i := 0; i < n; i++ {
		m[i] = i
	}
	return m
}

// intKeys returns the keys 0..n-1, used to drive the bulk-retrieval benchmarks.
func intKeys(n int) []int {
	return slices.Generate(n, slices.NumericIdentityGenerator[int])
}

func BenchmarkFromKeys(b *testing.B) {
	for _, bm := range ladder {
		keys := intKeys(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = maps.FromKeys(keys, 0)
			}
		})
	}
}

func BenchmarkClear(b *testing.B) {
	for _, bm := range ladder {
		// Build the map once per size cell, then each iteration clears it and
		// refills it (reusing the map's retained capacity, so no reallocation),
		// measuring a clear+refill cycle. The old code rebuilt the map under
		// StopTimer every iteration; a fresh make+populate costs far more than the
		// delete loop (especially at small sizes), so wall-time ≈ b.N × O(size)
		// and blew up CI — the hour-plus of issue #112, which is also why this had
		// to cap at 1000. Clear has no O(1) inverse, but both halves of the cycle
		// are O(size), so the full ladder is back and the linear scaling is
		// faithful; the refill is timed rather than excluded with a per-iteration
		// b.StopTimer() because that reads memstats under -benchmem and, at the
		// millions of iterations the small cells reach, would itself blow up.
		m := intMap(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				maps.Clear(m)
				for k := 0; k < bm.n; k++ {
					m[k] = k
				}
			}
		})
	}
}

func BenchmarkContainsValue(b *testing.B) {
	for _, bm := range ladder {
		m := intMap(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Search for a value that is absent so the full map is scanned,
				// measuring the worst case rather than an early hit.
				_ = maps.ContainsValue(m, -1)
			}
		})
	}
}

func BenchmarkCopy(b *testing.B) {
	for _, bm := range ladder {
		m := intMap(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = maps.Copy(m)
			}
		})
	}
}

func BenchmarkGetMany(b *testing.B) {
	for _, bm := range ladder {
		m := intMap(bm.n)
		keys := intKeys(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = maps.GetMany(m, keys)
			}
		})
	}
}

func BenchmarkGetManyOrDefault(b *testing.B) {
	for _, bm := range ladder {
		m := intMap(bm.n)
		keys := intKeys(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = maps.GetManyOrDefault(m, keys, 0)
			}
		})
	}
}

func BenchmarkGetOrDefault(b *testing.B) {
	for _, bm := range ladder {
		m := intMap(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = maps.GetOrDefault(m, i%bm.n, -1)
			}
		})
	}
}

func BenchmarkItems(b *testing.B) {
	for _, bm := range ladder {
		m := intMap(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = maps.Items(m)
			}
		})
	}
}

func BenchmarkKeys(b *testing.B) {
	for _, bm := range ladder {
		m := intMap(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = maps.Keys(m)
			}
		})
	}
}

func BenchmarkValues(b *testing.B) {
	for _, bm := range ladder {
		m := intMap(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = maps.Values(m)
			}
		})
	}
}

func BenchmarkFilter(b *testing.B) {
	for _, bm := range ladder {
		m := intMap(bm.n)
		fn := func(key, value int) bool { return value%2 == 0 }
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = maps.Filter(m, fn)
			}
		})
	}
}

func BenchmarkMap(b *testing.B) {
	for _, bm := range ladder {
		m := intMap(bm.n)
		fn := func(key, value int) (int, int) { return key, value * 2 }
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = maps.Map(m, fn)
			}
		})
	}
}

func BenchmarkUpdate(b *testing.B) {
	for _, bm := range ladder {
		m := intMap(bm.n)
		update := intMap(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = maps.Update(m, update)
			}
		})
	}
}
