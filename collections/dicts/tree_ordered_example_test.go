package dicts_test

import (
	"fmt"
	"sync"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

func ExampleTree_All() {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 3, Value: "three"},
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	for key, value := range tree.All() {
		fmt.Printf("%d=%s\n", key, value)
	}
	// Output:
	// 1=one
	// 2=two
	// 3=three
}

func ExampleTree_Backward() {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 3, Value: "three"},
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	for key := range tree.Backward() {
		fmt.Println(key)
	}
	// Output:
	// 3
	// 2
	// 1
}

func ExampleTree_Floor() {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 10, Value: "ten"},
		dicts.Pair[int, string]{Key: 20, Value: "twenty"},
		dicts.Pair[int, string]{Key: 30, Value: "thirty"},
	)

	key, value, ok := tree.Floor(25)
	fmt.Printf("floor(25) = %d=%s (%v)\n", key, value, ok)
	// Output:
	// floor(25) = 20=twenty (true)
}

func ExampleTree_Ceiling() {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 10, Value: "ten"},
		dicts.Pair[int, string]{Key: 20, Value: "twenty"},
		dicts.Pair[int, string]{Key: 30, Value: "thirty"},
	)

	key, value, ok := tree.Ceiling(25)
	fmt.Printf("ceiling(25) = %d=%s (%v)\n", key, value, ok)
	// Output:
	// ceiling(25) = 30=thirty (true)
}

func ExampleTree_Range() {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 1, Value: "a"},
		dicts.Pair[int, string]{Key: 2, Value: "b"},
		dicts.Pair[int, string]{Key: 3, Value: "c"},
		dicts.Pair[int, string]{Key: 4, Value: "d"},
	)

	for _, pair := range tree.Range(2, 3) {
		fmt.Printf("%d=%s\n", pair.Key, pair.Value)
	}
	// Output:
	// 2=b
	// 3=c
}

func ExampleTree_RangeAll() {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 1, Value: "a"},
		dicts.Pair[int, string]{Key: 2, Value: "b"},
		dicts.Pair[int, string]{Key: 3, Value: "c"},
		dicts.Pair[int, string]{Key: 4, Value: "d"},
	)

	for key := range tree.RangeAll(2, 3) {
		fmt.Println(key)
	}
	// Output:
	// 2
	// 3
}

func ExampleTree_Min() {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 30, Value: "thirty"},
		dicts.Pair[int, string]{Key: 10, Value: "ten"},
		dicts.Pair[int, string]{Key: 20, Value: "twenty"},
	)

	key, value, _ := tree.Min()
	fmt.Printf("min = %d=%s\n", key, value)
	// Output:
	// min = 10=ten
}

func ExampleTree_Max() {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 30, Value: "thirty"},
		dicts.Pair[int, string]{Key: 10, Value: "ten"},
		dicts.Pair[int, string]{Key: 20, Value: "twenty"},
	)

	key, value, _ := tree.Max()
	fmt.Printf("max = %d=%s\n", key, value)
	// Output:
	// max = 30=thirty
}

func ExampleConcurrentTree() {
	tree := dicts.NewConcurrentTree[int, int]()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			tree.PutInPlace(n, n*n)
		}(i)
	}
	wg.Wait()

	min, minSquare, _ := tree.Min()
	max, maxSquare, _ := tree.Max()
	fmt.Printf("entries=%d min=%d:%d max=%d:%d\n", tree.Length(), min, minSquare, max, maxSquare)
	// Output:
	// entries=100 min=0:0 max=99:9801
}
