package channels_test

import (
	"context"
	"fmt"
	"github.com/pickeringtech/go-collections/channels"
	"reflect"
	"testing"
	"time"
)

func ExampleFilter() {
	ctx := context.Background()
	input := channels.FromSlice(ctx, []string{"hello", "everyone", "world", "goodness", "gracious"})
	output := channels.Filter(ctx, input, func(element string) bool {
		return len(element) > 5
	})

	// Capture results in a slice.
	results := channels.CollectAsSlice(output)

	// Print results.
	fmt.Printf("Results: %v", results)
	// Output: Results: [everyone goodness gracious]
}

func TestFilter(t *testing.T) {
	type args[T any] struct {
		input <-chan T
		fn    channels.FilterFunc[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []string
	}
	ctx := context.Background()
	tests := []testCase[string]{
		{
			name: "filters out words with 5 characters or less",
			args: args[string]{
				input: channels.FromSlice(ctx, []string{"hello", "everyone", "world", "goodness", "gracious"}),
				fn: func(element string) bool {
					return len(element) > 5
				},
			},
			want: []string{"everyone", "goodness", "gracious"},
		},
		{
			name: "empty input provides nil output",
			args: args[string]{
				input: channels.FromSlice(ctx, []string{}),
				fn: func(element string) bool {
					return len(element) > 5
				},
			},
			want: nil,
		},
		{
			name: "nil input provides nil output",
			args: args[string]{
				input: channels.FromSlice[string](ctx, nil),
				fn: func(element string) bool {
					return len(element) > 5
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := channels.Filter(ctx, tt.args.input, tt.args.fn)
			got := channels.CollectAsSlice(output)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFilterCancellation asserts that cancelling the context tears the Filter goroutine down: it closes the output
// channel and returns even though the input channel never sends a value or closes.
func TestFilterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan int) // never written to, never closed
	output := channels.Filter(ctx, input, func(int) bool { return true })

	cancel()

	select {
	case _, ok := <-output:
		if ok {
			t.Fatal("Filter() emitted a value after cancellation, want closed channel")
		}
	case <-time.After(time.Second):
		t.Fatal("Filter() goroutine did not exit after cancellation")
	}
}
