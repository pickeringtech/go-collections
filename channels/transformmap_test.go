package channels_test

import (
	"context"
	"fmt"
	"github.com/pickeringtech/go-collections/channels"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func ExampleMap() {
	ctx := context.Background()
	input := channels.FromSlice(ctx, []string{"one", "two", "three", "four", "five"})
	output := channels.Map(ctx, input, func(s string) int {
		return len(s)
	})

	// Capture results in a slice.
	results := channels.CollectAsSlice(output)

	// Print results.
	fmt.Printf("Results: %v", results)
	// Output: Results: [3 3 5 4 4]
}

func TestMap(t *testing.T) {
	type args[I any, O any] struct {
		input <-chan I
		fn    channels.MapFunc[I, O]
	}
	type testCase[I any, O any] struct {
		name string
		args args[I, O]
		want []O
	}
	ctx := context.Background()
	tests := []testCase[string, int]{
		{
			name: "transforms elements correctly",
			args: args[string, int]{
				input: channels.FromSlice[string](ctx, []string{"one", "two", "three", "four", "five"}),
				fn: func(s string) int {
					return len(s)
				},
			},
			want: []int{3, 3, 5, 4, 4},
		},
		{
			name: "empty input provides nil output",
			args: args[string, int]{
				input: channels.FromSlice[string](ctx, []string{}),
				fn: func(s string) int {
					return len(s)
				},
			},
			want: nil,
		},
		{
			name: "nil input provides nil output",
			args: args[string, int]{
				input: channels.FromSlice[string](ctx, nil),
				fn: func(s string) int {
					return len(s)
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := channels.Map(ctx, tt.args.input, tt.args.fn)
			got := channels.CollectAsSlice(output)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMapCancellation asserts that cancelling the context tears the Map goroutine down: it closes the output channel
// and returns even though the input channel never sends a value or closes, so the goroutine is reclaimed deterministically.
func TestMapCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int) // never written to, never closed
	output := channels.Map(ctx, input, func(i int) int { return i })

	cancel()

	select {
	case _, ok := <-output:
		if ok {
			t.Fatal("Map() emitted a value after cancellation, want closed channel")
		}
	case <-time.After(time.Second):
		t.Fatal("Map() goroutine did not exit after cancellation")
	}
}

// TestMapCancellationWhileSending covers the other cancellation path: the goroutine
// has read a value and is blocked trying to deliver the result downstream (the
// output is unbuffered and unread). Cancelling unblocks that send, so the goroutine
// abandons the value and exits rather than leaking.
func TestMapCancellationWhileSending(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int, 1)
	input <- 1 // a value is ready, so the goroutine proceeds to the send

	output := channels.Map(ctx, input, func(i int) int { return i * 2 })

	// Yield generously so the goroutine consumes the input (context still live, so
	// the read wins the select) and parks in the blocked send before we cancel.
	for i := 0; i < 1000; i++ {
		runtime.Gosched()
	}
	cancel()

	select {
	case v, ok := <-output:
		if ok {
			t.Fatalf("Map() delivered %d after cancellation, want closed channel", v)
		}
	case <-time.After(time.Second):
		t.Fatal("Map() goroutine did not exit after cancellation while sending")
	}
}
