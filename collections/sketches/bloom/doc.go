// Package bloom provides a generic Bloom filter: a compact, probabilistic set
// that answers membership queries in space independent of element size, at the
// cost of a tunable false-positive rate.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/collections/sketches/bloom"
//
//	f, err := bloom.New[string](1_000_000, 0.01) // 1M items, 1% false positives
//	if err != nil {
//		// only on a nonsensical capacity / rate
//	}
//	f.Add("alice")
//	f.Contains("alice") // true
//	f.Contains("carol") // false — never added, and Contains never lies about that
//
// # The guarantee
//
// A Bloom filter never returns a false negative: if Contains reports false, the
// element was definitely never added. It may return a false positive — report
// an element as present when it was not — with a probability close to the rate
// the filter was sized for. There is no way to remove an element or to list the
// contents; the structure only ever accumulates.
//
// # Accuracy and memory
//
// New sizes the filter from the two parameters using the standard optimal
// formulas:
//
//	m = ceil(-n·ln p / (ln 2)^2)   bits
//	k = round((m/n)·ln 2)          hash functions
//
// where n is the expected item count and p the target false-positive rate. The
// memory cost is therefore about -1.44·log2(p) bits per element — roughly 9.6
// bits per element at p=0.01 — regardless of how large each element is. Adding
// more than n elements is allowed but raises the false-positive rate above the
// target; EstimatedFalsePositiveRate reports the rate at the current fill level
// and ApproxCount estimates how many distinct elements have been added.
//
// # Hashing
//
// Hashing is seeded and deterministic, so a filter's behaviour is reproducible
// and two filters with the same seed and capacity are mergeable. k bit
// positions are derived from two base hashes by Kirsch–Mitzenmacher double
// hashing, which preserves the false-positive bound while computing only two
// hashes per operation. Use WithSeed to vary the seed (for example to make
// adversarial collisions hard) and WithHasher to supply a custom hash for exotic
// key types or cross-process reproducibility.
//
// # Mergeability
//
// Merge folds one filter into another by OR-ing their bit arrays, so a filter
// built per shard or per worker can be combined into one covering the union —
// the basis for parallel and distributed aggregation. Both filters must share
// capacity and seed.
//
// # Thread safety
//
// Filter is not safe for concurrent use. ConcurrentFilter wraps it with a
// read-write mutex — membership queries take a read lock, mutations the write
// lock — and is the variant to reach for when multiple goroutines touch the same
// filter. To merge a ConcurrentFilter into another, pass its Snapshot.
package bloom
