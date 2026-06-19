package countmin

// Frequency is the behavioural surface shared by Sketch and ConcurrentSketch,
// so callers can write code that accepts either variant. Merge takes the plain
// *Sketch[T] deliberately (see the note on bloom.Membership): the source must be
// a value the caller is not concurrently mutating.
type Frequency[T comparable] interface {
	// Add records one occurrence of value.
	Add(value T)

	// AddCount records n occurrences of value.
	AddCount(value T, n uint64)

	// Estimate returns the approximate count of value (never an under-count).
	Estimate(value T) uint64

	// Merge folds other's counters into the receiver.
	Merge(other *Sketch[T]) error

	// Clear resets every counter.
	Clear()

	// Total returns the sum of all counts added (N).
	Total() uint64

	// Width returns the number of columns per row (w).
	Width() uint64

	// Depth returns the number of rows (d).
	Depth() uint64
}
