package bloom

import "testing"

// TestNew_BitArrayIsWholeWords pins the documented memory footprint: the bit
// array is BitSize() bits rounded up to whole 64-bit words, and no larger. It
// inspects the array itself rather than measuring heap growth around New,
// which would also count whatever the runtime, the race detector or the test
// framework allocates concurrently.
func TestNew_BitArrayIsWholeWords(t *testing.T) {
	type args struct {
		n int
		p float64
	}
	tests := []struct {
		name      string
		args      args
		wantBits  uint64
		wantWords int
	}{
		// m = ceil(-1000·ln(0.01)/ln²2) = 9586 bits → 150 words.
		{name: "a thousand items at one percent needs 150 words", args: args{n: 1000, p: 0.01}, wantBits: 9586, wantWords: 150},
		// m = 2 bits still occupies one whole word.
		{name: "a two-bit filter needs one word", args: args{n: 1, p: 0.5}, wantBits: 2, wantWords: 1},
		// m = ceil(-ln(5.6e-14)/ln²2) = 64 bits fills one word exactly, so
		// rounding up must not add a spare word.
		{name: "a 64-bit filter fills exactly one word", args: args{n: 1, p: 5.6e-14}, wantBits: 64, wantWords: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New[int](tt.args.n, tt.args.p)
			if err != nil {
				t.Fatalf("New(%d, %v): unexpected error %v", tt.args.n, tt.args.p, err)
			}
			if f.BitSize() != tt.wantBits {
				t.Errorf("BitSize() = %d, want %d", f.BitSize(), tt.wantBits)
			}
			if len(f.bitsArr) != tt.wantWords || cap(f.bitsArr) != tt.wantWords {
				t.Errorf("bit array has len %d, cap %d; want %d whole words", len(f.bitsArr), cap(f.bitsArr), tt.wantWords)
			}
		})
	}
}
