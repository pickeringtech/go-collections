// Package relational provides the data-engineering relational primitives —
// GroupBy and aggregate, the four joins, pivot/unpivot, and partition — as
// free functions over Go's native slices and maps. It is the "split-apply-
// combine" toolkit: reshape rows into groups, reduce each group, stitch tables
// together by key, and flip between long and wide layouts, without writing the
// bookkeeping loops by hand.
//
// Every function here follows the library's contracts: it never mutates its
// input, and nil or empty input yields a non-nil empty result so callers can
// range or len it without a nil check.
//
// # Quick Start
//
//	import (
//		"github.com/pickeringtech/go-collections/relational"
//		"github.com/pickeringtech/go-collections/stats"
//	)
//
//	type Order struct {
//		Dept   string
//		Amount float64
//	}
//	orders := []Order{
//		{"books", 10}, {"books", 30}, {"toys", 5},
//	}
//
//	// GROUP BY dept, then AVG(amount) — the aggregate pipeline. The grouping
//	// keeps every row; Aggregate reduces each group with any reducer of shape
//	// func([]T) (R, bool), so stats.Mean plugs straight in via AggregateBy.
//	byDept := relational.GroupBy(orders, func(o Order) string { return o.Dept })
//	avg := relational.AggregateBy(byDept,
//		func(o Order) float64 { return o.Amount },
//		stats.Mean[float64],
//	)
//	// avg = map[books:20 toys:5]
//
//	// JOIN orders to a per-dept label table by dept key.
//	type Label struct {
//		Dept string
//		Name string
//	}
//	labels := []Label{{"books", "Books Dept"}, {"toys", "Toys Dept"}}
//	joined := relational.InnerJoin(orders, labels,
//		func(o Order) string { return o.Dept },
//		func(l Label) string { return l.Dept },
//	)
//	// each joined element pairs an Order with its matching Label.
//
// This Quick Start is compiled and run as Example_quickStart in the package's
// test suite, so it is guaranteed to track the real API.
//
// # Core Operations
//
// Grouping and aggregation:
//   - GroupBy / GroupBySeq: partition values into key→[]value buckets (first-seen order).
//   - CountBy: per-key counts, cheaper than GroupBy when you only need sizes.
//   - Aggregate / AggregateBy: reduce each group with a func([]T) (R, bool) reducer.
//
// Joins (all many-to-many, emitting the cross product of matching rows):
//   - InnerJoin: matched rows only.
//   - LeftJoin / RightJoin: matched rows plus unmatched rows from one side.
//   - FullOuterJoin: matched rows plus unmatched rows from both sides.
//
// Reshaping:
//   - Pivot / Unpivot: convert between long-format Cells and a wide nested map.
//   - Partition: split a slice into predicate-matching and non-matching halves.
//
// # Integration
//
// Aggregate is intentionally the shape of the stats reducers: any function of
// shape func([]T) (R, bool) — stats.Sum, stats.Mean, stats.MinMax-style helpers,
// or your own — is a valid AggregateFunc, and a reducer that reports ok==false
// for a group causes that group to be omitted from the result, preserving the
// stats (result, ok) idiom end to end. For grouping that needs richer
// key→values structure (set-valued or ordered multimaps, concurrent access),
// reach for collections/multimaps; relational stays on native maps so its
// results compose directly with the maps and slices packages.
package relational
