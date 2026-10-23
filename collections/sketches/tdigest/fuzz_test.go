package tdigest_test

import (
	"math"
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
