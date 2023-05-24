package slices_test

import (
	"github.com/pickeringtech/go-collectionutil/slices"
	"reflect"
	"testing"
)

func TestAllMatch(t *testing.T) {
	type args struct {
		input []int
		fun   slices.FindFunc[int]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all elements pass test",
			args: args{
				input: []int{1, 2, 3, 4, 5},
				fun: func(element int) bool {
					return element < 6
				},
			},
			want: true,
		},
		{
			name: "some elements fail test",
			args: args{
				input: []int{1, 2, 3, 4, 5},
				fun: func(element int) bool {
					return element < 4
				},
			},
			want: false,
		},
		{
			name: "nil input results in true",
			args: args{
				input: nil,
				fun: func(element int) bool {
					return element < 4
				},
			},
			want: true,
		},
		{
			name: "empty input results in true",
			args: args{
				input: []int{},
				fun: func(element int) bool {
					return element < 4
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.AllMatch(tt.args.input, tt.args.fun)
			if got != tt.want {
				t.Errorf("AllMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyMatch(t *testing.T) {
	type args struct {
		input []int
		fun   slices.FindFunc[int]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "finds a match",
			args: args{
				input: []int{1, 3, 5, 7, 8},
				fun: func(element int) bool {
					return element%2 == 0
				},
			},
			want: true,
		},
		{
			name: "no match found",
			args: args{
				input: []int{1, 3, 5, 7},
				fun: func(element int) bool {
					return element%2 == 0
				},
			},
			want: false,
		},
		{
			name: "no match is found on nil input",
			args: args{
				input: nil,
				fun: func(element int) bool {
					return element%2 == 0
				},
			},
			want: false,
		},
		{
			name: "no match is found on empty input",
			args: args{
				input: []int{},
				fun: func(element int) bool {
					return element%2 == 0
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.AnyMatch(tt.args.input, tt.args.fun)
			if got != tt.want {
				t.Errorf("AnyMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFind(t *testing.T) {
	type args struct {
		input []int
		fun   slices.FindFunc[int]
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
		wantOk     bool
	}{
		{
			name: "selects the expected element",
			args: args{
				input: []int{2, 4, 6, 8, 10},
				fun: func(element int) bool {
					return element > 4
				},
			},
			wantResult: 6,
			wantOk:     true,
		},
		{
			name: "nil input results in zero value and boolean false",
			args: args{
				input: nil,
				fun: func(element int) bool {
					return element > 4
				},
			},
			wantResult: 0,
			wantOk:     false,
		},
		{
			name: "empty input results in zero value and boolean false",
			args: args{
				input: []int{},
				fun: func(element int) bool {
					return element > 4
				},
			},
			wantResult: 0,
			wantOk:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotOk := slices.Find(tt.args.input, tt.args.fun)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("FindAny() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotOk != tt.wantOk {
				t.Errorf("FindAny() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestFindIndex(t *testing.T) {
	type args[T any] struct {
		input []T
		fun   slices.FindFunc[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want int
	}
	tests := []testCase[int]{
		{
			name: "finds expected element index",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
				fun: func(a int) bool {
					return a > 2
				},
			},
			want: 2,
		},
		{
			name: "no match results in -1",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
				fun: func(a int) bool {
					return a > 10
				},
			},
			want: -1,
		},
		{
			name: "nil input results in -1",
			args: args[int]{
				input: nil,
				fun: func(a int) bool {
					return a > 10
				},
			},
			want: -1,
		},
		{
			name: "empty input results in -1",
			args: args[int]{
				input: nil,
				fun: func(a int) bool {
					return a > 10
				},
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.FindIndex(tt.args.input, tt.args.fun)
			if got != tt.want {
				t.Errorf("FindIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindLast(t *testing.T) {
	type args[T any] struct {
		input []T
		fun   slices.FindFunc[T]
	}
	type testCase[T any] struct {
		name       string
		args       args[T]
		wantResult T
		wantOk     bool
	}
	tests := []testCase[int]{
		{
			name: "finds the last entry matching the test function",
			args: args[int]{
				input: []int{5, 4, 3, 2, 1},
				fun: func(a int) bool {
					return a > 3
				},
			},
			wantResult: 4,
			wantOk:     true,
		},
		{
			name: "no match causes zero value and falsy boolean returns",
			args: args[int]{
				input: []int{5, 4, 3, 2, 1},
				fun: func(a int) bool {
					return a > 10
				},
			},
			wantResult: 0,
			wantOk:     false,
		},
		{
			name: "nil input causes zero value and falsy boolean returns",
			args: args[int]{
				input: nil,
				fun: func(a int) bool {
					return a > 10
				},
			},
			wantResult: 0,
			wantOk:     false,
		},
		{
			name: "empty input causes zero value and falsy boolean returns",
			args: args[int]{
				input: []int{},
				fun: func(a int) bool {
					return a > 10
				},
			},
			wantResult: 0,
			wantOk:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotOk := slices.FindLast(tt.args.input, tt.args.fun)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("FindLast() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotOk != tt.wantOk {
				t.Errorf("FindLast() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestFindLastIndex(t *testing.T) {
	type args[T any] struct {
		input []T
		fun   slices.FindFunc[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want int
	}
	tests := []testCase[int]{
		{
			name: "finds the last entry matching the test function",
			args: args[int]{
				input: []int{5, 4, 3, 2, 1},
				fun: func(a int) bool {
					return a > 3
				},
			},
			want: 1,
		},
		{
			name: "no match causes zero value and falsy boolean returns",
			args: args[int]{
				input: []int{5, 4, 3, 2, 1},
				fun: func(a int) bool {
					return a > 10
				},
			},
			want: -1,
		},
		{
			name: "nil input causes zero value and falsy boolean returns",
			args: args[int]{
				input: nil,
				fun: func(a int) bool {
					return a > 10
				},
			},
			want: -1,
		},
		{
			name: "empty input causes zero value and falsy boolean returns",
			args: args[int]{
				input: []int{},
				fun: func(a int) bool {
					return a > 10
				},
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.FindLastIndex(tt.args.input, tt.args.fun)
			if got != tt.want {
				t.Errorf("FindLastIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirst(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
		wantOk     bool
	}{
		{
			name: "finds the first element",
			args: args{
				input: []int{1, 2, 3, 4},
			},
			wantResult: 1,
			wantOk:     true,
		},
		{
			name: "returns zero value and false if nil input",
			args: args{
				input: nil,
			},
			wantResult: 0,
			wantOk:     false,
		},
		{
			name: "returns zero value and false if empty input",
			args: args{
				input: []int{},
			},
			wantResult: 0,
			wantOk:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotOk := slices.First(tt.args.input)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("FindFirst() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotOk != tt.wantOk {
				t.Errorf("FindFirst() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestIncludes(t *testing.T) {
	type args[T comparable] struct {
		input []T
		value T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "finds included value",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
				value: 3,
			},
			want: true,
		},
		{
			name: "does not find if value is not included",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
				value: 6,
			},
			want: false,
		},
		{
			name: "nil input provides falsy return",
			args: args[int]{
				input: nil,
				value: 6,
			},
			want: false,
		},
		{
			name: "empty input provides falsy return",
			args: args[int]{
				input: []int{},
				value: 6,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Includes(tt.args.input, tt.args.value)
			if got != tt.want {
				t.Errorf("Includes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexOf(t *testing.T) {
	type args[T comparable] struct {
		input []T
		value T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want int
	}
	tests := []testCase[int]{
		{
			name: "finds index of element in input",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
				value: 3,
			},
			want: 2,
		},
		{
			name: "not finding value results in -1",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
				value: 6,
			},
			want: -1,
		},
		{
			name: "nil input results in -1",
			args: args[int]{
				input: nil,
				value: 6,
			},
			want: -1,
		},
		{
			name: "empty input results in -1",
			args: args[int]{
				input: []int{},
				value: 6,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.IndexOf(tt.args.input, tt.args.value)
			if got != tt.want {
				t.Errorf("IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
