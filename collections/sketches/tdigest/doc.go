// Package tdigest provides a t-digest: a streaming, mergeable sketch that
// estimates quantiles (and the inverse CDF) of a stream of float64 values in
// memory that scales with a compression parameter and is independent of how
// many values have been seen.
//
// It is the approximate, bounded-memory counterpart of the exact
// [github.com/pickeringtech/go-collections/stats] quantile functions: where
// stats.Quantile sorts the whole sample, a Digest keeps only a small set of
// weighted centroids. Use it when the stream is too large to hold or sort, or
// when per-shard summaries must be merged.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/collections/sketches/tdigest"
//
//	d, _ := tdigest.New() // DefaultCompression (100)
//	for _, v := range stream {
//	    d.Add(v)
//	}
//	p99, ok := d.Percentile(99) // estimated 99th percentile
//	_ = p99
//	_ = ok
//
// # How it works
//
// A Digest summarises the distribution as centroids — weighted means of nearby
// points — kept sorted by value. Incoming points are buffered and, once enough
// accumulate, merged into the centroids in a single sorted sweep ("merging"
// t-digest). A scale function bounds each centroid's weight as a function of its
// quantile, so centroids stay small at the tails and may grow large in the
// middle. Quantile then walks the centroids, interpolating between their
// cumulative-weight midpoints (and between the extreme centroids and the exact
// observed min/max). No randomness is involved.
//
// # Accuracy and memory
//
// The retained centroid count is roughly proportional to the compression
// parameter and independent of the stream length, so memory is bounded no
// matter how long the stream runs. Higher compression keeps more centroids,
// costing memory but tightening the estimates. The scale function deliberately
// makes error smallest at the tails (q near 0 or 1, e.g. p99/p999 latency
// percentiles — the common case for monitoring) and largest in the middle of
// the distribution. q=0 and q=1 return the exact observed minimum and maximum.
//
// # Float64, not generic
//
// Unlike the comparable-typed sketches in the sibling packages (bloom, countmin,
// hll), Digest is deliberately not generic over a type parameter: it operates on
// float64 only, matching the value domain of the stats quantile functions it
// approximates. There is no meaningful "quantile" of an arbitrary comparable
// type, so a type parameter would add no value.
//
// # Mergeability
//
// Merge folds one digest's centroids into another, so per-shard or per-worker
// digests combine into one covering the union of all their streams — the basis
// for parallel and distributed quantile aggregation. Both digests must share the
// same compression; a mismatch (or a nil argument) returns an error wrapping
// ErrInvalidConfig and leaves the receiver unchanged.
//
// # Determinism and order-dependence
//
// A Digest uses no randomness, so a fixed sequence of operations always produces
// the same result. However, the retained centroids — and therefore the
// estimates — can depend on the order of Add and Merge calls: different orderings
// of the same data yield close but not bit-identical quantiles.
//
// # Thread safety
//
// Digest is not safe for concurrent use. ConcurrentDigest (via NewConcurrent)
// wraps it with a read-write mutex. To merge a ConcurrentDigest into another,
// pass its Snapshot.
package tdigest
