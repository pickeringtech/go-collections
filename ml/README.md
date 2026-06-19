# ml — Machine-Learning Utilities

The `ml` umbrella holds go-collections' machine-learning and data-engineering
utilities: pure, dependency-light functions over slices that compose with the
[`stats`](../stats) package. The `ml` package itself holds no code — the
building blocks live in focused sub-packages, grouped by concern.

## Metrics

The first family is evaluation metrics under [`ml/metrics`](./metrics/), one
package per problem type:

| Package                                     | Scores                | Headline metrics                                              |
| ------------------------------------------- | --------------------- | ------------------------------------------------------------ |
| **[regression](./metrics/regression/)**     | continuous targets    | MSE, RMSE, MAE, MAPE, R²                                      |
| **[classification](./metrics/classification/)** | discrete labels   | accuracy, confusion matrix, precision/recall/F1, ROC/AUC, log-loss |
| **[clustering](./metrics/clustering/)**     | clusterings (no labels) | silhouette score                                           |
| **[ranking](./metrics/ranking/)**           | ordered result lists  | DCG, NDCG, (mean) average precision                          |

## Conventions

Every function is a **pure function over slices** that returns `(result, ok)`
rather than panicking or returning an `error`. `ok` is `false` for input that
cannot be summarised — empty, mismatched lengths, or a mathematically undefined
degenerate case — and each function's doc states its exact rejection policy.
Numeric summaries are `float64`, inputs are never mutated, and the reductions
route through [`stats`](../stats) rather than reimplementing summation.
