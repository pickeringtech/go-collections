// Package stats provides numeric summary statistics over slices of numbers.
//
// Functions are free functions taking a []T where T is constraints.Numeric.
// Each returns a result plus a bool reporting whether that result is
// meaningful, following the library's "(result, ok)" empty contract: ok is
// false when the input is empty and, for ordering-sensitive operations, when
// the input contains a NaN.
//
// # Quantiles
//
//	import "github.com/pickeringtech/go-collections/stats"
//
//	data := []float64{1, 2, 3, 4, 5}
//
//	med, _ := stats.Quantile(data, 0.5)   // 3 — the median
//	p90, _ := stats.Percentile(data, 90)  // 4.6
//	qs, _ := stats.Quartiles(data)        // {Q1:2, Q2:3, Q3:4}
//	iqr, _ := stats.IQR(data)             // 2
//
// # Interpolation
//
// When the requested rank falls between two samples the value is interpolated.
// The default everywhere is Linear ("type 7" in Hyndman & Fan), which matches
// numpy.percentile's default. QuantileWith and PercentileWith accept an
// explicit InterpolationMethod (Linear, Lower, Higher, Nearest, Midpoint) for
// callers who need a specific convention.
//
// # NaN policy
//
// A NaN has no defined ordering, so any quantile of a sample that contains one
// is undefined. Rather than return a silently-wrong number, the quantile
// functions report ok=false when the input contains a NaN. Integer inputs can
// never be NaN, so this only affects float inputs.
//
// # Ownership
//
// Inputs are never mutated. A function that needs sorted data copies the input
// first, consistent with the library's ownership-isolation direction.
package stats
