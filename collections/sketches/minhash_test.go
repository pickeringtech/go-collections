package sketches_test

import (
	"math"
	randv2 "math/rand/v2"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches"
)

func deterministicRng() *randv2.Rand {
	return randv2.New(randv2.NewPCG(42, 0))
}

func TestNewMinHash_nilRngIsDeterministic(t *testing.T) {
	// Two sketches constructed with nil rng and the same numHashes must be
	// compatible (same permutation family) and produce identical estimates.
	a := sketches.NewMinHash[string](64, nil)
	b := sketches.NewMinHash[string](64, nil)

	a.Add("hello")
	b.Add("hello")

	est, ok := sketches.EstimatedJaccard(a, b)
	if !ok {
		t.Fatalf("EstimatedJaccard returned ok=false for identically configured sketches")
	}
	if math.Abs(est-1.0) > 1e-9 {
		t.Fatalf("identical sets: EstimatedJaccard = %v, want 1.0", est)
	}
}

func TestEstimatedJaccard_identicalSets(t *testing.T) {
	rng := deterministicRng()
	a := sketches.NewMinHash[string](128, rng)
	rng2 := deterministicRng()
	b := sketches.NewMinHash[string](128, rng2)

	words := []string{"the", "quick", "brown", "fox", "jumps"}
	for _, w := range words {
		a.Add(w)
		b.Add(w)
	}

	est, ok := sketches.EstimatedJaccard(a, b)
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	if math.Abs(est-1.0) > 1e-9 {
		t.Fatalf("identical sets: got %v, want 1.0", est)
	}
}

func TestEstimatedJaccard_disjointSets(t *testing.T) {
	rng := deterministicRng()
	a := sketches.NewMinHash[string](256, rng)
	rng2 := deterministicRng()
	b := sketches.NewMinHash[string](256, rng2)

	for _, w := range []string{"a", "b", "c", "d", "e"} {
		a.Add(w)
	}
	for _, w := range []string{"f", "g", "h", "i", "j"} {
		b.Add(w)
	}

	est, ok := sketches.EstimatedJaccard(a, b)
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	// Disjoint sets have Jaccard = 0; allow generous tolerance for randomness.
	if est > 0.1 {
		t.Fatalf("disjoint sets: got %v, want ~0.0", est)
	}
}

func TestEstimatedJaccard_partialOverlap(t *testing.T) {
	rng := deterministicRng()
	a := sketches.NewMinHash[int](512, rng)
	rng2 := deterministicRng()
	b := sketches.NewMinHash[int](512, rng2)

	// a = {0..99}, b = {50..149}: exact Jaccard = 50/150 ≈ 0.333
	for i := 0; i < 100; i++ {
		a.Add(i)
	}
	for i := 50; i < 150; i++ {
		b.Add(i)
	}

	est, ok := sketches.EstimatedJaccard(a, b)
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	// MinHash with 512 hashes should be within ~5% of the true value.
	exact := 50.0 / 150.0
	if math.Abs(est-exact) > 0.07 {
		t.Fatalf("partial overlap: got %v, want ~%v (±0.07)", est, exact)
	}
}

func TestEstimatedJaccard_mismatchedNumHashes(t *testing.T) {
	a := sketches.NewMinHash[string](64, nil)
	b := sketches.NewMinHash[string](128, nil)

	_, ok := sketches.EstimatedJaccard(a, b)
	if ok {
		t.Fatalf("expected ok=false for mismatched numHashes, got true")
	}
}

func TestEstimatedJaccard_mismatchedPermFamily(t *testing.T) {
	// Two sketches with different rng produce different permutation families.
	rng1 := randv2.New(randv2.NewPCG(1, 0))
	rng2 := randv2.New(randv2.NewPCG(999, 0))

	a := sketches.NewMinHash[string](64, rng1)
	b := sketches.NewMinHash[string](64, rng2)

	// There is a tiny probabilistic chance that all coefficients agree — but
	// with 64 pairs drawn from different seeds it is astronomically unlikely.
	_, ok := sketches.EstimatedJaccard(a, b)
	if ok {
		t.Fatalf("expected ok=false for mismatched perm family, got true")
	}
}

func TestMinHash_emptySignature(t *testing.T) {
	m := sketches.NewMinHash[string](16, nil)
	sig := m.Signature()
	if len(sig) != 16 {
		t.Fatalf("Signature() len = %v, want 16", len(sig))
	}
	// All positions should be MaxUint64 before any element is added.
	for i, v := range sig {
		if v != math.MaxUint64 {
			t.Fatalf("Signature()[%d] = %v, want MaxUint64", i, v)
		}
	}
}

func TestMinHash_signatureIsCopy(t *testing.T) {
	m := sketches.NewMinHash[string](4, nil)
	m.Add("hello")

	sig1 := m.Signature()
	// Mutate the returned slice — the sketch must not be affected.
	for i := range sig1 {
		sig1[i] = 0
	}

	m.Add("world")
	sig2 := m.Signature()

	// sig2 must differ from sig1 (it was not clobbered by the mutation).
	changed := false
	for i := range sig2 {
		if sig2[i] != sig1[i] {
			changed = true
			break
		}
	}
	if !changed {
		t.Fatalf("Signature() appears to return the internal slice; expected a copy")
	}
}

func TestEstimatedJaccardInRange(t *testing.T) {
	a := sketches.NewMinHash[int](64, nil)
	b := sketches.NewMinHash[int](64, nil)

	for i := 0; i < 50; i++ {
		a.Add(i)
	}
	for i := 25; i < 75; i++ {
		b.Add(i)
	}

	est, ok := sketches.EstimatedJaccard(a, b)
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	if est < 0 || est > 1 {
		t.Fatalf("EstimatedJaccard = %v, want in [0,1]", est)
	}
}

func TestNewMinHash_zeroNumHashesClampedToOne(t *testing.T) {
	// numHashes == 0 must be clamped to 1; the sketch must be usable and
	// EstimatedJaccard must not return NaN or panic.
	a := sketches.NewMinHash[string](0, nil)
	b := sketches.NewMinHash[string](0, nil)

	a.Add("hello")
	b.Add("hello")

	est, ok := sketches.EstimatedJaccard(a, b)
	if !ok {
		t.Fatalf("EstimatedJaccard ok = false, want true for clamped sketches")
	}
	if math.IsNaN(est) {
		t.Fatalf("EstimatedJaccard = NaN, want a valid float in [0,1]")
	}
	if est < 0 || est > 1 {
		t.Fatalf("EstimatedJaccard = %v, want in [0,1]", est)
	}
}

func TestNewMinHash_negativeNumHashesClampedToOne(t *testing.T) {
	// numHashes < 0 must not panic; it is clamped to 1.
	a := sketches.NewMinHash[string](-10, nil)
	b := sketches.NewMinHash[string](-10, nil)

	a.Add("world")
	b.Add("world")

	est, ok := sketches.EstimatedJaccard(a, b)
	if !ok {
		t.Fatalf("EstimatedJaccard ok = false, want true for clamped sketches")
	}
	if math.IsNaN(est) {
		t.Fatalf("EstimatedJaccard = NaN, want a valid float in [0,1]")
	}
	if est < 0 || est > 1 {
		t.Fatalf("EstimatedJaccard = %v, want in [0,1]", est)
	}
}

func TestEstimatedJaccard_nilSketches(t *testing.T) {
	valid := sketches.NewMinHash[string](16, nil)

	cases := []struct {
		name string
		a, b *sketches.MinHash[string]
	}{
		{"both nil", nil, nil},
		{"first nil", nil, valid},
		{"second nil", valid, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			est, ok := sketches.EstimatedJaccard(tc.a, tc.b)
			if ok {
				t.Fatalf("EstimatedJaccard ok = true, want false for nil sketch")
			}
			if est != 0 {
				t.Fatalf("EstimatedJaccard = %v, want 0 for nil sketch", est)
			}
		})
	}
}

func TestEstimatedJaccard_zeroValueSketch(t *testing.T) {
	// A zero-value sketch has no permutations; it must not yield (NaN, true).
	valid := sketches.NewMinHash[string](16, nil)

	cases := []struct {
		name string
		a, b *sketches.MinHash[string]
	}{
		{"both zero-value", &sketches.MinHash[string]{}, &sketches.MinHash[string]{}},
		{"first zero-value", &sketches.MinHash[string]{}, valid},
		{"second zero-value", valid, &sketches.MinHash[string]{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			est, ok := sketches.EstimatedJaccard(tc.a, tc.b)
			if ok {
				t.Fatalf("EstimatedJaccard ok = true, want false for zero-value sketch")
			}
			if est != 0 {
				t.Fatalf("EstimatedJaccard = %v, want 0 for zero-value sketch", est)
			}
		})
	}
}
