package countmin_test

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/countmin"
)

func mustNew[T comparable](t *testing.T, eps, delta float64, opts ...countmin.Option[T]) *countmin.Sketch[T] {
	t.Helper()
	s, err := countmin.New(eps, delta, opts...)
	if err != nil {
		t.Fatalf("New(%v, %v): unexpected error %v", eps, delta, err)
	}
	return s
}

func TestNew_InvalidConfig(t *testing.T) {
	tests := []struct {
		name       string
		eps, delta float64
		// wantMsg is a fragment the error must carry, so the caller is told
		// which bound is wrong rather than only that something is.
		wantMsg string
	}{
		// A zero bound would also be caught by the counter cap (it needs
		// infinitely many counters), but it must be reported as out of range.
		{"zero epsilon is reported as out of range", 0, 0.01, "epsilon must be in (0,1)"},
		{"epsilon one is reported as out of range", 1, 0.01, "epsilon must be in (0,1)"},
		{"negative epsilon is reported as out of range", -0.1, 0.01, "epsilon must be in (0,1)"},
		{"zero delta is reported as out of range", 0.01, 0, "delta must be in (0,1)"},
		{"delta one is reported as out of range", 0.01, 1, "delta must be in (0,1)"},
		// Absurd bounds: w*d exceeds the allocation cap, so New rejects it
		// instead of overflowing or panicking on allocation.
		{"astronomically tight epsilon exceeds the counter cap", 1e-18, 0.01, "counter limit"},
		// w = ceil(e/1e-12) = 2718281828460 and d = ceil(ln 100) = 5, so the
		// table needs w·d = 13591409142300 counters — over the 2^40 cap even
		// though w alone is under it — and the error reports that product.
		{"bounds whose product exceeds the counter cap report the product", 1e-12, 0.01, "need 13591409142300 counters"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := countmin.New[string](tc.eps, tc.delta)
			if !errors.Is(err, countmin.ErrInvalidConfig) {
				t.Fatalf("error = %v, want ErrInvalidConfig", err)
			}
			if !strings.Contains(err.Error(), tc.wantMsg) {
				t.Errorf("error = %q, want it to contain %q", err, tc.wantMsg)
			}
		})
	}
}

func TestSketch_NeverUnderEstimates(t *testing.T) {
	s := mustNew[int](t, 0.001, 0.001)
	truth := make(map[int]uint64)
	// Add a skewed stream: small keys appear far more often.
	for i := 0; i < 100000; i++ {
		key := i % 500
		s.Add(key)
		truth[key]++
	}
	for key, want := range truth {
		got := s.Estimate(key)
		if got < want {
			t.Fatalf("Estimate(%d) = %d under true count %d; Count-Min must never under-estimate", key, got, want)
		}
	}
}

func TestSketch_ErrorWithinBound(t *testing.T) {
	const eps = 0.001
	s := mustNew[int](t, eps, 0.001)
	truth := make(map[int]uint64)
	for i := 0; i < 100000; i++ {
		key := i % 1000
		s.Add(key)
		truth[key]++
	}
	bound := uint64(eps*float64(s.Total())) + 1
	over := 0
	for key, want := range truth {
		got := s.Estimate(key)
		if got-want > bound {
			over++
		}
	}
	// At most delta fraction may exceed the bound; allow a little slack.
	if over > len(truth)/100+5 {
		t.Errorf("%d/%d estimates exceeded the epsilon·N bound %d", over, len(truth), bound)
	}
}

func TestSketch_AddCount(t *testing.T) {
	s := mustNew[string](t, 0.01, 0.01)
	s.AddCount("x", 1000)
	if got := s.Estimate("x"); got < 1000 {
		t.Errorf("Estimate after AddCount(1000) = %d, want >= 1000", got)
	}
	if got := s.Total(); got != 1000 {
		t.Errorf("Total() = %d, want 1000", got)
	}
}

func TestSketch_AbsentElementSmall(t *testing.T) {
	s := mustNew[int](t, 0.001, 0.001)
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	// An element never added estimates at most the epsilon·N error, often 0.
	if got := s.Estimate(99999); got > uint64(0.001*float64(s.Total()))+1 {
		t.Errorf("Estimate(absent) = %d, larger than the error bound", got)
	}
}

func TestSketch_Merge(t *testing.T) {
	a := mustNew[string](t, 0.01, 0.01)
	b := mustNew[string](t, 0.01, 0.01)
	a.AddCount("x", 10)
	b.AddCount("x", 5)
	b.AddCount("y", 7)
	if err := a.Merge(b); err != nil {
		t.Fatalf("Merge: %v", err)
	}
	if got := a.Estimate("x"); got < 15 {
		t.Errorf("merged Estimate(x) = %d, want >= 15", got)
	}
	if got := a.Estimate("y"); got < 7 {
		t.Errorf("merged Estimate(y) = %d, want >= 7", got)
	}
	if got := a.Total(); got != 22 {
		t.Errorf("merged Total() = %d, want 22", got)
	}
}

// TestSketch_MergeAddsCountsExactly checks that Merge adds counters exactly
// when neither side is near saturation, whichever side holds the larger count.
// Each sketch holds only "x", so every counter is either 0 or x's count and the
// estimate is exact. The saturation guard must only fire when the sum would
// actually overflow, never merely because the receiver's count is the larger.
func TestSketch_MergeAddsCountsExactly(t *testing.T) {
	type args struct {
		receiver, other uint64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{name: "a larger receiver count adds the smaller one", args: args{receiver: 10, other: 5}, want: 15},
		{name: "a smaller receiver count adds the larger one", args: args{receiver: 5, other: 10}, want: 15},
		{name: "equal counts double", args: args{receiver: 7, other: 7}, want: 14},
		{name: "an empty receiver takes the other count", args: args{receiver: 0, other: 7}, want: 7},
		{name: "an empty other leaves the receiver count", args: args{receiver: 7, other: 0}, want: 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := mustNew[string](t, 0.01, 0.01)
			b := mustNew[string](t, 0.01, 0.01)
			a.AddCount("x", tt.args.receiver)
			b.AddCount("x", tt.args.other)
			err := a.Merge(b)
			if err != nil {
				t.Fatalf("Merge: %v", err)
			}
			got := a.Estimate("x")
			if got != tt.want {
				t.Errorf("merged Estimate(x) = %d, want %d", got, tt.want)
			}
			total := a.Total()
			if total != tt.want {
				t.Errorf("merged Total() = %d, want %d", total, tt.want)
			}
		})
	}
}

func TestSketch_MergeMismatch(t *testing.T) {
	a := mustNew[int](t, 0.01, 0.01)
	tests := map[string]*countmin.Sketch[int]{
		"different epsilon": mustNew[int](t, 0.02, 0.01),
		"different delta":   mustNew[int](t, 0.01, 0.1),
		"different seed":    mustNew[int](t, 0.01, 0.01, countmin.WithSeed[int](99)),
	}
	for name, b := range tests {
		t.Run(name, func(t *testing.T) {
			if err := a.Merge(b); !errors.Is(err, countmin.ErrInvalidConfig) {
				t.Errorf("Merge error = %v, want ErrInvalidConfig", err)
			}
		})
	}
}

func TestSketch_Clear(t *testing.T) {
	s := mustNew[int](t, 0.01, 0.01)
	s.AddCount(1, 100)
	s.Clear()
	if got := s.Estimate(1); got != 0 {
		t.Errorf("after Clear, Estimate(1) = %d, want 0", got)
	}
	if got := s.Total(); got != 0 {
		t.Errorf("after Clear, Total() = %d, want 0", got)
	}
}

func TestSketch_Dimensions(t *testing.T) {
	// w = ceil(e/0.01) = 272, d = ceil(ln(1/0.01)) = 5.
	s := mustNew[int](t, 0.01, 0.01)
	if got := s.Width(); got != 272 {
		t.Errorf("Width() = %d, want 272", got)
	}
	if got := s.Depth(); got != 5 {
		t.Errorf("Depth() = %d, want 5", got)
	}
}

func TestSketch_CustomHasher(t *testing.T) {
	calls := 0
	h := func(seed uint64, v int) uint64 {
		calls++
		return uint64(v) * 2654435761
	}
	s := mustNew[int](t, 0.01, 0.01, countmin.WithHasher[int](h))
	s.AddCount(1, 3)
	if got := s.Estimate(1); got < 3 {
		t.Errorf("custom hasher: Estimate(1) = %d, want >= 3", got)
	}
	if calls == 0 {
		t.Error("custom hasher was never called")
	}
}

func TestSketch_CounterSaturation(t *testing.T) {
	s := mustNew[string](t, 0.01, 0.01)
	s.AddCount("x", math.MaxUint64) // counter and total fill to the brim
	s.AddCount("x", 5)              // both must saturate, not wrap
	if got := s.Estimate("x"); got != math.MaxUint64 {
		t.Errorf("saturated Estimate(x) = %d, want MaxUint64", got)
	}
	if got := s.Total(); got != math.MaxUint64 {
		t.Errorf("saturated Total() = %d, want MaxUint64", got)
	}
}

func TestSketch_MergeNil(t *testing.T) {
	s := mustNew[int](t, 0.01, 0.01)
	if err := s.Merge(nil); !errors.Is(err, countmin.ErrInvalidConfig) {
		t.Errorf("Merge(nil) error = %v, want ErrInvalidConfig", err)
	}
}

func TestSketch_MergeSaturation(t *testing.T) {
	a := mustNew[string](t, 0.01, 0.01)
	b := mustNew[string](t, 0.01, 0.01)
	a.AddCount("x", math.MaxUint64)
	b.AddCount("x", math.MaxUint64)
	if err := a.Merge(b); err != nil {
		t.Fatalf("Merge: %v", err)
	}
	if got := a.Estimate("x"); got != math.MaxUint64 {
		t.Errorf("merged saturated Estimate(x) = %d, want MaxUint64", got)
	}
	if got := a.Total(); got != math.MaxUint64 {
		t.Errorf("merged saturated Total() = %d, want MaxUint64", got)
	}
}

func ExampleSketch() {
	s, _ := countmin.New[string](0.001, 0.01)
	for _, word := range []string{"go", "go", "rust", "go", "rust", "zig"} {
		s.Add(word)
	}
	fmt.Println("go:", s.Estimate("go"))
	fmt.Println("rust:", s.Estimate("rust"))
	fmt.Println("zig:", s.Estimate("zig"))
	// Output:
	// go: 3
	// rust: 2
	// zig: 1
}

// TestSketch_SignedZeroNeverUnderCounts guards the never-under-count guarantee
// for floats: +0 == -0, so their occurrences are one element's count.
func TestSketch_SignedZeroNeverUnderCounts(t *testing.T) {
	s := mustNew[float64](t, 0.01, 0.01)
	s.AddCount(0, 10)
	s.AddCount(math.Copysign(0, -1), 5)
	got := s.Estimate(math.Copysign(0, -1))
	if got < 15 {
		t.Errorf("Estimate(-0) = %d after 10×(+0) and 5×(-0), want >= 15", got)
	}
}

// TestSketch_FailureRateWithinDelta checks the (epsilon, delta) guarantee
// empirically on a narrow table: epsilon=0.17, delta=1e-4 gives a 16×10
// sketch. One heavy hitter carries 90% of N, so an absent key's estimate
// exceeds epsilon·N only if every row maps it onto the heavy hitter's column.
// With independent rows that is a (1/16)^10 event; when every row's column was
// (h1 + row·h2) mod w, a key matching the heavy hitter's (h1, h2) mod w collided
// in all rows at once, flooring the failure rate near 1/128 — 77× delta.
func TestSketch_FailureRateWithinDelta(t *testing.T) {
	const (
		epsilon = 0.17
		delta   = 1e-4
		queries = 200_000
	)
	s := mustNew[uint64](t, epsilon, delta)
	if s.Width() != 16 || s.Depth() != 10 {
		t.Fatalf("dimensions = %d×%d, want 16×10", s.Width(), s.Depth())
	}
	s.AddCount(1<<62, 900)
	for i := uint64(0); i < 100; i++ {
		s.Add(1<<61 + i)
	}
	bound := epsilon * float64(s.Total())
	failures := 0
	for i := uint64(0); i < queries; i++ {
		// Keys never added: their true count is 0.
		if float64(s.Estimate(1<<40+i)) > bound {
			failures++
		}
	}
	rate := float64(failures) / queries
	if rate > delta {
		t.Errorf("P(overshoot > epsilon·N) = %.5f over %d absent keys, want <= delta %v", rate, queries, delta)
	}
}

// TestNew_SizeCapFollowsAddressSpace checks that New rejects, with
// ErrInvalidConfig, a counter table that could not be allocated on this
// platform, rather than panicking in make. The fixed 2^40-counter cap assumed a
// 64-bit address space; on a 32-bit platform epsilon=1e-9 needs ~1.4e10
// counters (108 GB), under that cap but far beyond what make can address. The
// tighter cap only binds on 32-bit platforms; the 64-bit cap is exercised by
// TestNew_InvalidConfig.
func TestNew_SizeCapFollowsAddressSpace(t *testing.T) {
	if strconv.IntSize != 32 {
		t.Skip("the address-space cap is tighter than the 2^40-counter cap only on 32-bit platforms")
	}
	_, err := countmin.New[int](1e-9, 0.01)
	if !errors.Is(err, countmin.ErrInvalidConfig) {
		t.Errorf("New(1e-9, 0.01) error = %v, want ErrInvalidConfig", err)
	}
}
