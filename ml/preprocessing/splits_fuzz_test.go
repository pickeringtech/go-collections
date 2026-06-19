package preprocessing_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

// FuzzTrainTestSplit asserts that, for any seed and fraction, the train and test
// sets are a disjoint cover of the input and that the same seed reproduces the
// split exactly.
func FuzzTrainTestSplit(f *testing.F) {
	f.Add(20, uint8(128), int64(1))
	f.Add(1, uint8(0), int64(7))
	f.Add(100, uint8(255), int64(42))

	f.Fuzz(func(t *testing.T, n int, rawFrac uint8, seed int64) {
		if n < 1 || n > 4096 {
			return
		}
		input := make([]int, n)
		for i := range input {
			input[i] = i
		}
		frac := float64(rawFrac) / 255

		train, test, ok := preprocessing.TrainTestSplit(input, frac, preprocessing.NewRand(seed))
		if !ok {
			t.Fatalf("ok = false for valid frac %v", frac)
		}

		// Disjoint cover: train ∪ test is a permutation of the input.
		combined := append(append([]int{}, train...), test...)
		sort.Ints(combined)
		if !reflect.DeepEqual(combined, input) {
			t.Fatalf("train+test is not a permutation of the input")
		}

		// Reproducible: the same seed yields the same partition.
		train2, test2, _ := preprocessing.TrainTestSplit(input, frac, preprocessing.NewRand(seed))
		if !reflect.DeepEqual(train, train2) || !reflect.DeepEqual(test, test2) {
			t.Fatalf("same seed produced different splits")
		}
	})
}

// FuzzKFold asserts that the folds are a disjoint cover of the input and that
// fold sizes differ by at most one.
func FuzzKFold(f *testing.F) {
	f.Add(10, 3, int64(1))
	f.Add(1, 1, int64(2))
	f.Add(50, 7, int64(9))

	f.Fuzz(func(t *testing.T, n int, k int, seed int64) {
		if n < 1 || n > 4096 || k < 1 || k > n {
			return
		}
		input := make([]int, n)
		for i := range input {
			input[i] = i
		}

		folds, ok := preprocessing.KFold(input, k, preprocessing.NewRand(seed))
		if !ok || len(folds) != k {
			t.Fatalf("KFold = (%d folds, %v), want (%d, true)", len(folds), ok, k)
		}

		var all []int
		minSize, maxSize := n, 0
		for _, fold := range folds {
			all = append(all, fold...)
			if len(fold) < minSize {
				minSize = len(fold)
			}
			if len(fold) > maxSize {
				maxSize = len(fold)
			}
		}
		if maxSize-minSize > 1 {
			t.Fatalf("fold sizes differ by more than one: %d..%d", minSize, maxSize)
		}
		sort.Ints(all)
		if !reflect.DeepEqual(all, input) {
			t.Fatalf("folds are not a disjoint cover of the input")
		}
	})
}
