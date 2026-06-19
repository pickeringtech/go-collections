package ranking

import (
	"math"
	"sort"

	"github.com/pickeringtech/go-collections/stats"
)

// nonFinite reports whether x is NaN or ±Inf.
func nonFinite(x float64) bool {
	return math.IsNaN(x) || math.IsInf(x, 0)
}

// effectiveK clamps a requested cutoff to the list length: k <= 0 (or k beyond
// the list) means "use the whole list".
func effectiveK(k, n int) int {
	if k <= 0 || k > n {
		return n
	}
	return k
}

// dcgInOrder sums the discounted gains relᵢ / log₂(i+2) over the first k
// positions of the already-ordered relevances.
func dcgInOrder(relevances []float64, k int) float64 {
	var sum float64
	for i := 0; i < k; i++ {
		sum += relevances[i] / math.Log2(float64(i+2))
	}
	return sum
}

// DCG returns the discounted cumulative gain of relevances, taken in the given
// order — Σ relᵢ / log₂(i+2) over the first k items — together with an ok flag.
// Later positions are discounted, so the same relevances earlier in the list
// score higher. k is the cutoff rank; k <= 0 (or k larger than the list) scores
// the whole list.
//
// ok is false (and the result is 0) when relevances is empty or any relevance
// is non-finite (NaN/±Inf).
func DCG(relevances []float64, k int) (float64, bool) {
	if len(relevances) == 0 {
		return 0, false
	}
	for _, r := range relevances {
		if nonFinite(r) {
			return 0, false
		}
	}
	return dcgInOrder(relevances, effectiveK(k, len(relevances))), true
}

// NDCG returns the normalised discounted cumulative gain at cutoff k, together
// with an ok flag. Items are ranked by scores (highest first); NDCG is the DCG
// of their true relevances in that ranking divided by the ideal DCG — the DCG
// of the relevances sorted into their best possible order. The result lies in
// [0, 1], with 1 for a perfect ranking. k <= 0 scores the whole list.
//
// ok is false (and the result is 0) when the inputs cannot be summarised:
//   - trueRelevance is empty, or len(trueRelevance) != len(scores);
//   - any relevance or score is non-finite (NaN/±Inf);
//   - the ideal DCG is 0 (every relevance is 0), so NDCG is undefined.
func NDCG(trueRelevance, scores []float64, k int) (float64, bool) {
	n := len(trueRelevance)
	if n == 0 || n != len(scores) {
		return 0, false
	}
	for i := 0; i < n; i++ {
		if nonFinite(trueRelevance[i]) || nonFinite(scores[i]) {
			return 0, false
		}
	}

	// Rank the relevances by descending score (stable, so equal scores keep
	// their input order).
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(a, b int) bool { return scores[order[a]] > scores[order[b]] })
	ranked := make([]float64, n)
	for rank, idx := range order {
		ranked[rank] = trueRelevance[idx]
	}

	ideal := append([]float64(nil), trueRelevance...)
	sort.Sort(sort.Reverse(sort.Float64Slice(ideal)))

	cutoff := effectiveK(k, n)
	idcg := dcgInOrder(ideal, cutoff)
	if idcg == 0 {
		return 0, false
	}
	return dcgInOrder(ranked, cutoff) / idcg, true
}

// AveragePrecision returns the average precision of a single ranked result
// list, together with an ok flag. ranked[i] reports whether the item at rank
// i+1 is relevant. AP is the mean of the precision-at-k values taken at each
// relevant position, which rewards placing relevant items near the top.
//
// ok is false (and the result is 0) when ranked is empty or contains no
// relevant items (AP would divide by zero relevant documents).
func AveragePrecision(ranked []bool) (float64, bool) {
	if len(ranked) == 0 {
		return 0, false
	}

	var hits int
	var sum float64
	for i, relevant := range ranked {
		if relevant {
			hits++
			sum += float64(hits) / float64(i+1) // precision@(i+1)
		}
	}
	if hits == 0 {
		return 0, false
	}
	return sum / float64(hits), true
}

// MeanAveragePrecision returns the mean of AveragePrecision over several ranked
// result lists (one per query), together with an ok flag. It is the standard
// summary of ranking quality across a set of queries.
//
// ok is false (and the result is 0) when queries is empty or any query is
// itself undefined under AveragePrecision (empty or with no relevant items), so
// the mean is taken over a well-defined set.
func MeanAveragePrecision(queries [][]bool) (float64, bool) {
	if len(queries) == 0 {
		return 0, false
	}

	aps := make([]float64, len(queries))
	for i, q := range queries {
		ap, ok := AveragePrecision(q)
		if !ok {
			return 0, false
		}
		aps[i] = ap
	}
	return stats.Mean(aps)
}
