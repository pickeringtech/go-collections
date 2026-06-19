package classification_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/ml/metrics/classification"
)

const epsilon = 1e-9

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}

// The canonical 3-class fixture used across the averaging tests.
var (
	yTrue3 = []int{0, 1, 2, 2, 1, 0, 1, 2}
	yPred3 = []int{0, 2, 1, 2, 1, 0, 1, 2}
)

func TestAccuracy(t *testing.T) {
	tests := []struct {
		name   string
		yTrue  []int
		yPred  []int
		want   float64
		wantOK bool
	}{
		{"three of four wrong-free", []int{1, 2, 3, 4}, []int{1, 2, 3, 9}, 0.75, true},
		{"perfect", []int{1, 1, 1}, []int{1, 1, 1}, 1, true},
		{"none correct", []int{1, 2}, []int{2, 1}, 0, true},
		{"empty undefined", []int{}, []int{}, 0, false},
		{"length mismatch undefined", []int{1, 2}, []int{1}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := classification.Accuracy(tt.yTrue, tt.yPred)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && !approxEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccuracyOnStrings(t *testing.T) {
	got, ok := classification.Accuracy(
		[]string{"cat", "dog", "cat"},
		[]string{"cat", "cat", "cat"},
	)
	if !ok || !approxEqual(got, 2.0/3) {
		t.Fatalf("got %v %v, want 0.6667 true", got, ok)
	}
}

func TestConfusionMatrix(t *testing.T) {
	cm, ok := classification.ConfusionMatrix(yTrue3, yPred3)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if got := len(cm.Labels()); got != 3 {
		t.Fatalf("labels = %d, want 3", got)
	}
	// Diagonal: correctly classified counts.
	if got := cm.Count(0, 0); got != 2 {
		t.Errorf("Count(0,0) = %d, want 2", got)
	}
	if got := cm.Count(2, 2); got != 2 {
		t.Errorf("Count(2,2) = %d, want 2", got)
	}
	// Off-diagonal confusions: true 1 predicted 2, and true 2 predicted 1.
	if got := cm.Count(1, 2); got != 1 {
		t.Errorf("Count(1,2) = %d, want 1", got)
	}
	if got := cm.Count(2, 1); got != 1 {
		t.Errorf("Count(2,1) = %d, want 1", got)
	}
	// Unknown labels count as zero, never panic.
	if got := cm.Count(9, 9); got != 0 {
		t.Errorf("Count(9,9) = %d, want 0", got)
	}

	if _, ok := classification.ConfusionMatrix([]int{}, []int{}); ok {
		t.Error("empty: ok = true, want false")
	}
}

func TestPrecisionRecallF1Averaging(t *testing.T) {
	tests := []struct {
		name string
		fn   func([]int, []int, classification.Averaging) (float64, bool)
		avg  classification.Averaging
		want float64
	}{
		{"macro precision", classification.Precision[int], classification.Macro, 0.7777777777777777},
		{"macro recall", classification.Recall[int], classification.Macro, 0.7777777777777777},
		{"macro f1", classification.F1[int], classification.Macro, 0.7777777777777777},
		{"micro f1 equals accuracy", classification.F1[int], classification.Micro, 0.75},
		{"micro precision equals accuracy", classification.Precision[int], classification.Micro, 0.75},
		{"weighted f1", classification.F1[int], classification.Weighted, 0.75},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.fn(yTrue3, yPred3, tt.avg)
			if !ok {
				t.Fatal("ok = false, want true")
			}
			if !approxEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrecisionRejectsBadInput(t *testing.T) {
	if _, ok := classification.F1([]int{}, []int{}, classification.Macro); ok {
		t.Error("empty: ok = true, want false")
	}
	if _, ok := classification.F1([]int{1, 2}, []int{1}, classification.Macro); ok {
		t.Error("length mismatch: ok = true, want false")
	}
	// Unrecognised averaging strategy is rejected.
	if _, ok := classification.F1(yTrue3, yPred3, classification.Averaging(99)); ok {
		t.Error("bad averaging: ok = true, want false")
	}
}

func TestAveragingString(t *testing.T) {
	cases := map[classification.Averaging]string{
		classification.Macro:        "macro",
		classification.Micro:        "micro",
		classification.Weighted:     "weighted",
		classification.Averaging(7): "unknown",
	}
	for avg, want := range cases {
		if got := avg.String(); got != want {
			t.Errorf("Averaging(%d).String() = %q, want %q", avg, got, want)
		}
	}
}

func TestPrecisionRecallRejectEmpty(t *testing.T) {
	// Exercise the early-out on each top-level entry point, not just F1.
	if _, ok := classification.Precision([]int{}, []int{}, classification.Macro); ok {
		t.Error("Precision empty: ok = true, want false")
	}
	if _, ok := classification.Recall([]int{1}, []int{1, 2}, classification.Macro); ok {
		t.Error("Recall mismatch: ok = true, want false")
	}
}

func TestPerClassDegenerateContributesZero(t *testing.T) {
	// Label 2 is predicted once but never true, so its recall and F1 are 0
	// (zero denominators), exercising the degenerate per-class branches.
	yTrue := []int{0, 0, 1}
	yPred := []int{0, 2, 1}

	macroR, ok := classification.Recall(yTrue, yPred, classification.Macro)
	if !ok || !approxEqual(macroR, 0.5) { // (0.5 + 1 + 0) / 3
		t.Errorf("macro recall = %v %v, want 0.5 true", macroR, ok)
	}
	macroF1, ok := classification.F1(yTrue, yPred, classification.Macro)
	if !ok || !approxEqual(macroF1, 5.0/9) { // (2/3 + 1 + 0) / 3
		t.Errorf("macro F1 = %v %v, want %v true", macroF1, ok, 5.0/9)
	}
}

func TestBinaryMetrics(t *testing.T) {
	// positive = 1. tp=2 (idx 2,3-ish), build explicitly:
	yTrue := []int{1, 0, 1, 1, 0, 0}
	yPred := []int{1, 0, 0, 1, 1, 0}
	// positive class 1: tp = predicted 1 & true 1 -> idx0,idx3 = 2
	//                   fp = predicted 1 & true 0 -> idx4 = 1
	//                   fn = predicted !=1 & true 1 -> idx2 = 1
	p, ok := classification.PrecisionBinary(yTrue, yPred, 1)
	if !ok || !approxEqual(p, 2.0/3) {
		t.Errorf("PrecisionBinary = %v %v, want 0.6667 true", p, ok)
	}
	r, ok := classification.RecallBinary(yTrue, yPred, 1)
	if !ok || !approxEqual(r, 2.0/3) {
		t.Errorf("RecallBinary = %v %v, want 0.6667 true", r, ok)
	}
	f, ok := classification.F1Binary(yTrue, yPred, 1)
	if !ok || !approxEqual(f, 2.0/3) {
		t.Errorf("F1Binary = %v %v, want 0.6667 true", f, ok)
	}
}

func TestBinaryDegenerate(t *testing.T) {
	// No positive predictions -> precision defined as 0, ok true.
	p, ok := classification.PrecisionBinary([]int{1, 1}, []int{0, 0}, 1)
	if !ok || p != 0 {
		t.Errorf("PrecisionBinary = %v %v, want 0 true", p, ok)
	}
	if _, ok := classification.F1Binary([]int{}, []int{}, 1); ok {
		t.Error("empty: ok = true, want false")
	}
	if _, ok := classification.PrecisionBinary([]int{1}, []int{1, 0}, 1); ok {
		t.Error("PrecisionBinary mismatch: ok = true, want false")
	}
	if _, ok := classification.RecallBinary([]int{}, []int{}, 1); ok {
		t.Error("RecallBinary empty: ok = true, want false")
	}
}
