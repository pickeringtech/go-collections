# Stats - Numeric Summary Statistics

The `stats` package provides summary statistics over slices of numbers. Functions
are plain free functions over `[]T` (any `constraints.Numeric` type) and return a
result plus an `ok` bool, following the library's `(result, ok)` empty contract.

## Quick Start

```go
import "github.com/pickeringtech/go-collections/stats"

data := []float64{1, 2, 3, 4, 5}

med, _ := stats.Quantile(data, 0.5)   // 3   — the median
p90, _ := stats.Percentile(data, 90)  // 4.6
qs, _ := stats.Quartiles(data)        // {Q1:2, Q2:3, Q3:4}
iqr, _ := stats.IQR(data)             // 2
```

## Quantiles & Percentiles

- `Quantile(input, q)` for `q ∈ [0, 1]`.
- `Percentile(input, p)` for `p ∈ [0, 100]` — exactly `Quantile(input, p/100)`.
- `Quartiles(input)` returns a `QuartileSet{Q1, Q2, Q3}`.
- `IQR(input)` returns `Q3 - Q1`.

The median is simply `Quantile(input, 0.5)`.

## Interpolation

When the requested rank falls between two samples, the value is interpolated.
The default everywhere is **Linear** ("type 7" in Hyndman & Fan, 1996), which
matches `numpy.percentile`'s default — the convention most users expect.

`QuantileWith` / `PercentileWith` accept an explicit `InterpolationMethod`:

| Method     | Behaviour                                             |
|------------|-------------------------------------------------------|
| `Linear`   | Linear interpolation between the two samples (default)|
| `Lower`    | The lower of the two bracketing samples               |
| `Higher`   | The higher of the two bracketing samples              |
| `Nearest`  | The sample nearest the desired rank (ties round up)   |
| `Midpoint` | The mean of the two bracketing samples                |

```go
v, _ := stats.QuantileWith(data, 0.9, stats.Midpoint) // 4.5
```

## NaN policy

A `NaN` has no defined ordering, so any quantile of a sample containing one is
undefined. Rather than return a silently-wrong number, the functions report
`ok == false` when the input contains a `NaN`. Integer inputs can never be
`NaN`, so this only affects float inputs.

## Ownership

Inputs are **never mutated**. A function that needs sorted data copies the input
first, consistent with the library's ownership-isolation direction.
