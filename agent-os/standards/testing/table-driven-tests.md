# Table-Driven Tests

All `Test` functions are table-driven. Standard library only — no testify.

```go
func TestFilter_Strings(t *testing.T) {
	type args struct {
		input []string
		fun   slices.FilterFunc[string]
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "filters by length",
			args: args{input: []string{"a", "ab", "abc"}, fun: func(e string) bool { return len(e) > 1 }},
			want: []string{"ab", "abc"},
		},
		{name: "nil input results in empty output", args: args{input: nil}, want: []string{}},
		{name: "empty input results in empty output", args: args{input: []string{}}, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Filter(tt.args.input, tt.args.fun)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

- Group inputs in a local `type args struct`.
- Cases are `[]struct{ name string; args args; want ... }`.
- Run each case via `t.Run(tt.name, ...)`.
- Compare with `reflect.DeepEqual`; report with `t.Errorf("Func() = %v, want %v", got, tt.want)`.
- `name` reads as a sentence describing the scenario.
