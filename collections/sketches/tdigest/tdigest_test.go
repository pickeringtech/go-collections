package tdigest_test

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/tdigest"
	"github.com/pickeringtech/go-collections/stats"
)

// buildDigest feeds values into a fresh Digest at the given compression.
func buildDigest(t *testing.T, compression float64, values []float64) *tdigest.Digest {
	t.Helper()
	d, err := tdigest.New(tdigest.WithCompression(compression))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, v := range values {
		d.Add(v)
	}
	return d
}

// dataRange returns max-min of values (0 for empty).
func dataRange(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	lo, hi := values[0], values[0]
	for _, v := range values {
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	return hi - lo
}

// tolerance returns the absolute error allowed for the q-quantile of data with
// the given range. The band is a fraction of the data range, tighter at the
// tails (q<=0.05 or q>=0.95) where the t-digest is most accurate and more
// generous in the middle.
func tolerance(rng, q float64) float64 {
	frac := 0.06
	if q <= 0.05 || q >= 0.95 {
		frac = 0.02
	}
	return frac*rng + 1e-9
}

func TestNew_DefaultCompression(t *testing.T) {
	d, err := tdigest.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if d.Compression() != tdigest.DefaultCompression {
		t.Errorf("Compression() = %v, want %v", d.Compression(), tdigest.DefaultCompression)
	}
}

func TestWithCompression_Validation(t *testing.T) {
	type args struct {
		compression float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "positive compression is valid", args: args{compression: 50}, wantErr: false},
		{name: "small positive compression is valid", args: args{compression: 1}, wantErr: false},
		{name: "zero compression is rejected", args: args{compression: 0}, wantErr: true},
		{name: "negative compression is rejected", args: args{compression: -10}, wantErr: true},
		{name: "NaN compression is rejected", args: args{compression: math.NaN()}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tdigest.New(tdigest.WithCompression(tt.args.compression))
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Fatalf("New error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, tdigest.ErrInvalidConfig) {
				t.Errorf("error = %v, want wrapping ErrInvalidConfig", err)
			}
		})
	}
}

func TestDigest_QuantileEmpty(t *testing.T) {
	d := buildDigest(t, 100, nil)
	if got, ok := d.Quantile(0.5); ok {
		t.Errorf("Quantile on empty = (%v, %v), want (_, false)", got, ok)
	}
	if got, ok := d.Percentile(50); ok {
		t.Errorf("Percentile on empty = (%v, %v), want (_, false)", got, ok)
	}
	if got, ok := d.CDF(0); ok {
		t.Errorf("CDF on empty = (%v, %v), want (_, false)", got, ok)
	}
	if got, ok := d.Min(); ok {
		t.Errorf("Min on empty = (%v, %v), want (_, false)", got, ok)
	}
	if got, ok := d.Max(); ok {
		t.Errorf("Max on empty = (%v, %v), want (_, false)", got, ok)
	}
	if d.Count() != 0 {
		t.Errorf("Count on empty = %v, want 0", d.Count())
	}
}

func TestDigest_QuantileRangeRejected(t *testing.T) {
	d := buildDigest(t, 100, []float64{1, 2, 3})
	type args struct {
		q float64
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "q below 0 rejected", args: args{q: -0.1}},
		{name: "q above 1 rejected", args: args{q: 1.1}},
		{name: "q NaN rejected", args: args{q: math.NaN()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := d.Quantile(tt.args.q); ok {
				t.Errorf("Quantile(%v) = (%v, %v), want (_, false)", tt.args.q, got, ok)
			}
		})
	}
}

func TestPercentile_RangeRejected(t *testing.T) {
	d := buildDigest(t, 100, []float64{1, 2, 3})
	for _, p := range []float64{-1, 101, math.NaN()} {
		if got, ok := d.Percentile(p); ok {
			t.Errorf("Percentile(%v) = (%v, %v), want (_, false)", p, got, ok)
		}
	}
}

func TestDigest_SingleElement(t *testing.T) {
	d := buildDigest(t, 100, []float64{42})
	for _, q := range []float64{0, 0.25, 0.5, 0.75, 1} {
		got, ok := d.Quantile(q)
		if !ok || got != 42 {
			t.Errorf("Quantile(%v) = (%v, %v), want (42, true)", q, got, ok)
		}
	}
}

func TestDigest_AllEqual(t *testing.T) {
	values := make([]float64, 1000)
	for i := range values {
		values[i] = 7
	}
	d := buildDigest(t, 100, values)
	for _, q := range []float64{0, 0.5, 1} {
		got, ok := d.Quantile(q)
		if !ok || got != 7 {
			t.Errorf("Quantile(%v) = (%v, %v), want (7, true)", q, got, ok)
		}
	}
}

func TestDigest_EdgeQuantilesAreExactExtremes(t *testing.T) {
	values := []float64{5, 1, 9, 3, 7, 2, 8, 4, 6}
	d := buildDigest(t, 100, values)
	if got, ok := d.Quantile(0); !ok || got != 1 {
		t.Errorf("Quantile(0) = (%v, %v), want (1, true)", got, ok)
	}
	if got, ok := d.Quantile(1); !ok || got != 9 {
		t.Errorf("Quantile(1) = (%v, %v), want (9, true)", got, ok)
	}
	if minV, ok := d.Min(); !ok || minV != 1 {
		t.Errorf("Min() = (%v, %v), want (1, true)", minV, ok)
	}
	if maxV, ok := d.Max(); !ok || maxV != 9 {
		t.Errorf("Max() = (%v, %v), want (9, true)", maxV, ok)
	}
}

func TestDigest_NonFiniteIgnored(t *testing.T) {
	d := buildDigest(t, 100, []float64{1, 2, 3})
	before := d.Count()
	d.Add(math.NaN())
	d.Add(math.Inf(1))
	d.Add(math.Inf(-1))
	d.AddWeighted(5, math.NaN())
	d.AddWeighted(5, 0)
	d.AddWeighted(5, -1)
	d.AddWeighted(math.NaN(), 2)
	if d.Count() != before {
		t.Errorf("Count after non-finite adds = %v, want %v (all ignored)", d.Count(), before)
	}
}

// TestDigest_QuantileOracle compares the digest's quantile estimates against the
// exact stats.Quantile oracle across several distributions and quantiles.
func TestDigest_QuantileOracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	uniform := make([]float64, 20000)
	for i := range uniform {
		uniform[i] = rng.Float64() * 1000
	}
	normalish := make([]float64, 20000)
	for i := range normalish {
		// Sum of uniforms approximates a bell curve.
		s := 0.0
		for j := 0; j < 12; j++ {
			s += rng.Float64()
		}
		normalish[i] = (s - 6) * 50
	}
	skewed := make([]float64, 20000)
	for i := range skewed {
		skewed[i] = math.Exp(rng.Float64() * 6) // heavy right tail
	}

	type args struct {
		data []float64
		q    float64
	}
	quantiles := []float64{0.01, 0.05, 0.25, 0.5, 0.75, 0.95, 0.99}
	datasets := map[string][]float64{
		"uniform":   uniform,
		"normalish": normalish,
		"skewed":    skewed,
	}

	var tests []struct {
		name string
		args args
	}
	for dname, data := range datasets {
		for _, q := range quantiles {
			tests = append(tests, struct {
				name string
				args args
			}{
				name: dname,
				args: args{data: data, q: q},
			})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := buildDigest(t, 200, tt.args.data)
			exact, ok := stats.Quantile(tt.args.data, tt.args.q)
			if !ok {
				t.Fatalf("stats.Quantile oracle failed for q=%v", tt.args.q)
			}
			got, gotOK := d.Quantile(tt.args.q)
			if !gotOK {
				t.Fatalf("Quantile(%v) ok=false", tt.args.q)
			}
			tol := tolerance(dataRange(tt.args.data), tt.args.q)
			if math.Abs(got-exact) > tol {
				t.Errorf("Quantile(%v) = %v, exact = %v, |diff| = %v > tol %v",
					tt.args.q, got, exact, math.Abs(got-exact), tol)
			}
		})
	}
}

func TestPercentile_MatchesQuantile(t *testing.T) {
	values := make([]float64, 5000)
	for i := range values {
		values[i] = float64(i)
	}
	d := buildDigest(t, 100, values)
	for _, p := range []float64{0, 25, 50, 75, 100} {
		pv, ok := d.Percentile(p)
		if !ok {
			t.Fatalf("Percentile(%v) ok=false", p)
		}
		qv, _ := d.Quantile(p / 100)
		if pv != qv {
			t.Errorf("Percentile(%v) = %v, Quantile(%v) = %v, want equal", p, pv, p/100, qv)
		}
	}
}

func TestDigest_CDFRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewPCG(7, 11))
	values := make([]float64, 20000)
	for i := range values {
		values[i] = rng.Float64() * 100
	}
	d := buildDigest(t, 200, values)

	// CDF below min is 0, at/above max is 1.
	if c, ok := d.CDF(-5); !ok || c != 0 {
		t.Errorf("CDF below min = (%v, %v), want (0, true)", c, ok)
	}
	maxV, _ := d.Max()
	if c, ok := d.CDF(maxV); !ok || c != 1 {
		t.Errorf("CDF(max) = (%v, %v), want (1, true)", c, ok)
	}
	if c, ok := d.CDF(math.NaN()); ok {
		t.Errorf("CDF(NaN) = (%v, %v), want (_, false)", c, ok)
	}

	// CDF(Quantile(q)) ≈ q.
	for _, q := range []float64{0.1, 0.3, 0.5, 0.7, 0.9} {
		v, _ := d.Quantile(q)
		c, ok := d.CDF(v)
		if !ok {
			t.Fatalf("CDF(%v) ok=false", v)
		}
		if math.Abs(c-q) > 0.03 {
			t.Errorf("CDF(Quantile(%v)) = %v, want ≈ %v", q, c, q)
		}
	}
}

// TestDigest_QuantileInterpolationBranches drives the extreme-but-not-edge
// quantiles so both the "before the first centroid centre" and "after the last
// centroid centre" interpolation paths execute on a multi-centroid digest.
func TestDigest_QuantileInterpolationBranches(t *testing.T) {
	values := make([]float64, 5000)
	for i := range values {
		values[i] = float64(i)
	}
	d := buildDigest(t, 100, values)

	// A target below the first centroid's centre exercises the i==0 path; a
	// target beyond the last centroid's centre exercises the post-loop path.
	lo, _ := d.Quantile(1e-6)
	if lo < 0 || lo > 50 {
		t.Errorf("Quantile(1e-6) = %v, want near the minimum (0)", lo)
	}
	hi, _ := d.Quantile(1 - 1e-6)
	if hi < 4950 || hi > 4999 {
		t.Errorf("Quantile(1-1e-6) = %v, want near the maximum (4999)", hi)
	}
}

// TestDigest_CDFInterpolationBranches drives the CDF's first-centroid,
// mid-centroid and last-centroid (post-loop) interpolation paths. The min/max
// are placed far below/above the body so a query just inside each extreme falls
// before the first centroid mean and after the last centroid mean respectively.
func TestDigest_CDFInterpolationBranches(t *testing.T) {
	var values []float64
	values = append(values, -1000) // a lone minimum well below the body
	for i := 0; i < 5000; i++ {
		values = append(values, float64(i))
	}
	values = append(values, 1e6) // a lone maximum well above the body
	d := buildDigest(t, 100, values)

	// A value between the minimum and the first body centroid exercises i==0.
	if c, ok := d.CDF(-500); !ok || c <= 0 || c >= 0.1 {
		t.Errorf("CDF(-500) = (%v, %v), want a small positive fraction", c, ok)
	}
	// A mid value.
	if c, ok := d.CDF(2500); !ok || c <= 0.3 || c >= 0.7 {
		t.Errorf("CDF(2500) = (%v, %v), want ≈ 0.5", c, ok)
	}
	// A value between the last body centroid and the maximum exercises the
	// post-loop last-centroid path.
	if c, ok := d.CDF(500000); !ok || c <= 0.9 || c >= 1 {
		t.Errorf("CDF(500000) = (%v, %v), want a large fraction below 1", c, ok)
	}
}

func TestDigest_CDFSingleElement(t *testing.T) {
	d := buildDigest(t, 100, []float64{42})
	// All weight at 42: below it 0, at/above it 1.
	if c, ok := d.CDF(41); !ok || c != 0 {
		t.Errorf("CDF(41) = (%v, %v), want (0, true)", c, ok)
	}
	if c, ok := d.CDF(42); !ok || c != 1 {
		t.Errorf("CDF(42) = (%v, %v), want (1, true)", c, ok)
	}
}

// TestDigest_ExtremeWeightsNoOverflow guards the interpolation against
// intermediate-overflow: with centroid weights and values near the float64
// limit, the naive (target-r0)*(v1-v0)/(r1-r0) form overflows to ±Inf before
// the divide. Quantile and CDF must still return finite, in-range estimates.
func TestDigest_ExtremeWeightsNoOverflow(t *testing.T) {
	d := buildDigest(t, 100, nil)
	d.AddWeighted(0, 1e300)
	for i := 0; i < 200; i++ {
		d.AddWeighted(float64(i+1), 1)
	}
	d.AddWeighted(1e9, 1e308)

	minV, _ := d.Min()
	maxV, _ := d.Max()
	for q := 0.0; q <= 1.0; q += 0.001 {
		got, ok := d.Quantile(q)
		if !ok {
			t.Fatalf("Quantile(%v) ok=false", q)
		}
		if math.IsNaN(got) || math.IsInf(got, 0) {
			t.Fatalf("Quantile(%v) = %v, want finite", q, got)
		}
		if got < minV || got > maxV {
			t.Fatalf("Quantile(%v) = %v outside [%v,%v]", q, got, minV, maxV)
		}
	}
	for _, x := range []float64{1, 100, 1e3, 1e6, 5e8, 1e9 - 1} {
		c, ok := d.CDF(x)
		if !ok {
			t.Fatalf("CDF(%v) ok=false", x)
		}
		if math.IsNaN(c) || math.IsInf(c, 0) || c < 0 || c > 1 {
			t.Fatalf("CDF(%v) = %v, want a finite fraction in [0,1]", x, c)
		}
	}
}

func TestDigest_Clear(t *testing.T) {
	d := buildDigest(t, 100, []float64{1, 2, 3, 4, 5})
	d.Clear()
	if d.Count() != 0 {
		t.Errorf("Count after Clear = %v, want 0", d.Count())
	}
	if _, ok := d.Quantile(0.5); ok {
		t.Error("Quantile after Clear should be ok=false")
	}
	if d.Compression() != 100 {
		t.Errorf("Compression after Clear = %v, want 100", d.Compression())
	}
	// Reusable after Clear.
	d.Add(99)
	if got, ok := d.Quantile(0.5); !ok || got != 99 {
		t.Errorf("Quantile after Clear+Add = (%v, %v), want (99, true)", got, ok)
	}
}

func TestDigest_AddWeightedEquivalence(t *testing.T) {
	// AddWeighted(x, k) should behave like Add(x) repeated k times: the total
	// weight is identical and the quantile estimates agree within tolerance.
	// (They are not bit-identical because the unit-weight stream produces more
	// centroids to merge, so the asserted band is statistical, not exact.)
	var expanded []float64
	weighted := buildDigest(t, 100, nil)
	for v := 0; v < 100; v++ {
		w := float64(v%5 + 1)
		weighted.AddWeighted(float64(v), w)
		for k := 0; k < int(w); k++ {
			expanded = append(expanded, float64(v))
		}
	}
	repeated := buildDigest(t, 100, expanded)

	if weighted.Count() != repeated.Count() {
		t.Errorf("weighted Count = %v, repeated Count = %v", weighted.Count(), repeated.Count())
	}
	rng := dataRange(expanded)
	for _, q := range []float64{0.1, 0.5, 0.9} {
		wq, _ := weighted.Quantile(q)
		rq, _ := repeated.Quantile(q)
		if math.Abs(wq-rq) > tolerance(rng, q) {
			t.Errorf("Quantile(%v): weighted %v vs repeated %v exceeds tolerance", q, wq, rq)
		}
	}
}

func TestDigest_Merge(t *testing.T) {
	rng := rand.New(rand.NewPCG(3, 5))
	a := make([]float64, 10000)
	b := make([]float64, 10000)
	for i := range a {
		a[i] = rng.Float64() * 500
		b[i] = 500 + rng.Float64()*500
	}
	union := append(append([]float64{}, a...), b...)

	da := buildDigest(t, 200, a)
	db := buildDigest(t, 200, b)
	if err := da.Merge(db); err != nil {
		t.Fatalf("Merge: %v", err)
	}
	dUnion := buildDigest(t, 200, union)

	for _, q := range []float64{0.05, 0.25, 0.5, 0.75, 0.95} {
		merged, _ := da.Quantile(q)
		direct, _ := dUnion.Quantile(q)
		tol := tolerance(dataRange(union), q)
		if math.Abs(merged-direct) > tol {
			t.Errorf("merge(A,B).Quantile(%v) = %v vs digest(A∪B) = %v, |diff| %v > tol %v",
				q, merged, direct, math.Abs(merged-direct), tol)
		}
	}
	// Merge updates min/max across the union (exact extremes).
	unionMin, unionMax, _ := stats.MinMax(union)
	if minV, _ := da.Min(); minV != unionMin {
		t.Errorf("merged Min() = %v, want %v", minV, unionMin)
	}
	if maxV, _ := da.Max(); maxV != unionMax {
		t.Errorf("merged Max() = %v, want %v", maxV, unionMax)
	}
}

func TestDigest_MergeErrors(t *testing.T) {
	d := buildDigest(t, 100, []float64{1, 2, 3})
	snapshot := buildDigest(t, 100, []float64{1, 2, 3})
	beforeCount := d.Count()

	// Nil argument.
	if err := d.Merge(nil); !errors.Is(err, tdigest.ErrInvalidConfig) {
		t.Errorf("Merge(nil) error = %v, want ErrInvalidConfig", err)
	}
	if d.Count() != beforeCount {
		t.Errorf("receiver mutated after nil Merge: Count = %v, want %v", d.Count(), beforeCount)
	}

	// Compression mismatch.
	other := buildDigest(t, 50, []float64{4, 5, 6})
	if err := d.Merge(other); !errors.Is(err, tdigest.ErrInvalidConfig) {
		t.Errorf("Merge(mismatch) error = %v, want ErrInvalidConfig", err)
	}
	if d.Count() != beforeCount {
		t.Errorf("receiver mutated after mismatched Merge: Count = %v, want %v", d.Count(), beforeCount)
	}

	// A valid same-compression merge still works after the failures.
	if err := d.Merge(snapshot); err != nil {
		t.Errorf("valid Merge after errors: %v", err)
	}
	if d.Count() != beforeCount*2 {
		t.Errorf("Count after valid Merge = %v, want %v", d.Count(), beforeCount*2)
	}
}

func TestDigest_MergeEmpty(t *testing.T) {
	d := buildDigest(t, 100, []float64{1, 2, 3})
	empty := buildDigest(t, 100, nil)
	if err := d.Merge(empty); err != nil {
		t.Fatalf("Merge(empty): %v", err)
	}
	if d.Count() != 3 {
		t.Errorf("Count after merging empty = %v, want 3", d.Count())
	}
	// Merging into an empty digest also works.
	target := buildDigest(t, 100, nil)
	src := buildDigest(t, 100, []float64{10, 20, 30})
	if err := target.Merge(src); err != nil {
		t.Fatalf("Merge into empty: %v", err)
	}
	if got, ok := target.Quantile(0); !ok || got != 10 {
		t.Errorf("Quantile(0) after merge into empty = (%v, %v), want (10, true)", got, ok)
	}
}
