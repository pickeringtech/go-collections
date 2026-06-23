package streaming

import "github.com/pickeringtech/go-collections/constraints"

// RunningMean maintains the arithmetic mean of an unbounded stream of float64
// values in O(1) memory, updated one element at a time. It accumulates the mean
// incrementally (mean += (x − mean)/n) rather than summing then dividing, so it
// never holds a large running total that could overflow or lose precision on a
// long stream.
//
// RunningMean is single-threaded; see the package documentation on thread
// safety.
type RunningMean struct {
	n    int
	mean float64
}

// NewRunningMean creates an empty RunningMean. Before any Add, Result reports
// ok == false.
func NewRunningMean() *RunningMean {
	return &RunningMean{}
}

// Add folds one value into the running mean in O(1).
func (m *RunningMean) Add(x float64) {
	m.n++
	m.mean += (x - m.mean) / float64(m.n)
}

// Result returns the current mean. The second return value is false — and the
// mean 0 — when no value has been added yet.
func (m *RunningMean) Result() (float64, bool) {
	if m.n == 0 {
		return 0, false
	}
	return m.mean, true
}

// Count returns the number of values added so far.
func (m *RunningMean) Count() int {
	return m.n
}

// RunningVariance maintains the mean and variance of an unbounded stream of
// float64 values in O(1) memory using Welford's online algorithm, accumulating
// n, the running mean, and M2 (the sum of squared deviations from that mean).
//
// This is an incremental reimplementation of the recurrence stats.welford runs
// over a slice. The streaming form is deliberate: stats.welford is an
// unexported, batch-over-slice helper that cannot be fed one element at a time
// without first buffering the whole stream, which defeats the point of a
// bounded-memory aggregate. The recurrence — and therefore the
// SampleVariance/PopulationVariance contracts below — intentionally match
// stats.SampleVariance and stats.PopulationVariance exactly, so a RunningVariance
// over a stream agrees with stats over the same data buffered into a slice.
// Nothing in the stats package is exported or modified to share this.
//
// As in stats, non-finite inputs (NaN/Inf) propagate to a non-finite result
// rather than being silently dropped. RunningVariance is single-threaded; see
// the package documentation on thread safety.
type RunningVariance struct {
	n    int
	mean float64
	m2   float64
}

// NewRunningVariance creates an empty RunningVariance.
func NewRunningVariance() *RunningVariance {
	return &RunningVariance{}
}

// Add folds one value into the running statistics in O(1) using Welford's
// update.
func (v *RunningVariance) Add(x float64) {
	v.n++
	delta := x - v.mean
	v.mean += delta / float64(v.n)
	delta2 := x - v.mean
	v.m2 += delta * delta2
}

// Mean returns the current mean. The second return value is false — and the
// mean 0 — when no value has been added yet.
func (v *RunningVariance) Mean() (float64, bool) {
	if v.n == 0 {
		return 0, false
	}
	return v.mean, true
}

// SampleVariance returns the sample variance (M2 / (n−1), with Bessel's
// correction), matching the contract of stats.SampleVariance: it returns
// ok == false for fewer than two elements, where sample variance is undefined.
func (v *RunningVariance) SampleVariance() (float64, bool) {
	if v.n < 2 {
		return 0, false
	}
	return v.m2 / float64(v.n-1), true
}

// PopulationVariance returns the population variance (M2 / n, no Bessel's
// correction), matching the contract of stats.PopulationVariance: it returns
// ok == false for empty input, and a defined variance of 0 for a single element.
func (v *RunningVariance) PopulationVariance() (float64, bool) {
	if v.n < 1 {
		return 0, false
	}
	return v.m2 / float64(v.n), true
}

// Count returns the number of values added so far.
func (v *RunningVariance) Count() int {
	return v.n
}

// EWMA maintains an exponentially weighted moving average of a stream of float64
// values in O(1) memory. Each new value is blended into the average with weight
// alpha, so recent values dominate and older ones decay geometrically:
// result = alpha*x + (1−alpha)*result. A larger alpha tracks the stream more
// responsively; a smaller alpha smooths it more heavily.
//
// The average is primed on the first Add (it is set to that first value rather
// than blended into a zero), so the series is not biased toward zero at the
// start. EWMA is single-threaded; see the package documentation on thread
// safety.
type EWMA struct {
	alpha  float64
	value  float64
	primed bool
}

// NewEWMA creates an EWMA with the given smoothing factor. alpha is clamped into
// the half-open interval (0, 1]: a value <= 0 (or NaN) becomes a small positive
// default and a value > 1 is capped at 1, so the average is always well-defined.
// Before the first Add, Result reports ok == false.
func NewEWMA(alpha float64) *EWMA {
	return &EWMA{alpha: clampAlpha(alpha)}
}

// clampAlpha forces alpha into (0, 1]. Non-positive or NaN values (NaN fails
// every comparison, so it falls through to the default) become a small positive
// default; values above 1 are capped at 1.
func clampAlpha(alpha float64) float64 {
	if alpha > 1 {
		return 1
	}
	if alpha > 0 {
		return alpha
	}
	return 0.5
}

// Add folds one value into the moving average in O(1). The first Add primes the
// average to x; subsequent Adds blend x in with weight alpha.
func (e *EWMA) Add(x float64) {
	if !e.primed {
		e.value = x
		e.primed = true
		return
	}
	e.value += e.alpha * (x - e.value)
}

// Result returns the current moving average. The second return value is false —
// and the average 0 — before the first Add.
func (e *EWMA) Result() (float64, bool) {
	if !e.primed {
		return 0, false
	}
	return e.value, true
}

// RunningMinMax tracks the smallest and largest element of an unbounded stream
// of an ordered type in O(1) memory, updated one element at a time. It mirrors
// stats.MinMax over a slice: the same stream of data yields the same lo/hi pair.
//
// As with stats.MinMax, NaN has no defined ordering, so a NaN-contaminated float
// stream yields an unspecified min/max. RunningMinMax is single-threaded; see
// the package documentation on thread safety.
type RunningMinMax[T constraints.Ordered] struct {
	lo    T
	hi    T
	empty bool
}

// NewRunningMinMax creates an empty RunningMinMax. Before any Add, Result reports
// ok == false.
func NewRunningMinMax[T constraints.Ordered]() *RunningMinMax[T] {
	return &RunningMinMax[T]{empty: true}
}

// Add folds one element into the running extremes in O(1).
func (mm *RunningMinMax[T]) Add(x T) {
	if mm.empty {
		mm.lo, mm.hi = x, x
		mm.empty = false
		return
	}
	if x < mm.lo {
		mm.lo = x
	}
	if x > mm.hi {
		mm.hi = x
	}
}

// Result returns the smallest and largest elements seen so far, mirroring
// stats.MinMax. The third return value is false — and both results the zero
// value of T — when no element has been added yet.
func (mm *RunningMinMax[T]) Result() (lo, hi T, ok bool) {
	if mm.empty {
		var zero T
		return zero, zero, false
	}
	return mm.lo, mm.hi, true
}
