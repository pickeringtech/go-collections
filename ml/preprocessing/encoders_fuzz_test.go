package preprocessing_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

// bytesToInts turns the fuzzer's []byte into []int categories.
func bytesToInts(b []byte) []int {
	out := make([]int, len(b))
	for i, v := range b {
		out[i] = int(v)
	}
	return out
}

// FuzzOneHotEncoder asserts that every row of the one-hot encoding of the
// training data has exactly one 1 (each training value is a known category),
// the column count equals the number of distinct categories, and the hot column
// matches the label encoder's code for that value.
func FuzzOneHotEncoder(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{1})
	f.Add([]byte{1, 2, 2, 3})
	f.Add([]byte{9, 9, 9})

	f.Fuzz(func(t *testing.T, data []byte) {
		cats := bytesToInts(data)

		oneHot := preprocessing.NewOneHotEncoder[int]().Fit(cats)
		labels := preprocessing.NewLabelEncoder[int]().Fit(cats)

		rows, ok := oneHot.Transform(cats)
		codes, okLabels := labels.Transform(cats)
		if len(cats) == 0 {
			if ok || okLabels {
				t.Fatalf("empty Fit reported fitted")
			}
			return
		}
		if !ok || !okLabels {
			t.Fatalf("Transform not-ok for fitted encoder")
		}

		nCats := len(oneHot.Categories())
		for i, row := range rows {
			if len(row) != nCats {
				t.Fatalf("row %d has %d columns, want %d", i, len(row), nCats)
			}
			sum := 0.0
			hot := -1
			for c, v := range row {
				sum += v
				if v == 1 {
					hot = c
				}
			}
			if sum != 1 {
				t.Fatalf("row %d sums to %v, want 1", i, sum)
			}
			if hot != codes[i] {
				t.Fatalf("row %d hot column %d != label code %d", i, hot, codes[i])
			}
		}
	})
}

// FuzzLabelEncoder asserts that encoding then inverse-encoding the training data
// reproduces it exactly (a round-trip invariant).
func FuzzLabelEncoder(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{1})
	f.Add([]byte{3, 1, 2, 1})
	f.Add([]byte{5, 5, 5})

	f.Fuzz(func(t *testing.T, data []byte) {
		cats := bytesToInts(data)
		enc := preprocessing.NewLabelEncoder[int]().Fit(cats)

		codes, ok := enc.Transform(cats)
		if len(cats) == 0 {
			if ok {
				t.Fatalf("empty Fit reported fitted")
			}
			return
		}
		if !ok {
			t.Fatalf("Transform not-ok for fitted encoder")
		}

		back, ok := enc.InverseTransform(codes)
		if !ok {
			t.Fatalf("InverseTransform not-ok for in-range codes")
		}
		if len(back) != len(cats) {
			t.Fatalf("round-trip length %d, want %d", len(back), len(cats))
		}
		for i := range cats {
			if back[i] != cats[i] {
				t.Fatalf("round-trip[%d] = %v, want %v", i, back[i], cats[i])
			}
		}
	})
}
