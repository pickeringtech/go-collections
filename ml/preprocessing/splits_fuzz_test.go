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

// FuzzStratifiedSplit asserts StratifiedSplit's guarantees for arbitrary label
// layouts: train and test are a disjoint cover of the input; the test size is
// exactly TrainTestSplit's; every label with two or more rows has a test row
// when the test size allows one each, and a train row when the train size
// does; and the same seed reproduces the split.
func FuzzStratifiedSplit(f *testing.F) {
	f.Add([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, uint8(128), int64(1))
	f.Add([]byte{0, 0, 0, 0, 0, 0, 1, 1, 1, 1}, uint8(77), int64(2))
	f.Add([]byte{0, 1, 1, 2, 2, 2, 2, 2, 2, 2}, uint8(200), int64(3))

	f.Fuzz(func(t *testing.T, rawLabels []byte, rawFrac uint8, seed int64) {
		if len(rawLabels) == 0 || len(rawLabels) > 512 {
			return
		}
		labels := make([]byte, len(rawLabels))
		sizes := map[byte]int{}
		for i, b := range rawLabels {
			labels[i] = b % 8 // a small alphabet makes repeated labels likely
			sizes[labels[i]]++
		}
		input := benchInts(len(labels))
		frac := float64(rawFrac) / 255

		train, test, ok := preprocessing.StratifiedSplit(input, labels, frac, preprocessing.NewRand(seed))
		if !ok {
			t.Fatalf("ok = false for valid frac %v", frac)
		}
		combined := append(append([]int{}, train...), test...)
		sort.Ints(combined)
		if !reflect.DeepEqual(combined, input) {
			t.Fatalf("train+test is not a permutation of the input")
		}
		_, wantTest, _ := preprocessing.TrainTestSplit(input, frac, preprocessing.NewRand(seed))
		if len(test) != len(wantTest) {
			t.Fatalf("len(test) = %d, want TrainTestSplit's %d", len(test), len(wantTest))
		}

		multi := 0
		for _, size := range sizes {
			if size >= 2 {
				multi++
			}
		}
		inTest := map[byte]int{}
		for _, idx := range test {
			inTest[labels[idx]]++
		}
		for label, size := range sizes {
			if size < 2 {
				continue
			}
			if len(test) >= multi && inTest[label] == 0 {
				t.Fatalf("label %d (%d rows) has no test row though test has %d >= %d", label, size, len(test), multi)
			}
			if len(train) >= multi && inTest[label] == size {
				t.Fatalf("label %d (%d rows) has no train row though train has %d >= %d", label, size, len(train), multi)
			}
		}

		train2, test2, _ := preprocessing.StratifiedSplit(input, labels, frac, preprocessing.NewRand(seed))
		if !reflect.DeepEqual(train, train2) || !reflect.DeepEqual(test, test2) {
			t.Fatalf("same seed produced different splits")
		}
	})
}
