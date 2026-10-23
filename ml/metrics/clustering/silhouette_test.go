package clustering_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/ml/metrics/clustering"
)

const epsilon = 1e-9

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}

// Two tight, well-separated clusters on a line.
var (
	wellSeparated = [][]float64{{0, 0}, {0.5, 0}, {10, 0}, {10.5, 0}}
	twoClusters   = []int{0, 0, 1, 1}
)

func TestSilhouetteScore(t *testing.T) {
	got, ok := clustering.SilhouetteScore(wellSeparated, twoClusters)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if !approxEqual(got, 0.9499687304565353) {
		t.Errorf("got %v, want %v", got, 0.9499687304565353)
	}
}

func TestSilhouetteSamples(t *testing.T) {
	tests := []struct {
		name   string
		points [][]float64
		labels []int
		want   []float64
	}{
		{
			name:   "two pairs of points",
			points: wellSeparated,
			labels: twoClusters,
			want:   []float64{0.9512195121951219, 0.9487179487179487, 0.9487179487179487, 0.9512195121951219},
		},
		{
			// With pairs, every point's cohesion is a single distance; a cluster
			// of three makes it a true mean of two. For the point at 0:
			// a = (2+4)/2 = 3, b = (10+14)/2 = 12, so s = (12−3)/12 = 0.75.
			name:   "a cluster of three averages its intra-cluster distances",
			points: [][]float64{{0}, {2}, {4}, {10}, {14}},
			labels: []int{0, 0, 0, 1, 1},
			want:   []float64{0.75, 0.8, 0.625, 0.5, 2.0 / 3.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := clustering.SilhouetteSamples(tt.points, tt.labels)
			if !ok {
				t.Fatal("ok = false, want true")
			}
			if len(got) != len(tt.want) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.want))
			}
			for i := range tt.want {
				if !approxEqual(got[i], tt.want[i]) {
					t.Errorf("sample %d = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSilhouetteOverlappingClustersScoreLow(t *testing.T) {
	// Interleaved points get a much lower (here negative) score than the
	// well-separated arrangement.
	points := [][]float64{{0, 0}, {10, 0}, {0.5, 0}, {10.5, 0}}
	labels := []int{0, 1, 1, 0} // deliberately mis-assigned
	got, ok := clustering.SilhouetteScore(points, labels)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if got >= 0 {
		t.Errorf("got %v, want a negative score for mis-assigned clusters", got)
	}
}

func TestSilhouetteSingletonCluster(t *testing.T) {
	// A lone point in its own cluster contributes a silhouette of 0.
	points := [][]float64{{0, 0}, {0.5, 0}, {10, 0}}
	labels := []int{0, 0, 1}
	samples, ok := clustering.SilhouetteSamples(points, labels)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if samples[2] != 0 {
		t.Errorf("singleton sample = %v, want 0", samples[2])
	}
}

func TestSilhouetteCoincidentPoints(t *testing.T) {
	// Every point sits at the same location across both clusters, so a == b == 0
	// and the silhouette of each sample is defined as 0 (max(a,b) == 0).
	points := [][]float64{{0, 0}, {0, 0}, {0, 0}, {0, 0}}
	labels := []int{0, 0, 1, 1}
	got, ok := clustering.SilhouetteScore(points, labels)
	if !ok || got != 0 {
		t.Errorf("got %v %v, want 0 true", got, ok)
	}
}

func TestSilhouetteRejectsBadInput(t *testing.T) {
	tests := []struct {
		name   string
		points [][]float64
		labels []int
	}{
		{"too few points", [][]float64{{0, 0}}, []int{0}},
		{"length mismatch", [][]float64{{0, 0}, {1, 1}}, []int{0}},
		{"single cluster", [][]float64{{0, 0}, {1, 1}}, []int{0, 0}},
		{"every point its own cluster", [][]float64{{0, 0}, {1, 1}}, []int{0, 1}},
		{"ragged rows", [][]float64{{0, 0}, {1}}, []int{0, 1}},
		{"non-finite coordinate", [][]float64{{0, 0}, {math.Inf(1), 0}, {1, 1}}, []int{0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := clustering.SilhouetteScore(tt.points, tt.labels)
			if ok {
				t.Error("ok = true, want false")
			}
		})
	}
}

func TestSilhouetteScoreWithCustomMetric(t *testing.T) {
	// Manhattan (L1) distance as a custom metric.
	manhattan := func(a, b []float64) float64 {
		var sum float64
		for i := range a {
			sum += math.Abs(a[i] - b[i])
		}
		return sum
	}
	got, ok := clustering.SilhouetteScoreWith(wellSeparated, twoClusters, manhattan)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	// On a line, L1 and L2 coincide, so the score matches the Euclidean one.
	if !approxEqual(got, 0.9499687304565353) {
		t.Errorf("got %v, want %v", got, 0.9499687304565353)
	}
}

func TestSilhouetteWithNilMetric(t *testing.T) {
	s, ok := clustering.SilhouetteSamplesWith(wellSeparated, twoClusters, nil)
	if ok || s != nil {
		t.Errorf("SilhouetteSamplesWith(nil) = %v %v, want nil false", s, ok)
	}
	score, ok := clustering.SilhouetteScoreWith(wellSeparated, twoClusters, nil)
	if ok || score != 0 {
		t.Errorf("SilhouetteScoreWith(nil) = %v %v, want 0 false", score, ok)
	}
}

func TestSilhouetteRejectsMisbehavingMetric(t *testing.T) {
	tests := []struct {
		name string
		dist clustering.DistanceFunc
	}{
		{"NaN", func(a, b []float64) float64 { return math.NaN() }},
		{"Inf", func(a, b []float64) float64 { return math.Inf(1) }},
		{"negative", func(a, b []float64) float64 { return -1 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, ok := clustering.SilhouetteSamplesWith(wellSeparated, twoClusters, tt.dist)
			if ok || s != nil {
				t.Errorf("SilhouetteSamplesWith = %v %v, want nil false", s, ok)
			}
			score, ok := clustering.SilhouetteScoreWith(wellSeparated, twoClusters, tt.dist)
			if ok || score != 0 {
				t.Errorf("SilhouetteScoreWith = %v %v, want 0 false", score, ok)
			}
		})
	}
}

func TestSilhouetteRejectsOverflowingMetric(t *testing.T) {
	// A finite, non-negative metric whose accumulated sum overflows to +Inf.
	// Point 0's own cluster has two other members, so MaxFloat64 + MaxFloat64
	// overflows before any mean is taken — which would otherwise produce a NaN
	// coefficient while still reporting ok == true.
	points := [][]float64{{0, 0}, {1, 0}, {2, 0}, {10, 0}, {11, 0}}
	labels := []int{0, 0, 0, 1, 1}
	huge := func(a, b []float64) float64 { return math.MaxFloat64 }

	s, ok := clustering.SilhouetteSamplesWith(points, labels, huge)
	if ok || s != nil {
		t.Errorf("SilhouetteSamplesWith = %v %v, want nil false", s, ok)
	}
	score, ok := clustering.SilhouetteScoreWith(points, labels, huge)
	if ok || score != 0 {
		t.Errorf("SilhouetteScoreWith = %v %v, want 0 false", score, ok)
	}
}

func TestNoMutation(t *testing.T) {
	points := [][]float64{{0, 0}, {0.5, 0}, {10, 0}, {10.5, 0}}
	snapshot := make([][]float64, len(points))
	for i := range points {
		snapshot[i] = append([]float64(nil), points[i]...)
	}
	_, _ = clustering.SilhouetteScore(points, twoClusters)
	for i := range points {
		for j := range points[i] {
			if points[i][j] != snapshot[i][j] {
				t.Fatal("input was mutated")
			}
		}
	}
}
