// Package stats turns slices of numbers into statistical summaries — basic
// reductions, weighted and specialised means, the quantile family, covariance
// and correlation today, with the wider numeric surface (variance, …) landing
// alongside it pre-1.0. It is also the home for value-rescaling transforms such
// as normalization and standardization.
//
// It is the home for "summarise numbers into a statistic" operations, which
// almost always return float64 (transforms return a fresh []float64). A few
// reductions stay exact in T because the result never leaves the input's domain
// (Product, Range, CumulativeSum, MinMax). The sibling slices package keeps
// operations about slice structure and element ordering (Sort).
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/stats"
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
//	_ = wm
//	_ = gm
//	_ = hm
//
// This Quick Start is compiled and run as Example_quickStart in the package's
// test suite, so it is guaranteed to track the real API.
//
// # Basic operations
//
// Alongside the means and quantiles, the package provides the everyday
// reductions: Product, Range and the running CumulativeSum (each exact in T);
// Median (float64; an even-length sample averages its two middle elements);
// Mode (the most frequent value or values, over any comparable type); the
// single-pass MinMax and the index reductions ArgMin/ArgMax; and Clamp /
// ClampAll, which constrain values to a closed interval. See each function's
// doc and runnable Example for the exact contract.
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
// Sums are accumulated with Kahan compensated summation, and variance,
// covariance and correlation use Welford's online algorithm, so large or
// near-constant inputs do not lose precision to naive floating-point round-off.
//
// Non-finite inputs (NaN, ±Inf) are handled per operation, documented on each
// function. The means, the quantile family, and the basic reductions that carry
// an ok flag (Product, Range, Median, Mode) reject them (ok == false), since the
// resulting statistic would be undefined. The variance/covariance/correlation
// family and the transforms instead let them propagate — the result is
// non-finite with ok == true — so a NaN in the data surfaces as a NaN statistic
// rather than a plausible-looking wrong number, never silently dropped. Two
// kinds of operation handle them differently again: the ordering reductions
// MinMax/ArgMin/ArgMax work over constraints.Ordered (which also accepts
// strings, where NaN has no meaning), so like the standard library they leave a
// NaN-contaminated float result unspecified; and the functions that return a
// bare slice with no ok flag (CumulativeSum, ClampAll) have no way to signal a
// rejection, so non-finite values flow through them per IEEE-754.
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
// Median is Quantile(input, 0.5).
package stats
