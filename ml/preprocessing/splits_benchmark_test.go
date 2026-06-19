package preprocessing_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func BenchmarkShuffle(b *testing.B) {
	for _, n := range benchSizes {
		data := benchInts(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = preprocessing.Shuffle(data, preprocessing.NewRand(1))
			}
		})
	}
}

func BenchmarkTrainTestSplit(b *testing.B) {
	for _, n := range benchSizes {
		data := benchInts(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _ = preprocessing.TrainTestSplit(data, 0.2, preprocessing.NewRand(1))
			}
		})
	}
}

func BenchmarkKFold(b *testing.B) {
	for _, n := range benchSizes {
		data := benchInts(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.KFold(data, 5, preprocessing.NewRand(1))
			}
		})
	}
}

func BenchmarkStratifiedSplit(b *testing.B) {
	for _, n := range benchSizes {
		data := benchInts(n)
		labels := make([]int, n)
		for i := range labels {
			labels[i] = i % 4
		}
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _ = preprocessing.StratifiedSplit(data, labels, 0.2, preprocessing.NewRand(1))
			}
		})
	}
}
