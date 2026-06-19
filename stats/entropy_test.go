package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func TestEntropy(t *testing.T) {
	t.Run("uniform identical sample has zero entropy", func(t *testing.T) {
		got, ok := stats.Entropy([]int{4, 4, 4, 4})
		if !ok || got != 0 {
			t.Fatalf("Entropy = %v, %v; want 0, true", got, ok)
		}
		// Guard against returning negative zero from -(p*log2 p) with p=1.
		if math.Signbit(got) {
			t.Errorf("Entropy returned negative zero")
		}
	})

	t.Run("two equiprobable values is one bit", func(t *testing.T) {
		got, ok := stats.Entropy([]int{0, 0, 1, 1})
		if !ok || !approxEqual(got, 1) {
			t.Fatalf("Entropy = %v, %v; want 1, true", got, ok)
		}
	})

	t.Run("four equiprobable values is two bits", func(t *testing.T) {
		got, ok := stats.Entropy([]string{"a", "b", "c", "d"})
		if !ok || !approxEqual(got, 2) {
			t.Fatalf("Entropy = %v, %v; want 2, true", got, ok)
		}
	})

	t.Run("empty is undefined", func(t *testing.T) {
		got, ok := stats.Entropy([]int{})
		if ok {
			t.Errorf("ok = true (%v), want false", got)
		}
	})

	t.Run("non-finite float is rejected", func(t *testing.T) {
		got, ok := stats.Entropy([]float64{1, math.NaN(), 2})
		if ok {
			t.Errorf("ok = true (%v), want false", got)
		}
	})
}

func TestGini(t *testing.T) {
	t.Run("pure sample has zero impurity", func(t *testing.T) {
		got, ok := stats.Gini([]int{4, 4, 4, 4})
		if !ok || got != 0 {
			t.Fatalf("Gini = %v, %v; want 0, true", got, ok)
		}
	})

	t.Run("two equiprobable values is one half", func(t *testing.T) {
		got, ok := stats.Gini([]int{0, 0, 1, 1})
		if !ok || !approxEqual(got, 0.5) {
			t.Fatalf("Gini = %v, %v; want 0.5, true", got, ok)
		}
	})

	t.Run("four equiprobable values is three quarters", func(t *testing.T) {
		got, ok := stats.Gini([]string{"a", "b", "c", "d"})
		if !ok || !approxEqual(got, 0.75) {
			t.Fatalf("Gini = %v, %v; want 0.75, true", got, ok)
		}
	})

	t.Run("empty is undefined", func(t *testing.T) {
		got, ok := stats.Gini([]int{})
		if ok {
			t.Errorf("ok = true (%v), want false", got)
		}
	})

	t.Run("non-finite float is rejected", func(t *testing.T) {
		got, ok := stats.Gini([]float64{1, math.Inf(1), 2})
		if ok {
			t.Errorf("ok = true (%v), want false", got)
		}
	})
}
