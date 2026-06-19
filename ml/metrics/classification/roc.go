package classification

import (
	"math"
	"sort"

	"github.com/pickeringtech/go-collections/stats"
)

// Point is one operating point on a ROC curve: the false-positive and
// true-positive rates achieved at a given score Threshold (samples scoring at
// or above Threshold are classified positive).
type Point struct {
	FPR       float64
	TPR       float64
	Threshold float64
}

// Curve is a receiver-operating-characteristic curve: a sequence of operating
// Points ordered from the strictest threshold (the (0,0) corner, Threshold
// +Inf) to the most permissive (the (1,1) corner).
type Curve struct {
	Points []Point
}

// finiteScores reports whether every score is finite (no NaN, no ±Inf).
func finiteScores(scores []float64) bool {
	for _, s := range scores {
		if math.IsNaN(s) || math.IsInf(s, 0) {
			return false
		}
	}
	return true
}

// countClasses returns the number of positive and negative samples for the
// designated positive label.
func countClasses[T comparable](yTrue []T, positive T) (pos, neg int) {
	for _, y := range yTrue {
		if y == positive {
			pos++
		} else {
			neg++
		}
	}
	return pos, neg
}

// ROCCurve builds the ROC curve for a binary classifier from its true labels
// and the scores it assigned to the positive class, together with an ok flag. A
// higher score means more confidence in the positive label. The curve sweeps
// the decision threshold from strict to permissive; consecutive samples sharing
// a score collapse to a single operating point.
//
// ok is false (and the curve is the zero Curve) when the inputs cannot be
// summarised:
//   - yTrue is empty, or len(yTrue) != len(scores);
//   - any score is non-finite (NaN or ±Inf);
//   - the data is single-class (no positives or no negatives), for which a TPR
//     or FPR is undefined.
func ROCCurve[T comparable](yTrue []T, scores []float64, positive T) (Curve, bool) {
	if len(yTrue) == 0 || len(yTrue) != len(scores) {
		return Curve{}, false
	}
	if !finiteScores(scores) {
		return Curve{}, false
	}

	pos, neg := countClasses(yTrue, positive)
	if pos == 0 || neg == 0 {
		return Curve{}, false
	}

	order := make([]int, len(scores))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(a, b int) bool { return scores[order[a]] > scores[order[b]] })

	points := []Point{{FPR: 0, TPR: 0, Threshold: math.Inf(1)}}
	var tp, fp int
	for k := 0; k < len(order); k++ {
		i := order[k]
		if yTrue[i] == positive {
			tp++
		} else {
			fp++
		}
		// Emit a point only at the end of a run of equal scores, so tied
		// samples share a single threshold.
		tiedWithNext := k+1 < len(order) && scores[order[k+1]] == scores[i]
		if tiedWithNext {
			continue
		}
		points = append(points, Point{
			FPR:       float64(fp) / float64(neg),
			TPR:       float64(tp) / float64(pos),
			Threshold: scores[i],
		})
	}
	return Curve{Points: points}, true
}

// tieAveragedRanks returns the ascending ranks of scores (1-based), assigning
// tied scores their average rank.
func tieAveragedRanks(scores []float64) []float64 {
	n := len(scores)
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(a, b int) bool { return scores[order[a]] < scores[order[b]] })

	ranks := make([]float64, n)
	i := 0
	for i < n {
		j := i
		for j+1 < n && scores[order[j+1]] == scores[order[i]] {
			j++
		}
		// Positions i..j (0-based) occupy ranks i+1..j+1; share their average.
		avg := float64(i+j+2) / 2
		for k := i; k <= j; k++ {
			ranks[order[k]] = avg
		}
		i = j + 1
	}
	return ranks
}

// AUC returns the area under the ROC curve for a binary classifier, together
// with an ok flag. It is computed from the Mann–Whitney U statistic — the
// probability that a random positive sample scores above a random negative one
// — using tie-averaged ranks, which is exact and avoids accumulating the
// trapezoid error of integrating the curve directly. An AUC of 1 is a perfect
// ranker, 0.5 is random, and below 0.5 is worse than random.
//
// ok is false (and the result is 0) under the same conditions as ROCCurve:
// empty or unequal-length input, any non-finite score, or single-class data.
func AUC[T comparable](yTrue []T, scores []float64, positive T) (float64, bool) {
	if len(yTrue) == 0 || len(yTrue) != len(scores) {
		return 0, false
	}
	if !finiteScores(scores) {
		return 0, false
	}

	pos, neg := countClasses(yTrue, positive)
	if pos == 0 || neg == 0 {
		return 0, false
	}

	ranks := tieAveragedRanks(scores)
	var sumRanksPos float64
	for i := range yTrue {
		if yTrue[i] == positive {
			sumRanksPos += ranks[i]
		}
	}

	// U = ΣranksPos − nPos(nPos+1)/2; AUC = U / (nPos·nNeg).
	u := sumRanksPos - float64(pos)*float64(pos+1)/2
	return u / (float64(pos) * float64(neg)), true
}

// LogLoss returns the binary logistic loss (cross-entropy) of probabilistic
// predictions, together with an ok flag. probs[i] is the predicted probability
// that sample i belongs to the positive label; the loss is the mean of
// −log(p) over positive samples and −log(1−p) over the rest. Lower is better,
// with 0 the unreachable ideal of perfect, confident predictions.
//
// Probabilities are clamped to [1e-15, 1−1e-15] before the logarithm so a
// confident-but-wrong prediction yields a large finite penalty rather than
// +Inf. ok is false (and the result is 0) when the inputs cannot be summarised:
// yTrue is empty, len(yTrue) != len(probs), or any probability is NaN or lies
// outside [0, 1].
//
// This is the binary case; multiclass log-loss over a probability matrix is not
// yet provided.
func LogLoss[T comparable](yTrue []T, probs []float64, positive T) (float64, bool) {
	if len(yTrue) == 0 || len(yTrue) != len(probs) {
		return 0, false
	}

	const eps = 1e-15
	losses := make([]float64, len(yTrue))
	for i := range yTrue {
		p := probs[i]
		if math.IsNaN(p) || p < 0 || p > 1 {
			return 0, false
		}
		p = math.Min(math.Max(p, eps), 1-eps)
		if yTrue[i] == positive {
			losses[i] = -math.Log(p)
		} else {
			losses[i] = -math.Log(1 - p)
		}
	}
	return stats.Mean(losses)
}
