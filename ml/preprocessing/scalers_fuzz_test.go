package preprocessing_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

// FuzzStandardScaler asserts the no-leakage invariant: parameters are captured
// at Fit time and reused verbatim, so (a) Transform is a pure function of the
// stored mean/stddev, (b) transforming the same data twice is identical, and
// (c) Fit never mutates the training slice.
func FuzzStandardScaler(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{42})
	f.Add([]byte{1, 2, 3, 4, 5})
	f.Add([]byte{7, 7, 7})

	f.Fuzz(func(t *testing.T, data []byte) {
		train := bytesToFloats(data)

		snapshot := make([]float64, len(train))
		copy(snapshot, train)

		scaler := preprocessing.NewStandardScaler().Fit(train)

		// Empty training data leaves the scaler unfitted.
		if len(train) == 0 {
			if _, ok := scaler.Transform([]float64{1}); ok {
				t.Fatalf("unfitted scaler reported ok")
			}
			return
		}

		mean := scaler.Mean()
		std := scaler.StdDev()

		got, ok := scaler.Transform(train)
		if !ok {
			t.Fatalf("Transform reported not-ok for fitted scaler")
		}
		for i, v := range train {
			want := 0.0
			if std != 0 {
				want = (v - mean) / std
			}
			if !floatsClose(got[i], want) {
				t.Fatalf("Transform[%d] = %v, want %v (mean=%v std=%v)", i, got[i], want, mean, std)
			}
		}

		// Re-transforming yields identical output: parameters are frozen.
		again, _ := scaler.Transform(train)
		if !floatSlicesClose(got, again) {
			t.Fatalf("Transform not stable across calls: %v vs %v", got, again)
		}

		// Training slice is untouched.
		if !floatSlicesClose(train, snapshot) {
			t.Fatalf("Fit mutated training data: %v != %v", train, snapshot)
		}
	})
}

// FuzzMinMaxScaler asserts that fitted min-max output stays in [0, 1] for the
// training data itself and that Fit does not mutate its input.
func FuzzMinMaxScaler(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{5})
	f.Add([]byte{1, 9, 3, 7})
	f.Add([]byte{4, 4, 4})

	f.Fuzz(func(t *testing.T, data []byte) {
		train := bytesToFloats(data)
		snapshot := make([]float64, len(train))
		copy(snapshot, train)

		got, ok := preprocessing.NewMinMaxScaler().FitTransform(train)
		if len(train) == 0 {
			if ok {
				t.Fatalf("empty FitTransform reported ok")
			}
			return
		}
		if !ok {
			t.Fatalf("FitTransform reported not-ok for non-empty data")
		}
		for i, v := range got {
			if v < 0 || v > 1 {
				t.Fatalf("scaled[%d] = %v outside [0, 1]", i, v)
			}
		}
		if !floatSlicesClose(train, snapshot) {
			t.Fatalf("FitTransform mutated input: %v != %v", train, snapshot)
		}
	})
}
