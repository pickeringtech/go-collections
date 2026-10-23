package sketchhash_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/internal/sketchhash"
)

// TestHash64_Deterministic checks that the fast-path key types hash to the same
// value on repeated calls with the same seed — the property Merge and the
// golden tests rely on.
func TestHash64_Deterministic(t *testing.T) {
	tests := []struct {
		name string
		fn   func() (uint64, uint64)
	}{
		{"string", func() (uint64, uint64) { return sketchhash.Hash64(7, "hello"), sketchhash.Hash64(7, "hello") }},
		{"int", func() (uint64, uint64) { return sketchhash.Hash64(7, 42), sketchhash.Hash64(7, 42) }},
		{"int64", func() (uint64, uint64) { return sketchhash.Hash64(7, int64(42)), sketchhash.Hash64(7, int64(42)) }},
		{"uint64", func() (uint64, uint64) { return sketchhash.Hash64(7, uint64(42)), sketchhash.Hash64(7, uint64(42)) }},
		{"float64", func() (uint64, uint64) { return sketchhash.Hash64(7, 3.14), sketchhash.Hash64(7, 3.14) }},
		{"struct", func() (uint64, uint64) {
			type pt struct{ X, Y int }
			return sketchhash.Hash64(7, pt{1, 2}), sketchhash.Hash64(7, pt{1, 2})
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, b := tc.fn()
			if a != b {
				t.Errorf("Hash64 not deterministic: %d != %d", a, b)
			}
		})
	}
}

// TestHash64_AllKindFastPaths exercises every fast-path branch of the type
// switch so each numeric kind is hashed through its dedicated case, and checks
// the result is stable on repeat.
func TestHash64_AllKindFastPaths(t *testing.T) {
	tests := []struct {
		name string
		fn   func(uint64) uint64
	}{
		{"int8", func(s uint64) uint64 { return sketchhash.Hash64(s, int8(-7)) }},
		{"int16", func(s uint64) uint64 { return sketchhash.Hash64(s, int16(-7)) }},
		{"int32", func(s uint64) uint64 { return sketchhash.Hash64(s, int32(-7)) }},
		{"int64", func(s uint64) uint64 { return sketchhash.Hash64(s, int64(-7)) }},
		{"int", func(s uint64) uint64 { return sketchhash.Hash64(s, -7) }},
		{"uint", func(s uint64) uint64 { return sketchhash.Hash64(s, uint(7)) }},
		{"uint8", func(s uint64) uint64 { return sketchhash.Hash64(s, uint8(7)) }},
		{"uint16", func(s uint64) uint64 { return sketchhash.Hash64(s, uint16(7)) }},
		{"uint32", func(s uint64) uint64 { return sketchhash.Hash64(s, uint32(7)) }},
		{"uint64", func(s uint64) uint64 { return sketchhash.Hash64(s, uint64(7)) }},
		{"uintptr", func(s uint64) uint64 { return sketchhash.Hash64(s, uintptr(7)) }},
		{"float32", func(s uint64) uint64 { return sketchhash.Hash64(s, float32(1.5)) }},
		{"float64", func(s uint64) uint64 { return sketchhash.Hash64(s, 1.5) }},
		{"string", func(s uint64) uint64 { return sketchhash.Hash64(s, "x") }},
		{"comparable-fallback", func(s uint64) uint64 {
			type pt struct{ X, Y int }
			return sketchhash.Hash64(s, pt{1, 2})
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			first := tc.fn(3)
			second := tc.fn(3)
			if first != second {
				t.Errorf("%s: Hash64 not deterministic (%d != %d)", tc.name, first, second)
			}
		})
	}
}

// TestHash64_SeedSensitivity checks that changing the seed changes the hash, so
// distinct seeds give distinct sketches.
func TestHash64_SeedSensitivity(t *testing.T) {
	if sketchhash.Hash64(1, "x") == sketchhash.Hash64(2, "x") {
		t.Error("different seeds produced the same hash for the same value")
	}
}

// TestHash64_ValueSensitivity checks the obvious: distinct values usually hash
// apart. (A collision here would still be valid, but for these inputs none is
// expected.)
func TestHash64_ValueSensitivity(t *testing.T) {
	if sketchhash.Hash64(1, "a") == sketchhash.Hash64(1, "b") {
		t.Error("distinct values collided unexpectedly")
	}
	if sketchhash.Hash64(1, 1) == sketchhash.Hash64(1, 2) {
		t.Error("distinct ints collided unexpectedly")
	}
}

// TestPair_SecondLaneOdd checks the double-hash invariant: the second lane is
// always odd (hence non-zero), so it can step through every table slot.
func TestPair_SecondLaneOdd(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_, h2 := sketchhash.Pair(uint64(i), i)
		if h2&1 == 0 {
			t.Fatalf("second lane not odd for i=%d: %d", i, h2)
		}
	}
}

// TestPair_LanesDiffer checks the two lanes are independent for the same value.
func TestPair_LanesDiffer(t *testing.T) {
	collisions := 0
	for i := 0; i < 1000; i++ {
		h1, h2 := sketchhash.Pair(0, i)
		if h1 == h2 {
			collisions++
		}
	}
	if collisions > 1 {
		t.Errorf("lanes collided %d/1000 times; expected ~0", collisions)
	}
}

// TestHash64_Distribution is a coarse sanity check that the low bits used for
// table indexing are spread across buckets rather than clumped.
func TestHash64_Distribution(t *testing.T) {
	const buckets = 16
	const n = 100000
	counts := make([]int, buckets)
	for i := 0; i < n; i++ {
		counts[sketchhash.Hash64(0, i)%buckets]++
	}
	expected := n / buckets
	for b, c := range counts {
		// Allow a generous ±20% band; this only catches gross skew.
		if c < expected*8/10 || c > expected*12/10 {
			t.Errorf("bucket %d count %d outside ±20%% of %d", b, c, expected)
		}
	}
}

// TestHash64_EqualFloatsHashEqual checks that keys equal under == hash equal:
// +0 and -0 compare equal (and are the same Go map key), so a sketch must not
// tell them apart — otherwise a Bloom filter given +0 reports -0 absent (a
// false negative). Every NaN is canonicalised to one bit pattern, so all NaNs
// hash alike whatever their payload.
func TestHash64_EqualFloatsHashEqual(t *testing.T) {
	negZero := math.Copysign(0, -1)
	quietNaN := math.NaN()
	payloadNaN := math.Float64frombits(0x7FF0_0000_0000_0001)
	negNaN := math.Copysign(math.NaN(), -1)
	tests := []struct {
		name string
		a, b uint64
	}{
		{"float64 +0 and -0", sketchhash.Hash64(3, 0.0), sketchhash.Hash64(3, negZero)},
		{"float32 +0 and -0", sketchhash.Hash64(3, float32(0)), sketchhash.Hash64(3, float32(negZero))},
		{"float64 NaN payloads", sketchhash.Hash64(3, quietNaN), sketchhash.Hash64(3, payloadNaN)},
		{"float64 negative NaN", sketchhash.Hash64(3, quietNaN), sketchhash.Hash64(3, negNaN)},
		{"float32 NaN payloads", sketchhash.Hash64(3, float32(quietNaN)), sketchhash.Hash64(3, math.Float32frombits(0x7F80_0001))},
		{"any-typed +0 and -0", sketchhash.Hash64[any](3, 0.0), sketchhash.Hash64[any](3, negZero)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.a != tt.b {
				t.Errorf("Hash64 = %#x and %#x, want equal", tt.a, tt.b)
			}
		})
	}
}

// TestHash64_SeedsAreIndependent checks that distinct seeds give unrelated
// hash functions. Folding the seed into the bits the first input byte is XORed
// into made Hash64(s, x) == Hash64(0, x^s) for every s < 256, so seeds 0..255
// produced the identical hash multiset over any 256-aligned run of integers
// (and "independent" HyperLogLog trials were one trial repeated). Across 256
// seeds and 1024 keys a random function collides with seed 0's hash set with
// probability ~2^-34, so any shared value is a structural relation.
func TestHash64_SeedsAreIndependent(t *testing.T) {
	const keys = 1024
	reference := make(map[uint64]bool, keys)
	for i := 0; i < keys; i++ {
		reference[sketchhash.Hash64(0, i)] = true
	}
	shared := 0
	for seed := uint64(1); seed < 256; seed++ {
		for i := 0; i < keys; i++ {
			if reference[sketchhash.Hash64(seed, i)] {
				shared++
			}
		}
	}
	if shared != 0 {
		t.Errorf("%d hashes under seeds 1..255 also occur under seed 0, want 0", shared)
	}
}

// TestHash64_SeedNotXoredIntoInput pins the specific relations the old seed
// folding produced, for each key family and for both double-hash lanes.
func TestHash64_SeedNotXoredIntoInput(t *testing.T) {
	h1a, h2a := sketchhash.Pair(1, 6)
	h1b, h2b := sketchhash.Pair(0, 6^1)
	tests := []struct {
		name string
		a, b uint64
	}{
		{"int: Hash64(s, x) vs Hash64(0, x^s)", sketchhash.Hash64(5, 40), sketchhash.Hash64(0, 40^5)},
		{"string: Hash64(1, \"a\") vs Hash64(0, \"`\")", sketchhash.Hash64(1, "a"), sketchhash.Hash64(0, "`")},
		{"float64 low seed bits", sketchhash.Hash64(1, math.Float64frombits(2)), sketchhash.Hash64(0, math.Float64frombits(3))},
		{"Pair first lane", h1a, h1b},
		{"Pair second lane", h2a, h2b},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.a == tt.b {
				t.Errorf("hashes collide (%#x): the seed is not mixed independently of the input", tt.a)
			}
		})
	}
}
