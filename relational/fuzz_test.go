package relational_test

import (
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/relational"
)

// FuzzGroupBy checks the multiset round-trip oracle: flattening all group
// slices must reproduce the input as a multiset (same elements with the same
// multiplicities), because GroupBy partitions without dropping or duplicating.
func FuzzGroupBy(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{})
	f.Add([]byte{1, 2, 3, 4, 5})

	f.Fuzz(func(t *testing.T, data []byte) {
		input := make([]uint8, len(data))
		copy(input, data)
		key := func(b uint8) uint8 { return b % 4 }

		groups := relational.GroupBy(input, key)

		// Every element lands in the bucket its key selects, and nowhere else.
		flat := []uint8{}
		for k, vs := range groups {
			for _, v := range vs {
				if key(v) != k {
					t.Fatalf("value %d in wrong bucket %d", v, k)
				}
			}
			flat = append(flat, vs...)
		}

		// Flattened groups must equal the input as a multiset.
		gotSorted := append([]uint8(nil), flat...)
		wantSorted := append([]uint8(nil), input...)
		sort.Slice(gotSorted, func(i, j int) bool { return gotSorted[i] < gotSorted[j] })
		sort.Slice(wantSorted, func(i, j int) bool { return wantSorted[i] < wantSorted[j] })
		if len(gotSorted) != len(wantSorted) {
			t.Fatalf("flattened length = %d, want %d", len(gotSorted), len(wantSorted))
		}
		for i := range gotSorted {
			if gotSorted[i] != wantSorted[i] {
				t.Fatalf("multiset mismatch at %d: got %d, want %d", i, gotSorted[i], wantSorted[i])
			}
		}
	})
}

// FuzzCountBy checks that CountBy agrees with GroupBy: the count for each key
// equals the length of that key's GroupBy slice.
func FuzzCountBy(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{1, 2, 3, 4, 5, 6})

	f.Fuzz(func(t *testing.T, data []byte) {
		input := make([]uint8, len(data))
		copy(input, data)
		key := func(b uint8) uint8 { return b % 5 }

		counts := relational.CountBy(input, key)
		groups := relational.GroupBy(input, key)

		if len(counts) != len(groups) {
			t.Fatalf("CountBy key count = %d, GroupBy key count = %d", len(counts), len(groups))
		}
		for k, vs := range groups {
			if counts[k] != len(vs) {
				t.Fatalf("CountBy[%d] = %d, want len(group) = %d", k, counts[k], len(vs))
			}
		}
	})
}

// FuzzPivotUnpivot checks the round-trip oracle: unpivoting a collision-free
// pivot reproduces the original (row, col, value) observations as a set. The
// input is consumed in (row, col, value) triples and de-duplicated on
// (row, col) so the pivot is collision-free.
func FuzzPivotUnpivot(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{0, 0, 1, 0, 1, 2, 1, 0, 3})

	f.Fuzz(func(t *testing.T, data []byte) {
		type triple struct {
			row, col, val uint8
		}
		// Build collision-free observations keyed by (row, col).
		seen := map[[2]uint8]bool{}
		rows := []triple{}
		for i := 0; i+2 < len(data); i += 3 {
			rc := [2]uint8{data[i], data[i+1]}
			if seen[rc] {
				continue
			}
			seen[rc] = true
			rows = append(rows, triple{data[i], data[i+1], data[i+2]})
		}

		wide := relational.Pivot(rows,
			func(tr triple) uint8 { return tr.row },
			func(tr triple) uint8 { return tr.col },
			func(tr triple) uint8 { return tr.val },
		)
		cells := relational.Unpivot(wide)

		if len(cells) != len(rows) {
			t.Fatalf("round-trip cell count = %d, want %d", len(cells), len(rows))
		}
		// Compare as sets keyed on (row, col).
		want := map[[2]uint8]uint8{}
		for _, r := range rows {
			want[[2]uint8{r.row, r.col}] = r.val
		}
		for _, c := range cells {
			rc := [2]uint8{c.Row, c.Col}
			v, ok := want[rc]
			if !ok || v != c.Value {
				t.Fatalf("round-trip cell (%d,%d)=%d not in original (ok=%v want=%d)", c.Row, c.Col, c.Value, ok, v)
			}
		}
	})
}

// FuzzInnerJoin is a differential test against a naive O(n*m) double loop: the
// fast hash-indexed InnerJoin must emit the same pairs (as a multiset) as the
// reference cross product.
func FuzzInnerJoin(f *testing.F) {
	f.Add([]byte(nil), []byte(nil))
	f.Add([]byte{1, 1, 2}, []byte{1, 2, 2})

	f.Fuzz(func(t *testing.T, ldata, rdata []byte) {
		left := make([]uint8, len(ldata))
		copy(left, ldata)
		right := make([]uint8, len(rdata))
		copy(right, rdata)
		key := func(b uint8) uint8 { return b % 3 }

		got := relational.InnerJoin(left, right, key, key)

		// Naive reference: every (l, r) where keys match.
		type pair struct {
			l, r uint8
		}
		wantCounts := map[pair]int{}
		wantTotal := 0
		for _, l := range left {
			for _, r := range right {
				if key(l) == key(r) {
					wantCounts[pair{l, r}]++
					wantTotal++
				}
			}
		}

		if len(got) != wantTotal {
			t.Fatalf("InnerJoin emitted %d pairs, naive reference %d", len(got), wantTotal)
		}
		gotCounts := map[pair]int{}
		for _, p := range got {
			if !p.LeftOK || !p.RightOK {
				t.Fatalf("InnerJoin emitted a non-matched pair: %+v", p)
			}
			gotCounts[pair{p.Left, p.Right}]++
		}
		for p, n := range wantCounts {
			if gotCounts[p] != n {
				t.Fatalf("pair (%d,%d) count = %d, want %d", p.l, p.r, gotCounts[p], n)
			}
		}
	})
}
