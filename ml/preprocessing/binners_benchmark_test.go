package preprocessing_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func BenchmarkFixedWidthBinner(b *testing.B) {
	for _, n := range benchSizes {
		data := benchFloats(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.NewFixedWidthBinner(10).FitTransform(data)
			}
		})
	}
}

func BenchmarkQuantileBinner(b *testing.B) {
	for _, n := range benchSizes {
		data := benchFloats(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.NewQuantileBinner(10).FitTransform(data)
			}
		})
	}
}
