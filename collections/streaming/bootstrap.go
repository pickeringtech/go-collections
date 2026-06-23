package streaming

import "math/rand/v2"

// Bootstrap draws a bootstrap resample of input: a new slice of the same length
// formed by sampling len(input) elements uniformly at random WITH REPLACEMENT,
// so a given element may appear several times or not at all. This is the
// fundamental building block of the statistical bootstrap, used to estimate the
// sampling distribution of a statistic by recomputing it over many resamples.
//
// rng drives the sampling; passing nil uses a deterministic default generator
// (equivalent to NewRand(0)), so a resample is reproducible unless you supply
// your own source — matching the seeding contract of the rest of this package.
//
// The input is never mutated. For nil or empty input the result is a non-nil
// empty slice. The returned slice is freshly allocated, so the caller may mutate
// it freely. Use BootstrapN to draw several independent resamples at once.
func Bootstrap[T any](input []T, rng *rand.Rand) []T {
	r := randOrDefault(rng)
	out := make([]T, len(input))
	if len(input) == 0 {
		return out
	}
	for i := range out {
		out[i] = input[r.IntN(len(input))]
	}
	return out
}

// BootstrapN draws count independent bootstrap resamples of input, each of
// length len(input) and sampled uniformly with replacement (see Bootstrap).
// Successive resamples share one generator, so the whole batch is reproducible
// from a single seed.
//
// rng drives the sampling; passing nil uses a deterministic default generator
// (equivalent to NewRand(0)). The input is never mutated. For count <= 0 the
// result is a non-nil empty slice; otherwise it holds exactly count resamples,
// each itself a non-nil slice (empty when input is nil or empty).
func BootstrapN[T any](input []T, count int, rng *rand.Rand) [][]T {
	if count <= 0 {
		return [][]T{}
	}
	r := randOrDefault(rng)
	out := make([][]T, 0, count)
	for i := 0; i < count; i++ {
		out = append(out, Bootstrap(input, r))
	}
	return out
}
