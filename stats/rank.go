package stats

import "github.com/pickeringtech/go-collections/constraints"

// PercentileOfScore returns the percentile rank of score within input — the
// percentage of values less than or equal to score, in [0, 100]. It is the
// inverse of Percentile: where Percentile maps a percentage to a value,
// PercentileOfScore maps a value to its standing in the sample. A score at or
// above every element ranks 100; below every element it ranks 0.
//
// This is the "weak" definition (the fraction of the sample the score weakly
// dominates, x ≤ score). The result is float64 with an ok flag, and input is
// never mutated nor required to be sorted.
//
// ok is false (and the result 0) when the rank is undefined: when input is
// empty, when score is non-finite (NaN or ±Inf), or when input contains a
// non-finite value — a NaN has no ordering and so cannot be ranked, so per the
// quantile family's rejection policy such inputs are reported rather than
// silently miscounted.
func PercentileOfScore[T constraints.Numeric](input []T, score T) (float64, bool) {
	if len(input) == 0 {
		return 0, false
	}
	s := float64(score)
	if nonFinite(s) {
		return 0, false
	}
	count := 0
	for _, v := range input {
		f := float64(v)
		if nonFinite(f) {
			return 0, false
		}
		if f <= s {
			count++
		}
	}
	return float64(count) / float64(len(input)) * 100, true
}
