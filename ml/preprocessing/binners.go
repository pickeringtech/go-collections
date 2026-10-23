package preprocessing

import (
	"math"
	"sort"

	"github.com/pickeringtech/go-collections/stats"
)

// binIndex returns the index of the bin that value falls into, given the
// internal cut points edges (length nBins-1, ascending). Bins are left-closed,
// right-open: a value equal to a cut point falls into the upper bin. Values
// below the first cut point are bin 0; values at or above the last cut point are
// the top bin. The result is always in [0, len(edges)]. value must be finite: a
// NaN compares false against every edge and would land in the top bin.
func binIndex(edges []float64, value float64) int {
	return sort.Search(len(edges), func(i int) bool { return value < edges[i] })
}

// binAll is the binners' shared Transform: each value of input mapped to its
// bin, gated on fitted. ok is false when the binner is unfitted or any value is
// non-finite — NaN has no order and so no bin, and an infinity has no finite
// bin either — matching the reject policy the binners apply at Fit time.
func binAll(fitted bool, edges, input []float64) ([]int, bool) {
	if !fitted {
		return nil, false
	}
	out := make([]int, len(input))
	for i, v := range input {
		if !finite(v) {
			return nil, false
		}
		out[i] = binIndex(edges, v)
	}
	return out, true
}

// finite reports whether v is neither NaN nor ±Inf.
func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

// FixedWidthBinner discretizes values into nBins equal-width bins spanning the
// training range [min, max]. The cut points are learned at Fit time (by
// arithmetic over the learned min/max) and frozen for every subsequent
// Transform. Finite test values outside the training range fall into the first
// or last bin.
//
// Like QuantileBinner it REJECTS non-finite values: training data containing a
// NaN or ±Inf has no finite range and leaves the binner unfitted, and
// Transform reports ok == false for a NaN or ±Inf input value, which has no
// bin.
type FixedWidthBinner struct {
	nBins  int
	edges  []float64
	fitted bool
}

// NewFixedWidthBinner returns an unfitted FixedWidthBinner that will produce
// nBins bins. nBins must be at least 1; otherwise Fit leaves it unfitted.
func NewFixedWidthBinner(nBins int) *FixedWidthBinner {
	return &FixedWidthBinner{nBins: nBins}
}

// Fit learns the equal-width cut points from the range of data and returns the
// receiver so calls can be chained. It leaves the binner unfitted when nBins is
// below 1, data is empty, or any value is non-finite. A degenerate range
// (min == max) yields a single populated bin. A finite range whose width
// overflows float64 (such as {−math.MaxFloat64, math.MaxFloat64}) still yields
// finite cut points: they are computed in half-units, which is exact at that
// magnitude.
//
// Fitting always starts from scratch: a rejected Fit discards whatever an
// earlier Fit learned, leaving the binner unfitted (Transform reports
// ok == false) rather than transforming with stale parameters.
func (b *FixedWidthBinner) Fit(data []float64) *FixedWidthBinner {
	*b = FixedWidthBinner{nBins: b.nBins} // a rejected Fit leaves the binner unfitted
	if b.nBins < 1 {
		return b
	}
	for _, v := range data {
		if !finite(v) {
			return b
		}
	}
	lo, hi, ok := stats.MinMax(data)
	if !ok {
		return b
	}

	if hi == lo {
		// Degenerate range: every value collapses into a single bin.
		b.edges = []float64{}
		b.fitted = true
		return b
	}
	// hi − lo can only overflow for values of opposite sign near
	// ±MaxFloat64; halving such values is exact, so work in half-units there.
	scale := 1.0
	if math.IsInf(hi-lo, 1) {
		scale = 0.5
	}
	width := (hi*scale - lo*scale) / float64(b.nBins)
	edges := make([]float64, 0, b.nBins-1)
	for i := 1; i < b.nBins; i++ {
		edges = append(edges, (lo*scale+width*float64(i))/scale)
	}
	b.edges = edges
	b.fitted = true
	return b
}

// Transform maps each value in input to its bin index, returning an []int (one
// index per value, in [0, nBins)) and an ok flag — false if the binner is
// unfitted or any value is non-finite (NaN or ±Inf).
func (b *FixedWidthBinner) Transform(input []float64) ([]int, bool) {
	return binAll(b.fitted, b.edges, input)
}

// FitTransform fits on input and immediately transforms it. Equivalent to
// Fit(input).Transform(input).
func (b *FixedWidthBinner) FitTransform(input []float64) ([]int, bool) {
	return b.Fit(input).Transform(input)
}

// Edges returns a copy of the internal cut points (length nBins-1) learned at
// Fit time. It is meaningful only after a successful Fit.
func (b *FixedWidthBinner) Edges() []float64 {
	out := make([]float64, len(b.edges))
	copy(out, b.edges)
	return out
}

// QuantileBinner discretizes values into nBins bins of approximately equal
// population, with cut points placed at the i/nBins quantiles of the training
// data (via stats.Quantile). The cut points are learned at Fit time and frozen
// for every subsequent Transform.
//
// Following the quantile family, it REJECTS non-finite values: Fit leaves the
// binner unfitted when the training data contains a NaN or ±Inf, and Transform
// reports ok == false for a NaN or ±Inf input value, which has no bin.
type QuantileBinner struct {
	nBins  int
	edges  []float64
	fitted bool
}

// NewQuantileBinner returns an unfitted QuantileBinner that will produce nBins
// bins. nBins must be at least 1; otherwise Fit leaves it unfitted.
func NewQuantileBinner(nBins int) *QuantileBinner {
	return &QuantileBinner{nBins: nBins}
}

// Fit learns the quantile cut points from data and returns the receiver so calls
// can be chained. It leaves the binner unfitted when nBins is below 1, data is
// empty, or any value is non-finite.
//
// Fitting always starts from scratch: a rejected Fit discards whatever an
// earlier Fit learned, leaving the binner unfitted (Transform reports
// ok == false) rather than transforming with stale parameters.
func (b *QuantileBinner) Fit(data []float64) *QuantileBinner {
	*b = QuantileBinner{nBins: b.nBins} // a rejected Fit leaves the binner unfitted
	if b.nBins < 1 {
		return b
	}
	edges := make([]float64, 0, b.nBins-1)
	for i := 1; i < b.nBins; i++ {
		q := float64(i) / float64(b.nBins)
		edge, ok := stats.Quantile(data, q)
		if !ok {
			return b
		}
		edges = append(edges, edge)
	}
	// nBins == 1 needs no cut points but is still a valid (single-bin) fit;
	// guard the empty-data case the quantile loop would otherwise skip.
	if b.nBins == 1 {
		_, ok := stats.Quantile(data, 0)
		if !ok {
			return b
		}
	}
	b.edges = edges
	b.fitted = true
	return b
}

// Transform maps each value in input to its bin index, returning an []int (one
// index per value, in [0, nBins)) and an ok flag — false if the binner is
// unfitted or any value is non-finite (NaN or ±Inf).
func (b *QuantileBinner) Transform(input []float64) ([]int, bool) {
	return binAll(b.fitted, b.edges, input)
}

// FitTransform fits on input and immediately transforms it. Equivalent to
// Fit(input).Transform(input).
func (b *QuantileBinner) FitTransform(input []float64) ([]int, bool) {
	return b.Fit(input).Transform(input)
}

// Edges returns a copy of the internal quantile cut points (length nBins-1)
// learned at Fit time. It is meaningful only after a successful Fit.
func (b *QuantileBinner) Edges() []float64 {
	out := make([]float64, len(b.edges))
	copy(out, b.edges)
	return out
}
