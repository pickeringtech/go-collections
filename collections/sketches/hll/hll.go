package hll

import (
	"errors"
	"fmt"
	"math"
	"math/bits"

	"github.com/pickeringtech/go-collections/collections/sketches/internal/sketchhash"
)

// ErrInvalidConfig is returned by New when the requested precision is out of
// range, and is the root of the error Merge returns when two sketches have
// different precision, so callers can test with errors.Is.
var ErrInvalidConfig = errors.New("hll: invalid configuration")

const (
	// MinPrecision is the smallest allowed precision, giving m = 2^4 = 16
	// registers. Precision p sets m = 2^p and a standard error of about
	// 1.04/sqrt(m).
	MinPrecision = 4
	// MaxPrecision is the largest allowed precision, giving m = 2^18 = 262144
	// registers.
	MaxPrecision = 18
	// DefaultPrecision gives 2^14 = 16384 registers (~16 KB) and a standard
	// error of about 0.81%.
	DefaultPrecision = 14
)

// Sketch is a HyperLogLog: it estimates the number of distinct elements
// (cardinality) of a stream using a fixed, tiny amount of memory — a few
// kilobytes for billions of distinct items — by observing the longest run of
// leading zero bits seen in a hash, bucketed across many registers.
//
// A Sketch is not safe for concurrent use. Wrap one with NewConcurrent for a
// goroutine-safe variant.
//
// The zero value is not usable; construct a Sketch with New.
type Sketch[T comparable] struct {
	registers []uint8 // m = 1<<precision registers holding observed ranks
	precision int     // p: register index is the top p bits of the hash
	seed      uint64
	hasher    func(seed uint64, value T) uint64
}

// Interface guard.
var _ Cardinality[string] = (*Sketch[string])(nil)

// New creates a HyperLogLog with DefaultPrecision, adjustable via WithPrecision.
// It returns an error wrapping ErrInvalidConfig if the resulting precision is
// outside [MinPrecision, MaxPrecision].
func New[T comparable](opts ...Option[T]) (*Sketch[T], error) {
	s := &Sketch[T]{
		precision: DefaultPrecision,
		hasher:    sketchhash.Hash64[T],
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.precision < MinPrecision || s.precision > MaxPrecision {
		return nil, fmt.Errorf("%w: precision must be in [%d,%d], got %d",
			ErrInvalidConfig, MinPrecision, MaxPrecision, s.precision)
	}
	s.registers = make([]uint8, 1<<s.precision)
	return s, nil
}

// Option configures a Sketch at construction time.
type Option[T comparable] func(*Sketch[T])

// WithPrecision sets the register-index width p, trading memory for accuracy:
// m = 2^p registers, standard error about 1.04/sqrt(m). p must be in
// [MinPrecision, MaxPrecision].
func WithPrecision[T comparable](p int) Option[T] {
	return func(s *Sketch[T]) { s.precision = p }
}

// WithSeed sets the hashing seed. Two sketches must share a seed (and the same
// precision) to be merged. The default seed is 0.
func WithSeed[T comparable](seed uint64) Option[T] {
	return func(s *Sketch[T]) { s.seed = seed }
}

// WithHasher overrides the hash function for custom key types or cross-process
// reproducibility. The function must be deterministic in (seed, value) and
// distribute its output bits well, since HyperLogLog reads both the high index
// bits and the low-bit zero-run.
func WithHasher[T comparable](hasher func(seed uint64, value T) uint64) Option[T] {
	return func(s *Sketch[T]) {
		if hasher != nil {
			s.hasher = hasher
		}
	}
}

// Add records an element. Adding the same element again never changes the
// estimate, so Count reflects the number of distinct elements seen.
func (s *Sketch[T]) Add(value T) {
	x := s.hasher(s.seed, value)
	p := s.precision
	idx := x >> (64 - p)
	// Shift the index bits out, then set a guard bit just below the remaining
	// field so the leading-zero count is bounded by (64-p) even when the rest
	// is all zeros.
	rest := (x << p) | (1 << (p - 1))
	// LeadingZeros64 returns a value in [0,64]; +1 keeps the rank within a
	// register byte, so this conversion never overflows.
	rank := uint8(bits.LeadingZeros64(rest) + 1) // #nosec G115 -- rank is in [1,64], always fits uint8
	if rank > s.registers[idx] {
		s.registers[idx] = rank
	}
}

// Count returns the estimated number of distinct elements added.
func (s *Sketch[T]) Count() uint64 {
	m := float64(len(s.registers))
	sum := 0.0
	zeros := 0
	for _, r := range s.registers {
		sum += math.Ldexp(1, -int(r)) // 2^-r
		if r == 0 {
			zeros++
		}
	}
	estimate := alpha(m) * m * m / sum
	// Small-range correction: linear counting is more accurate while many
	// registers are still empty. With 64-bit hashes the large-range correction
	// is unnecessary and is omitted.
	if estimate <= 2.5*m && zeros > 0 {
		estimate = m * math.Log(m/float64(zeros))
	}
	return uint64(estimate + 0.5)
}

// Merge folds other into s by taking the register-wise maximum, so s then
// estimates the cardinality of the union of both streams. Both sketches must
// have identical precision and seed; otherwise Merge returns an error wrapping
// ErrInvalidConfig and leaves s unchanged.
func (s *Sketch[T]) Merge(other *Sketch[T]) error {
	if other == nil {
		return fmt.Errorf("%w: cannot merge a nil sketch", ErrInvalidConfig)
	}
	if s.precision != other.precision || s.seed != other.seed {
		return fmt.Errorf("%w: sketches differ (precision=%d/%d seed=%d/%d)",
			ErrInvalidConfig, s.precision, other.precision, s.seed, other.seed)
	}
	for i, r := range other.registers {
		if r > s.registers[i] {
			s.registers[i] = r
		}
	}
	return nil
}

// Clear resets every register, returning the sketch to an estimated cardinality
// of zero while keeping its precision and seed.
func (s *Sketch[T]) Clear() {
	for i := range s.registers {
		s.registers[i] = 0
	}
}

// Precision returns the register-index width (p).
func (s *Sketch[T]) Precision() int { return s.precision }

// RegisterCount returns the number of registers (m = 2^p).
func (s *Sketch[T]) RegisterCount() int { return len(s.registers) }

// StandardError returns the expected relative error of Count, about
// 1.04/sqrt(m).
func (s *Sketch[T]) StandardError() float64 {
	return 1.04 / math.Sqrt(float64(len(s.registers)))
}

// alpha is the HyperLogLog bias-correction constant for m registers.
func alpha(m float64) float64 {
	switch m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	default:
		return 0.7213 / (1 + 1.079/m)
	}
}
