package preprocessing_test

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

// sortedCopy returns a sorted copy of s, for multiset comparison.
func sortedCopy(s []int) []int {
	out := make([]int, len(s))
	copy(out, s)
	sort.Ints(out)
	return out
}

func TestTrainTestSplit(t *testing.T) {
	input := benchInts(10)
	train, test, ok := preprocessing.TrainTestSplit(input, 0.3, preprocessing.NewRand(1))
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	if len(test) != 3 || len(train) != 7 {
		t.Fatalf("sizes train=%d test=%d, want 7/3", len(train), len(test))
	}
	// Train and test are a disjoint cover of the input.
	combined := append(append([]int{}, train...), test...)
	if !reflect.DeepEqual(sortedCopy(combined), benchInts(10)) {
		t.Fatalf("train+test = %v, want a permutation of 0..9", sortedCopy(combined))
	}
}

func TestTrainTestSplitDeterministic(t *testing.T) {
	input := benchInts(20)
	tr1, te1, _ := preprocessing.TrainTestSplit(input, 0.25, preprocessing.NewRand(7))
	tr2, te2, _ := preprocessing.TrainTestSplitSeed(input, 0.25, 7)
	if !reflect.DeepEqual(tr1, tr2) || !reflect.DeepEqual(te1, te2) {
		t.Fatalf("same seed produced different splits")
	}
}

func TestTrainTestSplitFractions(t *testing.T) {
	input := benchInts(5)
	train, test, _ := preprocessing.TrainTestSplit(input, 0, preprocessing.NewRand(1))
	if len(train) != 5 || len(test) != 0 {
		t.Fatalf("frac 0: train=%d test=%d, want 5/0", len(train), len(test))
	}
	train, test, _ = preprocessing.TrainTestSplit(input, 1, preprocessing.NewRand(1))
	if len(train) != 0 || len(test) != 5 {
		t.Fatalf("frac 1: train=%d test=%d, want 0/5", len(train), len(test))
	}
}

func TestTrainTestSplitInvalid(t *testing.T) {
	if _, _, ok := preprocessing.TrainTestSplit([]int{}, 0.5, nil); ok {
		t.Fatalf("empty input reported ok")
	}
	if _, _, ok := preprocessing.TrainTestSplit([]int{1, 2}, -0.1, nil); ok {
		t.Fatalf("negative frac reported ok")
	}
	if _, _, ok := preprocessing.TrainTestSplit([]int{1, 2}, 1.1, nil); ok {
		t.Fatalf("frac > 1 reported ok")
	}
}

func TestKFold(t *testing.T) {
	input := benchInts(10)
	folds, ok := preprocessing.KFold(input, 3, preprocessing.NewRand(1))
	if !ok || len(folds) != 3 {
		t.Fatalf("KFold = (%d folds, %v), want (3, true)", len(folds), ok)
	}
	// Fold sizes are near-equal: 4, 3, 3.
	sizes := []int{len(folds[0]), len(folds[1]), len(folds[2])}
	if !reflect.DeepEqual(sizes, []int{4, 3, 3}) {
		t.Fatalf("fold sizes = %v, want [4 3 3]", sizes)
	}
	// The folds are a disjoint cover of the input.
	var all []int
	for _, f := range folds {
		all = append(all, f...)
	}
	if !reflect.DeepEqual(sortedCopy(all), benchInts(10)) {
		t.Fatalf("folds union = %v, want 0..9", sortedCopy(all))
	}
}

func TestKFoldInvalid(t *testing.T) {
	if _, ok := preprocessing.KFold([]int{}, 2, nil); ok {
		t.Fatalf("empty input reported ok")
	}
	if _, ok := preprocessing.KFold([]int{1, 2}, 0, nil); ok {
		t.Fatalf("k=0 reported ok")
	}
	if _, ok := preprocessing.KFold([]int{1, 2}, 3, nil); ok {
		t.Fatalf("k>n reported ok")
	}
}

func TestStratifiedSplitPreservesProportions(t *testing.T) {
	// 6 "A" and 4 "B"; with frac 0.5 each class contributes half to test.
	input := make([]int, 10)
	labels := make([]string, 10)
	for i := range input {
		input[i] = i
		if i < 6 {
			labels[i] = "A"
		} else {
			labels[i] = "B"
		}
	}
	train, test, ok := preprocessing.StratifiedSplit(input, labels, 0.5, preprocessing.NewRand(3))
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	// Count class membership in test.
	countA, countB := 0, 0
	for _, v := range test {
		if v < 6 {
			countA++
		} else {
			countB++
		}
	}
	if countA != 3 || countB != 2 {
		t.Fatalf("test class counts A=%d B=%d, want 3/2", countA, countB)
	}
	if len(train)+len(test) != 10 {
		t.Fatalf("train+test = %d, want 10", len(train)+len(test))
	}
}

func TestStratifiedSplitInvalid(t *testing.T) {
	if _, _, ok := preprocessing.StratifiedSplit([]int{1, 2}, []string{"A"}, 0.5, nil); ok {
		t.Fatalf("mismatched lengths reported ok")
	}
	if _, _, ok := preprocessing.StratifiedSplit([]int{}, []string{}, 0.5, nil); ok {
		t.Fatalf("empty input reported ok")
	}
}

func TestShuffleIsPermutation(t *testing.T) {
	input := benchInts(50)
	got := preprocessing.Shuffle(input, preprocessing.NewRand(9))
	if reflect.DeepEqual(got, input) {
		t.Fatalf("Shuffle returned the input order (vanishingly unlikely)")
	}
	if !reflect.DeepEqual(sortedCopy(got), benchInts(50)) {
		t.Fatalf("Shuffle is not a permutation of the input")
	}
}

func TestShuffleDeterministic(t *testing.T) {
	input := benchInts(20)
	a := preprocessing.Shuffle(input, preprocessing.NewRand(5))
	b := preprocessing.ShuffleSeed(input, 5)
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("same seed produced different shuffles")
	}
}

func TestShuffleZeroValueOutput(t *testing.T) {
	got := preprocessing.Shuffle[int](nil, preprocessing.NewRand(1))
	if got == nil || len(got) != 0 {
		t.Fatalf("Shuffle(nil) = %v, want non-nil empty", got)
	}
}

func TestSplitsDoNotMutateInput(t *testing.T) {
	input := benchInts(10)
	original := benchInts(10)
	preprocessing.Shuffle(input, preprocessing.NewRand(1))
	preprocessing.TrainTestSplit(input, 0.3, preprocessing.NewRand(1))
	preprocessing.KFold(input, 3, preprocessing.NewRand(1))
	if !reflect.DeepEqual(input, original) {
		t.Fatalf("input mutated to %v", input)
	}
}

// TestStratifiedSplitRejectsNaNLabels guards against NaN labels being dropped:
// NaN is unequal to itself, so each NaN would form its own unmatchable group
// and its row would silently vanish from both train and test.
func TestStratifiedSplitRejectsNaNLabels(t *testing.T) {
	nan := math.NaN()
	type point struct{ x, y float64 }
	tests := []struct {
		name string
		fn   func() bool
	}{
		{name: "float label", fn: func() bool {
			_, _, ok := preprocessing.StratifiedSplit([]int{1, 2, 3, 4}, []float64{1, nan, 2, 1}, 0.5, preprocessing.NewRand(1))
			return ok
		}},
		{name: "every label NaN", fn: func() bool {
			_, _, ok := preprocessing.StratifiedSplit([]int{1, 2}, []float64{nan, nan}, 0.5, preprocessing.NewRand(1))
			return ok
		}},
		{name: "struct label carrying a NaN", fn: func() bool {
			_, _, ok := preprocessing.StratifiedSplit([]int{1, 2}, []point{{1, 2}, {nan, 2}}, 0.5, preprocessing.NewRand(1))
			return ok
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fn() {
				t.Errorf("ok = true, want false")
			}
		})
	}
}

// labelsOf repeats each label by its count, in order: labelsOf("A", 2, "B", 1)
// is [A A B].
func labelsOf(spec ...any) []string {
	var out []string
	for i := 0; i < len(spec); i += 2 {
		for j := 0; j < spec[i+1].(int); j++ {
			out = append(out, spec[i].(string))
		}
	}
	return out
}

// testCountsByLabel tallies how many of each label landed in test, where test
// holds indices into labels.
func testCountsByLabel(labels []string, test []int) map[string]int {
	counts := map[string]int{}
	for _, idx := range test {
		counts[labels[idx]]++
	}
	return counts
}

// TestStratifiedSplitTotalMatchesTrainTestSplit pins the overall test size: it
// is round(testFrac·n), exactly what TrainTestSplit draws, however the rows
// divide into labels. Rounding each label's share independently instead put
// all ten singleton labels in test at 0.5, and none at 0.3.
func TestStratifiedSplitTotalMatchesTrainTestSplit(t *testing.T) {
	singletons := labelsOf("a", 1, "b", 1, "c", 1, "d", 1, "e", 1, "f", 1, "g", 1, "h", 1, "i", 1, "j", 1)
	tests := []struct {
		name   string
		labels []string
		frac   float64
	}{
		{name: "ten singleton labels at one half", labels: singletons, frac: 0.5},
		{name: "ten singleton labels at three tenths", labels: singletons, frac: 0.3},
		{name: "three equal labels at one half", labels: labelsOf("A", 3, "B", 3, "C", 3), frac: 0.5},
		{name: "skewed labels at fifteen percent", labels: labelsOf("A", 7, "B", 2, "C", 1), frac: 0.15},
		{name: "everything in test", labels: labelsOf("A", 3, "B", 2), frac: 1},
		{name: "nothing in test", labels: labelsOf("A", 3, "B", 2), frac: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := benchInts(len(tt.labels))
			train, test, ok := preprocessing.StratifiedSplit(input, tt.labels, tt.frac, preprocessing.NewRand(7))
			_, wantTest, _ := preprocessing.TrainTestSplit(input, tt.frac, preprocessing.NewRand(7))
			if !ok || len(test) != len(wantTest) || len(train)+len(test) != len(input) {
				t.Errorf("StratifiedSplit = %d train / %d test, %v; want %d test, true", len(train), len(test), ok, len(wantTest))
			}
		})
	}
}

// TestStratifiedSplitApportionment checks the per-label test counts: a
// proportional largest-remainder share, adjusted so every label with two or
// more rows reaches both sides when the totals allow.
func TestStratifiedSplitApportionment(t *testing.T) {
	tests := []struct {
		name   string
		labels []string
		frac   float64
		want   map[string]int
	}{
		{name: "exact proportional shares", labels: labelsOf("A", 6, "B", 4), frac: 0.5, want: map[string]int{"A": 3, "B": 2}},
		{name: "largest remainder takes the leftover row", labels: labelsOf("A", 5, "B", 4, "C", 1), frac: 0.3, want: map[string]int{"A": 2, "B": 1, "C": 0}},
		{name: "a tied remainder still covers both multi-row labels", labels: labelsOf("A", 7, "B", 2, "C", 1), frac: 0.15, want: map[string]int{"A": 1, "B": 1, "C": 0}},
		{name: "a small label still reaches test", labels: labelsOf("A", 98, "B", 2), frac: 0.1, want: map[string]int{"A": 9, "B": 1}},
		{name: "a small label still keeps a train row", labels: labelsOf("A", 2, "B", 10), frac: 0.75, want: map[string]int{"A": 1, "B": 8}},
		{name: "the most over-allocated label donates the test row", labels: labelsOf("A", 60, "B", 40, "C", 2), frac: 0.1, want: map[string]int{"A": 5, "B": 4, "C": 1}},
		{name: "the most under-allocated label receives the test row", labels: labelsOf("A", 2, "B", 10, "C", 9), frac: 0.8, want: map[string]int{"A": 1, "B": 8, "C": 8}},
		{name: "too few test rows to cover every label", labels: labelsOf("A", 2, "B", 2, "C", 2), frac: 0.3, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, test, ok := preprocessing.StratifiedSplit(benchInts(len(tt.labels)), tt.labels, tt.frac, preprocessing.NewRand(3))
			if !ok {
				t.Fatalf("ok = false, want true")
			}
			got := testCountsByLabel(tt.labels, test)
			if tt.want == nil {
				return
			}
			for label, want := range tt.want {
				if got[label] != want {
					t.Errorf("label %s has %d test rows, want %d (all: %v)", label, got[label], want, got)
				}
			}
		})
	}
}

// TestStratifiedSplitDeterministic checks the same seed reproduces the split
// exactly, ties and adjustments included.
func TestStratifiedSplitDeterministic(t *testing.T) {
	labels := labelsOf("a", 1, "b", 1, "c", 1, "d", 1, "A", 4, "B", 2)
	input := benchInts(len(labels))
	trainA, testA, _ := preprocessing.StratifiedSplitSeed(input, labels, 0.4, 11)
	trainB, testB, _ := preprocessing.StratifiedSplit(input, labels, 0.4, preprocessing.NewRand(11))
	if !reflect.DeepEqual(trainA, trainB) || !reflect.DeepEqual(testA, testB) {
		t.Errorf("same seed produced different splits: %v/%v vs %v/%v", trainA, testA, trainB, testB)
	}
}
