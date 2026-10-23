// Package sketchhash provides the seeded 64-bit hashing shared by the
// probabilistic sketches in collections/sketches (Bloom, Count-Min,
// HyperLogLog).
//
// Two sketches can only be merged when they hash identically, so the hashing
// has to be reproducible from a stored seed rather than randomised per process.
// Hash64 therefore takes a numeric seed and, for the key types sketches are
// used with in practice — strings, the integer kinds and the floats — produces
// a value that depends only on (seed, value). Those are the keys tests assert
// on, so golden values are stable across runs. The seed is mixed into the whole
// hash state before any input is consumed, so distinct seeds behave as
// independent hash functions (as repeated HyperLogLog trials or a Count-Min
// sketch's rows require), not as relabellings of one another.
//
// For any other comparable key type Hash64 falls back to hash/maphash's
// Comparable (Go 1.24+), keyed by a process-stable seed created once at init.
// That keeps every sketch in a single process mutually mergeable, but the
// values are not reproducible across processes for those exotic key types —
// pass a custom hasher to a sketch if you need that.
package sketchhash

import (
	"hash/maphash"
	"math"
)

const (
	fnvPrime64 = 1099511628211

	// lane is the golden-ratio odd constant used to derive an independent
	// second hash lane for Kirsch–Mitzenmacher double hashing.
	lane = 0x9E3779B97F4A7C15
)

// comparableSeed keys the maphash.Comparable fallback. It is created once so
// every sketch in this process hashes exotic key types the same way (a
// prerequisite for Merge); it is intentionally not derived from the caller's
// numeric seed, whose basis (see seedBasis) is folded in separately by Hash64.
var comparableSeed = maphash.MakeSeed()

// Hash64 returns a 64-bit hash of v keyed by seed. For strings, the integer
// kinds and the floats the result depends only on (seed, v) and is stable
// across runs and processes; for other comparable types it is stable only
// within a single process (see the package doc).
//
// Floats are hashed by value, not raw bits, so keys equal under == hash
// equal: -0 hashes as +0. NaN is never equal to anything, so it is given one
// fixed meaning instead: every NaN, whatever its sign or payload, hashes
// alike, and a sketch therefore treats all NaNs as a single element.
func Hash64[T comparable](seed uint64, v T) uint64 {
	basis := seedBasis(seed)
	switch x := any(v).(type) {
	case string:
		return mix64(fnvString(basis, x))
	case float32:
		return mix64(fnvUint(basis, float32Bits(x)))
	case float64:
		return mix64(fnvUint(basis, float64Bits(x)))
	default:
		u, ok := asUint64(v)
		if ok {
			return mix64(fnvUint(basis, u))
		}
		return mix64(maphash.Comparable(comparableSeed, v) ^ basis)
	}
}

// seedBasis derives the hashing start state from seed by running it through
// the SplitMix64 finalizer, a bijection that spreads every seed bit across the
// whole word. Starting FNV-1a from the raw seed instead (offset ^ seed) put the
// seed in exactly the bits the first input byte is XORed into, so
// Hash64(s, x) == Hash64(0, x^s) for every s < 256 and those seeds were the
// same hash function on any 256-aligned run of integer keys. Distinct seeds now
// start from unrelated states, so they give independent hash functions.
func seedBasis(seed uint64) uint64 {
	return mix64(seed)
}

// asUint64 returns the bit pattern of v as a uint64 for the integer kinds,
// reporting ok=false for any other type so Hash64 can fall back. It is split out
// from Hash64 to keep each function's branching modest. Signed kinds are widened
// to int64 (a safe conversion) and reinterpreted by signedBits; unsigned kinds
// widen directly.
func asUint64[T comparable](v T) (uint64, bool) {
	switch x := any(v).(type) {
	case int:
		return signedBits(int64(x)), true
	case int8:
		return signedBits(int64(x)), true
	case int16:
		return signedBits(int64(x)), true
	case int32:
		return signedBits(int64(x)), true
	case int64:
		return signedBits(x), true
	case uint:
		return uint64(x), true
	case uint8:
		return uint64(x), true
	case uint16:
		return uint64(x), true
	case uint32:
		return uint64(x), true
	case uint64:
		return x, true
	case uintptr:
		return uint64(x), true
	default:
		return 0, false
	}
}

// Canonical NaN bit patterns: every NaN is hashed as one of these, whatever
// its sign or payload.
const (
	canonicalNaN32 = 0x7FC0_0000
	canonicalNaN64 = 0x7FF8_0000_0000_0001
)

// float64Bits returns the bit pattern Hash64 hashes for f. +0 and -0 compare
// equal (and are the same Go map key), so -0 is mapped to +0's all-zero
// pattern; otherwise a Bloom filter given one would report the other absent.
// NaN never compares equal to anything, so "equal keys hash equal" says
// nothing about it; every NaN is mapped to one canonical pattern so that a
// sketch treats all NaNs as a single element instead of splitting them by
// payload.
func float64Bits(f float64) uint64 {
	if f == 0 {
		return 0
	}
	if math.IsNaN(f) {
		return canonicalNaN64
	}
	return math.Float64bits(f)
}

// float32Bits is float64Bits for float32: -0 hashes as +0 and every NaN as one
// canonical pattern.
func float32Bits(f float32) uint64 {
	if f == 0 {
		return 0
	}
	if math.IsNaN(float64(f)) {
		return canonicalNaN32
	}
	return uint64(math.Float32bits(f))
}

// signedBits reinterprets a signed 64-bit integer as its uint64 bit pattern.
// The two's-complement reinterpretation is exactly what hashing wants, so the
// sign change carries no information loss.
func signedBits(x int64) uint64 {
	return uint64(x) // #nosec G115 -- intentional bit reinterpretation for hashing; sign/overflow is harmless here
}

// Pair returns two independent 64-bit hashes of v for double hashing. The
// second lane is forced odd (and therefore non-zero) so that, combined modulo a
// table size, it can reach every slot.
func Pair[T comparable](seed uint64, v T) (h1, h2 uint64) {
	h1 = Hash64(seed, v)
	h2 = Hash64(seed^lane, v) | 1
	return h1, h2
}

// fnvString hashes s with FNV-1a, starting from basis (see seedBasis) in
// place of the standard offset basis.
func fnvString(basis uint64, s string) uint64 {
	h := basis
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= fnvPrime64
	}
	return h
}

// fnvUint hashes the eight bytes of x with FNV-1a, starting from basis (see
// seedBasis) in place of the standard offset basis.
func fnvUint(basis uint64, x uint64) uint64 {
	h := basis
	for i := 0; i < 8; i++ {
		h ^= x & 0xff
		h *= fnvPrime64
		x >>= 8
	}
	return h
}

// mix64 is the SplitMix64 finalizer. FNV-1a mixes its low bits well but leaves
// the high bits weak; the finalizer avalanches the whole word so the value is
// safe to slice for register indices (HyperLogLog) and double-hash lanes.
func mix64(z uint64) uint64 {
	z += 0x9E3779B97F4A7C15
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}
