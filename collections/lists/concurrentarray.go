package lists

import (
	"sort"
	"sync"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
	"github.com/pickeringtech/go-collections/slices"
)

// ConcurrentArray is a slice-backed implementation of MutableList that is safe
// for concurrent use. Every operation is guarded by a sync.Mutex.
//
// ConcurrentArray must not be copied after first use; copying after construction
// produces an independent lock over shared backing data, which breaks the
// thread-safety contract. go vet reports any such copy.
type ConcurrentArray[T any] struct {
	_        nocopy.NoCopy
	elements []T
	// version counts writes to elements. Every mutating method increments it
	// under the write lock, so FilterInPlace can tell whether anything wrote
	// to the list while its predicate ran outside the lock.
	version uint64
	lock    sync.Mutex
}

// NewConcurrentArray creates a new ConcurrentArray seeded with the given
// elements, preserving their order. The elements are copied, so the list owns
// its backing array and is unaffected by later mutations of the caller's slice
// (which would otherwise bypass the lock).
func NewConcurrentArray[T any](elements ...T) *ConcurrentArray[T] {
	return &ConcurrentArray[T]{
		elements: slices.Copy(elements),
	}
}

// Interface guards
var _ List[int] = &ConcurrentArray[int]{}
var _ MutableList[int] = &ConcurrentArray[int]{}

// AllMatch returns true if every element satisfies the predicate fun (vacuously
// true for an empty list). The predicate is evaluated after the lock is
// released, against a point-in-time snapshot taken under the lock, so it may
// safely call back into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) AllMatch(fun func(T) bool) bool {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	return slices.AllMatch(snapshot, fun)
}

// AnyMatch returns true if at least one element satisfies the predicate fun. The
// predicate is evaluated after the lock is released, against a point-in-time
// snapshot taken under the lock, so it may safely call back into the
// collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) AnyMatch(fun func(T) bool) bool {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	return slices.AnyMatch(snapshot, fun)
}

// NoneMatch returns true if no element satisfies the predicate fun (vacuously
// true for an empty list). The predicate is evaluated after the lock is
// released, against a point-in-time snapshot taken under the lock, so it may
// safely call back into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) NoneMatch(fun func(T) bool) bool {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	return !slices.AnyMatch(snapshot, fun)
}

// Dequeue returns the first element, whether one was present, and a new List
// (independent of the receiver's backing array) with that element removed,
// without modifying the receiver. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Dequeue() (T, bool, List[T]) {
	a.lock.Lock()
	defer a.lock.Unlock()

	// Operate on a copy so the returned List is independent of the receiver's
	// backing array (PopFront returns a sub-slice of its input).
	res, ok, rest := slices.PopFront(slices.Copy(a.elements))
	return res, ok, NewArray(rest...)
}

// DequeueInPlace removes and returns the first element, reporting whether one
// was present, modifying the receiver. It is safe for concurrent use.
// It advances over the array's own storage in O(1) without allocating.
func (a *ConcurrentArray[T]) DequeueInPlace() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.version++

	res, ok, rest := ownedPopFront(a.elements)
	a.elements = rest
	return res, ok
}

// Enqueue returns a new List (independent of the receiver's backing array)
// with element appended to the end, without modifying the receiver. It is safe
// for concurrent use.
func (a *ConcurrentArray[T]) Enqueue(element T) List[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	// PushCopy copies and appends in a single allocation, so the new List is
	// independent of the receiver (plain Push could append into shared capacity).
	return NewArray(slices.PushCopy(a.elements, element)...)
}

// EnqueueInPlace appends element to the end of the receiver. It is safe for
// concurrent use.
// It appends to the array's own storage, so it is amortised O(1).
func (a *ConcurrentArray[T]) EnqueueInPlace(element T) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.version++

	a.elements = append(a.elements, element)
}

// Filter returns a new List containing only the elements for which fun returns
// true, without modifying the receiver. The predicate is evaluated after the
// lock is released, against a point-in-time snapshot taken under the lock, so it
// may safely call back into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Filter(fun func(T) bool) List[T] {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	return NewArray(slices.Filter(snapshot, fun)...)
}

// FilterInPlace retains only the elements for which fn returns true, modifying
// the receiver. The predicate is called once per element, in order, after the
// lock is released, against a point-in-time snapshot taken under the lock, so
// it may safely call back into the collection. It is safe for concurrent use.
//
// If nothing wrote to the list while the predicate ran, the list is left
// holding exactly the elements the predicate kept, in their original order —
// the same result as Array.FilterInPlace. If a write did land in that window
// (from another goroutine, or from the predicate itself), removal falls back to
// a multiset diff against the current contents so the concurrent write is
// preserved: each rejected element removes the first deeply-equal
// (reflect.DeepEqual) occurrence from the list as it stands at apply time. In
// that fallback a rejected element that is not DeepEqual to itself (such as
// NaN) is kept, and when the list holds several deeply-equal elements the first
// is removed rather than the one rejected.
func (a *ConcurrentArray[T]) FilterInPlace(fn func(T) bool) {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	version := a.version
	a.lock.Unlock()

	kept, rejected := partitionElements(snapshot, fn)

	a.lock.Lock()
	if a.version == version {
		a.elements = kept
	} else {
		a.elements = removeEachDeepEqual(a.elements, rejected)
	}
	a.version++
	a.lock.Unlock()
}

// Find returns the first element for which fun returns true and whether such an
// element was found. The predicate is evaluated after the lock is released,
// against a point-in-time snapshot taken under the lock, so it may safely call
// back into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Find(fun func(T) bool) (T, bool) {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	return slices.Find(snapshot, fun)
}

// FindIndex returns the index of the first element for which fun returns true,
// or -1 if none match. The predicate is evaluated after the lock is released,
// against a point-in-time snapshot taken under the lock, so it may safely call
// back into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) FindIndex(fun func(T) bool) int {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	return slices.FindIndex(snapshot, fun)
}

// ForEach calls fun once for each element in order. fun is invoked after the
// lock is released, against a point-in-time snapshot taken under the lock, so
// fun may safely call back into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) ForEach(fun EachFunc[T]) {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	for _, element := range snapshot {
		fun(element)
	}
}

// ForEachWithIndex calls fun once for each element in order, passing the
// element's index and value. fun is invoked after the lock is released, against
// a point-in-time snapshot taken under the lock, so fun may safely call back
// into the collection. It is safe for concurrent use.
func (a *ConcurrentArray[T]) ForEachWithIndex(fun IndexedEachFunc[T]) {
	a.lock.Lock()
	snapshot := slices.Copy(a.elements)
	a.lock.Unlock()

	for idx, element := range snapshot {
		fun(idx, element)
	}
}

// Get returns the element at index and true, or defaultValue and false if the
// index is out of bounds. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Get(index int, defaultValue T) (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if index < 0 || index >= len(a.elements) {
		return defaultValue, false
	}
	return a.elements[index], true
}

// AsSlice returns a copy of the elements as a new slice, independent of the
// receiver's backing array. It is safe for concurrent use.
func (a *ConcurrentArray[T]) AsSlice() []T {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Copy(a.elements)
}

// Insert returns a new List (independent of the receiver) with the given
// elements inserted at index, without modifying the receiver. It is safe for
// concurrent use.
func (a *ConcurrentArray[T]) Insert(index int, element ...T) List[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	// Operate on a copy so the returned List is independent of the receiver and
	// the receiver's backing array is never mutated by the insert.
	return NewArray(slices.Insert(slices.Copy(a.elements), index, element...)...)
}

// InsertInPlace inserts the given elements at index, modifying the receiver. It
// is safe for concurrent use.
func (a *ConcurrentArray[T]) InsertInPlace(index int, element ...T) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.version++

	a.elements = slices.Insert(a.elements, index, element...)
}

// Length returns the number of elements in the list. It is safe for concurrent
// use.
func (a *ConcurrentArray[T]) Length() int {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Length(a.elements)
}

// IsEmpty returns true if the list contains no elements. It is safe for
// concurrent use.
func (a *ConcurrentArray[T]) IsEmpty() bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.Length(a.elements) == 0
}

// RemoveAt returns a new List (independent of the receiver's backing array)
// with the element at index removed, without modifying the receiver. If index
// is out of bounds the elements are returned unchanged. It is safe for
// concurrent use.
func (a *ConcurrentArray[T]) RemoveAt(index int) List[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	return NewArray(deleteOwned(slices.Copy(a.elements), index)...)
}

// Remove returns a new List (independent of the receiver's backing array) with
// the first element deeply equal to element removed, without modifying the
// receiver. If no element matches, the elements are returned unchanged. It is
// safe for concurrent use.
func (a *ConcurrentArray[T]) Remove(element T) List[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	elements := slices.Copy(a.elements)
	return NewArray(deleteOwned(elements, indexOfDeepEqual(elements, element))...)
}

// RemoveAtInPlace removes the element at index, returning it and whether the
// index was in bounds, modifying the receiver. It is safe for concurrent use.
// It shifts the array's own storage in place rather than allocating.
func (a *ConcurrentArray[T]) RemoveAtInPlace(index int) (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.version++

	if index < 0 || index >= len(a.elements) {
		var zero T
		return zero, false
	}
	removed := a.elements[index]
	a.elements = ownedRemoveAt(a.elements, index)
	return removed, true
}

// RemoveInPlace removes the first element deeply equal to element, reporting
// whether an element was removed, modifying the receiver. It is safe for
// concurrent use.
// It shifts the array's own storage in place rather than allocating.
func (a *ConcurrentArray[T]) RemoveInPlace(element T) bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.version++

	index := indexOfDeepEqual(a.elements, element)
	if index < 0 {
		return false
	}
	a.elements = ownedRemoveAt(a.elements, index)
	return true
}

// Clear removes all elements from the list. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Clear() {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.version++

	a.elements = nil
}

// PeekEnd returns the last element without removing it, and whether one was
// present. It is safe for concurrent use.
func (a *ConcurrentArray[T]) PeekEnd() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.PeekEnd(a.elements)
}

// PeekFront returns the first element without removing it, and whether one was
// present. It is safe for concurrent use.
func (a *ConcurrentArray[T]) PeekFront() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return slices.PeekFront(a.elements)
}

// Pop returns the last element, whether one was present, and a new List
// (independent of the receiver's backing array) with that element removed,
// without modifying the receiver. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Pop() (T, bool, List[T]) {
	a.lock.Lock()
	defer a.lock.Unlock()

	// Operate on a copy so the returned List is independent of the receiver's
	// backing array (Pop returns a sub-slice of its input).
	res, ok, rest := slices.Pop(slices.Copy(a.elements))
	return res, ok, NewArray(rest...)
}

// PopInPlace removes and returns the last element, reporting whether one was
// present, modifying the receiver. It is safe for concurrent use.
// It shrinks the array's own storage in O(1) without allocating, keeping the
// capacity for later pushes.
func (a *ConcurrentArray[T]) PopInPlace() (T, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.version++

	res, ok, rest := ownedPopBack(a.elements)
	a.elements = rest
	return res, ok
}

// Push returns a new List (independent of the receiver's backing array) with
// element appended to the end, without modifying the receiver. It is safe for
// concurrent use.
func (a *ConcurrentArray[T]) Push(element T) List[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	// PushCopy copies and appends in a single allocation, so the new List is
	// independent of the receiver (plain Push could append into shared capacity).
	return NewArray(slices.PushCopy(a.elements, element)...)
}

// PushInPlace appends element to the end of the receiver. It is safe for
// concurrent use.
// It appends to the array's own storage, so it is amortised O(1).
func (a *ConcurrentArray[T]) PushInPlace(element T) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.version++

	a.elements = append(a.elements, element)
}

// Sort returns a new List sorted according to the less-than function lessThan,
// without modifying the receiver. The sort is stable: elements that compare
// equal keep their original relative order. It is safe for concurrent use.
func (a *ConcurrentArray[T]) Sort(lessThan func(T, T) bool) List[T] {
	a.lock.Lock()
	defer a.lock.Unlock()

	sorted := slices.Copy(a.elements)
	sort.SliceStable(sorted, func(i, j int) bool { return lessThan(sorted[i], sorted[j]) })
	return NewArray(sorted...)
}

// SortInPlace sorts the receiver's elements according to the less-than function
// lessThan. The sort is stable: elements that compare equal keep their
// original relative order. It is safe for concurrent use.
func (a *ConcurrentArray[T]) SortInPlace(lessThan func(T, T) bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.version++

	sort.SliceStable(a.elements, func(i, j int) bool { return lessThan(a.elements[i], a.elements[j]) })
}

// partitionElements calls fn once per snapshot element, in order, and splits
// the elements into those it kept (in their original order) and those it
// rejected. It is the evaluation phase shared by every concurrent list's
// FilterInPlace, run outside the lock.
func partitionElements[T any](snapshot []T, fn func(T) bool) (kept, rejected []T) {
	for _, element := range snapshot {
		if fn(element) {
			kept = append(kept, element)
		} else {
			rejected = append(rejected, element)
		}
	}
	return kept, rejected
}

// removeEachDeepEqual removes from elements, for each rejected element, the
// first occurrence deeply equal to it, and returns the result. It is the
// conservative apply phase the concurrent arrays' FilterInPlace uses when a
// write raced the predicate.
func removeEachDeepEqual[T any](elements, rejected []T) []T {
	for _, element := range rejected {
		index := indexOfDeepEqual(elements, element)
		if index >= 0 {
			elements = ownedRemoveAt(elements, index)
		}
	}
	return elements
}

// applyFilter is the apply phase shared by the concurrent linked lists'
// FilterInPlace; the caller holds the write lock over list. kept and rejected
// are the snapshot elements the predicate kept and rejected, and unchanged
// reports whether nothing wrote to the list since the snapshot was taken.
//
// When unchanged, the list still equals the snapshot, so its contents are
// replaced with kept — exactly what the non-concurrent FilterInPlace leaves.
// Clear and PushInPlace keep the list's own shape (circular or not). Otherwise
// each rejected element removes the first deeply-equal occurrence from the
// list as it now stands, preserving the concurrent write.
func applyFilter[T any](list MutableList[T], kept, rejected []T, unchanged bool) {
	if len(rejected) == 0 {
		return // nothing was rejected, so there is nothing to remove either way
	}
	if unchanged {
		list.Clear()
		for _, element := range kept {
			list.PushInPlace(element)
		}
		return
	}
	for _, element := range rejected {
		list.RemoveInPlace(element)
	}
}

// ownedPopBack removes and returns the last element of a slice the caller owns
// (the concurrent arrays' own storage), reporting whether there was one. It
// works in place in O(1): the vacated slot is zeroed so the backing array does
// not keep the element reachable, and the capacity is kept for later pushes.
func ownedPopBack[T any](elements []T) (T, bool, []T) {
	var zero T
	last := len(elements) - 1
	if last < 0 {
		return zero, false, elements
	}
	element := elements[last]
	elements[last] = zero
	return element, true, elements[:last]
}

// ownedPopFront removes and returns the first element of a slice the caller
// owns, reporting whether there was one. It works in place in O(1), zeroing the
// vacated slot; removing the only element rewinds to the start of the backing
// array, so a queue that keeps draining to empty reuses its capacity instead of
// walking off the end of it.
func ownedPopFront[T any](elements []T) (T, bool, []T) {
	var zero T
	if len(elements) == 0 {
		return zero, false, elements
	}
	element := elements[0]
	elements[0] = zero
	if len(elements) == 1 {
		return element, true, elements[:0]
	}
	return element, true, elements[1:]
}

// ownedRemoveAt removes the element at index, which must be in bounds, from a
// slice the caller owns. It shifts the tail left in place and zeroes the
// vacated final slot rather than allocating a new slice.
func ownedRemoveAt[T any](elements []T, index int) []T {
	copy(elements[index:], elements[index+1:])
	last := len(elements) - 1
	var zero T
	elements[last] = zero
	return elements[:last]
}
