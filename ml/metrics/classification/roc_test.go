package classification_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/ml/metrics/classification"
)

func TestAUC(t *testing.T) {
	tests := []struct {
		name   string
		yTrue  []int
		scores []float64
		pos    int
		want   float64
		wantOK bool
	}{
		{
			name:   "textbook example",
			yTrue:  []int{0, 0, 1, 1},
			scores: []float64{0.1, 0.4, 0.35, 0.8},
			pos:    1,
			want:   0.75,
			wantOK: true,
		},
		{
			name:   "perfect ranker",
			yTrue:  []int{0, 0, 1, 1},
			scores: []float64{0.1, 0.2, 0.8, 0.9},
			pos:    1,
			want:   1,
			wantOK: true,
		},
		{
			name:   "tie-averaged ranks",
			yTrue:  []int{0, 1, 1, 0},
			scores: []float64{0.5, 0.5, 0.6, 0.4},
			pos:    1,
			want:   0.875,
			wantOK: true,
		},
		{
			name:   "single class is undefined",
			yTrue:  []int{1, 1, 1},
			scores: []float64{0.2, 0.5, 0.9},
			pos:    1,
			wantOK: false,
		},
		{
			name:   "non-finite score rejected",
			yTrue:  []int{0, 1},
			scores: []float64{0.2, math.Inf(1)},
			pos:    1,
			wantOK: false,
		},
		{
			name:   "length mismatch rejected",
			yTrue:  []int{0, 1},
			scores: []float64{0.2},
			pos:    1,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := classification.AUC(tt.yTrue, tt.scores, tt.pos)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && !approxEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestROCCurve(t *testing.T) {
	yTrue := []int{0, 0, 1, 1}
	scores := []float64{0.1, 0.4, 0.35, 0.8}
	curve, ok := classification.ROCCurve(yTrue, scores, 1)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	// First point is the strict-threshold (0,0) corner; last is (1,1).
	first := curve.Points[0]
	if first.FPR != 0 || first.TPR != 0 || !math.IsInf(first.Threshold, 1) {
		t.Errorf("first point = %+v, want {0 0 +Inf}", first)
	}
	last := curve.Points[len(curve.Points)-1]
	if !approxEqual(last.FPR, 1) || !approxEqual(last.TPR, 1) {
		t.Errorf("last point = %+v, want FPR=1 TPR=1", last)
	}
	// The curve must be monotonic non-decreasing in both FPR and TPR.
	for i := 1; i < len(curve.Points); i++ {
		if curve.Points[i].FPR < curve.Points[i-1].FPR-epsilon ||
			curve.Points[i].TPR < curve.Points[i-1].TPR-epsilon {
			t.Fatalf("curve not monotonic at %d: %+v -> %+v",
				i, curve.Points[i-1], curve.Points[i])
		}
	}

	if _, ok := classification.ROCCurve([]int{1, 1}, []float64{0.2, 0.5}, 1); ok {
		t.Error("single class: ok = true, want false")
	}
	if _, ok := classification.ROCCurve([]int{0, 1}, []float64{0.2}, 1); ok {
		t.Error("length mismatch: ok = true, want false")
	}
	if _, ok := classification.ROCCurve([]int{0, 1}, []float64{0.2, math.NaN()}, 1); ok {
		t.Error("non-finite score: ok = true, want false")
	}
}

func TestROCCurveCollapsesTies(t *testing.T) {
	// Two samples share score 0.5: they must collapse to one operating point,
	// so the curve has the (0,0) corner plus one point per distinct score.
	yTrue := []int{0, 1, 1, 0}
	scores := []float64{0.5, 0.5, 0.6, 0.4}
	curve, ok := classification.ROCCurve(yTrue, scores, 1)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	// Distinct scores: 0.6, 0.5, 0.4 -> 3 points, plus the +Inf corner = 4.
	if len(curve.Points) != 4 {
		t.Errorf("points = %d, want 4", len(curve.Points))
	}
}

func TestLogLoss(t *testing.T) {
	tests := []struct {
		name   string
		yTrue  []int
		probs  []float64
		pos    int
		want   float64
		wantOK bool
	}{
		{
			name:   "confident and mostly right",
			yTrue:  []int{1, 0, 1, 1},
			probs:  []float64{0.9, 0.1, 0.8, 0.7},
			pos:    1,
			want:   0.19763488164214868,
			wantOK: true,
		},
		{
			name:   "probability above one rejected",
			yTrue:  []int{1, 0},
			probs:  []float64{1.2, 0.1},
			pos:    1,
			wantOK: false,
		},
		{
			name:   "NaN probability rejected",
			yTrue:  []int{1, 0},
			probs:  []float64{math.NaN(), 0.1},
			pos:    1,
			wantOK: false,
		},
		{
			name:   "length mismatch rejected",
			yTrue:  []int{1, 0},
			probs:  []float64{0.9},
			pos:    1,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := classification.LogLoss(tt.yTrue, tt.probs, tt.pos)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && !approxEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogLossClampsCertainMistake(t *testing.T) {
	// A confident, wrong prediction (p=0 for a positive) must yield a large but
	// finite loss rather than +Inf.
	got, ok := classification.LogLoss([]int{1}, []float64{0}, 1)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if math.IsInf(got, 0) || math.IsNaN(got) {
		t.Fatalf("loss = %v, want a large finite value", got)
	}
	if got < 30 {
		t.Errorf("loss = %v, want a large penalty (>30)", got)
	}
}
