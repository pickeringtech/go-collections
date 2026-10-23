package stats

import "math"

// nonFinite reports whether x is NaN or ±Inf — values for which the means in
// this package are undefined and so are rejected.
func nonFinite(x float64) bool {
	return math.IsNaN(x) || math.IsInf(x, 0)
}
