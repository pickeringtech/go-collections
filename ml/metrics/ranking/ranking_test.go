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
	_, ok = ranking.DCG(nil, 0)
	if ok {
		t.Error("empty: ok = true, want false")
	}
	_, ok = ranking.DCG([]float64{1, math.NaN()}, 0)
	if ok {
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
	_, ok := ranking.NDCG([]float64{}, []float64{}, 0)
	if ok {
		t.Error("empty: ok = true, want false")
	}
	_, ok = ranking.NDCG([]float64{1, 2}, []float64{1}, 0)
	if ok {
		t.Error("length mismatch: ok = true, want false")
	}
	// All-zero relevance has zero ideal DCG: NDCG is undefined.
	_, ok = ranking.NDCG([]float64{0, 0, 0}, []float64{3, 2, 1}, 0)
	if ok {
		t.Error("zero relevance: ok = true, want false")
	}
	// Non-finite relevance and non-finite score are both rejected.
	_, ok = ranking.NDCG([]float64{1, math.Inf(1)}, []float64{2, 1}, 0)
	if ok {
		t.Error("non-finite relevance: ok = true, want false")
	}
	_, ok = ranking.NDCG([]float64{1, 2}, []float64{2, math.NaN()}, 0)
	if ok {
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
	_, ok = ranking.AveragePrecision([]bool{false, false})
	if ok {
		t.Error("no relevant: ok = true, want false")
	}
	_, ok = ranking.AveragePrecision(nil)
	if ok {
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
	_, ok = ranking.MeanAveragePrecision([][]bool{{true}, {false}})
	if ok {
		t.Error("degenerate query: ok = true, want false")
	}
	_, ok = ranking.MeanAveragePrecision(nil)
	if ok {
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

// TestNDCGRejectsNegativeRelevance pins NDCG to its documented [0, 1] range:
// with a negative gain the ideal DCG no longer bounds the ranked DCG (NDCG of
// {1, −1} ranked worst-first was −1), so, like scikit-learn, negative
// relevance is rejected.
func TestNDCGRejectsNegativeRelevance(t *testing.T) {
	tests := []struct {
		name      string
		relevance []float64
		scores    []float64
	}{
		{name: "one negative gain", relevance: []float64{1, -1}, scores: []float64{1, 2}},
		{name: "all negative gains", relevance: []float64{-1, -2, -3}, scores: []float64{3, 2, 1}},
		{name: "negative gain beyond the cutoff", relevance: []float64{3, 2, -1}, scores: []float64{3, 2, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ranking.NDCG(tt.relevance, tt.scores, 2)
			if ok {
				t.Errorf("NDCG() = %v, true; want ok == false", got)
			}
		})
	}
}

// TestNDCGTiesKeepInputOrder pins the documented tie-break: items with equal
// scores are ranked in their input order.
func TestNDCGTiesKeepInputOrder(t *testing.T) {
	// Both items score 1, so the less relevant first item is ranked first:
	// DCG = 1 + 3/log2(3), ideal = 3 + 1/log2(3).
	got, ok := ranking.NDCG([]float64{1, 3}, []float64{1, 1}, 0)
	want := (1 + 3/math.Log2(3)) / (3 + 1/math.Log2(3))
	if !ok || !approxEqual(got, want) {
		t.Errorf("NDCG() = %v, %v; want %v, true", got, ok, want)
	}
}

// TestDCGIsCompensated guards the discounted sum against round-off: each later
// relevance is chosen so its discounted gain is exactly 1, which a naive running
// sum after a gain of 1e17 would absorb and lose.
func TestDCGIsCompensated(t *testing.T) {
	relevances := make([]float64, 17)
	relevances[0] = 1e17
	for i := 1; i < len(relevances); i++ {
		relevances[i] = math.Log2(float64(i + 2)) // gain / discount == 1 exactly
	}
	got, ok := ranking.DCG(relevances, 0)
	if !ok || got != 1e17+16 {
		t.Errorf("DCG() = %v, %v; want %v, true", got, ok, float64(1e17+16))
	}
}
