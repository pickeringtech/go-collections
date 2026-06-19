package classification_test

import (
	"testing"

	"github.com/pickeringtech/go-collections/ml/metrics/classification"
)

func benchLabels(n int) (yTrue, yPred []int) {
	yTrue = make([]int, n)
	yPred = make([]int, n)
	for i := range yTrue {
		yTrue[i] = i % 5
		yPred[i] = (i + i/5) % 5 // a deterministic mix of hits and misses
	}
	return yTrue, yPred
}

func benchScores(n int) (yTrue []int, scores []float64) {
	yTrue = make([]int, n)
	scores = make([]float64, n)
	for i := range yTrue {
		yTrue[i] = i % 2
		scores[i] = float64(i%100) / 100
	}
	return yTrue, scores
}

func BenchmarkAccuracy(b *testing.B) {
	yTrue, yPred := benchLabels(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = classification.Accuracy(yTrue, yPred)
	}
}

func BenchmarkF1Macro(b *testing.B) {
	yTrue, yPred := benchLabels(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = classification.F1(yTrue, yPred, classification.Macro)
	}
}

func BenchmarkAUC(b *testing.B) {
	yTrue, scores := benchScores(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = classification.AUC(yTrue, scores, 1)
	}
}
