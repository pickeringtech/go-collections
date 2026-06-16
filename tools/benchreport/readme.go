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
	// Search for the end marker only *after* the start marker, so a stray
	// MarkerEnd earlier in the document can't be mistaken for the region's end
	// (and an end marker that only appears before the start is correctly an
	// error, not silent corruption).
	afterStart := si + len(MarkerStart)
	rel := strings.Index(doc[afterStart:], MarkerEnd)
	if rel < 0 {
		return "", fmt.Errorf("end marker %q not found after start marker in README", MarkerEnd)
	}
	ei := afterStart + rel

	before := doc[:si]
	after := doc[ei+len(MarkerEnd):]
	// region always ends in a newline, so MarkerEnd lands on its own line.
	return before + MarkerStart + "\n\n" + region + MarkerEnd + after, nil
}
