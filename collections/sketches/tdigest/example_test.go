package tdigest_test

import (
	"fmt"
	"sync"

	"github.com/pickeringtech/go-collections/collections/sketches/tdigest"
)

func ExampleDigest() {
	d, _ := tdigest.New()
	for i := 1; i <= 1000; i++ {
		d.Add(float64(i))
	}
	// The median of 1..1000 is ~500; the t-digest estimate is close.
	median, _ := d.Quantile(0.5)
	within := median > 490 && median < 510
	fmt.Println("median within [490,510]:", within)
	// Output:
	// median within [490,510]: true
}

func ExampleDigest_Merge() {
	west, _ := tdigest.New()
	east, _ := tdigest.New()
	for i := 0; i < 1000; i++ {
		west.Add(float64(i))        // 0..999
		east.Add(float64(i + 1000)) // 1000..1999
	}

	_ = west.Merge(east) // west now summarises both shards
	min, _ := west.Min()
	max, _ := west.Max()
	fmt.Println(min, max)
	// Output:
	// 0 1999
}

func ExampleConcurrentDigest() {
	c, _ := tdigest.NewConcurrent()

	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 2500; i++ {
				c.Add(float64(g*2500 + i))
			}
		}(g)
	}
	wg.Wait()

	// 10000 values 0..9999; the 99th percentile is near 9900.
	p99, _ := c.Percentile(99)
	within := p99 > 9700 && p99 < 10000
	fmt.Println("count:", c.Count())
	fmt.Println("p99 within [9700,10000):", within)
	// Output:
	// count: 10000
	// p99 within [9700,10000): true
}
