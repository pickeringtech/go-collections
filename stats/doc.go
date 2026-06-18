// Package stats is the home for summarising collections of values into
// statistics. It owns the implementations; thin accessors elsewhere (for
// example slices.NumericSlice) only delegate here, so there is a single source
// of truth for each operation.
//
// Operations whose result is naturally a real number return float64 (Median,
// Mean, …); operations that summarise without leaving the input's domain are
// exact in T (Product, Range, CumulativeSum, MinMax). The ordering reductions
// MinMax/ArgMin/ArgMax accept any constraints.Ordered type, so they work on
// strings as well as numbers.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/stats"
//
//	xs := []int{3, 1, 4, 1, 5}
//
//	prod, _ := stats.Product(xs)       // 60, true
//	med, _ := stats.Median(xs)         // 3, true   (float64)
//	rng, _ := stats.Range(xs)          // 4, true   (5 - 1)
//	modes, _ := stats.Mode(xs)         // []int{1}, true
//	lo, hi, _ := stats.MinMax(xs)      // 1, 5, true
//	cum := stats.CumulativeSum(xs)     // []int{3, 4, 8, 9, 14}
//	clamped := stats.Clamp(7, 0, 5)    // 5
//
//	_, _, _, _, _, _, _, _ = prod, med, rng, modes, lo, hi, cum, clamped
//
// This Quick Start is compiled and run as Example_quickStart in the package's
// test suite, so it is guaranteed to track the real API.
//
// # Conventions
//
// Empty input is undefined, so functions that summarise into a single value
// return (result, bool) in the library's (result, ok) idiom rather than a
// silent zero. The ok flag is false for empty input and for input the function
// cannot summarise (see each function's doc for its specific rejection policy).
//
// Non-finite inputs (NaN, ±Inf) make an ordering or aggregate statistic
// undefined, so every (value, ok) function rejects them by returning ok=false.
// The two functions that return a slice rather than a single value
// (CumulativeSum and ClampAll) have no ok flag to carry that signal, so they
// instead let non-finite values flow through per IEEE-754 and document it.
//
// Return types are exact in T for the type-preserving reductions (Product,
// Range, CumulativeSum, MinMax) and for the comparable Mode; Median returns
// float64 because the midpoint of an even-length sample need not be
// representable in T. The ordering reductions (MinMax, ArgMin, ArgMax) operate
// on constraints.Ordered, which has no NaN concept for strings, so they follow
// the standard library and leave a NaN-contaminated float result unspecified
// rather than rejecting it.
package stats
