package regression_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/metrics/regression"
)

func benchPair(n int) (yTrue, yPred []float64) {
	yTrue = make([]float64, n)
	yPred = make([]float64, n)
	for i := range yTrue {
		yTrue[i] = float64(i%100) + 1
		yPred[i] = float64(i%100) + 1.5 // a constant 0.5 offset
	}
	return yTrue, yPred
}

func BenchmarkMeanSquaredError(b *testing.B) {
	yTrue, yPred := benchPair(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = regression.MeanSquaredError(yTrue, yPred)
	}
}

func BenchmarkMeanAbsoluteError(b *testing.B) {
	yTrue, yPred := benchPair(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = regression.MeanAbsoluteError(yTrue, yPred)
	}
}

func BenchmarkRSquared(b *testing.B) {
	yTrue, yPred := benchPair(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = regression.RSquared(yTrue, yPred)
	}
}
