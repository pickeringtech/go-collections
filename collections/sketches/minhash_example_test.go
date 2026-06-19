package sketches_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/collections/sketches"
)

func ExampleNewMinHash() {
	// nil rng → deterministic default seed.
	m := sketches.NewMinHash[string](128, nil)
	m.Add("hello")
	m.Add("world")

	sig := m.Signature()
	fmt.Println(len(sig))
	// Output: 128
}

func ExampleEstimatedJaccard() {
	// Two identical sets have Jaccard similarity 1.
	a := sketches.NewMinHash[string](256, nil)
	b := sketches.NewMinHash[string](256, nil)

	for _, word := range []string{"the", "quick", "brown", "fox"} {
		a.Add(word)
		b.Add(word)
	}

	est, ok := sketches.EstimatedJaccard(a, b)
	fmt.Printf("%.1f %v", est, ok)
	// Output: 1.0 true
}

func ExampleEstimatedJaccard_mismatch() {
	// Sketches with different numHashes are incompatible.
	a := sketches.NewMinHash[string](64, nil)
	b := sketches.NewMinHash[string](128, nil)

	_, ok := sketches.EstimatedJaccard(a, b)
	fmt.Println(ok)
	// Output: false
}
