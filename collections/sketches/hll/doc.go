// Package hll provides a generic HyperLogLog: it estimates the number of
// distinct elements (cardinality) of a stream using a fixed, tiny amount of
// memory — a few kilobytes for billions of distinct items.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/collections/sketches/hll"
//
//	s, err := hll.New[string]() // default precision 14 (~16 KB, ~0.81% error)
//	if err != nil {
//		// only on an out-of-range WithPrecision
//	}
//	for _, visitor := range stream {
//		s.Add(visitor)
//	}
//	s.Count() // estimated distinct visitors
//
// # How it works
//
// Each element is hashed; the top p bits choose one of m = 2^p registers and the
// remaining bits contribute the length of their leading zero run. Each register
// keeps the longest run it has seen, and the harmonic mean of the registers,
// scaled by a bias-correction constant, estimates the cardinality. Re-adding an
// element never changes the estimate, so Count reflects distinct elements only.
//
// # Accuracy and memory
//
// Precision p (see WithPrecision) trades memory for accuracy: m = 2^p registers
// of one byte each, with a standard error of about 1.04/sqrt(m). The default
// precision 14 gives 16384 registers (~16 KB) and roughly 0.81% error. Memory is
// fixed by p alone — it does not grow with the stream — and StandardError
// reports the expected relative error. For small cardinalities the estimator
// switches to linear counting, which is more accurate while many registers are
// still empty.
//
// # Hashing
//
// Hashing is seeded and deterministic. HyperLogLog reads both the high index
// bits and the low-bit zero-run, so a well-distributed hash matters; the default
// hasher applies a strong finalizer for this reason. Use WithSeed to vary the
// seed and WithHasher for a custom hash. Two sketches must share precision and
// seed to be merged.
//
// # Mergeability
//
// Merge takes the register-wise maximum of two sketches, yielding a sketch that
// estimates the cardinality of the union — the basis for parallel and
// distributed distinct-counting. Merging is exact with respect to the union: the
// merged registers equal those of a single sketch fed both streams.
//
// # Thread safety
//
// Sketch is not safe for concurrent use. ConcurrentSketch wraps it with a
// read-write mutex — Count takes a read lock, Add/Merge/Clear the write lock. To
// merge a ConcurrentSketch into another, pass its Snapshot.
package hll
