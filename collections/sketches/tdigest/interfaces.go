package tdigest

// Quantiles is the behavioural surface shared by Digest and ConcurrentDigest,
// so callers can write code that accepts either variant. Merge takes the plain
// *Digest deliberately (mirroring the sibling sketches): the source must be a
// value the caller is not concurrently mutating — pass a ConcurrentDigest's
// Snapshot for the safe equivalent.
type Quantiles interface {
	// Add records a single value with weight 1, ignoring non-finite values.
	Add(x float64)

	// AddWeighted records a value carrying the given weight, ignoring
	// non-finite values or non-positive weights.
	AddWeighted(x float64, w float64)

	// Quantile returns the estimated q-quantile for q in [0, 1].
	Quantile(q float64) (float64, bool)

	// Percentile returns the estimated p-th percentile for p in [0, 100].
	Percentile(p float64) (float64, bool)

	// CDF returns the estimated fraction of the distribution at or below x.
	CDF(x float64) (float64, bool)

	// Merge folds other's centroids into the receiver (union of streams).
	Merge(other *Digest) error

	// Clear empties the digest, keeping its compression.
	Clear()

	// Count returns the total weight added.
	Count() float64

	// Min returns the smallest value ever added.
	Min() (float64, bool)

	// Max returns the largest value ever added.
	Max() (float64, bool)

	// Compression returns the compression parameter.
	Compression() float64
}
