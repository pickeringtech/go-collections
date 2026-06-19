// Package sketches is documented in doc.go.
package sketches

import (
	"hash/maphash"
	"math"
	randv2 "math/rand/v2"
)

// defaultSeed is the PCG seed used when NewMinHash receives a nil rng.
// Keeping it fixed ensures that two sketches constructed with nil in the same
// program run share the same permutation coefficients and are compatible via
// EstimatedJaccard.
const defaultSeed uint64 = 0x9e3779b97f4a7c15

// globalSeed is the maphash seed shared by all MinHash sketches. It is
// initialised once per program run via maphash.MakeSeed. Because all sketches
// use this seed, the hash function h(x) is identical across sketches created
// in the same run, and per-permutation salting is handled entirely by the
// (a, b) universal-hash coefficients. This removes any need for unsafe.
var globalSeed = maphash.MakeSeed()

// permCoeffs holds the (a, b) coefficients for one universal-hash permutation.
// The hash for element x is: (a*h(x) + b) mod 2^64, where h(x) is computed
// via maphash.Comparable using globalSeed. Taking the minimum of these values
// across all added elements gives one MinHash component.
type permCoeffs struct {
	a, b uint64
}

// MinHash is a MinHash sketch for any comparable type T. It maintains a
// fixed-size signature of uint64 minimums — one per hash function in the
// permutation family — that can be compared with EstimatedJaccard to
// approximate the Jaccard similarity between two sets.
//
// MinHash is stateful and NOT goroutine-safe. A thread-safe ConcurrentMinHash
// variant is planned for a later issue.
type MinHash[T comparable] struct {
	// perms holds numHashes permutation coefficients. All are derived from the
	// same rng, so two sketches built from the same rng (or both nil) produce
	// equal perms and are compatible via EstimatedJaccard.
	perms []permCoeffs
	mins  []uint64
}

// NewMinHash constructs a new MinHash[T] sketch with numHashes hash functions.
// A larger numHashes improves accuracy but uses more memory.
//
// numHashes must be at least 1. If a value less than 1 is provided it is
// clamped to 1 so that the sketch is always valid and EstimatedJaccard never
// produces NaN.
//
// rng provides the random source used to derive the (a, b) permutation
// coefficients. Passing nil selects a hard-coded deterministic default seed,
// making sketches with the same numHashes constructed in the same program run
// always produce the same permutation family — and therefore always comparable
// via EstimatedJaccard.
//
// The returned sketch is empty; call Add to insert elements.
func NewMinHash[T comparable](numHashes int, rng *randv2.Rand) *MinHash[T] {
	if numHashes < 1 {
		numHashes = 1
	}
	if rng == nil {
		rng = randv2.New(randv2.NewPCG(defaultSeed, 0))
	}
	perms := make([]permCoeffs, numHashes)
	for i := range perms {
		a := rng.Uint64()
		// Ensure a is odd for full-period behaviour in the linear congruential
		// construction; b can be any value.
		if a%2 == 0 {
			a++
		}
		perms[i] = permCoeffs{a: a, b: rng.Uint64()}
	}
	mins := make([]uint64, numHashes)
	for i := range mins {
		mins[i] = math.MaxUint64
	}
	return &MinHash[T]{
		perms: perms,
		mins:  mins,
	}
}

// Add incorporates element into the sketch, updating the running minimum for
// each permutation. Add never mutates element; it reads it only to compute
// its hash.
func (m *MinHash[T]) Add(element T) {
	h := maphash.Comparable(globalSeed, element)
	for i, p := range m.perms {
		v := p.a*h + p.b
		if v < m.mins[i] {
			m.mins[i] = v
		}
	}
}

// Signature returns a copy of the current MinHash signature — the vector of
// per-permutation minimum hash values accumulated so far. The returned slice
// is a fresh copy; callers may read it without synchronisation.
func (m *MinHash[T]) Signature() []uint64 {
	out := make([]uint64, len(m.mins))
	copy(out, m.mins)
	return out
}

// EstimatedJaccard estimates the Jaccard similarity between the element sets
// represented by two MinHash sketches by computing the fraction of signature
// positions at which their minimum hash values agree.
//
// It returns ok == false when the two sketches are incompatible: they must have
// the same number of hash functions (numHashes) and the same permutation
// coefficients. Sketches created with different rng instances will typically
// have mismatched parameters and are not comparable.
//
// The estimate is in [0, 1].
func EstimatedJaccard[T comparable](a, b *MinHash[T]) (float64, bool) {
	if len(a.perms) != len(b.perms) {
		return 0, false
	}
	for i := range a.perms {
		if a.perms[i] != b.perms[i] {
			return 0, false
		}
	}
	matches := 0
	for i := range a.mins {
		if a.mins[i] == b.mins[i] {
			matches++
		}
	}
	return float64(matches) / float64(len(a.mins)), true
}
