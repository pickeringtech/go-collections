package heaps_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/collections/heaps"
)

func ExampleNewMin() {
	pq := heaps.NewMin(5, 1, 3, 2, 4)

	top, _ := pq.Peek()
	fmt.Println("smallest:", top)
	fmt.Println("drained:", pq.AsSortedSlice())

	// Output:
	// smallest: 1
	// drained: [1 2 3 4 5]
}

func ExampleNewMax() {
	pq := heaps.NewMax(5, 1, 3, 2, 4)

	top, _ := pq.Peek()
	fmt.Println("largest:", top)
	fmt.Println("drained:", pq.AsSortedSlice())

	// Output:
	// largest: 5
	// drained: [5 4 3 2 1]
}

func ExampleNew() {
	type task struct {
		name     string
		priority int
	}
	// Highest priority leaves the heap first.
	pq := heaps.New(func(a, b task) bool { return a.priority > b.priority })
	pq.PushInPlace(task{"email", 1})
	pq.PushInPlace(task{"deploy", 10})
	pq.PushInPlace(task{"review", 5})

	for next := range pq.Drain() {
		fmt.Println(next.name)
	}

	// Output:
	// deploy
	// review
	// email
}

func ExampleBinary_PopInPlace() {
	pq := heaps.NewMin(3, 1, 2)

	for {
		v, ok := pq.PopInPlace()
		if !ok {
			break
		}
		fmt.Print(v, " ")
	}
	fmt.Println()

	// Output:
	// 1 2 3
}

func ExampleBinary_Push() {
	original := heaps.NewMin(2, 4)

	bigger := original.Push(1)

	fmt.Println("original length:", original.Length())
	top, _ := bigger.Peek()
	fmt.Println("new heap smallest:", top)

	// Output:
	// original length: 2
	// new heap smallest: 1
}

func ExampleBinary_Drain() {
	pq := heaps.NewMin(30, 10, 20)

	for v := range pq.Drain() {
		fmt.Println(v)
	}
	// Drain does not consume the heap.
	fmt.Println("length after drain:", pq.Length())

	// Output:
	// 10
	// 20
	// 30
	// length after drain: 3
}

func ExampleNewConcurrentMin() {
	pq := heaps.NewConcurrentMin(5, 1, 3)
	pq.PushInPlace(0)

	v, _ := pq.PopInPlace()
	fmt.Println("highest priority:", v)

	// Output:
	// highest priority: 0
}
