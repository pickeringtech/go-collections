// Package ml is the umbrella for go-collections' machine-learning and
// data-engineering utilities — pure, dependency-light functions that compose
// with the stats package. It holds no code of its own; the building blocks live
// in focused sub-packages, grouped by concern.
//
// # Taxonomy
//
// One family scores how far apart two vectors are:
//
//   - ml/distance — distance metrics where lower means closer: Euclidean,
//     Manhattan, Minkowski, Cosine, Hamming and Levenshtein.
//   - ml/similarity — similarity metrics where higher means more alike: cosine
//     similarity and the set-overlap measures Jaccard, Dice and Overlap.
//
// A second family is evaluation metrics, under ml/metrics, one package per
// problem type:
//
//   - ml/metrics/regression — MSE, RMSE, MAE, MAPE and R² for continuous targets.
//   - ml/metrics/classification — accuracy, the confusion matrix,
//     precision/recall/F1 (macro/micro/weighted), and the probabilistic ROC/AUC
//     and log-loss for discrete labels.
//   - ml/metrics/clustering — the silhouette coefficient, scoring a clustering
//     without ground-truth labels.
//   - ml/metrics/ranking — DCG/NDCG and (mean) average precision for ordered
//     result lists.
//
// Further ml/<concern> families (for example preprocessing) land here pre-1.0
// alongside the metrics.
//
// # Conventions
//
// Every function in this family is a pure function over slices that returns its
// result with an ok flag in the library's (result, ok) idiom rather than
// panicking or returning an error: ok is false for input that cannot be
// summarised (empty, mismatched lengths, or a mathematically undefined
// degenerate case), and each function's doc spells out its exact rejection
// policy. Numeric summaries are float64, inputs are never mutated, and the
// reductions route through the stats package rather than reimplementing
// summation. See each sub-package's doc for details.
package ml
