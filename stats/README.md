# Stats - Correct Numeric Summaries

The `stats` package summarizes collections of numbers into statistics. It is the
home for operations that reduce slices of numbers to descriptive figures —
variance, standard deviation, covariance, correlation — and for value-rescaling
transforms such as normalization and standardization. The companion `slices`
package owns slice *structure* and element *ordering* (`Min`/`Max`/sorting);
`stats` owns the numeric *summaries*. One operation lives in exactly one place.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/stats"

xs := []float64{1, 2, 3, 4, 5}
ys := []float64{2, 4, 6, 8, 10}

r, _ := stats.Correlation(xs, ys)   // 1 — perfectly linear
z, _ := stats.Standardize(xs)       // zero mean, unit variance
ma, _ := stats.MovingAverage(xs, 2) // [1.5 2.5 3.5 4.5] — full windows only
```

## Means

| Function                       | Returns           | Undefined (`ok == false`) when                          |
| ------------------------------ | ----------------- | ------------------------------------------------------- |
| `WeightedMean(values, weights)`| `(float64, bool)` | empty, differing lengths, negative/zero-sum/NaN weights |
| `GeometricMean(values)`        | `(float64, bool)` | empty, or any non-positive / non-finite value           |
| `HarmonicMean(values)`         | `(float64, bool)` | empty, or any non-positive / non-finite value           |

Sums are accumulated with Kahan compensated summation. The means **reject**
non-finite input (`ok == false`); the relational stats and transforms below
instead let non-finite values propagate (see Conventions).

## Relational statistics

Two equal-length series in, one `float64` out (with an `ok` flag). Mismatched
lengths are rejected with `ok == false`.

| Function                          | Returns           | Undefined (`ok == false`) when                          |
| --------------------------------- | ----------------- | ------------------------------------------------------- |
| `PopulationCovariance(x, y)`      | `(float64, bool)` | empty or differing lengths                              |
| `SampleCovariance(x, y)`          | `(float64, bool)` | fewer than two pairs, or differing lengths              |
| `Correlation(x, y)`               | `(float64, bool)` | fewer than two pairs, differing lengths, constant input |

`Correlation` is Pearson's coefficient in `[−1, 1]`. The `n`/`n−1` factors cancel
in the ratio, so sample and population conventions give the same value — hence a
single function.

## Transforms

Rescale a series into a fresh `[]float64` (the input is never mutated), with an
`ok` flag for undefined input.

| Function                     | Returns             | Behaviour                                                                      |
| ---------------------------- | ------------------- | ------------------------------------------------------------------------------ |
| `Normalize(input)`           | `([]float64, bool)` | min-max scaling to `[0, 1]`; constant *finite* input maps to all-zeros; empty → `false` |
| `Standardize(input)`         | `([]float64, bool)` | z-score `(x − mean) / popStdDev`; zero-spread *finite* input → all-zeros; empty → `false` |
| `MovingAverage(input, w)`    | `([]float64, bool)` | rolling mean over **full windows only**; result length `len−w+1`                |

The all-zeros result for constant/zero-spread input holds for *finite* input;
non-finite values (NaN/Inf) propagate per the package's policy (see Conventions).

### `MovingAverage` edge handling (explicit)

- `w < 1` is invalid → `ok == false`.
- `w > len(input)` cannot form a full window → `ok == false` (covers empty input).
- `w == len(input)` yields a single value: the mean of the whole input.
- Partial leading windows are **not** produced — every output value is the mean
  of exactly `w` elements, so none is a weaker average over fewer points.
- Computed with an incremental running sum, so it is `O(len(input))` regardless
  of window size. A consequence is that a non-finite value propagates to its
  window and every subsequent window; clean `NaN`/`Inf` first if you need strict
  per-window locality.

All operations accept any `constraints.Numeric` slice (`[]int`, `[]float64`, …).

## Advanced

A tier of richer operations, all keeping the package conventions (see below).

| Function                            | Returns             | Undefined (`ok == false`) when                                   |
| ----------------------------------- | ------------------- | ---------------------------------------------------------------- |
| `LinearRegression(x, y)`            | `(LineFit, bool)`   | empty, differing lengths, < 2 points, or constant `x`/`y`        |
| `Histogram(input, bins)`            | `([]Bin, bool)`     | empty, `bins < 1`, non-finite, or zero-width (constant) range    |
| `Skewness(input)`                   | `(float64, bool)`   | empty or constant (zero variance)                                |
| `Kurtosis(input)` (excess)          | `(float64, bool)`   | empty or constant (zero variance)                                |
| `Entropy(input)` (Shannon, bits)    | `(float64, bool)`   | empty, or non-finite float element                               |
| `Gini(input)` (impurity)            | `(float64, bool)`   | empty, or non-finite float element                               |
| `PercentileOfScore(input, score)`   | `(float64, bool)`   | empty, non-finite score, or non-finite element                   |
| `Dot(a, b)`                         | `(float64, bool)`   | empty or differing lengths                                       |
| `Norm(a)` (L2)                      | `(float64, bool)`   | empty                                                            |
| `EuclideanDistance(a, b)`           | `(float64, bool)`   | empty or differing lengths                                       |
| `CosineSimilarity(a, b)`            | `(float64, bool)`   | empty, differing lengths, or a zero vector                       |

- `LinearRegression` returns `LineFit{Slope, Intercept, R2}` with a `Predict(x)`
  method, so residuals are `yᵢ − fit.Predict(xᵢ)`.
- `Histogram` returns equal-width `Bin{Min, Max, Count}` buckets; the final bin's
  upper bound is inclusive so the maximum value is counted, and the counts always
  sum to `len(input)`.
- `Entropy` and `Gini` summarise the distribution of a **categorical** sample over
  any `comparable` type (strings, ints, …), not just numbers.
- The moment-based stats (regression, skewness, kurtosis) and the vector ops
  **propagate** non-finite input; the categorical measures and `PercentileOfScore`
  **reject** it (see Conventions).

## Conventions

These conventions are deliberate and apply uniformly across the package:

- **Numerical stability.** Variance, covariance and correlation use
  [Welford's online algorithm](https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance#Welford's_online_algorithm),
  never the naive `Σxy − ΣxΣy/n`, which loses catastrophic precision on large or
  near-constant magnitudes.
- **Return type.** Scalar summaries return `float64`; transforms return a new
  `[]float64`. Both are paired with an `ok bool`.
- **Empty/edge contract.** Statistics on undefined input return `ok == false`
  rather than a silent zero. Sample variants are undefined for fewer than two
  elements; population variants only for empty input.
- **NaN/Inf policy.** Non-finite inputs propagate: the result is non-finite and
  `ok == true`. Values are never silently filtered out, so a `NaN` in the data
  surfaces as a `NaN` statistic rather than a plausible-looking wrong number.
- **Sample vs population.** Both variants are offered where
  [Bessel's correction](https://en.wikipedia.org/wiki/Bessel%27s_correction)
  applies, named unambiguously, so the choice is always the caller's.
