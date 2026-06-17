package collections_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/collections"
)

// ExamplePair verifies that the quick-start examples in doc.go, README.md and
// collections/README.md compile as written, using collections.Pair with only
// the facade package imported.
func ExamplePair() {
	users := collections.NewDict(
		collections.Pair[int, string]{Key: 1, Value: "Alice"},
		collections.Pair[int, string]{Key: 2, Value: "Bob"},
	)

	name, ok := users.Get(1, "")
	fmt.Println(name, ok)
	// Output: Alice true
}
