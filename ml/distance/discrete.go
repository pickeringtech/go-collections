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

// Levenshtein computes the Levenshtein edit distance between two strings:
// the minimum number of single-character operations (insertion, deletion, or
// substitution) required to transform a into b.
//
// The implementation operates over Unicode rune sequences, so multi-byte UTF-8
// characters (e.g. accented letters, emoji) are treated as single units rather
// than individual bytes.
//
// Examples:
//
//	Levenshtein("kitten", "sitting") → 3
//	Levenshtein("", "abc")           → 3
//	Levenshtein("abc", "abc")        → 0
func Levenshtein(a, b string) int {
	ra := []rune(a)
	rb := []rune(b)
	lenA := len(ra)
	lenB := len(rb)

	if lenA == 0 {
		return lenB
	}
	if lenB == 0 {
		return lenA
	}

	// Use two rows of the DP matrix to keep memory O(min(lenA, lenB)).
	// We ensure b is the shorter string so the row length is minimal.
	if lenA < lenB {
		ra, rb = rb, ra
		lenA, lenB = lenB, lenA
	}

	prev := make([]int, lenB+1)
	curr := make([]int, lenB+1)

	for j := 0; j <= lenB; j++ {
		prev[j] = j
	}

	for i := 1; i <= lenA; i++ {
		curr[0] = i
		for j := 1; j <= lenB; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			del := prev[j] + 1
			ins := curr[j-1] + 1
			sub := prev[j-1] + cost
			curr[j] = min(del, ins, sub)
		}
		prev, curr = curr, prev
	}

	return prev[lenB]
}
