package bloom_test

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/bloom"
)

func mustNew[T comparable](t *testing.T, n int, p float64, opts ...bloom.Option[T]) *bloom.Filter[T] {
	t.Helper()
	f, err := bloom.New(n, p, opts...)
	if err != nil {
		t.Fatalf("New(%d, %v): unexpected error %v", n, p, err)
	}
	return f
}

func TestNew_InvalidConfig(t *testing.T) {
	tests := []struct {
		name string
		n    int
		p    float64
		// wantMsg is a fragment the error must carry, so the caller is told
		// which argument is wrong rather than only that something is.
		wantMsg string
	}{
		{"zero items", 0, 0.01, "expectedItems must be positive"},
		{"negative items", -5, 0.01, "expectedItems must be positive"},
		// A zero rate would also be caught by the size cap (it needs infinitely
		// many bits), but it must be reported as an out-of-range rate.
		{"zero rate", 100, 0, "falsePositiveRate must be in (0,1)"},
		{"one rate", 100, 1, "falsePositiveRate must be in (0,1)"},
		{"negative rate", 100, -0.1, "falsePositiveRate must be in (0,1)"},
		{"rate above one", 100, 1.5, "falsePositiveRate must be in (0,1)"},
		// Absurd capacity: the optimal bit count exceeds the allocation cap, so
		// New rejects it instead of overflowing or panicking on allocation.
		// (Inputs stay within a 32-bit int so the test builds on all GOARCHes.)
		{"astronomically large", 1_000_000_000, 1e-300, "bit limit"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := bloom.New[string](tc.n, tc.p)
			if !errors.Is(err, bloom.ErrInvalidConfig) {
				t.Fatalf("New(%d, %v) error = %v, want ErrInvalidConfig", tc.n, tc.p, err)
			}
			if !strings.Contains(err.Error(), tc.wantMsg) {
				t.Errorf("New(%d, %v) error = %q, want it to contain %q", tc.n, tc.p, err, tc.wantMsg)
			}
		})
	}
}

func TestFilter_NoFalseNegatives(t *testing.T) {
	f := mustNew[int](t, 1000, 0.01)
	for i := 0; i < 1000; i++ {
		f.Add(i)
	}
	for i := 0; i < 1000; i++ {
		if !f.Contains(i) {
			t.Fatalf("Contains(%d) = false after Add; Bloom filters must never report a false negative", i)
		}
	}
}

func TestFilter_EmptyContainsNothing(t *testing.T) {
	f := mustNew[string](t, 100, 0.01)
	if f.Contains("absent") {
		t.Error("empty filter reported membership")
	}
}

func TestFilter_FalsePositiveRateNearTarget(t *testing.T) {
	const n = 10000
	const target = 0.01
	f := mustNew[int](t, n, target)
	for i := 0; i < n; i++ {
		f.Add(i)
	}
	// Query n elements that were never added (disjoint range) and measure the
	// observed false-positive rate.
	fp := 0
	for i := n; i < 2*n; i++ {
		if f.Contains(i) {
			fp++
		}
	}
	rate := float64(fp) / float64(n)
	// Allow up to ~3x the target to absorb sampling noise and the double-hash
	// approximation; a gross blow-up still fails.
	if rate > target*3 {
		t.Errorf("observed false-positive rate %.4f exceeds 3x target %.4f", rate, target)
	}
}

func TestFilter_Sizing(t *testing.T) {
	// For n=1000, p=0.01 the optimal sizing is m≈9585 bits, k≈7.
	f := mustNew[int](t, 1000, 0.01)
	if got := f.HashCount(); got != 7 {
		t.Errorf("HashCount() = %d, want 7", got)
	}
	if got := f.BitSize(); got < 9000 || got > 10000 {
		t.Errorf("BitSize() = %d, want ~9585", got)
	}
}

func TestFilter_Merge(t *testing.T) {
	a := mustNew[int](t, 1000, 0.01)
	b := mustNew[int](t, 1000, 0.01)
	for i := 0; i < 500; i++ {
		a.Add(i)
	}
	for i := 500; i < 1000; i++ {
		b.Add(i)
	}
	if err := a.Merge(b); err != nil {
		t.Fatalf("Merge: unexpected error %v", err)
	}
	for i := 0; i < 1000; i++ {
		if !a.Contains(i) {
			t.Errorf("after merge, Contains(%d) = false", i)
		}
	}
}

func TestFilter_MergeMismatch(t *testing.T) {
	a := mustNew[int](t, 1000, 0.01)
	tests := map[string]*bloom.Filter[int]{
		"different capacity": mustNew[int](t, 2000, 0.01),
		"different rate":     mustNew[int](t, 1000, 0.02),
		"different seed":     mustNew[int](t, 1000, 0.01, bloom.WithSeed[int](99)),
	}
	for name, b := range tests {
		t.Run(name, func(t *testing.T) {
			if err := a.Merge(b); !errors.Is(err, bloom.ErrInvalidConfig) {
				t.Errorf("Merge error = %v, want ErrInvalidConfig", err)
			}
		})
	}
}

func TestFilter_MergeNil(t *testing.T) {
	a := mustNew[int](t, 100, 0.01)
	if err := a.Merge(nil); !errors.Is(err, bloom.ErrInvalidConfig) {
		t.Errorf("Merge(nil) error = %v, want ErrInvalidConfig", err)
	}
}

func TestFilter_Clear(t *testing.T) {
	f := mustNew[int](t, 100, 0.01)
	for i := 0; i < 50; i++ {
		f.Add(i)
	}
	f.Clear()
	for i := 0; i < 50; i++ {
		if f.Contains(i) {
			t.Errorf("after Clear, Contains(%d) = true", i)
		}
	}
	if got := f.ApproxCount(); got != 0 {
		t.Errorf("after Clear, ApproxCount() = %d, want 0", got)
	}
}

func TestFilter_ApproxCount(t *testing.T) {
	const n = 5000
	f := mustNew[int](t, n, 0.01)
	for i := 0; i < n; i++ {
		f.Add(i)
	}
	got := f.ApproxCount()
	// The estimator should land within 5% of the true distinct count.
	lo, hi := uint64(float64(n)*0.95), uint64(float64(n)*1.05)
	if got < lo || got > hi {
		t.Errorf("ApproxCount() = %d, want within [%d,%d]", got, lo, hi)
	}
}

func TestFilter_EstimatedFalsePositiveRate(t *testing.T) {
	f := mustNew[int](t, 1000, 0.01)
	if got := f.EstimatedFalsePositiveRate(); got != 0 {
		t.Errorf("empty filter FPR = %v, want 0", got)
	}
	for i := 0; i < 1000; i++ {
		f.Add(i)
	}
	// At design capacity the load-based estimate should be near the target.
	if got := f.EstimatedFalsePositiveRate(); got < 0.002 || got > 0.05 {
		t.Errorf("FPR at capacity = %v, want roughly 0.01", got)
	}
}

func TestFilter_CustomHasher(t *testing.T) {
	// A trivial hasher that ignores the seed; proves WithHasher is wired in.
	calls := 0
	h := func(seed uint64, v int) uint64 {
		calls++
		return uint64(v) * 2654435761
	}
	f := mustNew[int](t, 100, 0.01, bloom.WithHasher[int](h))
	f.Add(1)
	if !f.Contains(1) {
		t.Error("custom hasher: Contains(1) = false after Add")
	}
	if calls == 0 {
		t.Error("custom hasher was never called")
	}
}

func TestFilter_LooseRateClampsHashes(t *testing.T) {
	// A very loose target rate drives the rounded optimal k below 1; it must
	// clamp to at least one hash function.
	f := mustNew[int](t, 100, 0.9)
	if got := f.HashCount(); got < 1 {
		t.Errorf("HashCount() = %d, want >= 1", got)
	}
}

func TestFilter_ApproxCountSaturates(t *testing.T) {
	// A two-bit filter (n=1, p=0.5) saturates once both bits are set; the
	// estimator then reports the maximum rather than diverging.
	f := mustNew[int](t, 1, 0.5)
	for i := 0; i < 200; i++ {
		f.Add(i)
	}
	if f.BitSize() != 2 {
		t.Fatalf("expected a 2-bit filter, got %d bits", f.BitSize())
	}
	if got := f.ApproxCount(); got != math.MaxUint64 {
		t.Errorf("saturated ApproxCount() = %d, want MaxUint64", got)
	}
}

func ExampleFilter() {
	f, _ := bloom.New[string](1000, 0.01)
	f.Add("alice")
	f.Add("bob")

	fmt.Println(f.Contains("alice"))
	fmt.Println(f.Contains("carol")) // never added → definitely false
	// Output:
	// true
	// false
}

func ExampleFilter_Merge() {
	west, _ := bloom.New[string](1000, 0.01)
	east, _ := bloom.New[string](1000, 0.01)
	west.Add("user-1")
	east.Add("user-2")

	_ = west.Merge(east) // west now covers both shards
	fmt.Println(west.Contains("user-1"), west.Contains("user-2"))
	// Output:
	// true true
}

// TestFilter_SignedZeroNoFalseNegative guards the no-false-negatives guarantee
// for floats: +0 == -0, so after adding one the filter must report the other
// present. Every NaN is treated as the same element (see sketchhash).
func TestFilter_SignedZeroNoFalseNegative(t *testing.T) {
	negZero := math.Copysign(0, -1)
	tests := []struct {
		name     string
		contains func() bool
	}{
		{"float64 -0 after +0", func() bool {
			f := mustNew[float64](t, 1000, 0.01)
			f.Add(0)
			return f.Contains(negZero)
		}},
		{"float64 +0 after -0", func() bool {
			f := mustNew[float64](t, 1000, 0.01)
			f.Add(negZero)
			return f.Contains(0)
		}},
		{"float32 -0 after +0", func() bool {
			f := mustNew[float32](t, 1000, 0.01)
			f.Add(0)
			return f.Contains(float32(negZero))
		}},
		{"NaN with another payload", func() bool {
			f := mustNew[float64](t, 1000, 0.01)
			f.Add(math.NaN())
			return f.Contains(math.Float64frombits(0x7FF0_0000_0000_0001))
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.contains() {
				t.Error("Contains = false, want true (an equal value was added)")
			}
		})
	}
}

// TestFilter_InterfaceKeysDistinguishTypes checks that a Filter[any] does not
// report a value of another dynamic type as present just because it has the
// same numeric value: any(int(5)) and any(int64(5)) are distinct keys under ==.
// The rate is set low enough that a false positive here is not plausible.
func TestFilter_InterfaceKeysDistinguishTypes(t *testing.T) {
	f := mustNew[any](t, 100, 1e-9)
	f.Add(int(5))
	f.Add(1.0)
	tests := []struct {
		name  string
		probe any
		want  bool
	}{
		{"same type and value is present", int(5), true},
		{"float64 is present", 1.0, true},
		{"int64 of the same value is absent", int64(5), false},
		{"uint8 of the same value is absent", uint8(5), false},
		{"uintptr of the same value is absent", uintptr(5), false},
		{"uint64 bit pattern of the float is absent", math.Float64bits(1.0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.Contains(tt.probe)
			if got != tt.want {
				t.Errorf("Contains(%T(%v)) = %v, want %v", tt.probe, tt.probe, got, tt.want)
			}
		})
	}
}

// TestNew_SizeCapFollowsAddressSpace checks that New rejects, with
// ErrInvalidConfig, a filter whose bit array could not be allocated on this
// platform, rather than panicking in make. The fixed 2^40-bit cap assumed a
// 64-bit address space; on a 32-bit platform 2e9 items at a 1e-12 rate need
// ~1.2e11 bits (14 GB), under that cap but far beyond what make can address.
// The tighter cap only binds on 32-bit platforms; the 64-bit cap is exercised
// by TestNew_InvalidConfig.
func TestNew_SizeCapFollowsAddressSpace(t *testing.T) {
	if strconv.IntSize != 32 {
		t.Skip("the address-space cap is tighter than the 2^40-bit cap only on 32-bit platforms")
	}
	_, err := bloom.New[int](2_000_000_000, 1e-12)
	if !errors.Is(err, bloom.ErrInvalidConfig) {
		t.Errorf("New(2e9, 1e-12) error = %v, want ErrInvalidConfig", err)
	}
}
