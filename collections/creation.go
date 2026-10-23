package collections

import (
	"iter"

	"github.com/pickeringtech/go-collections/collections/deques"
	"github.com/pickeringtech/go-collections/collections/dicts"
	"github.com/pickeringtech/go-collections/collections/heaps"
	"github.com/pickeringtech/go-collections/collections/lists"
	"github.com/pickeringtech/go-collections/collections/lru"
	"github.com/pickeringtech/go-collections/collections/multimaps"
	"github.com/pickeringtech/go-collections/collections/sets"
	"github.com/pickeringtech/go-collections/constraints"
)

// NewList creates a List backed by an array (slice) with the given values. It
// returns the MutableList interface so the in-place operations (PushInPlace,
// InsertInPlace, ...) are reachable; MutableList embeds List, so the immutable
// API remains available too.
func NewList[T any](values ...T) lists.MutableList[T] {
	return lists.NewArray(values...)
}

// NewConcurrentList creates a thread-safe List backed by an array with the given
// values. It returns the MutableList interface so the in-place operations are
// reachable.
func NewConcurrentList[T any](values ...T) lists.MutableList[T] {
	return lists.NewConcurrentArray[T](values...)
}

// NewConcurrentRWList creates a thread-safe List backed by an array, optimised
// for concurrent reads (RWMutex-guarded), with the given values. It returns the
// MutableList interface so the in-place operations are reachable.
func NewConcurrentRWList[T any](values ...T) lists.MutableList[T] {
	return lists.NewConcurrentRWArray[T](values...)
}

// NewQueue creates a Queue (FIFO) backed by an array with the given values. It
// returns the MutableQueue interface so the in-place operations (EnqueueInPlace,
// DequeueInPlace) are reachable; MutableQueue embeds Queue, so the immutable API
// remains available too.
func NewQueue[T any](values ...T) lists.MutableQueue[T] {
	return lists.NewArray(values...)
}

// NewConcurrentRWQueue creates a thread-safe Queue (FIFO) backed by an array,
// optimised for concurrent reads (RWMutex-guarded), with the given values. It
// returns the MutableQueue interface so the in-place operations are reachable.
func NewConcurrentRWQueue[T any](values ...T) lists.MutableQueue[T] {
	return lists.NewConcurrentRWArray[T](values...)
}

// NewConcurrentQueue creates a thread-safe Queue (FIFO) backed by an array with
// the given values. It returns the MutableQueue interface so the in-place
// operations are reachable.
func NewConcurrentQueue[T any](values ...T) lists.MutableQueue[T] {
	return lists.NewConcurrentArray[T](values...)
}

// NewStack creates a Stack (LIFO) backed by an array with the given values. It
// returns the MutableStack interface so the in-place operations (PushInPlace,
// PopInPlace) are reachable; MutableStack embeds Stack, so the immutable API
// remains available too.
func NewStack[T any](values ...T) lists.MutableStack[T] {
	return lists.NewArray(values...)
}

// NewConcurrentStack creates a thread-safe Stack (LIFO) backed by an array with
// the given values. It returns the MutableStack interface so the in-place
// operations are reachable.
func NewConcurrentStack[T any](values ...T) lists.MutableStack[T] {
	return lists.NewConcurrentArray[T](values...)
}

// NewConcurrentRWStack creates a thread-safe Stack (LIFO) backed by an array,
// optimised for concurrent reads (RWMutex-guarded), with the given values. It
// returns the MutableStack interface so the in-place operations are reachable.
func NewConcurrentRWStack[T any](values ...T) lists.MutableStack[T] {
	return lists.NewConcurrentRWArray[T](values...)
}

// NewDeque creates an unbounded Deque (double-ended queue) backed by a ring
// buffer with the given values (values[0] becomes the front). It returns the
// MutableDeque interface so the in-place operations (InsertInPlace,
// RemoveAtInPlace, ...) are reachable; MutableDeque embeds Deque, so the
// immutable API remains available too.
func NewDeque[T any](values ...T) deques.MutableDeque[T] {
	return deques.NewRingBuffer[T](values...)
}

// NewConcurrentDeque creates a thread-safe, unbounded Deque (mutex-guarded) backed by a ring buffer with the given values. It returns the MutableDeque interface so the in-place operations are reachable.
func NewConcurrentDeque[T any](values ...T) deques.MutableDeque[T] {
	return deques.NewConcurrentRingBuffer[T](values...)
}

// NewConcurrentRWDeque creates a thread-safe, unbounded Deque optimised for concurrent reads (RWMutex-guarded) with the given values. It returns the MutableDeque interface so the in-place operations are reachable.
func NewConcurrentRWDeque[T any](values ...T) deques.MutableDeque[T] {
	return deques.NewConcurrentRWRingBuffer[T](values...)
}

// NewBoundedDeque creates a bounded (circular) Deque with the given capacity and overflow policy, seeded with the given values. It returns the MutableDeque interface so the in-place operations are reachable.
func NewBoundedDeque[T any](capacity int, policy deques.OverflowPolicy, values ...T) deques.MutableDeque[T] {
	return deques.NewBoundedRingBuffer[T](capacity, policy, values...)
}

// NewBoundedConcurrentDeque creates a thread-safe, bounded Deque (mutex-guarded) with the given capacity and overflow policy, seeded with the given values. It returns the MutableDeque interface so the in-place operations are reachable.
func NewBoundedConcurrentDeque[T any](capacity int, policy deques.OverflowPolicy, values ...T) deques.MutableDeque[T] {
	return deques.NewBoundedConcurrentRingBuffer[T](capacity, policy, values...)
}

// NewBoundedConcurrentRWDeque creates a thread-safe, bounded Deque optimised for concurrent reads (RWMutex-guarded) with the given capacity and overflow policy, seeded with the given values. It returns the MutableDeque interface so the in-place operations are reachable.
func NewBoundedConcurrentRWDeque[T any](capacity int, policy deques.OverflowPolicy, values ...T) deques.MutableDeque[T] {
	return deques.NewBoundedConcurrentRWRingBuffer[T](capacity, policy, values...)
}

// NewDict creates a Dict backed by a hash map with the given entries. It returns
// the MutableDict interface so the in-place operations (PutInPlace,
// UpdateInPlace, ...) are reachable; MutableDict embeds Dict, so the immutable
// API remains available too.
func NewDict[K comparable, V any](entries ...dicts.Pair[K, V]) dicts.MutableDict[K, V] {
	return dicts.NewHash[K, V](entries...)
}

// NewConcurrentDict creates a thread-safe Dict (mutex-guarded) with the given
// entries. It returns the MutableDict interface so the in-place operations are
// reachable; use UpdateInPlace for race-free read-modify-write (e.g. counters).
func NewConcurrentDict[K comparable, V any](entries ...dicts.Pair[K, V]) dicts.MutableDict[K, V] {
	return dicts.NewConcurrentHash[K, V](entries...)
}

// NewConcurrentRWDict creates a thread-safe Dict optimised for concurrent reads
// (RWMutex-guarded) with the given entries. It returns the MutableDict interface
// so the in-place operations are reachable; use UpdateInPlace for race-free
// read-modify-write (e.g. counters).
func NewConcurrentRWDict[K comparable, V any](entries ...dicts.Pair[K, V]) dicts.MutableDict[K, V] {
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

// NewLinkedList creates a List backed by a singly linked list with the given elements. It returns the MutableList interface so the in-place operations are reachable.
func NewLinkedList[T any](elements ...T) lists.MutableList[T] {
	return lists.NewLinked[T](elements...)
}

// NewConcurrentLinkedList creates a thread-safe List backed by a singly linked list with the given elements. It returns the MutableList interface so the in-place operations are reachable.
func NewConcurrentLinkedList[T any](elements ...T) lists.MutableList[T] {
	return lists.NewConcurrentLinked[T](elements...)
}

// NewConcurrentRWLinkedList creates a thread-safe List backed by a singly linked list, optimised for concurrent reads, with the given elements. It returns the MutableList interface so the in-place operations are reachable.
func NewConcurrentRWLinkedList[T any](elements ...T) lists.MutableList[T] {
	return lists.NewConcurrentRWLinked[T](elements...)
}

// NewListMultimap creates a list-backed Multimap (one key to many ordered,
// possibly-duplicate values) with the given entries. V may be any type. It
// returns the MutableMultimap interface so the in-place operations (PutInPlace,
// RemoveInPlace, ...) are reachable; MutableMultimap embeds Multimap, so the
// immutable API remains available too.
func NewListMultimap[K comparable, V any](entries ...multimaps.Entry[K, V]) multimaps.MutableMultimap[K, V] {
	return multimaps.NewListMultimap(entries...)
}

// NewConcurrentListMultimap creates a thread-safe list-backed Multimap (mutex-guarded) with the given entries. V may be any type. It returns the MutableMultimap interface so the in-place operations are reachable.
func NewConcurrentListMultimap[K comparable, V any](entries ...multimaps.Entry[K, V]) multimaps.MutableMultimap[K, V] {
	return multimaps.NewConcurrentListMultimap(entries...)
}

// NewConcurrentRWListMultimap creates a thread-safe list-backed Multimap optimised for concurrent reads (RWMutex-guarded) with the given entries. V may be any type. It returns the MutableMultimap interface so the in-place operations are reachable.
func NewConcurrentRWListMultimap[K comparable, V any](entries ...multimaps.Entry[K, V]) multimaps.MutableMultimap[K, V] {
	return multimaps.NewConcurrentRWListMultimap(entries...)
}

// NewSetMultimap creates a set-backed Multimap (one key to many distinct values) with the given entries. It returns the MutableMultimap interface so the in-place operations are reachable.
func NewSetMultimap[K comparable, V comparable](entries ...multimaps.Entry[K, V]) multimaps.MutableMultimap[K, V] {
	return multimaps.NewSetMultimap(entries...)
}

// NewConcurrentSetMultimap creates a thread-safe set-backed Multimap (mutex-guarded) with the given entries. It returns the MutableMultimap interface so the in-place operations are reachable.
func NewConcurrentSetMultimap[K comparable, V comparable](entries ...multimaps.Entry[K, V]) multimaps.MutableMultimap[K, V] {
	return multimaps.NewConcurrentSetMultimap(entries...)
}

// NewConcurrentRWSetMultimap creates a thread-safe set-backed Multimap optimised for concurrent reads (RWMutex-guarded) with the given entries. It returns the MutableMultimap interface so the in-place operations are reachable.
func NewConcurrentRWSetMultimap[K comparable, V comparable](entries ...multimaps.Entry[K, V]) multimaps.MutableMultimap[K, V] {
	return multimaps.NewConcurrentRWSetMultimap(entries...)
}

// NewDoublyLinkedList creates a List backed by a doubly linked list with the given elements. It returns the MutableList interface so the in-place operations are reachable.
func NewDoublyLinkedList[T any](elements ...T) lists.MutableList[T] {
	return lists.NewDoublyLinked[T](elements...)
}

// NewConcurrentDoublyLinkedList creates a thread-safe List backed by a doubly linked list with the given elements. It returns the MutableList interface so the in-place operations are reachable.
func NewConcurrentDoublyLinkedList[T any](elements ...T) lists.MutableList[T] {
	return lists.NewConcurrentDoublyLinked[T](elements...)
}

// NewConcurrentRWDoublyLinkedList creates a thread-safe List backed by a doubly linked list, optimised for concurrent reads, with the given elements. It returns the MutableList interface so the in-place operations are reachable.
func NewConcurrentRWDoublyLinkedList[T any](elements ...T) lists.MutableList[T] {
	return lists.NewConcurrentRWDoublyLinked[T](elements...)
}

// NewHeap creates a Heap (priority queue) ordered by the given comparator, seeded with the given values. The values are heapified in O(n).
func NewHeap[T any](less heaps.LessFunc[T], values ...T) heaps.Heap[T] {
	return heaps.New(less, values...)
}

// NewMinHeap creates a min-heap over an ordered type (the smallest element leaves the heap first), seeded with the given values.
func NewMinHeap[T constraints.Ordered](values ...T) heaps.Heap[T] {
	return heaps.NewMin(values...)
}

// NewMaxHeap creates a max-heap over an ordered type (the largest element leaves the heap first), seeded with the given values.
func NewMaxHeap[T constraints.Ordered](values ...T) heaps.Heap[T] {
	return heaps.NewMax(values...)
}

// NewConcurrentHeap creates a thread-safe Heap (mutex-guarded) ordered by the given comparator, seeded with the given values.
func NewConcurrentHeap[T any](less heaps.LessFunc[T], values ...T) heaps.Heap[T] {
	return heaps.NewConcurrent(less, values...)
}

// NewConcurrentMinHeap creates a thread-safe min-heap (mutex-guarded) over an ordered type, seeded with the given values.
func NewConcurrentMinHeap[T constraints.Ordered](values ...T) heaps.Heap[T] {
	return heaps.NewConcurrentMin(values...)
}

// NewConcurrentMaxHeap creates a thread-safe max-heap (mutex-guarded) over an ordered type, seeded with the given values.
func NewConcurrentMaxHeap[T constraints.Ordered](values ...T) heaps.Heap[T] {
	return heaps.NewConcurrentMax(values...)
}

// NewConcurrentRWHeap creates a thread-safe Heap optimised for concurrent reads (RWMutex-guarded) ordered by the given comparator, seeded with the given values.
func NewConcurrentRWHeap[T any](less heaps.LessFunc[T], values ...T) heaps.Heap[T] {
	return heaps.NewConcurrentRW(less, values...)
}

// NewConcurrentRWMinHeap creates a thread-safe min-heap optimised for concurrent reads (RWMutex-guarded) over an ordered type, seeded with the given values.
func NewConcurrentRWMinHeap[T constraints.Ordered](values ...T) heaps.Heap[T] {
	return heaps.NewConcurrentRWMin(values...)
}

// NewConcurrentRWMaxHeap creates a thread-safe max-heap optimised for concurrent reads (RWMutex-guarded) over an ordered type, seeded with the given values.
func NewConcurrentRWMaxHeap[T constraints.Ordered](values ...T) heaps.Heap[T] {
	return heaps.NewConcurrentRWMax(values...)
}

// NewLRU creates a bounded least-recently-used cache holding at most capacity entries; inserting beyond that evicts the least-recently-used entry. A capacity below 1 is treated as 1. Configure optional behaviour (an eviction callback, seed entries) with lru.Option values. It returns a MutableCache because an LRU is inherently stateful — its defining recency-marking read, Get, is a mutation.
func NewLRU[K comparable, V any](capacity int, opts ...lru.Option[K, V]) lru.MutableCache[K, V] {
	return lru.NewLRU[K, V](capacity, opts...)
}

// NewConcurrentLRU creates a thread-safe LRU cache (mutex-guarded) bounded to capacity entries. It accepts the same lru.Option values as NewLRU.
func NewConcurrentLRU[K comparable, V any](capacity int, opts ...lru.Option[K, V]) lru.MutableCache[K, V] {
	return lru.NewConcurrentLRU[K, V](capacity, opts...)
}

// NewConcurrentRWLRU creates a thread-safe LRU cache optimised for concurrent reads (RWMutex-guarded) bounded to capacity entries. It accepts the same lru.Option values as NewLRU.
func NewConcurrentRWLRU[K comparable, V any](capacity int, opts ...lru.Option[K, V]) lru.MutableCache[K, V] {
	return lru.NewConcurrentLRURW[K, V](capacity, opts...)
}

// ListFromSeq creates a List backed by an array from the values produced by seq,
// preserving their order. It is the inbound counterpart to the List.Values
// iterator. It returns the MutableList interface so the in-place operations are
// reachable.
func ListFromSeq[T any](seq iter.Seq[T]) lists.MutableList[T] {
	return lists.FromSeq(seq)
}

// DictFromSeq2 creates a Dict backed by a hash map from the key/value pairs produced
// by seq. When seq yields the same key more than once, the last value wins.
func DictFromSeq2[K comparable, V any](seq iter.Seq2[K, V]) dicts.Dict[K, V] {
	return dicts.FromSeq2(seq)
}

// SetFromSeq creates a Set backed by a hash map from the elements produced by seq,
// collapsing duplicates as any set does.
func SetFromSeq[T comparable](seq iter.Seq[T]) sets.Set[T] {
	return sets.FromSeq(seq)
}

// DequeFromSeq creates an unbounded Deque backed by a ring buffer from the values
// produced by seq (the first value becomes the front). It returns the
// MutableDeque interface so the in-place operations are reachable.
func DequeFromSeq[T any](seq iter.Seq[T]) deques.MutableDeque[T] {
	return deques.FromSeq(seq)
}

// HeapFromSeq creates a Heap (priority queue) ordered by less from the values
// produced by seq.
func HeapFromSeq[T any](less heaps.LessFunc[T], seq iter.Seq[T]) heaps.Heap[T] {
	return heaps.FromSeq(less, seq)
}

// ListMultimapFromSeq2 creates a list-backed Multimap from the key/value pairs
// produced by seq, preserving order and keeping duplicate values. V may be any
// type. It returns the MutableMultimap interface so the in-place operations are
// reachable.
func ListMultimapFromSeq2[K comparable, V any](seq iter.Seq2[K, V]) multimaps.MutableMultimap[K, V] {
	return multimaps.ListMultimapFromSeq2(seq)
}

// SetMultimapFromSeq2 creates a set-backed Multimap from the key/value pairs produced
// by seq, collapsing duplicate values bound to the same key. It returns the
// MutableMultimap interface so the in-place operations are reachable.
func SetMultimapFromSeq2[K comparable, V comparable](seq iter.Seq2[K, V]) multimaps.MutableMultimap[K, V] {
	return multimaps.SetMultimapFromSeq2(seq)
}
