package regression

import (
	"math"

	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/stats"
)

// MeanSquaredError returns the mean squared error (MSE) between the true values
// yTrue and the predicted values yPred — the average of the squared residuals
// (yTrueᵢ − yPredᵢ)² — together with an ok flag.
//
// ok is false (and the result is 0) when the inputs cannot be summarised:
//   - yTrue is empty, or len(yTrue) != len(yPred);
//   - any residual is non-finite (NaN or ±Inf), for which the mean would be
//     undefined.
//
// The squared residuals are reduced with stats.Mean, so they accumulate with
// Kahan compensated summation and inherit the means family's reject-non-finite
// policy.
func MeanSquaredError[T constraints.Numeric](yTrue, yPred []T) (float64, bool) {
	if len(yTrue) == 0 || len(yTrue) != len(yPred) {
		return 0, false
	}

	squared := make([]float64, len(yTrue))
	for i := range yTrue {
		d := float64(yTrue[i]) - float64(yPred[i])
		squared[i] = d * d
	}
	return stats.Mean(squared)
}

// RootMeanSquaredError returns the root mean squared error (RMSE) — the square
// root of MeanSquaredError — together with an ok flag. RMSE is in the same
// units as the inputs, which often makes it easier to interpret than MSE.
//
// It mirrors MeanSquaredError's contract: ok is false for empty or
// unequal-length input and when any residual is non-finite.
func RootMeanSquaredError[T constraints.Numeric](yTrue, yPred []T) (float64, bool) {
	mse, ok := MeanSquaredError(yTrue, yPred)
	if !ok {
		return 0, false
	}
	return math.Sqrt(mse), true
}

// MeanAbsoluteError returns the mean absolute error (MAE) between yTrue and
// yPred — the average of the absolute residuals |yTrueᵢ − yPredᵢ| — together
// with an ok flag. Unlike MSE, MAE weights every error linearly, so it is less
// sensitive to outliers.
//
// ok is false (and the result is 0) when yTrue is empty, len(yTrue) !=
// len(yPred), or any residual is non-finite (NaN or ±Inf). The absolute
// residuals are reduced with stats.Mean.
func MeanAbsoluteError[T constraints.Numeric](yTrue, yPred []T) (float64, bool) {
	if len(yTrue) == 0 || len(yTrue) != len(yPred) {
		return 0, false
	}

	absolute := make([]float64, len(yTrue))
	for i := range yTrue {
		absolute[i] = math.Abs(float64(yTrue[i]) - float64(yPred[i]))
	}
	return stats.Mean(absolute)
}

// MeanAbsolutePercentageError returns the mean absolute percentage error
// (MAPE) between yTrue and yPred — the average of |（yTrueᵢ − yPredᵢ) / yTrueᵢ|
// — together with an ok flag. The result is a fraction: multiply by 100 for a
// percentage (a MAPE of 0.2 means the predictions are off by 20% on average).
//
// MAPE is undefined when any true value is zero (it would divide by zero), so
// ok is false (and the result is 0) in that case, as it is for empty or
// unequal-length input and for any non-finite term.
func MeanAbsolutePercentageError[T constraints.Numeric](yTrue, yPred []T) (float64, bool) {
	if len(yTrue) == 0 || len(yTrue) != len(yPred) {
		return 0, false
	}

	percentage := make([]float64, len(yTrue))
	for i := range yTrue {
		actual := float64(yTrue[i])
		if actual == 0 {
			return 0, false
		}
		percentage[i] = math.Abs((actual - float64(yPred[i])) / actual)
	}
	return stats.Mean(percentage)
}

// RSquared returns the coefficient of determination (R²) of a set of
// predictions — 1 − SS_res/SS_tot, where SS_res is the sum of squared
// residuals and SS_tot is the total sum of squares of yTrue about its mean —
// together with an ok flag. R² is 1 for a perfect fit, 0 for a model no better
// than predicting the mean, and negative for a model that does worse than that.
//
// It is computed as 1 − MeanSquaredError / stats.PopulationVariance(yTrue),
// since both share the same 1/n factor.
//
// ok is false (and the result is 0) when the inputs cannot be summarised:
//   - yTrue is empty, or len(yTrue) != len(yPred);
//   - any residual is non-finite (NaN or ±Inf);
//   - yTrue has zero variance (every true value is identical), for which R² is
//     undefined because SS_tot is zero.
func RSquared[T constraints.Numeric](yTrue, yPred []T) (float64, bool) {
	mse, ok := MeanSquaredError(yTrue, yPred)
	if !ok {
		return 0, false
	}

	variance, ok := stats.PopulationVariance(yTrue)
	if !ok || variance == 0 {
		return 0, false
	}
	return 1 - mse/variance, true
}
