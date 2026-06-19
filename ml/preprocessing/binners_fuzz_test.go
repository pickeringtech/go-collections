package preprocessing_test

import (
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

// FuzzFixedWidthBinner asserts that every bin index is within [0, nBins) and
// that binning sorted input yields non-decreasing indices (bins are monotonic
// in value).
func FuzzFixedWidthBinner(f *testing.F) {
	f.Add([]byte(nil), uint8(4))
	f.Add([]byte{5}, uint8(3))
	f.Add([]byte{1, 9, 3, 7}, uint8(5))
	f.Add([]byte{4, 4, 4}, uint8(2))

	f.Fuzz(func(t *testing.T, data []byte, rawBins uint8) {
		nBins := int(rawBins%16) + 1 // 1..16
		values := bytesToFloats(data)

		binner := preprocessing.NewFixedWidthBinner(nBins).Fit(values)
		got, ok := binner.Transform(values)
		if len(values) == 0 {
			if ok {
				t.Fatalf("empty Fit reported fitted")
			}
			return
		}
		if !ok {
			t.Fatalf("Transform not-ok for fitted binner")
		}
		assertBinInvariants(t, values, got, nBins)
	})
}

// FuzzQuantileBinner asserts the same invariants for quantile binning.
func FuzzQuantileBinner(f *testing.F) {
	f.Add([]byte(nil), uint8(4))
	f.Add([]byte{5}, uint8(3))
	f.Add([]byte{1, 9, 3, 7}, uint8(5))
	f.Add([]byte{4, 4, 4}, uint8(2))

	f.Fuzz(func(t *testing.T, data []byte, rawBins uint8) {
		nBins := int(rawBins%16) + 1
		values := bytesToFloats(data)

		binner := preprocessing.NewQuantileBinner(nBins).Fit(values)
		got, ok := binner.Transform(values)
		if len(values) == 0 {
			if ok {
				t.Fatalf("empty Fit reported fitted")
			}
			return
		}
		if !ok {
			t.Fatalf("Transform not-ok for fitted binner")
		}
		assertBinInvariants(t, values, got, nBins)
	})
}

// assertBinInvariants checks bin indices are in range and monotonic in value.
func assertBinInvariants(t *testing.T, values []float64, bins []int, nBins int) {
	t.Helper()
	if len(bins) != len(values) {
		t.Fatalf("len(bins) = %d, want %d", len(bins), len(values))
	}
	for i, b := range bins {
		if b < 0 || b >= nBins {
			t.Fatalf("bin[%d] = %d outside [0, %d)", i, b, nBins)
		}
	}

	// Pair values with their bins, sort by value, and check bins never decrease.
	type pair struct {
		v float64
		b int
	}
	pairs := make([]pair, len(values))
	for i := range values {
		pairs[i] = pair{values[i], bins[i]}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].v < pairs[j].v })
	for i := 1; i < len(pairs); i++ {
		if pairs[i].b < pairs[i-1].b {
			t.Fatalf("non-monotonic bins: value %v -> %d after value %v -> %d",
				pairs[i].v, pairs[i].b, pairs[i-1].v, pairs[i-1].b)
		}
	}
}
