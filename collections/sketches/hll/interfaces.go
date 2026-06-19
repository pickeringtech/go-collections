package hll

// Cardinality is the behavioural surface shared by Sketch and ConcurrentSketch,
// so callers can write code that accepts either variant. Merge takes the plain
// *Sketch[T] deliberately (see the note on bloom.Membership): the source must be
// a value the caller is not concurrently mutating.
type Cardinality[T comparable] interface {
	// Add records an element.
	Add(value T)

	// Count returns the estimated number of distinct elements added.
	Count() uint64

	// Merge folds other's registers into the receiver (register-wise max).
	Merge(other *Sketch[T]) error

	// Clear resets every register.
	Clear()

	// Precision returns the register-index width (p).
	Precision() int

	// RegisterCount returns the number of registers (m = 2^p).
	RegisterCount() int

	// StandardError returns the expected relative error of Count.
	StandardError() float64
}
