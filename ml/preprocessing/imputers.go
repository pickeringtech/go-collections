package preprocessing

import (
	"math"

	"github.com/pickeringtech/go-collections/stats"
)

// MissingFunc reports whether v should be treated as a missing value that needs
// imputing. It lets the caller decide what "missing" means — a NaN, a sentinel
// such as -1, an empty string, and so on.
type MissingFunc[T any] func(v T) bool

// IsNaN is a MissingFunc[float64] that treats NaN as missing. It is the default
// used by the float imputers (MeanImputer, MedianImputer) when constructed with
// a nil MissingFunc.
func IsNaN(v float64) bool {
	return math.IsNaN(v)
}

// nonMissingFloats returns the elements of input for which isMissing is false.
func nonMissingFloats(input []float64, isMissing MissingFunc[float64]) []float64 {
	out := make([]float64, 0, len(input))
	for _, v := range input {
		if !isMissing(v) {
			out = append(out, v)
		}
	}
	return out
}

// MeanImputer fills missing values with the arithmetic mean of the non-missing
// training values. The fill value is learned at Fit time (via stats.Mean) and
// frozen for every subsequent Transform.
//
// Following the means family, Fit rejects training data whose non-missing
// values are empty or non-finite, leaving the imputer unfitted.
type MeanImputer struct {
	isMissing MissingFunc[float64]
	fill      float64
	fitted    bool
}

// NewMeanImputer returns an unfitted MeanImputer. A nil isMissing defaults to
// IsNaN (NaN values are treated as missing).
func NewMeanImputer(isMissing MissingFunc[float64]) *MeanImputer {
	if isMissing == nil {
		isMissing = IsNaN
	}
	return &MeanImputer{isMissing: isMissing}
}

// Fit learns the mean of the non-missing values in data and returns the
// receiver so calls can be chained. If every value is missing, or the
// non-missing values are non-finite, the imputer is left unfitted.
func (im *MeanImputer) Fit(data []float64) *MeanImputer {
	fill, ok := stats.Mean(nonMissingFloats(data, im.isMissing))
	if !ok {
		return im
	}
	im.fill = fill
	im.fitted = true
	return im
}

// Transform returns a fresh copy of input with every missing value replaced by
// the learned fill value, together with an ok flag (false if unfitted).
func (im *MeanImputer) Transform(input []float64) ([]float64, bool) {
	return imputeFloats(im.fitted, im.isMissing, im.fill, input)
}

// FitTransform fits on input and immediately transforms it. Equivalent to
// Fit(input).Transform(input).
func (im *MeanImputer) FitTransform(input []float64) ([]float64, bool) {
	return im.Fit(input).Transform(input)
}

// Fill returns the fill value learned at Fit time. It is meaningful only after
// a successful Fit.
func (im *MeanImputer) Fill() float64 {
	return im.fill
}

// MedianImputer fills missing values with the median of the non-missing
// training values. The fill value is learned at Fit time (via stats.Median) and
// frozen for every subsequent Transform.
//
// Following the quantile family, Fit rejects training data whose non-missing
// values are empty or non-finite, leaving the imputer unfitted.
type MedianImputer struct {
	isMissing MissingFunc[float64]
	fill      float64
	fitted    bool
}

// NewMedianImputer returns an unfitted MedianImputer. A nil isMissing defaults
// to IsNaN (NaN values are treated as missing).
func NewMedianImputer(isMissing MissingFunc[float64]) *MedianImputer {
	if isMissing == nil {
		isMissing = IsNaN
	}
	return &MedianImputer{isMissing: isMissing}
}

// Fit learns the median of the non-missing values in data and returns the
// receiver so calls can be chained. If every value is missing, or the
// non-missing values are non-finite, the imputer is left unfitted.
func (im *MedianImputer) Fit(data []float64) *MedianImputer {
	fill, ok := stats.Median(nonMissingFloats(data, im.isMissing))
	if !ok {
		return im
	}
	im.fill = fill
	im.fitted = true
	return im
}

// Transform returns a fresh copy of input with every missing value replaced by
// the learned fill value, together with an ok flag (false if unfitted).
func (im *MedianImputer) Transform(input []float64) ([]float64, bool) {
	return imputeFloats(im.fitted, im.isMissing, im.fill, input)
}

// FitTransform fits on input and immediately transforms it. Equivalent to
// Fit(input).Transform(input).
func (im *MedianImputer) FitTransform(input []float64) ([]float64, bool) {
	return im.Fit(input).Transform(input)
}

// Fill returns the fill value learned at Fit time. It is meaningful only after
// a successful Fit.
func (im *MedianImputer) Fill() float64 {
	return im.fill
}

// imputeFloats is the shared float imputation step: a fresh, non-mutating copy
// of input with missing values replaced by fill, gated on fitted.
func imputeFloats(fitted bool, isMissing MissingFunc[float64], fill float64, input []float64) ([]float64, bool) {
	if !fitted {
		return nil, false
	}
	out := make([]float64, len(input))
	for i, v := range input {
		if isMissing(v) {
			out[i] = fill
			continue
		}
		out[i] = v
	}
	return out, true
}

// ModeImputer fills missing values with the most frequent of the non-missing
// training values, learned at Fit time (via stats.Mode) and frozen for every
// subsequent Transform. When several values tie for most frequent, the one that
// appears first in the training data wins.
//
// It is generic over any comparable type. Because stats.Mode rejects non-finite
// values, fitting float categories that contain NaN/Inf among the non-missing
// values leaves the imputer unfitted.
type ModeImputer[T comparable] struct {
	isMissing MissingFunc[T]
	fill      T
	fitted    bool
}

// NewModeImputer returns an unfitted ModeImputer. A nil isMissing means no value
// is considered missing (Transform is then the identity).
func NewModeImputer[T comparable](isMissing MissingFunc[T]) *ModeImputer[T] {
	return &ModeImputer[T]{isMissing: isMissing}
}

// Fit learns the modal non-missing value of data and returns the receiver so
// calls can be chained. If every value is missing (or rejected by stats.Mode)
// the imputer is left unfitted.
func (im *ModeImputer[T]) Fit(data []T) *ModeImputer[T] {
	present := make([]T, 0, len(data))
	for _, v := range data {
		if im.isMissing != nil && im.isMissing(v) {
			continue
		}
		present = append(present, v)
	}
	modes, ok := stats.Mode(present)
	if !ok || len(modes) == 0 {
		return im
	}
	im.fill = modes[0]
	im.fitted = true
	return im
}

// Transform returns a fresh copy of input with every missing value replaced by
// the learned modal value, together with an ok flag (false if unfitted).
func (im *ModeImputer[T]) Transform(input []T) ([]T, bool) {
	if !im.fitted {
		return nil, false
	}
	out := make([]T, len(input))
	for i, v := range input {
		if im.isMissing != nil && im.isMissing(v) {
			out[i] = im.fill
			continue
		}
		out[i] = v
	}
	return out, true
}

// FitTransform fits on input and immediately transforms it. Equivalent to
// Fit(input).Transform(input).
func (im *ModeImputer[T]) FitTransform(input []T) ([]T, bool) {
	return im.Fit(input).Transform(input)
}

// Fill returns the modal value learned at Fit time. It is meaningful only after
// a successful Fit.
func (im *ModeImputer[T]) Fill() T {
	return im.fill
}

// ConstantImputer fills missing values with a fixed constant supplied at
// construction. It needs no training data, so it is ready to Transform
// immediately; Fit is provided only for symmetry with the other estimators and
// is a no-op.
type ConstantImputer[T any] struct {
	isMissing MissingFunc[T]
	fill      T
}

// NewConstantImputer returns a ConstantImputer that replaces missing values with
// fill. A nil isMissing means no value is considered missing (Transform is then
// the identity).
func NewConstantImputer[T any](fill T, isMissing MissingFunc[T]) *ConstantImputer[T] {
	return &ConstantImputer[T]{isMissing: isMissing, fill: fill}
}

// Fit is a no-op that returns the receiver; a ConstantImputer learns nothing
// from data and is always ready to Transform.
func (im *ConstantImputer[T]) Fit(_ []T) *ConstantImputer[T] {
	return im
}

// Transform returns a fresh copy of input with every missing value replaced by
// the constant fill value. The ok flag is always true.
func (im *ConstantImputer[T]) Transform(input []T) ([]T, bool) {
	out := make([]T, len(input))
	for i, v := range input {
		if im.isMissing != nil && im.isMissing(v) {
			out[i] = im.fill
			continue
		}
		out[i] = v
	}
	return out, true
}

// FitTransform transforms input directly (Fit is a no-op for a ConstantImputer).
func (im *ConstantImputer[T]) FitTransform(input []T) ([]T, bool) {
	return im.Transform(input)
}

// Fill returns the constant fill value.
func (im *ConstantImputer[T]) Fill() T {
	return im.fill
}
