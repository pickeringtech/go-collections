// Command word-frequency reads text from stdin (or from files named as
// arguments), tokenises it, counts how often each word occurs, and prints the
// most frequent words.
//
// It demonstrates a realistic slices + maps flow:
//   - slices.Filter trims the token stream down to interesting words;
//   - a map tallies occurrences, and maps.Items pulls the tallies out;
//   - slices.Sort ranks them, with a tie-breaker that keeps the output
//     deterministic even though map iteration order is not.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/pickeringtech/go-collections/maps"
	"github.com/pickeringtech/go-collections/slices"
)

func main() {
	topN := flag.Int("n", 5, "how many of the most frequent words to print")
	minLen := flag.Int("min", 1, "ignore words shorter than this many runes")
	flag.Parse()

	text, err := readInput(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, "word-frequency:", err)
		os.Exit(1)
	}

	for _, line := range rank(text, *topN, *minLen) {
		fmt.Println(line)
	}
}

// readInput returns the concatenation of every named file, or stdin when no
// files are given.
func readInput(paths []string) (string, error) {
	if len(paths) == 0 {
		data, err := io.ReadAll(os.Stdin)
		return string(data), err
	}
	var b strings.Builder
	for _, p := range paths {
		data, err := os.ReadFile(p) //nolint:gosec // example reads user-named files by design
		if err != nil {
			return "", err
		}
		b.Write(data)
		b.WriteByte('\n')
	}
	return b.String(), nil
}

// rank tokenises text and returns up to topN "word: count" lines, ordered by
// descending count and then alphabetically so ties resolve deterministically.
func rank(text string, topN, minLen int) []string {
	tokens := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	// slices: keep only the words long enough to be worth counting.
	words := slices.Filter(tokens, func(w string) bool {
		return len([]rune(w)) >= minLen
	})

	// map: tally occurrences, then lift the entries out for ranking.
	counts := map[string]int{}
	for _, w := range words {
		counts[w]++
	}
	entries := maps.Items(counts)

	// slices: rank by count desc, then word asc. The alphabetical tie-break is
	// what makes the output reproducible despite random map iteration order.
	ranked := slices.Sort(entries, func(a, b maps.Entry[string, int]) bool {
		if a.Value != b.Value {
			return a.Value > b.Value
		}
		return a.Key < b.Key
	})

	end := min(topN, len(ranked))
	out := make([]string, 0, end)
	for _, e := range slices.SubSlice(ranked, 0, end) {
		out = append(out, fmt.Sprintf("%s: %d", e.Key, e.Value))
	}
	return out
}
