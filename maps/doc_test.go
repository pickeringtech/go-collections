package maps_test

import (
	"fmt"

	"github.com/pickeringtech/go-collections/maps"
)

// Example_quickStart is the runnable twin of the package godoc Quick Start. Keep
// the two in sync: `go test` compiles and output-checks this, which is what
// guarantees the documented entry-point API actually exists and behaves as shown.
func Example_quickStart() {
	inventory := map[string]int{
		"apples":  50,
		"oranges": 30,
		"bananas": 20,
	}

	// Filter for low-stock items.
	lowStock := maps.Filter(inventory, func(item string, count int) bool {
		return count < 40
	})

	// Map rebuilds the map and can change the key, the value, or both; here it
	// doubles each value and keeps the key.
	doubled := maps.Map(inventory, func(item string, count int) (string, int) {
		return item, count * 2
	})

	for item, count := range lowStock {
		fmt.Printf("low stock: %s=%d\n", item, count)
	}
	for item, count := range doubled {
		fmt.Printf("doubled: %s=%d\n", item, count)
	}
	// Unordered output:
	// low stock: oranges=30
	// low stock: bananas=20
	// doubled: apples=100
	// doubled: oranges=60
	// doubled: bananas=40
}
