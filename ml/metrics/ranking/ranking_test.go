package ranking_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/ml/metrics/ranking"
)

const epsilon = 1e-9

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}

func TestDCG(t *testing.T) {
	rels := []float64{3, 2, 3, 0, 1, 2}
	got, ok := ranking.DCG(rels, 0)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if !approxEqual(got, 6.861126688593501) {
		t.Errorf("got %v, want %v", got, 6.861126688593501)
	}

	// Empty and non-finite inputs are undefined.
	if _, ok := ranking.DCG(nil, 0); ok {
		t.Error("empty: ok = true, want false")
	}
	if _, ok := ranking.DCG([]float64{1, math.NaN()}, 0); ok {
		t.Error("NaN: ok = true, want false")
	}
}

func TestNDCG(t *testing.T) {
	trueRelevance := []float64{3, 2, 3, 0, 1, 2}
	scores := []float64{6, 5, 4, 3, 2, 1} // identity order

	tests := []struct {
		name string
		k    int
		want float64
	}{
		{"full list", 0, 0.9608081943360616},
		{"cutoff at 3", 3, 0.9777813616305049},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ranking.NDCG(trueRelevance, scores, tt.k)
			if !ok {
				t.Fatal("ok = false, want true")
			}
			if !approxEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNDCGPerfectRankingIsOne(t *testing.T) {
	// Scores already agree with relevance order, so NDCG is 1.
	got, ok := ranking.NDCG([]float64{3, 2, 1}, []float64{0.9, 0.5, 0.1}, 0)
	if !ok || !approxEqual(got, 1) {
		t.Errorf("got %v %v, want 1 true", got, ok)
	}
}

func TestNDCGRejectsBadInput(t *testing.T) {
	if _, ok := ranking.NDCG([]float64{}, []float64{}, 0); ok {
		t.Error("empty: ok = true, want false")
	}
	if _, ok := ranking.NDCG([]float64{1, 2}, []float64{1}, 0); ok {
		t.Error("length mismatch: ok = true, want false")
	}
	// All-zero relevance has zero ideal DCG: NDCG is undefined.
	if _, ok := ranking.NDCG([]float64{0, 0, 0}, []float64{3, 2, 1}, 0); ok {
		t.Error("zero relevance: ok = true, want false")
	}
	// Non-finite relevance and non-finite score are both rejected.
	if _, ok := ranking.NDCG([]float64{1, math.Inf(1)}, []float64{2, 1}, 0); ok {
		t.Error("non-finite relevance: ok = true, want false")
	}
	if _, ok := ranking.NDCG([]float64{1, 2}, []float64{2, math.NaN()}, 0); ok {
		t.Error("non-finite score: ok = true, want false")
	}
}

func TestAveragePrecision(t *testing.T) {
	got, ok := ranking.AveragePrecision([]bool{true, false, true, false, true})
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if !approxEqual(got, 0.7555555555555555) {
		t.Errorf("got %v, want %v", got, 0.7555555555555555)
	}

	// All relevant items at the top scores a perfect 1.
	perfect, ok := ranking.AveragePrecision([]bool{true, true, false, false})
	if !ok || !approxEqual(perfect, 1) {
		t.Errorf("got %v %v, want 1 true", perfect, ok)
	}

	// No relevant items is undefined.
	if _, ok := ranking.AveragePrecision([]bool{false, false}); ok {
		t.Error("no relevant: ok = true, want false")
	}
	if _, ok := ranking.AveragePrecision(nil); ok {
		t.Error("empty: ok = true, want false")
	}
}

func TestMeanAveragePrecision(t *testing.T) {
	queries := [][]bool{
		{true, false, true, false, true},
		{false, true, true, false},
	}
	got, ok := ranking.MeanAveragePrecision(queries)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if !approxEqual(got, 0.6694444444444444) {
		t.Errorf("got %v, want %v", got, 0.6694444444444444)
	}

	// A degenerate query (no relevant items) makes the whole mean undefined.
	if _, ok := ranking.MeanAveragePrecision([][]bool{{true}, {false}}); ok {
		t.Error("degenerate query: ok = true, want false")
	}
	if _, ok := ranking.MeanAveragePrecision(nil); ok {
		t.Error("empty: ok = true, want false")
	}
}

func TestDCGNoMutation(t *testing.T) {
	rels := []float64{3, 1, 2}
	snapshot := append([]float64(nil), rels...)
	_, _ = ranking.NDCG(rels, []float64{0.1, 0.9, 0.5}, 0)
	for i := range rels {
		if rels[i] != snapshot[i] {
			t.Fatal("NDCG mutated its input")
		}
	}
}
