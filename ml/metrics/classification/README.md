# Classification Metrics

The `classification` package scores discrete-label predictions against their
true labels — accuracy, the confusion matrix, precision/recall/F1 with macro,
micro and weighted averaging, and the probabilistic metrics ROC/AUC and
log-loss. It is part of the [`ml/metrics`](../) family; the sibling
[`regression`](../regression) package scores continuous values.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/ml/metrics/classification"

yTrue := []int{0, 1, 2, 2, 1, 0, 1, 2}
yPred := []int{0, 2, 1, 2, 1, 0, 1, 2}

acc, _ := classification.Accuracy(yTrue, yPred)                 // 0.75
f1, _ := classification.F1(yTrue, yPred, classification.Macro)  // 0.7778
cm, _ := classification.ConfusionMatrix(yTrue, yPred)
hits := cm.Count(2, 2)                                          // 2

// Binary, score-based metrics.
labels := []int{0, 0, 1, 1}
scores := []float64{0.1, 0.4, 0.35, 0.8}
auc, _ := classification.AUC(labels, scores, 1)                 // 0.75
```

## Label metrics

Labels may be of any comparable type (compared with `==`).

| Function                                   | Returns             | Notes                                                   |
| ------------------------------------------ | ------------------- | ------------------------------------------------------- |
| `Accuracy(yTrue, yPred)`                   | `(float64, bool)`   | fraction of exact matches                               |
| `ConfusionMatrix(yTrue, yPred)`            | `(Matrix[T], bool)` | counts grid; `Labels()` and `Count(true, pred)` methods |
| `Precision/Recall/F1(yTrue, yPred, avg)`   | `(float64, bool)`   | multiclass, combined under an `Averaging`               |
| `PrecisionBinary/RecallBinary/F1Binary(…)` | `(float64, bool)`   | binary, against a designated positive label             |

`Averaging` selects how per-class scores combine:

- **`Macro`** — unweighted mean across classes (rare classes count fully).
- **`Micro`** — pool the counts first; for single-label data this equals the
  accuracy.
- **`Weighted`** — per-class scores weighted by each class's support.

Degenerate per-class cases are **defined, not rejected**: no predictions →
precision 0; no true samples → recall 0; precision + recall 0 → F1 0.

## Score metrics (binary)

The model emits a score/probability for the positive label rather than a hard
label.

| Function                          | Returns           | Undefined (`ok == false`) when                          |
| --------------------------------- | ----------------- | ------------------------------------------------------- |
| `ROCCurve(yTrue, scores, pos)`    | `(Curve, bool)`   | empty, differing lengths, non-finite score, single class |
| `AUC(yTrue, scores, pos)`         | `(float64, bool)` | as ROCCurve                                             |
| `LogLoss(yTrue, probs, pos)`      | `(float64, bool)` | empty, differing lengths, any prob NaN or ∉ [0, 1]      |

`AUC` uses tie-averaged ranks (the Mann–Whitney U statistic) — exact, and free
of the trapezoid error of integrating the curve. `LogLoss` clamps probabilities
to `[1e-15, 1−1e-15]`, so a confident mistake is a large finite penalty rather
than `+Inf`. (Multiclass log-loss over a probability matrix is not yet provided.)

## Conventions

Every function returns `(result, ok)` rather than panicking or returning an
`error`; `ok` is `false` (result zero) for input that cannot be summarised.
Inputs are never mutated. The log-loss reduction routes through
[`stats.Mean`](../../../stats).
