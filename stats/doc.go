// Package stats turns slices of numbers into statistical summaries —
// weighted and specialised means today, with the wider numeric surface
// (variance, quantiles, correlation, …) landing alongside it pre-1.0.
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
// # Conventions
//
// Empty input is undefined, so every function returns (float64, bool) in the
// library's (result, ok) idiom rather than a silent zero. The ok flag is false
// for empty input and for input the function cannot summarise (see each
// function's doc for its specific rejection policy).
//
// Sums are accumulated with Kahan compensated summation so large inputs do not
// lose precision to naive floating-point round-off.
//
// Non-finite inputs (NaN, ±Inf) are rejected: any function that encounters one
// returns ok=false, because the resulting statistic would be undefined.
package stats
