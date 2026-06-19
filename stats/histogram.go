package stats

import "github.com/pickeringtech/go-collections/constraints"

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
// buckets in ascending order together with an
// ok flag. The bucket widths are (max−min)/bins; each value falls into
// ⌊(x−min)/width⌋, with the maximum value placed in the last bin (whose upper
// bound is inclusive). The returned counts therefore always sum to len(input).
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

	lo := float64(input[0])
	hi := lo
	for _, v := range input {
		f := float64(v)
		if nonFinite(f) {
			return nil, false
		}
		if f < lo {
			lo = f
		}
		if f > hi {
			hi = f
		}
	}

	span := hi - lo
	if span == 0 {
		// Every value is identical: a zero-width range gives no boundaries to
		// bucket against, so the histogram is undefined.
		return nil, false
	}

	width := span / float64(bins)
	out := make([]Bin, bins)
	for i := range out {
		out[i] = Bin{Min: lo + float64(i)*width, Max: lo + float64(i+1)*width}
	}
	// Pin the last bin's upper bound to the exact maximum so floating-point
	// drift in lo + bins*width cannot leave it just shy of (and excluding) hi.
	out[bins-1].Max = hi

	for _, v := range input {
		idx := int((float64(v) - lo) / width)
		if idx >= bins {
			// The maximum value lands at idx == bins; fold it into the last
			// (inclusive) bin rather than overflowing the slice.
			idx = bins - 1
		}
		out[idx].Count++
	}
	return out, true
}
