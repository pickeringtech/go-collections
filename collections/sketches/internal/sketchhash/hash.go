// Package sketchhash provides the seeded 64-bit hashing shared by the
// probabilistic sketches in collections/sketches (Bloom, Count-Min,
// HyperLogLog).
//
// Two sketches can only be merged when they hash identically, so the hashing
// has to be reproducible from a stored seed rather than randomised per process.
// Hash64 therefore takes a numeric seed and, for the key types sketches are
// used with in practice — strings, the integer kinds and the floats — produces
// a value that depends only on (seed, value). Those are the keys tests assert
// on, so golden values are stable across runs.
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
	fnvOffset64 = 14695981039346656037
	fnvPrime64  = 1099511628211

	// lane is the golden-ratio odd constant used to derive an independent
	// second hash lane for Kirsch–Mitzenmacher double hashing.
	lane = 0x9E3779B97F4A7C15
)

// comparableSeed keys the maphash.Comparable fallback. It is created once so
// every sketch in this process hashes exotic key types the same way (a
// prerequisite for Merge); it is intentionally not derived from the caller's
// numeric seed, which is folded in separately by Hash64.
var comparableSeed = maphash.MakeSeed()

// Hash64 returns a 64-bit hash of v keyed by seed. For strings, the integer
// kinds and the floats the result depends only on (seed, v) and is stable
// across runs and processes; for other comparable types it is stable only
// within a single process (see the package doc).
func Hash64[T comparable](seed uint64, v T) uint64 {
	switch x := any(v).(type) {
	case string:
		return mix64(fnvString(seed, x))
	case float32:
		return mix64(fnvUint(seed, uint64(math.Float32bits(x))))
	case float64:
		return mix64(fnvUint(seed, math.Float64bits(x)))
	default:
		if u, ok := asUint64(v); ok {
			return mix64(fnvUint(seed, u))
		}
		return mix64(maphash.Comparable(comparableSeed, v) ^ seed)
	}
}

// asUint64 returns the bit pattern of v as a uint64 for the integer kinds,
// reporting ok=false for any other type so Hash64 can fall back. It is split out
// from Hash64 to keep each function's branching modest.
func asUint64[T comparable](v T) (uint64, bool) {
	switch x := any(v).(type) {
	case int:
		return uint64(x), true
	case int8:
		return uint64(x), true
	case int16:
		return uint64(x), true
	case int32:
		return uint64(x), true
	case int64:
		return uint64(x), true
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

// Pair returns two independent 64-bit hashes of v for double hashing. The
// second lane is forced odd (and therefore non-zero) so that, combined modulo a
// table size, it can reach every slot.
func Pair[T comparable](seed uint64, v T) (h1, h2 uint64) {
	h1 = Hash64(seed, v)
	h2 = Hash64(seed^lane, v) | 1
	return h1, h2
}

// fnvString folds seed into the FNV-1a basis and hashes s.
func fnvString(seed uint64, s string) uint64 {
	h := uint64(fnvOffset64) ^ seed
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= fnvPrime64
	}
	return h
}

// fnvUint folds seed into the FNV-1a basis and hashes the eight bytes of x.
func fnvUint(seed uint64, x uint64) uint64 {
	h := uint64(fnvOffset64) ^ seed
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
