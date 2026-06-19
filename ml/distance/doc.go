// Package distance provides distance metrics for vectors and sequences.
// All DISTANCE functions follow the convention: lower values mean items are
// more alike (closer), and zero means identical.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/ml/distance"
//
//	a := []float64{1, 2, 3}
//	b := []float64{4, 6, 3}
//
//	// Continuous vector distances
//	e, ok := distance.Euclidean(a, b)    // 5.0, true
//	m, ok := distance.Manhattan(a, b)    // 7.0, true
//	w, ok := distance.Minkowski(a, b, 3) // ~3.27, true
//	c, ok := distance.CosineDistance(a, b) // 1 - cosine_similarity, true
//
//	// Discrete / sequence distances
//	h, ok := distance.Hamming([]string{"a","b","c"}, []string{"a","x","c"}) // 1, true
//	lev := distance.Levenshtein("kitten", "sitting") // 3
//
//	_ = e
//	_ = m
//	_ = w
//	_ = c
//	_ = h
//	_ = lev
//
// # Distance functions (DISTANCE — lower = closer)
//
// Continuous vector distances — accept any constraints.Numeric element type,
// return (float64, bool). All return ok == false for empty input or
// mismatched-length vectors:
//
//   - Euclidean: the straight-line L2 distance between two points.
//   - Manhattan: the L1 (taxicab) distance — sum of absolute differences.
//   - Minkowski: the Lp generalisation; p<1 returns ok==false (not a metric).
//   - CosineDistance: 1 minus the cosine similarity; in [0, 2]; delegates to
//     ml/similarity.CosineSimilarity.
//
// Discrete distances:
//
//   - Hamming: the number of positions at which two equal-length sequences
//     differ. Accepts any comparable element type; returns (int, bool).
//   - Levenshtein: the minimum number of single-character edits (insert,
//     delete, substitute) to turn string a into string b. Operates over runes
//     so it handles multi-byte Unicode correctly.
//
// # NaN/Inf policy
//
// Non-finite inputs propagate to non-finite results with ok == true for all
// continuous vector distances, matching the stats package's NaN/Inf policy.
package distance
