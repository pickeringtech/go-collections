// Package preprocessing provides deterministic feature-engineering primitives —
// scalers, encoders, imputers, binners and dataset splits — built on a strict
// fit/transform contract so that train and test data are processed with the
// same learned parameters and never leak into one another.
//
// It is the stateful companion to the stateless stats package: where stats
// summarises a slice into a statistic, preprocessing learns parameters from a
// training slice (Fit) and then applies those frozen parameters to any slice
// (Transform). The underlying maths is reused from stats wherever it already
// exists (Mean, PopulationStdDev, Median, IQR, MinMax, Mode, Quantile).
//
// # Fit/Transform Contract — read this first
//
// Every estimator learns from training data and is then applied unchanged to
// new data. THE ENTIRE POINT is to avoid leakage: parameters are captured once,
// at Fit time, and Transform reuses them verbatim. NEVER re-fit on your test or
// validation data.
//
//	scaler := preprocessing.NewStandardScaler()
//	scaler.Fit(trainData)                       // learn mean/stddev from TRAIN only
//	trainZ, _ := scaler.Transform(trainData)    // apply train params to train
//	testZ, _ := scaler.Transform(testData)      // apply the SAME train params to test
//
// Calling Transform on an estimator that has not been fitted returns ok ==
// false (the library's (result, ok) idiom) rather than panicking. FitTransform
// is a convenience for the common train-time case where you fit and transform
// the same slice in one call — it is the ONLY place fit and transform share a
// slice on purpose.
//
// Fit returns the receiver so calls can be chained:
//
//	z, ok := preprocessing.NewStandardScaler().Fit(train).Transform(test)
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/ml/preprocessing"
//
//	// Scale features to zero mean / unit variance using train parameters.
//	scaler := preprocessing.NewStandardScaler().Fit([]float64{2, 4, 4, 4, 5, 5, 7, 9})
//	z, _ := scaler.Transform([]float64{5}) // [0], 5 is the train mean
//
//	// One-hot encode categories in a stable, sorted column order.
//	enc := preprocessing.NewOneHotEncoder[string]().Fit([]string{"b", "a", "a", "c"})
//	rows, _ := enc.Transform([]string{"a"}) // [[1 0 0]], columns a,b,c
//
//	// Fill missing values with the train mean.
//	imp := preprocessing.NewMeanImputer(nil).Fit([]float64{1, 2, 3})
//	filled, _ := imp.Transform([]float64{math.NaN(), 2}) // [2 2]
//
//	// Reproducible train/test split.
//	rng := preprocessing.NewRand(42)
//	train, test, _ := preprocessing.TrainTestSplit([]int{0, 1, 2, 3, 4}, 0.4, rng)
//
//	_ = z
//	_ = rows
//	_ = filled
//	_, _ = train, test
//
// This Quick Start is compiled and run as Example_quickStart in the package's
// test suite, so it is guaranteed to track the real API.
//
// # Conventions
//
//   - The (result, ok) idiom: every Transform returns an ok flag. ok == false
//     means the estimator was not fitted, or the input was empty/invalid; see
//     each type's doc for its specific rejection policy.
//   - Non-mutating: Transform and the split functions never modify their input;
//     they return fresh slices (a non-nil empty slice for empty, non-nil input).
//   - Single implementation + thin accessors: estimators expose accessors
//     (Categories, Edges, …) onto the parameters learned at Fit time.
//   - Determinism: the split functions take an explicit *math/rand/v2 generator
//     (or a seed via the *Seed wrappers), so the same seed reproduces exactly.
//
// # NaN/Inf policy
//
// The policy is inherited per operation from the stats function it reuses, and
// is therefore split by family rather than uniform across the package:
//
//   - StandardScaler and MinMaxScaler PROPAGATE non-finite values, matching
//     stats.Standardize / stats.Normalize: a NaN/Inf in the data flows into the
//     result rather than being rejected.
//   - RobustScaler, the median/mode/mean imputers and the quantile binner REJECT
//     non-finite values (Fit returns an unfitted estimator), matching the
//     means/quantile family in stats.
//
// See each type's doc comment for the exact contract.
package preprocessing
