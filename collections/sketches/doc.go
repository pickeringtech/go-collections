// Package sketches provides probabilistic data-sketching structures for
// approximate set operations. Sketches trade exact answers for very low memory
// and sub-linear time — useful when exact set sizes or intersections are
// prohibitively large.
//
// # Quick Start
//
//	import (
//	    "github.com/pickeringtech/go-collections/collections/sketches"
//	    "math/rand/v2"
//	)
//
//	// Build two MinHash sketches (nil rng → deterministic default seed).
//	a := sketches.NewMinHash[string](128, nil)
//	b := sketches.NewMinHash[string](128, nil)
//
//	for _, word := range []string{"the", "quick", "brown", "fox"} {
//	    a.Add(word)
//	}
//	for _, word := range []string{"the", "lazy", "brown", "dog"} {
//	    b.Add(word)
//	}
//
//	// Estimate the Jaccard similarity of the two sets.
//	est, ok := sketches.EstimatedJaccard(a, b)
//	// ok == true; est ≈ 0.5 (exact Jaccard = 2/6 ≈ 0.333; MinHash has variance)
//	_ = est
//	_ = ok
//
// # MinHash
//
// MinHash[T] estimates the Jaccard similarity between two sets by maintaining
// a fixed-size signature of uint64 hash minimums — one per hash function in
// the permutation family. Adding the same elements to two MinHash sketches
// configured with the same numHashes and the same rng produces comparable
// signatures whose collision rate approximates the Jaccard coefficient.
//
// The accuracy of the estimate improves with the number of hash functions
// (numHashes): 128 gives roughly ±7% error at 95% confidence; 256 gives ±5%.
//
// # Seeding and reproducibility
//
// NewMinHash accepts a *rand.Rand (from math/rand/v2) as its last parameter.
// Passing nil selects a hard-coded deterministic default seed, so the same
// numHashes always produces the same permutation family. This makes sketches
// portable across process restarts and across goroutines that build independent
// sketches that will later be compared.
//
// # Goroutine safety
//
// MinHash is stateful and NOT goroutine-safe. Do not call Add or Signature
// concurrently on the same instance. A thread-safe ConcurrentMinHash variant
// is planned for a later issue.
//
// # EstimatedJaccard
//
// EstimatedJaccard compares two sketches' signatures element-wise and returns
// the fraction of positions where the minimum hash values agree. It returns
// ok == false when the two sketches were built with different numHashes or
// different permutation families (mismatched construction parameters).
package sketches
