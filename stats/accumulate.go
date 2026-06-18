package stats

import "math"

// kahan is a Kahan (compensated) summation accumulator. It tracks a running
// compensation term so that adding many values does not lose low-order bits to
// floating-point round-off the way a naive running total does. Use it via the
// zero value: var k kahan; k.add(x); ...; total := k.sum.
type kahan struct {
	sum, c float64
}

// add folds x into the running total, carrying the round-off error forward in
// the compensation term so it is not lost on the next addition.
func (k *kahan) add(x float64) {
	y := x - k.c
	t := k.sum + y
	k.c = (t - k.sum) - y
	k.sum = t
}

// nonFinite reports whether x is NaN or ±Inf — values for which the means in
// this package are undefined and so are rejected.
func nonFinite(x float64) bool {
	return math.IsNaN(x) || math.IsInf(x, 0)
}
