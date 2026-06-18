package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

// bytesToFloats turns the fuzzer's []byte into a []float64 we can feed to the
// quantile functions. Each byte becomes a finite, non-negative value (0-255),
// so the inputs are always NaN-free and ordering invariants are easy to reason
// about.
func bytesToFloats(b []byte) []float64 {
	if b == nil {
		return nil
	}
	out := make([]float64, len(b))
	for i, v := range b {
		out[i] = float64(v)
	}
	return out
}

// FuzzQuantiles asserts the invariants that hold for any finite sample:
// quantiles stay within [min, max], are monotonic in q, Percentile agrees with
// Quantile, the quartiles are ordered with IQR = Q3-Q1 >= 0, the input is never
// mutated, out-of-range q is rejected, and a NaN poisons the result.
func FuzzQuantiles(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{42})
	f.Add([]byte{1, 2, 3, 4, 5})
	f.Add([]byte{5, 3, 1, 4, 2}) // unsorted
	f.Add([]byte{7, 7, 7})       // duplicates

	f.Fuzz(func(t *testing.T, data []byte) {
		input := bytesToFloats(data)

		// Guard against mutation by snapshotting before any call.
		snapshot := make([]float64, len(input))
		copy(snapshot, input)

		// Empty input: every entry point reports not-ok.
		if len(input) == 0 {
			if _, ok := stats.Quantile(input, 0.5); ok {
				t.Fatalf("Quantile(empty) reported ok")
			}
			if _, ok := stats.Percentile(input, 50); ok {
				t.Fatalf("Percentile(empty) reported ok")
			}
			if _, ok := stats.Quartiles(input); ok {
				t.Fatalf("Quartiles(empty) reported ok")
			}
			if _, ok := stats.IQR(input); ok {
				t.Fatalf("IQR(empty) reported ok")
			}
			return
		}

		var lo, hi float64 = input[0], input[0]
		for _, v := range input {
			if v < lo {
				lo = v
			}
			if v > hi {
				hi = v
			}
		}

		// q=0 is the minimum, q=1 the maximum.
		min, ok := stats.Quantile(input, 0)
		if !ok || min != lo {
			t.Fatalf("Quantile(_, 0) = (%v, %v), want (%v, true)", min, ok, lo)
		}
		max, ok := stats.Quantile(input, 1)
		if !ok || max != hi {
			t.Fatalf("Quantile(_, 1) = (%v, %v), want (%v, true)", max, ok, hi)
		}

		// Sweep q across [0,1]: every value is bracketed by [lo,hi], monotonic
		// non-decreasing, and Percentile agrees exactly with Quantile.
		prev := math.Inf(-1)
		for i := 0; i <= 20; i++ {
			q := float64(i) / 20
			got, ok := stats.Quantile(input, q)
			if !ok {
				t.Fatalf("Quantile(_, %v) reported not-ok", q)
			}
			if got < lo || got > hi {
				t.Fatalf("Quantile(_, %v) = %v outside [%v, %v]", q, got, lo, hi)
			}
			if got < prev {
				t.Fatalf("Quantile not monotonic at q=%v: %v < %v", q, got, prev)
			}
			prev = got

			pct, ok := stats.Percentile(input, q*100)
			if !ok || pct != got {
				t.Fatalf("Percentile(_, %v) = (%v, %v), want (%v, true)", q*100, pct, ok, got)
			}
		}

		// Out-of-range q/p are rejected.
		if _, ok := stats.Quantile(input, -0.01); ok {
			t.Fatalf("Quantile(_, -0.01) reported ok")
		}
		if _, ok := stats.Quantile(input, 1.01); ok {
			t.Fatalf("Quantile(_, 1.01) reported ok")
		}
		if _, ok := stats.Percentile(input, 100.01); ok {
			t.Fatalf("Percentile(_, 100.01) reported ok")
		}

		// Quartiles are ordered and bracketed; Q2 is the median; IQR = Q3-Q1.
		qs, ok := stats.Quartiles(input)
		if !ok {
			t.Fatalf("Quartiles reported not-ok for non-empty input")
		}
		if !(lo <= qs.Q1 && qs.Q1 <= qs.Q2 && qs.Q2 <= qs.Q3 && qs.Q3 <= hi) {
			t.Fatalf("quartiles out of order/range: lo=%v %+v hi=%v", lo, qs, hi)
		}
		median, _ := stats.Quantile(input, 0.5)
		if qs.Q2 != median {
			t.Fatalf("Quartiles.Q2 = %v, want median %v", qs.Q2, median)
		}
		iqr, ok := stats.IQR(input)
		if !ok || iqr != qs.Q3-qs.Q1 {
			t.Fatalf("IQR = (%v, %v), want (%v, true)", iqr, ok, qs.Q3-qs.Q1)
		}
		if iqr < 0 {
			t.Fatalf("IQR = %v, want >= 0", iqr)
		}

		// A NaN anywhere poisons every quantile result.
		poisoned := append([]float64{math.NaN()}, input...)
		if _, ok := stats.Quantile(poisoned, 0.5); ok {
			t.Fatalf("Quantile reported ok for NaN-contaminated input")
		}
		if _, ok := stats.IQR(poisoned); ok {
			t.Fatalf("IQR reported ok for NaN-contaminated input")
		}

		// The caller's slice must be untouched throughout.
		for i := range snapshot {
			if input[i] != snapshot[i] {
				t.Fatalf("input mutated at %d: %v != %v", i, input[i], snapshot[i])
			}
		}
	})
}
