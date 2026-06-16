package sets_test

import (
	"fmt"
	"sync"

	"github.com/pickeringtech/go-collections/collections/sets"
)

func ExampleTreeSet_All() {
	s := sets.NewTreeSet(3, 1, 2, 1)

	for element := range s.All() {
		fmt.Println(element)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleTreeSet_Backward() {
	s := sets.NewTreeSet(3, 1, 2)

	for element := range s.Backward() {
		fmt.Println(element)
	}
	// Output:
	// 3
	// 2
	// 1
}

func ExampleTreeSet_Floor() {
	s := sets.NewTreeSet(10, 20, 30)

	element, ok := s.Floor(25)
	fmt.Printf("floor(25) = %d (%v)\n", element, ok)
	// Output:
	// floor(25) = 20 (true)
}

func ExampleTreeSet_Ceiling() {
	s := sets.NewTreeSet(10, 20, 30)

	element, ok := s.Ceiling(25)
	fmt.Printf("ceiling(25) = %d (%v)\n", element, ok)
	// Output:
	// ceiling(25) = 30 (true)
}

func ExampleTreeSet_Range() {
	s := sets.NewTreeSet(1, 2, 3, 4, 5)

	fmt.Println(s.Range(2, 4))
	// Output:
	// [2 3 4]
}

func ExampleTreeSet_Intersection() {
	a := sets.NewTreeSet(1, 2, 3, 4)
	b := sets.NewTreeSet(3, 4, 5, 6)

	fmt.Println(a.Intersection(b).AsSlice())
	// Output:
	// [3 4]
}

func ExampleConcurrentTreeSet() {
	s := sets.NewConcurrentTreeSet[int]()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			s.AddInPlace(n % 10)
		}(i)
	}
	wg.Wait()

	min, _ := s.Min()
	max, _ := s.Max()
	fmt.Printf("size=%d min=%d max=%d\n", s.Length(), min, max)
	// Output:
	// size=10 min=0 max=9
}
