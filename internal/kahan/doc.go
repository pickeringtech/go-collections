// Package kahan provides the module's single compensated-summation
// accumulator, shared by stats and the ml packages so that every float64
// summary in the library rounds the same way.
//
// It is internal: the accumulator is an implementation detail of the
// statistics, not part of the public API.
//
// # Quick Start
//
// The zero value is an empty sum. Add folds in each term; Total reads the
// compensated result.
//
//	var total Sum
//	for _, x := range []float64{1e17, 1, 2, -1e17} {
//		total.Add(x)
//	}
//	// Result: total.Total() == 3 (a naive running sum returns 0)
//
// # Why Use Compensated Summation?
//
// Native approach — a naive running total loses the low-order bits of every
// term that is small next to the total, and never gets them back:
//
//	var sum float64
//	for _, x := range values { sum += x }
//
// With this package the lost bits are carried in a separate compensation term
// and added back when the total is read:
//
//	var total Sum
//	for _, x := range values { total.Add(x) }
//	sum := total.Total()
//
// The algorithm is Neumaier's improvement of Kahan summation (Kahan–Babuška):
// unlike classic Kahan it stays exact when a term is larger than the running
// total, which is what lets a sliding-window sum subtract a large value that is
// leaving the window without permanently losing the small ones it absorbed.
package kahan
