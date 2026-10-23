package streaming_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/streaming"
	"github.com/pickeringtech/go-collections/slices"
)

// aggregateBenchmarkInputs returns the scaling ladder of stream lengths used by
// the online-aggregate benchmarks.
func aggregateBenchmarkInputs() []struct {
	name string
	sli  []float64
} {
	gen := func(n int) []float64 {
		return slices.Generate(n, func(i int) float64 { return float64(i) })
	}
	return []struct {
		name string
		sli  []float64
	}{
		{name: "3", sli: gen(3)},
		{name: "10", sli: gen(10)},
		{name: "100", sli: gen(100)},
		{name: "1_000", sli: gen(1_000)},
		{name: "10_000", sli: gen(10_000)},
		{name: "100_000", sli: gen(100_000)},
		{name: "1_000_000", sli: gen(1_000_000)},
	}
}

func BenchmarkRunningMean_Add(b *testing.B) {
	for _, bm := range aggregateBenchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := streaming.NewRunningMean()
				for _, v := range bm.sli {
					m.Add(v)
				}
			}
		})
	}
}

func BenchmarkRunningVariance_Add(b *testing.B) {
	for _, bm := range aggregateBenchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v := streaming.NewRunningVariance()
				for _, x := range bm.sli {
					v.Add(x)
				}
			}
		})
	}
}

func BenchmarkEWMA_Add(b *testing.B) {
	for _, bm := range aggregateBenchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				e := streaming.NewEWMA(0.3)
				for _, v := range bm.sli {
					e.Add(v)
				}
			}
		})
	}
}

func BenchmarkRunningMinMax_Add(b *testing.B) {
	for _, bm := range aggregateBenchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				mm := streaming.NewRunningMinMax[float64]()
				for _, v := range bm.sli {
					mm.Add(v)
				}
			}
		})
	}
}
