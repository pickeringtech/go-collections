package lru_test

import (
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lru"
)

// concurrentFactory builds the two thread-safe variants for race testing.
type concurrentFactory struct {
	name string
	make func(capacity int, opts ...lru.Option[int, int]) lru.MutableCache[int, int]
}

func concurrentFactories() []concurrentFactory {
	return []concurrentFactory{
		{"ConcurrentLRU", func(c int, o ...lru.Option[int, int]) lru.MutableCache[int, int] {
			return lru.NewConcurrentLRU(c, o...)
		}},
		{"ConcurrentLRURW", func(c int, o ...lru.Option[int, int]) lru.MutableCache[int, int] {
			return lru.NewConcurrentLRURW(c, o...)
		}},
	}
}

// TestConcurrent_NoRaces hammers every kind of operation from many goroutines.
// Run with -race, it asserts there are no data races or panics, and that the
// cache never exceeds its capacity bound.
func TestConcurrent_NoRaces(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			const capacity = 64
			cache := f.make(capacity)

			var wg sync.WaitGroup
			const workers = 16
			const iterations = 500
			for w := 0; w < workers; w++ {
				wg.Add(1)
				go func(base int) {
					defer wg.Done()
					for i := 0; i < iterations; i++ {
						key := (base + i) % 128
						cache.PutInPlace(key, key)
						cache.Get(key)
						cache.Peek(key)
						cache.Contains(key)
						cache.Length()
						cache.IsEmpty()
						cache.Capacity()
						cache.Keys()
						cache.Values()
						cache.Items()
						cache.AsMap()
						cache.ForEach(func(int, int) {})
						for range cache.All() {
						}
						_ = cache.Put(key, key)
						_ = cache.Remove(key)
						cache.RemoveInPlace(key)
					}
				}(w * 7)
			}
			wg.Wait()

			if cache.Length() > capacity {
				t.Errorf("Length %d exceeds capacity %d", cache.Length(), capacity)
			}
		})
	}
}

// TestConcurrent_ParallelWritesRespectCapacity confirms a flood of distinct
// concurrent writes still leaves the cache at exactly its bound.
func TestConcurrent_ParallelWritesRespectCapacity(t *testing.T) {
	for _, f := range concurrentFactories() {
		t.Run(f.name, func(t *testing.T) {
			const capacity = 100
			cache := f.make(capacity)

			var wg sync.WaitGroup
			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func(n int) {
					defer wg.Done()
					cache.PutInPlace(n, n)
				}(i)
			}
			wg.Wait()

			if cache.Length() != capacity {
				t.Errorf("Length() = %d, want %d", cache.Length(), capacity)
			}
		})
	}
}
