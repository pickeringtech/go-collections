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
//	lev := distance.StringLevenshtein("kitten", "sitting") // 3
//	tok := distance.Levenshtein([]string{"a", "b"}, []string{"a", "c"}) // 1
//
//	_ = e
//	_ = m
//	_ = w
//	_ = c
//	_ = h
//	_ = lev
//	_ = tok
//
// # Distance functions (DISTANCE — lower = closer)
//
// Continuous vector distances — accept any constraints.Numeric element type,
// return (float64, bool). All return ok == false for empty input or
// mismatched-length vectors:
//
//   - Euclidean: the straight-line L2 distance between two points.
//   - Manhattan: the L1 (taxicab) distance — sum of absolute differences.
//   - Minkowski: the Lp generalisation; a non-finite p or p<1 returns
//     ok==false (not a valid metric parameter).
//   - CosineDistance: 1 minus the cosine similarity; in [0, 2]; delegates to
//     ml/similarity.CosineSimilarity.
//
// Discrete distances:
//
//   - Hamming: the number of positions at which two equal-length sequences
//     differ. Accepts any comparable element type; returns (int, bool).
//   - Levenshtein: the minimum number of single-element edits (insert,
//     delete, substitute) to turn sequence a into sequence b. Accepts any
//     comparable element type; returns int (always defined).
//   - StringLevenshtein: Levenshtein over the runes of two strings, so it
//     handles multi-byte Unicode correctly.
//
// # NaN/Inf policy
//
// Non-finite inputs propagate to non-finite results with ok == true for all
// continuous vector distances, matching the stats package's NaN/Inf policy.
package distance
