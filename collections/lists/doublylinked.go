package lists

// doublyNode represents a single node in the doubly linked list.
type doublyNode[T any] struct {
	value T
	next  *doublyNode[T]
	prev  *doublyNode[T]
}

// DoublyLinked is a doubly linked list implementation.
type DoublyLinked[T any] struct {
	head       *doublyNode[T]
	tail       *doublyNode[T]
	size       int
	isCircular bool
}

// NewDoublyLinked creates a new doubly linked list with the given values.
func NewDoublyLinked[T any](values ...T) *DoublyLinked[T] {
	dl := &DoublyLinked[T]{}
	for _, value := range values {
		dl.PushInPlace(value)
	}
	return dl
}

// NewDoublyLinkedCircular creates a new circular doubly linked list with the given values.
func NewDoublyLinkedCircular[T any](values ...T) *DoublyLinked[T] {
	dl := &DoublyLinked[T]{isCircular: true}
	for _, value := range values {
		dl.PushInPlace(value)
	}
	return dl
}

// Interface guards to ensure DoublyLinked implements the required interfaces
var _ List[int] = &DoublyLinked[int]{}
var _ MutableList[int] = &DoublyLinked[int]{}

// AllMatch returns true if all elements satisfy the given predicate.
func (dl *DoublyLinked[T]) AllMatch(fn func(T) bool) bool {
	current := dl.head
	for current != nil {
		if !fn(current.value) {
			return false
		}
		current = current.next
		if dl.isCircular && current == dl.head {
			break
		}
	}
	return true
}

// AnyMatch returns true if any element satisfies the given predicate.
func (dl *DoublyLinked[T]) AnyMatch(fn func(T) bool) bool {
	current := dl.head
	for current != nil {
		if fn(current.value) {
			return true
		}
		current = current.next
		if dl.isCircular && current == dl.head {
			break
		}
	}
	return false
}

// Find returns the first element that satisfies the given predicate.
func (dl *DoublyLinked[T]) Find(fn func(T) bool) (T, bool) {
	current := dl.head
	for current != nil {
		if fn(current.value) {
			return current.value, true
		}
		current = current.next
		if dl.isCircular && current == dl.head {
			break
		}
	}
	var zero T
	return zero, false
}

// FindIndex returns the index of the first element that satisfies the given predicate.
func (dl *DoublyLinked[T]) FindIndex(fn func(T) bool) int {
	current := dl.head
	index := 0
	for current != nil {
		if fn(current.value) {
			return index
		}
		current = current.next
		index++
		if dl.isCircular && current == dl.head {
			break
		}
	}
	return -1
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func (dl *DoublyLinked[T]) Filter(fn func(T) bool) []T {
	var result []T
	current := dl.head
	for current != nil {
		if fn(current.value) {
			result = append(result, current.value)
		}
		current = current.next
		if dl.isCircular && current == dl.head {
			break
		}
	}
	return result
}

// FilterInPlace removes elements that don't satisfy the predicate.
func (dl *DoublyLinked[T]) FilterInPlace(fn func(T) bool) {
	current := dl.head
	for current != nil {
		next := current.next
		if !fn(current.value) {
			dl.removeNode(current)
		}
		current = next
		if dl.isCircular && current == dl.head {
			break
		}
	}
}

// removeNode removes a specific node from the list.
func (dl *DoublyLinked[T]) removeNode(node *doublyNode[T]) {
	if node == nil {
		return
	}

	if node.prev != nil {
		node.prev.next = node.next
	} else {
		dl.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		dl.tail = node.prev
	}

	if dl.isCircular && dl.head != nil && dl.tail != nil {
		dl.head.prev = dl.tail
		dl.tail.next = dl.head
	}

	dl.size--
}

// Get returns the element at the given index, or defaultValue if out of bounds.
func (dl *DoublyLinked[T]) Get(index int, defaultValue T) T {
	if index < 0 || index >= dl.size {
		return defaultValue
	}

	// Optimize by starting from head or tail
	var current *doublyNode[T]
	if index < dl.size/2 {
		// Start from head
		current = dl.head
		for i := 0; i < index; i++ {
			current = current.next
		}
	} else {
		// Start from tail
		current = dl.tail
		for i := dl.size - 1; i > index; i-- {
			current = current.prev
		}
	}

	return current.value
}

// Length returns the number of elements in the list.
func (dl *DoublyLinked[T]) Length() int {
	return dl.size
}

// ForEach executes the given function for each element.
func (dl *DoublyLinked[T]) ForEach(fn EachFunc[T]) {
	current := dl.head
	for current != nil {
		fn(current.value)
		current = current.next
		if dl.isCircular && current == dl.head {
			break
		}
	}
}

// ForEachWithIndex executes the given function for each element with its index.
func (dl *DoublyLinked[T]) ForEachWithIndex(fn IndexedEachFunc[T]) {
	current := dl.head
	index := 0
	for current != nil {
		fn(index, current.value)
		current = current.next
		index++
		if dl.isCircular && current == dl.head {
			break
		}
	}
}

// GetAsSlice returns the list as a slice.
func (dl *DoublyLinked[T]) GetAsSlice() []T {
	result := make([]T, 0, dl.size)
	current := dl.head
	for current != nil {
		result = append(result, current.value)
		current = current.next
		if dl.isCircular && current == dl.head {
			break
		}
	}
	return result
}

// Insert creates a new slice with elements inserted at the given index.
func (dl *DoublyLinked[T]) Insert(index int, elements ...T) []T {
	slice := dl.GetAsSlice()
	if index < 0 || index > len(slice) {
		return slice
	}

	result := make([]T, 0, len(slice)+len(elements))
	result = append(result, slice[:index]...)
	result = append(result, elements...)
	result = append(result, slice[index:]...)
	return result
}

// InsertInPlace inserts elements at the given index.
func (dl *DoublyLinked[T]) InsertInPlace(index int, elements ...T) {
	if index < 0 || index > dl.size {
		return
	}

	for i, element := range elements {
		dl.insertAt(index+i, element)
	}
}

// insertAt inserts a single element at the given index.
func (dl *DoublyLinked[T]) insertAt(index int, element T) {
	newNode := &doublyNode[T]{value: element}

	if dl.size == 0 {
		// Empty list
		dl.head = newNode
		dl.tail = newNode
		if dl.isCircular {
			newNode.next = newNode
			newNode.prev = newNode
		}
		dl.size++
		return
	}

	if index == 0 {
		// Insert at beginning
		newNode.next = dl.head
		dl.head.prev = newNode
		dl.head = newNode
		if dl.isCircular {
			newNode.prev = dl.tail
			dl.tail.next = newNode
		}
		dl.size++
		return
	}

	if index == dl.size {
		// Insert at end
		newNode.prev = dl.tail
		dl.tail.next = newNode
		dl.tail = newNode
		if dl.isCircular {
			newNode.next = dl.head
			dl.head.prev = newNode
		}
		dl.size++
		return
	}

	// Insert in middle
	current := dl.head
	for i := 0; i < index; i++ {
		current = current.next
	}

	newNode.next = current
	newNode.prev = current.prev
	current.prev.next = newNode
	current.prev = newNode
	dl.size++
}

// Sort returns a new sorted slice.
func (dl *DoublyLinked[T]) Sort(lessThan func(T, T) bool) []T {
	slice := dl.GetAsSlice()
	// Simple bubble sort for demonstration
	for i := 0; i < len(slice); i++ {
		for j := 0; j < len(slice)-1-i; j++ {
			if lessThan(slice[j+1], slice[j]) {
				slice[j], slice[j+1] = slice[j+1], slice[j]
			}
		}
	}
	return slice
}

// SortInPlace sorts the list in place.
func (dl *DoublyLinked[T]) SortInPlace(lessThan func(T, T) bool) {
	if dl.size <= 1 {
		return
	}

	// Convert to slice, sort, and rebuild
	sorted := dl.Sort(lessThan)
	dl.head = nil
	dl.tail = nil
	dl.size = 0
	for _, element := range sorted {
		dl.PushInPlace(element)
	}
}

// Push adds an element to the end and returns a new slice.
func (dl *DoublyLinked[T]) Push(element T) []T {
	slice := dl.GetAsSlice()
	return append(slice, element)
}

// PushInPlace adds an element to the end.
func (dl *DoublyLinked[T]) PushInPlace(element T) {
	dl.insertAt(dl.size, element)
}

// Pop removes and returns the last element.
func (dl *DoublyLinked[T]) Pop() (T, bool, []T) {
	slice := dl.GetAsSlice()
	if len(slice) == 0 {
		var zero T
		return zero, false, slice
	}
	return slice[len(slice)-1], true, slice[:len(slice)-1]
}

// PopInPlace removes and returns the last element.
func (dl *DoublyLinked[T]) PopInPlace() (T, bool) {
	if dl.size == 0 {
		var zero T
		return zero, false
	}

	value := dl.tail.value
	dl.removeNode(dl.tail)
	return value, true
}

// PeekEnd returns the last element without removing it.
func (dl *DoublyLinked[T]) PeekEnd() (T, bool) {
	if dl.tail == nil {
		var zero T
		return zero, false
	}
	return dl.tail.value, true
}

// Enqueue adds an element to the end and returns a new slice.
func (dl *DoublyLinked[T]) Enqueue(element T) []T {
	return dl.Push(element)
}

// EnqueueInPlace adds an element to the end.
func (dl *DoublyLinked[T]) EnqueueInPlace(element T) {
	dl.PushInPlace(element)
}

// Dequeue removes and returns the first element.
func (dl *DoublyLinked[T]) Dequeue() (T, bool, []T) {
	slice := dl.GetAsSlice()
	if len(slice) == 0 {
		var zero T
		return zero, false, slice
	}
	return slice[0], true, slice[1:]
}

// DequeueInPlace removes and returns the first element.
func (dl *DoublyLinked[T]) DequeueInPlace() (T, bool) {
	if dl.size == 0 {
		var zero T
		return zero, false
	}

	value := dl.head.value
	dl.removeNode(dl.head)
	return value, true
}

// PeekFront returns the first element without removing it.
func (dl *DoublyLinked[T]) PeekFront() (T, bool) {
	if dl.head == nil {
		var zero T
		return zero, false
	}
	return dl.head.value, true
}
