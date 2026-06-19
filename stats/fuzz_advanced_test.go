package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

// FuzzVectorOps asserts the algebraic invariants of the vector operations over
// finite, non-negative inputs (each byte 0-255): Dot and EuclideanDistance are
// symmetric, a vector's distance to itself is zero, Norm is non-negative, and a
// defined CosineSimilarity stays within [-1, 1].
func FuzzVectorOps(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4})
	f.Add([]byte{0, 0, 5, 5})
	f.Add([]byte{255, 0, 1, 128})

	f.Fuzz(func(t *testing.T, data []byte) {
		all := bytesToFloats(data)
		// Split into two equal-length vectors; an odd trailing byte is dropped.
		half := len(all) / 2
		if half == 0 {
			return
		}
		a := all[:half]
		b := all[half : 2*half]

		dotAB, okAB := stats.Dot(a, b)
		dotBA, okBA := stats.Dot(b, a)
		if okAB != okBA || dotAB != dotBA {
			t.Fatalf("Dot not symmetric: (%v,%v) vs (%v,%v)", dotAB, okAB, dotBA, okBA)
		}

		dist, ok := stats.EuclideanDistance(a, b)
		if !ok || dist < 0 {
			t.Fatalf("EuclideanDistance = %v, %v; want >= 0, true", dist, ok)
		}
		distRev, _ := stats.EuclideanDistance(b, a)
		if dist != distRev {
			t.Fatalf("EuclideanDistance not symmetric: %v vs %v", dist, distRev)
		}
		self, _ := stats.EuclideanDistance(a, a)
		if self != 0 {
			t.Fatalf("EuclideanDistance(a, a) = %v, want 0", self)
		}

		norm, ok := stats.Norm(a)
		if !ok || norm < 0 {
			t.Fatalf("Norm = %v, %v; want >= 0, true", norm, ok)
		}

		// CosineSimilarity is undefined for a zero vector (ok == false); when
		// defined it is a cosine and must lie within [-1, 1] (allowing a little
		// floating-point slack).
		cos, ok := stats.CosineSimilarity(a, b)
		if ok {
			if cos < -1-1e-9 || cos > 1+1e-9 {
				t.Fatalf("CosineSimilarity = %v, outside [-1, 1]", cos)
			}
		}
	})
}

// FuzzDistribution asserts the invariants of the categorical distribution
// measures: Shannon entropy is in [0, log2(k)] and Gini impurity is in [0, 1),
// for k distinct values, and both are exactly zero only when the sample is pure.
func FuzzDistribution(f *testing.F) {
	f.Add([]byte{1, 1, 1})
	f.Add([]byte{1, 2, 3, 4})
	f.Add([]byte{7, 7, 8, 8})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 {
			_, ok := stats.Entropy(data)
			if ok {
				t.Fatalf("Entropy(empty) reported ok")
			}
			_, ok = stats.Gini(data)
			if ok {
				t.Fatalf("Gini(empty) reported ok")
			}
			return
		}

		distinct := map[byte]struct{}{}
		for _, v := range data {
			distinct[v] = struct{}{}
		}
		k := len(distinct)

		h, ok := stats.Entropy(data)
		if !ok {
			t.Fatalf("Entropy reported not-ok for non-empty byte input")
		}
		if h < 0 || h > math.Log2(float64(k))+1e-9 {
			t.Fatalf("Entropy = %v, outside [0, log2(%d)=%v]", h, k, math.Log2(float64(k)))
		}

		g, ok := stats.Gini(data)
		if !ok {
			t.Fatalf("Gini reported not-ok for non-empty byte input")
		}
		if g < 0 || g >= 1 {
			t.Fatalf("Gini = %v, outside [0, 1)", g)
		}

		// A pure (single-value) sample has zero uncertainty in both measures;
		// any diversity makes both strictly positive.
		if k == 1 && (h != 0 || g != 0) {
			t.Fatalf("pure sample: Entropy=%v Gini=%v, want 0/0", h, g)
		}
		if k > 1 && (h <= 0 || g <= 0) {
			t.Fatalf("diverse sample (k=%d): Entropy=%v Gini=%v, want > 0", k, h, g)
		}
	})
}
