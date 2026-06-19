package preprocessing

import "github.com/pickeringtech/go-collections/stats"

// StandardScaler rescales features to zero mean and unit variance using the
// z-score (x − mean) / stddev, where mean and the population standard deviation
// are learned at Fit time and frozen for every subsequent Transform.
//
// It mirrors stats.Standardize's semantics: a feature with zero spread in the
// training data transforms to all zeros (every value sits at the mean), and
// non-finite values (NaN/Inf) PROPAGATE rather than being rejected.
type StandardScaler struct {
	mean   float64
	stdDev float64
	fitted bool
}

// NewStandardScaler returns an unfitted StandardScaler. Call Fit before
// Transform.
func NewStandardScaler() *StandardScaler {
	return &StandardScaler{}
}

// Fit learns the mean and population standard deviation from data and returns
// the receiver so calls can be chained. Empty input leaves the scaler unfitted.
// Non-finite training data is not rejected — it propagates into the learned
// parameters and therefore into later Transform output.
func (s *StandardScaler) Fit(data []float64) *StandardScaler {
	sum, ok := stats.Sum(data)
	if !ok {
		return s
	}
	// stats.Sum/PopulationStdDev propagate non-finite values, matching
	// stats.Standardize; we deliberately avoid stats.Mean here because it
	// rejects them.
	s.mean = sum / float64(len(data))
	s.stdDev, _ = stats.PopulationStdDev(data)
	s.fitted = true
	return s
}

// Transform applies the learned z-score to input, returning a fresh []float64
// and an ok flag. ok is false if the scaler is unfitted; otherwise input of any
// length (including empty) transforms successfully. When the learned standard
// deviation is zero, every value maps to 0.
func (s *StandardScaler) Transform(input []float64) ([]float64, bool) {
	if !s.fitted {
		return nil, false
	}
	out := make([]float64, len(input))
	if s.stdDev == 0 {
		return out, true
	}
	for i, v := range input {
		out[i] = (v - s.mean) / s.stdDev
	}
	return out, true
}

// FitTransform fits on input and immediately transforms it — the common
// train-time convenience. It is equivalent to Fit(input).Transform(input).
func (s *StandardScaler) FitTransform(input []float64) ([]float64, bool) {
	return s.Fit(input).Transform(input)
}

// Mean returns the mean learned at Fit time. It is meaningful only after a
// successful Fit.
func (s *StandardScaler) Mean() float64 {
	return s.mean
}

// StdDev returns the population standard deviation learned at Fit time. It is
// meaningful only after a successful Fit.
func (s *StandardScaler) StdDev() float64 {
	return s.stdDev
}

// MinMaxScaler rescales features to the range [0, 1] using min-max scaling
// (x − min) / (max − min), where min and max are learned at Fit time and frozen
// for every subsequent Transform. Test values outside the training range map
// outside [0, 1].
//
// It mirrors stats.Normalize's semantics: a feature with a degenerate training
// range (max == min) transforms to all zeros, and non-finite values (NaN/Inf)
// PROPAGATE rather than being rejected.
type MinMaxScaler struct {
	min    float64
	max    float64
	fitted bool
}

// NewMinMaxScaler returns an unfitted MinMaxScaler. Call Fit before Transform.
func NewMinMaxScaler() *MinMaxScaler {
	return &MinMaxScaler{}
}

// Fit learns the minimum and maximum from data and returns the receiver so
// calls can be chained. Empty input leaves the scaler unfitted.
func (s *MinMaxScaler) Fit(data []float64) *MinMaxScaler {
	lo, hi, ok := stats.MinMax(data)
	if !ok {
		return s
	}
	s.min = lo
	s.max = hi
	s.fitted = true
	return s
}

// Transform applies the learned min-max scaling to input, returning a fresh
// []float64 and an ok flag. ok is false if the scaler is unfitted. When the
// learned range is degenerate (max == min), every value maps to 0.
func (s *MinMaxScaler) Transform(input []float64) ([]float64, bool) {
	if !s.fitted {
		return nil, false
	}
	out := make([]float64, len(input))
	span := s.max - s.min
	if span == 0 {
		return out, true
	}
	for i, v := range input {
		out[i] = (v - s.min) / span
	}
	return out, true
}

// FitTransform fits on input and immediately transforms it — the common
// train-time convenience. It is equivalent to Fit(input).Transform(input).
func (s *MinMaxScaler) FitTransform(input []float64) ([]float64, bool) {
	return s.Fit(input).Transform(input)
}

// Min returns the minimum learned at Fit time. It is meaningful only after a
// successful Fit.
func (s *MinMaxScaler) Min() float64 {
	return s.min
}

// Max returns the maximum learned at Fit time. It is meaningful only after a
// successful Fit.
func (s *MinMaxScaler) Max() float64 {
	return s.max
}

// RobustScaler rescales features using statistics that are robust to outliers:
// (x − median) / IQR, where the median and interquartile range (IQR = Q3 − Q1)
// are learned at Fit time and frozen for every subsequent Transform.
//
// Unlike the mean/variance-based scalers, RobustScaler follows the quantile
// family's policy and REJECTS non-finite training data: a NaN/Inf anywhere in
// the data leaves the scaler unfitted. When the learned IQR is zero, every
// value maps to 0.
type RobustScaler struct {
	median float64
	iqr    float64
	fitted bool
}

// NewRobustScaler returns an unfitted RobustScaler. Call Fit before Transform.
func NewRobustScaler() *RobustScaler {
	return &RobustScaler{}
}

// Fit learns the median and interquartile range from data and returns the
// receiver so calls can be chained. Empty input, or input containing any
// non-finite value, leaves the scaler unfitted.
func (s *RobustScaler) Fit(data []float64) *RobustScaler {
	// Quartiles gives the median (Q2) and IQR (Q3-Q1) in a single rejection
	// check, so there is no separate, unreachable failure path for the IQR.
	qs, ok := stats.Quartiles(data)
	if !ok {
		return s
	}
	s.median = qs.Q2
	s.iqr = qs.Q3 - qs.Q1
	s.fitted = true
	return s
}

// Transform applies the learned robust scaling to input, returning a fresh
// []float64 and an ok flag. ok is false if the scaler is unfitted. When the
// learned IQR is zero, every value maps to 0.
func (s *RobustScaler) Transform(input []float64) ([]float64, bool) {
	if !s.fitted {
		return nil, false
	}
	out := make([]float64, len(input))
	if s.iqr == 0 {
		return out, true
	}
	for i, v := range input {
		out[i] = (v - s.median) / s.iqr
	}
	return out, true
}

// FitTransform fits on input and immediately transforms it — the common
// train-time convenience. It is equivalent to Fit(input).Transform(input).
func (s *RobustScaler) FitTransform(input []float64) ([]float64, bool) {
	return s.Fit(input).Transform(input)
}

// Median returns the median learned at Fit time. It is meaningful only after a
// successful Fit.
func (s *RobustScaler) Median() float64 {
	return s.median
}

// IQR returns the interquartile range learned at Fit time. It is meaningful
// only after a successful Fit.
func (s *RobustScaler) IQR() float64 {
	return s.iqr
}
