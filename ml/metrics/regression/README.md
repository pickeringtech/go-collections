# Regression Metrics

The `regression` package scores continuous-valued predictions against their
true targets — the standard regression error metrics as pure functions over a
pair of equal-length slices. It is part of the [`ml/metrics`](../) family; the
sibling [`classification`](../classification) package scores discrete labels.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/ml/metrics/regression"

yTrue := []float64{3, -0.5, 2, 7}
yPred := []float64{2.5, 0, 2, 8}

mse, _ := regression.MeanSquaredError(yTrue, yPred)      // 0.375
rmse, _ := regression.RootMeanSquaredError(yTrue, yPred) // 0.6124
mae, _ := regression.MeanAbsoluteError(yTrue, yPred)     // 0.5
r2, _ := regression.RSquared(yTrue, yPred)               // 0.9486
```

## Metrics

| Function                              | Returns           | Undefined (`ok == false`) when                                |
| ------------------------------------- | ----------------- | ------------------------------------------------------------- |
| `MeanSquaredError(yTrue, yPred)`      | `(float64, bool)` | empty, differing lengths, or any non-finite term              |
| `RootMeanSquaredError(yTrue, yPred)`  | `(float64, bool)` | as MSE (it is `√MSE`)                                          |
| `MeanAbsoluteError(yTrue, yPred)`     | `(float64, bool)` | empty, differing lengths, or any non-finite term              |
| `MeanAbsolutePercentageError(…)`      | `(float64, bool)` | as above, **plus** any true value is zero                     |
| `RSquared(yTrue, yPred)`              | `(float64, bool)` | as MSE, **plus** `yTrue` has zero variance (constant target)  |

- **MSE / RMSE** penalise large errors heavily (squared loss). RMSE is in the
  inputs' units, which usually makes it easier to interpret.
- **MAE** weights every error linearly, so it is more robust to outliers.
- **MAPE** is scale-free and returned as a **fraction** — multiply by 100 for a
  percentage (a MAPE of `0.2` means predictions are off by 20% on average).
- **R²** is the coefficient of determination: 1 for a perfect fit, 0 for a model
  no better than predicting the mean, and negative for one that does worse.

## Conventions

Every function is generic over `constraints.Numeric` (so integer series work)
and returns `(float64, bool)` in the library's `(result, ok)` idiom rather than
panicking or returning an `error`. `ok` is `false` — and the result `0` — for
input that cannot be summarised. Inputs are never mutated.

The error reductions route through the [`stats`](../../../stats) package
(`stats.Mean`, `stats.PopulationVariance`) rather than reimplementing summation,
inheriting its Kahan compensated summation and its policy of **rejecting**
non-finite (NaN/±Inf) values.
