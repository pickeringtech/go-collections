package lists

type node[T any] struct {
	value  T
	next   *node[T]
	linked *Linked[T]
}

type Linked[T any] struct {
	head       *node[T]
	tail       *node[T]
	isCircular bool
}

func NewLinked[T any](values ...T) *Linked[T] {
	linked := &Linked[T]{}

	for _, value := range values {
		linked.insertAtEnd(value)
	}

	return linked
}

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

// Filter returns a new slice containing only elements that satisfy the predicate.
func (l *Linked[T]) Filter(fn func(T) bool) []T {
	var result []T
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
	return result
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

// Get returns the element at the given index, or defaultValue if out of bounds.
func (l *Linked[T]) Get(index int, defaultValue T) T {
	if index < 0 {
		return defaultValue
	}

	current := l.head
	currentIndex := 0
	for current != nil {
		if currentIndex == index {
			return current.value
		}
		current = current.next
		currentIndex++
		if l.isCircular && current == l.head {
			break
		}
	}
	return defaultValue
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

// GetAsSlice returns the list as a slice.
func (l *Linked[T]) GetAsSlice() []T {
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

// Insert creates a new slice with elements inserted at the given index.
func (l *Linked[T]) Insert(index int, elements ...T) []T {
	slice := l.GetAsSlice()
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
func (l *Linked[T]) InsertInPlace(index int, elements ...T) {
	if index < 0 {
		return
	}

	// Insert at beginning
	if index == 0 {
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
		return
	}

	// Find insertion point
	current := l.head
	currentIndex := 0
	for current != nil && currentIndex < index-1 {
		current = current.next
		currentIndex++
		if l.isCircular && current == l.head {
			return // Invalid index in circular list
		}
	}

	if current == nil {
		// Index beyond list, append to end
		for _, element := range elements {
			l.insertAtEnd(element)
		}
		return
	}

	// Insert elements
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

// Sort returns a new sorted slice.
func (l *Linked[T]) Sort(lessThan func(T, T) bool) []T {
	slice := l.GetAsSlice()
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
func (l *Linked[T]) SortInPlace(lessThan func(T, T) bool) {
	if l.head == nil || l.head.next == nil {
		return
	}

	// Convert to slice, sort, and rebuild
	sorted := l.Sort(lessThan)
	l.head = nil
	l.tail = nil
	for _, element := range sorted {
		l.insertAtEnd(element)
	}
}

// Push adds an element to the end and returns a new slice.
func (l *Linked[T]) Push(element T) []T {
	slice := l.GetAsSlice()
	return append(slice, element)
}

// PushInPlace adds an element to the end.
func (l *Linked[T]) PushInPlace(element T) {
	l.insertAtEnd(element)
}

// Pop removes and returns the last element.
func (l *Linked[T]) Pop() (T, bool, []T) {
	slice := l.GetAsSlice()
	if len(slice) == 0 {
		var zero T
		return zero, false, slice
	}
	return slice[len(slice)-1], true, slice[:len(slice)-1]
}

// PopInPlace removes and returns the last element.
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

// Enqueue adds an element to the end and returns a new slice.
func (l *Linked[T]) Enqueue(element T) []T {
	return l.Push(element)
}

// EnqueueInPlace adds an element to the end.
func (l *Linked[T]) EnqueueInPlace(element T) {
	l.PushInPlace(element)
}

// Dequeue removes and returns the first element.
func (l *Linked[T]) Dequeue() (T, bool, []T) {
	slice := l.GetAsSlice()
	if len(slice) == 0 {
		var zero T
		return zero, false, slice
	}
	return slice[0], true, slice[1:]
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
