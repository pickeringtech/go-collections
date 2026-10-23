package streaming_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/streaming"
	"github.com/pickeringtech/go-collections/slices"
)

// reservoirBenchmarkInputs returns the scaling ladder of stream lengths. Every
// element past the first k is a replacement candidate, so longer streams
// exercise the steady-state Add path.
func reservoirBenchmarkInputs() []struct {
	name string
	sli  []int
} {
	return []struct {
		name string
		sli  []int
	}{
		{name: "10", sli: slices.Generate(10, func(i int) int { return i })},
		{name: "100", sli: slices.Generate(100, func(i int) int { return i })},
		{name: "1_000", sli: slices.Generate(1_000, func(i int) int { return i })},
		{name: "10_000", sli: slices.Generate(10_000, func(i int) int { return i })},
		{name: "100_000", sli: slices.Generate(100_000, func(i int) int { return i })},
		{name: "1_000_000", sli: slices.Generate(1_000_000, func(i int) int { return i })},
	}
}

func BenchmarkReservoir_Add(b *testing.B) {
	for _, bm := range reservoirBenchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := streaming.NewReservoir[int](16, streaming.NewRand(1))
				for _, v := range bm.sli {
					r.Add(v)
				}
			}
		})
	}
}

func BenchmarkWeightedReservoir_Add(b *testing.B) {
	for _, bm := range reservoirBenchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := streaming.NewWeightedReservoir[int](16, streaming.NewRand(1))
				for _, v := range bm.sli {
					r.Add(v, float64(v)+1)
				}
			}
		})
	}
}
