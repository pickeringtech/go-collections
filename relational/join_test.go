package relational_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/relational"
)

type lrow struct {
	key int
	v   string
}

type rrow struct {
	key int
	v   string
}

func lkey(l lrow) int { return l.key }
func rkey(r rrow) int { return r.key }

type joinFn func([]lrow, []rrow, func(lrow) int, func(rrow) int) []relational.JoinPair[lrow, rrow]

func TestInnerJoin(t *testing.T) {
	tests := []struct {
		name  string
		join  joinFn
		left  []lrow
		right []rrow
		want  []relational.JoinPair[lrow, rrow]
	}{
		{
			name:  "nil inputs yield non-nil empty result",
			join:  relational.InnerJoin[int, lrow, rrow],
			left:  nil,
			right: nil,
			want:  []relational.JoinPair[lrow, rrow]{},
		},
		{
			name:  "no matches yields empty result",
			join:  relational.InnerJoin[int, lrow, rrow],
			left:  []lrow{{1, "a"}},
			right: []rrow{{2, "b"}},
			want:  []relational.JoinPair[lrow, rrow]{},
		},
		{
			name:  "one-to-one match",
			join:  relational.InnerJoin[int, lrow, rrow],
			left:  []lrow{{1, "a"}},
			right: []rrow{{1, "x"}},
			want: []relational.JoinPair[lrow, rrow]{
				{Left: lrow{1, "a"}, Right: rrow{1, "x"}, LeftOK: true, RightOK: true},
			},
		},
		{
			name:  "many-to-many cross product",
			join:  relational.InnerJoin[int, lrow, rrow],
			left:  []lrow{{1, "a"}, {1, "b"}},
			right: []rrow{{1, "x"}, {1, "y"}},
			want: []relational.JoinPair[lrow, rrow]{
				{Left: lrow{1, "a"}, Right: rrow{1, "x"}, LeftOK: true, RightOK: true},
				{Left: lrow{1, "a"}, Right: rrow{1, "y"}, LeftOK: true, RightOK: true},
				{Left: lrow{1, "b"}, Right: rrow{1, "x"}, LeftOK: true, RightOK: true},
				{Left: lrow{1, "b"}, Right: rrow{1, "y"}, LeftOK: true, RightOK: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.join(tt.left, tt.right, lkey, rkey)
			assertPairs(t, got, tt.want)
		})
	}
}

func TestLeftJoin(t *testing.T) {
	tests := []struct {
		name  string
		join  joinFn
		left  []lrow
		right []rrow
		want  []relational.JoinPair[lrow, rrow]
	}{
		{
			name:  "nil inputs yield non-nil empty result",
			join:  relational.LeftJoin[int, lrow, rrow],
			left:  nil,
			right: nil,
			want:  []relational.JoinPair[lrow, rrow]{},
		},
		{
			name:  "unmatched left row emitted with zero right",
			join:  relational.LeftJoin[int, lrow, rrow],
			left:  []lrow{{1, "a"}, {2, "b"}},
			right: []rrow{{1, "x"}},
			want: []relational.JoinPair[lrow, rrow]{
				{Left: lrow{1, "a"}, Right: rrow{1, "x"}, LeftOK: true, RightOK: true},
				{Left: lrow{2, "b"}, Right: rrow{}, LeftOK: true, RightOK: false},
			},
		},
		{
			name:  "many-to-many for matched key",
			join:  relational.LeftJoin[int, lrow, rrow],
			left:  []lrow{{1, "a"}},
			right: []rrow{{1, "x"}, {1, "y"}},
			want: []relational.JoinPair[lrow, rrow]{
				{Left: lrow{1, "a"}, Right: rrow{1, "x"}, LeftOK: true, RightOK: true},
				{Left: lrow{1, "a"}, Right: rrow{1, "y"}, LeftOK: true, RightOK: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.join(tt.left, tt.right, lkey, rkey)
			assertPairs(t, got, tt.want)
		})
	}
}

func TestRightJoin(t *testing.T) {
	tests := []struct {
		name  string
		join  joinFn
		left  []lrow
		right []rrow
		want  []relational.JoinPair[lrow, rrow]
	}{
		{
			name:  "nil inputs yield non-nil empty result",
			join:  relational.RightJoin[int, lrow, rrow],
			left:  nil,
			right: nil,
			want:  []relational.JoinPair[lrow, rrow]{},
		},
		{
			name:  "unmatched right row emitted with zero left",
			join:  relational.RightJoin[int, lrow, rrow],
			left:  []lrow{{1, "a"}},
			right: []rrow{{1, "x"}, {2, "y"}},
			want: []relational.JoinPair[lrow, rrow]{
				{Left: lrow{1, "a"}, Right: rrow{1, "x"}, LeftOK: true, RightOK: true},
				{Left: lrow{}, Right: rrow{2, "y"}, LeftOK: false, RightOK: true},
			},
		},
		{
			name:  "many-to-many for matched key",
			join:  relational.RightJoin[int, lrow, rrow],
			left:  []lrow{{1, "a"}, {1, "b"}},
			right: []rrow{{1, "x"}},
			want: []relational.JoinPair[lrow, rrow]{
				{Left: lrow{1, "a"}, Right: rrow{1, "x"}, LeftOK: true, RightOK: true},
				{Left: lrow{1, "b"}, Right: rrow{1, "x"}, LeftOK: true, RightOK: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.join(tt.left, tt.right, lkey, rkey)
			assertPairs(t, got, tt.want)
		})
	}
}

func TestFullOuterJoin(t *testing.T) {
	tests := []struct {
		name  string
		join  joinFn
		left  []lrow
		right []rrow
		want  []relational.JoinPair[lrow, rrow]
	}{
		{
			name:  "nil inputs yield non-nil empty result",
			join:  relational.FullOuterJoin[int, lrow, rrow],
			left:  nil,
			right: nil,
			want:  []relational.JoinPair[lrow, rrow]{},
		},
		{
			name:  "matched, unmatched left, unmatched right",
			join:  relational.FullOuterJoin[int, lrow, rrow],
			left:  []lrow{{1, "a"}, {2, "b"}},
			right: []rrow{{1, "x"}, {3, "z"}},
			want: []relational.JoinPair[lrow, rrow]{
				{Left: lrow{1, "a"}, Right: rrow{1, "x"}, LeftOK: true, RightOK: true},
				{Left: lrow{2, "b"}, Right: rrow{}, LeftOK: true, RightOK: false},
				{Left: lrow{}, Right: rrow{3, "z"}, LeftOK: false, RightOK: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.join(tt.left, tt.right, lkey, rkey)
			assertPairs(t, got, tt.want)
		})
	}
}

// assertPairs requires got be a non-nil slice equal to want.
func assertPairs(t *testing.T, got, want []relational.JoinPair[lrow, rrow]) {
	t.Helper()
	if got == nil {
		t.Fatalf("join returned nil slice")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("join() = %v, want %v", got, want)
	}
}

func TestJoinsDoNotMutateInput(t *testing.T) {
	left := []lrow{{1, "a"}, {2, "b"}}
	right := []rrow{{1, "x"}}
	leftSnap := []lrow{{1, "a"}, {2, "b"}}
	rightSnap := []rrow{{1, "x"}}
	_ = relational.FullOuterJoin(left, right, lkey, rkey)
	if !reflect.DeepEqual(left, leftSnap) || !reflect.DeepEqual(right, rightSnap) {
		t.Errorf("join mutated input slices")
	}
}
