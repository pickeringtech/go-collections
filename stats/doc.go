// Package stats turns slices of numbers into statistical summaries — sums and
// basic reductions, arithmetic/weighted/specialised means, the quantile family,
// covariance and correlation today, with the wider numeric surface (variance, …)
// landing
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
// The float64 summaries accumulate with Kahan compensated summation, and
// variance, covariance and correlation use Welford's online algorithm, so large
// or near-constant inputs do not lose precision to naive floating-point
// round-off. The exact-in-T reduction Sum instead accumulates in the input type
// T, so an integer sum is exact (no float round-off) at the cost of possible
// overflow on very large inputs.
//
// function. The exact-in-T arithmetic reductions (Sum, Product and the running
// CumulativeSum) let them propagate per IEEE arithmetic — a NaN in the data
// yields a NaN result (ok == true where there is an ok flag) — as do the
// variance/covariance/correlation family and the transforms, so the value
// surfaces rather than being silently dropped. The float64 summaries (the means
// and the quantile family, including Median) reject them (ok == false), since
// the statistic would be undefined; Range and Mode reject them too, because a
// min/max spread and a frequency count are undefined once ordering or equality
// breaks down on non-finite data. The ordering reductions MinMax/ArgMin/ArgMax
// work over constraints.Ordered (which also accepts strings, where NaN has no
// meaning), so like the standard library they leave a NaN-contaminated float
// result unspecified; and ClampAll, returning a bare slice with no ok flag,
// lets non-finite values flow through per IEEE-754.
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
//
// # Advanced
//
// Beyond the everyday summaries the package offers a tier of advanced
// operations, each keeping the same conventions (return type, (result, ok)
// idiom, numerical stability and the family-appropriate NaN/Inf policy):
//
//   - LinearRegression fits an ordinary-least-squares line (Slope/Intercept/R²,
//     with a Predict method for fitted values and residuals).
//   - Histogram buckets a sample into equal-width Bins over its [min, max] range.
//   - Skewness and Kurtosis report the distribution's asymmetry and tail weight
//     from stably-accumulated higher central moments (Kurtosis is excess, so a
//     normal distribution is 0).
//   - Entropy (Shannon, in bits) and Gini (impurity) summarise the distribution
//     of a categorical sample over any comparable type, not just numbers.
//   - PercentileOfScore is the inverse of Percentile: the percentage of values
//     at or below a given score.
//   - Dot, Norm, EuclideanDistance and CosineSimilarity provide the vector
//     operations common to ML-adjacent work. These are the module's canonical
//     vector geometry: the ml/distance, ml/similarity and clustering metrics
//     delegate here rather than carrying their own copies.
//
// As elsewhere, the moment-based statistics (regression, skewness, kurtosis) and
// the vector operations let non-finite values propagate (ok == true with a
// non-finite result), while the categorical measures (Entropy, Gini) and
// PercentileOfScore reject them (ok == false), matching the package's family
// split. Regression, correlation, skewness, kurtosis and cosine similarity are
// additionally undefined for a constant/zero input (zero variance or a zero
// vector), and report ok == false.
package stats
