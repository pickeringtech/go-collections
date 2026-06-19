package relational_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/relational"
)

type record struct {
	row string
	col string
	val int
}

func TestPivot(t *testing.T) {
	rowKey := func(r record) string { return r.row }
	colKey := func(r record) string { return r.col }
	value := func(r record) int { return r.val }

	tests := []struct {
		name string
		rows []record
		want map[string]map[string]int
	}{
		{
			name: "nil input yields non-nil empty map",
			rows: nil,
			want: map[string]map[string]int{},
		},
		{
			name: "empty input yields non-nil empty map",
			rows: []record{},
			want: map[string]map[string]int{},
		},
		{
			name: "reshapes long to wide",
			rows: []record{
				{"r1", "c1", 1}, {"r1", "c2", 2}, {"r2", "c1", 3},
			},
			want: map[string]map[string]int{
				"r1": {"c1": 1, "c2": 2},
				"r2": {"c1": 3},
			},
		},
		{
			name: "collision is last-write-wins",
			rows: []record{
				{"r1", "c1", 1}, {"r1", "c1", 99},
			},
			want: map[string]map[string]int{
				"r1": {"c1": 99},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := relational.Pivot(tt.rows, rowKey, colKey, value)
			if got == nil {
				t.Fatalf("Pivot returned nil map")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pivot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnpivot(t *testing.T) {
	tests := []struct {
		name string
		wide map[string]map[string]int
		want []relational.Cell[string, string, int]
	}{
		{
			name: "nil input yields non-nil empty slice",
			wide: nil,
			want: []relational.Cell[string, string, int]{},
		},
		{
			name: "empty input yields non-nil empty slice",
			wide: map[string]map[string]int{},
			want: []relational.Cell[string, string, int]{},
		},
		{
			name: "flattens wide to long",
			wide: map[string]map[string]int{
				"r1": {"c1": 1, "c2": 2},
				"r2": {"c1": 3},
			},
			want: []relational.Cell[string, string, int]{
				{Row: "r1", Col: "c1", Value: 1},
				{Row: "r1", Col: "c2", Value: 2},
				{Row: "r2", Col: "c1", Value: 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := relational.Unpivot(tt.wide)
			if got == nil {
				t.Fatalf("Unpivot returned nil slice")
			}
			// Map iteration order is undefined, so sort before comparing.
			sortCells(got)
			sortCells(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unpivot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func sortCells(cells []relational.Cell[string, string, int]) {
	sort.Slice(cells, func(i, j int) bool {
		if cells[i].Row != cells[j].Row {
			return cells[i].Row < cells[j].Row
		}
		return cells[i].Col < cells[j].Col
	})
}
