package kahan

import "math"

// Sum is a compensated summation accumulator (Neumaier's variant of Kahan
// summation). It keeps the round-off lost by each addition in a separate
// compensation term and adds it back in Total, so a long sum is accurate to
// within a few ulps of the exact result instead of drifting with every term.
//
// The zero value is an empty sum (Total returns 0).
//
// Non-finite values follow IEEE-754 arithmetic: once the running total is ±Inf
// (from an infinite term or from overflow) it stays infinite, and adding the
// opposite infinity, or any NaN, makes it NaN.
type Sum struct {
	sum, c float64
}

// Add folds x into the running total, accumulating the rounding error of the
// addition in the compensation term.
func (s *Sum) Add(x float64) {
	t := s.sum + x
	if math.Abs(s.sum) >= math.Abs(x) {
		s.c += (s.sum - t) + x // the low-order bits of x were lost
	} else {
		s.c += (x - t) + s.sum // the low-order bits of the old total were lost
	}
	s.sum = t
}

// Total returns the compensated sum of every value added so far.
func (s *Sum) Total() float64 {
	if math.IsInf(s.sum, 0) {
		// Once the running total is infinite there is no round-off left to
		// recover, and the error terms (Inf − Inf) have made c NaN or infinite.
		// The IEEE total is the answer; checking here rather than in Add keeps
		// the per-term hot path branch-light. (c can only go non-finite after
		// the total has: while every partial sum is finite, so is each error
		// term.)
		return s.sum
	}
	return s.sum + s.c
}
