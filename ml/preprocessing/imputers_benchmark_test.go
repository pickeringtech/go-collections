package preprocessing_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/preprocessing"
)

func BenchmarkMeanImputer(b *testing.B) {
	for _, n := range benchSizes {
		data := benchFloats(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.NewMeanImputer(nil).FitTransform(data)
			}
		})
	}
}

func BenchmarkMedianImputer(b *testing.B) {
	for _, n := range benchSizes {
		data := benchFloats(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.NewMedianImputer(nil).FitTransform(data)
			}
		})
	}
}

func BenchmarkModeImputer(b *testing.B) {
	isMissing := func(v string) bool { return v == "" }
	for _, n := range benchSizes {
		data := benchStrings(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.NewModeImputer(isMissing).FitTransform(data)
			}
		})
	}
}

func BenchmarkConstantImputer(b *testing.B) {
	isMissing := func(v string) bool { return v == "" }
	for _, n := range benchSizes {
		data := benchStrings(n)
		b.Run(sizeName(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = preprocessing.NewConstantImputer("x", isMissing).Transform(data)
			}
		})
	}
}
