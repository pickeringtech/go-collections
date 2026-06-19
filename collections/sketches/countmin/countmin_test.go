package countmin_test

import (
	"errors"
	"fmt"
	"math"
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
	}{
		{"zero epsilon", 0, 0.01},
		{"epsilon one", 1, 0.01},
		{"negative epsilon", -0.1, 0.01},
		{"zero delta", 0.01, 0},
		{"delta one", 0.01, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := countmin.New[string](tc.eps, tc.delta)
			if !errors.Is(err, countmin.ErrInvalidConfig) {
				t.Errorf("error = %v, want ErrInvalidConfig", err)
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
