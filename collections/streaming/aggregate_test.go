package streaming_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/collections/streaming"
	"github.com/pickeringtech/go-collections/stats"
)

const floatTolerance = 1e-9

func TestRunningMean_ResultAndCount(t *testing.T) {
	type args struct {
		values []float64
	}
	tests := []struct {
		name    string
		args    args
		wantVal float64
		wantOK  bool
		wantN   int
	}{
		{name: "empty stream not ok", args: args{values: nil}, wantVal: 0, wantOK: false, wantN: 0},
		{name: "single element", args: args{values: []float64{42}}, wantVal: 42, wantOK: true, wantN: 1},
		{name: "even numbers", args: args{values: []float64{2, 4, 6, 8}}, wantVal: 5, wantOK: true, wantN: 4},
		{name: "negatives", args: args{values: []float64{-1, -2, -3}}, wantVal: -2, wantOK: true, wantN: 3},
		{name: "duplicates", args: args{values: []float64{3, 3, 3, 3}}, wantVal: 3, wantOK: true, wantN: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := streaming.NewRunningMean()
			for _, v := range tt.args.values {
				m.Add(v)
			}
			gotVal, gotOK := m.Result()
			if gotOK != tt.wantOK {
				t.Fatalf("Result() ok = %v, want %v", gotOK, tt.wantOK)
			}
			if gotOK && math.Abs(gotVal-tt.wantVal) > floatTolerance {
				t.Errorf("Result() = %v, want %v", gotVal, tt.wantVal)
			}
			if m.Count() != tt.wantN {
				t.Errorf("Count() = %d, want %d", m.Count(), tt.wantN)
			}
		})
	}
}

// TestRunningMean_MatchesStatsOracle cross-checks the streaming mean against the
// batch stats.Mean over the same data.
func TestRunningMean_MatchesStatsOracle(t *testing.T) {
	data := []float64{3.5, -2.1, 9.0, 0.0, 4.4, 4.4, -100.0, 7.7}
	m := streaming.NewRunningMean()
	for i := range data {
		m.Add(data[i])
		got, ok := m.Result()
		want, wantOK := stats.Mean(data[:i+1])
		if ok != wantOK {
			t.Fatalf("after %d adds: ok = %v, stats ok = %v", i+1, ok, wantOK)
		}
		if math.Abs(got-want) > 1e-6 {
			t.Errorf("after %d adds: Result() = %v, stats.Mean = %v", i+1, got, want)
		}
	}
}

func TestRunningVariance_Contracts(t *testing.T) {
	type args struct {
		values []float64
	}
	tests := []struct {
		name       string
		args       args
		wantSample bool
		wantPop    bool
		wantMean   bool
	}{
		{name: "empty: nothing defined", args: args{values: nil}, wantSample: false, wantPop: false, wantMean: false},
		{name: "single: pop ok, sample not", args: args{values: []float64{5}}, wantSample: false, wantPop: true, wantMean: true},
		{name: "two: all ok", args: args{values: []float64{5, 7}}, wantSample: true, wantPop: true, wantMean: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := streaming.NewRunningVariance()
			for _, x := range tt.args.values {
				v.Add(x)
			}
			if _, ok := v.SampleVariance(); ok != tt.wantSample {
				t.Errorf("SampleVariance() ok = %v, want %v", ok, tt.wantSample)
			}
			if _, ok := v.PopulationVariance(); ok != tt.wantPop {
				t.Errorf("PopulationVariance() ok = %v, want %v", ok, tt.wantPop)
			}
			if _, ok := v.Mean(); ok != tt.wantMean {
				t.Errorf("Mean() ok = %v, want %v", ok, tt.wantMean)
			}
			if v.Count() != len(tt.args.values) {
				t.Errorf("Count() = %d, want %d", v.Count(), len(tt.args.values))
			}
		})
	}
}

// TestRunningVariance_MatchesStatsOracle cross-checks the streaming variance and
// mean against stats.SampleVariance / stats.PopulationVariance / stats.Mean over
// the same data, fed incrementally.
func TestRunningVariance_MatchesStatsOracle(t *testing.T) {
	data := []float64{2, 4, 4, 4, 5, 5, 7, 9, -3.3, 100.25}
	v := streaming.NewRunningVariance()
	for i := range data {
		v.Add(data[i])
		prefix := data[:i+1]

		gotSample, gotSampleOK := v.SampleVariance()
		wantSample, wantSampleOK := stats.SampleVariance(prefix)
		if gotSampleOK != wantSampleOK {
			t.Fatalf("after %d adds: sample ok = %v, stats ok = %v", i+1, gotSampleOK, wantSampleOK)
		}
		if gotSampleOK && math.Abs(gotSample-wantSample) > 1e-6 {
			t.Errorf("after %d adds: SampleVariance() = %v, stats = %v", i+1, gotSample, wantSample)
		}

		gotPop, gotPopOK := v.PopulationVariance()
		wantPop, wantPopOK := stats.PopulationVariance(prefix)
		if gotPopOK != wantPopOK {
			t.Fatalf("after %d adds: pop ok = %v, stats ok = %v", i+1, gotPopOK, wantPopOK)
		}
		if gotPopOK && math.Abs(gotPop-wantPop) > 1e-6 {
			t.Errorf("after %d adds: PopulationVariance() = %v, stats = %v", i+1, gotPop, wantPop)
		}

		gotMean, _ := v.Mean()
		wantMean, _ := stats.Mean(prefix)
		if math.Abs(gotMean-wantMean) > 1e-6 {
			t.Errorf("after %d adds: Mean() = %v, stats = %v", i+1, gotMean, wantMean)
		}
	}
}

func TestEWMA_Result(t *testing.T) {
	type args struct {
		alpha  float64
		values []float64
	}
	tests := []struct {
		name    string
		args    args
		wantVal float64
		wantOK  bool
	}{
		{name: "no values: not ok", args: args{alpha: 0.5, values: nil}, wantVal: 0, wantOK: false},
		{name: "single value primes", args: args{alpha: 0.5, values: []float64{10}}, wantVal: 10, wantOK: true},
		{name: "alpha 0.5 over three", args: args{alpha: 0.5, values: []float64{10, 20, 30}}, wantVal: 22.5, wantOK: true},
		{name: "alpha 1 tracks latest", args: args{alpha: 1, values: []float64{10, 20, 30}}, wantVal: 30, wantOK: true},
		{name: "alpha clamped from above", args: args{alpha: 5, values: []float64{10, 20}}, wantVal: 20, wantOK: true},
		{name: "alpha clamped from zero primes only", args: args{alpha: 0, values: []float64{10}}, wantVal: 10, wantOK: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := streaming.NewEWMA(tt.args.alpha)
			for _, v := range tt.args.values {
				e.Add(v)
			}
			gotVal, gotOK := e.Result()
			if gotOK != tt.wantOK {
				t.Fatalf("Result() ok = %v, want %v", gotOK, tt.wantOK)
			}
			if gotOK && math.Abs(gotVal-tt.wantVal) > floatTolerance {
				t.Errorf("Result() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

// TestEWMA_AlphaClampedIntoUnitInterval asserts that any out-of-range alpha
// still yields a usable, finite average — never NaN/Inf or a degenerate series.
func TestEWMA_AlphaClampedIntoUnitInterval(t *testing.T) {
	alphas := []float64{-1, 0, math.NaN(), 1.5, math.Inf(1)}
	for _, a := range alphas {
		e := streaming.NewEWMA(a)
		for _, v := range []float64{1, 2, 3, 4} {
			e.Add(v)
		}
		got, ok := e.Result()
		if !ok {
			t.Errorf("alpha %v: Result() not ok after adds", a)
		}
		if math.IsNaN(got) || math.IsInf(got, 0) {
			t.Errorf("alpha %v: Result() = %v, want finite", a, got)
		}
	}
}

func TestRunningMinMax_Result(t *testing.T) {
	type args struct {
		values []int
	}
	tests := []struct {
		name   string
		args   args
		wantLo int
		wantHi int
		wantOK bool
	}{
		{name: "empty: not ok", args: args{values: nil}, wantLo: 0, wantHi: 0, wantOK: false},
		{name: "single element", args: args{values: []int{7}}, wantLo: 7, wantHi: 7, wantOK: true},
		{name: "ascending", args: args{values: []int{1, 2, 3, 4}}, wantLo: 1, wantHi: 4, wantOK: true},
		{name: "descending", args: args{values: []int{4, 3, 2, 1}}, wantLo: 1, wantHi: 4, wantOK: true},
		{name: "duplicates", args: args{values: []int{5, 5, 5}}, wantLo: 5, wantHi: 5, wantOK: true},
		{name: "negatives", args: args{values: []int{-3, 0, -10, 8}}, wantLo: -10, wantHi: 8, wantOK: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mm := streaming.NewRunningMinMax[int]()
			for _, v := range tt.args.values {
				mm.Add(v)
			}
			lo, hi, ok := mm.Result()
			if ok != tt.wantOK {
				t.Fatalf("Result() ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && (lo != tt.wantLo || hi != tt.wantHi) {
				t.Errorf("Result() = (%d, %d), want (%d, %d)", lo, hi, tt.wantLo, tt.wantHi)
			}
		})
	}
}

// TestRunningMinMax_MatchesStatsOracle cross-checks the streaming extremes
// against stats.MinMax over the same data.
func TestRunningMinMax_MatchesStatsOracle(t *testing.T) {
	data := []int{5, 1, 9, 3, 7, 2, 8, -4, 100}
	mm := streaming.NewRunningMinMax[int]()
	for i := range data {
		mm.Add(data[i])
		gotLo, gotHi, ok := mm.Result()
		wantLo, wantHi, wantOK := stats.MinMax(data[:i+1])
		if ok != wantOK || gotLo != wantLo || gotHi != wantHi {
			t.Errorf("after %d adds: (%d, %d, %v), stats = (%d, %d, %v)",
				i+1, gotLo, gotHi, ok, wantLo, wantHi, wantOK)
		}
	}
}

// TestRunningMinMax_Strings exercises the Ordered constraint on a non-numeric
// type to confirm the generic bound is satisfied.
func TestRunningMinMax_Strings(t *testing.T) {
	mm := streaming.NewRunningMinMax[string]()
	for _, s := range []string{"banana", "apple", "cherry"} {
		mm.Add(s)
	}
	lo, hi, ok := mm.Result()
	if !ok || lo != "apple" || hi != "cherry" {
		t.Errorf("Result() = (%q, %q, %v), want (apple, cherry, true)", lo, hi, ok)
	}
}
