package slices

import (
	"reflect"
	"testing"
)

func TestFindFirst(t *testing.T) {
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
			gotResult, gotOk := FindFirst(tt.args.input)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("FindFirst() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotOk != tt.wantOk {
				t.Errorf("FindFirst() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestFindAny(t *testing.T) {
	type args struct {
		input []int
		fun   FindFunc[int]
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
			gotResult, gotOk := FindAny(tt.args.input, tt.args.fun)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("FindAny() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotOk != tt.wantOk {
				t.Errorf("FindAny() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestAnyMatch(t *testing.T) {
	type args struct {
		input []int
		fun   FindFunc[int]
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
			got := AnyMatch(tt.args.input, tt.args.fun)
			if got != tt.want {
				t.Errorf("AnyMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllMatch(t *testing.T) {
	type args struct {
		input []int
		fun   FindFunc[int]
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
			if got := AllMatch(tt.args.input, tt.args.fun); got != tt.want {
				t.Errorf("AllMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
