package bloom

import (
	"errors"
	"fmt"
	"math"
	"math/bits"

	"github.com/pickeringtech/go-collections/collections/sketches/internal/sketchhash"
)

// ErrInvalidConfig is returned by New when the requested capacity or
// false-positive rate is out of range. It is also the root of the errors
// returned by Merge on a configuration mismatch, so callers can test with
// errors.Is.
var ErrInvalidConfig = errors.New("bloom: invalid configuration")

// Filter is a Bloom filter: a compact, probabilistic set that answers
// membership queries with a tunable false-positive rate and, crucially, no
// false negatives — if Contains returns false the element was definitely never
// added. Adding an element never removes another, and the structure cannot
// shrink, so a Bloom filter trades the ability to enumerate or delete for
// memory that is a small constant per element regardless of element size.
//
// A Filter is not safe for concurrent use. Wrap one with NewConcurrent for a
// goroutine-safe variant.
//
// The zero value is not usable; construct a Filter with New.
type Filter[T comparable] struct {
	bitsArr []uint64 // bit array, packed 64 bits per word
	m       uint64   // number of bits
	k       uint64   // number of hash functions
	seed    uint64   // hashing seed; filters must share it to be merged
	hasher  func(seed uint64, value T) uint64
	added   uint64 // count of Add calls (used for load reporting)
}

// Interface guard.
var _ Membership[string] = (*Filter[string])(nil)

// New creates a Bloom filter sized for expectedItems insertions at the target
// falsePositiveRate (for example 0.01 for 1%). It returns an error wrapping
// ErrInvalidConfig if expectedItems <= 0 or falsePositiveRate is not in the
// open interval (0, 1).
//
// The actual memory used is optimalBits(expectedItems, falsePositiveRate) bits,
// rounded up to a whole 64-bit word. Inserting more than expectedItems elements
// is allowed but degrades the false-positive rate beyond the target.
func New[T comparable](expectedItems int, falsePositiveRate float64, opts ...Option[T]) (*Filter[T], error) {
	if expectedItems <= 0 {
		return nil, fmt.Errorf("%w: expectedItems must be positive, got %d", ErrInvalidConfig, expectedItems)
	}
	if falsePositiveRate <= 0 || falsePositiveRate >= 1 || math.IsNaN(falsePositiveRate) {
		return nil, fmt.Errorf("%w: falsePositiveRate must be in (0,1), got %v", ErrInvalidConfig, falsePositiveRate)
	}

	m := optimalBits(expectedItems, falsePositiveRate)
	k := optimalHashes(m, uint64(expectedItems))

	f := &Filter[T]{
		bitsArr: make([]uint64, (m+63)/64),
		m:       m,
		k:       k,
		hasher:  sketchhash.Hash64[T],
	}
	for _, opt := range opts {
		opt(f)
	}
	return f, nil
}

// Option configures a Filter at construction time.
type Option[T comparable] func(*Filter[T])

// WithSeed sets the hashing seed. Two filters must share a seed (and the same
// capacity) to be merged. The default seed is 0, so filters built with the same
// capacity are mergeable out of the box.
func WithSeed[T comparable](seed uint64) Option[T] {
	return func(f *Filter[T]) { f.seed = seed }
}

// WithHasher overrides the hash function, for full control over hashing of
// custom key types or for cross-process reproducibility. The function must be
// deterministic in (seed, value).
func WithHasher[T comparable](hasher func(seed uint64, value T) uint64) Option[T] {
	return func(f *Filter[T]) {
		if hasher != nil {
			f.hasher = hasher
		}
	}
}

// Add inserts an element into the filter. Adding the same element twice is a
// no-op the second time.
func (f *Filter[T]) Add(value T) {
	h1, h2 := f.pair(value)
	for i := uint64(0); i < f.k; i++ {
		pos := (h1 + i*h2) % f.m
		f.bitsArr[pos>>6] |= 1 << (pos & 63)
	}
	f.added++
}

// Contains reports whether value may be in the set. A true result is
// probabilistic — it is correct except for a false-positive rate near the one
// the filter was sized for. A false result is exact: the element was never
// added.
func (f *Filter[T]) Contains(value T) bool {
	h1, h2 := f.pair(value)
	for i := uint64(0); i < f.k; i++ {
		pos := (h1 + i*h2) % f.m
		if f.bitsArr[pos>>6]&(1<<(pos&63)) == 0 {
			return false
		}
	}
	return true
}

// pair derives the two double-hashing lanes for value.
func (f *Filter[T]) pair(value T) (uint64, uint64) {
	h1 := f.hasher(f.seed, value)
	h2 := f.hasher(f.seed^0x9E3779B97F4A7C15, value) | 1
	return h1, h2
}

// Merge folds other into f by OR-ing their bit arrays, so f then reports
// membership for the union of both inputs. Both filters must have identical
// capacity (bit size and hash count) and seed; otherwise Merge returns an error
// wrapping ErrInvalidConfig and leaves f unchanged.
func (f *Filter[T]) Merge(other *Filter[T]) error {
	if other == nil {
		return fmt.Errorf("%w: cannot merge a nil filter", ErrInvalidConfig)
	}
	if f.m != other.m || f.k != other.k || f.seed != other.seed {
		return fmt.Errorf("%w: filters differ (m=%d/%d k=%d/%d seed=%d/%d)",
			ErrInvalidConfig, f.m, other.m, f.k, other.k, f.seed, other.seed)
	}
	for i := range f.bitsArr {
		f.bitsArr[i] |= other.bitsArr[i]
	}
	f.added += other.added
	return nil
}

// Clear removes every element, returning the filter to its empty state while
// keeping its capacity and seed.
func (f *Filter[T]) Clear() {
	for i := range f.bitsArr {
		f.bitsArr[i] = 0
	}
	f.added = 0
}

// BitSize returns the number of bits in the underlying array (m).
func (f *Filter[T]) BitSize() uint64 { return f.m }

// HashCount returns the number of hash functions used per element (k).
func (f *Filter[T]) HashCount() uint64 { return f.k }

// ApproxCount estimates the number of distinct elements added, from the count
// of set bits. It is approximate and saturates as the filter fills; for a
// filter loaded near or beyond its design capacity prefer a HyperLogLog.
func (f *Filter[T]) ApproxCount() uint64 {
	set := f.popcount()
	if set == 0 {
		return 0
	}
	if set >= f.m {
		// Fully saturated: the estimator diverges, so report capacity.
		return math.MaxUint64
	}
	// Swamidass & Baldi estimator: n ≈ -(m/k)·ln(1 - X/m).
	n := -(float64(f.m) / float64(f.k)) * math.Log(1-float64(set)/float64(f.m))
	return uint64(math.Round(n))
}

// EstimatedFalsePositiveRate returns the false-positive probability at the
// current fill level, computed from the fraction of bits set: (X/m)^k. This
// reflects how full the filter actually is, which rises above the design target
// once more than the expected number of elements have been added.
func (f *Filter[T]) EstimatedFalsePositiveRate() float64 {
	frac := float64(f.popcount()) / float64(f.m)
	return math.Pow(frac, float64(f.k))
}

// popcount returns the number of set bits across the whole array.
func (f *Filter[T]) popcount() uint64 {
	var n int
	for _, w := range f.bitsArr {
		n += bits.OnesCount64(w)
	}
	return uint64(n) // #nosec G115 -- popcount over m bits is in [0, m]; never negative
}

// optimalBits returns the bit-array size m that minimises memory for n items at
// false-positive rate p: m = ceil(-n·ln p / (ln 2)^2). With n >= 1 and
// p in (0,1) the numerator is positive, so m >= 1.
func optimalBits(n int, p float64) uint64 {
	return uint64(math.Ceil(-float64(n) * math.Log(p) / (math.Ln2 * math.Ln2)))
}

// optimalHashes returns the hash-function count k that minimises the
// false-positive rate for a filter of m bits holding n items:
// k = round((m/n)·ln 2), clamped to at least 1 (the rounding can reach 0 when
// the requested rate is very loose).
func optimalHashes(m, n uint64) uint64 {
	k := uint64(math.Round(float64(m) / float64(n) * math.Ln2))
	if k < 1 {
		return 1
	}
	return k
}
