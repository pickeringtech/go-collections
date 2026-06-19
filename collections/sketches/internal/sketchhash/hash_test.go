package sketchhash

import "testing"

// TestHash64_Deterministic checks that the fast-path key types hash to the same
// value on repeated calls with the same seed — the property Merge and the
// golden tests rely on.
func TestHash64_Deterministic(t *testing.T) {
	tests := []struct {
		name string
		fn   func() (uint64, uint64)
	}{
		{"string", func() (uint64, uint64) { return Hash64(7, "hello"), Hash64(7, "hello") }},
		{"int", func() (uint64, uint64) { return Hash64(7, 42), Hash64(7, 42) }},
		{"int64", func() (uint64, uint64) { return Hash64(7, int64(42)), Hash64(7, int64(42)) }},
		{"uint64", func() (uint64, uint64) { return Hash64(7, uint64(42)), Hash64(7, uint64(42)) }},
		{"float64", func() (uint64, uint64) { return Hash64(7, 3.14), Hash64(7, 3.14) }},
		{"struct", func() (uint64, uint64) {
			type pt struct{ X, Y int }
			return Hash64(7, pt{1, 2}), Hash64(7, pt{1, 2})
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
		{"int8", func(s uint64) uint64 { return Hash64(s, int8(-7)) }},
		{"int16", func(s uint64) uint64 { return Hash64(s, int16(-7)) }},
		{"int32", func(s uint64) uint64 { return Hash64(s, int32(-7)) }},
		{"int64", func(s uint64) uint64 { return Hash64(s, int64(-7)) }},
		{"int", func(s uint64) uint64 { return Hash64(s, -7) }},
		{"uint", func(s uint64) uint64 { return Hash64(s, uint(7)) }},
		{"uint8", func(s uint64) uint64 { return Hash64(s, uint8(7)) }},
		{"uint16", func(s uint64) uint64 { return Hash64(s, uint16(7)) }},
		{"uint32", func(s uint64) uint64 { return Hash64(s, uint32(7)) }},
		{"uint64", func(s uint64) uint64 { return Hash64(s, uint64(7)) }},
		{"uintptr", func(s uint64) uint64 { return Hash64(s, uintptr(7)) }},
		{"float32", func(s uint64) uint64 { return Hash64(s, float32(1.5)) }},
		{"float64", func(s uint64) uint64 { return Hash64(s, 1.5) }},
		{"string", func(s uint64) uint64 { return Hash64(s, "x") }},
		{"comparable-fallback", func(s uint64) uint64 {
			type pt struct{ X, Y int }
			return Hash64(s, pt{1, 2})
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
	if Hash64(1, "x") == Hash64(2, "x") {
		t.Error("different seeds produced the same hash for the same value")
	}
}

// TestHash64_ValueSensitivity checks the obvious: distinct values usually hash
// apart. (A collision here would still be valid, but for these inputs none is
// expected.)
func TestHash64_ValueSensitivity(t *testing.T) {
	if Hash64(1, "a") == Hash64(1, "b") {
		t.Error("distinct values collided unexpectedly")
	}
	if Hash64(1, 1) == Hash64(1, 2) {
		t.Error("distinct ints collided unexpectedly")
	}
}

// TestPair_SecondLaneOdd checks the double-hash invariant: the second lane is
// always odd (hence non-zero), so it can step through every table slot.
func TestPair_SecondLaneOdd(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_, h2 := Pair(uint64(i), i)
		if h2&1 == 0 {
			t.Fatalf("second lane not odd for i=%d: %d", i, h2)
		}
	}
}

// TestPair_LanesDiffer checks the two lanes are independent for the same value.
func TestPair_LanesDiffer(t *testing.T) {
	collisions := 0
	for i := 0; i < 1000; i++ {
		h1, h2 := Pair(0, i)
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
		counts[Hash64(0, i)%buckets]++
	}
	expected := n / buckets
	for b, c := range counts {
		// Allow a generous ±20% band; this only catches gross skew.
		if c < expected*8/10 || c > expected*12/10 {
			t.Errorf("bucket %d count %d outside ±20%% of %d", b, c, expected)
		}
	}
}
