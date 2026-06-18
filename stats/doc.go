// Package stats turns slices of numbers into statistical summaries — basic
// reductions, weighted and specialised means, and the quantile family today,
// with the wider numeric surface (variance, correlation, …) landing alongside
// it pre-1.0.
//
// It is the home for "summarise numbers into a statistic" operations and owns
// their implementations; thin accessors elsewhere (for example
// slices.NumericSlice) only delegate here, so there is a single source of truth
// per operation. Operations whose result is naturally a real number return
// float64 (Mean, Median, …); operations that summarise without leaving the
// input's domain are exact in T (Product, Range, CumulativeSum, MinMax). The
// sibling slices package keeps operations about slice structure and element
// ordering (Sort).
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
// Empty input is undefined, so functions that summarise into a single value
// return (result, bool) in the library's (result, ok) idiom rather than a
// silent zero. The ok flag is false for empty input and for input the function
// cannot summarise (see each function's doc for its specific rejection policy).
//
// Sums are accumulated with Kahan compensated summation so large inputs do not
// lose precision to naive floating-point round-off.
//
// Non-finite inputs (NaN, ±Inf) are rejected: any (value, ok) function that
// encounters one returns ok=false, because the resulting statistic would be
// undefined. Two kinds of operation stand outside that rule. The ordering
// reductions MinMax/ArgMin/ArgMax operate on constraints.Ordered, which has no
// NaN concept for strings, so they follow the standard library and leave a
// NaN-contaminated float result unspecified rather than rejecting it. The two
// slice-returning functions (CumulativeSum and ClampAll) have no ok flag to
// carry a rejection, so non-finite values flow through them per IEEE-754.
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
