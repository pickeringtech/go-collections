package preprocessing_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
	"github.com/pickeringtech/go-collections/stats"
)

func TestStandardScaler(t *testing.T) {
	// Train data has mean 5 and population stddev 2, so z = (x − 5) / 2.
	train := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	tests := []struct {
		name  string
		input []float64
		want  []float64
		ok    bool
	}{
		{name: "train data standardizes", input: train, want: []float64{-1.5, -0.5, -0.5, -0.5, 0, 0, 1, 2}, ok: true},
		{name: "test data reuses train params", input: []float64{5, 7}, want: []float64{0, 1}, ok: true},
		{name: "empty input is fine once fitted", input: []float64{}, want: []float64{}, ok: true},
		{name: "nil input is fine once fitted", input: nil, want: []float64{}, ok: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := preprocessing.NewStandardScaler().Fit(train)
			got, ok := s.Transform(tc.input)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if !floatSlicesClose(got, tc.want) {
				t.Fatalf("Transform = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestStandardScalerMatchesStandardize guards the doc's promise that the
// scaler mirrors stats.Standardize: an uncompensated Σx/n mean drifts on many
// large values (by 0.02σ here), shifting every z-score. Both now share the
// same mean and standard deviation, so the results agree exactly.
func TestStandardScalerMatchesStandardize(t *testing.T) {
	data := make([]float64, 1_000_000)
	for i := range data {
		data[i] = 1e9 + float64(i%3) + 0.1
	}
	got, ok := preprocessing.NewStandardScaler().FitTransform(data)
	want, _ := stats.Standardize(data)
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	worst := 0.0
	for i := range got {
		worst = math.Max(worst, math.Abs(got[i]-want[i]))
	}
	if worst != 0 {
		t.Errorf("max |scaler − stats.Standardize| = %v, want 0", worst)
	}
}

func TestStandardScalerUnfitted(t *testing.T) {
	got, ok := preprocessing.NewStandardScaler().Transform([]float64{1, 2, 3})
	if ok || got != nil {
		t.Fatalf("unfitted Transform = (%v, %v), want (nil, false)", got, ok)
	}
}

func TestStandardScalerConstantTrainMapsToZero(t *testing.T) {
	got, ok := preprocessing.NewStandardScaler().FitTransform([]float64{3, 3, 3})
	if !ok || !floatSlicesClose(got, []float64{0, 0, 0}) {
		t.Fatalf("FitTransform = (%v, %v), want ([0 0 0], true)", got, ok)
	}
}

func TestStandardScalerEmptyFitLeavesUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewStandardScaler().Fit(nil).Transform([]float64{1}); ok {
		t.Fatalf("Transform after empty Fit reported ok")
	}
}

func TestStandardScalerAccessors(t *testing.T) {
	s := preprocessing.NewStandardScaler().Fit([]float64{2, 4, 4, 4, 5, 5, 7, 9})
	if !floatsClose(s.Mean(), 5) {
		t.Fatalf("Mean() = %v, want 5", s.Mean())
	}
	if !floatsClose(s.StdDev(), 2) {
		t.Fatalf("StdDev() = %v, want 2", s.StdDev())
	}
}

func TestMinMaxScaler(t *testing.T) {
	train := []float64{2, 4, 6, 8, 10}
	tests := []struct {
		name  string
		input []float64
		want  []float64
		ok    bool
	}{
		{name: "train data scales to unit interval", input: train, want: []float64{0, 0.25, 0.5, 0.75, 1}, ok: true},
		{name: "test value above train max exceeds 1", input: []float64{12}, want: []float64{1.25}, ok: true},
		{name: "test value below train min is negative", input: []float64{0}, want: []float64{-0.25}, ok: true},
		{name: "empty input is fine once fitted", input: []float64{}, want: []float64{}, ok: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := preprocessing.NewMinMaxScaler().Fit(train).Transform(tc.input)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if !floatSlicesClose(got, tc.want) {
				t.Fatalf("Transform = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestMinMaxScalerMirrorsNormalize checks the scaler inherits Normalize's
// semantics: the learned range skips NaN (so a NaN only affects its own
// position) and a finite range whose width overflows still rescales.
func TestMinMaxScalerMirrorsNormalize(t *testing.T) {
	nan, maxF := math.NaN(), math.MaxFloat64
	tests := []struct {
		name  string
		train []float64
		input []float64
		want  []float64
	}{
		{name: "leading NaN in training data", train: []float64{nan, 0, 10}, input: []float64{nan, 0, 5, 10}, want: []float64{nan, 0, 0.5, 1}},
		{name: "overflowing training range", train: []float64{-maxF, maxF}, input: []float64{-maxF, 0, maxF}, want: []float64{0, 0.5, 1}},
		{name: "all-NaN training data propagates", train: []float64{nan, nan}, input: []float64{1, 2}, want: []float64{nan, nan}},
		{name: "NaN beside a degenerate range", train: []float64{nan, 5, 5}, input: []float64{nan, 5, 7}, want: []float64{nan, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := preprocessing.NewMinMaxScaler().Fit(tt.train).Transform(tt.input)
			if !ok || !floatSlicesClose(got, tt.want) {
				t.Errorf("Transform() = %v, %v; want %v, true", got, ok, tt.want)
			}
			norm, _ := stats.Normalize(tt.train)
			mirrored, _ := preprocessing.NewMinMaxScaler().FitTransform(tt.train)
			if !floatSlicesClose(mirrored, norm) {
				t.Errorf("FitTransform() = %v, want stats.Normalize's %v", mirrored, norm)
			}
		})
	}
}

func TestMinMaxScalerDegenerateRangeMapsToZero(t *testing.T) {
	got, ok := preprocessing.NewMinMaxScaler().FitTransform([]float64{7, 7, 7})
	if !ok || !floatSlicesClose(got, []float64{0, 0, 0}) {
		t.Fatalf("FitTransform = (%v, %v), want ([0 0 0], true)", got, ok)
	}
}

func TestMinMaxScalerUnfitted(t *testing.T) {
	if _, ok := preprocessing.NewMinMaxScaler().Transform([]float64{1}); ok {
		t.Fatalf("unfitted Transform reported ok")
	}
}

func TestMinMaxScalerAccessors(t *testing.T) {
	s := preprocessing.NewMinMaxScaler().Fit([]float64{2, 4, 6, 8, 10})
	if !floatsClose(s.Min(), 2) || !floatsClose(s.Max(), 10) {
		t.Fatalf("Min/Max = %v/%v, want 2/10", s.Min(), s.Max())
	}
}

func TestRobustScaler(t *testing.T) {
	// Median 5, Q1 3, Q3 7 (linear interpolation), so IQR 4: r = (x − 5) / 4.
	train := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	got, ok := preprocessing.NewRobustScaler().Fit(train).Transform([]float64{5, 9, 1})
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	want := []float64{0, 1, -1}
	if !floatSlicesClose(got, want) {
		t.Fatalf("Transform = %v, want %v", got, want)
	}
}

func TestRobustScalerRejectsNonFiniteTrain(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, ok := preprocessing.NewRobustScaler().Fit([]float64{1, 2, bad}).Transform([]float64{1}); ok {
			t.Fatalf("Fit accepted non-finite %v", bad)
		}
	}
}

func TestRobustScalerDegenerateIQRMapsToZero(t *testing.T) {
	got, ok := preprocessing.NewRobustScaler().FitTransform([]float64{4, 4, 4, 4})
	if !ok || !floatSlicesClose(got, []float64{0, 0, 0, 0}) {
		t.Fatalf("FitTransform = (%v, %v), want ([0 0 0 0], true)", got, ok)
	}
}

func TestRobustScalerAccessors(t *testing.T) {
	s := preprocessing.NewRobustScaler().Fit([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	if !floatsClose(s.Median(), 5) || !floatsClose(s.IQR(), 4) {
		t.Fatalf("Median/IQR = %v/%v, want 5/4", s.Median(), s.IQR())
	}
}

// StandardScaler and MinMaxScaler propagate non-finite values, matching the
// stats transforms they mirror.
func TestScalersPropagateNonFinite(t *testing.T) {
	stdGot, ok := preprocessing.NewStandardScaler().FitTransform([]float64{1, 2, math.Inf(1), 4})
	if !ok {
		t.Fatalf("StandardScaler ok = false, want true")
	}
	finite := false
	for _, v := range stdGot {
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			finite = true
		}
	}
	if finite {
		t.Fatalf("StandardScaler did not propagate non-finite: %v", stdGot)
	}

	mmGot, ok := preprocessing.NewMinMaxScaler().FitTransform([]float64{1, 2, math.NaN(), 4})
	if !ok {
		t.Fatalf("MinMaxScaler ok = false, want true")
	}
	if !math.IsNaN(mmGot[2]) {
		t.Fatalf("MinMaxScaler did not propagate NaN: %v", mmGot)
	}
}

// Transform must not mutate the caller's input slice.
func TestScalersDoNotMutateInput(t *testing.T) {
	input := []float64{1, 2, 3, 4}
	original := []float64{1, 2, 3, 4}
	preprocessing.NewStandardScaler().FitTransform(input)
	preprocessing.NewMinMaxScaler().FitTransform(input)
	preprocessing.NewRobustScaler().FitTransform(input)
	if !floatSlicesClose(input, original) {
		t.Fatalf("input mutated to %v, want %v", input, original)
	}
}

// The scalers satisfy the Transformer contract.
var (
	_ preprocessing.Transformer[float64, []float64] = (*preprocessing.StandardScaler)(nil)
	_ preprocessing.Transformer[float64, []float64] = (*preprocessing.MinMaxScaler)(nil)
	_ preprocessing.Transformer[float64, []float64] = (*preprocessing.RobustScaler)(nil)
)
