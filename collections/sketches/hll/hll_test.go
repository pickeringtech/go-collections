package hll_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/hll"
)

func mustNew[T comparable](t *testing.T, opts ...hll.Option[T]) *hll.Sketch[T] {
	t.Helper()
	s, err := hll.New(opts...)
	if err != nil {
		t.Fatalf("New: unexpected error %v", err)
	}
	return s
}

// relErr returns the relative error of got against the true value.
func relErr(got uint64, truth int) float64 {
	return math.Abs(float64(got)-float64(truth)) / float64(truth)
}

func TestNew_InvalidPrecision(t *testing.T) {
	tests := []struct {
		name string
		p    int
	}{
		{"too low", hll.MinPrecision - 1},
		{"too high", hll.MaxPrecision + 1},
		{"zero", 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := hll.New[int](hll.WithPrecision[int](tc.p))
			if !errors.Is(err, hll.ErrInvalidConfig) {
				t.Errorf("error = %v, want ErrInvalidConfig", err)
			}
		})
	}
}

func TestSketch_Defaults(t *testing.T) {
	s := mustNew[int](t)
	if got := s.Precision(); got != hll.DefaultPrecision {
		t.Errorf("Precision() = %d, want %d", got, hll.DefaultPrecision)
	}
	if got := s.RegisterCount(); got != 1<<hll.DefaultPrecision {
		t.Errorf("RegisterCount() = %d, want %d", got, 1<<hll.DefaultPrecision)
	}
}

func TestSketch_EmptyCountZero(t *testing.T) {
	s := mustNew[int](t)
	if got := s.Count(); got != 0 {
		t.Errorf("empty Count() = %d, want 0", got)
	}
}

func TestSketch_AccuracyAcrossScales(t *testing.T) {
	for _, n := range []int{100, 1000, 10000, 100000} {
		t.Run(fmt.Sprintf("n_%d", n), func(t *testing.T) {
			s := mustNew[int](t)
			for i := 0; i < n; i++ {
				s.Add(i)
			}
			got := s.Count()
			// Allow ~3 standard errors plus the small-range slack.
			tolerance := 3*s.StandardError() + 0.01
			if e := relErr(got, n); e > tolerance {
				t.Errorf("Count() = %d for n=%d, relative error %.4f exceeds tolerance %.4f", got, n, e, tolerance)
			}
		})
	}
}

func TestSketch_Idempotent(t *testing.T) {
	s := mustNew[int](t)
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	before := s.Count()
	// Re-adding the same elements must not change the distinct-count estimate.
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	if after := s.Count(); after != before {
		t.Errorf("Count changed after re-adding duplicates: %d -> %d", before, after)
	}
}

func TestSketch_Merge(t *testing.T) {
	a := mustNew[int](t)
	b := mustNew[int](t)
	// Disjoint halves; the union has 20000 distinct elements.
	for i := 0; i < 10000; i++ {
		a.Add(i)
	}
	for i := 10000; i < 20000; i++ {
		b.Add(i)
	}
	if err := a.Merge(b); err != nil {
		t.Fatalf("Merge: %v", err)
	}
	got := a.Count()
	tolerance := 3*a.StandardError() + 0.01
	if e := relErr(got, 20000); e > tolerance {
		t.Errorf("merged Count() = %d, relative error %.4f exceeds tolerance %.4f", got, e, tolerance)
	}
}

func TestSketch_MergeOverlap(t *testing.T) {
	a := mustNew[int](t)
	b := mustNew[int](t)
	// Fully overlapping; the union is still 10000 distinct.
	for i := 0; i < 10000; i++ {
		a.Add(i)
		b.Add(i)
	}
	if err := a.Merge(b); err != nil {
		t.Fatalf("Merge: %v", err)
	}
	got := a.Count()
	tolerance := 3*a.StandardError() + 0.01
	if e := relErr(got, 10000); e > tolerance {
		t.Errorf("merged Count() = %d, relative error %.4f exceeds tolerance %.4f", got, e, tolerance)
	}
}

func TestSketch_MergeMismatch(t *testing.T) {
	a := mustNew[int](t)
	tests := map[string]*hll.Sketch[int]{
		"different precision": mustNew[int](t, hll.WithPrecision[int](12)),
		"different seed":      mustNew[int](t, hll.WithSeed[int](99)),
	}
	for name, b := range tests {
		t.Run(name, func(t *testing.T) {
			if err := a.Merge(b); !errors.Is(err, hll.ErrInvalidConfig) {
				t.Errorf("Merge error = %v, want ErrInvalidConfig", err)
			}
		})
	}
}

func TestSketch_Clear(t *testing.T) {
	s := mustNew[int](t)
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	s.Clear()
	if got := s.Count(); got != 0 {
		t.Errorf("after Clear, Count() = %d, want 0", got)
	}
}

func TestSketch_LowPrecisionStillWorks(t *testing.T) {
	s := mustNew[int](t, hll.WithPrecision[int](hll.MinPrecision))
	for i := 0; i < 50; i++ {
		s.Add(i)
	}
	// Coarse precision, but the estimate should be in the right ballpark.
	if got := s.Count(); got < 20 || got > 100 {
		t.Errorf("Count() = %d for 50 distinct at min precision, want roughly 50", got)
	}
}

func TestSketch_CustomHasher(t *testing.T) {
	calls := 0
	h := func(seed uint64, v int) uint64 {
		calls++
		return uint64(v) * 2654435761
	}
	s := mustNew[int](t, hll.WithHasher[int](h))
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	if calls == 0 {
		t.Error("custom hasher was never called")
	}
	if s.Count() == 0 {
		t.Error("Count() = 0 after 1000 adds with custom hasher")
	}
}

func TestSketch_MergeNil(t *testing.T) {
	s := mustNew[int](t)
	if err := s.Merge(nil); !errors.Is(err, hll.ErrInvalidConfig) {
		t.Errorf("Merge(nil) error = %v, want ErrInvalidConfig", err)
	}
}

// TestSketch_AlphaConstants exercises the small-m bias-correction constants by
// estimating at the precisions whose register counts (16, 32, 64) have special
// alpha values.
func TestSketch_AlphaConstants(t *testing.T) {
	for _, p := range []int{4, 5, 6} { // m = 16, 32, 64
		t.Run(fmt.Sprintf("p_%d", p), func(t *testing.T) {
			s := mustNew[int](t, hll.WithPrecision[int](p))
			for i := 0; i < 200; i++ {
				s.Add(i)
			}
			// Coarse registers, but the estimate must be a sane positive number.
			if got := s.Count(); got == 0 {
				t.Errorf("Count() = 0 at precision %d", p)
			}
		})
	}
}

func ExampleSketch() {
	s, _ := hll.New[string]()
	for i := 0; i < 100000; i++ {
		s.Add(fmt.Sprintf("user-%d", i))
	}
	// Estimate is within ~1% of the true 100000 distinct users.
	within := math.Abs(float64(s.Count())-100000)/100000 < 0.02
	fmt.Println("within 2%:", within)
	// Output:
	// within 2%: true
}
