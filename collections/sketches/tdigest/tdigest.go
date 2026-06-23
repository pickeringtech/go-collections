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

// DefaultCompression is the compression used when WithCompression is not given.
// It sets the rough upper bound on the number of retained centroids and so
// trades memory for accuracy; 100 is the value recommended by Dunning's
// reference implementation as a good default.
const DefaultCompression = 100

// Digest is a t-digest: a streaming, mergeable sketch that estimates quantiles
// of a distribution of float64 values in space that scales with the compression
// parameter and is independent of the number of values seen. It is the
// approximate, bounded-memory counterpart of [stats.Quantile]: where the latter
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
	compression float64    // the delta parameter bounding centroid count
}

// centroid is a single cluster in the digest: the weighted mean of the points
// it absorbed.
type centroid struct {
	mean   float64
	weight float64
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

// WithCompression sets the compression parameter (delta), which bounds the
// number of retained centroids to roughly proportional to c: larger c keeps
// more centroids, costing memory but tightening the quantile estimates. c must
// be strictly positive.
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
// added w times. A non-finite x, or a weight that is not strictly positive
// (including NaN), is silently ignored — see Add for the rationale.
func (d *Digest) AddWeighted(x float64, w float64) {
	if nonFinite(x) || w <= 0 || math.IsNaN(w) {
		return
	}
	d.buffer = append(d.buffer, centroid{mean: x, weight: w})
	d.count += w
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
			// Absorb next into cur (weighted-mean update).
			cur.mean += (next.mean - cur.mean) * next.weight / (cur.weight + next.weight)
			cur.weight += next.weight
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
// its centre, and the observed maximum at rank count. The ranks are strictly
// increasing, so any interpolation between adjacent knots is well defined. The
// caller must ensure at least one centroid exists.
func (d *Digest) knots() []knot {
	ks := make([]knot, 0, len(d.centroids)+2)
	ks = append(ks, knot{value: d.min, rank: 0})
	weightSoFar := 0.0
	for _, c := range d.centroids {
		ks = append(ks, knot{value: c.mean, rank: weightSoFar + c.weight/2})
		weightSoFar += c.weight
	}
	ks = append(ks, knot{value: d.max, rank: d.count})
	return ks
}

// interpolate linearly maps target from the source interval [r0, r1] onto the
// destination interval [v0, v1]. Every call site supplies r1 > r0 (consecutive
// knot ranks, or values, are strictly increasing), so there is no
// divide-by-zero to guard against. The fraction is formed (and clamped to a
// finite ratio) before scaling by the destination span, so neither term can
// overflow to ±Inf even when the operands approach math.MaxFloat64 (e.g. a
// centroid carrying a weight near the float64 limit) — unlike the algebraically
// equal (target-r0)*(v1-v0)/(r1-r0), whose product can overflow first.
func interpolate(target, r0, v0, r1, v1 float64) float64 {
	return v0 + (target-r0)/(r1-r0)*(v1-v0)
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

// CDF returns the estimated fraction of the distribution that lies at or below
// x — the inverse of Quantile. The second return value is false (and the result
// 0) when the digest is empty or x is non-finite. Values below the observed
// minimum return 0 and values at or above the maximum return 1.
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
	if x <= d.min {
		// Below the minimum (and, by the check above, strictly below the
		// maximum): nothing lies at or below x.
		return 0, true
	}
	// d.min < x < d.max here, so x lies strictly inside the knot value range.
	// CDF is the inverse of Quantile: interpolate rank as a function of value in
	// the segment whose upper knot's value first reaches x, then normalise to a
	// fraction. (A single-centroid digest has min == max and is fully handled by
	// the bounds checks above, so there are at least two knot segments here.)
	ks := d.knots()
	i := 1
	for i < len(ks)-1 && ks[i].value < x {
		i++
	}
	// Knot values are non-decreasing and may repeat (e.g. the minimum knot and a
	// first centroid whose mean equals it); when the bracketing segment is
	// degenerate its rank is simply the shared knot's rank.
	rank := ks[i].rank
	if ks[i].value > ks[i-1].value {
		rank = interpolate(x, ks[i-1].value, ks[i-1].rank, ks[i].value, ks[i].rank)
	}
	return rank / d.count, true
}

// Merge folds other into the receiver: it absorbs other's centroids so the
// receiver then summarises the union of both streams. Both digests must have
// been built with identical compression; otherwise Merge returns an error
// wrapping ErrInvalidConfig and leaves the receiver unchanged. Merging a nil
// digest is also an error. Because the receiver is unchanged on error, callers
// can retry after correcting the configuration.
func (d *Digest) Merge(other *Digest) error {
	if other == nil {
		return fmt.Errorf("%w: cannot merge a nil digest", ErrInvalidConfig)
	}
	if d.compression != other.compression {
		return fmt.Errorf("%w: digests differ (compression=%v/%v)",
			ErrInvalidConfig, d.compression, other.compression)
	}
	// Pull other's buffered points and centroids into the receiver's buffer,
	// then let compress rebuild the summary. other is read, not mutated.
	d.buffer = append(d.buffer, other.centroids...)
	d.buffer = append(d.buffer, other.buffer...)
	d.count += other.count
	if other.min < d.min {
		d.min = other.min
	}
	if other.max > d.max {
		d.max = other.max
	}
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
