# ml/metrics — Evaluation Metrics

Evaluation metrics as pure functions over slices, one package per problem type.
Part of the [`ml`](../) umbrella; composes with [`stats`](../../stats).

| Package                                  | Scores                  | Headline metrics                                              |
| ---------------------------------------- | ----------------------- | ------------------------------------------------------------ |
| **[regression](./regression/)**          | continuous targets      | MSE, RMSE, MAE, MAPE, R²                                      |
| **[classification](./classification/)**  | discrete labels         | accuracy, confusion matrix, precision/recall/F1, ROC/AUC, log-loss |
| **[clustering](./clustering/)**          | clusterings (no labels) | silhouette score                                             |
| **[ranking](./ranking/)**                | ordered result lists    | DCG, NDCG, (mean) average precision                          |

All functions return `(result, ok)`; `ok` is `false` for empty, mismatched, or
degenerate input. See each package's README and godoc for the exact contracts.
