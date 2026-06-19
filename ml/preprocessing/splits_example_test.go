package preprocessing_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func ExampleTrainTestSplit() {
	// 30% of the ten elements go to the test set; the seed makes it reproducible.
	rng := preprocessing.NewRand(42)
	train, test, ok := preprocessing.TrainTestSplit([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 0.3, rng)
	fmt.Println(len(train), len(test), ok)
	// Output: 7 3 true
}

func ExampleKFold() {
	// Ten elements into three folds of near-equal size: 4, 3, 3.
	folds, ok := preprocessing.KFold([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 3, preprocessing.NewRand(42))
	fmt.Println(len(folds), len(folds[0]), len(folds[1]), len(folds[2]), ok)
	// Output: 3 4 3 3 true
}

func ExampleShuffle() {
	// Shuffle returns a reordered copy; the same seed always gives the same order.
	got := preprocessing.ShuffleSeed([]int{1, 2, 3, 4, 5}, 1)
	fmt.Println(got)
	// Output: [2 1 3 4 5]
}
