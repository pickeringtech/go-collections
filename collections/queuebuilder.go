package collections

import "github.com/pickeringtech/go-collections/collections/lists"

// QueueBuilder fluently configures and constructs a Queue, allowing the
// concurrency characteristics and initial values to be chosen before building.
type QueueBuilder[T any] struct {
	isConcurrent bool
	isRW         bool
	values       []T
}

// NewQueueBuilder creates an empty QueueBuilder that builds a plain (non-concurrent) Queue by default.
func NewQueueBuilder[T any]() QueueBuilder[T] {
	return QueueBuilder[T]{}
}

// Concurrent marks the builder to produce a thread-safe Queue, returning the updated builder.
func (b QueueBuilder[T]) Concurrent() QueueBuilder[T] {
	b.isConcurrent = true
	return b
}

// RW marks the builder to favour concurrent reads (read-write locking), returning the updated builder.
func (b QueueBuilder[T]) RW() QueueBuilder[T] {
	b.isRW = true
	return b
}

// Add appends the given values to those the built Queue will contain, returning the updated builder.
func (b QueueBuilder[T]) Add(values ...T) QueueBuilder[T] {
	b.values = append(b.values, values...)
	return b
}

// Build constructs the Queue described by the builder.
func (b QueueBuilder[T]) Build() lists.Queue[T] {
	if b.isConcurrent && b.isRW {
		return NewConcurrentRWQueue[T](b.values...)
	}
	if b.isConcurrent {
		return NewConcurrentQueue[T](b.values...)
	}
	return NewQueue[T](b.values...)
}
