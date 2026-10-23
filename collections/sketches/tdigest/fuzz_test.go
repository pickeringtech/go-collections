package tdigest_test

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/tdigest"
	"github.com/pickeringtech/go-collections/stats"
)

// FuzzDigest_BoundedError checks that for any newline-separated list of numbers
// the t-digest's quantile estimates stay within a generous band of the exact
// quantiles computed by stats.Quantile (the oracle). The band is wide because
// fuzzing produces small inputs where the t-digest's relative error is larger;
// the point is to catch gross misbehaviour (a panic, a wildly wrong estimate,
// or an estimate outside [min,max]), not to assert the analytic bound on tiny
// inputs.
func FuzzDigest_BoundedError(f *testing.F) {
	f.Add("1\n2\n3")
	f.Add("")
	f.Add("5\n5\n5\n5")
	f.Add("-10\n0\n10\n1000\n-1000")
	f.Add("0.5\n1.5\n2.5\n3.5\n4.5\n5.5")

	f.Fuzz(func(t *testing.T, blob string) {
		var values []float64
		for _, tok := range strings.Split(blob, "\n") {
			tok = strings.TrimSpace(tok)
			if tok == "" {
				continue
			}
			v, err := strconv.ParseFloat(tok, 64)
			if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
				continue
			}
			values = append(values, v)
		}

		d, err := tdigest.New()
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		for _, v := range values {
			d.Add(v)
		}

		if len(values) == 0 {
			if _, ok := d.Quantile(0.5); ok {
				t.Fatal("Quantile on empty digest returned ok=true")
			}
			return
		}

		lo, hi, _ := stats.MinMax(values)
		rng := hi - lo

		for _, q := range []float64{0, 0.25, 0.5, 0.75, 1} {
			got, ok := d.Quantile(q)
			if !ok {
				t.Fatalf("Quantile(%v) ok=false for %d values", q, len(values))
			}
			// Every estimate must lie within the observed [min,max].
			if got < lo-1e-6 || got > hi+1e-6 {
				t.Fatalf("Quantile(%v) = %v outside [%v,%v]", q, got, lo, hi)
			}
			exact, _ := stats.Quantile(values, q)
			// Very generous band: on the tiny inputs fuzzing produces, the
			// t-digest snaps quantiles to a centroid mean rather than
			// interpolating like the exact oracle, so the gap can be a large
			// fraction of the (small) data range. Half the range plus an
			// absolute slack catches gross misbehaviour without flagging this
			// expected small-n coarseness.
			tol := 0.5*rng + 1e-6
			if math.Abs(got-exact) > tol {
				t.Fatalf("Quantile(%v) = %v, exact = %v, |diff| %v > tol %v",
					q, got, exact, math.Abs(got-exact), tol)
			}
		}
	})
}

// FuzzDigest_CDF checks CDF against the exact empirical CDF (the fraction of
// values at or below x). Below a couple of hundred unit-weight values the
// default compression never merges two values into one centroid, so every
// centroid is a point mass and CDF must equal the empirical CDF exactly — at
// every observed value (ties included) and between them. Larger inputs must
// still give a non-decreasing CDF in [0, 1].
func FuzzDigest_CDF(f *testing.F) {
	f.Add([]byte{1, 2, 3})
	f.Add([]byte{})
	f.Add([]byte{5, 5, 5, 5})
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
	f.Add([]byte{1, 5, 5, 5, 5, 9})

	f.Fuzz(func(t *testing.T, raw []byte) {
		values := make([]float64, len(raw))
		for i, b := range raw {
			values[i] = float64(b % 32) // a small domain so ties are common
		}
		d, err := tdigest.New()
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		for _, v := range values {
			d.Add(v)
		}
		if len(values) == 0 {
			_, ok := d.CDF(0)
			if ok {
				t.Fatal("CDF on empty digest returned ok=true")
			}
			return
		}
		sorted := append([]float64(nil), values...)
		sort.Float64s(sorted)
		prev := 0.0
		for x := -1.0; x <= 33; x += 0.5 {
			got, ok := d.CDF(x)
			if !ok || got < prev || got > 1 {
				t.Fatalf("CDF(%v) = (%v, %v) after %v; want non-decreasing in [0, 1]", x, got, ok, prev)
			}
			prev = got
			if len(values) < 200 {
				want := float64(sort.Search(len(sorted), func(i int) bool { return sorted[i] > x })) / float64(len(sorted))
				if got != want {
					t.Fatalf("CDF(%v) = %v over %d values, want exact %v", x, got, len(values), want)
				}
			}
		}
	})
}
