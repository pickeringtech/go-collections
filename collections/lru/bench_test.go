package lru_test

import (
	"fmt"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lru"
)

// ladder is the fixed benchmark scaling ladder used across the repo.
var ladder = []int{3, 10, 100, 1_000, 10_000, 100_000, 1_000_000}

// buildLRU returns a plain LRU pre-filled with size sequential int entries, big
// enough to hold them all without eviction.
func buildLRU(size int) *lru.LRU[int, int] {
	cache := lru.NewLRU[int, int](size)
	for i := 0; i < size; i++ {
		cache.PutInPlace(i, i)
	}
	return cache
}

func BenchmarkLRU_Get(b *testing.B) {
	for _, size := range ladder {
		b.Run(fmt.Sprintf("%d_elements", size), func(b *testing.B) {
			cache := buildLRU(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = cache.Get(i % size)
			}
		})
	}
}

func BenchmarkLRU_Peek(b *testing.B) {
	for _, size := range ladder {
		b.Run(fmt.Sprintf("%d_elements", size), func(b *testing.B) {
			cache := buildLRU(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = cache.Peek(i % size)
			}
		})
	}
}

func BenchmarkLRU_PutInPlace(b *testing.B) {
	for _, size := range ladder {
		b.Run(fmt.Sprintf("%d_elements", size), func(b *testing.B) {
			// A cache bounded to size, so steady-state inserts each evict one
			// entry — the realistic hot path for an LRU at capacity.
			cache := buildLRU(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cache.PutInPlace(size+i, i)
			}
		})
	}
}

func BenchmarkLRU_RemoveInPlace(b *testing.B) {
	for _, size := range ladder {
		b.Run(fmt.Sprintf("%d_elements", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				cache := buildLRU(size)
				b.StartTimer()
				cache.RemoveInPlace(i % size)
			}
		})
	}
}

func BenchmarkLRU_Put(b *testing.B) {
	// Immutable Put copies the whole cache, so it scales linearly — kept on the
	// same ladder to make that cost visible against the O(1) in-place variant.
	for _, size := range ladder {
		b.Run(fmt.Sprintf("%d_elements", size), func(b *testing.B) {
			cache := buildLRU(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = cache.Put(size+i, i)
			}
		})
	}
}

// BenchmarkGet_ByImplementation contrasts the per-operation cost of the three
// implementations so the locking overhead of the concurrent variants is visible.
func BenchmarkGet_ByImplementation(b *testing.B) {
	impls := []struct {
		name string
		make func(capacity int) lru.MutableCache[int, int]
	}{
		{"LRU", func(c int) lru.MutableCache[int, int] { return lru.NewLRU[int, int](c) }},
		{"ConcurrentLRU", func(c int) lru.MutableCache[int, int] { return lru.NewConcurrentLRU[int, int](c) }},
		{"ConcurrentLRURW", func(c int) lru.MutableCache[int, int] { return lru.NewConcurrentLRURW[int, int](c) }},
	}
	const size = 1_000
	for _, impl := range impls {
		b.Run(impl.name, func(b *testing.B) {
			cache := impl.make(size)
			for i := 0; i < size; i++ {
				cache.PutInPlace(i, i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = cache.Get(i % size)
			}
		})
	}
}
