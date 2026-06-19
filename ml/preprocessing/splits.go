package preprocessing

import "math/rand/v2"

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
	if len(input) == 0 || testFrac < 0 || testFrac > 1 {
		return nil, nil, false
	}
	indices := shuffledIndices(len(input), randOrDefault(rng))
	testN := int(testFrac*float64(len(input)) + 0.5)

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
// but preserves class proportions: within each label group, testFrac of the
// elements go to test. labels must be parallel to input. Both returned slices
// are fresh copies; input is never modified. rng drives the partition; passing
// nil uses a deterministic default.
//
// ok is false when input is empty, len(input) != len(labels), or testFrac is
// outside [0, 1].
func StratifiedSplit[T any, L comparable](input []T, labels []L, testFrac float64, rng *rand.Rand) (train, test []T, ok bool) {
	if len(input) == 0 || len(input) != len(labels) || testFrac < 0 || testFrac > 1 {
		return nil, nil, false
	}
	generator := randOrDefault(rng)

	// Group element indices by label, preserving first-seen label order so the
	// partition is deterministic for a given seed.
	order := make([]L, 0)
	groups := make(map[L][]int)
	for i, label := range labels {
		existing, seen := groups[label]
		if !seen {
			order = append(order, label)
		}
		groups[label] = append(existing, i)
	}

	train = make([]T, 0, len(input))
	test = make([]T, 0, len(input))
	for _, label := range order {
		members := groups[label]
		generator.Shuffle(len(members), func(i, j int) { members[i], members[j] = members[j], members[i] })
		testN := int(testFrac*float64(len(members)) + 0.5)
		for pos, idx := range members {
			if pos < testN {
				test = append(test, input[idx])
				continue
			}
			train = append(train, input[idx])
		}
	}
	return train, test, true
}

// StratifiedSplitSeed is StratifiedSplit using a generator seeded from seed.
func StratifiedSplitSeed[T any, L comparable](input []T, labels []L, testFrac float64, seed int64) (train, test []T, ok bool) {
	return StratifiedSplit(input, labels, testFrac, NewRand(seed))
}
