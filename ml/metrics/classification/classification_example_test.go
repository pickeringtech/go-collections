package classification_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/ml/metrics/classification"
)

// Example_quickStart is the runnable twin of the package godoc overview. Keep
// the two in sync: `go test` compiles and output-checks this, which guarantees
// the documented classification API actually exists and behaves as shown.
func Example_quickStart() {
	yTrue := []int{0, 1, 2, 2, 1, 0, 1, 2}
	yPred := []int{0, 2, 1, 2, 1, 0, 1, 2}

	acc, _ := classification.Accuracy(yTrue, yPred)
	f1, _ := classification.F1(yTrue, yPred, classification.Macro)
	cm, _ := classification.ConfusionMatrix(yTrue, yPred)

	labels := []int{0, 0, 1, 1}
	scores := []float64{0.1, 0.4, 0.35, 0.8}
	auc, _ := classification.AUC(labels, scores, 1)

	fmt.Printf("acc=%.2f macroF1=%.4f hits(2,2)=%d auc=%.2f",
		acc, f1, cm.Count(2, 2), auc)
	// Output: acc=0.75 macroF1=0.7778 hits(2,2)=2 auc=0.75
}

func ExampleAccuracy() {
	acc, ok := classification.Accuracy(
		[]string{"cat", "dog", "cat", "bird"},
		[]string{"cat", "dog", "bird", "bird"},
	)
	fmt.Printf("%.2f %v", acc, ok)
	// Output: 0.75 true
}

func ExampleF1() {
	yTrue := []int{0, 1, 2, 2, 1, 0, 1, 2}
	yPred := []int{0, 2, 1, 2, 1, 0, 1, 2}

	macro, _ := classification.F1(yTrue, yPred, classification.Macro)
	micro, _ := classification.F1(yTrue, yPred, classification.Micro)

	// Micro F1 equals the accuracy for single-label classification.
	fmt.Printf("macro=%.4f micro=%.2f", macro, micro)
	// Output: macro=0.7778 micro=0.75
}

func ExampleLogLoss() {
	loss, ok := classification.LogLoss(
		[]int{1, 0, 1, 1},
		[]float64{0.9, 0.1, 0.8, 0.7},
		1,
	)
	fmt.Printf("%.4f %v", loss, ok)
	// Output: 0.1976 true
}
