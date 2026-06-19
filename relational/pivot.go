package relational

// Cell is a single long-format observation: a row key, a column key and the
// value at their intersection. It is the unit Unpivot emits and the natural
// shape of tidy/long data — one row per (row, column, value) — which Pivot
// reshapes into wide form and Unpivot reverses.
type Cell[K comparable, K2 comparable, V any] struct {
	Row   K
	Col   K2
	Value V
}

// Pivot reshapes long-format rows into a wide nested map indexed first by row
// key then by column key. It is the spreadsheet "pivot table" / SQL crosstab
// move: rows that share a rowKey become one map entry, and each contributes a
// column keyed by colKey holding value. This turns a flat list of observations
// into a row-by-column grid you can index directly as wide[row][col].
//
// Collisions are last-write-wins: if two input rows map to the same
// (rowKey, colKey) the later one (in input order) overwrites the earlier, just
// as repeated map assignment would. This keeps Pivot a pure reshape with no
// hidden aggregation; if you need to combine colliding values, GroupBy +
// Aggregate them first, then Pivot the aggregated rows. Unpivot reverses a
// collision-free Pivot exactly (round-trip).
//
// The input slice is never mutated. Empty or nil input yields a non-nil empty
// outer map; every inner map created is likewise non-nil.
func Pivot[R any, K comparable, K2 comparable, V any](rows []R, rowKey func(R) K, colKey func(R) K2, value func(R) V) map[K]map[K2]V {
	wide := map[K]map[K2]V{}
	for _, row := range rows {
		r := rowKey(row)
		inner, ok := wide[r]
		if !ok {
			inner = map[K2]V{}
			wide[r] = inner
		}
		inner[colKey(row)] = value(row) // last write wins on (row, col) collision
	}
	return wide
}

// Unpivot flattens a wide nested map back into long-format Cells, one per
// (row, column) entry — the inverse of Pivot. Use it to round-trip wide data
// back to tidy form, or to feed a grid into the row-oriented GroupBy/Join
// pipeline.
//
// Cell order is not specified, because Go map iteration order is randomised;
// callers that need a stable order should sort the result. A nil or empty input
// yields a non-nil empty slice.
func Unpivot[K comparable, K2 comparable, V any](wide map[K]map[K2]V) []Cell[K, K2, V] {
	cells := []Cell[K, K2, V]{}
	for row, inner := range wide {
		for col, v := range inner {
			cells = append(cells, Cell[K, K2, V]{Row: row, Col: col, Value: v})
		}
	}
	return cells
}
