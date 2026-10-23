// Package sketches provides probabilistic data-sketching structures — compact
// summaries that trade exact answers for very low memory and sub-linear time,
// useful when exact computation over large sets or streams is prohibitive.
//
// The package itself holds MinHash, a set-similarity sketch (documented below).
// The streaming sketches each live in their own sub-package:
//
//   - [bloom] — approximate set membership. "Have I seen this before?" with no
//     false negatives and a tunable false-positive rate, in bits per element
//     independent of element size.
//   - [countmin] — approximate frequency counts. "How often have I seen this?"
//     never under-reporting, in space independent of the stream's cardinality.
//   - [hll] — approximate cardinality. "How many distinct things have I seen?"
//     in a few kilobytes, even for billions of distinct items.
//   - [tdigest] — approximate quantiles. "What's the 99th-percentile latency?"
//     over an unbounded stream, in memory set by a compression parameter.
//
// The streaming sketches share a common design: bounded, configurable accuracy
// with documented error bounds; mergeability (Merge) for parallel and
// distributed aggregation; and a delegating, read-write-mutex-guarded Concurrent
// variant alongside a plain type that is not safe for concurrent use. The
// comparable-typed sketches (bloom, countmin, hll) additionally take seeded,
// pluggable hashing (WithSeed/WithHasher) so behaviour is reproducible; tdigest
// is float64-only (it has no keys to hash) and approximates the value domain of
// the stats quantile functions rather than being generic over an element type.
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
// numHashes always produces the same permutation family within a program run.
// This makes sketches comparable across goroutines that build independent
// sketches in the same process. Note that element hashing uses a process-local
// seed (maphash.MakeSeed), so sketches are not portable across process restarts
// — two runs with identical elements will produce different signatures.
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
//
// [bloom]: https://pkg.go.dev/github.com/pickeringtech/go-collections/collections/sketches/bloom
// [countmin]: https://pkg.go.dev/github.com/pickeringtech/go-collections/collections/sketches/countmin
// [hll]: https://pkg.go.dev/github.com/pickeringtech/go-collections/collections/sketches/hll
// [tdigest]: https://pkg.go.dev/github.com/pickeringtech/go-collections/collections/sketches/tdigest
package sketches
