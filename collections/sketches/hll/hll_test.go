package hll_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/collections/sketches/hll"
)

func mustNew[T comparable](t *testing.T, opts ...hll.Option[T]) *hll.Sketch[T] {
	t.Helper()
	s, err := hll.New(opts...)
	if err != nil {
		t.Fatalf("New: unexpected error %v", err)
	}
	return s
}

// relErr returns the relative error of got against the true value.
func relErr(got uint64, truth int) float64 {
	return math.Abs(float64(got)-float64(truth)) / float64(truth)
}

func TestNew_InvalidPrecision(t *testing.T) {
	tests := []struct {
		name string
		p    int
	}{
		{"too low", hll.MinPrecision - 1},
		{"too high", hll.MaxPrecision + 1},
		{"zero", 0},
		// Guards the int-typed precision: a large value must not wrap into the
		// valid range (e.g. 260 mod 256 = 4) and be silently accepted.
		{"far above range", 260},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := hll.New[int](hll.WithPrecision[int](tc.p))
			if !errors.Is(err, hll.ErrInvalidConfig) {
				t.Errorf("error = %v, want ErrInvalidConfig", err)
			}
		})
	}
}

func TestNew_PrecisionBoundsAreInclusive(t *testing.T) {
	tests := []struct {
		name          string
		p             int
		wantRegisters int
	}{
		{name: "minimum precision is accepted", p: hll.MinPrecision, wantRegisters: 16},
		{name: "maximum precision is accepted", p: hll.MaxPrecision, wantRegisters: 262144},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := hll.New[int](hll.WithPrecision[int](tt.p))
			if err != nil {
				t.Fatalf("New(WithPrecision(%d)) error = %v, want nil", tt.p, err)
			}
			if s.RegisterCount() != tt.wantRegisters {
				t.Errorf("RegisterCount() = %d, want %d", s.RegisterCount(), tt.wantRegisters)
			}
		})
	}
}

func TestSketch_Defaults(t *testing.T) {
	s := mustNew[int](t)
	if got := s.Precision(); got != hll.DefaultPrecision {
		t.Errorf("Precision() = %d, want %d", got, hll.DefaultPrecision)
	}
	if got := s.RegisterCount(); got != 1<<hll.DefaultPrecision {
		t.Errorf("RegisterCount() = %d, want %d", got, 1<<hll.DefaultPrecision)
	}
}

func TestSketch_EmptyCountZero(t *testing.T) {
	s := mustNew[int](t)
	if got := s.Count(); got != 0 {
		t.Errorf("empty Count() = %d, want 0", got)
	}
}

func TestSketch_AccuracyAcrossScales(t *testing.T) {
	for _, n := range []int{100, 1000, 10000, 100000} {
		t.Run(fmt.Sprintf("n_%d", n), func(t *testing.T) {
			s := mustNew[int](t)
			for i := 0; i < n; i++ {
				s.Add(i)
			}
			got := s.Count()
			// Allow ~3 standard errors plus the small-range slack.
			tolerance := 3*s.StandardError() + 0.01
			if e := relErr(got, n); e > tolerance {
				t.Errorf("Count() = %d for n=%d, relative error %.4f exceeds tolerance %.4f", got, n, e, tolerance)
			}
		})
	}
}

func TestSketch_Idempotent(t *testing.T) {
	s := mustNew[int](t)
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	before := s.Count()
	// Re-adding the same elements must not change the distinct-count estimate.
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	if after := s.Count(); after != before {
		t.Errorf("Count changed after re-adding duplicates: %d -> %d", before, after)
	}
}

func TestSketch_Merge(t *testing.T) {
	a := mustNew[int](t)
	b := mustNew[int](t)
	// Disjoint halves; the union has 20000 distinct elements.
	for i := 0; i < 10000; i++ {
		a.Add(i)
	}
	for i := 10000; i < 20000; i++ {
		b.Add(i)
	}
	if err := a.Merge(b); err != nil {
		t.Fatalf("Merge: %v", err)
	}
	got := a.Count()
	tolerance := 3*a.StandardError() + 0.01
	if e := relErr(got, 20000); e > tolerance {
		t.Errorf("merged Count() = %d, relative error %.4f exceeds tolerance %.4f", got, e, tolerance)
	}
}

func TestSketch_MergeOverlap(t *testing.T) {
	a := mustNew[int](t)
	b := mustNew[int](t)
	// Fully overlapping; the union is still 10000 distinct.
	for i := 0; i < 10000; i++ {
		a.Add(i)
		b.Add(i)
	}
	if err := a.Merge(b); err != nil {
		t.Fatalf("Merge: %v", err)
	}
	got := a.Count()
	tolerance := 3*a.StandardError() + 0.01
	if e := relErr(got, 10000); e > tolerance {
		t.Errorf("merged Count() = %d, relative error %.4f exceeds tolerance %.4f", got, e, tolerance)
	}
}

func TestSketch_MergeMismatch(t *testing.T) {
	a := mustNew[int](t)
	tests := map[string]*hll.Sketch[int]{
		"different precision": mustNew[int](t, hll.WithPrecision[int](12)),
		"different seed":      mustNew[int](t, hll.WithSeed[int](99)),
	}
	for name, b := range tests {
		t.Run(name, func(t *testing.T) {
			if err := a.Merge(b); !errors.Is(err, hll.ErrInvalidConfig) {
				t.Errorf("Merge error = %v, want ErrInvalidConfig", err)
			}
		})
	}
}

func TestSketch_Clear(t *testing.T) {
	s := mustNew[int](t)
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	s.Clear()
	if got := s.Count(); got != 0 {
		t.Errorf("after Clear, Count() = %d, want 0", got)
	}
}

func TestSketch_LowPrecisionStillWorks(t *testing.T) {
	s := mustNew[int](t, hll.WithPrecision[int](hll.MinPrecision))
	for i := 0; i < 50; i++ {
		s.Add(i)
	}
	// Coarse precision, but the estimate should be in the right ballpark.
	if got := s.Count(); got < 20 || got > 100 {
		t.Errorf("Count() = %d for 50 distinct at min precision, want roughly 50", got)
	}
}

func TestSketch_CustomHasher(t *testing.T) {
	calls := 0
	h := func(seed uint64, v int) uint64 {
		calls++
		return uint64(v) * 2654435761
	}
	s := mustNew[int](t, hll.WithHasher[int](h))
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}
	if calls == 0 {
		t.Error("custom hasher was never called")
	}
	if s.Count() == 0 {
		t.Error("Count() = 0 after 1000 adds with custom hasher")
	}
}

func TestSketch_MergeNil(t *testing.T) {
	s := mustNew[int](t)
	if err := s.Merge(nil); !errors.Is(err, hll.ErrInvalidConfig) {
		t.Errorf("Merge(nil) error = %v, want ErrInvalidConfig", err)
	}
}

// TestSketch_AlphaConstants exercises the small-m bias-correction constants by
// estimating at the precisions whose register counts (16, 32, 64) have special
// alpha values.
func TestSketch_AlphaConstants(t *testing.T) {
	for _, p := range []int{4, 5, 6} { // m = 16, 32, 64
		t.Run(fmt.Sprintf("p_%d", p), func(t *testing.T) {
			s := mustNew[int](t, hll.WithPrecision[int](p))
			for i := 0; i < 200; i++ {
				s.Add(i)
			}
			// Coarse registers, but the estimate must be a sane positive number.
			if got := s.Count(); got == 0 {
				t.Errorf("Count() = 0 at precision %d", p)
			}
		})
	}
}

func ExampleSketch() {
	s, _ := hll.New[string]()
	for i := 0; i < 100000; i++ {
		s.Add(fmt.Sprintf("user-%d", i))
	}
	// Estimate is within ~1% of the true 100000 distinct users.
	within := math.Abs(float64(s.Count())-100000)/100000 < 0.02
	fmt.Println("within 2%:", within)
	// Output:
	// within 2%: true
}

// TestSketch_SignedZeroIsOneElement checks that +0 and -0 (equal under ==) are
// counted as one distinct element, and every NaN payload as another single one.
func TestSketch_SignedZeroIsOneElement(t *testing.T) {
	s := mustNew[float64](t)
	s.Add(0)
	s.Add(math.Copysign(0, -1))
	s.Add(math.NaN())
	s.Add(math.Float64frombits(0x7FF0_0000_0000_0001))
	s.Add(math.Copysign(math.NaN(), -1))
	got := s.Count()
	if got != 2 {
		t.Errorf("Count() = %d for {+0, -0, NaN×3}, want 2", got)
	}
}

// TestSketch_SeedsGiveIndependentEstimates runs the same stream through sketches
// with seeds 0..63 and checks the estimates behave like independent trials:
// their mean relative error is near zero and their spread is near the
// documented StandardError. When seeds were not mixed properly, seeds 0..255
// gave bit-identical registers for 256-aligned integer keys, so every "trial"
// returned the same estimate (zero spread, and a mean error equal to one
// trial's error).
func TestSketch_SeedsGiveIndependentEstimates(t *testing.T) {
	const (
		precision = 10 // m = 1024, StandardError ≈ 3.25%
		n         = 4096
		trials    = 64
	)
	var sum, sumSq float64
	var se float64
	for seed := uint64(0); seed < trials; seed++ {
		s := mustNew[int](t, hll.WithPrecision[int](precision), hll.WithSeed[int](seed))
		for i := 0; i < n; i++ {
			s.Add(i)
		}
		rel := (float64(s.Count()) - n) / n
		sum += rel
		sumSq += rel * rel
		se = s.StandardError()
	}
	mean := sum / trials
	spread := math.Sqrt(sumSq/trials - mean*mean)
	// The mean of 64 independent trials has standard error se/8; allow 4 of them.
	if math.Abs(mean) > 4*se/8 {
		t.Errorf("mean relative error across seeds = %+.4f, want within ±%.4f", mean, 4*se/8)
	}
	if spread < se/2 || spread > 2*se {
		t.Errorf("spread of estimates across seeds = %.4f, want within [%.4f, %.4f] of StandardError %.4f",
			spread, se/2, 2*se, se)
	}
}

// TestSketch_InterfaceKeysCountedByType checks that a Sketch[any] counts values
// of different dynamic types as distinct, as == does: int(i), int64(i),
// uint32(i) and float64(i) are four keys, not one.
func TestSketch_InterfaceKeysCountedByType(t *testing.T) {
	s := mustNew[any](t)
	for i := 0; i < 1000; i++ {
		s.Add(i) // dynamic type int
		s.Add(int64(i))
		s.Add(uint32(i))
		s.Add(float64(i))
	}
	got := s.Count()
	tolerance := 3*s.StandardError() + 0.01
	e := relErr(got, 4000)
	if e > tolerance {
		t.Errorf("Count() = %d for 4000 distinct any keys, relative error %.4f exceeds %.4f", got, e, tolerance)
	}
}

// identityHasher feeds a value straight through as its own hash, so a test can
// choose exactly which register an element lands in and with what rank.
func identityHasher(_ uint64, v uint64) uint64 { return v }

// hashWithRank returns a hash that a precision-p sketch files under register
// idx with rank r: the top p bits are the index and the remaining 64-p bits
// hold r-1 leading zeros then a one. r = 65-p is the all-zero remainder, whose
// rank is capped by the guard bit at 64-p leading zeros plus one.
func hashWithRank(p int, idx uint64, r int) uint64 {
	h := idx << (64 - p)
	if r < 65-p {
		h |= 1 << (64 - p - r)
	}
	return h
}

// sketchWithRanks builds a precision-p sketch whose register i holds ranks[i]
// (0 leaves the register empty).
func sketchWithRanks(t *testing.T, p int, ranks []int) *hll.Sketch[uint64] {
	t.Helper()
	if len(ranks) != 1<<p {
		t.Fatalf("got %d ranks for precision %d, want %d", len(ranks), p, 1<<p)
	}
	s := mustNew[uint64](t, hll.WithPrecision[uint64](p), hll.WithHasher[uint64](identityHasher))
	for idx, r := range ranks {
		if r > 0 {
			s.Add(hashWithRank(p, uint64(idx), r))
		}
	}
	return s
}

// repeatRank returns n copies of rank r.
func repeatRank(n, r int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = r
	}
	return out
}

// rawEstimate is the textbook HyperLogLog estimate alpha·m²/Σ2^-r over the
// given register ranks.
func rawEstimate(alpha float64, ranks []int) float64 {
	m := float64(len(ranks))
	sum := 0.0
	for _, r := range ranks {
		sum += math.Ldexp(1, -r)
	}
	return alpha * m * m / sum
}

// TestSketch_CountUsesRawEstimateWhenNoRegisterIsEmpty pins Count against the
// textbook estimator on hand-built registers: with no empty register linear
// counting does not apply, so Count is exactly alpha_m·m²/Σ2^-r with the bias
// constant alpha_m of Flajolet et al. (0.673, 0.697 and 0.709 for m = 16, 32,
// 64; 0.7213/(1+1.079/m) beyond). It also pins the rank of an all-zero hash
// remainder at 65-p, the cap the guard bit sets.
func TestSketch_CountUsesRawEstimateWhenNoRegisterIsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		p     int
		ranks []int
		alpha float64
	}{
		{name: "sixteen registers use alpha 0.673", p: 4, ranks: repeatRank(16, 10), alpha: 0.673},
		{name: "thirty-two registers use alpha 0.697", p: 5, ranks: repeatRank(32, 10), alpha: 0.697},
		{name: "sixty-four registers use alpha 0.709", p: 6, ranks: repeatRank(64, 10), alpha: 0.709},
		{name: "128 registers use the asymptotic alpha", p: 7, ranks: repeatRank(128, 10), alpha: 0.7213 / (1 + 1.079/128)},
		{name: "1024 registers use the asymptotic alpha", p: 10, ranks: repeatRank(1024, 10), alpha: 0.7213 / (1 + 1.079/1024)},
		// Every register at rank 1 estimates 21.536, below the 2.5m = 40
		// linear-counting threshold, but with no empty register linear
		// counting (m·ln(m/0)) is undefined and must not be used.
		{name: "a small estimate with no empty register is not linear-counted", p: 4, ranks: repeatRank(16, 1), alpha: 0.673},
		{
			name: "an all-zero hash remainder takes the capped rank 65-p",
			p:    4,
			// Register 15 gets the all-zero remainder (rank 61); the others sit
			// at rank 55 so that its 2^-61 term moves the estimate visibly.
			ranks: append(repeatRank(15, 55), 61),
			alpha: 0.673,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sketchWithRanks(t, tt.p, tt.ranks)
			got := s.Count()
			want := uint64(rawEstimate(tt.alpha, tt.ranks) + 0.5)
			if got != want {
				t.Errorf("Count() = %d, want %d", got, want)
			}
		})
	}
}

// TestSketch_CountLinearCountsAtTheThreshold pins the small-range switch at
// its boundary: linear counting applies while the raw estimate is at most
// 2.5m (inclusive) and some register is empty. The ranks below make the raw
// estimate exactly 160 = 2.5·64 in float64 with one empty register, so Count
// must be the linear-counting m·ln(m/zeros) = 64·ln 64 ≈ 266.17, not 160.
func TestSketch_CountLinearCountsAtTheThreshold(t *testing.T) {
	ranks := []int{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		3, 6, 7, 9, 17, 20, 21, 22, 24, 26, 29, 32, 34, 36, 40, 41, 47,
		0, // the one empty register
	}
	raw := rawEstimate(0.709, ranks)
	if raw != 160 {
		t.Fatalf("fixture raw estimate = %v, want exactly 2.5m = 160", raw)
	}
	s := sketchWithRanks(t, 6, ranks)
	got := s.Count()
	if got != 266 {
		t.Errorf("Count() = %d at a raw estimate of exactly 2.5m, want linear counting's 266", got)
	}
}
