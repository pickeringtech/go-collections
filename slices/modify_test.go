package slices_test

import (
	"github.com/pickeringtech/go-collectionutil/slices"
	"reflect"
	"testing"
)

func TestConcatenate(t *testing.T) {
	type args struct {
		inputA []int
		inputB []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "inputs are joined consecutively",
			args: args{
				inputA: []int{1, 2, 3},
				inputB: []int{4, 5, 6},
			},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "nil first input means result is second input",
			args: args{
				inputA: nil,
				inputB: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "empty first input means result is second input",
			args: args{
				inputA: []int{},
				inputB: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "nil second input means result is first input",
			args: args{
				inputA: []int{1, 2, 3},
				inputB: nil,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "empty second input means result is first input",
			args: args{
				inputA: []int{1, 2, 3},
				inputB: []int{},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Concatenate(tt.args.inputA, tt.args.inputB)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Concatenate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "duplicates the input into a new slice",
			args: args{
				input: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "nil input provides nil output",
			args: args{
				input: nil,
			},
			want: nil,
		},
		{
			name: "empty input provides nil output",
			args: args{
				input: []int{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Copy(tt.args.input)
			tt.args.input = append(tt.args.input, 45)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		input []int
		index int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "removes the element at the specified index",
			args: args{
				input: []int{1, 2, 3, 4},
				index: 2,
			},
			want: []int{1, 2, 4},
		},
		{
			name: "removes the element at the last index",
			args: args{
				input: []int{1, 2, 3, 4},
				index: 3,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "removes the zeroth element",
			args: args{
				input: []int{1, 2, 3, 4},
				index: 0,
			},
			want: []int{2, 3, 4},
		},
		{
			name: "if index is beyond range the slice is not modified",
			args: args{
				input: []int{1, 2, 3, 4},
				index: 4,
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "if index is below zero the slice is not modified",
			args: args{
				input: []int{1, 2, 3, 4},
				index: -1,
			},
			want: []int{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Delete(tt.args.input, tt.args.index)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFill(t *testing.T) {
	type args[T any] struct {
		input []T
		value T
	}
	type testCase[T any] struct {
		name                   string
		args                   args[T]
		want                   []T
		ensureInputIsUnchanged bool
	}
	tests := []testCase[any]{
		{
			name: "fills every value",
			args: args[any]{
				input: []any{1, 2, 3, 4},
				value: 10,
			},
			want:                   []any{10, 10, 10, 10},
			ensureInputIsUnchanged: true,
		},
		{
			name: "nil input causes nil output",
			args: args[any]{
				input: nil,
				value: 10,
			},
			want:                   nil,
			ensureInputIsUnchanged: true,
		},
		{
			name: "empty input causes nil output",
			args: args[any]{
				input: []any{},
				value: 10,
			},
			want:                   nil,
			ensureInputIsUnchanged: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalInput := slices.Copy(tt.args.input)
			got := slices.Fill(tt.args.input, tt.args.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fill() = %v, want %v", got, tt.want)
			}
			if tt.ensureInputIsUnchanged && !reflect.DeepEqual(tt.args.input, originalInput) {
				t.Errorf("Fill() modified original input - original %v, modified input %v", originalInput, tt.args.input)
			}
		})
	}
}

func TestFillFrom(t *testing.T) {
	type args[T any] struct {
		input     []T
		value     T
		fromIndex int
	}
	type testCase[T any] struct {
		name                   string
		args                   args[T]
		want                   []T
		ensureInputIsUnchanged bool
	}
	tests := []testCase[any]{
		{
			name: "fills from a given index",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: 2,
			},
			want: []any{1, 2, 10, 10, 10},
		},
		{
			name: "if from index is beyond slice length it is unchanged",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: 5,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "if from index is below zero the slice is unchanged",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: -1,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "nil input results in nil output",
			args: args[any]{
				input:     nil,
				value:     10,
				fromIndex: -1,
			},
			want: nil,
		},
		{
			name: "empty input results in nil output",
			args: args[any]{
				input:     []any{},
				value:     10,
				fromIndex: -1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalInput := slices.Copy(tt.args.input)
			got := slices.FillFrom(tt.args.input, tt.args.value, tt.args.fromIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FillFrom() = %v, want %v", got, tt.want)
			}
			if tt.ensureInputIsUnchanged && !reflect.DeepEqual(tt.args.input, originalInput) {
				t.Errorf("FillFrom() modified original input - original %v, modified input %v", originalInput, tt.args.input)
			}
		})
	}
}

func TestFillFromTo(t *testing.T) {
	type args[T any] struct {
		input     []T
		value     T
		fromIndex int
		toIndex   int
	}
	type testCase[T any] struct {
		name                   string
		args                   args[T]
		want                   []T
		ensureInputIsUnchanged bool
	}
	tests := []testCase[any]{
		{
			name: "fills a range within a slice with a specified value",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: 2,
				toIndex:   4,
			},
			want: []any{1, 2, 10, 10, 5},
		},
		{
			name: "from index larger than to index causes no change",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: 4,
				toIndex:   2,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "negative from index causes no change",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: -1,
				toIndex:   2,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "to index beyond length of input causes no change",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: 0,
				toIndex:   6,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "nil input results in nil output",
			args: args[any]{
				input:     nil,
				value:     10,
				fromIndex: 0,
				toIndex:   6,
			},
			want: nil,
		},
		{
			name: "empty input results in nil output",
			args: args[any]{
				input:     []any{},
				value:     10,
				fromIndex: 0,
				toIndex:   6,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalInput := slices.Copy(tt.args.input)
			got := slices.FillFromTo(tt.args.input, tt.args.value, tt.args.fromIndex, tt.args.toIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FillFromTo() = %v, want %v", got, tt.want)
			}
			if tt.ensureInputIsUnchanged && !reflect.DeepEqual(tt.args.input, originalInput) {
				t.Errorf("FillFromTo() modified original input - original %v, modified input %v", originalInput, tt.args.input)
			}
		})
	}
}

func TestFillTo(t *testing.T) {
	type args[T any] struct {
		input   []T
		value   T
		toIndex int
	}
	type testCase[T any] struct {
		name                   string
		args                   args[T]
		want                   []T
		ensureInputIsUnchanged bool
	}
	tests := []testCase[any]{
		{
			name: "fills a range within a slice with a specified value",
			args: args[any]{
				input:   []any{1, 2, 3, 4, 5},
				value:   10,
				toIndex: 4,
			},
			want: []any{10, 10, 10, 10, 5},
		},
		{
			name: "negative to index causes no change",
			args: args[any]{
				input:   []any{1, 2, 3, 4, 5},
				value:   10,
				toIndex: -1,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "to index beyond length of input causes no change",
			args: args[any]{
				input:   []any{1, 2, 3, 4, 5},
				value:   10,
				toIndex: 6,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "nil input results in nil output",
			args: args[any]{
				input:   nil,
				value:   10,
				toIndex: 6,
			},
			want: nil,
		},
		{
			name: "empty input results in nil output",
			args: args[any]{
				input:   []any{},
				value:   10,
				toIndex: 6,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalInput := slices.Copy(tt.args.input)
			got := slices.FillTo(tt.args.input, tt.args.value, tt.args.toIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FillTo() = %v, want %v", got, tt.want)
			}
			if tt.ensureInputIsUnchanged && !reflect.DeepEqual(tt.args.input, originalInput) {
				t.Errorf("FillFromTo() modified original input - original %v, modified input %v", originalInput, tt.args.input)
			}
		})
	}
}

func TestJoinToString(t *testing.T) {
	type args[T any] struct {
		input     []T
		separator string
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want string
	}
	tests := []testCase[any]{
		{
			name: "joins correctly using separator",
			args: args[any]{
				input:     []any{"Earth", "Wind", "Fire", "Water"},
				separator: "-and-",
			},
			want: "Earth-and-Wind-and-Fire-and-Water",
		},
		{
			name: "joins correctly using separator and varying types",
			args: args[any]{
				input:     []any{"Earth", 10, "Fire", []string{"a", "b", "c"}},
				separator: "-and-",
			},
			want: "Earth-and-10-and-Fire-and-[a b c]",
		},
		{
			name: "nil input results in empty output",
			args: args[any]{
				input:     nil,
				separator: "-and-",
			},
			want: "",
		},
		{
			name: "empty input results in empty output",
			args: args[any]{
				input:     []any{},
				separator: "-and-",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.JoinToString(tt.args.input, tt.args.separator)
			if got != tt.want {
				t.Errorf("JoinToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPop(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantSli []int
	}{
		{
			name: "gets the last element from the input returning the smaller slice",
			args: args{
				input: []int{1, 2, 3, 4},
			},
			want:    4,
			wantSli: []int{1, 2, 3},
		},
		{
			name: "nil input provides zero value for type and nil resulting slice",
			args: args{
				input: nil,
			},
			want:    0,
			wantSli: nil,
		},
		{
			name: "empty input provides zero value for type and nil resulting slice",
			args: args{
				input: []int{},
			},
			want:    0,
			wantSli: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotSli := slices.Pop(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pop() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(gotSli, tt.wantSli) {
				t.Errorf("Pop() gotSli = %v, want %v", gotSli, tt.wantSli)
			}
		})
	}
}

func TestPopFront(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name             string
		args             args
		wantFirstElement int
		wantNewSlice     []int
	}{
		{
			name: "first element is removed from input slice and returned",
			args: args{
				input: []int{5, 4, 3, 2, 1},
			},
			wantFirstElement: 5,
			wantNewSlice:     []int{4, 3, 2, 1},
		},
		{
			name: "nil input provides zero value output and nil resulting slice",
			args: args{
				input: nil,
			},
			wantFirstElement: 0,
			wantNewSlice:     nil,
		},
		{
			name: "empty input provides zero value output and nil resulting slice",
			args: args{
				input: []int{},
			},
			wantFirstElement: 0,
			wantNewSlice:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFirstElement, gotNewSlice := slices.PopFront(tt.args.input)
			if !reflect.DeepEqual(gotFirstElement, tt.wantFirstElement) {
				t.Errorf("PopFront() gotFirstElement = %v, want %v", gotFirstElement, tt.wantFirstElement)
			}
			if !reflect.DeepEqual(gotNewSlice, tt.wantNewSlice) {
				t.Errorf("PopFront() gotNewSlice = %v, want %v", gotNewSlice, tt.wantNewSlice)
			}
		})
	}
}

func TestPush(t *testing.T) {
	type args struct {
		input       []int
		newElements []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "pushes the new elements to the end of the input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: []int{4, 5, 6},
			},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "nil input slice results in only the new elements",
			args: args{
				input:       nil,
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "empty input slice results in only the new elements",
			args: args{
				input:       []int{},
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "nil new elements results in only original input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: nil,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "empty new elements results in only original input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: []int{},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Push(tt.args.input, tt.args.newElements...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Push() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPushFront(t *testing.T) {
	type args struct {
		input       []int
		newElements []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "adds the new elements to the front of the input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6, 1, 2, 3},
		},
		{
			name: "nil new elements results in original input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: nil,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "empty new elements results in original input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: []int{},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "nil input slice results in only new elements",
			args: args{
				input:       nil,
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "empty input slice results in only new elements",
			args: args{
				input:       []int{},
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.PushFront(tt.args.input, tt.args.newElements...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PushFront() = %v, want %v", got, tt.want)
			}
		})
	}
}
