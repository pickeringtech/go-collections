package streaming_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/streaming"
	"github.com/pickeringtech/go-collections/slices"
)

// topKBenchmarkInputs returns the scaling ladder of stream lengths, each a
// descending sequence so every Add after the first k must displace the heap
// minimum — the worst case for top-k.
func topKBenchmarkInputs() []struct {
	name string
	sli  []int
} {
	return []struct {
		name string
		sli  []int
	}{
		{name: "3", sli: []int{3, 2, 1}},
		{name: "10", sli: slices.Generate(10, func(i int) int { return 10 - i })},
		{name: "100", sli: slices.Generate(100, func(i int) int { return 100 - i })},
		{name: "1_000", sli: slices.Generate(1_000, func(i int) int { return 1_000 - i })},
		{name: "10_000", sli: slices.Generate(10_000, func(i int) int { return 10_000 - i })},
		{name: "100_000", sli: slices.Generate(100_000, func(i int) int { return 100_000 - i })},
		{name: "1_000_000", sli: slices.Generate(1_000_000, func(i int) int { return 1_000_000 - i })},
	}
}

func BenchmarkTopK_Add(b *testing.B) {
	for _, bm := range topKBenchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				top := streaming.NewTopKOrdered[int](16)
				for _, v := range bm.sli {
					top.Add(v)
				}
			}
		})
	}
}

func BenchmarkTopK_Result(b *testing.B) {
	for _, bm := range topKBenchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			top := streaming.NewTopKOrdered[int](16)
			for _, v := range bm.sli {
				top.Add(v)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = top.Result()
			}
		})
	}
}
