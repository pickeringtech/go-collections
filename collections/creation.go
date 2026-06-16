package collections

import (
	"github.com/pickeringtech/go-collections/collections/deques"
	"github.com/pickeringtech/go-collections/collections/dicts"
	"github.com/pickeringtech/go-collections/collections/lists"
	"github.com/pickeringtech/go-collections/collections/sets"
)

// NewList creates a List backed by an array (slice) with the given values.
func NewList[T any](values ...T) lists.List[T] {
	return lists.NewArray(values...)
}

// NewConcurrentList creates a thread-safe List backed by an array with the given values.
func NewConcurrentList[T any](values ...T) lists.List[T] {
	return lists.NewConcurrentArray[T](values...)
}

//func NewConcurrentRWList[T any](values ...T) lists.List[T] {
//	return lists.NewConcurrentArrayRW[T](values...)
//}

// NewQueue creates a Queue (FIFO) backed by an array with the given values.
func NewQueue[T any](values ...T) lists.Queue[T] {
	return lists.NewArray(values...)
}

//func NewConcurrentRWQueue[T any](values ...T) lists.List[T] {
//	return lists.NewConcurrentArrayRW[T](values...)
//}

// NewConcurrentQueue creates a thread-safe Queue (FIFO) backed by an array with the given values.
func NewConcurrentQueue[T any](values ...T) lists.Queue[T] {
	return lists.NewConcurrentArray[T](values...)
}

// NewStack creates a Stack (LIFO) backed by an array with the given values.
func NewStack[T any](values ...T) lists.Stack[T] {
	return lists.NewArray(values...)
}

// NewConcurrentStack creates a thread-safe Stack (LIFO) backed by an array with the given values.
func NewConcurrentStack[T any](values ...T) lists.Stack[T] {
	return lists.NewConcurrentArray[T](values...)
}

//func NewConcurrentRWStack[T any](values ...T) lists.List[T] {
//	return lists.NewConcurrentArrayRW[T](values...)
//}

// NewDeque creates an unbounded Deque (double-ended queue) backed by a ring buffer with the given values (values[0] becomes the front).
func NewDeque[T any](values ...T) deques.Deque[T] {
	return deques.NewRingBuffer[T](values...)
}

// NewConcurrentDeque creates a thread-safe, unbounded Deque (mutex-guarded) backed by a ring buffer with the given values.
func NewConcurrentDeque[T any](values ...T) deques.Deque[T] {
	return deques.NewConcurrentRingBuffer[T](values...)
}

// NewConcurrentRWDeque creates a thread-safe, unbounded Deque optimised for concurrent reads (RWMutex-guarded) with the given values.
func NewConcurrentRWDeque[T any](values ...T) deques.Deque[T] {
	return deques.NewConcurrentRWRingBuffer[T](values...)
}

// NewBoundedDeque creates a bounded (circular) Deque with the given capacity and overflow policy, seeded with the given values.
func NewBoundedDeque[T any](capacity int, policy deques.OverflowPolicy, values ...T) deques.Deque[T] {
	return deques.NewBoundedRingBuffer[T](capacity, policy, values...)
}

// NewBoundedConcurrentDeque creates a thread-safe, bounded Deque (mutex-guarded) with the given capacity and overflow policy, seeded with the given values.
func NewBoundedConcurrentDeque[T any](capacity int, policy deques.OverflowPolicy, values ...T) deques.Deque[T] {
	return deques.NewBoundedConcurrentRingBuffer[T](capacity, policy, values...)
}

// NewBoundedConcurrentRWDeque creates a thread-safe, bounded Deque optimised for concurrent reads (RWMutex-guarded) with the given capacity and overflow policy, seeded with the given values.
func NewBoundedConcurrentRWDeque[T any](capacity int, policy deques.OverflowPolicy, values ...T) deques.Deque[T] {
	return deques.NewBoundedConcurrentRWRingBuffer[T](capacity, policy, values...)
}

// NewDict creates a Dict backed by a hash map with the given entries.
func NewDict[K comparable, V any](entries ...dicts.Pair[K, V]) dicts.Dict[K, V] {
	return dicts.NewHash[K, V](entries...)
}

// NewConcurrentDict creates a thread-safe Dict (mutex-guarded) with the given entries.
func NewConcurrentDict[K comparable, V any](entries ...dicts.Pair[K, V]) dicts.Dict[K, V] {
	return dicts.NewConcurrentHash[K, V](entries...)
}

// NewConcurrentRWDict creates a thread-safe Dict optimised for concurrent reads (RWMutex-guarded) with the given entries.
func NewConcurrentRWDict[K comparable, V any](entries ...dicts.Pair[K, V]) dicts.Dict[K, V] {
	return dicts.NewConcurrentHashRW[K, V](entries...)
}

// NewSet creates a Set backed by a hash map with the given elements.
func NewSet[T comparable](elements ...T) sets.Set[T] {
	return sets.NewHash[T](elements...)
}

// NewConcurrentSet creates a thread-safe Set (mutex-guarded) with the given elements.
func NewConcurrentSet[T comparable](elements ...T) sets.Set[T] {
	return sets.NewConcurrentHash[T](elements...)
}

// NewConcurrentRWSet creates a thread-safe Set optimised for concurrent reads (RWMutex-guarded) with the given elements.
func NewConcurrentRWSet[T comparable](elements ...T) sets.Set[T] {
	return sets.NewConcurrentHashRW[T](elements...)
}

// NewLinkedList creates a List backed by a singly linked list with the given elements.
func NewLinkedList[T any](elements ...T) lists.List[T] {
	return lists.NewLinked[T](elements...)
}

// NewConcurrentLinkedList creates a thread-safe List backed by a singly linked list with the given elements.
func NewConcurrentLinkedList[T any](elements ...T) lists.List[T] {
	return lists.NewConcurrentLinked[T](elements...)
}

// NewConcurrentRWLinkedList creates a thread-safe List backed by a singly linked list, optimised for concurrent reads, with the given elements.
func NewConcurrentRWLinkedList[T any](elements ...T) lists.List[T] {
	return lists.NewConcurrentRWLinked[T](elements...)
}

// NewDoublyLinkedList creates a List backed by a doubly linked list with the given elements.
func NewDoublyLinkedList[T any](elements ...T) lists.List[T] {
	return lists.NewDoublyLinked[T](elements...)
}

// NewConcurrentDoublyLinkedList creates a thread-safe List backed by a doubly linked list with the given elements.
func NewConcurrentDoublyLinkedList[T any](elements ...T) lists.List[T] {
	return lists.NewConcurrentDoublyLinked[T](elements...)
}

// NewConcurrentRWDoublyLinkedList creates a thread-safe List backed by a doubly linked list, optimised for concurrent reads, with the given elements.
func NewConcurrentRWDoublyLinkedList[T any](elements ...T) lists.List[T] {
	return lists.NewConcurrentRWDoublyLinked[T](elements...)
}
