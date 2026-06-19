package bloom

// Membership is the behavioural surface shared by Filter and ConcurrentFilter,
// so callers can write code that accepts either variant. Merge takes the plain
// *Filter[T] (not the interface) deliberately: a concurrent filter locks itself
// and reads the other under no lock, so the source must be a value the caller
// is not concurrently mutating.
type Membership[T comparable] interface {
	// Add inserts an element.
	Add(value T)

	// Contains reports whether value may be present (no false negatives).
	Contains(value T) bool

	// Merge folds other's contents into the receiver.
	Merge(other *Filter[T]) error

	// Clear empties the filter, keeping its capacity and seed.
	Clear()

	// ApproxCount estimates the number of distinct elements added.
	ApproxCount() uint64

	// EstimatedFalsePositiveRate returns the false-positive probability at the
	// current fill level.
	EstimatedFalsePositiveRate() float64

	// BitSize returns the underlying bit-array size (m).
	BitSize() uint64

	// HashCount returns the number of hash functions per element (k).
	HashCount() uint64
}
