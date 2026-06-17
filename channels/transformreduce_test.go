package channels_test

import (
	"context"
	"fmt"
	"github.com/pickeringtech/go-collections/channels"
	"reflect"
	"testing"
	"time"
)

func ExampleReduce() {
	ctx := context.Background()
	input := channels.FromSlice(ctx, []int{1, 2, 3, 4, 5})

	// Creates a new pipeline which totals the input channel.
	pipeline := channels.NewPipeline[int, int](ctx, input, func(ctx context.Context, input <-chan int) <-chan int {
		return channels.Reduce(ctx, input, func(accumulator int, element int) int {
			return accumulator + element
		})
	})

	// Capture results in a slice.
	results := pipeline.CollectAsSlice()

	fmt.Printf("Results: %v", results)
	// Output: Results: [15]
}

func TestReduce(t *testing.T) {
	type args[I any, O any] struct {
		input <-chan I
		fn    channels.ReduceFunc[I, O]
	}
	type testCase[I any, O any] struct {
		name string
		args args[I, O]
		want []O
	}
	ctx := context.Background()
	tests := []testCase[int, int]{
		{
			name: "totals correctly",
			args: args[int, int]{
				input: channels.FromSlice(ctx, []int{1, 2, 3, 4, 5}),
				fn:    func(a int, b int) int { return a + b },
			},
			want: []int{15},
		},
		{
			name: "empty input results in zero output",
			args: args[int, int]{
				input: channels.FromSlice(ctx, []int{}),
				fn:    func(a int, b int) int { return a + b },
			},
			want: []int{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCh := channels.Reduce(ctx, tt.args.input, tt.args.fn)
			got := channels.CollectAsSlice(gotCh)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestReduceCancellation asserts that cancelling the context tears the Reduce goroutine down: it abandons the partial
// accumulation, closes the output channel without emitting, and returns even though the input never closes.
func TestReduceCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int) // never written to, never closed
	output := channels.Reduce(ctx, input, func(acc, el int) int { return acc + el })

	cancel()

	select {
	case _, ok := <-output:
		if ok {
			t.Fatal("Reduce() emitted a value after cancellation, want closed channel")
		}
	case <-time.After(time.Second):
		t.Fatal("Reduce() goroutine did not exit after cancellation")
	}
}
