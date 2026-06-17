package collections_test

import (
	"fmt"
	"slices"

	"github.com/pickeringtech/go-collections/collections"
)

// Example_quickStart is the runnable twin of the package godoc overview. Keep
// the two in sync: `go test` compiles and output-checks this, which is what
// guarantees the documented facade API actually exists and behaves as shown.
func Example_quickStart() {
	// Dicts — key/value mappings with rich operations beyond native maps.
	users := collections.NewDict(
		collections.Pair[int, string]{Key: 1, Value: "Alice"},
		collections.Pair[int, string]{Key: 2, Value: "Bob"},
	)
	name, _ := users.Get(1, "")

	// Sets — unique elements with mathematical operations.
	perms := collections.NewSet("read", "write")

	// Lists — ordered sequences; immutable transforms return a new List.
	nums := collections.NewList(1, 2, 3, 4)
	evens := nums.Filter(func(n int) bool { return n%2 == 0 })

	fmt.Println(name, perms.Contains("write"), evens.AsSlice())
	// Output: Alice true [2 4]
}

// Example_iteratorBridge is the runnable twin of the README's Go 1.23+ iterator
// bridge section: ListFromSeq, DictFromSeq2 and SetFromSeq build collections
// from any range-over-func iterator.
func Example_iteratorBridge() {
	list := collections.ListFromSeq(slices.Values([]int{3, 1, 2}))

	set := collections.SetFromSeq(slices.Values([]string{"a", "b", "a"}))

	fmt.Println(list.AsSlice(), set.Contains("a"), set.Length())
	// Output: [3 1 2] true 2
}
