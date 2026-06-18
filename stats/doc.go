// Package stats turns slices of numbers into statistical summaries — sums,
// arithmetic/weighted/specialised means and the quantile family today, with the
// wider numeric surface (variance, correlation, …) landing alongside it pre-1.0.
//
// It is the home for "summarise numbers into a statistic" operations, which
// almost always return float64. The sibling slices package keeps operations
// about slice structure and element ordering (Min/Max/Sort); stats does not
// duplicate those.
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
// Empty input is undefined, so every function returns (float64, bool) in the
// library's (result, ok) idiom rather than a silent zero. The ok flag is false
// for empty input and for input the function cannot summarise (see each
// function's doc for its specific rejection policy).
//
// Operations fall into two tiers with deliberately different policies:
//
//   - float64 summaries (Mean, WeightedMean, the specialised means, …) sum with
//     Kahan compensated summation so large inputs do not lose precision to naive
//     round-off, and they reject non-finite inputs (NaN, ±Inf) with ok=false
//     because the resulting statistic would be undefined.
//   - exact-in-T reductions (Sum) accumulate in the input type T, so integer
//     results are exact (no float round-off) at the cost of possible overflow on
//     very large inputs, and non-finite float values propagate through the total
//     per IEEE arithmetic rather than being rejected. Each such function
//     documents this in its own doc comment.
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
