package streaming

import (
	"math"
	"math/rand/v2"

	"github.com/pickeringtech/go-collections/collections/heaps"
)

// NewRand returns a *math/rand/v2.Rand seeded deterministically from seed, so
// that the same seed always reproduces the same sample — the basis for
// reproducible reservoir sampling.
//
// A deterministic, non-cryptographic generator is exactly the requirement here
// (reproducibility, not security), and the int64 seed is reinterpreted to
// uint64 bits where every bit pattern is a valid seed, so the conversion's
// wraparound on negative seeds is intentional. Both are flagged by gosec
// (G404 weak RNG, G115 integer conversion) and deliberately suppressed.
func NewRand(seed int64) *rand.Rand {
	return rand.New(rand.NewPCG(uint64(seed), uint64(seed))) // #nosec G404,G115
}

// randOrDefault returns rng, or a deterministically seeded generator when rng is
// nil, so the reservoir constructors are always well-defined and reproducible.
func randOrDefault(rng *rand.Rand) *rand.Rand {
	if rng != nil {
		return rng
	}
	return NewRand(0)
}

// Reservoir maintains a uniform random sample of up to k elements drawn from an
// unbounded stream, in O(k) memory regardless of how many elements are fed in.
// After any number of Add calls, every element seen so far is equally likely to
// be in the sample: each of the n elements is retained with probability
// min(1, k/n). This is Vitter's Algorithm R.
//
// Sampling is without replacement — no element appears in the sample twice
// (though equal-valued elements fed from the stream are each independent
// candidates). The sample is exact in its probabilities, not approximate.
//
// Reservoir is seeded for reproducibility: the same seed and the same stream
// always produce the same sample. It is single-threaded; see the package
// documentation on thread safety.
type Reservoir[T any] struct {
	k   int
	rng *rand.Rand
	// seen counts the elements fed so far. It is an int64, not an int, so it
	// cannot overflow after 2^31-1 elements on 32-bit platforms.
	seen  int64
	items []T
}

// NewReservoir creates a Reservoir that maintains a uniform random sample of up
// to k elements. rng drives the sampling; passing nil uses a deterministic
// default generator (equivalent to NewRand(0)), so sampling is reproducible
// unless you supply your own source. For k <= 0 the reservoir retains nothing
// and Result is always empty.
func NewReservoir[T any](k int, rng *rand.Rand) *Reservoir[T] {
	return &Reservoir[T]{
		k:   k,
		rng: randOrDefault(rng),
	}
}

// Add feeds one element into the stream. While fewer than k elements have been
// seen it is always retained; thereafter the n-th element (n > k) replaces a
// uniformly chosen incumbent with probability k/n, preserving the uniform-sample
// invariant. The element count is 64-bit on every platform, so the sample stays
// uniform for streams of up to math.MaxInt64 elements. Add is O(1) and a no-op
// when k <= 0.
func (r *Reservoir[T]) Add(element T) {
	if r.k <= 0 {
		return
	}
	r.seen++
	if len(r.items) < r.k {
		r.items = append(r.items, element)
		return
	}
	j := r.rng.Int64N(r.seen)
	if j < int64(r.k) {
		r.items[j] = element
	}
}

// Result returns a copy of the current sample as a non-nil slice. It does not
// modify the Reservoir, so Add may continue afterwards, and the caller may
// mutate the returned slice freely. The order of the sample is unspecified and
// must not be relied upon. The length is min(k, number of elements fed).
func (r *Reservoir[T]) Result() []T {
	out := make([]T, len(r.items))
	copy(out, r.items)
	return out
}

// Len returns the number of elements currently in the sample, between 0 and
// max(k, 0) — it is always 0 when k <= 0.
func (r *Reservoir[T]) Len() int {
	return len(r.items)
}

// weightedItem pairs a stream element with its random sampling key (see
// samplingKey); keeping the k largest keys yields a weighted sample without
// replacement (Efraimidis & Spirakis, A-Res).
type weightedItem[T any] struct {
	key  float64
	item T
}

// WeightedReservoir maintains a weighted random sample of up to k elements drawn
// from an unbounded stream, in O(k) memory. Each element is added with a
// positive weight, and an element's probability of being retained grows with its
// weight relative to the rest of the stream. This is the A-Res algorithm of
// Efraimidis & Spirakis: each element is assigned the random key u^(1/weight)
// and the k elements with the largest keys are kept, backed by a size-bounded
// min-heap (collections/heaps) so Add is O(log k). The key is evaluated in an
// order-preserving log-log form (see Add) so that every positive finite weight
// works at any scale: only the ratios between weights affect the sample.
//
// Sampling is without replacement. The sample is exact in its probabilities, not
// approximate. WeightedReservoir is seeded for reproducibility and is
// single-threaded; see the package documentation on thread safety.
type WeightedReservoir[T any] struct {
	k    int
	rng  *rand.Rand
	heap *heaps.Binary[weightedItem[T]]
}

// NewWeightedReservoir creates a WeightedReservoir that maintains a weighted
// random sample of up to k elements. rng drives the sampling; passing nil uses a
// deterministic default generator (equivalent to NewRand(0)). For k <= 0 the
// reservoir retains nothing and Result is always empty.
func NewWeightedReservoir[T any](k int, rng *rand.Rand) *WeightedReservoir[T] {
	return &WeightedReservoir[T]{
		k:   k,
		rng: randOrDefault(rng),
		heap: heaps.New(func(a, b weightedItem[T]) bool {
			return a.key < b.key
		}),
	}
}

// Add feeds one element into the stream with the given weight. Weights must be
// strictly positive and finite; an element whose weight is zero, negative, NaN
// or ±Inf has no well-defined sampling probability and is ignored (an infinite
// weight would tie with every other infinite weight). Larger weights make an
// element proportionally more likely to be retained, and scaling every weight by
// the same factor leaves the sampling distribution unchanged — any positive
// finite magnitude, subnormal to math.MaxFloat64, is handled. Add is O(log k)
// and a no-op when k <= 0.
func (r *WeightedReservoir[T]) Add(element T, weight float64) {
	if r.k <= 0 || weight <= 0 || math.IsNaN(weight) || math.IsInf(weight, 1) {
		return
	}
	entry := weightedItem[T]{key: samplingKey(r.rng.Float64(), weight), item: element}
	if r.heap.Length() < r.k {
		r.heap.PushInPlace(entry)
		return
	}
	lowest, _ := r.heap.Peek()
	if lowest.key < entry.key {
		r.heap.PopInPlace()
		r.heap.PushInPlace(entry)
	}
}

// samplingKey returns the A-Res key for an element of the given weight from a
// uniform draw u in [0, 1), in a form that orders exactly as u^(1/weight) but
// never underflows or overflows. Evaluated directly, u^(1/weight) underflows to
// 0 for small weights and rounds to 1 for huge ones, and the resulting ties let
// the earliest elements win; even log(u)/weight overflows to -Inf for subnormal
// weights. Since log is increasing, u^(1/w) orders like log(u)/w = -E/w with
// E = -log(u) > 0, which orders like log(w) - log(E): both terms are bounded
// (|log w| <= ~745 for any positive finite w, and log E is in about
// [-36.7, 3.6] for the u that Float64 returns). u == 0 gives E = +Inf and the
// key -Inf, the lowest possible, matching 0^(1/w) = 0.
func samplingKey(u, weight float64) float64 {
	return math.Log(weight) - math.Log(-math.Log(u))
}

// Result returns a copy of the current sample as a non-nil slice, ordered by
// descending sampling key (most strongly retained first). It does not modify the
// WeightedReservoir, so Add may continue afterwards. The length is min(k, number
// of positively-weighted elements fed).
func (r *WeightedReservoir[T]) Result() []T {
	sorted := r.heap.AsSortedSlice()
	out := make([]T, len(sorted))
	for i, entry := range sorted {
		out[len(sorted)-1-i] = entry.item
	}
	return out
}

// Len returns the number of elements currently in the sample, between 0 and
// max(k, 0) — it is always 0 when k <= 0.
func (r *WeightedReservoir[T]) Len() int {
	return r.heap.Length()
}
