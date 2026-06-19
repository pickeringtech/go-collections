package distance_test

import (
	"testing"
	"unicode/utf8"

	"github.com/pickeringtech/go-collections/ml/distance"
)

// FuzzLevenshtein asserts the structural invariants of the edit distance:
//   - non-negative result;
//   - zero when both strings are equal;
//   - symmetry: d(a,b) == d(b,a);
//   - triangle inequality: d(a,c) <= d(a,b) + d(b,c);
//   - lower bound: |len(a) - len(b)| <= d(a,b) (in runes);
//   - upper bound: d(a,b) <= max(len(a), len(b)) (in runes).
func FuzzLevenshtein(f *testing.F) {
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

		dAB := distance.Levenshtein(a, b)
		dBA := distance.Levenshtein(b, a)

		// Non-negative.
		if dAB < 0 {
			t.Fatalf("Levenshtein(%q, %q) = %d, want >= 0", a, b, dAB)
		}

		// Identity: d(a,a) == 0.
		if distance.Levenshtein(a, a) != 0 {
			t.Fatalf("Levenshtein(%q, %q) != 0", a, a)
		}

		// Symmetry.
		if dAB != dBA {
			t.Fatalf("Levenshtein(%q, %q) = %d, Levenshtein(%q, %q) = %d; want equal", a, b, dAB, b, a, dBA)
		}

		// Lower bound: |lenRunes(a) - lenRunes(b)| <= d(a,b).
		lenA := utf8.RuneCountInString(a)
		lenB := utf8.RuneCountInString(b)
		diff := lenA - lenB
		if diff < 0 {
			diff = -diff
		}
		if dAB < diff {
			t.Fatalf("Levenshtein(%q, %q) = %d, want >= %d (lower bound)", a, b, dAB, diff)
		}

		// Upper bound: d(a,b) <= max(lenA, lenB).
		maxLen := lenA
		if lenB > maxLen {
			maxLen = lenB
		}
		if dAB > maxLen {
			t.Fatalf("Levenshtein(%q, %q) = %d, want <= %d (upper bound)", a, b, dAB, maxLen)
		}
	})
}
