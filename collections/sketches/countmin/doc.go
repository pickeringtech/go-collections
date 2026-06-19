// Package countmin provides a generic Count-Min sketch: a structure that tracks
// approximate frequency counts over a stream in space far smaller than an exact
// counter table would need.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/collections/sketches/countmin"
//
//	s, err := countmin.New[string](0.001, 0.01) // overshoot ≤ 0.001·N, 99% confidence
//	if err != nil {
//		// only on a nonsensical epsilon / delta
//	}
//	s.Add("/index.html")
//	s.AddCount("/api/v1", 42)
//	s.Estimate("/index.html") // approximate hit count, never an under-count
//
// # The guarantee
//
// Estimate never under-reports: it returns the true count plus a one-sided
// error. With the sizing below, the overshoot is at most epsilon·N (N being the
// total count added) with probability at least 1-delta. This makes Count-Min a
// good fit for heavy-hitter and frequency questions over streams too large to
// count exactly, where over-counting a rare item is acceptable but missing a
// frequent one is not.
//
// # Accuracy and memory
//
// New sizes the counter table from the two error parameters:
//
//	w = ceil(e/epsilon)   columns per row
//	d = ceil(ln(1/delta)) rows
//
// The table holds w·d 64-bit counters regardless of how many distinct elements
// pass through — memory depends only on the target accuracy, not the stream's
// cardinality. Tightening epsilon widens each row; tightening delta adds rows.
// Counters saturate at the maximum uint64 rather than overflowing.
//
// # Hashing
//
// Hashing is seeded and deterministic. Each of the d rows is indexed by an
// independent column derived from two base hashes via Kirsch–Mitzenmacher
// double hashing. Use WithSeed to vary the seed and WithHasher to supply a
// custom hash for exotic key types or cross-process reproducibility.
//
// # Mergeability
//
// Merge adds one sketch's counters into another element-wise, so per-shard or
// per-worker sketches combine into one covering the whole stream — the basis for
// parallel and distributed aggregation. Both sketches must share dimensions and
// seed.
//
// # Thread safety
//
// Sketch is not safe for concurrent use. ConcurrentSketch wraps it with a
// read-write mutex — Estimate takes a read lock, the mutating operations the
// write lock. To merge a ConcurrentSketch into another, pass its Snapshot.
package countmin
