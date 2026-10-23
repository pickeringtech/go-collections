package tdigest

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

// ErrInvalidConfig is returned by New/NewConcurrent when the requested
// compression is out of range, and is the root of the error Merge returns when
// two digests were built with different compression, so callers can test with
// errors.Is.
var ErrInvalidConfig = errors.New("tdigest: invalid configuration")

// ErrWeightOverflow is returned by Merge when the two digests' combined total
// weight would exceed math.MaxFloat64. The receiver is left unchanged, as it is
// for every Merge error.
var ErrWeightOverflow = errors.New("tdigest: total weight would overflow float64")

// DefaultCompression is the compression used when WithCompression is not given.
// It scales the number of retained centroids, and so trades memory for
// accuracy; 100 is the value recommended by Dunning's reference implementation
// as a good default. At this compression a digest holds roughly 280 centroids
// after 1e3 values, 450 after 1e4 and 900 after 1e7 (see WithCompression).
const DefaultCompression = 100

// Digest is a t-digest: a streaming, mergeable sketch that estimates quantiles
// of a distribution of float64 values in space that scales with the compression
// parameter and grows only logarithmically with the number of values seen (see
// WithCompression). It is the approximate, small-memory counterpart of
// [stats.Quantile]: where the latter
// sorts the whole sample exactly, a Digest keeps only a small set of weighted
// centroids whose density follows a scale function, so error is smallest at the
// tails (q near 0 or 1) and largest in the middle.
//
// Unlike the comparable-typed sketches in the sibling packages (bloom, countmin,
// hll), Digest is deliberately not generic: it operates on float64 only,
// matching the value domain of the stats quantile functions it approximates.
//
// A Digest is deterministic given a fixed sequence of Add/AddWeighted/Merge
// calls — it uses no randomness — but the retained centroids, and therefore the
// estimates, can depend on insertion and merge order. Different orderings of the
// same data yield close but not bit-identical results.
//
// A Digest is not safe for concurrent use. Wrap one with NewConcurrent for a
// goroutine-safe variant.
//
// The zero value is not usable; construct a Digest with New.
type Digest struct {
	centroids   []centroid // sorted by mean, ascending; the compressed summary
	buffer      []centroid // unmerged points awaiting the next compression
	count       float64    // total weight across all centroids (and the buffer)
	min         float64    // smallest value ever added (+Inf when empty)
	max         float64    // largest value ever added (-Inf when empty)
	compression float64    // the delta parameter scaling the centroid count
}

// centroid is a single cluster in the digest: the weighted mean of the points
// it absorbed.
type centroid struct {
	mean   float64
	weight float64
	// point records that every value the centroid absorbed equals mean, so its
	// whole weight sits exactly at mean (a point mass) rather than being spread
	// around it. CDF needs this to count repeated values as "at or below" x.
	point bool
}

// bufferFactor sizes the unmerged buffer relative to compression. A larger
// buffer amortises the O(n log n) compression cost over more Add calls; this
// multiple is the one used by Dunning's MergingDigest.
const bufferFactor = 5

// Interface guard.
var _ Quantiles = (*Digest)(nil)

// New creates a t-digest with DefaultCompression, adjustable via
// WithCompression. It returns an error wrapping ErrInvalidConfig if the
// resulting compression is not strictly positive (zero, negative or NaN).
func New(opts ...Option) (*Digest, error) {
	d := &Digest{
		compression: DefaultCompression,
		min:         math.Inf(1),
		max:         math.Inf(-1),
	}
	for _, opt := range opts {
		opt(d)
	}
	if d.compression <= 0 || math.IsNaN(d.compression) {
		return nil, fmt.Errorf("%w: compression must be > 0, got %v", ErrInvalidConfig, d.compression)
	}
	return d, nil
}

// Option configures a Digest at construction time.
type Option func(*Digest)

// WithCompression sets the compression parameter (delta), which scales the
// number of retained centroids: larger c keeps more centroids, costing memory
// but tightening the quantile estimates. c must be strictly positive.
//
// The centroid count is proportional to c but not independent of the stream
// length: the q(1-q) scale function keeps tail centroids small however long the
// stream runs, so the count grows logarithmically with the number of values n —
// each tenfold increase in n adds about 1.5·c centroids (at c = 100, roughly 280
// after 1e3 values and 900 after 1e7). On top of that, up to 5·c values are
// buffered between compressions.
func WithCompression(c float64) Option {
	return func(d *Digest) { d.compression = c }
}

// Add records a single value with weight 1. Non-finite values (NaN or ±Inf) are
// silently ignored: Add has no error channel, and a non-finite value has no
// meaningful place in a quantile summary (it matches the rejection that
// [stats.Quantile] performs by returning ok=false).
func (d *Digest) Add(x float64) {
	d.AddWeighted(x, 1)
}

// AddWeighted records a value carrying the given weight, as though it had been
// added w times. The call is silently ignored — see Add for the rationale —
// when x is non-finite, when w is not a strictly positive finite number (zero,
// negative, NaN or ±Inf), or when adding w would push the digest's total
// weight (Count) past math.MaxFloat64. Rejecting the overflow keeps Count
// finite, so the quantile queries never divide by an infinite total.
func (d *Digest) AddWeighted(x float64, w float64) {
	if nonFinite(x) || nonFinite(w) || w <= 0 {
		return
	}
	total := d.count + w
	if math.IsInf(total, 1) {
		return
	}
	d.buffer = append(d.buffer, centroid{mean: x, weight: w, point: true})
	d.count = total
	if x < d.min {
		d.min = x
	}
	if x > d.max {
		d.max = x
	}
	if float64(len(d.buffer)) > bufferFactor*d.compression {
		d.compress()
	}
}

// compress folds the unmerged buffer into the retained centroids. It sorts the
// union by mean and sweeps left to right, accumulating points into the current
// centroid until adding the next would push that centroid past the weight bound
// the scale function permits at its quantile, at which point a new centroid is
// started. This is the standard "merging" t-digest update.
func (d *Digest) compress() {
	if len(d.buffer) == 0 {
		return
	}
	merged := make([]centroid, 0, len(d.centroids)+len(d.buffer))
	merged = append(merged, d.centroids...)
	merged = append(merged, d.buffer...)
	d.buffer = d.buffer[:0]
	sort.Slice(merged, func(i, j int) bool { return merged[i].mean < merged[j].mean })

	out := make([]centroid, 0, len(merged))
	totalWeight := d.count
	cur := merged[0]
	soFar := 0.0 // weight of all centroids already finalised in out
	for i := 1; i < len(merged); i++ {
		next := merged[i]
		// Quantile at the right edge of the proposed combined centroid.
		q := (soFar + cur.weight + next.weight) / totalWeight
		limit := totalWeight * sizeBound(q, d.compression)
		if cur.weight+next.weight <= limit {
			// Absorb next into cur: the combined mean lies next.weight/W of the
			// way from cur.mean to next.mean, where W is the combined weight.
			// lerp keeps it finite even when the means sit near ±MaxFloat64.
			cur.point = cur.point && next.point && cur.mean == next.mean
			cur.weight += next.weight
			cur.mean = lerp(cur.mean, next.mean, next.weight/cur.weight)
		} else {
			out = append(out, cur)
			soFar += cur.weight
			cur = next
		}
	}
	out = append(out, cur)
	d.centroids = out
}

// sizeBound returns the maximum fraction of the total weight a single centroid
// may carry at quantile q (which compress always supplies in (0, 1]), derived
// from the q(1-q) scale function: centroids may be large in the middle of the
// distribution and must be small at the tails, which is what gives the t-digest
// its high accuracy for extreme quantiles.
func sizeBound(q float64, compression float64) float64 {
	return 4 * q * (1 - q) / compression
}

// Quantile returns the estimated q-quantile for q in [0, 1], mirroring the edge
// semantics of [stats.Quantile]: q=0 is the minimum, q=1 the maximum and q=0.5
// the median. The second return value is false (and the result 0) when the
// digest is empty or when q is outside [0, 1] or NaN.
func (d *Digest) Quantile(q float64) (float64, bool) {
	if math.IsNaN(q) || q < 0 || q > 1 {
		return 0, false
	}
	d.compress()
	if len(d.centroids) == 0 {
		return 0, false
	}
	if len(d.centroids) == 1 {
		return d.centroids[0].mean, true
	}
	if q == 0 {
		return d.min, true
	}
	if q == 1 {
		return d.max, true
	}

	// Target cumulative weight (rank) we are looking for. Walk the knot
	// sequence and interpolate value as a function of rank in the segment that
	// brackets the target.
	target := q * d.count
	ks := d.knots()
	// q in (0,1) puts target in (0, count); find the segment whose upper knot's
	// rank first reaches target. The last segment ends at rank count >= target,
	// so the final iteration always matches.
	i := 1
	for i < len(ks)-1 && ks[i].rank < target {
		i++
	}
	return interpolate(target, ks[i-1].rank, ks[i-1].value, ks[i].rank, ks[i].value), true
}

// knot is a point on the digest's empirical distribution: the value at a given
// cumulative weight (rank).
type knot struct {
	value float64
	rank  float64
}

// knots returns the rank/value control points used by Quantile and CDF: the
// observed minimum at rank 0, each centroid's mean at the cumulative weight of
// its centre, and the observed maximum at rank count. The ranks are
// non-decreasing (strictly increasing unless subnormal weights round a
// half-weight away; interpolate tolerates the resulting zero-width segments).
// The caller must ensure at least one centroid exists.
func (d *Digest) knots() []knot {
	ks := make([]knot, 0, len(d.centroids)+2)
	ks = append(ks, knot{value: d.min, rank: 0})
	weightSoFar := 0.0
	for _, c := range d.centroids {
		// float64(...) forces c.weight/2 to round on its own: without it Go may
		// fuse the halving into the addition (an FMA on arm64, ppc64le, s390x),
		// which rounds differently. With subnormal weights that difference moves
		// whole ranks, so the conversion keeps results identical on every
		// platform. The same guard appears wherever a product feeds a sum below.
		ks = append(ks, knot{value: c.mean, rank: weightSoFar + float64(c.weight/2)})
		weightSoFar += c.weight
	}
	ks = append(ks, knot{value: d.max, rank: d.count})
	return ks
}

// interpolate linearly maps target from the source interval [r0, r1] onto the
// destination interval [v0, v1]. Every call site supplies r1 >= r0 with target
// in [r0, r1] (adjacent knot ranks, or values, that bracket target). Both steps
// — the fraction along the source interval and the point along the destination
// one — are computed without overflowing to ±Inf even when an interval spans
// more than math.MaxFloat64 (e.g. -1.7e308 to 1.7e308): see fraction and lerp.
func interpolate(target, r0, v0, r1, v1 float64) float64 {
	return lerp(v0, v1, fraction(target, r0, r1))
}

// fraction returns how far target lies along [r0, r1] (r1 >= r0), clamped to
// [0, 1]. A zero-width interval yields 0: knot ranks can coincide when weights
// are so small (subnormal) that halving one rounds it away, and q·Count can then
// underflow onto that shared rank, which would otherwise divide 0 by 0. When the
// span r1-r0 overflows to +Inf, the ratio is formed from halved operands
// instead, which cannot overflow (both endpoints are then far from the
// subnormal range, so halving loses nothing).
func fraction(target, r0, r1 float64) float64 {
	span := r1 - r0
	if span == 0 {
		return 0
	}
	f := (target - r0) / span
	if math.IsInf(span, 0) {
		f = (float64(target/2) - float64(r0/2)) / (float64(r1/2) - float64(r0/2))
	}
	return min(max(f, 0), 1)
}

// lerp returns the point a fraction frac in [0, 1] of the way from v0 to v1:
// v0 + frac·(v1-v0). When the span v1-v0 overflows to ±Inf (v0 and v1 near
// ±math.MaxFloat64 with opposite signs) it is evaluated on halved operands
// instead, as 2·(v0/2 + frac·(v1/2 - v0/2)), which cannot overflow. The result
// is clamped into the closed interval between v0 and v1, so rounding can never
// step outside it; in particular lerp(v, v, f) is exactly v.
func lerp(v0, v1, frac float64) float64 {
	span := v1 - v0
	point := v0 + float64(frac*span)
	if math.IsInf(span, 0) {
		point = 2 * (float64(v0/2) + float64(frac*(float64(v1/2)-float64(v0/2))))
	}
	return min(max(point, min(v0, v1)), max(v0, v1))
}

// Percentile returns the estimated p-th percentile for p in [0, 100]. It is
// exactly Quantile(p/100); see Quantile for the empty/range/NaN contract and
// mirrors [stats.Percentile].
func (d *Digest) Percentile(p float64) (float64, bool) {
	if math.IsNaN(p) || p < 0 || p > 100 {
		return 0, false
	}
	return d.Quantile(p / 100)
}

// CDF returns the estimated fraction of the total weight at or below x. The
// second return value is false (and the result 0) when the digest is empty or x
// is non-finite. Otherwise the result is in [0, 1] and non-decreasing in x:
//
//   - x below Min() gives 0, and x at or above Max() gives 1.
//   - Weight known to sit exactly at a value counts in full once x reaches that
//     value: every centroid that has only ever absorbed copies of one value (an
//     unmerged Add or AddWeighted, or a run of equal values) is a point mass, so
//     repeated values — including a repeated minimum — are counted as at or
//     below x. In particular CDF(Min()) is the fraction of weight in point
//     masses at the minimum. With a compression of 4 or more (the default is
//     100) the smallest observation is never merged into a wider centroid, so
//     this always includes it; below 4 it can be, and CDF(Min()) may then be 0.
//   - Weight in a centroid that merged several distinct values is assumed
//     spread evenly either side of its mean: half of it counts as below the
//     mean, and the CDF rises linearly between adjacent centroid means (and
//     between the extreme centroids and Min()/Max()).
//
// Equal-mean centroids are pooled before this rule is applied. CDF is an
// approximate inverse of Quantile: CDF(Quantile(q)) ≈ q away from point masses.
func (d *Digest) CDF(x float64) (float64, bool) {
	if nonFinite(x) {
		return 0, false
	}
	d.compress()
	if len(d.centroids) == 0 {
		return 0, false
	}
	if x >= d.max {
		// At or above the observed maximum the whole distribution lies at or
		// below x. Checked before the minimum so an all-equal digest (min == max)
		// reports 1 at that single value: every observation is at or below it.
		return 1, true
	}
	if x < d.min {
		return 0, true
	}
	// d.min <= x < d.max here. Take the last knot at or below x: if it sits at
	// x, its rank — the highest at that value, so point masses at x count in
	// full — is the answer; otherwise x lies strictly inside the segment to the
	// next knot, whose value exceeds x (the final knot is the maximum, > x).
	ks := d.cdfKnots()
	j := 0
	for ks[j+1].value <= x {
		j++
	}
	rank := ks[j].rank
	if ks[j].value < x {
		rank = interpolate(x, ks[j].value, ks[j].rank, ks[j+1].value, ks[j+1].rank)
	}
	// Summing the centroid weights in sorted order can round a hair past the
	// separately accumulated count, so clamp the fraction.
	return min(rank/d.count, 1), true
}

// cdfKnots returns the control points of CDF's estimate as (value, rank) pairs
// with non-decreasing values and ranks: the minimum at rank 0, two knots per
// run of equal-mean centroids, and the maximum at rank count. A run at mean m,
// preceded by weight below, holding spread weight s (centroids that merged
// distinct values) and point weight p (centroids holding only m), contributes
// (m, below + s/2) — the CDF just below m, half the spread weight having
// arrived — and (m, below + s/2 + p) — the CDF at m, after the jump of the
// point mass. The CDF is linear between consecutive knots of different value.
// The caller must ensure at least one centroid exists.
func (d *Digest) cdfKnots() []knot {
	ks := make([]knot, 0, 2*len(d.centroids)+2)
	ks = append(ks, knot{value: d.min, rank: 0})
	below := 0.0
	for i := 0; i < len(d.centroids); {
		mean := d.centroids[i].mean
		spread, point := 0.0, 0.0
		for ; i < len(d.centroids) && d.centroids[i].mean == mean; i++ {
			if d.centroids[i].point {
				point += d.centroids[i].weight
			} else {
				spread += d.centroids[i].weight
			}
		}
		ks = append(ks,
			knot{value: mean, rank: below + float64(spread/2)},
			knot{value: mean, rank: below + float64(spread/2) + point})
		below += spread + point
	}
	return append(ks, knot{value: d.max, rank: d.count})
}

// Merge folds other into the receiver: it absorbs other's centroids so the
// receiver then summarises the union of both streams. Both digests must have
// been built with identical compression; otherwise Merge returns an error
// wrapping ErrInvalidConfig and leaves the receiver unchanged. Merging a nil
// digest is also an error (wrapping ErrInvalidConfig), as is a merge whose
// combined total weight would exceed math.MaxFloat64 (wrapping
// ErrWeightOverflow). Because the receiver is unchanged on error, callers
// can retry after correcting the configuration. Merging a digest into itself
// is allowed and is equivalent to merging an identical copy: every value then
// counts twice.
func (d *Digest) Merge(other *Digest) error {
	if other == nil {
		return fmt.Errorf("%w: cannot merge a nil digest", ErrInvalidConfig)
	}
	if d.compression != other.compression {
		return fmt.Errorf("%w: digests differ (compression=%v/%v)",
			ErrInvalidConfig, d.compression, other.compression)
	}
	if math.IsInf(d.count+other.count, 1) {
		return fmt.Errorf("%w: cannot combine total weights %v and %v",
			ErrWeightOverflow, d.count, other.count)
	}
	// Snapshot everything read from other before the receiver is touched:
	// other may be the receiver itself (d.Merge(d)), in which case appending to
	// d.buffer would also grow the other.buffer read afterwards. other is read,
	// not mutated.
	incoming := make([]centroid, 0, len(other.centroids)+len(other.buffer))
	incoming = append(incoming, other.centroids...)
	incoming = append(incoming, other.buffer...)
	count, lo, hi := other.count, other.min, other.max

	// Pull the snapshot into the receiver's buffer, then let compress rebuild
	// the summary.
	d.buffer = append(d.buffer, incoming...)
	d.count += count
	d.min = min(d.min, lo)
	d.max = max(d.max, hi)
	d.compress()
	return nil
}

// Clear empties the digest, returning it to the state a freshly constructed one
// would have, while keeping its compression.
func (d *Digest) Clear() {
	d.centroids = d.centroids[:0]
	d.buffer = d.buffer[:0]
	d.count = 0
	d.min = math.Inf(1)
	d.max = math.Inf(-1)
}

// Count returns the total weight added (the sum of the weights of every value,
// 1 per Add). It is a float64 because AddWeighted accepts fractional weights.
func (d *Digest) Count() float64 { return d.count }

// Min returns the smallest value ever added. The second return value is false
// (and the result 0) when the digest is empty.
func (d *Digest) Min() (float64, bool) {
	if d.count == 0 {
		return 0, false
	}
	return d.min, true
}

// Max returns the largest value ever added. The second return value is false
// (and the result 0) when the digest is empty.
func (d *Digest) Max() (float64, bool) {
	if d.count == 0 {
		return 0, false
	}
	return d.max, true
}

// Compression returns the compression parameter the digest was built with.
func (d *Digest) Compression() float64 { return d.compression }

// nonFinite reports whether f is NaN or ±Inf.
func nonFinite(f float64) bool {
	return math.IsNaN(f) || math.IsInf(f, 0)
}
