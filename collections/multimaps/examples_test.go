package multimaps_test

import (
	"fmt"
	"sort"
	"sync"

	"github.com/pickeringtech/go-collections/collections/multimaps"
)

// ExampleListMultimap demonstrates the list-backed multimap, which preserves
// insertion order and keeps duplicate values.
func ExampleListMultimap() {
	orders := multimaps.NewListMultimap[string, string]()
	orders.PutInPlace("alice", "book")
	orders.PutInPlace("alice", "pen")
	orders.PutInPlace("alice", "book") // duplicate kept

	fmt.Println(orders.Get("alice"))
	fmt.Println("entries:", orders.Length())
	fmt.Println("keys:", orders.KeyCount())
	// Output:
	// [book pen book]
	// entries: 3
	// keys: 1
}

// ExampleSetMultimap demonstrates the set-backed multimap, which discards
// duplicate values per key.
func ExampleSetMultimap() {
	tags := multimaps.NewSetMultimap[string, string]()
	tags.PutInPlace("doc1", "go")
	tags.PutInPlace("doc1", "go") // duplicate ignored
	tags.PutInPlace("doc1", "testing")

	fmt.Println("entries:", tags.Length())
	fmt.Println("has go:", tags.ContainsEntry("doc1", "go"))
	// Output:
	// entries: 2
	// has go: true
}

// ExampleNewListMultimap demonstrates seeding a multimap with entries.
func ExampleNewListMultimap() {
	scores := multimaps.NewListMultimap(
		multimaps.Entry[string, int]{Key: "alice", Value: 10},
		multimaps.Entry[string, int]{Key: "alice", Value: 20},
		multimaps.Entry[string, int]{Key: "bob", Value: 5},
	)

	fmt.Println(scores.Get("alice"))
	fmt.Println(scores.Get("bob"))
	// Output:
	// [10 20]
	// [5]
}

// ExampleListMultimap_Filter demonstrates filtering entries; keys left with no
// surviving values are dropped.
func ExampleListMultimap_Filter() {
	m := multimaps.NewListMultimap[string, int]()
	m.PutAllInPlace("a", 1, 2, 3, 4)

	even := m.Filter(func(_ string, value int) bool { return value%2 == 0 })
	fmt.Println(even.Get("a"))
	// Output:
	// [2 4]
}

// ExampleListMultimap_All demonstrates iterating over every entry with
// range-over-func.
func ExampleListMultimap_All() {
	m := multimaps.NewListMultimap[string, int]()
	m.PutAllInPlace("a", 1, 2, 3)

	sum := 0
	for _, value := range m.All() {
		sum += value
	}
	fmt.Println("sum:", sum)
	// Output:
	// sum: 6
}

// ExampleSetMultimap_Keys demonstrates listing the distinct keys.
func ExampleSetMultimap_Keys() {
	m := multimaps.NewSetMultimap[string, int]()
	m.PutInPlace("a", 1)
	m.PutInPlace("b", 2)
	m.PutInPlace("c", 3)

	keys := m.Keys()
	sort.Strings(keys)
	fmt.Println(keys)
	// Output:
	// [a b c]
}

// ExampleConcurrentListMultimap demonstrates safe concurrent writes from many
// goroutines. Run the suite with -race to verify thread safety.
func ExampleConcurrentListMultimap() {
	m := multimaps.NewConcurrentListMultimap[string, int]()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			m.PutInPlace("events", n)
		}(i)
	}
	wg.Wait()

	fmt.Println("entries:", m.Length())
	// Output:
	// entries: 100
}
