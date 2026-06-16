package lru_test

import (
	"fmt"
	"sync"

	"github.com/pickeringtech/go-collections/collections/lru"
)

func ExampleNewLRU() {
	// A cache bounded to two entries.
	cache := lru.NewLRU[string, int](2)

	cache.PutInPlace("a", 1)
	cache.PutInPlace("b", 2)
	cache.PutInPlace("c", 3) // exceeds capacity → evicts the least-recently-used "a"

	fmt.Println("contains a:", cache.Contains("a"))
	fmt.Println("keys:", cache.Keys())

	// Output:
	// contains a: false
	// keys: [c b]
}

func ExampleLRU_Get() {
	cache := lru.NewLRU[string, int](2, lru.WithEntries(
		lru.Pair[string, int]{Key: "a", Value: 1},
		lru.Pair[string, int]{Key: "b", Value: 2},
	))

	// Get marks "a" as recently used, so "b" becomes the eviction candidate.
	value, found := cache.Get("a")
	fmt.Printf("a=%d found=%t\n", value, found)

	cache.PutInPlace("c", 3) // evicts "b", not "a"
	fmt.Println("contains b:", cache.Contains("b"))
	fmt.Println("keys:", cache.Keys())

	// Output:
	// a=1 found=true
	// contains b: false
	// keys: [c a]
}

func ExampleLRU_Peek() {
	cache := lru.NewLRU[string, int](2, lru.WithEntries(
		lru.Pair[string, int]{Key: "a", Value: 1},
		lru.Pair[string, int]{Key: "b", Value: 2},
	))

	// Peek reads without affecting recency, so "a" stays the eviction candidate.
	value, found := cache.Peek("a")
	fmt.Printf("a=%d found=%t\n", value, found)

	cache.PutInPlace("c", 3) // evicts "a"
	fmt.Println("contains a:", cache.Contains("a"))

	// Output:
	// a=1 found=true
	// contains a: false
}

func ExampleWithOnEvict() {
	cache := lru.NewLRU[string, int](1, lru.WithOnEvict(func(key string, value int) {
		fmt.Printf("evicted %s=%d\n", key, value)
	}))

	cache.PutInPlace("a", 1)
	cache.PutInPlace("b", 2) // evicts "a", firing the callback

	// Output:
	// evicted a=1
}

func ExampleLRU_All() {
	cache := lru.NewLRU[string, int](3, lru.WithEntries(
		lru.Pair[string, int]{Key: "a", Value: 1},
		lru.Pair[string, int]{Key: "b", Value: 2},
		lru.Pair[string, int]{Key: "c", Value: 3},
	))

	// Iterate most- to least-recently-used.
	for key, value := range cache.All() {
		fmt.Printf("%s=%d\n", key, value)
	}

	// Output:
	// c=3
	// b=2
	// a=1
}

func ExampleLRU_Put() {
	base := lru.NewLRU[string, int](2, lru.WithEntries(
		lru.Pair[string, int]{Key: "a", Value: 1},
		lru.Pair[string, int]{Key: "b", Value: 2},
	))

	// Put returns a new cache; the original is untouched.
	updated := base.Put("c", 3)

	fmt.Println("base has a:", base.Contains("a"))
	fmt.Println("updated has a:", updated.Contains("a"))

	// Output:
	// base has a: true
	// updated has a: false
}

func ExampleNewConcurrentLRU() {
	cache := lru.NewConcurrentLRU[int, int](100)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			cache.PutInPlace(n, n*n)
		}(i)
	}
	wg.Wait()

	value, found := cache.Get(7)
	fmt.Printf("7*7=%d found=%t\n", value, found)
	fmt.Println("length:", cache.Length())

	// Output:
	// 7*7=49 found=true
	// length: 50
}
