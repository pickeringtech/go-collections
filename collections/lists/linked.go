package lists

import "reflect"

type node[T any] struct {
	value  T
	next   *node[T]
	linked *Linked[T]
}

// Linked is a singly linked list implementation that provides O(1) operations
// at both ends and rich data manipulation capabilities. It's perfect for implementing
// stacks, queues, and ordered sequences where you primarily access elements from the ends.
//
// Linked supports both immutable operations (returning slices) and mutable operations
// (modifying the list in place) to fit different programming styles.
//
// Example usage:
//
//	// Create a task queue
//	tasks := lists.NewLinked("design", "implement", "test")
//
//	// Stack operations (LIFO)
//	tasks.PushInPlace("deploy")
//	lastTask, found := tasks.PopInPlace()  // "deploy"
//
//	// Queue operations (FIFO)
//	tasks.EnqueueInPlace("monitor")
//	firstTask, found := tasks.DequeueInPlace()  // "design"
//
//	// Rich operations
//	longTasks := tasks.Filter(func(task string) bool {
//		return len(task) > 4
//	})
type Linked[T any] struct {
	head       *node[T]
	tail       *node[T]
	isCircular bool
}

// NewLinked creates a new singly linked list with the given elements.
// Elements are added in the order provided.
//
// Example:
//
//	// Empty list
//	empty := lists.NewLinked[int]()
//
//	// List with initial elements
//	numbers := lists.NewLinked(1, 2, 3, 4, 5)
//
//	// List from slice
//	items := []string{"apple", "banana", "cherry"}
//	fruits := lists.NewLinked(items...)
func NewLinked[T any](values ...T) *Linked[T] {
	linked := &Linked[T]{}

	for _, value := range values {
		linked.insertAtEnd(value)
	}

	return linked
}

// NewLinkedCircular creates a new circular singly linked list with the given
// elements. In a circular list the tail's next pointer references the head, so
// iteration wraps around; the list's own methods stop after a single pass.
// Elements are added in the order provided.
func NewLinkedCircular[T any](values ...T) *Linked[T] {
	linked := &Linked[T]{
		isCircular: true,
	}

	for _, value := range values {
		linked.insertAtEnd(value)
	}

	return linked
}

func (l *Linked[T]) insertAtEnd(value T) {
	newNode := &node[T]{
		value:  value,
		linked: l,
		next:   nil,
	}

	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		l.tail = newNode
	}

	if l.isCircular {
		l.tail.next = l.head
	}
}

// Interface guards to ensure Linked implements the required interfaces
var _ List[int] = &Linked[int]{}
var _ MutableList[int] = &Linked[int]{}

// AllMatch returns true if all elements satisfy the given predicate.
func (l *Linked[T]) AllMatch(fn func(T) bool) bool {
	current := l.head
	for current != nil {
		if !fn(current.value) {
			return false
		}
		current = current.next
		if l.isCircular && current == l.head {
			break
		}
	}
	return true
}

// AnyMatch returns true if any element satisfies the given predicate.
func (l *Linked[T]) AnyMatch(fn func(T) bool) bool {
	current := l.head
	for current != nil {
		if fn(current.value) {
			return true
		}
		current = current.next
		if l.isCircular && current == l.head {
			break
		}
	}
	return false
}

// NoneMatch returns true if no element satisfies the given predicate.
func (l *Linked[T]) NoneMatch(fn func(T) bool) bool {
	return !l.AnyMatch(fn)
}

// Find returns the first element that satisfies the given predicate.
func (l *Linked[T]) Find(fn func(T) bool) (T, bool) {
	current := l.head
	for current != nil {
		if fn(current.value) {
			return current.value, true
		}
		current = current.next
		if l.isCircular && current == l.head {
			break
		}
	}
	var zero T
	return zero, false
}

// FindIndex returns the index of the first element that satisfies the given predicate.
func (l *Linked[T]) FindIndex(fn func(T) bool) int {
	current := l.head
	index := 0
	for current != nil {
		if fn(current.value) {
			return index
		}
		current = current.next
		index++
		if l.isCircular && current == l.head {
			break
		}
	}
	return -1
}

// Filter returns a new List containing only elements that satisfy the predicate.
func (l *Linked[T]) Filter(fn func(T) bool) List[T] {
	// Initialise non-nil so an empty result yields an initialised, non-nil empty
	// List, matching slices.Filter and the slice-backed Array implementation.
	result := []T{}
	current := l.head
	for current != nil {
		if fn(current.value) {
			result = append(result, current.value)
		}
		current = current.next
		if l.isCircular && current == l.head {
			break
		}
	}
	return NewArray(result...)
}

// FilterInPlace removes elements that don't satisfy the predicate.
func (l *Linked[T]) FilterInPlace(fn func(T) bool) {
	if l.head == nil {
		return
	}

	// Handle head nodes that don't match
	for l.head != nil && !fn(l.head.value) {
		l.head = l.head.next
		if l.isCircular && l.head == l.tail {
			l.head = nil
			l.tail = nil
			return
		}
	}

	if l.head == nil {
		l.tail = nil
		return
	}

	// Handle remaining nodes
	current := l.head
	for current.next != nil && (!l.isCircular || current.next != l.head) {
		if !fn(current.next.value) {
			if current.next == l.tail {
				l.tail = current
			}
			current.next = current.next.next
		} else {
			current = current.next
		}
	}

	if l.isCircular && l.tail != nil {
		l.tail.next = l.head
	}
}

// Get returns the element at the given index and true, or defaultValue and
// false if the index is out of bounds.
func (l *Linked[T]) Get(index int, defaultValue T) (T, bool) {
	if index < 0 {
		return defaultValue, false
	}

	current := l.head
	currentIndex := 0
	for current != nil {
		if currentIndex == index {
			return current.value, true
		}
		current = current.next
		currentIndex++
		if l.isCircular && current == l.head {
			break
		}
	}
	return defaultValue, false
}

// Length returns the number of elements in the list.
func (l *Linked[T]) Length() int {
	if l.head == nil {
		return 0
	}

	count := 0
	current := l.head
	for current != nil {
		count++
		current = current.next
		if l.isCircular && current == l.head {
			break
		}
	}
	return count
}

// IsEmpty returns true if the list contains no elements.
func (l *Linked[T]) IsEmpty() bool {
	return l.head == nil
}

// RemoveAt returns a new List with the element at index removed, without
// modifying the receiver. If index is out of bounds the elements are returned
// unchanged. AsSlice already allocates a fresh slice, so the element is
// deleted in place on it without a second copy.
func (l *Linked[T]) RemoveAt(index int) List[T] {
	return NewArray(deleteOwned(l.AsSlice(), index)...)
}

// Remove returns a new List with the first element deeply equal to element
// removed, without modifying the receiver. If no element matches, the elements
// are returned unchanged.
func (l *Linked[T]) Remove(element T) List[T] {
	slice := l.AsSlice()
	return NewArray(deleteOwned(slice, indexOfDeepEqual(slice, element))...)
}

// RemoveAtInPlace removes the element at index, returning it and whether the
// index was in bounds, modifying the receiver.
func (l *Linked[T]) RemoveAtInPlace(index int) (T, bool) {
	return l.removeFirst(func(i int, _ T) bool { return i == index })
}

// RemoveInPlace removes the first element deeply equal to element, reporting
// whether an element was removed, modifying the receiver.
func (l *Linked[T]) RemoveInPlace(element T) bool {
	_, ok := l.removeFirst(func(_ int, value T) bool {
		return reflect.DeepEqual(value, element)
	})
	return ok
}

// removeFirst unlinks the first node whose index and value satisfy match,
// returning the removed value and whether a node was removed. It maintains the
// circular invariant when the list is circular.
func (l *Linked[T]) removeFirst(match func(index int, value T) bool) (T, bool) {
	var zero T
	if l.head == nil {
		return zero, false
	}

	if match(0, l.head.value) {
		value := l.head.value
		if l.head == l.tail {
			l.head = nil
			l.tail = nil
		} else {
			l.head = l.head.next
			if l.isCircular {
				l.tail.next = l.head
			}
		}
		return value, true
	}

	prev := l.head
	current := l.head.next
	index := 1
	for current != nil && (!l.isCircular || current != l.head) {
		if match(index, current.value) {
			value := current.value
			prev.next = current.next
			if current == l.tail {
				l.tail = prev
				if l.isCircular {
					l.tail.next = l.head
				}
			}
			return value, true
		}
		prev = current
		current = current.next
		index++
	}
	return zero, false
}

// Clear removes all elements from the list.
func (l *Linked[T]) Clear() {
	l.head = nil
	l.tail = nil
}

// ForEach executes the given function for each element.
func (l *Linked[T]) ForEach(fn EachFunc[T]) {
	current := l.head
	for current != nil {
		fn(current.value)
		current = current.next
		if l.isCircular && current == l.head {
			break
		}
	}
}

// ForEachWithIndex executes the given function for each element with its index.
func (l *Linked[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	current := l.head
	index := 0
	for current != nil {
		fn(index, current.value)
		current = current.next
		index++
		if l.isCircular && current == l.head {
			break
		}
	}
}

// AsSlice returns the list as a slice.
func (l *Linked[T]) AsSlice() []T {
	var result []T
	current := l.head
	for current != nil {
		result = append(result, current.value)
		current = current.next
		if l.isCircular && current == l.head {
			break
		}
	}
	return result
}

// Insert creates a new List with elements inserted at the given index. The
// index may range over 0 <= index <= Length(): an index equal to the length
// appends. An out-of-range index leaves the list's elements unchanged.
func (l *Linked[T]) Insert(index int, elements ...T) List[T] {
	slice := l.AsSlice()
	if index < 0 || index > len(slice) {
		return NewArray(slice...)
	}

	result := make([]T, 0, len(slice)+len(elements))
	result = append(result, slice[:index]...)
	result = append(result, elements...)
	result = append(result, slice[index:]...)
	return NewArray(result...)
}

// InsertInPlace inserts elements at the given index. The index may range over
// 0 <= index <= Length(): an index equal to the length appends. An out-of-range
// index leaves the list untouched.
func (l *Linked[T]) InsertInPlace(index int, elements ...T) {
	if index < 0 {
		return
	}

	if index == 0 {
		l.insertAtHead(elements)
		return
	}

	// Walk to the node after which the elements should be spliced in.
	current, ok := l.findNodeBefore(index)
	if !ok {
		return // Invalid index in circular list — leave the list untouched.
	}
	if current == nil {
		// Index is beyond the end of the list (index > len): out of range, so
		// leave the list untouched. Appending here is reserved for index == len,
		// which resolves to the tail node and falls through to insertAfter below.
		return
	}

	l.insertAfter(current, elements)
}

// insertAtHead prepends elements (preserving their order) to the front of the list.
func (l *Linked[T]) insertAtHead(elements []T) {
	for i := len(elements) - 1; i >= 0; i-- {
		newNode := &node[T]{
			value:  elements[i],
			next:   l.head,
			linked: l,
		}
		l.head = newNode
		if l.tail == nil {
			l.tail = newNode
		}
	}
	if l.isCircular && l.tail != nil {
		l.tail.next = l.head
	}
}

// findNodeBefore walks to the node at index-1. It returns (nil, true) when the
// index is beyond the list (caller should append) and (nil, false) when a
// circular list wraps back to the head before reaching the index (invalid).
func (l *Linked[T]) findNodeBefore(index int) (*node[T], bool) {
	current := l.head
	currentIndex := 0
	for current != nil && currentIndex < index-1 {
		current = current.next
		currentIndex++
		if l.isCircular && current == l.head {
			return nil, false
		}
	}
	return current, true
}

// insertAfter splices elements in immediately after the given node.
func (l *Linked[T]) insertAfter(current *node[T], elements []T) {
	for _, element := range elements {
		newNode := &node[T]{
			value:  element,
			next:   current.next,
			linked: l,
		}
		current.next = newNode
		if current == l.tail {
			l.tail = newNode
		}
		current = newNode
	}
	if l.isCircular && l.tail != nil {
		l.tail.next = l.head
	}
}

// Sort returns a new sorted List.
func (l *Linked[T]) Sort(lessThan func(T, T) bool) List[T] {
	slice := l.AsSlice()
	// Simple bubble sort for demonstration
	for i := 0; i < len(slice); i++ {
		for j := 0; j < len(slice)-1-i; j++ {
			if lessThan(slice[j+1], slice[j]) {
				slice[j], slice[j+1] = slice[j+1], slice[j]
			}
		}
	}
	return NewArray(slice...)
}

// SortInPlace sorts the list in place.
func (l *Linked[T]) SortInPlace(lessThan func(T, T) bool) {
	if l.head == nil || l.head.next == nil {
		return
	}

	// Convert to slice, sort, and rebuild
	sorted := l.Sort(lessThan).AsSlice()
	l.head = nil
	l.tail = nil
	for _, element := range sorted {
		l.insertAtEnd(element)
	}
}

// Push adds an element to the end and returns a new List.
func (l *Linked[T]) Push(element T) List[T] {
	slice := l.AsSlice()
	return NewArray(append(slice, element)...)
}

// PushInPlace adds an element to the end of the list (stack operation).
// This is a mutable operation that modifies the list in place.
// Use this for implementing stacks (LIFO - Last In, First Out).
//
// Example:
//
//	stack := lists.NewLinked[int]()
//	stack.PushInPlace(1)
//	stack.PushInPlace(2)
//	stack.PushInPlace(3)
//	// Stack now contains: [1, 2, 3] (3 is at the top)
func (l *Linked[T]) PushInPlace(element T) {
	l.insertAtEnd(element)
}

// Pop removes and returns the last element (stack operation).
// This is an immutable operation that returns a new List.
// Returns the removed element, whether it was found, and the new List.
//
// Example:
//
//	stack := lists.NewLinked(1, 2, 3)
//	element, found, newList := stack.Pop()
//	// element: 3, found: true, newList: [1, 2]
//	// Original stack unchanged
func (l *Linked[T]) Pop() (T, bool, List[T]) {
	slice := l.AsSlice()
	if len(slice) == 0 {
		var zero T
		return zero, false, NewArray(slice...)
	}
	return slice[len(slice)-1], true, NewArray(slice[:len(slice)-1]...)
}

// PopInPlace removes and returns the last element (stack operation).
// This is a mutable operation that modifies the list in place.
// Returns the removed element and whether it was found.
//
// Example:
//
//	stack := lists.NewLinked(1, 2, 3)
//	element, found := stack.PopInPlace()
//	// element: 3, found: true
//	// Stack now contains: [1, 2]
func (l *Linked[T]) PopInPlace() (T, bool) {
	if l.head == nil {
		var zero T
		return zero, false
	}

	if l.head == l.tail {
		// Only one element
		value := l.head.value
		l.head = nil
		l.tail = nil
		return value, true
	}

	// Find second to last node
	current := l.head
	for current.next != l.tail {
		current = current.next
	}

	value := l.tail.value
	l.tail = current
	l.tail.next = nil
	if l.isCircular {
		l.tail.next = l.head
	}
	return value, true
}

// PeekEnd returns the last element without removing it.
func (l *Linked[T]) PeekEnd() (T, bool) {
	if l.tail == nil {
		var zero T
		return zero, false
	}
	return l.tail.value, true
}

// Enqueue adds an element to the end and returns a new List.
func (l *Linked[T]) Enqueue(element T) List[T] {
	return l.Push(element)
}

// EnqueueInPlace adds an element to the end.
func (l *Linked[T]) EnqueueInPlace(element T) {
	l.PushInPlace(element)
}

// Dequeue removes and returns the first element.
func (l *Linked[T]) Dequeue() (T, bool, List[T]) {
	slice := l.AsSlice()
	if len(slice) == 0 {
		var zero T
		return zero, false, NewArray(slice...)
	}
	return slice[0], true, NewArray(slice[1:]...)
}

// DequeueInPlace removes and returns the first element.
func (l *Linked[T]) DequeueInPlace() (T, bool) {
	if l.head == nil {
		var zero T
		return zero, false
	}

	value := l.head.value
	l.head = l.head.next
	if l.head == nil {
		l.tail = nil
	} else if l.isCircular {
		l.tail.next = l.head
	}
	return value, true
}

// PeekFront returns the first element without removing it.
func (l *Linked[T]) PeekFront() (T, bool) {
	if l.head == nil {
		var zero T
		return zero, false
	}
	return l.head.value, true
}
