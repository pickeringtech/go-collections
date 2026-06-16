// Command ordered-processing runs a list of integers through three ordered
// views of the same data and prints each result.
//
// It demonstrates the lists package end-to-end on a mutable Array:
//   - a stack reverses the input (LIFO push/pop);
//   - a queue replays it in arrival order (FIFO enqueue/dequeue);
//   - an in-place sort orders it ascending.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pickeringtech/go-collections/collections/lists"
	"github.com/pickeringtech/go-collections/slices"
)

func main() {
	raw := flag.String("nums", "5,3,8,1,9,2", "comma-separated integers to process")
	flag.Parse()

	nums, err := parse(*raw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ordered-processing:", err)
		os.Exit(1)
	}

	fmt.Printf("input:    %v\n", nums)
	fmt.Printf("reversed: %v\n", reverseViaStack(nums))
	fmt.Printf("fifo:     %v\n", replayViaQueue(nums))
	fmt.Printf("sorted:   %v\n", sortAscending(nums))
}

// reverseViaStack pushes every number onto a stack and pops them off, yielding
// the input in reverse (LIFO) order.
func reverseViaStack(nums []int) []int {
	stack := lists.NewArray[int]()
	for _, n := range nums {
		stack.PushInPlace(n)
	}
	out := make([]int, 0, len(nums))
	for !stack.IsEmpty() {
		v, _ := stack.PopInPlace()
		out = append(out, v)
	}
	return out
}

// replayViaQueue enqueues every number and dequeues them, yielding the input in
// its original arrival (FIFO) order.
func replayViaQueue(nums []int) []int {
	queue := lists.NewArray(nums...)
	out := make([]int, 0, len(nums))
	for !queue.IsEmpty() {
		v, _ := queue.DequeueInPlace()
		out = append(out, v)
	}
	return out
}

// sortAscending sorts a list of numbers in place and returns the result.
func sortAscending(nums []int) []int {
	list := lists.NewArray(nums...)
	list.SortInPlace(slices.AscendingSortFunc[int])
	return list.AsSlice()
}

// parse converts "5, 3 ,8" into [5 3 8], trimming spaces and dropping blanks.
func parse(s string) ([]int, error) {
	out := []int{}
	for _, field := range strings.Split(s, ",") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		n, err := strconv.Atoi(field)
		if err != nil {
			return nil, fmt.Errorf("%q is not an integer", field)
		}
		out = append(out, n)
	}
	return out, nil
}
