package stats_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/stats"
)

func ExampleDot() {
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}

	dot, ok := stats.Dot(a, b)
	fmt.Printf("%.1f %v", dot, ok)
	// Output: 32.0 true
}

func ExampleDot_mismatch() {
	// Length mismatch makes the dot product undefined.
	_, ok := stats.Dot([]float64{1, 2, 3}, []float64{1, 2})
	fmt.Println(ok)
	// Output: false
}

func ExampleNorm() {
	// Classic 3-4-5 right triangle — the hypotenuse is 5.
	v := []float64{3, 4}

	n, ok := stats.Norm(v)
	fmt.Printf("%.1f %v", n, ok)
	// Output: 5.0 true
}

func ExampleNorm_empty() {
	// The norm of an empty vector is undefined.
	_, ok := stats.Norm([]float64{})
	fmt.Println(ok)
	// Output: false
}
