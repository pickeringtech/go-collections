package stats

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
)

// Bin is one bucket of a Histogram: the half-open interval [Min, Max) it covers
// and the Count of input values that fell into it. The single exception is the
// final bin of a Histogram, whose upper bound is inclusive ([Min, Max]) so that
// the largest input value is counted rather than dropped.
type Bin struct {
	Min, Max float64
	Count    int
}

// Histogram buckets input into bins equal-width buckets spanning the data's
// [min, max] range (so bins is the requested bucket count), returning the
// buckets in ascending order together with an ok flag.
//
// The edges are min + i·(max−min)/bins for i = 0 … bins, with the last pinned
// to exactly max. Each value is counted in the bin whose reported interval
// contains it — [Min, Max) for every bin but the last, which is [Min, Max] —
// so the counts always agree with the edges a caller sees, and they always sum
// to len(input). A range whose width max−min overflows float64 (for example
// {−math.MaxFloat64, math.MaxFloat64}) is still bucketed: the edges are then
// computed in half-units, which is exact at that magnitude.
//
// It returns ok == false when a histogram is undefined: when input is empty,
// when bins < 1, or when the data cannot be ordered into a finite range —
// either because a value is non-finite (NaN/±Inf) or because every value is
// identical (a zero-width range has no meaningful bin boundaries). This matches
// Range and Mode, which likewise reject non-finite data because a min/max
// spread is undefined once ordering breaks down; clean such values beforehand
// if you need to bucket them.
func Histogram[T constraints.Numeric](input []T, bins int) ([]Bin, bool) {
	if len(input) == 0 || bins < 1 {
		return nil, false
	}
	lo, hi, ok := finiteSpan(input)
	if !ok {
		return nil, false
	}

	// hi − lo can only overflow for values of opposite sign near ±MaxFloat64;
	// halving such values is exact, so work in half-units there to keep the
	// width and every edge finite. For any other data scale is 1 and the edges
	// are exactly lo + i·width.
	scale := 1.0
	if math.IsInf(hi-lo, 1) {
		scale = 0.5
	}
	width := (hi*scale - lo*scale) / float64(bins)

	out := make([]Bin, bins)
	for i := range out {
		out[i].Min = (lo*scale + float64(i)*width) / scale
		if i > 0 {
			out[i-1].Max = out[i].Min
		}
	}
	// Pin the last bin's upper bound to the exact maximum so floating-point
	// drift in lo + bins*width cannot leave it just shy of (and excluding) hi.
	out[bins-1].Max = hi

	for _, v := range input {
		x := float64(v)
		out[locateBin(out, x, (x*scale-lo*scale)/width)].Count++
	}
	return out, true
}

// finiteSpan returns the minimum and maximum of input, with ok == false when a
// value is non-finite or every value is identical (a zero-width range).
func finiteSpan[T constraints.Numeric](input []T) (lo, hi float64, ok bool) {
	lo = float64(input[0])
	hi = lo
	for _, v := range input {
		f := float64(v)
		if nonFinite(f) {
			return 0, 0, false
		}
		lo = min(lo, f)
		hi = max(hi, f)
	}
	return lo, hi, lo != hi
}

// locateBin returns the index of the bin whose reported interval holds x —
// [Min, Max) for every bin but the last, which is [Min, Max] — starting from
// the arithmetic estimate pos = (x−min)/width. That estimate and the reported
// edges (min + i·width) round independently, so the estimate can land one bin
// off; it is then stepped onto the bin that actually contains x, which keeps
// the counts consistent with the edges.
func locateBin(bins []Bin, x, pos float64) int {
	last := len(bins) - 1
	i := 0
	if pos > 0 {
		i = int(min(pos, float64(last)))
	}
	for i > 0 && x < bins[i].Min {
		i--
	}
	for i < last && x >= bins[i].Max {
		i++
	}
	return i
}
