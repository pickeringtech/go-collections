// Package regression scores continuous-valued predictions against their true
// targets — the standard regression error metrics as pure functions over a
// pair of equal-length slices.
//
// It is part of the ml/metrics family (see the ml umbrella package). The
// sibling classification package scores discrete labels; this package scores
// real numbers.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/ml/metrics/regression"
//
//	yTrue := []float64{3, -0.5, 2, 7}
//	yPred := []float64{2.5, 0.0, 2, 8}
//
//	mse, _ := regression.MeanSquaredError(yTrue, yPred)      // 0.375
//	rmse, _ := regression.RootMeanSquaredError(yTrue, yPred) // 0.6124...
//	mae, _ := regression.MeanAbsoluteError(yTrue, yPred)     // 0.5
//	r2, _ := regression.RSquared(yTrue, yPred)               // 0.9486...
//
//	_ = mse
//	_ = rmse
//	_ = mae
//	_ = r2
//
// This Quick Start is compiled and run as Example_quickStart in the package's
// test suite, so it is guaranteed to track the real API.
//
// # Metrics
//
//   - MeanSquaredError (MSE) and RootMeanSquaredError (RMSE) — squared-error
//     loss; RMSE is in the inputs' units. Both penalise large errors heavily.
//   - MeanAbsoluteError (MAE) — linear loss, robust to outliers.
//   - MeanAbsolutePercentageError (MAPE) — scale-free relative error, returned
//     as a fraction (×100 for a percentage).
//   - RSquared (R²) — the coefficient of determination, the fraction of target
//     variance the predictions explain.
//
// # Conventions
//
// Every function takes (yTrue, yPred) of any constraints.Numeric type and
// returns (float64, bool) in the library's (result, ok) idiom rather than
// panicking or returning an error. ok is false — and the result is 0 — when the
// inputs cannot be summarised: yTrue is empty, the two slices differ in length,
// or a term is non-finite (NaN/±Inf). MAPE additionally rejects input where any
// true value is zero, and RSquared additionally rejects input where yTrue has
// zero variance; both cases are mathematically undefined. Inputs are never
// mutated.
//
// The error reductions route through the stats package (stats.Mean,
// stats.PopulationVariance) rather than reimplementing summation, so they
// inherit its Kahan compensated summation and its means-family policy of
// rejecting non-finite values.
package regression
