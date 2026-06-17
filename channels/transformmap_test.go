package channels_test

import (
	"context"
	"fmt"
	"github.com/pickeringtech/go-collections/channels"
	"reflect"
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
