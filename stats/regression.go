package stats

import "github.com/pickeringtech/go-collections/constraints"

// LineFit holds the result of an ordinary-least-squares straight-line fit:
// the line y = Slope·x + Intercept that minimises the sum of squared vertical
// residuals, together with R2 — the coefficient of determination, the fraction
// of the variance in y explained by the fit, in [0, 1].
//
// Predict evaluates the fitted line at an x value, so callers can derive
// fitted values and residuals (yᵢ − Predict(xᵢ)) without re-deriving the
// coefficients.
type LineFit struct {
	Slope     float64
	Intercept float64
	R2        float64
}

// Predict returns the fitted line's value at x, i.e. Slope·x + Intercept. Pair
// it with the observed yᵢ to obtain a residual: yᵢ − fit.Predict(xᵢ).
func (l LineFit) Predict(x float64) float64 {
	return l.Slope*x + l.Intercept
}

// LinearRegression fits an ordinary-least-squares straight line to the paired
// series x and y, returning the slope, intercept and R² bundled in a LineFit.
//
// The coefficients come from one numerically-stable Welford pass (see
// accumulate): Slope = cov(x,y)/var(x), Intercept = ȳ − Slope·x̄, and
// R² = cov(x,y)² / (var(x)·var(y)) — equivalently the squared Pearson
// correlation. Computing them from accumulated moments rather than the textbook
// Σxy − ΣxΣy/n form keeps the fit accurate on large or near-constant magnitudes.
//
// It returns ok == false when the fit is undefined: when the inputs are empty,
// of differing lengths, have fewer than two points, or when either x or y is
// constant (zero variance — a vertical/horizontal degenerate with no slope to
// fit or no variance to explain). This mirrors Correlation, which is likewise
// undefined for a constant series. Non-finite inputs (NaN/Inf) propagate to a
// non-finite LineFit with ok == true, consistent with the
// variance/covariance/correlation family.
func LinearRegression[T constraints.Numeric](x, y []T) (LineFit, bool) {
	if len(x) != len(y) || len(x) < 2 {
		return LineFit{}, false
	}
	m := accumulate(x, y)
	if m.m2X == 0 || m.m2Y == 0 {
		// A constant x has no spread to fit a slope against; a constant y has no
		// variance for the fit to explain (R² would be 0/0). Either way the
		// regression is undefined. NaN/Inf make the moments non-finite (not
		// zero), so they fall through and propagate as documented.
		return LineFit{}, false
	}
	slope := m.c / m.m2X
	return LineFit{
		Slope:     slope,
		Intercept: m.meanY - slope*m.meanX,
		R2:        (m.c * m.c) / (m.m2X * m.m2Y),
	}, true
}
