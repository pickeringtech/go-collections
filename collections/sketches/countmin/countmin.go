package countmin

import (
	"errors"
	"fmt"
	"math"

	"github.com/pickeringtech/go-collections/collections/sketches/internal/sketchhash"
)

// ErrInvalidConfig is returned by New when the requested error bounds are out of
// range, and is the root of the error Merge returns on a configuration
// mismatch, so callers can test with errors.Is.
var ErrInvalidConfig = errors.New("countmin: invalid configuration")

// maxCounters caps the number of counters (w*d) New will allocate. Bounds tight
// enough to need more than this — from a minuscule epsilon or delta — are almost
// certainly a misconfiguration, and the unbounded product risks integer
// overflow, so New rejects them. 1<<48 counters is 2 PiB, far beyond any real
// in-memory sketch.
const maxCounters = 1 << 48

// Sketch is a Count-Min sketch: it tracks approximate frequencies of elements
// in sublinear space. Estimate never under-reports a count — it returns the
// true frequency plus a bounded, one-sided error — which makes it the right
// tool for "how often have I seen this?" over a stream too large to hold exact
// counters for.
//
// A Sketch is not safe for concurrent use. Wrap one with NewConcurrent for a
// goroutine-safe variant.
//
// The zero value is not usable; construct a Sketch with New.
type Sketch[T comparable] struct {
	counts []uint64 // d*w counters, row-major
	w      uint64   // columns per row
	d      uint64   // rows (depth)
	seed   uint64
	hasher func(seed uint64, value T) uint64
	total  uint64 // total count added across all elements
}

// Interface guard.
var _ Frequency[string] = (*Sketch[string])(nil)

// New creates a Count-Min sketch whose estimates overshoot the true count by at
// most epsilon·N (where N is the total count added) with probability at least
// 1-delta. Smaller epsilon and delta cost more memory:
//
//	w = ceil(e/epsilon)   columns
//	d = ceil(ln(1/delta)) rows
//
// New returns an error wrapping ErrInvalidConfig if epsilon or delta is not in
// the open interval (0, 1).
func New[T comparable](epsilon, delta float64, opts ...Option[T]) (*Sketch[T], error) {
	if epsilon <= 0 || epsilon >= 1 || math.IsNaN(epsilon) {
		return nil, fmt.Errorf("%w: epsilon must be in (0,1), got %v", ErrInvalidConfig, epsilon)
	}
	if delta <= 0 || delta >= 1 || math.IsNaN(delta) {
		return nil, fmt.Errorf("%w: delta must be in (0,1), got %v", ErrInvalidConfig, delta)
	}

	// With epsilon, delta in (0,1): e/epsilon > e > 2 so w >= 3, and
	// ln(1/delta) > 0 so d >= 1 — no clamping is needed. Compute the dimensions
	// in floating point first so an absurd configuration is rejected before the
	// (potentially overflowing) cast and allocation.
	wFloat := math.Ceil(math.E / epsilon)
	dFloat := math.Ceil(math.Log(1 / delta))
	if wFloat*dFloat > maxCounters {
		return nil, fmt.Errorf("%w: these bounds need %.0f counters, exceeding the %d-counter limit; raise epsilon or delta",
			ErrInvalidConfig, wFloat*dFloat, int64(maxCounters))
	}
	w := uint64(wFloat)
	d := uint64(dFloat)

	s := &Sketch[T]{
		counts: make([]uint64, w*d),
		w:      w,
		d:      d,
		hasher: sketchhash.Hash64[T],
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

// Option configures a Sketch at construction time.
type Option[T comparable] func(*Sketch[T])

// WithSeed sets the hashing seed. Two sketches must share a seed (and the same
// dimensions) to be merged. The default seed is 0.
func WithSeed[T comparable](seed uint64) Option[T] {
	return func(s *Sketch[T]) { s.seed = seed }
}

// WithHasher overrides the hash function for custom key types or cross-process
// reproducibility. The function must be deterministic in (seed, value).
func WithHasher[T comparable](hasher func(seed uint64, value T) uint64) Option[T] {
	return func(s *Sketch[T]) {
		if hasher != nil {
			s.hasher = hasher
		}
	}
}

// Add records one occurrence of value.
func (s *Sketch[T]) Add(value T) { s.AddCount(value, 1) }

// AddCount records n occurrences of value. Counters saturate at the maximum
// uint64 rather than overflowing.
func (s *Sketch[T]) AddCount(value T, n uint64) {
	h1, h2 := s.pair(value)
	for row := uint64(0); row < s.d; row++ {
		col := (h1 + row*h2) % s.w
		i := row*s.w + col
		if s.counts[i] > math.MaxUint64-n {
			s.counts[i] = math.MaxUint64
		} else {
			s.counts[i] += n
		}
	}
	if s.total > math.MaxUint64-n {
		s.total = math.MaxUint64
	} else {
		s.total += n
	}
}

// Estimate returns the approximate number of times value has been added. The
// result is never less than the true count and overshoots by at most epsilon·N
// with probability at least 1-delta.
func (s *Sketch[T]) Estimate(value T) uint64 {
	h1, h2 := s.pair(value)
	est := uint64(math.MaxUint64)
	for row := uint64(0); row < s.d; row++ {
		col := (h1 + row*h2) % s.w
		if c := s.counts[row*s.w+col]; c < est {
			est = c
		}
	}
	return est
}

// pair derives the two double-hashing lanes for value.
func (s *Sketch[T]) pair(value T) (uint64, uint64) {
	h1 := s.hasher(s.seed, value)
	h2 := s.hasher(s.seed^0x9E3779B97F4A7C15, value) | 1
	return h1, h2
}

// Merge adds other's counters into s element-wise, so s then estimates
// frequencies over the combined stream. Both sketches must have identical
// dimensions and seed; otherwise Merge returns an error wrapping
// ErrInvalidConfig and leaves s unchanged.
//
// Both sketches must also use the same hash function. A custom hasher supplied
// via WithHasher cannot be compared for equality, so a hasher mismatch is not
// detected — merging sketches built with different hashers silently produces
// meaningless results.
func (s *Sketch[T]) Merge(other *Sketch[T]) error {
	if other == nil {
		return fmt.Errorf("%w: cannot merge a nil sketch", ErrInvalidConfig)
	}
	if s.w != other.w || s.d != other.d || s.seed != other.seed {
		return fmt.Errorf("%w: sketches differ (w=%d/%d d=%d/%d seed=%d/%d)",
			ErrInvalidConfig, s.w, other.w, s.d, other.d, s.seed, other.seed)
	}
	// Equal dimensions guarantee equal counts length, since the length is w*d
	// and the fields are unexported.
	for i := range s.counts {
		if s.counts[i] > math.MaxUint64-other.counts[i] {
			s.counts[i] = math.MaxUint64
		} else {
			s.counts[i] += other.counts[i]
		}
	}
	if s.total > math.MaxUint64-other.total {
		s.total = math.MaxUint64
	} else {
		s.total += other.total
	}
	return nil
}

// Clear resets every counter to zero, keeping dimensions and seed.
func (s *Sketch[T]) Clear() {
	for i := range s.counts {
		s.counts[i] = 0
	}
	s.total = 0
}

// Total returns the sum of all counts added (N), the scale of the epsilon·N
// error term.
func (s *Sketch[T]) Total() uint64 { return s.total }

// Width returns the number of columns per row (w).
func (s *Sketch[T]) Width() uint64 { return s.w }

// Depth returns the number of rows (d).
func (s *Sketch[T]) Depth() uint64 { return s.d }
