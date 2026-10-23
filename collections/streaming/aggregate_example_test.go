package streaming_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/collections/streaming"
)

func ExampleRunningMean() {
	mean := streaming.NewRunningMean()
	for _, v := range []float64{2, 4, 6, 8} {
		mean.Add(v)
	}

	avg, ok := mean.Result()
	fmt.Println(avg, ok, mean.Count())
	// Output: 5 true 4
}

func ExampleRunningVariance() {
	rv := streaming.NewRunningVariance()
	for _, v := range []float64{2, 4, 4, 4, 5, 5, 7, 9} {
		rv.Add(v)
	}

	sample, _ := rv.SampleVariance()
	pop, _ := rv.PopulationVariance()
	fmt.Printf("sample=%.4f population=%.4f\n", sample, pop)
	// Output: sample=4.5714 population=4.0000
}

func ExampleEWMA() {
	ewma := streaming.NewEWMA(0.5)
	for _, v := range []float64{10, 20, 30} {
		ewma.Add(v)
	}

	avg, ok := ewma.Result()
	fmt.Println(avg, ok)
	// Output: 22.5 true
}

func ExampleRunningMinMax() {
	mm := streaming.NewRunningMinMax[int]()
	for _, v := range []int{5, 1, 9, 3, 7} {
		mm.Add(v)
	}

	lo, hi, ok := mm.Result()
	fmt.Println(lo, hi, ok)
	// Output: 1 9 true
}
