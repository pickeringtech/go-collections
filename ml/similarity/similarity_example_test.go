package similarity_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/collections/sets"
	"github.com/pickeringtech/go-collections/ml/similarity"
)

func ExampleDotProduct() {
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}

	dot, ok := similarity.DotProduct(a, b)
	fmt.Printf("%.1f %v", dot, ok)
	// Output: 32.0 true
}

func ExampleCosineSimilarity() {
	// Identical direction — cosine is 1 regardless of magnitude.
	a := []float64{1, 2, 3}
	b := []float64{2, 4, 6}

	cos, ok := similarity.CosineSimilarity(a, b)
	fmt.Printf("%.1f %v", cos, ok)
	// Output: 1.0 true
}

func ExampleCosineSimilarity_orthogonal() {
	// Orthogonal vectors — cosine is 0.
	a := []float64{1, 0}
	b := []float64{0, 1}

	cos, ok := similarity.CosineSimilarity(a, b)
	fmt.Printf("%.1f %v", cos, ok)
	// Output: 0.0 true
}

func ExampleCosineSimilarity_zeroVector() {
	// A zero vector has no direction — similarity is undefined.
	_, ok := similarity.CosineSimilarity([]float64{0, 0}, []float64{1, 2})
	fmt.Println(ok)
	// Output: false
}

func ExampleJaccard() {
	a := sets.NewHash("a", "b", "c", "d")
	b := sets.NewHash("b", "c", "d", "e")

	// |{b,c,d}| / |{a,b,c,d,e}| = 3/5 = 0.6
	j := similarity.Jaccard(a, b)
	fmt.Printf("%.1f", j)
	// Output: 0.6
}

func ExampleJaccard_empty() {
	// Both empty sets → Jaccard returns 0.
	j := similarity.Jaccard(sets.NewHash[string](), sets.NewHash[string]())
	fmt.Printf("%.1f", j)
	// Output: 0.0
}

func ExampleDice() {
	a := sets.NewHash("a", "b", "c", "d")
	b := sets.NewHash("b", "c", "d", "e")

	// 2*|{b,c,d}| / (4+4) = 6/8 = 0.75
	d := similarity.Dice(a, b)
	fmt.Printf("%.2f", d)
	// Output: 0.75
}

func ExampleOverlap() {
	// a is a subset of b — Overlap returns 1.
	a := sets.NewHash("a", "b")
	b := sets.NewHash("a", "b", "c")

	o := similarity.Overlap(a, b)
	fmt.Printf("%.1f", o)
	// Output: 1.0
}
