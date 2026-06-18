// Package stats summarizes collections of numbers into statistics, correctly.
//
// It is the home for operations that reduce slices of numbers to descriptive
// figures — variance, standard deviation, covariance, correlation — and for
// value-rescaling transforms such as normalization and standardization. The
// companion slices package owns slice structure and element ordering
// (Min/Max/sorting); stats owns the numeric summaries. One operation lives in
// exactly one place.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/stats"
//
//	xs := []float64{1, 2, 3, 4, 5}
//	ys := []float64{2, 4, 6, 8, 10}
//
//	r, _ := stats.Correlation(xs, ys)       // 1 — perfectly linear
//	z, _ := stats.Standardize(xs)           // zero mean, unit variance
//	ma, _ := stats.MovingAverage(xs, 2)     // rolling means over full windows
//
//	fmt.Printf("r=%.1f z0=%.4f ma=%v", r, z[0], ma)
//	// r=1.0 z0=-1.4142 ma=[1.5 2.5 3.5 4.5]
//
// # Conventions
//
// These are the deliberate, locked-in conventions every operation in this
// package follows:
//
//   - Numerical stability. Variance, covariance and correlation use Welford's
//     online algorithm, never the naive Σxy − ΣxΣy/n, which loses catastrophic
//     precision on large or near-constant magnitudes.
//   - Return type. Scalar summaries return float64, paired with an ok bool.
//     Transforms that rescale a series return a new []float64, also paired with
//     an ok bool; the input is never mutated.
//   - Empty/edge contract. Statistics on undefined input return ok == false
//     rather than a silent zero. Sample variants are undefined for fewer than
//     two elements; population variants are undefined only for empty input.
//   - NaN/Inf policy. Non-finite inputs propagate: the result is non-finite and
//     ok == true. Values are never silently filtered out, so a NaN in the data
//     surfaces as a NaN statistic rather than a plausible-looking wrong number.
//   - Sample vs population. Both variants are offered where Bessel's correction
//     applies (variance, standard deviation, covariance), named unambiguously
//     so the choice is always the caller's.
package stats
