package distance_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/ml/distance"
)

func ExampleEuclidean() {
	// Classic 3-4-5 triangle.
	a := []float64{0, 0}
	b := []float64{3, 4}

	d, ok := distance.Euclidean(a, b)
	fmt.Printf("%.1f %v", d, ok)
	// Output: 5.0 true
}

func ExampleEuclidean_mismatch() {
	_, ok := distance.Euclidean([]float64{1, 2, 3}, []float64{1, 2})
	fmt.Println(ok)
	// Output: false
}

func ExampleManhattan() {
	// Taxicab distance from origin to (3, 4).
	a := []float64{0, 0}
	b := []float64{3, 4}

	d, ok := distance.Manhattan(a, b)
	fmt.Printf("%.1f %v", d, ok)
	// Output: 7.0 true
}

func ExampleMinkowski() {
	// p=2 equals Euclidean distance.
	a := []float64{0, 0}
	b := []float64{3, 4}

	d, ok := distance.Minkowski(a, b, 2)
	fmt.Printf("%.1f %v", d, ok)
	// Output: 5.0 true
}

func ExampleMinkowski_pLessThanOne() {
	// p < 1 is not a valid metric — returns ok == false.
	_, ok := distance.Minkowski([]float64{1}, []float64{2}, 0.5)
	fmt.Println(ok)
	// Output: false
}

func ExampleCosineDistance() {
	// Orthogonal vectors have cosine similarity 0, so cosine distance = 1.
	a := []float64{1, 0}
	b := []float64{0, 1}

	d, ok := distance.CosineDistance(a, b)
	fmt.Printf("%.1f %v", d, ok)
	// Output: 1.0 true
}

func ExampleHamming() {
	a := []string{"a", "b", "c"}
	b := []string{"a", "x", "c"}

	d, ok := distance.Hamming(a, b)
	fmt.Printf("%d %v", d, ok)
	// Output: 1 true
}

func ExampleHamming_mismatch() {
	_, ok := distance.Hamming([]int{1, 2}, []int{1, 2, 3})
	fmt.Println(ok)
	// Output: false
}

func ExampleLevenshtein() {
	// Edit distance over word tokens: substitute "quick" → "slow" and
	// "fox" → "dog", then insert "!".
	a := []string{"the", "quick", "brown", "fox"}
	b := []string{"the", "slow", "brown", "dog", "!"}
	fmt.Println(distance.Levenshtein(a, b))
	// Output: 3
}

func ExampleLevenshtein_bytes() {
	// Over raw bytes, the two-byte UTF-8 "é" costs two edits; StringLevenshtein
	// counts it as one character.
	fmt.Println(distance.Levenshtein([]byte("café"), []byte("cafe")), distance.StringLevenshtein("café", "cafe"))
	// Output: 2 1
}

func ExampleStringLevenshtein() {
	d := distance.StringLevenshtein("kitten", "sitting")
	fmt.Println(d)
	// Output: 3
}

func ExampleStringLevenshtein_identical() {
	d := distance.StringLevenshtein("hello", "hello")
	fmt.Println(d)
	// Output: 0
}

func ExampleStringLevenshtein_empty() {
	// Distance from empty string to "abc" is 3 (3 insertions).
	d := distance.StringLevenshtein("", "abc")
	fmt.Println(d)
	// Output: 3
}
