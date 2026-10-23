package distance

// Hamming computes the Hamming distance between two equal-length sequences:
// the number of positions at which the corresponding elements differ. It is
// commonly used to measure the distance between bit strings, DNA sequences,
// or any fixed-width encoded items.
//
// It returns ok == false when the sequences have differing lengths — the
// Hamming distance is undefined for sequences that are not the same length.
// For empty equal-length inputs (both nil or both empty), it returns (0, true)
// — the distance between two empty sequences is zero.
func Hamming[T comparable](a, b []T) (int, bool) {
	if len(a) != len(b) {
		return 0, false
	}
	var count int
	for i := range a {
		if a[i] != b[i] {
			count++
		}
	}
	return count, true
}

// Levenshtein computes the Levenshtein edit distance between two sequences:
// the minimum number of single-element operations (insertion, deletion, or
// substitution) required to transform a into b. Elements may be of any
// comparable type (compared with ==), so it measures edits between token
// lists, DNA bases, byte strings or any other sequence; StringLevenshtein is
// the thin accessor for text.
//
// The distance is always defined: it is 0 for two empty (or nil) sequences,
// and len(b) when a is empty. Neither input is modified. It runs in
// O(len(a)·len(b)) time and keeps only two rows of the dynamic-programming
// matrix, so memory is O(min(len(a), len(b))).
//
// Examples:
//
//	Levenshtein([]int{1, 2, 3}, []int{1, 3})    → 1
//	Levenshtein([]string{"a", "b"}, []string{}) → 2
func Levenshtein[T comparable](a, b []T) int {
	// Keep b the shorter sequence so the DP rows are as short as possible.
	if len(a) < len(b) {
		a, b = b, a
	}
	if len(b) == 0 {
		return len(a)
	}

	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			del := prev[j] + 1
			ins := curr[j-1] + 1
			sub := prev[j-1] + cost
			curr[j] = min(del, ins, sub)
		}
		prev, curr = curr, prev
	}

	return prev[len(b)]
}

// StringLevenshtein computes the Levenshtein edit distance between two
// strings: the minimum number of single-character insertions, deletions or
// substitutions required to transform a into b.
//
// It operates over Unicode rune sequences, so multi-byte UTF-8 characters
// (e.g. accented letters, emoji) are treated as single units rather than
// individual bytes; it is exactly Levenshtein([]rune(a), []rune(b)). Use
// Levenshtein([]byte(a), []byte(b)) for a byte-level distance instead.
//
// Examples:
//
//	StringLevenshtein("kitten", "sitting") → 3
//	StringLevenshtein("", "abc")           → 3
//	StringLevenshtein("abc", "abc")        → 0
func StringLevenshtein(a, b string) int {
	return Levenshtein([]rune(a), []rune(b))
}
