package preprocessing_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

// FuzzMeanImputer asserts that imputing data which contains no missing values is
// the identity, the output length always matches the input, and Transform never
// mutates its input. Here -1 is the sentinel for missing.
func FuzzMeanImputer(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{5})
	f.Add([]byte{1, 2, 3, 4})

	isMissing := func(v float64) bool { return v == -1 }

	f.Fuzz(func(t *testing.T, data []byte) {
		// bytesToFloats yields values in [0, 255], so none equals the -1
		// sentinel: the data has no missing entries.
		values := bytesToFloats(data)
		imp := preprocessing.NewMeanImputer(isMissing).Fit(values)

		if len(values) == 0 {
			if _, ok := imp.Transform([]float64{1}); ok {
				t.Fatalf("empty Fit reported fitted")
			}
			return
		}

		snapshot := make([]float64, len(values))
		copy(snapshot, values)

		got, ok := imp.Transform(values)
		if !ok {
			t.Fatalf("Transform reported not-ok for fitted imputer")
		}
		if len(got) != len(values) {
			t.Fatalf("len(out) = %d, want %d", len(got), len(values))
		}
		// No value is missing, so imputation is the identity.
		if !floatSlicesClose(got, values) {
			t.Fatalf("identity imputation changed data: %v != %v", got, values)
		}
		if !floatSlicesClose(values, snapshot) {
			t.Fatalf("Transform mutated input: %v != %v", values, snapshot)
		}
	})
}
