package tdigest_test

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"sort"
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
	// Clear resets both extremes, so the next value becomes the minimum and
	// the maximum, exactly as on a fresh digest.
	minV, minOK := d.Min()
	if !minOK || minV != 99 {
		t.Errorf("Min after Clear+Add = (%v, %v), want (99, true)", minV, minOK)
	}
	maxV, maxOK := d.Max()
	if !maxOK || maxV != 99 {
		t.Errorf("Max after Clear+Add = (%v, %v), want (99, true)", maxV, maxOK)
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

// TestDigest_MergeSelf checks that d.Merge(d) behaves exactly like merging an
// identical copy. Merge used to append other's centroids to the receiver's
// buffer and then read other.buffer — by then the receiver's own, already
// extended buffer — so a self-merge absorbed the centroids three times over:
// 1000 adds then Merge(self) reported Count 2000 but carried weight 2501, and
// the median moved from ~500 to ~333.
func TestDigest_MergeSelf(t *testing.T) {
	values := make([]float64, 1000)
	for i := range values {
		values[i] = float64(i)
	}
	tests := []struct {
		name    string
		prepare func(d *tdigest.Digest)
	}{
		{"with unmerged buffer", func(*tdigest.Digest) {}},
		{"after compression", func(d *tdigest.Digest) { _, _ = d.Quantile(0.5) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := buildDigest(t, 100, values)
			viaCopy := buildDigest(t, 100, values)
			duplicate := buildDigest(t, 100, values)
			tt.prepare(self)
			tt.prepare(viaCopy)
			tt.prepare(duplicate)

			err := self.Merge(self)
			if err != nil {
				t.Fatalf("Merge(self): %v", err)
			}
			err = viaCopy.Merge(duplicate)
			if err != nil {
				t.Fatalf("Merge(copy): %v", err)
			}
			if self.Count() != 2000 {
				t.Errorf("Count after Merge(self) = %v, want 2000", self.Count())
			}
			for _, q := range []float64{0.01, 0.25, 0.5, 0.75, 0.99} {
				got, _ := self.Quantile(q)
				want, _ := viaCopy.Quantile(q)
				if got != want {
					t.Errorf("Quantile(%v) after Merge(self) = %v, want %v (as merging a copy)", q, got, want)
				}
			}
			median, _ := self.Quantile(0.5)
			if math.Abs(median-499.5) > 10 {
				t.Errorf("median after Merge(self) = %v, want ≈ 499.5", median)
			}
		})
	}
}

// assertFiniteQueries fails the test if any query on d returns ok=false or a
// non-finite or out-of-range result.
func assertFiniteQueries(t *testing.T, d *tdigest.Digest) {
	t.Helper()
	if math.IsInf(d.Count(), 0) || math.IsNaN(d.Count()) {
		t.Fatalf("Count() = %v, want finite", d.Count())
	}
	lo, _ := d.Min()
	hi, _ := d.Max()
	for _, q := range []float64{0, 0.01, 0.25, 0.5, 0.75, 0.99, 1} {
		got, ok := d.Quantile(q)
		if !ok || math.IsNaN(got) || got < lo || got > hi {
			t.Errorf("Quantile(%v) = (%v, %v), want a value in [%v, %v]", q, got, ok, lo, hi)
		}
	}
	// lo/2 + hi/2 rather than (lo+hi)/2, which overflows for extreme data.
	for _, x := range []float64{lo, lo/2 + hi/2, hi} {
		got, ok := d.CDF(x)
		if !ok || math.IsNaN(got) || got < 0 || got > 1 {
			t.Errorf("CDF(%v) = (%v, %v), want a fraction in [0, 1]", x, got, ok)
		}
	}
}

// TestDigest_AddWeightedRejectsUnusableWeights checks that AddWeighted ignores
// a weight that is not a positive finite number, and a weight that would push
// the total weight past math.MaxFloat64. An infinite weight, or a total that
// overflowed to +Inf, used to be accepted and turned every later query into
// NaN with ok=true.
func TestDigest_AddWeightedRejectsUnusableWeights(t *testing.T) {
	tests := []struct {
		name      string
		weight    float64
		wantCount float64
	}{
		{name: "positive infinite weight is ignored", weight: math.Inf(1), wantCount: 2},
		{name: "negative infinite weight is ignored", weight: math.Inf(-1), wantCount: 2},
		{name: "NaN weight is ignored", weight: math.NaN(), wantCount: 2},
		{name: "zero weight is ignored", weight: 0, wantCount: 2},
		{name: "negative weight is ignored", weight: -1, wantCount: 2},
		{name: "a weight overflowing the total is ignored", weight: math.MaxFloat64, wantCount: math.MaxFloat64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := buildDigest(t, 100, []float64{2, 3})
			if tt.weight == math.MaxFloat64 {
				// Fill the digest to the limit first: this weight still fits
				// (2 + MaxFloat64 rounds to MaxFloat64), the second does not.
				d = buildDigest(t, 100, nil)
				d.AddWeighted(2, math.MaxFloat64)
			}
			d.AddWeighted(1, tt.weight)
			if d.Count() != tt.wantCount {
				t.Errorf("Count() = %v, want %v", d.Count(), tt.wantCount)
			}
			// The rejected value 1 sits below everything already added, so it
			// must not become the minimum either (a zero weight leaves Count
			// unchanged, so Count alone cannot tell it was ignored).
			minV, ok := d.Min()
			if !ok || minV != 2 {
				t.Errorf("Min() = (%v, %v), want (2, true)", minV, ok)
			}
			assertFiniteQueries(t, d)
		})
	}
}

// TestDigest_MergeRejectsWeightOverflow checks that a Merge whose combined
// weight would exceed math.MaxFloat64 returns ErrWeightOverflow and leaves the
// receiver unchanged, rather than overflowing the total to +Inf.
func TestDigest_MergeRejectsWeightOverflow(t *testing.T) {
	heavy := func() *tdigest.Digest {
		d := buildDigest(t, 100, nil)
		d.AddWeighted(1, math.MaxFloat64/2)
		d.AddWeighted(2, math.MaxFloat64/2)
		return d
	}
	tests := []struct {
		name  string
		other func(receiver *tdigest.Digest) *tdigest.Digest
	}{
		{"another heavy digest", func(*tdigest.Digest) *tdigest.Digest { return heavy() }},
		{"the receiver itself", func(receiver *tdigest.Digest) *tdigest.Digest { return receiver }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := heavy()
			before := d.Count()
			beforeMedian, _ := d.Quantile(0.5)
			err := d.Merge(tt.other(d))
			if !errors.Is(err, tdigest.ErrWeightOverflow) {
				t.Fatalf("Merge error = %v, want ErrWeightOverflow", err)
			}
			if d.Count() != before {
				t.Errorf("Count after rejected Merge = %v, want %v (unchanged)", d.Count(), before)
			}
			median, _ := d.Quantile(0.5)
			if median != beforeMedian {
				t.Errorf("median after rejected Merge = %v, want %v (unchanged)", median, beforeMedian)
			}
			assertFiniteQueries(t, d)
		})
	}
}

// TestDigest_ExtremeFiniteValuesStayFinite checks that finite inputs near
// ±math.MaxFloat64 never produce a non-finite estimate. The centroid-mean
// update computed (next-cur)·w/W and the interpolation computed v0 + f·(v1-v0);
// both overflow when the operands span more than math.MaxFloat64 (or when a
// huge weight multiplies the gap), so {-1.7e308, 1.7e308} gave Quantile(0.5) =
// +Inf and CDF(1e308) = NaN with ok=true, and compression could store a ±Inf
// centroid mean.
func TestDigest_ExtremeFiniteValuesStayFinite(t *testing.T) {
	type args struct {
		compression float64
		values      []float64
		weight      float64
	}
	repeated := func(n int, xs ...float64) []float64 {
		out := make([]float64, 0, n*len(xs))
		for i := 0; i < n; i++ {
			out = append(out, xs...)
		}
		return out
	}
	tests := []struct {
		name       string
		args       args
		wantMedian float64
		medianTol  float64
	}{
		{
			name:       "two opposite extremes",
			args:       args{compression: 100, values: []float64{-1.7e308, 1.7e308}, weight: 1},
			wantMedian: 0,
		},
		{
			name:       "the float64 range end to end",
			args:       args{compression: 100, values: []float64{-math.MaxFloat64, 0, math.MaxFloat64}, weight: 1},
			wantMedian: 0,
		},
		{
			name:       "many opposite extremes merged into centroids",
			args:       args{compression: 100, values: repeated(500, -1.7e308, 1.7e308), weight: 1},
			wantMedian: 0,
			medianTol:  1.7e308,
		},
		{
			// Halving a subnormal rank loses it entirely, so the overflow-safe
			// path must only be taken when a span actually overflows.
			name:       "subnormal weights",
			args:       args{compression: 100, values: []float64{0, 1, 2, 3, 4}, weight: 5e-324},
			wantMedian: 2,
		},
		{
			name: "huge weights merged at low compression",
			args: args{compression: 1, values: []float64{0, 10, 20, 30, 40}, weight: math.MaxFloat64 / 8},
			// Every value carries the same weight, so the median is the middle value
			// up to the digest's coarseness at compression 1.
			wantMedian: 20,
			medianTol:  20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := buildDigest(t, tt.args.compression, nil)
			for _, v := range tt.args.values {
				d.AddWeighted(v, tt.args.weight)
			}
			assertFiniteQueries(t, d)
			median, _ := d.Quantile(0.5)
			if math.Abs(median-tt.wantMedian) > tt.medianTol {
				t.Errorf("Quantile(0.5) = %v, want %v ± %v", median, tt.wantMedian, tt.medianTol)
			}
			c, ok := d.CDF(1e308)
			if !ok || math.IsNaN(c) || c < 0 || c > 1 {
				t.Errorf("CDF(1e308) = (%v, %v), want a fraction in [0, 1]", c, ok)
			}
		})
	}
}

// exactCDF returns the fraction of values at or below x — the empirical CDF
// that Digest.CDF estimates. sorted must be in ascending order.
func exactCDF(sorted []float64, x float64) float64 {
	atOrBelow := sort.Search(len(sorted), func(i int) bool { return sorted[i] > x })
	return float64(atOrBelow) / float64(len(sorted))
}

// TestDigest_CDFAtRepeatedValues checks CDF's "fraction at or below x"
// contract where values repeat, including at the minimum. CDF used to read
// the rank of the first centroid at x's centre, leaving out the weight sitting
// exactly at x: 9×0 then 1 gave CDF(0) = 0 (true 0.9), and 1, 50×5, 9 gave
// CDF(5) ≈ 0.03 (true 51/52).
func TestDigest_CDFAtRepeatedValues(t *testing.T) {
	type point struct {
		value, weight float64
	}
	repeat := func(value float64, n int) []point {
		out := make([]point, n)
		for i := range out {
			out[i] = point{value: value, weight: 1}
		}
		return out
	}
	concat := func(parts ...[]point) []point {
		var out []point
		for _, p := range parts {
			out = append(out, p...)
		}
		return out
	}
	tests := []struct {
		name   string
		points []point
		x      float64
		want   float64
	}{
		{"repeated minimum", concat(repeat(0, 9), repeat(1, 1)), 0, 0.9},
		{"just above a repeated minimum", concat(repeat(0, 9), repeat(1, 1)), 0.5, 0.9},
		{"repeated interior value", concat(repeat(1, 1), repeat(5, 50), repeat(9, 1)), 5, 51.0 / 52},
		{"just below a repeated interior value", concat(repeat(1, 1), repeat(5, 50), repeat(9, 1)), 4.999, 1.0 / 52},
		{"singleton minimum", concat(repeat(1, 1), repeat(5, 50), repeat(9, 1)), 1, 1.0 / 52},
		{"heavily weighted interior point", []point{{1, 1}, {5, 1000}, {9, 1}}, 5, 1001.0 / 1002},
		{"distinct singletons", concat(repeat(1, 1), repeat(2, 1), repeat(3, 1), repeat(4, 1), repeat(5, 1)), 3, 0.6},
		{"between distinct singletons", concat(repeat(1, 1), repeat(2, 1), repeat(3, 1), repeat(4, 1), repeat(5, 1)), 3.5, 0.6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := buildDigest(t, 100, nil)
			for _, p := range tt.points {
				d.AddWeighted(p.value, p.weight)
			}
			got, ok := d.CDF(tt.x)
			if !ok || math.Abs(got-tt.want) > 1e-12 {
				t.Errorf("CDF(%v) = (%v, %v), want (%v, true)", tt.x, got, ok, tt.want)
			}
		})
	}
}

// TestDigest_CDFMatchesEmpiricalCDF compares CDF against the exact empirical
// CDF on large streams, with and without ties, at every distinct value of a
// discrete stream and at a grid of points of a continuous one. Ties are where
// the old midpoint-only estimate was furthest off: it dropped the whole
// weight sitting at x.
func TestDigest_CDFMatchesEmpiricalCDF(t *testing.T) {
	rng := rand.New(rand.NewPCG(5, 9))
	discrete := make([]float64, 20000)
	for i := range discrete {
		discrete[i] = float64(rng.IntN(20)) // 20 distinct values, ~5% of the weight each
	}
	continuous := make([]float64, 20000)
	for i := range continuous {
		continuous[i] = rng.NormFloat64()
	}
	mixed := make([]float64, 0, 20000)
	mixed = append(mixed, continuous[:15000]...)
	for i := 0; i < 5000; i++ {
		mixed = append(mixed, 0.25) // a 25% point mass inside a continuous body
	}
	gridOf := func(values []float64) []float64 {
		lo, hi, _ := stats.MinMax(values)
		grid := []float64{0.25}
		for i := 0; i <= 100; i++ {
			grid = append(grid, lo+(hi-lo)*float64(i)/100)
		}
		return grid
	}
	tests := []struct {
		name   string
		values []float64
		at     []float64
	}{
		{"discrete values", discrete, []float64{0, 1, 2, 5, 9.5, 10, 13, 18, 19}},
		{"continuous values", continuous, gridOf(continuous)},
		{"point mass in a continuous body", mixed, gridOf(mixed)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := buildDigest(t, 100, tt.values)
			sorted := append([]float64(nil), tt.values...)
			sort.Float64s(sorted)
			for _, x := range tt.at {
				got, ok := d.CDF(x)
				want := exactCDF(sorted, x)
				if !ok || math.Abs(got-want) > 0.02 {
					t.Errorf("CDF(%v) = (%v, %v), exact %v, |diff| > 0.02", x, got, ok, want)
				}
			}
		})
	}
}

// TestDigest_CDFIsMonotone checks that CDF never decreases as x grows — across
// point masses, wide centroids and the extremes — and stays in [0, 1].
func TestDigest_CDFIsMonotone(t *testing.T) {
	rng := rand.New(rand.NewPCG(2, 4))
	d := buildDigest(t, 50, nil)
	for i := 0; i < 5000; i++ {
		d.Add(math.Round(rng.NormFloat64()*4) / 2) // many ties on a half-integer grid
		d.Add(rng.NormFloat64() * 3)
	}
	other := buildDigest(t, 50, nil)
	for i := 0; i < 3000; i++ {
		other.Add(math.Round(rng.NormFloat64() * 3))
	}
	err := d.Merge(other)
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	prev := 0.0
	for x := -20.0; x <= 20; x += 0.01 {
		got, ok := d.CDF(x)
		if !ok || got < prev || got > 1 {
			t.Fatalf("CDF(%v) = (%v, %v) after CDF = %v just below it; want non-decreasing in [0, 1]", x, got, ok, prev)
		}
		prev = got
	}
}

// TestDigest_CDFAtMinimumCountsTheMinimum pins the documented x == Min()
// behaviour: from compression 4 upwards the smallest observation always stays
// a point mass, so CDF(Min()) is at least its weight's share, even after a
// Merge. (Below 4 the scale function lets it merge, and the doc says so.)
func TestDigest_CDFAtMinimumCountsTheMinimum(t *testing.T) {
	for _, compression := range []float64{4, 20, 100, 500} {
		t.Run(fmt.Sprintf("compression %v", compression), func(t *testing.T) {
			rng := rand.New(rand.NewPCG(uint64(compression), 1))
			a := buildDigest(t, compression, nil)
			b := buildDigest(t, compression, nil)
			for i := 0; i < 3000; i++ {
				a.Add(rng.NormFloat64())
				b.Add(rng.NormFloat64() + 1)
			}
			err := a.Merge(b)
			if err != nil {
				t.Fatalf("Merge: %v", err)
			}
			lo, _ := a.Min()
			got, ok := a.CDF(lo)
			if !ok || got < 1/a.Count() {
				t.Errorf("CDF(Min()) = (%v, %v), want at least 1/Count = %v", got, ok, 1/a.Count())
			}
		})
	}
}

// TestDigest_SmallDigestExactEstimates pins the merge criterion and the CDF
// interpolation on digests small enough to work out by hand. At compression 1,
// compressing 1, 2, 3, 4 (unit weights, total 4) absorbs 2 into 1 (combined
// weight 2 against a bound of 4·4·½·½ = 4) and then 3 (combined weight 3
// against 4·4·¾·¼ = 3: the bound is inclusive), giving one spread centroid
// (mean 2, weight 3) and a point at 4 (the bound at q = 1 is 0). The quantile
// knots are therefore (1, 0), (2, 1.5), (4, 3.5) and (4, 4), and the CDF rises
// linearly from 0 at 1 to 1.5/4 at 2, then towards 3/4 just below 4.
//
// -1.7e308, -1.6e308, 1.6e308, 1.7e308 compress the same way: the three lower
// values merge into a spread centroid at m = (-1.7e308 - 1.6e308 + 1.6e308)/3
// (weight 3) and 1.7e308 stays a point. The CDF segment between them spans more
// than math.MaxFloat64, so its fraction must be formed from halved operands:
// CDF(1e308) = (1.5 + 1.5·(1e308 - m)/(1.7e308 - m))/4.
func TestDigest_SmallDigestExactEstimates(t *testing.T) {
	quantile := func(q float64) func(d *tdigest.Digest) (float64, bool) {
		return func(d *tdigest.Digest) (float64, bool) { return d.Quantile(q) }
	}
	cdf := func(x float64) func(d *tdigest.Digest) (float64, bool) {
		return func(d *tdigest.Digest) (float64, bool) { return d.CDF(x) }
	}
	oneToFour := []float64{1, 2, 3, 4}
	extremes := []float64{-1.7e308, -1.6e308, 1.6e308, 1.7e308}
	type args struct {
		values   []float64
		estimate func(d *tdigest.Digest) (float64, bool)
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "a low quantile interpolates from the minimum to the spread centroid", args: args{values: oneToFour, estimate: quantile(0.1)}, want: 1 + 0.4/1.5},
		{name: "the median interpolates between the spread centroid and the point", args: args{values: oneToFour, estimate: quantile(0.5)}, want: 2.5},
		{name: "the upper quartile interpolates between the spread centroid and the point", args: args{values: oneToFour, estimate: quantile(0.75)}, want: 3.5},
		{name: "CDF between the minimum and the spread centroid interpolates from 0", args: args{values: oneToFour, estimate: cdf(1.5)}, want: 0.1875},
		{name: "CDF at the spread centroid's mean counts half its weight", args: args{values: oneToFour, estimate: cdf(2)}, want: 0.375},
		{name: "CDF between the spread centroid and the point interpolates", args: args{values: oneToFour, estimate: cdf(3)}, want: 0.5625},
		{name: "CDF just below the point stays on the interpolated segment", args: args{values: oneToFour, estimate: cdf(3.5)}, want: 0.65625},
		{name: "CDF across a segment wider than MaxFloat64 interpolates on halved operands", args: args{values: extremes, estimate: cdf(1e308)}, want: 0.6341911764705882},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := buildDigest(t, 1, tt.args.values)
			got, ok := tt.args.estimate(d)
			if !ok || math.Abs(got-tt.want) > 1e-12 {
				t.Errorf("estimate = (%v, %v), want (%v, true)", got, ok, tt.want)
			}
		})
	}
}

// TestDigest_FirstBatchIsOrderIndependent pins the buffering schedule: Add
// buffers values and compresses only once more than 5·compression are pending,
// so the first 5·compression+1 values are folded in one sorted batch and their
// insertion order cannot change the summary. Compressing any earlier — at
// exactly 5·compression values, or after every Add — would merge part of the
// batch on its own first and make the estimates depend on the order. The values
// are distinct, so the sort that starts each compression has no ties to break.
func TestDigest_FirstBatchIsOrderIndependent(t *testing.T) {
	type args struct {
		compression float64
		values      []float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "eleven values at compression 2 with the minimum added last",
			args: args{compression: 2, values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0}},
		},
		{
			name: "sixteen values at compression 3 with a low value added last",
			args: args{compression: 3, values: []float64{0, 1, 2, 3, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 4}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ascending := append([]float64(nil), tt.args.values...)
			sort.Float64s(ascending)
			inOrder := buildDigest(t, tt.args.compression, tt.args.values)
			sorted := buildDigest(t, tt.args.compression, ascending)
			for i := 1; i < 20; i++ {
				q := float64(i) / 20
				got, _ := inOrder.Quantile(q)
				want, _ := sorted.Quantile(q)
				if got != want {
					t.Errorf("Quantile(%v) = %v, want %v (as when added in ascending order)", q, got, want)
				}
			}
			for x := -0.5; x <= 16; x += 0.25 {
				got, _ := inOrder.CDF(x)
				want, _ := sorted.CDF(x)
				if got != want {
					t.Errorf("CDF(%v) = %v, want %v (as when added in ascending order)", x, got, want)
				}
			}
		})
	}
}
