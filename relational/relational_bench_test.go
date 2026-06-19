package relational_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/relational"
	"github.com/pickeringtech/go-collections/stats"
)

// ladder is the shared element-count matrix from
// agent-os/standards/testing/benchmark-scaling.md: every benchmark sub-benchmarks
// across these sizes via b.Run so the linear scaling is visible.
var ladder = []struct {
	name string
	n    int
}{
	{"3 elements", 3},
	{"10 elements", 10},
	{"100 elements", 100},
	{"1_000 elements", 1_000},
	{"10_000 elements", 10_000},
	{"100_000 elements", 100_000},
	{"1_000_000 elements", 1_000_000},
}

// benchKey spreads values over roughly n/4 keys so most groups carry several
// values, exercising the grouped (rather than one-per-key) shape.
func benchKey(n int) func(int) int {
	keys := n/4 + 1
	return func(v int) int { return v % keys }
}

func intSlice(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

func BenchmarkGroupBy(b *testing.B) {
	for _, bm := range ladder {
		input := intSlice(bm.n)
		key := benchKey(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = relational.GroupBy(input, key)
			}
		})
	}
}

func BenchmarkGroupBySeq(b *testing.B) {
	for _, bm := range ladder {
		input := intSlice(bm.n)
		key := benchKey(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = relational.GroupBySeq(sliceSeq(input), key)
			}
		})
	}
}

func BenchmarkCountBy(b *testing.B) {
	for _, bm := range ladder {
		input := intSlice(bm.n)
		key := benchKey(bm.n)
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = relational.CountBy(input, key)
			}
		})
	}
}

func BenchmarkAggregate(b *testing.B) {
	for _, bm := range ladder {
		groups := relational.GroupBy(intSlice(bm.n), benchKey(bm.n))
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = relational.Aggregate(groups, stats.Sum[int])
			}
		})
	}
}

func BenchmarkAggregateBy(b *testing.B) {
	for _, bm := range ladder {
		groups := relational.GroupBy(intSlice(bm.n), benchKey(bm.n))
		project := func(n int) float64 { return float64(n) }
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = relational.AggregateBy(groups, project, stats.Mean[float64])
			}
		})
	}
}

func BenchmarkInnerJoin(b *testing.B) {
	for _, bm := range ladder {
		left := intSlice(bm.n)
		right := intSlice(bm.n)
		id := func(n int) int { return n }
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = relational.InnerJoin(left, right, id, id)
			}
		})
	}
}

func BenchmarkLeftJoin(b *testing.B) {
	for _, bm := range ladder {
		left := intSlice(bm.n)
		right := intSlice(bm.n)
		id := func(n int) int { return n }
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = relational.LeftJoin(left, right, id, id)
			}
		})
	}
}

func BenchmarkRightJoin(b *testing.B) {
	for _, bm := range ladder {
		left := intSlice(bm.n)
		right := intSlice(bm.n)
		id := func(n int) int { return n }
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = relational.RightJoin(left, right, id, id)
			}
		})
	}
}

func BenchmarkFullOuterJoin(b *testing.B) {
	for _, bm := range ladder {
		left := intSlice(bm.n)
		right := intSlice(bm.n)
		id := func(n int) int { return n }
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = relational.FullOuterJoin(left, right, id, id)
			}
		})
	}
}

func BenchmarkPivot(b *testing.B) {
	for _, bm := range ladder {
		type cell struct {
			row, col, val int
		}
		rows := make([]cell, bm.n)
		cols := bm.n/4 + 1
		for i := range rows {
			rows[i] = cell{row: i % cols, col: i % 3, val: i}
		}
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = relational.Pivot(rows,
					func(c cell) int { return c.row },
					func(c cell) int { return c.col },
					func(c cell) int { return c.val },
				)
			}
		})
	}
}

func BenchmarkUnpivot(b *testing.B) {
	for _, bm := range ladder {
		wide := map[int]map[int]int{}
		cols := bm.n/4 + 1
		for i := 0; i < bm.n; i++ {
			r := i % cols
			inner, ok := wide[r]
			if !ok {
				inner = map[int]int{}
				wide[r] = inner
			}
			inner[i%3] = i
		}
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = relational.Unpivot(wide)
			}
		})
	}
}

func BenchmarkPartition(b *testing.B) {
	for _, bm := range ladder {
		input := intSlice(bm.n)
		pred := func(n int) bool { return n%2 == 0 }
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = relational.Partition(input, pred)
			}
		})
	}
}
