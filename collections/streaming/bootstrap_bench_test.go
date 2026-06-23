package streaming_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/streaming"
	"github.com/pickeringtech/go-collections/slices"
)

// bootstrapBenchmarkInputs returns the scaling ladder of input lengths used by
// the bootstrap benchmarks.
func bootstrapBenchmarkInputs() []struct {
	name string
	sli  []int
} {
	return []struct {
		name string
		sli  []int
	}{
		{name: "3", sli: slices.Generate(3, slices.NumericIdentityGenerator[int])},
		{name: "10", sli: slices.Generate(10, slices.NumericIdentityGenerator[int])},
		{name: "100", sli: slices.Generate(100, slices.NumericIdentityGenerator[int])},
		{name: "1_000", sli: slices.Generate(1_000, slices.NumericIdentityGenerator[int])},
		{name: "10_000", sli: slices.Generate(10_000, slices.NumericIdentityGenerator[int])},
		{name: "100_000", sli: slices.Generate(100_000, slices.NumericIdentityGenerator[int])},
		{name: "1_000_000", sli: slices.Generate(1_000_000, slices.NumericIdentityGenerator[int])},
	}
}

func BenchmarkBootstrap(b *testing.B) {
	for _, bm := range bootstrapBenchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			rng := streaming.NewRand(0)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = streaming.Bootstrap(bm.sli, rng)
			}
		})
	}
}

func BenchmarkBootstrapN(b *testing.B) {
	for _, bm := range bootstrapBenchmarkInputs() {
		b.Run(bm.name, func(b *testing.B) {
			rng := streaming.NewRand(0)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = streaming.BootstrapN(bm.sli, 10, rng)
			}
		})
	}
}
