package preprocessing_test

import (
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
