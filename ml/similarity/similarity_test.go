package similarity_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sets"
	"github.com/pickeringtech/go-collections/ml/similarity"
)

const floatTol = 1e-9

func floatsClose(a, b float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return math.IsNaN(a) && math.IsNaN(b)
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return a == b
	}
	return math.Abs(a-b) <= floatTol
}

func TestDotProduct(t *testing.T) {
	type args struct {
		a, b []float64
	}
	tests := []struct {
		name string
		args args
		want float64
		ok   bool
	}{
		{
			name: "basic dot product",
			args: args{a: []float64{1, 2, 3}, b: []float64{4, 5, 6}},
			want: 32,
			ok:   true,
		},
		{
			name: "orthogonal vectors",
			args: args{a: []float64{1, 0}, b: []float64{0, 1}},
			want: 0,
			ok:   true,
		},
		{
			name: "length mismatch is undefined",
			args: args{a: []float64{1, 2, 3}, b: []float64{1, 2}},
			want: 0,
			ok:   false,
		},
		{
			name: "empty is undefined",
			args: args{a: []float64{}, b: []float64{}},
			want: 0,
			ok:   false,
		},
		{
			name: "nil is undefined",
			args: args{a: nil, b: nil},
			want: 0,
			ok:   false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := similarity.DotProduct(tc.args.a, tc.args.b)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if ok && !floatsClose(got, tc.want) {
				t.Fatalf("DotProduct() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCosineSimilarity(t *testing.T) {
	type args struct {
		a, b []float64
	}
	tests := []struct {
		name string
		args args
		want float64
		ok   bool
	}{
		{
			name: "identical vectors",
			args: args{a: []float64{1, 2, 3}, b: []float64{1, 2, 3}},
			want: 1,
			ok:   true,
		},
		{
			name: "anti-parallel vectors",
			args: args{a: []float64{1, 0}, b: []float64{-1, 0}},
			want: -1,
			ok:   true,
		},
		{
			name: "orthogonal vectors",
			args: args{a: []float64{1, 0}, b: []float64{0, 1}},
			want: 0,
			ok:   true,
		},
		{
			name: "scaled vector same direction",
			args: args{a: []float64{1, 2, 3}, b: []float64{2, 4, 6}},
			want: 1,
			ok:   true,
		},
		{
			name: "length mismatch is undefined",
			args: args{a: []float64{1, 2}, b: []float64{1, 2, 3}},
			want: 0,
			ok:   false,
		},
		{
			name: "zero vector a is undefined",
			args: args{a: []float64{0, 0}, b: []float64{1, 2}},
			want: 0,
			ok:   false,
		},
		{
			name: "zero vector b is undefined",
			args: args{a: []float64{1, 2}, b: []float64{0, 0}},
			want: 0,
			ok:   false,
		},
		{
			name: "empty is undefined",
			args: args{a: []float64{}, b: []float64{}},
			want: 0,
			ok:   false,
		},
		{
			name: "nil is undefined",
			args: args{a: nil, b: nil},
			want: 0,
			ok:   false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := similarity.CosineSimilarity(tc.args.a, tc.args.b)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if ok && !floatsClose(got, tc.want) {
				t.Fatalf("CosineSimilarity() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCosineSimilarityNaNPropagates(t *testing.T) {
	a := []float64{1, math.NaN(), 3}
	b := []float64{4, 5, 6}
	got, ok := similarity.CosineSimilarity(a, b)
	if !ok {
		t.Fatalf("ok = false, want true (NaN propagates with ok == true)")
	}
	if !math.IsNaN(got) {
		t.Fatalf("CosineSimilarity() = %v, want NaN", got)
	}
}

func TestCosineSimilaritySymmetric(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}

	ab, okAB := similarity.CosineSimilarity(a, b)
	ba, okBA := similarity.CosineSimilarity(b, a)
	if !okAB || !okBA {
		t.Fatalf("ok = %v / %v, want both true", okAB, okBA)
	}
	if !floatsClose(ab, ba) {
		t.Fatalf("CosineSimilarity(a,b) = %v != CosineSimilarity(b,a) = %v", ab, ba)
	}
}

func TestJaccard(t *testing.T) {
	type args struct {
		a, b sets.Set[string]
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "identical sets",
			args: args{a: sets.NewHash("a", "b", "c"), b: sets.NewHash("a", "b", "c")},
			want: 1,
		},
		{
			name: "disjoint sets",
			args: args{a: sets.NewHash("a", "b"), b: sets.NewHash("c", "d")},
			want: 0,
		},
		{
			name: "partial overlap",
			args: args{a: sets.NewHash("a", "b", "c", "d"), b: sets.NewHash("b", "c", "d", "e")},
			want: 3.0 / 5.0, // |{b,c,d}| / |{a,b,c,d,e}|
		},
		{
			name: "subset is less than one",
			args: args{a: sets.NewHash("a", "b"), b: sets.NewHash("a", "b", "c")},
			want: 2.0 / 3.0,
		},
		{
			name: "both empty returns zero",
			args: args{a: sets.NewHash[string](), b: sets.NewHash[string]()},
			want: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := similarity.Jaccard(tc.args.a, tc.args.b)
			if !floatsClose(got, tc.want) {
				t.Errorf("Jaccard() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestJaccardSymmetric(t *testing.T) {
	a := sets.NewHash("a", "b", "c")
	b := sets.NewHash("b", "c", "d")
	ab := similarity.Jaccard(a, b)
	ba := similarity.Jaccard(b, a)
	if !floatsClose(ab, ba) {
		t.Fatalf("Jaccard(a,b) = %v != Jaccard(b,a) = %v", ab, ba)
	}
}

func TestDice(t *testing.T) {
	type args struct {
		a, b sets.Set[string]
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "identical sets",
			args: args{a: sets.NewHash("a", "b", "c"), b: sets.NewHash("a", "b", "c")},
			want: 1,
		},
		{
			name: "disjoint sets",
			args: args{a: sets.NewHash("a", "b"), b: sets.NewHash("c", "d")},
			want: 0,
		},
		{
			name: "partial overlap",
			args: args{a: sets.NewHash("a", "b", "c", "d"), b: sets.NewHash("b", "c", "d", "e")},
			want: 2.0 * 3.0 / (4.0 + 4.0), // 2*|{b,c,d}| / (4+4)
		},
		{
			name: "both empty returns zero",
			args: args{a: sets.NewHash[string](), b: sets.NewHash[string]()},
			want: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := similarity.Dice(tc.args.a, tc.args.b)
			if !floatsClose(got, tc.want) {
				t.Errorf("Dice() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestOverlap(t *testing.T) {
	type args struct {
		a, b sets.Set[string]
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "identical sets",
			args: args{a: sets.NewHash("a", "b", "c"), b: sets.NewHash("a", "b", "c")},
			want: 1,
		},
		{
			name: "subset returns one",
			args: args{a: sets.NewHash("a", "b"), b: sets.NewHash("a", "b", "c")},
			want: 1,
		},
		{
			name: "disjoint sets",
			args: args{a: sets.NewHash("a", "b"), b: sets.NewHash("c", "d")},
			want: 0,
		},
		{
			name: "partial overlap",
			args: args{a: sets.NewHash("a", "b", "c", "d"), b: sets.NewHash("b", "c", "d", "e")},
			want: 3.0 / 4.0,
		},
		{
			name: "both empty returns zero",
			args: args{a: sets.NewHash[string](), b: sets.NewHash[string]()},
			want: 0,
		},
		{
			name: "one empty returns zero",
			args: args{a: sets.NewHash[string](), b: sets.NewHash("a", "b")},
			want: 0,
		},
		{
			name: "b smaller than a",
			args: args{a: sets.NewHash("a", "b", "c"), b: sets.NewHash("a", "b")},
			want: 1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := similarity.Overlap(tc.args.a, tc.args.b)
			if !floatsClose(got, tc.want) {
				t.Errorf("Overlap() = %v, want %v", got, tc.want)
			}
		})
	}
}

// Verify that the functions are not mutating their inputs by checking
// slice identity before and after calls.
func TestCosineSimilarityDoesNotMutateInput(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}
	aSnap := []float64{1, 2, 3}
	bSnap := []float64{4, 5, 6}

	_, _ = similarity.CosineSimilarity(a, b)

	if !reflect.DeepEqual(a, aSnap) {
		t.Fatalf("CosineSimilarity mutated a: got %v, want %v", a, aSnap)
	}
	if !reflect.DeepEqual(b, bSnap) {
		t.Fatalf("CosineSimilarity mutated b: got %v, want %v", b, bSnap)
	}
}
