package preprocessing_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

// TestRejectedFitResetsEstimator pins the contract that invalid input "leaves
// the estimator unfitted" even after an earlier successful Fit: a rejected
// re-fit must discard the old parameters rather than silently keep
// transforming with them.
func TestRejectedFitResetsEstimator(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	isEmpty := func(s string) bool { return s == "" }
	tests := []struct {
		name string
		// refit fits valid data, then invalid data, and reports whether the
		// estimator still transforms and the learned parameters it exposes.
		refit func() (ok bool, params any)
		want  any
	}{
		{name: "StandardScaler", refit: func() (bool, any) {
			s := preprocessing.NewStandardScaler().Fit([]float64{1, 2, 3}).Fit(nil)
			_, ok := s.Transform([]float64{1})
			return ok, []float64{s.Mean(), s.StdDev()}
		}, want: []float64{0, 0}},
		{name: "MinMaxScaler", refit: func() (bool, any) {
			s := preprocessing.NewMinMaxScaler().Fit([]float64{1, 2, 3}).Fit([]float64{})
			_, ok := s.Transform([]float64{1})
			return ok, []float64{s.Min(), s.Max()}
		}, want: []float64{0, 0}},
		{name: "RobustScaler", refit: func() (bool, any) {
			s := preprocessing.NewRobustScaler().Fit([]float64{1, 2, 3, 4, 5}).Fit([]float64{nan, 100, 200})
			_, ok := s.Transform([]float64{1})
			return ok, []float64{s.Median(), s.IQR()}
		}, want: []float64{0, 0}},
		{name: "OneHotEncoder", refit: func() (bool, any) {
			e := preprocessing.NewOneHotEncoder[string]().Fit([]string{"a", "b"}).Fit(nil)
			_, ok := e.Transform([]string{"a"})
			return ok, e.Categories()
		}, want: []string{}},
		{name: "LabelEncoder", refit: func() (bool, any) {
			e := preprocessing.NewLabelEncoder[string]().Fit([]string{"a", "b"}).Fit(nil)
			_, ok := e.Transform([]string{"a"})
			_, invOK := e.InverseTransform([]int{0})
			return ok || invOK, e.Categories()
		}, want: []string{}},
		{name: "OrdinalEncoder with learned categories", refit: func() (bool, any) {
			e := preprocessing.NewOrdinalEncoder[string]().Fit([]string{"a", "b"}).Fit(nil)
			_, ok := e.Transform([]string{"a"})
			return ok, e.Categories()
		}, want: []string{}},
		{name: "TargetEncoder", refit: func() (bool, any) {
			e := preprocessing.NewTargetEncoder[string]().Fit([]string{"a", "b"}, []float64{1, 3}).Fit([]string{"a"}, []float64{nan})
			_, ok := e.Transform([]string{"a"})
			return ok, e.GlobalMean()
		}, want: 0.0},
		{name: "MeanImputer", refit: func() (bool, any) {
			im := preprocessing.NewMeanImputer(nil).Fit([]float64{1, 2, 3}).Fit([]float64{nan, nan})
			_, ok := im.Transform([]float64{nan})
			return ok, im.Fill()
		}, want: 0.0},
		{name: "MedianImputer", refit: func() (bool, any) {
			im := preprocessing.NewMedianImputer(nil).Fit([]float64{1, 2, 3}).Fit([]float64{inf, 1})
			_, ok := im.Transform([]float64{nan})
			return ok, im.Fill()
		}, want: 0.0},
		{name: "ModeImputer", refit: func() (bool, any) {
			im := preprocessing.NewModeImputer(isEmpty).Fit([]string{"x", "x", "y"}).Fit([]string{"", ""})
			_, ok := im.Transform([]string{""})
			return ok, im.Fill()
		}, want: ""},
		{name: "FixedWidthBinner", refit: func() (bool, any) {
			b := preprocessing.NewFixedWidthBinner(3).Fit([]float64{0, 3, 6}).Fit(nil)
			_, ok := b.Transform([]float64{1})
			return ok, b.Edges()
		}, want: []float64{}},
		{name: "QuantileBinner", refit: func() (bool, any) {
			b := preprocessing.NewQuantileBinner(2).Fit([]float64{1, 2, 3, 4}).Fit([]float64{inf, 5})
			_, ok := b.Transform([]float64{1})
			return ok, b.Edges()
		}, want: []float64{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, params := tt.refit()
			if ok {
				t.Errorf("Transform ok = true after a rejected re-fit, want false")
			}
			if !reflect.DeepEqual(params, tt.want) {
				t.Errorf("learned parameters = %v after a rejected re-fit, want %v", params, tt.want)
			}
		})
	}
}

// TestExplicitOrdinalEncoderIgnoresFit checks the one estimator whose Fit is a
// documented no-op keeps its caller-supplied order even when handed empty data.
func TestExplicitOrdinalEncoderIgnoresFit(t *testing.T) {
	e := preprocessing.NewOrdinalEncoder("low", "high").Fit(nil)
	got, ok := e.Transform([]string{"high"})
	if !ok || !reflect.DeepEqual(got, []int{1}) {
		t.Errorf("Transform() = %v, %v; want [1], true", got, ok)
	}
}
