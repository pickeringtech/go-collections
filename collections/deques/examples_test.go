package deques_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/collections/deques"
)

func ExampleNewRingBuffer() {
	d := deques.NewRingBuffer(1, 2, 3)
	d.PushFrontInPlace(0)
	d.PushBackInPlace(4)
	fmt.Println(d.AsSlice())
	// Output: [0 1 2 3 4]
}

func ExampleNewBoundedRingBuffer_overwriteOldest() {
	d := deques.NewBoundedRingBuffer[int](3, deques.OverwriteOldest)
	d.PushBackInPlace(1)
	d.PushBackInPlace(2)
	d.PushBackInPlace(3)
	d.PushBackInPlace(4) // evicts the front (1)
	fmt.Println(d.AsSlice())
	// Output: [2 3 4]
}

func ExampleNewBoundedRingBuffer_rejectWhenFull() {
	d := deques.NewBoundedRingBuffer[int](3, deques.RejectWhenFull)
	d.PushBackInPlace(1)
	d.PushBackInPlace(2)
	d.PushBackInPlace(3)
	accepted := d.PushBackInPlace(4) // rejected — deque is full
	fmt.Println(accepted, d.AsSlice())
	// Output: false [1 2 3]
}

func ExampleRingBuffer_PopFront() {
	d := deques.NewRingBuffer(1, 2, 3)
	value, ok, rest := d.PopFront()
	fmt.Println(value, ok, rest.AsSlice())
	fmt.Println(d.AsSlice()) // original unchanged
	// Output:
	// 1 true [2 3]
	// [1 2 3]
}

func ExampleRingBuffer_PopBackInPlace() {
	d := deques.NewRingBuffer(1, 2, 3)
	value, ok := d.PopBackInPlace()
	fmt.Println(value, ok, d.AsSlice())
	// Output: 3 true [1 2]
}

func ExampleRingBuffer_All() {
	d := deques.NewRingBuffer("a", "b", "c")
	for i, v := range d.All() {
		fmt.Printf("%d:%s ", i, v)
	}
	fmt.Println()
	// Output: 0:a 1:b 2:c
}

func ExampleRingBuffer_Backward() {
	d := deques.NewRingBuffer(1, 2, 3)
	for _, v := range d.Backward() {
		fmt.Print(v, " ")
	}
	fmt.Println()
	// Output: 3 2 1
}

func ExampleRingBuffer_Capacity() {
	unbounded := deques.NewRingBuffer[int]()
	bounded := deques.NewBoundedRingBuffer[int](5, deques.OverwriteOldest)
	fmt.Println(unbounded.Capacity(), bounded.Capacity())
	// Output: -1 5
}
