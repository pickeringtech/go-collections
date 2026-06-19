package preprocessing

import (
	"sort"

	"github.com/pickeringtech/go-collections/stats"
)

// binIndex returns the index of the bin that value falls into, given the
// internal cut points edges (length nBins-1, ascending). Bins are left-closed,
// right-open: a value equal to a cut point falls into the upper bin. Values
// below the first cut point are bin 0; values at or above the last cut point are
// the top bin. The result is always in [0, len(edges)].
func binIndex(edges []float64, value float64) int {
	return sort.Search(len(edges), func(i int) bool { return value < edges[i] })
}

// FixedWidthBinner discretizes values into nBins equal-width bins spanning the
// training range [min, max]. The cut points are learned at Fit time (by
// arithmetic over the learned min/max) and frozen for every subsequent
// Transform. Test values outside the training range fall into the first or last
// bin. Non-finite handling follows stats.MinMax over ordered floats.
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
// below 1 or data is empty. A degenerate range (min == max) yields a single
// populated bin.
func (b *FixedWidthBinner) Fit(data []float64) *FixedWidthBinner {
	if b.nBins < 1 {
		return b
	}
	lo, hi, ok := stats.MinMax(data)
	if !ok {
		return b
	}

	span := hi - lo
	if span <= 0 {
		// Degenerate range: every value collapses into a single bin.
		b.edges = []float64{}
		b.fitted = true
		return b
	}
	width := span / float64(b.nBins)
	edges := make([]float64, 0, b.nBins-1)
	for i := 1; i < b.nBins; i++ {
		edges = append(edges, lo+width*float64(i))
	}
	b.edges = edges
	b.fitted = true
	return b
}

// Transform maps each value in input to its bin index, returning an []int (one
// index per value, in [0, nBins)) and an ok flag (false if unfitted).
func (b *FixedWidthBinner) Transform(input []float64) ([]int, bool) {
	if !b.fitted {
		return nil, false
	}
	out := make([]int, len(input))
	for i, v := range input {
		out[i] = binIndex(b.edges, v)
	}
	return out, true
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
// Following the quantile family, Fit REJECTS non-finite training data, leaving
// the binner unfitted.
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
func (b *QuantileBinner) Fit(data []float64) *QuantileBinner {
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
// index per value, in [0, nBins)) and an ok flag (false if unfitted).
func (b *QuantileBinner) Transform(input []float64) ([]int, bool) {
	if !b.fitted {
		return nil, false
	}
	out := make([]int, len(input))
	for i, v := range input {
		out[i] = binIndex(b.edges, v)
	}
	return out, true
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
