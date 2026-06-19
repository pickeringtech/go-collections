// Package ranking scores ordered result lists — the information-retrieval
// metrics NDCG (with its DCG building block) and mean average precision — as
// pure functions over slices.
//
// It is part of the ml/metrics family (see the ml umbrella package). Where the
// classification package scores an unordered set of predictions, ranking cares
// about order: a relevant item near the top of the list is worth more than the
// same item further down.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/ml/metrics/ranking"
//
//	// Rank items by model score; judge against their true relevance grades.
//	trueRelevance := []float64{3, 2, 3, 0, 1, 2}
//	scores := []float64{6, 5, 4, 3, 2, 1}
//	ndcg, _ := ranking.NDCG(trueRelevance, scores, 0) // 0.9608
//
//	// Average precision over a single ranked list of relevant/not-relevant hits.
//	ap, _ := ranking.AveragePrecision([]bool{true, false, true, false, true}) // 0.7556
//
//	_ = ndcg
//	_ = ap
//
// This Quick Start is compiled and run as Example_quickStart in the package's
// test suite, so it is guaranteed to track the real API.
//
// # Metrics
//
//   - DCG — discounted cumulative gain, Σ relᵢ / log₂(i+2), the building block
//     for NDCG. A cutoff k <= 0 scores the whole list.
//   - NDCG — DCG of a score-ranked list divided by the ideal DCG, normalised to
//     [0, 1] (1 is a perfect ranking).
//   - AveragePrecision — the mean precision-at-k over the relevant positions of
//     one ranked list.
//   - MeanAveragePrecision — the mean of AveragePrecision across several queries.
//
// # Conventions
//
// Every function returns (result, ok) in the library's idiom rather than
// panicking or returning an error. ok is false — and the result the zero value
// — when the inputs cannot be summarised: empty input, mismatched lengths,
// non-finite relevances or scores, an all-zero ideal DCG (NDCG undefined), or a
// ranked list with no relevant items (average precision undefined). A degenerate
// query makes MeanAveragePrecision as a whole undefined, so the mean is always
// over a well-defined set. Inputs are never mutated, and the mean over queries
// routes through stats.Mean.
package ranking
