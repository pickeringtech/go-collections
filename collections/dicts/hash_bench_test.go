package dicts_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/collections/dicts"
	"testing"
)

func BenchmarkHash_Get(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Setup
			pairs := make([]dicts.Pair[int, string], size)
			for i := 0; i < size; i++ {
				pairs[i] = dicts.Pair[int, string]{Key: i, Value: fmt.Sprintf("value_%d", i)}
			}
			h := dicts.NewHash(pairs...)
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				key := i % size
				_, _ = h.Get(key, "default")
			}
		})
	}
}

func BenchmarkHash_Put(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Setup
			pairs := make([]dicts.Pair[int, string], size)
			for i := 0; i < size; i++ {
				pairs[i] = dicts.Pair[int, string]{Key: i, Value: fmt.Sprintf("value_%d", i)}
			}
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				h := dicts.NewHash(pairs...)
				_ = h.Put(size+i, fmt.Sprintf("new_value_%d", i))
			}
		})
	}
}

func BenchmarkHash_PutInPlace(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				// Setup fresh hash for each iteration
				pairs := make([]dicts.Pair[int, string], size)
				for j := 0; j < size; j++ {
					pairs[j] = dicts.Pair[int, string]{Key: j, Value: fmt.Sprintf("value_%d", j)}
				}
				h := dicts.NewHash(pairs...)
				b.StartTimer()
				
				h.PutInPlace(size+i, fmt.Sprintf("new_value_%d", i))
			}
		})
	}
}

func BenchmarkHash_Remove(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Setup
			pairs := make([]dicts.Pair[int, string], size)
			for i := 0; i < size; i++ {
				pairs[i] = dicts.Pair[int, string]{Key: i, Value: fmt.Sprintf("value_%d", i)}
			}
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				h := dicts.NewHash(pairs...)
				key := i % size
				_ = h.Remove(key)
			}
		})
	}
}

func BenchmarkHash_RemoveInPlace(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				// Setup fresh hash for each iteration
				pairs := make([]dicts.Pair[int, string], size)
				for j := 0; j < size; j++ {
					pairs[j] = dicts.Pair[int, string]{Key: j, Value: fmt.Sprintf("value_%d", j)}
				}
				h := dicts.NewHash(pairs...)
				b.StartTimer()
				
				key := i % size
				_, _ = h.RemoveInPlace(key)
			}
		})
	}
}

func BenchmarkHash_ForEach(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Setup
			pairs := make([]dicts.Pair[int, string], size)
			for i := 0; i < size; i++ {
				pairs[i] = dicts.Pair[int, string]{Key: i, Value: fmt.Sprintf("value_%d", i)}
			}
			h := dicts.NewHash(pairs...)
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				count := 0
				h.ForEach(func(key int, value string) {
					count++
				})
			}
		})
	}
}

func BenchmarkHash_Filter(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Setup
			pairs := make([]dicts.Pair[int, string], size)
			for i := 0; i < size; i++ {
				pairs[i] = dicts.Pair[int, string]{Key: i, Value: fmt.Sprintf("value_%d", i)}
			}
			h := dicts.NewHash(pairs...)
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				_ = h.Filter(func(key int, value string) bool {
					return key%2 == 0
				})
			}
		})
	}
}

func BenchmarkConcurrentHash_Get(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Setup
			pairs := make([]dicts.Pair[int, string], size)
			for i := 0; i < size; i++ {
				pairs[i] = dicts.Pair[int, string]{Key: i, Value: fmt.Sprintf("value_%d", i)}
			}
			h := dicts.NewConcurrentHash(pairs...)
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				key := i % size
				_, _ = h.Get(key, "default")
			}
		})
	}
}

func BenchmarkConcurrentHashRW_Get(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Setup
			pairs := make([]dicts.Pair[int, string], size)
			for i := 0; i < size; i++ {
				pairs[i] = dicts.Pair[int, string]{Key: i, Value: fmt.Sprintf("value_%d", i)}
			}
			h := dicts.NewConcurrentHashRW(pairs...)
			
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				key := i % size
				_, _ = h.Get(key, "default")
			}
		})
	}
}

// Comparison benchmark between different implementations
func BenchmarkComparison_Get(b *testing.B) {
	size := 1000
	pairs := make([]dicts.Pair[int, string], size)
	for i := 0; i < size; i++ {
		pairs[i] = dicts.Pair[int, string]{Key: i, Value: fmt.Sprintf("value_%d", i)}
	}
	
	b.Run("Hash", func(b *testing.B) {
		h := dicts.NewHash(pairs...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := i % size
			_, _ = h.Get(key, "default")
		}
	})
	
	b.Run("ConcurrentHash", func(b *testing.B) {
		h := dicts.NewConcurrentHash(pairs...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := i % size
			_, _ = h.Get(key, "default")
		}
	})
	
	b.Run("ConcurrentHashRW", func(b *testing.B) {
		h := dicts.NewConcurrentHashRW(pairs...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := i % size
			_, _ = h.Get(key, "default")
		}
	})
	
	b.Run("NativeMap", func(b *testing.B) {
		m := make(map[int]string, size)
		for _, pair := range pairs {
			m[pair.Key] = pair.Value
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := i % size
			if val, ok := m[key]; ok {
				_ = val
			} else {
				_ = "default"
			}
		}
	})
}
