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
//	j, ok := similarity.Jaccard(s1, s2)  // 3/5 = 0.6, true
//	d, ok := similarity.Dice(s1, s2)     // 2*3/(4+4) = 0.75, true
//	o, ok := similarity.Overlap(s1, s2)  // 3/4 = 0.75, true
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
//   - Dice: 2|A∩B| / (|A|+|B|) — weights intersection more than Jaccard.
//   - Overlap: |A∩B| / min(|A|,|B|) — measures containment; 1 when the
//     smaller set is a subset of the larger.
//
// Like the vector similarities, all three return (float64, bool): ok is false
// when the denominator is zero — both sets empty for Jaccard and Dice, either
// set empty for Overlap — so an undefined coefficient is never confused with
// the computed 0 of two disjoint sets.
package similarity
