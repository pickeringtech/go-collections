package slices_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/slices"
)

func ExampleChunk() {
	input := []int{1, 2, 3, 4, 5}
	output := slices.Chunk(input, 2)
	fmt.Printf("Output: %v\n", output)

	// Output: Output: [[1 2] [3 4] [5]]
}

func ExampleWindow() {
	input := []int{1, 2, 3, 4}
	output := slices.Window(input, 2)
	fmt.Printf("Output: %v\n", output)

	// Output: Output: [[1 2] [2 3] [3 4]]
}

func TestChunk(t *testing.T) {
	type args struct {
		input []int
		size  int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "even split",
			args: args{input: []int{1, 2, 3, 4}, size: 2},
			want: [][]int{{1, 2}, {3, 4}},
		},
		{
			name: "trailing partial chunk",
			args: args{input: []int{1, 2, 3, 4, 5}, size: 2},
			want: [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name: "size larger than input yields single chunk",
			args: args{input: []int{1, 2, 3}, size: 10},
			want: [][]int{{1, 2, 3}},
		},
		{
			name: "size of one chunks each element",
			args: args{input: []int{1, 2, 3}, size: 1},
			want: [][]int{{1}, {2}, {3}},
		},
		{
			name: "zero size yields empty output",
			args: args{input: []int{1, 2, 3}, size: 0},
			want: [][]int{},
		},
		{
			name: "negative size yields empty output",
			args: args{input: []int{1, 2, 3}, size: -1},
			want: [][]int{},
		},
		{
			name: "nil input yields non-nil empty output",
			args: args{input: nil, size: 2},
			want: [][]int{},
		},
		{
			name: "empty input yields non-nil empty output",
			args: args{input: []int{}, size: 2},
			want: [][]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Chunk(tt.args.input, tt.args.size)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindow(t *testing.T) {
	type args struct {
		input []int
		size  int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "width two slides one at a time",
			args: args{input: []int{1, 2, 3, 4}, size: 2},
			want: [][]int{{1, 2}, {2, 3}, {3, 4}},
		},
		{
			name: "width equal to length yields single window",
			args: args{input: []int{1, 2, 3}, size: 3},
			want: [][]int{{1, 2, 3}},
		},
		{
			name: "width larger than input yields empty output",
			args: args{input: []int{1, 2}, size: 3},
			want: [][]int{},
		},
		{
			name: "width of one yields singleton windows",
			args: args{input: []int{1, 2, 3}, size: 1},
			want: [][]int{{1}, {2}, {3}},
		},
		{
			name: "zero size yields empty output",
			args: args{input: []int{1, 2, 3}, size: 0},
			want: [][]int{},
		},
		{
			name: "negative size yields empty output",
			args: args{input: []int{1, 2, 3}, size: -1},
			want: [][]int{},
		},
		{
			name: "nil input yields non-nil empty output",
			args: args{input: nil, size: 2},
			want: [][]int{},
		},
		{
			name: "empty input yields non-nil empty output",
			args: args{input: []int{}, size: 2},
			want: [][]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Window(tt.args.input, tt.args.size)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Window() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestChunk_DoesNotMutateInputOnAppend asserts each chunk's capacity is clamped,
// so appending to one chunk cannot overwrite the next chunk or the input.
func TestChunk_DoesNotMutateInputOnAppend(t *testing.T) {
	input := []int{1, 2, 3, 4}
	chunks := slices.Chunk(input, 2)

	chunks[0] = append(chunks[0], 99)

	if !reflect.DeepEqual(input, []int{1, 2, 3, 4}) {
		t.Errorf("appending to a chunk mutated the input: %v", input)
	}
	if !reflect.DeepEqual(chunks[1], []int{3, 4}) {
		t.Errorf("appending to a chunk overwrote the next chunk: %v", chunks[1])
	}
}

// TestWindow_DoesNotMutateInputOnAppend asserts each window's capacity is
// clamped, so appending to one window cannot overwrite the overlapping data.
func TestWindow_DoesNotMutateInputOnAppend(t *testing.T) {
	input := []int{1, 2, 3, 4}
	windows := slices.Window(input, 2)

	windows[0] = append(windows[0], 99)

	if !reflect.DeepEqual(input, []int{1, 2, 3, 4}) {
		t.Errorf("appending to a window mutated the input: %v", input)
	}
	if !reflect.DeepEqual(windows[1], []int{2, 3}) {
		t.Errorf("appending to a window overwrote the next window: %v", windows[1])
	}
}
