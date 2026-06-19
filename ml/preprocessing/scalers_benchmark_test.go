package preprocessing_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

var benchSizes = []int{3, 10, 100, 1_000, 10_000, 100_000, 1_000_000}

func BenchmarkStandardScaler(b *testing.B) {
	for _, n := range benchSizes {
		data := benchFloats(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.NewStandardScaler().FitTransform(data)
			}
		})
	}
}

func BenchmarkMinMaxScaler(b *testing.B) {
	for _, n := range benchSizes {
		data := benchFloats(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.NewMinMaxScaler().FitTransform(data)
			}
		})
	}
}

func BenchmarkRobustScaler(b *testing.B) {
	for _, n := range benchSizes {
		data := benchFloats(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.NewRobustScaler().FitTransform(data)
			}
		})
	}
}
