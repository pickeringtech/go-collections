package collections

import "fmt"

func ExampleListBuilder_Build() {
	list := NewListBuilder[int]().
		Add(1, 2, 3).
		Concurrent().
		RW().
		Add(4, 5, 6).
		Build()

	fmt.Printf("list: %v\n", list.GetAsSlice())
	// Output: list: [1 2 3 4 5 6]
}
