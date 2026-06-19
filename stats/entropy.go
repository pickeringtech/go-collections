package stats

import "math"

// frequencies tallies how often each distinct value occurs in input. It returns
// the per-value counts and the total element count. A non-finite floating-point
// value (NaN/±Inf) makes ok false, mirroring Mode: NaN never compares equal to
// itself so it cannot be counted coherently, and ±Inf is rejected alongside it
// for a uniform categorical non-finite policy.
func frequencies[T comparable](input []T) (counts map[T]int, total int, ok bool) {
	if len(input) == 0 {
		return nil, 0, false
	}
	counts = make(map[T]int, len(input))
	for _, v := range input {
		if nonFiniteComparable(v) {
			return nil, 0, false
		}
		counts[v]++
	}
	return counts, len(input), true
}

// Entropy returns the Shannon entropy of input's value distribution, in bits
// (log base 2): H = −Σ pᵢ·log2(pᵢ), where pᵢ is the empirical probability of
// each distinct value. It measures the uncertainty of the distribution: 0 when
// every element is identical (no uncertainty), rising to log2(k) bits when k
// distinct values are equiprobable (maximum uncertainty).
//
// Entropy works on any comparable T (categories, strings, ints …), not just
// numbers — it summarises a distribution, not a magnitude. The probabilities
// are summed with Kahan compensated summation.
//
// ok is false (and the result 0) when the distribution is undefined: when input
// is empty, or, for floating-point element types, when input contains a
// non-finite value (NaN or ±Inf), which cannot be counted coherently — the same
// rejection policy as Mode.
func Entropy[T comparable](input []T) (float64, bool) {
	counts, total, ok := frequencies(input)
	if !ok {
		return 0, false
	}
	n := float64(total)
	var h kahan
	for _, c := range counts {
		p := float64(c) / n
		h.add(p * math.Log2(p))
	}
	// −Σ p·log2(p); negate the accumulated sum. The +0 guards against returning
	// a signed negative zero when every element is identical (a single term
	// p=1, log2(1)=0).
	return -h.sum + 0, true
}

// Gini returns the Gini impurity of input's value distribution:
// G = 1 − Σ pᵢ², where pᵢ is the empirical probability of each distinct value.
// It is the probability that two independently drawn elements have different
// values: 0 when every element is identical (pure), rising toward 1 − 1/k for k
// equiprobable values. It is the impurity measure used by decision-tree
// learners (CART), a sibling to Entropy as a distribution measure.
//
// Gini works on any comparable T, not just numbers. The squared probabilities
// are summed with Kahan compensated summation.
//
// ok is false (and the result 0) when the distribution is undefined: when input
// is empty, or, for floating-point element types, when input contains a
// non-finite value (NaN or ±Inf) — the same rejection policy as Mode.
func Gini[T comparable](input []T) (float64, bool) {
	counts, total, ok := frequencies(input)
	if !ok {
		return 0, false
	}
	n := float64(total)
	var sumSq kahan
	for _, c := range counts {
		p := float64(c) / n
		sumSq.add(p * p)
	}
	return 1 - sumSq.sum, true
}
