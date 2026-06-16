package multimaps_test

import (
	"strconv"
	"testing"

	"github.com/pickeringtech/go-collections/collections/multimaps"
)

var benchSizes = []struct {
	name string
	size int
}{
	{"3", 3},
	{"10", 10},
	{"100", 100},
	{"1_000", 1_000},
	{"10_000", 10_000},
	{"100_000", 100_000},
	{"1_000_000", 1_000_000},
}

// makeEntries builds n entries spread across roughly n/4 keys, so most keys
// carry several values — exercising the many-values-per-key shape of a multimap.
func makeEntries(n int) []multimaps.Entry[string, int] {
	entries := make([]multimaps.Entry[string, int], n)
	keys := n/4 + 1
	for i := 0; i < n; i++ {
		entries[i] = multimaps.Entry[string, int]{
			Key:   strconv.Itoa(i % keys),
			Value: i,
		}
	}
	return entries
}

func BenchmarkListMultimap_PutInPlace(b *testing.B) {
	for _, bm := range benchSizes {
		entries := makeEntries(bm.size)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := multimaps.NewListMultimap[string, int]()
				for _, e := range entries {
					m.PutInPlace(e.Key, e.Value)
				}
			}
		})
	}
}

func BenchmarkSetMultimap_PutInPlace(b *testing.B) {
	for _, bm := range benchSizes {
		entries := makeEntries(bm.size)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := multimaps.NewSetMultimap[string, int]()
				for _, e := range entries {
					m.PutInPlace(e.Key, e.Value)
				}
			}
		})
	}
}

func BenchmarkListMultimap_Get(b *testing.B) {
	for _, bm := range benchSizes {
		m := multimaps.NewListMultimap(makeEntries(bm.size)...)
		key := strconv.Itoa(bm.size / 8)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = m.Get(key)
			}
		})
	}
}

func BenchmarkSetMultimap_ContainsEntry(b *testing.B) {
	for _, bm := range benchSizes {
		m := multimaps.NewSetMultimap(makeEntries(bm.size)...)
		key := strconv.Itoa(bm.size / 8)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = m.ContainsEntry(key, 0)
			}
		})
	}
}

func BenchmarkListMultimap_Filter(b *testing.B) {
	for _, bm := range benchSizes {
		m := multimaps.NewListMultimap(makeEntries(bm.size)...)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = m.Filter(func(_ string, value int) bool { return value%2 == 0 })
			}
		})
	}
}
