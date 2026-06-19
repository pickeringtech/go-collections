// Package sketches is the umbrella for the library's probabilistic data
// structures — compact sketches that answer questions about massive streams in
// bounded memory, trading a little accuracy for a lot of space. Each sketch
// lives in its own sub-package; this package holds only the shared overview.
//
// # The family
//
//   - [bloom] — approximate set membership. "Have I seen this before?" with no
//     false negatives and a tunable false-positive rate, in bits per element
//     independent of element size.
//   - [countmin] — approximate frequency counts. "How often have I seen this?"
//     never under-reporting, in space independent of the stream's cardinality.
//   - [hll] — approximate cardinality. "How many distinct things have I seen?"
//     in a few kilobytes, even for billions of distinct items.
//
// (A streaming-quantiles sketch — t-digest — is co-designed with the stats
// quantile work and lives with that package rather than here.)
//
// # Shared design
//
// Every sketch in this family follows the same conventions:
//
//   - Bounded, configurable accuracy. Construction takes the error target
//     (false-positive rate, epsilon/delta, or precision); the documented bound
//     and the memory cost follow from it. Constructors return an error rather
//     than panic on nonsensical configuration.
//   - Seeded, pluggable hashing. Hashing is deterministic from a stored seed, so
//     behaviour is reproducible and golden tests are stable. WithSeed varies it
//     and WithHasher overrides it for custom key types.
//   - Mergeable. Merge combines two compatible sketches into one covering the
//     union of their inputs — the basis for parallel and distributed
//     aggregation. Compatibility (matching configuration and seed) is checked.
//   - A concurrent variant. Each plain type owns all the logic; the matching
//     Concurrent type wraps it with a read-write mutex and delegates, so there
//     is a single implementation to reason about. The plain types are not safe
//     for concurrent use.
//
// [bloom]: https://pkg.go.dev/github.com/pickeringtech/go-collections/collections/sketches/bloom
// [countmin]: https://pkg.go.dev/github.com/pickeringtech/go-collections/collections/sketches/countmin
// [hll]: https://pkg.go.dev/github.com/pickeringtech/go-collections/collections/sketches/hll
package sketches
