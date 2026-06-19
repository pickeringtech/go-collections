// Package classification scores discrete-label predictions against their true
// labels — accuracy, the confusion matrix, precision/recall/F1 with macro,
// micro and weighted averaging, and the probabilistic metrics ROC/AUC and
// log-loss — as pure functions over slices.
//
// It is part of the ml/metrics family (see the ml umbrella package). The
// sibling regression package scores continuous values; this package scores
// labels of any comparable type.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/ml/metrics/classification"
//
//	yTrue := []int{0, 1, 2, 2, 1, 0, 1, 2}
//	yPred := []int{0, 2, 1, 2, 1, 0, 1, 2}
//
//	acc, _ := classification.Accuracy(yTrue, yPred)               // 0.75
//	f1, _ := classification.F1(yTrue, yPred, classification.Macro) // 0.7778
//	cm, _ := classification.ConfusionMatrix(yTrue, yPred)
//	hits := cm.Count(2, 2)                                         // 2
//
//	// Binary, score-based metrics.
//	labels := []int{0, 0, 1, 1}
//	scores := []float64{0.1, 0.4, 0.35, 0.8}
//	auc, _ := classification.AUC(labels, scores, 1) // 0.75
//
//	_ = acc
//	_ = f1
//	_ = hits
//	_ = auc
//
// This Quick Start is compiled and run as Example_quickStart in the package's
// test suite, so it is guaranteed to track the real API.
//
// # Label metrics
//
// Accuracy, ConfusionMatrix, Precision, Recall and F1 take predicted labels of
// any comparable type (compared with ==). For multiclass problems Precision,
// Recall and F1 take an Averaging strategy:
//
//   - Macro — unweighted mean of the per-class scores (rare classes count as
//     much as common ones);
//   - Micro — pool the counts across classes first (for single-label data, all
//     three equal the accuracy);
//   - Weighted — per-class scores weighted by each class's support.
//
// PrecisionBinary, RecallBinary and F1Binary are the binary special cases: pass
// the label to treat as positive, and every other label is negative.
//
// # Score metrics
//
// ROCCurve, AUC and LogLoss operate on a binary problem where the model emits a
// score or probability for the positive label rather than a hard label. AUC is
// computed from tie-averaged ranks (the Mann–Whitney U statistic), and LogLoss
// clamps probabilities away from 0 and 1 so a confident mistake is a large
// finite penalty rather than +Inf.
//
// # Conventions
//
// Every function returns its result with an ok flag in the library's
// (result, ok) idiom rather than panicking or returning an error. ok is false —
// and the result the zero value — when the inputs cannot be summarised: empty
// input, mismatched lengths, an unrecognised Averaging, single-class data for
// the score metrics, or out-of-range / non-finite scores and probabilities.
// Degenerate per-class cases are defined rather than rejected: a class with no
// predictions has precision 0, a class with no true samples has recall 0, and a
// class with precision + recall of 0 has F1 0. Inputs are never mutated. The
// log-loss reduction routes through stats.Mean.
package classification
