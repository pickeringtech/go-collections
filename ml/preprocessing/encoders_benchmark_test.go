package preprocessing_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func BenchmarkOneHotEncoder(b *testing.B) {
	for _, n := range benchSizes {
		data := benchStrings(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.NewOneHotEncoder[string]().FitTransform(data)
			}
		})
	}
}

func BenchmarkLabelEncoder(b *testing.B) {
	for _, n := range benchSizes {
		data := benchStrings(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.NewLabelEncoder[string]().FitTransform(data)
			}
		})
	}
}

func BenchmarkOrdinalEncoder(b *testing.B) {
	for _, n := range benchSizes {
		data := benchStrings(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.NewOrdinalEncoder[string]().FitTransform(data)
			}
		})
	}
}

func BenchmarkTargetEncoder(b *testing.B) {
	for _, n := range benchSizes {
		data := benchStrings(n)
		target := benchFloats(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				enc := preprocessing.NewTargetEncoder[string]().Fit(data, target)
				_, _ = enc.Transform(data)
			}
		})
	}
}
