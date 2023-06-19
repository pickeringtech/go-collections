package slices_test

import (
	"fmt"
	"github.com/pickeringtech/go-collectionutil/slices"
	"reflect"
	"strconv"
	"testing"
)

func ExampleGenerate() {
	type numberAndSquare struct {
		number int
		square int
	}

	sli := slices.Generate[numberAndSquare](10, func(index int) numberAndSquare {
		return numberAndSquare{
			number: index,
			square: index * index,
		}
	})

	fmt.Printf("%v", sli)
	// Output: [{0 0} {1 1} {2 4} {3 9} {4 16} {5 25} {6 36} {7 49} {8 64} {9 81}]
}

func TestGenerate(t *testing.T) {
	type args[T any] struct {
		amount int
		fn     slices.GeneratorFunc[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[string]{
		{
			name: "generates correctly",
			args: args[string]{
				amount: 3,
				fn: func(index int) string {
					return strconv.Itoa(index * 3)
				},
			},
			want: []string{"0", "3", "6"},
		},
		{
			name: "amount 0 provides nil output",
			args: args[string]{
				amount: 0,
				fn: func(index int) string {
					return strconv.Itoa(index * 3)
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Generate(tt.args.amount, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
