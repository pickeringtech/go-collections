package relational_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/relational"
	"github.com/pickeringtech/go-collections/stats"
)

func TestAggregate(t *testing.T) {
	type args struct {
		groups map[string][]int
		aggFn  relational.AggregateFunc[int, int]
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{
			name: "nil groups yields non-nil empty map",
			args: args{groups: nil, aggFn: stats.Sum[int]},
			want: map[string]int{},
		},
		{
			name: "empty groups yields non-nil empty map",
			args: args{groups: map[string][]int{}, aggFn: stats.Sum[int]},
			want: map[string]int{},
		},
		{
			name: "sum per group",
			args: args{
				groups: map[string][]int{"a": {1, 2, 3}, "b": {10, 20}},
				aggFn:  stats.Sum[int],
			},
			want: map[string]int{"a": 6, "b": 30},
		},
		{
			name: "group with ok==false is omitted",
			args: args{
				// stats.Sum returns ok==false for an empty slice, so group "b"
				// is dropped while "a" survives.
				groups: map[string][]int{"a": {1, 2}, "b": {}},
				aggFn:  stats.Sum[int],
			},
			want: map[string]int{"a": 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := relational.Aggregate(tt.args.groups, tt.args.aggFn)
			if got == nil {
				t.Fatalf("Aggregate returned nil map")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

type order struct {
	dept   string
	amount float64
}

func TestAggregateBy(t *testing.T) {
	project := func(o order) float64 { return o.amount }
	type args struct {
		groups map[string][]order
	}
	tests := []struct {
		name  string
		args  args
		aggFn func([]float64) (float64, bool)
		want  map[string]float64
	}{
		{
			name:  "nil groups yields non-nil empty map",
			args:  args{groups: nil},
			aggFn: stats.Mean[float64],
			want:  map[string]float64{},
		},
		{
			name:  "empty groups yields non-nil empty map",
			args:  args{groups: map[string][]order{}},
			aggFn: stats.Mean[float64],
			want:  map[string]float64{},
		},
		{
			name: "mean of projected field per group",
			args: args{groups: map[string][]order{
				"books": {{"books", 10}, {"books", 30}},
				"toys":  {{"toys", 5}},
			}},
			aggFn: stats.Mean[float64],
			want:  map[string]float64{"books": 20, "toys": 5},
		},
		{
			name: "group with ok==false is omitted",
			args: args{groups: map[string][]order{
				"books": {{"books", 10}, {"books", 30}},
				"empty": {},
			}},
			aggFn: stats.Mean[float64],
			want:  map[string]float64{"books": 20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := relational.AggregateBy(tt.args.groups, project, tt.aggFn)
			if got == nil {
				t.Fatalf("AggregateBy returned nil map")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregateDoesNotMutateInput(t *testing.T) {
	groups := map[string][]int{"a": {1, 2, 3}}
	_ = relational.Aggregate(groups, stats.Sum[int])
	if !reflect.DeepEqual(groups, map[string][]int{"a": {1, 2, 3}}) {
		t.Errorf("Aggregate mutated input groups: %v", groups)
	}
}
