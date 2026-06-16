package main

import (
	"fmt"
	"strings"
)

// Markers delimiting the auto-generated preview region in the README. Only the
// span between them is rewritten; everything else is preserved byte-for-byte.
const (
	MarkerStart = "<!-- BENCH:START -->"
	MarkerEnd   = "<!-- BENCH:END -->"
)

// InjectRegion replaces the content between MarkerStart and MarkerEnd with the
// rendered preview, leaving the markers in place and the rest of the document
// untouched. It is idempotent: running it on its own output is a no-op, because
// the canonical layout it writes is exactly what it looks for.
func InjectRegion(doc, region string) (string, error) {
	si := strings.Index(doc, MarkerStart)
	if si < 0 {
		return "", fmt.Errorf("start marker %q not found in README", MarkerStart)
	}
	ei := strings.Index(doc, MarkerEnd)
	if ei < 0 {
		return "", fmt.Errorf("end marker %q not found in README", MarkerEnd)
	}
	if ei < si {
		return "", fmt.Errorf("end marker appears before start marker in README")
	}

	before := doc[:si]
	after := doc[ei+len(MarkerEnd):]
	// region always ends in a newline, so MarkerEnd lands on its own line.
	return before + MarkerStart + "\n\n" + region + MarkerEnd + after, nil
}
