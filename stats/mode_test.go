package stats_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestMode(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		want   []int
		wantOK bool
	}{
		{name: "single clear mode", input: []int{1, 2, 2, 3}, want: []int{2}, wantOK: true},
		{name: "ties returned in first-appearance order", input: []int{3, 1, 3, 1, 2}, want: []int{3, 1}, wantOK: true},
		{name: "all unique are all modes", input: []int{4, 5, 6}, want: []int{4, 5, 6}, wantOK: true},
		{name: "single value", input: []int{9}, want: []int{9}, wantOK: true},
		{name: "empty input", input: []int{}, wantOK: false},
		{name: "nil input", input: nil, wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.Mode(tt.input)
			if ok != tt.wantOK || (ok && !equalInts(got, tt.want)) {
				t.Fatalf("Mode(%v) = (%v, %v), want (%v, %v)", tt.input, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestModeOnStrings(t *testing.T) {
	got, ok := stats.Mode([]string{"a", "b", "a", "c", "b", "a"})
	if !ok || len(got) != 1 || got[0] != "a" {
		t.Fatalf(`Mode(strings) = (%v, %v), want (["a"], true)`, got, ok)
	}
}

func TestModeRejectsNonFinite(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		got, ok := stats.Mode([]float64{1, 1, bad, 2})
		if ok {
			t.Fatalf("Mode([1 1 %v 2]) = (%v, %v), want (nil, false)", bad, got, ok)
		}
	}
	// float32 inputs are checked too.
	_, ok := stats.Mode([]float32{1, float32(math.Inf(1)), 1})
	if ok {
		t.Fatalf("Mode(float32 with +Inf) should be ok=false")
	}
}

func TestModeFiniteFloatsStillWork(t *testing.T) {
	got, ok := stats.Mode([]float64{1.5, 2.5, 1.5})
	if !ok || len(got) != 1 || got[0] != 1.5 {
		t.Fatalf("Mode([1.5 2.5 1.5]) = (%v, %v), want ([1.5], true)", got, ok)
	}
}

// Named float types share float64/float32's non-finite values, so the
// categorical NaN/Inf rejection must not depend on the exact dynamic type.
type namedFloat float64
type namedFloat32 float32

// TestCategoricalRejectsNonFiniteNamedFloats guards the non-finite check
// against named float types (type F float64), which a type switch on
// float64/float32 misses, and against NaN-carrying non-float comparables.
func TestCategoricalRejectsNonFiniteNamedFloats(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name string
		fn   func() bool
	}{
		{name: "Mode of named float NaNs", fn: func() bool { _, ok := stats.Mode([]namedFloat{namedFloat(nan), namedFloat(nan)}); return ok }},
		{name: "Mode of a named float infinity", fn: func() bool { _, ok := stats.Mode([]namedFloat{1, namedFloat(inf)}); return ok }},
		{name: "Mode of a named float32 negative infinity", fn: func() bool { _, ok := stats.Mode([]namedFloat32{namedFloat32(-inf)}); return ok }},
		{name: "Entropy of named float NaNs", fn: func() bool { _, ok := stats.Entropy([]namedFloat{namedFloat(nan), namedFloat(nan)}); return ok }},
		{name: "Gini of named float NaNs", fn: func() bool { _, ok := stats.Gini([]namedFloat{namedFloat(nan), namedFloat(nan)}); return ok }},
		{name: "Gini of a named float32 infinity", fn: func() bool { _, ok := stats.Gini([]namedFloat32{2, namedFloat32(inf)}); return ok }},
		{name: "Mode of a complex NaN", fn: func() bool { _, ok := stats.Mode([]complex128{complex(nan, 0), 1}); return ok }},
		{name: "Mode of an interface-held infinity", fn: func() bool { _, ok := stats.Mode([]any{"a", inf}); return ok }},
		{name: "Mode of an interface-held named float32 infinity", fn: func() bool { _, ok := stats.Mode([]any{namedFloat32(-inf)}); return ok }},
		{name: "Mode of an interface-held NaN", fn: func() bool { _, ok := stats.Mode([]any{nan, 1}); return ok }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fn() {
				t.Errorf("ok = true, want false")
			}
		})
	}
}

// Finite named floats, and finite values behind an interface, are ordinary
// categories.
func TestModeNamedFloatsStillWork(t *testing.T) {
	got, ok := stats.Mode([]namedFloat{1, 2, 2, 3})
	if !ok || !reflect.DeepEqual(got, []namedFloat{2}) {
		t.Errorf("Mode() = %v, %v; want [2], true", got, ok)
	}
	mixed, ok := stats.Mode([]any{"a", 1.5, 1.5, 2})
	if !ok || !reflect.DeepEqual(mixed, []any{1.5}) {
		t.Errorf("Mode() = %v, %v; want [1.5], true", mixed, ok)
	}
}
