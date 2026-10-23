// Package similarity provides similarity metrics for vectors and sets —
// measures that return higher values for items that are more alike. It
// complements the ml/distance package, which returns lower values for items
// that are more alike.
//
// # Quick Start
//
//	import (
//	    "github.com/pickeringtech/go-collections/ml/similarity"
//	    "github.com/pickeringtech/go-collections/collections/sets"
//	)
//
//	// Vector similarity
//	a := []float64{1, 2, 3}
//	b := []float64{4, 5, 6}
//
//	dot, ok := similarity.DotProduct(a, b)     // 32.0, true
//	cos, ok := similarity.CosineSimilarity(a, b) // ~0.974, true
//
//	// Set similarity
//	s1 := sets.NewHash("a", "b", "c", "d")
//	s2 := sets.NewHash("b", "c", "d", "e")
//
//	j := similarity.Jaccard(s1, s2)  // 3/5 = 0.6
//	d := similarity.Dice(s1, s2)     // 2*3/(4+4) = 0.75
//	o := similarity.Overlap(s1, s2)  // 3/4 = 0.75
//
//	_ = dot
//	_ = cos
//	_ = j
//	_ = d
//	_ = o
//
// # Vector Similarity
//
// DotProduct delegates to stats.Dot, computing Σ aᵢ·bᵢ.
//
// CosineSimilarity measures the cosine of the angle between two vectors,
// normalised to [-1, 1]. It is independent of vector magnitude — only
// direction matters. A value of 1 means parallel (identical direction), 0 means
// perpendicular, and -1 means anti-parallel.
//
// Both functions follow the (float64, bool) idiom — ok == false for empty or
// mismatched-length inputs. CosineSimilarity also returns ok == false when
// either vector has zero magnitude (the cosine is undefined in that case).
// Non-finite inputs (NaN/Inf) propagate per the stats package's NaN/Inf policy.
//
// # Set Similarity
//
// Jaccard, Dice and Overlap all compose the Intersection and Length methods on
// the sets.Set[T] interface — |A∪B| is derived as |A|+|B|−|A∩B| rather than
// materialised, so no set algebra is reimplemented here.
//
//   - Jaccard: |A∩B| / |A∪B| — ranges from 0 (disjoint) to 1 (equal).
//     Empty∩Empty → 0.
//   - Dice: 2|A∩B| / (|A|+|B|) — weights intersection more than Jaccard.
//     Empty∩Empty → 0.
//   - Overlap: |A∩B| / min(|A|,|B|) — measures containment; 1 when the
//     smaller set is a subset of the larger. Zero denominator → 0.
package similarity
