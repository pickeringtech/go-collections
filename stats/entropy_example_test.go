package stats_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/stats"
)

func ExampleEntropy() {
	// A fair coin: two equiprobable outcomes carry exactly one bit of entropy.
	flips := []string{"heads", "tails", "heads", "tails"}

	h, ok := stats.Entropy(flips)
	fmt.Printf("%.1f bits %v", h, ok)
	// Output: 1.0 bits true
}

func ExampleGini() {
	// An evenly split label set has Gini impurity 1 - (0.5² + 0.5²) = 0.5.
	labels := []string{"spam", "ham", "spam", "ham"}

	g, ok := stats.Gini(labels)
	fmt.Printf("%.1f %v", g, ok)
	// Output: 0.5 true
}
