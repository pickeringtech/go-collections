# Relational - Split-Apply-Combine for Slices

The `relational` package provides the data-engineering relational primitives —
`GroupBy` and aggregate, the four joins, pivot/unpivot, and partition — as free
functions over Go's native slices and maps. It is the "split-apply-combine"
toolkit: reshape rows into groups, reduce each group, stitch tables together by
key, and flip between long and wide layouts, without writing the bookkeeping
loops by hand.

Every function follows the library's contracts: it never mutates its input, and
nil or empty input yields a non-nil empty result, so callers can `range` or
`len` it without a nil check.

## Quick Start

```go
import (
    "github.com/pickeringtech/go-collections/relational"
    "github.com/pickeringtech/go-collections/stats"
)

type Order struct {
    Dept   string
    Amount float64
}
orders := []Order{{"books", 10}, {"books", 30}, {"toys", 5}}

// GROUP BY dept, then AVG(amount). Aggregate takes any reducer of shape
// func([]T) (R, bool), so stats.Mean plugs straight in via AggregateBy.
byDept := relational.GroupBy(orders, func(o Order) string { return o.Dept })
avg := relational.AggregateBy(byDept,
    func(o Order) float64 { return o.Amount },
    stats.Mean[float64],
)
// avg = map[books:20 toys:5]

// JOIN orders to a per-dept label table by dept key.
type Label struct {
    Dept string
    Name string
}
labels := []Label{{"books", "Books Dept"}, {"toys", "Toys Dept"}}
joined := relational.InnerJoin(orders, labels,
    func(o Order) string { return o.Dept },
    func(l Label) string { return l.Dept },
)
```

## Grouping & aggregation

| Function                          | Returns        | Notes                                                   |
| --------------------------------- | -------------- | ------------------------------------------------------- |
| `GroupBy(input, keyFn)`           | `map[K][]V`    | buckets values by key; first-seen order within a group  |
| `GroupBySeq(seq, keyFn)`          | `map[K][]V`    | same, over an `iter.Seq` stream (consumed once)         |
| `CountBy(input, keyFn)`           | `map[K]int`    | per-key counts; cheaper than `GroupBy` for sizes only   |
| `Aggregate(groups, aggFn)`        | `map[K]R`      | reduce each group with a `func([]V) (R, bool)` reducer  |
| `AggregateBy(groups, proj, aggFn)`| `map[K]R`      | project each value to `N` first, then reduce            |

A group whose `aggFn` returns `ok == false` is **omitted** from the result map
(not stored as a zero value), preserving the `stats` `(result, ok)` idiom end to
end. A present key therefore always carries a defined result.

## Joins

All joins are **many-to-many**: if a key occurs `a` times left and `b` times
right, all `a*b` matching combinations are emitted. Each result row is a
`JoinPair[L, R]` carrying `Left`, `Right`, and `LeftOK`/`RightOK` flags (values,
not pointers — an unmatched side is the zero value with its `OK == false`).

| Function          | Emits                                                        |
| ----------------- | ----------------------------------------------------------- |
| `InnerJoin`       | matched rows only (`LeftOK && RightOK`)                      |
| `LeftJoin`        | matched rows + unmatched left rows (right zero, `RightOK=false`) |
| `RightJoin`       | matched rows + unmatched right rows (left zero, `LeftOK=false`)  |
| `FullOuterJoin`   | matched rows + unmatched left + unmatched right              |

Internally each join builds a private `map[K][]R` index, so the match is
`O(n+m)` rather than the `O(n*m)` of a nested loop.

## Reshaping

| Function                                | Returns              | Notes                                        |
| --------------------------------------- | -------------------- | -------------------------------------------- |
| `Pivot(rows, rowKey, colKey, value)`    | `map[K]map[K2]V`     | long → wide; collisions are last-write-wins  |
| `Unpivot(wide)`                         | `[]Cell[K, K2, V]`   | wide → long; inverse of a collision-free Pivot |
| `Partition(input, predicate)`           | `(matched, unmatched)` | split a slice into both halves in one pass |

`Pivot` does **no** aggregation: colliding `(row, col)` values overwrite
last-write-wins. To combine colliding values, `GroupBy` + `Aggregate` first,
then `Pivot` the aggregated rows. `Unpivot` reverses a collision-free `Pivot`
exactly (round-trip), but emits cells in unspecified order (sort if you need
stability).

## Integration

`Aggregate` is intentionally the shape of the `stats` reducers — any
`func([]T) (R, bool)` (`stats.Sum`, `stats.Mean`, your own) is a valid
`AggregateFunc`. For grouping that needs richer key→values structure
(set-valued / ordered multimaps, concurrent access) reach for
`collections/multimaps`; `relational` stays on native maps so its results
compose directly with the `maps` and `slices` packages.
