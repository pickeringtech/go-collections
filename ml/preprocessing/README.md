# ML Preprocessing

`ml/preprocessing` provides feature-engineering tools that learn parameters from
training data and reuse them on later data: scalers, categorical encoders,
imputers, binners, and dataset splits. The package has no third-party
dependencies.

```go
import "github.com/pickeringtech/go-collections/ml/preprocessing"
```

## Fit on training data, transform later data

Learn each estimator from the training partition, then use the same fitted
estimator for training, validation, and test data. Do not call `Fit` on the
test or validation partition: its statistics would leak into the transformation
and make evaluation less representative.

```go
package main

import (
	"fmt"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func main() {
	train := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	scaler := preprocessing.NewStandardScaler().Fit(train)

	// Transform test data with the mean and standard deviation learned from train.
	got, ok := scaler.Transform([]float64{5, 7})
	fmt.Println(got, ok)
	// Output: [0 1] true
}
```

`Fit` returns the estimator so calls can be chained. `Transform` returns a
result and an `ok` flag; for a learned estimator that has not been fitted, it
returns `nil, false` instead of panicking. An empty or invalid fit may leave an
estimator unfitted. `ConstantImputer` and an `OrdinalEncoder` constructed with
an explicit category order are ready to transform without learning from data.
Transforms return fresh slices and leave their inputs unchanged.

`FitTransform` is shorthand for fitting and transforming the same data. Use it
for training data only; use `Transform` for validation and test data.

## Available operations

| Family | Operations | Purpose |
| --- | --- | --- |
| Scalers | `StandardScaler`, `MinMaxScaler`, `RobustScaler` | Rescale numeric features using training-set statistics. |
| Encoders | `OneHotEncoder`, `OrdinalEncoder`, `LabelEncoder`, `TargetEncoder` | Map categories to indicator columns, codes, or training-target means. |
| Imputers | `MeanImputer`, `MedianImputer`, `ModeImputer`, `ConstantImputer` | Replace values selected by a missing-value predicate. |
| Binners | `FixedWidthBinner`, `QuantileBinner` | Convert numeric values to bins learned from training data. |
| Splits | `TrainTestSplit`, `StratifiedSplit`, `KFold`, `Shuffle` | Create reproducible partitions or shuffled copies. |

An unseen category maps to an all-zero row in `OneHotEncoder`, `-1` in
`LabelEncoder` and `OrdinalEncoder`, and the global training-target mean in
`TargetEncoder`. One-hot columns and learned encoder categories use sorted
order; supply categories to `OrdinalEncoder` when their order has meaning.

## Reproducible train/test split

Pass a seeded generator to make the same split reproducible. The split returns
fresh copies and does not modify the input.

```go
package main

import (
	"fmt"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func main() {
	records := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rng := preprocessing.NewRand(42)
	train, test, ok := preprocessing.TrainTestSplit(records, 0.3, rng)
	if !ok {
		// Handle empty input or an invalid test fraction.
	}
	fmt.Println(len(train), len(test), ok)
	// Output: 7 3 true
}
```

The package also provides `TrainTestSplitSeed`, `StratifiedSplit` for preserving
label proportions, and `KFold` for cross-validation. Each accepts an explicit
seed or generator; the seeded generator is deterministic and is not intended
for security-sensitive randomness.

## NaN and infinity

NaN and infinity handling is specific to each operation, not uniform across
the package. Read the relevant type comments for exact fit and transform
behavior: [scalers](./scalers.go), [encoders](./encoders.go),
[imputers](./imputers.go), and [binners](./binners.go). `MeanImputer` and
`MedianImputer` treat NaN as missing by default; pass a `MissingFunc` when a
different missing-value rule is needed.

## Further documentation

- [Package overview and quick start](./doc.go)
- [Executable examples](./example_test.go)
- [API reference](https://pkg.go.dev/github.com/pickeringtech/go-collections/ml/preprocessing)
