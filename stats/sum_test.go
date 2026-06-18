package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		want   int
		wantOK bool
	}{
		{
			name:   "sums positive and negative values",
			input:  []int{1, 2, -1, 3, 4, 5},
			want:   14,
			wantOK: true,
		},
		{
			name:   "single value",
			input:  []int{42},
			want:   42,
			wantOK: true,
		},
		{
			name:   "empty input is undefined",
			input:  []int{},
			want:   0,
			wantOK: false,
		},
		{
			name:   "nil input is undefined",
			input:  nil,
			want:   0,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.Sum(tt.input)
			if got != tt.want || ok != tt.wantOK {
				t.Errorf("Sum() = (%v, %v), want (%v, %v)", got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestSumFloat(t *testing.T) {
	got, ok := stats.Sum([]float64{0.1, 0.2, 0.3})
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if !approxEqual(got, 0.6) {
		t.Errorf("Sum() = %v, want ~0.6", got)
	}
}

// TestSumNonFinitePropagates documents the exact-in-T tier's policy: unlike the
// float64 summaries, Sum does not reject non-finite values — it propagates them
// per IEEE arithmetic and still reports ok=true for non-empty input.
func TestSumNonFinitePropagates(t *testing.T) {
	got, ok := stats.Sum([]float64{1, math.NaN(), 3})
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if !math.IsNaN(got) {
		t.Errorf("Sum() = %v, want NaN", got)
	}
}

func FuzzSum(f *testing.F) {
	f.Add(1, 2, 3)
	f.Add(0, 0, 0)
	f.Fuzz(func(t *testing.T, a, b, c int) {
		input := []int{a, b, c}
		sum, ok := stats.Sum(input)
		if !ok {
			t.Fatalf("ok = false for non-empty input %v", input)
		}
		// Sum is order-independent and equals the naive accumulation.
		if sum != a+b+c {
			t.Errorf("Sum() = %v, want %v", sum, a+b+c)
		}
		if rev, _ := stats.Sum([]int{c, b, a}); rev != sum {
			t.Errorf("Sum not order-independent: %v vs %v", sum, rev)
		}
	})
}
