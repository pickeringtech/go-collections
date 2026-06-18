// Package stats turns slices of numbers into statistical summaries — sums,
// arithmetic/weighted/specialised means, the quantile family, covariance and
// correlation today, with the wider numeric surface (variance, …) landing
// alongside it pre-1.0. It is also the home for value-rescaling transforms such
// as normalization and standardization.
//
// It is the home for "summarise numbers into a statistic" operations, which
// almost always return float64 (transforms return a fresh []float64). The
// sibling slices package keeps operations about slice structure and element
// ordering (Min/Max/Sort); stats does not duplicate those.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/stats"
//
//	data := []float64{10, 20, 30}
//
//	// Total and arithmetic mean — each with an ok flag (false for empty input).
//	total, _ := stats.Sum(data)  // 60, true
//	mean, _ := stats.Mean(data)  // 20, true
//
//	prices := []float64{10, 20, 30}
//	weights := []float64{1, 2, 3}
//
//	// Weighted mean — the result and an ok flag (false for empty/invalid input).
//	wm, ok := stats.WeightedMean(prices, weights) // 23.33..., true
//
//	// Specialised means for ratios and rates.
//	gm, _ := stats.GeometricMean([]float64{1, 10, 100}) // 10
//	hm, _ := stats.HarmonicMean([]float64{1, 2, 4})      // 1.714...
//
//	_ = total
//	_ = mean
//	_ = wm
//	_ = gm
//	_ = hm
//
// This Quick Start is compiled and run as Example_quickStart in the package's
// test suite, so it is guaranteed to track the real API.
//
// # Conventions
//
// Empty input is undefined, so every function returns an ok flag in the
// library's (result, ok) idiom rather than a silent zero. The ok flag is false
// for empty input and for input the function cannot summarise (see each
// function's doc for its specific rejection policy). Transforms (Normalize,
// Standardize, MovingAverage) follow the same idiom, returning ([]float64, bool)
// and never mutating their input.
//
// The float64 summaries accumulate with Kahan compensated summation, and
// variance, covariance and correlation use Welford's online algorithm, so large
// or near-constant inputs do not lose precision to naive floating-point
// round-off. The exact-in-T reduction Sum instead accumulates in the input type
// T, so an integer sum is exact (no float round-off) at the cost of possible
// overflow on very large inputs.
//
// Non-finite inputs (NaN, ±Inf) are handled per operation, documented on each
// function. The means and the quantile family reject them (ok == false), since
// the resulting statistic would be undefined. The variance/covariance/
// correlation family, the transforms, and the exact-in-T Sum instead let them
// propagate — the result is non-finite with ok == true — so a NaN in the data
// surfaces as a NaN statistic rather than a plausible-looking wrong number,
// never silently dropped.
//
// Where Bessel's correction applies (variance, standard deviation, covariance)
// both sample and population variants are offered, named unambiguously so the
// choice is always the caller's.
//
// # Quantiles
//
// Quantile/Percentile/Quartiles/IQR interpolate between samples when the
// requested rank falls between two values. The default everywhere is Linear
// ("type 7" in Hyndman & Fan), which matches numpy.percentile's default — the
// convention most users expect. QuantileWith/PercentileWith accept an explicit
// InterpolationMethod (Linear, Lower, Higher, Nearest, Midpoint). These
// functions sort a copy of the input, so the caller's slice is never mutated.
package stats
