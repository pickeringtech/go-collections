package stats_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/stats"
)

// Example_quickStart is the runnable twin of the package godoc overview. Keep
// the two in sync: `go test` compiles and output-checks this, which is what
// guarantees the documented stats API actually exists and behaves as shown.
func Example_quickStart() {
	xs := []int{3, 1, 4, 1, 5}

	prod, _ := stats.Product(xs)
	med, _ := stats.Median(xs)
	rng, _ := stats.Range(xs)
	modes, _ := stats.Mode(xs)
	lo, hi, _ := stats.MinMax(xs)
	cum := stats.CumulativeSum(xs)
	clamped := stats.Clamp(7, 0, 5)

	fmt.Printf("product=%d median=%.1f range=%d modes=%v min=%d max=%d cum=%v clamp=%d",
		prod, med, rng, modes, lo, hi, cum, clamped)
	// Output: product=60 median=3.0 range=4 modes=[1] min=1 max=5 cum=[3 4 8 9 14] clamp=5
}

func ExampleProduct() {
	p, ok := stats.Product([]int{2, 3, 4})
	fmt.Println(p, ok)
	// Output: 24 true
}

func ExampleMedian() {
	odd, _ := stats.Median([]int{3, 1, 2})
	even, _ := stats.Median([]int{1, 2, 3, 4})
	fmt.Printf("%.1f %.1f", odd, even)
	// Output: 2.0 2.5
}

func ExampleMode() {
	// Ties are returned in order of first appearance.
	modes, _ := stats.Mode([]int{3, 1, 3, 1, 2})
	fmt.Println(modes)
	// Output: [3 1]
}

func ExampleMinMax() {
	lo, hi, ok := stats.MinMax([]int{3, 1, 4, 1, 5})
	fmt.Println(lo, hi, ok)
	// Output: 1 5 true
}

func ExampleArgMin() {
	// Index of the smallest element; ties resolve to the first occurrence.
	i, _ := stats.ArgMin([]int{3, 1, 4, 1, 5})
	fmt.Println(i)
	// Output: 1
}

func ExampleArgMax() {
	i, _ := stats.ArgMax([]int{3, 1, 4, 1, 5})
	fmt.Println(i)
	// Output: 4
}

func ExampleRange() {
	r, ok := stats.Range([]int{3, 1, 4, 1, 5})
	fmt.Println(r, ok)
	// Output: 4 true
}

func ExampleCumulativeSum() {
	fmt.Println(stats.CumulativeSum([]int{3, 1, 4, 1, 5}))
	// Output: [3 4 8 9 14]
}

func ExampleClamp() {
	fmt.Println(stats.Clamp(-2, 0, 5), stats.Clamp(3, 0, 5), stats.Clamp(7, 0, 5))
	// Output: 0 3 5
}

func ExampleClampAll() {
	fmt.Println(stats.ClampAll([]int{-3, 0, 2, 9}, 0, 5))
	// Output: [0 0 2 5]
}
