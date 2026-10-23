package preprocessing

import (
	"math/rand/v2"
	"sort"
)

// NewRand returns a *math/rand/v2.Rand seeded deterministically from seed, so
// that the same seed always reproduces the same sequence — the basis for
// reproducible splits and shuffles.
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
// nil, so the split functions are always well-defined and reproducible.
func randOrDefault(rng *rand.Rand) *rand.Rand {
	if rng != nil {
		return rng
	}
	return NewRand(0)
}

// shuffledIndices returns the indices [0, n) in a freshly shuffled order using
// the Fisher-Yates shuffle driven by rng.
func shuffledIndices(n int, rng *rand.Rand) []int {
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}
	rng.Shuffle(n, func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })
	return indices
}

// Shuffle returns a randomly reordered copy of input, leaving input untouched.
// rng drives the ordering; passing nil uses a deterministic default. A nil or
// empty input yields a non-nil empty slice.
func Shuffle[T any](input []T, rng *rand.Rand) []T {
	out := make([]T, len(input))
	copy(out, input)
	randOrDefault(rng).Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
	return out
}

// ShuffleSeed is Shuffle using a generator seeded from seed.
func ShuffleSeed[T any](input []T, seed int64) []T {
	return Shuffle(input, NewRand(seed))
}

// TrainTestSplit randomly partitions input into a train and a test set, where
// testFrac (in [0, 1]) is the fraction of elements assigned to test. The number
// of test elements is round(testFrac * len(input)). Both returned slices are
// fresh copies; input is never modified. rng drives the partition; passing nil
// uses a deterministic default.
//
// ok is false when input is empty or testFrac is outside [0, 1].
func TrainTestSplit[T any](input []T, testFrac float64, rng *rand.Rand) (train, test []T, ok bool) {
	if len(input) == 0 || !validFraction(testFrac) {
		return nil, nil, false
	}
	indices := shuffledIndices(len(input), randOrDefault(rng))
	testN := testCount(len(input), testFrac)

	test = make([]T, 0, testN)
	train = make([]T, 0, len(input)-testN)
	for pos, idx := range indices {
		if pos < testN {
			test = append(test, input[idx])
			continue
		}
		train = append(train, input[idx])
	}
	return train, test, true
}

// validFraction reports whether frac lies in [0, 1]. It is written so that NaN,
// for which every comparison is false, is rejected rather than slipping past a
// pair of range checks and yielding an undefined split.
func validFraction(frac float64) bool {
	return frac >= 0 && frac <= 1
}

// testCount is the number of the n elements a split assigns to test:
// round(testFrac·n), rounding halves up. TrainTestSplit and StratifiedSplit
// share it so the two always agree on the test size.
func testCount(n int, testFrac float64) int {
	return int(testFrac*float64(n) + 0.5)
}

// TrainTestSplitSeed is TrainTestSplit using a generator seeded from seed.
func TrainTestSplitSeed[T any](input []T, testFrac float64, seed int64) (train, test []T, ok bool) {
	return TrainTestSplit(input, testFrac, NewRand(seed))
}

// KFold randomly partitions input into k folds of near-equal size, returning the
// folds as fresh copies; input is never modified. The first len(input) mod k
// folds are one element larger. Each fold serves as a test set, with the union
// of the others as the corresponding train set. rng drives the partition;
// passing nil uses a deterministic default.
//
// ok is false when input is empty or k is outside [1, len(input)].
func KFold[T any](input []T, k int, rng *rand.Rand) (folds [][]T, ok bool) {
	if len(input) == 0 || k < 1 || k > len(input) {
		return nil, false
	}
	indices := shuffledIndices(len(input), randOrDefault(rng))

	folds = make([][]T, k)
	base := len(input) / k
	remainder := len(input) % k
	pos := 0
	for f := 0; f < k; f++ {
		size := base
		if f < remainder {
			size++
		}
		fold := make([]T, 0, size)
		for i := 0; i < size; i++ {
			fold = append(fold, input[indices[pos]])
			pos++
		}
		folds[f] = fold
	}
	return folds, true
}

// KFoldSeed is KFold using a generator seeded from seed.
func KFoldSeed[T any](input []T, k int, seed int64) (folds [][]T, ok bool) {
	return KFold(input, k, NewRand(seed))
}

// StratifiedSplit partitions input into train and test sets like TrainTestSplit,
// but preserves class proportions. labels must be parallel to input. Both
// returned slices are fresh copies; input is never modified. rng drives the
// partition; passing nil uses a deterministic default, and the same generator
// state always reproduces the same split.
//
// The test set holds exactly round(testFrac·len(input)) elements — the same
// count TrainTestSplit draws for the same length and testFrac. Those rows are
// apportioned across the labels by the largest-remainder method: a label with
// k of the n rows first receives ⌊T·k/n⌋ of the T test rows, and the rows left
// over go one each to the labels with the largest remainders (T·k mod n). Ties
// between equal remainders are broken by a random order of the labels drawn
// from rng, so when many labels tie (several singletons, say) which of them
// reach test varies with the seed rather than always favouring the first-seen.
//
// The proportional shares are then adjusted, one row at a time, so that every
// label with two or more rows appears on both sides whenever the totals permit:
//   - if T is at least the number of labels with two or more rows, every such
//     label has at least one test row;
//   - if len(input) − T is at least that number, every such label keeps at
//     least one train row.
//
// Each adjustment moves a single row from (or to) the label whose test count
// most exceeds (or most falls short of) its exact proportional share, ties
// again in the random label order. A label with a single row necessarily
// lands on exactly one side. Within each label the rows are shuffled before
// they are divided, and the output lists labels in first-seen order.
//
// ok is false when input is empty, len(input) != len(labels), testFrac is
// outside [0, 1], or any label is NaN. A NaN label (detected as a label unequal
// to itself, so this covers NaN floats and NaN-carrying complex or struct
// labels) cannot be grouped — it matches no other label, not even itself — so
// rather than silently dropping its row the split is rejected.
func StratifiedSplit[T any, L comparable](input []T, labels []L, testFrac float64, rng *rand.Rand) (train, test []T, ok bool) {
	if len(input) == 0 || len(input) != len(labels) || !validFraction(testFrac) {
		return nil, nil, false
	}
	groups, ok := groupByLabel(labels)
	if !ok {
		return nil, nil, false
	}
	generator := randOrDefault(rng)

	testN := testCount(len(input), testFrac)
	counts := apportion(groups, testN, generator.Perm(len(groups)))

	train = make([]T, 0, len(input)-testN)
	test = make([]T, 0, testN)
	for g, members := range groups {
		generator.Shuffle(len(members), func(i, j int) { members[i], members[j] = members[j], members[i] })
		for pos, idx := range members {
			if pos < counts[g] {
				test = append(test, input[idx])
				continue
			}
			train = append(train, input[idx])
		}
	}
	return train, test, true
}

// groupByLabel returns the indices of labels grouped by label value, groups in
// first-seen order so the partition is deterministic for a given seed. ok is
// false if any label is NaN — unequal to itself, it would form a fresh,
// unmatchable map key every time it appears.
func groupByLabel[L comparable](labels []L) (groups [][]int, ok bool) {
	index := make(map[L]int)
	for i, label := range labels {
		if label != label {
			return nil, false
		}
		g, seen := index[label]
		if !seen {
			g = len(groups)
			index[label] = g
			groups = append(groups, nil)
		}
		groups[g] = append(groups[g], i)
	}
	return groups, true
}

// apportionment is the working state of StratifiedSplit's allocation of testN
// test rows across label groups of the given sizes (summing to n). priority is
// a random permutation of the group indices that breaks every tie.
type apportionment struct {
	sizes    []int
	counts   []int // test rows allocated to each group
	testN, n int
	priority []int
}

// apportion returns how many of each group's rows go to test: the
// largest-remainder share of testN, then adjusted so multi-row groups reach
// both sides where the totals permit (see StratifiedSplit).
func apportion(groups [][]int, testN int, priority []int) []int {
	a := apportionment{sizes: make([]int, len(groups)), counts: make([]int, len(groups)), testN: testN, priority: priority}
	for g, members := range groups {
		a.sizes[g] = len(members)
		a.n += len(members)
	}
	a.largestRemainder()
	a.coverBothSides()
	return a.counts
}

// largestRemainder gives each group ⌊testN·size/n⌋ test rows, then hands the
// leftover rows to the groups with the largest remainders testN·size mod n,
// ties in priority order. The arithmetic is exact in int64 (n² fits even on
// 32-bit platforms), so equal remainders really are equal.
func (a *apportionment) largestRemainder() {
	remainders := make([]int64, len(a.sizes))
	assigned := 0
	for g, size := range a.sizes {
		share := int64(a.testN) * int64(size)
		a.counts[g] = int(share / int64(a.n))
		remainders[g] = share % int64(a.n)
		assigned += a.counts[g]
	}
	order := append([]int(nil), a.priority...)
	sort.SliceStable(order, func(i, j int) bool { return remainders[order[i]] > remainders[order[j]] })
	// The leftover is fewer than the groups with a non-zero remainder, so only
	// those receive a row and none exceeds its size.
	for _, g := range order[:a.testN-assigned] {
		a.counts[g]++
	}
}

// coverBothSides moves single rows between groups so that every group with two
// or more rows has a test row (when testN allows one each) and a train row
// (when n − testN does). When a total allows it a willing donor or recipient
// always exists — otherwise the rows on that side would number fewer than the
// multi-row groups — and neither pass undoes the other: the test pass only
// lifts groups from zero test rows and takes from groups with rows to spare,
// and the train pass mirrors it.
func (a *apportionment) coverBothSides() {
	multi := 0
	for _, size := range a.sizes {
		if size >= 2 {
			multi++
		}
	}
	if a.testN >= multi {
		for _, g := range a.priority {
			if a.sizes[g] >= 2 && a.counts[g] == 0 {
				a.move(a.pick(a.canGiveTestRow, 1), g)
			}
		}
	}
	if a.n-a.testN >= multi {
		for _, g := range a.priority {
			if a.sizes[g] >= 2 && a.counts[g] == a.sizes[g] {
				a.move(g, a.pick(a.canTakeTestRow, -1))
			}
		}
	}
}

// move transfers one test row from group from to group to.
func (a *apportionment) move(from, to int) {
	a.counts[from]--
	a.counts[to]++
}

// excess is how far group g's test count exceeds its exact proportional share
// testN·size/n, scaled by n to stay in exact integers.
func (a *apportionment) excess(g int) int64 {
	return int64(a.counts[g])*int64(a.n) - int64(a.testN)*int64(a.sizes[g])
}

// pick returns the eligible group with the largest excess (sign 1) or the
// smallest (sign −1), ties in priority order. coverBothSides only asks when an
// eligible group is guaranteed to exist.
func (a *apportionment) pick(eligible func(g int) bool, sign int64) int {
	best := -1
	for _, g := range a.priority {
		if eligible(g) && (best < 0 || sign*a.excess(g) > sign*a.excess(best)) {
			best = g
		}
	}
	return best
}

// canGiveTestRow reports whether group g can hand a test row to another group
// without losing its own last one (a single-row group has nothing to keep).
func (a *apportionment) canGiveTestRow(g int) bool {
	return canSpare(a.counts[g], a.sizes[g])
}

// canTakeTestRow reports whether group g can take one more test row without
// losing its last train row (a single-row group has nothing to keep).
func (a *apportionment) canTakeTestRow(g int) bool {
	return canSpare(a.sizes[g]-a.counts[g], a.sizes[g])
}

// canSpare reports whether a group of size rows holding have rows on one side
// can give one up and still be represented there — or is a single row that
// cannot be on both sides anyway.
func canSpare(have, size int) bool {
	return have >= 2 || (size == 1 && have == 1)
}

// StratifiedSplitSeed is StratifiedSplit using a generator seeded from seed.
func StratifiedSplitSeed[T any, L comparable](input []T, labels []L, testFrac float64, seed int64) (train, test []T, ok bool) {
	return StratifiedSplit(input, labels, testFrac, NewRand(seed))
}
