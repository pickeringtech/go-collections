package distance_test

import (
	"testing"
	"unicode/utf8"

	"github.com/pickeringtech/go-collections/ml/distance"
)

// FuzzStringLevenshtein asserts the structural invariants of the edit distance:
//   - non-negative result;
//   - zero when both strings are equal;
//   - symmetry: d(a,b) == d(b,a);
//   - triangle inequality: d(a,c) <= d(a,b) + d(b,c);
//   - lower bound: |len(a) - len(b)| <= d(a,b) (in runes);
//   - upper bound: d(a,b) <= max(len(a), len(b)) (in runes).
func FuzzStringLevenshtein(f *testing.F) {
	f.Add("", "")
	f.Add("", "abc")
	f.Add("abc", "")
	f.Add("abc", "abc")
	f.Add("kitten", "sitting")
	f.Add("saturday", "sunday")
	f.Add("a", "b")
	f.Add("café", "cafe")

	f.Fuzz(func(t *testing.T, a, b string) {
		// Only test valid UTF-8 to keep the rune-count bound meaningful.
		if !utf8.ValidString(a) || !utf8.ValidString(b) {
			return
		}

		dAB := distance.StringLevenshtein(a, b)
		dBA := distance.StringLevenshtein(b, a)

		// Non-negative.
		if dAB < 0 {
			t.Fatalf("StringLevenshtein(%q, %q) = %d, want >= 0", a, b, dAB)
		}

		// Identity: d(a,a) == 0.
		if distance.StringLevenshtein(a, a) != 0 {
			t.Fatalf("StringLevenshtein(%q, %q) != 0", a, a)
		}

		// Symmetry.
		if dAB != dBA {
			t.Fatalf("StringLevenshtein(%q, %q) = %d, StringLevenshtein(%q, %q) = %d; want equal", a, b, dAB, b, a, dBA)
		}

		// Lower bound: |lenRunes(a) - lenRunes(b)| <= d(a,b).
		lenA := utf8.RuneCountInString(a)
		lenB := utf8.RuneCountInString(b)
		diff := lenA - lenB
		if diff < 0 {
			diff = -diff
		}
		if dAB < diff {
			t.Fatalf("StringLevenshtein(%q, %q) = %d, want >= %d (lower bound)", a, b, dAB, diff)
		}

		// Upper bound: d(a,b) <= max(lenA, lenB).
		maxLen := lenA
		if lenB > maxLen {
			maxLen = lenB
		}
		if dAB > maxLen {
			t.Fatalf("StringLevenshtein(%q, %q) = %d, want <= %d (upper bound)", a, b, dAB, maxLen)
		}
	})
}

// referenceLevenshtein is the textbook full-matrix edit distance, kept
// deliberately simple as a differential oracle for the two-row implementation.
func referenceLevenshtein(a, b []byte) int {
	d := make([][]int, len(a)+1)
	for i := range d {
		d[i] = make([]int, len(b)+1)
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			d[i][j] = min(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
		}
	}
	return d[len(a)][len(b)]
}

// FuzzLevenshtein checks the generic edit distance against the full-matrix
// reference on arbitrary byte sequences, plus symmetry and the agreement of
// StringLevenshtein with the generic distance over runes.
func FuzzLevenshtein(f *testing.F) {
	f.Add([]byte(""), []byte(""))
	f.Add([]byte("abc"), []byte(""))
	f.Add([]byte("kitten"), []byte("sitting"))
	f.Add([]byte("café"), []byte("cafe"))
	f.Add([]byte{0, 1, 2, 3}, []byte{3, 2, 1, 0})

	f.Fuzz(func(t *testing.T, a, b []byte) {
		if len(a) > 64 || len(b) > 64 {
			return // keep the quadratic reference cheap
		}
		got := distance.Levenshtein(a, b)
		want := referenceLevenshtein(a, b)
		if got != want {
			t.Fatalf("Levenshtein(%v, %v) = %d, want %d (reference)", a, b, got, want)
		}
		rev := distance.Levenshtein(b, a)
		if rev != got {
			t.Fatalf("Levenshtein not symmetric: %d vs %d", got, rev)
		}
		s := distance.StringLevenshtein(string(a), string(b))
		if s != distance.Levenshtein([]rune(string(a)), []rune(string(b))) {
			t.Fatalf("StringLevenshtein(%q, %q) = %d disagrees with Levenshtein over runes", a, b, s)
		}
	})
}
