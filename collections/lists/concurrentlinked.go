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
		linked.Insert(value)
	}

	return linked
}

func NewLinkedCircular[T any](values ...T) *Linked[T] {
	linked := &Linked[T]{
		isCircular: true,
	}

	for _, value := range values {
		linked.Insert(value)
	}

	return linked
}

func (l *Linked[T]) Insert(value T) {
	newNode := &node[T]{
		value:  value,
		linked: l,
		next:   l.tail,
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
