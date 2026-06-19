package preprocessing

import (
	"sort"

	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/stats"
)

// isNaN reports whether v is a floating-point NaN. It works for any ordered type
// because only NaN is unequal to itself; for non-float types it is always false.
// NaN cannot be used as a category: it never compares equal to itself, so it
// would defeat map deduplication and equality/binary-search lookups. The
// encoders therefore drop NaN categories at Fit time, and a NaN at Transform
// time falls through to the unseen-value path (all-zero row, code -1, or the
// global target mean).
func isNaN[C constraints.Ordered](v C) bool {
	return v != v
}

// sortedUnique returns the distinct, non-NaN values of input in ascending order.
// It is the basis for the encoders' stable, documented column ordering.
func sortedUnique[C constraints.Ordered](input []C) []C {
	seen := make(map[C]struct{}, len(input))
	out := make([]C, 0, len(input))
	for _, v := range input {
		if isNaN(v) {
			continue
		}
		_, ok := seen[v]
		if ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// indexOf returns the position of target in the ascending categories slice, or
// -1 if absent, using binary search.
func indexOf[C constraints.Ordered](categories []C, target C) int {
	i := sort.Search(len(categories), func(i int) bool { return categories[i] >= target })
	if i < len(categories) && categories[i] == target {
		return i
	}
	return -1
}

// OneHotEncoder maps each category to a row of indicator columns: a 1 in the
// column for that category and 0 elsewhere. Columns are ordered by sorted
// category value, a stable ordering that does not depend on row order in the
// training data. A value unseen at Fit time transforms to an all-zero row.
type OneHotEncoder[C constraints.Ordered] struct {
	categories []C
	fitted     bool
}

// NewOneHotEncoder returns an unfitted OneHotEncoder. Call Fit before Transform.
func NewOneHotEncoder[C constraints.Ordered]() *OneHotEncoder[C] {
	return &OneHotEncoder[C]{}
}

// Fit learns the sorted set of distinct categories from data and returns the
// receiver so calls can be chained. Empty input leaves the encoder unfitted.
func (e *OneHotEncoder[C]) Fit(data []C) *OneHotEncoder[C] {
	if len(data) == 0 {
		return e
	}
	e.categories = sortedUnique(data)
	e.fitted = true
	return e
}

// Transform maps each value in input to its one-hot row, returning a
// [][]float64 (one row per input value, one column per learned category in
// sorted order) and an ok flag (false if unfitted). An unseen value yields an
// all-zero row.
func (e *OneHotEncoder[C]) Transform(input []C) ([][]float64, bool) {
	if !e.fitted {
		return nil, false
	}
	out := make([][]float64, len(input))
	for i, v := range input {
		row := make([]float64, len(e.categories))
		idx := indexOf(e.categories, v)
		if idx >= 0 {
			row[idx] = 1
		}
		out[i] = row
	}
	return out, true
}

// FitTransform fits on input and immediately transforms it. Equivalent to
// Fit(input).Transform(input).
func (e *OneHotEncoder[C]) FitTransform(input []C) ([][]float64, bool) {
	return e.Fit(input).Transform(input)
}

// Categories returns a copy of the learned categories in column order. It is
// meaningful only after a successful Fit.
func (e *OneHotEncoder[C]) Categories() []C {
	out := make([]C, len(e.categories))
	copy(out, e.categories)
	return out
}

// LabelEncoder maps each category to an integer code in [0, k), assigned by
// sorted category value, and supports the inverse mapping. It is the
// general-purpose category↔integer encoding (commonly used for target labels).
// A value unseen at Fit time transforms to -1.
type LabelEncoder[C constraints.Ordered] struct {
	categories []C
	fitted     bool
}

// NewLabelEncoder returns an unfitted LabelEncoder. Call Fit before Transform.
func NewLabelEncoder[C constraints.Ordered]() *LabelEncoder[C] {
	return &LabelEncoder[C]{}
}

// Fit learns the sorted set of distinct categories from data and returns the
// receiver so calls can be chained. Empty input leaves the encoder unfitted.
func (e *LabelEncoder[C]) Fit(data []C) *LabelEncoder[C] {
	if len(data) == 0 {
		return e
	}
	e.categories = sortedUnique(data)
	e.fitted = true
	return e
}

// Transform maps each value in input to its integer code, returning an []int
// and an ok flag (false if unfitted). An unseen value maps to -1.
func (e *LabelEncoder[C]) Transform(input []C) ([]int, bool) {
	if !e.fitted {
		return nil, false
	}
	out := make([]int, len(input))
	for i, v := range input {
		out[i] = indexOf(e.categories, v)
	}
	return out, true
}

// FitTransform fits on input and immediately transforms it. Equivalent to
// Fit(input).Transform(input).
func (e *LabelEncoder[C]) FitTransform(input []C) ([]int, bool) {
	return e.Fit(input).Transform(input)
}

// InverseTransform maps integer codes back to their categories, returning a
// []C and an ok flag. ok is false if the encoder is unfitted or any code is
// out of the range [0, k).
func (e *LabelEncoder[C]) InverseTransform(codes []int) ([]C, bool) {
	if !e.fitted {
		return nil, false
	}
	out := make([]C, len(codes))
	for i, code := range codes {
		if code < 0 || code >= len(e.categories) {
			return nil, false
		}
		out[i] = e.categories[code]
	}
	return out, true
}

// Categories returns a copy of the learned categories indexed by their code. It
// is meaningful only after a successful Fit.
func (e *LabelEncoder[C]) Categories() []C {
	out := make([]C, len(e.categories))
	copy(out, e.categories)
	return out
}

// OrdinalEncoder maps each category to an integer code reflecting a meaningful
// order. Unlike LabelEncoder, the order is the caller's to define: pass the
// categories from lowest to highest to NewOrdinalEncoder. When constructed with
// no explicit categories, Fit falls back to learning them in sorted order. A
// value unseen at encode time transforms to -1.
type OrdinalEncoder[C constraints.Ordered] struct {
	categories []C
	explicit   bool
	fitted     bool
}

// NewOrdinalEncoder returns an OrdinalEncoder. When categories are supplied they
// define the code order (deduplicated, order preserved) and the encoder is
// ready to Transform immediately. When none are supplied, call Fit to learn them
// in sorted order.
func NewOrdinalEncoder[C constraints.Ordered](categories ...C) *OrdinalEncoder[C] {
	if len(categories) == 0 {
		return &OrdinalEncoder[C]{}
	}
	seen := make(map[C]struct{}, len(categories))
	order := make([]C, 0, len(categories))
	for _, v := range categories {
		if isNaN(v) {
			continue
		}
		_, ok := seen[v]
		if ok {
			continue
		}
		seen[v] = struct{}{}
		order = append(order, v)
	}
	return &OrdinalEncoder[C]{categories: order, explicit: true, fitted: true}
}

// Fit learns the sorted set of distinct categories from data and returns the
// receiver so calls can be chained. When the encoder was constructed with an
// explicit category order, Fit keeps that order and is a no-op. Empty input on
// an encoder without an explicit order leaves it unfitted.
func (e *OrdinalEncoder[C]) Fit(data []C) *OrdinalEncoder[C] {
	if e.explicit {
		return e
	}
	if len(data) == 0 {
		return e
	}
	e.categories = sortedUnique(data)
	e.fitted = true
	return e
}

// Transform maps each value in input to its ordinal code, returning an []int
// and an ok flag (false if unfitted). An unseen value maps to -1. Codes for an
// explicit order follow that order; otherwise they follow sorted order.
func (e *OrdinalEncoder[C]) Transform(input []C) ([]int, bool) {
	if !e.fitted {
		return nil, false
	}
	out := make([]int, len(input))
	for i, v := range input {
		out[i] = e.codeOf(v)
	}
	return out, true
}

// codeOf returns the code for v, or -1 if absent. Explicit (caller-ordered)
// categories are not sorted, so a linear scan is used; learned categories are
// sorted and use binary search.
func (e *OrdinalEncoder[C]) codeOf(v C) int {
	if e.explicit {
		for i, c := range e.categories {
			if c == v {
				return i
			}
		}
		return -1
	}
	return indexOf(e.categories, v)
}

// FitTransform fits on input and immediately transforms it. Equivalent to
// Fit(input).Transform(input).
func (e *OrdinalEncoder[C]) FitTransform(input []C) ([]int, bool) {
	return e.Fit(input).Transform(input)
}

// Categories returns a copy of the categories indexed by their code. It is
// meaningful only after a successful Fit (or immediately, when constructed with
// an explicit order).
func (e *OrdinalEncoder[C]) Categories() []C {
	out := make([]C, len(e.categories))
	copy(out, e.categories)
	return out
}

// TargetEncoder replaces each category with the mean of a numeric target over
// the training rows of that category — a compact, leakage-aware alternative to
// one-hot for high-cardinality features. The per-category means and a global
// fallback mean are learned at Fit time. A value unseen at Fit time transforms
// to the global mean.
type TargetEncoder[C constraints.Ordered] struct {
	means      map[C]float64
	globalMean float64
	fitted     bool
}

// NewTargetEncoder returns an unfitted TargetEncoder. Call Fit before Transform.
func NewTargetEncoder[C constraints.Ordered]() *TargetEncoder[C] {
	return &TargetEncoder[C]{}
}

// Fit learns each category's mean target and the global mean target, returning
// the receiver so calls can be chained. It leaves the encoder unfitted when
// categories and target differ in length, when input is empty, or when the
// target contains non-finite values (per the means family).
func (e *TargetEncoder[C]) Fit(categories []C, target []float64) *TargetEncoder[C] {
	if len(categories) != len(target) || len(categories) == 0 {
		return e
	}
	global, ok := stats.Mean(target)
	if !ok {
		return e
	}

	grouped := make(map[C][]float64, len(categories))
	for i, c := range categories {
		if isNaN(c) {
			// NaN cannot be a map key that matches itself; skip it so a NaN
			// category falls through to the global mean at Transform time.
			continue
		}
		grouped[c] = append(grouped[c], target[i])
	}
	means := make(map[C]float64, len(grouped))
	for c, values := range grouped {
		mean, _ := stats.Mean(values)
		means[c] = mean
	}

	e.means = means
	e.globalMean = global
	e.fitted = true
	return e
}

// Transform maps each value in input to its learned target mean, returning a
// []float64 and an ok flag (false if unfitted). An unseen value maps to the
// global mean.
func (e *TargetEncoder[C]) Transform(input []C) ([]float64, bool) {
	if !e.fitted {
		return nil, false
	}
	out := make([]float64, len(input))
	for i, v := range input {
		mean, ok := e.means[v]
		if !ok {
			mean = e.globalMean
		}
		out[i] = mean
	}
	return out, true
}

// GlobalMean returns the global target mean learned at Fit time, used as the
// fallback for unseen categories. It is meaningful only after a successful Fit.
func (e *TargetEncoder[C]) GlobalMean() float64 {
	return e.globalMean
}
