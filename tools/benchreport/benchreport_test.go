package main

import (
	"strings"
	"testing"
)

// sampleCSV is a trimmed but structurally faithful benchstat -format=csv output:
// two packages, three unit tables each, a geomean row, a cross-impl benchmark
// without a size sub-benchmark (Comparison_Get/Hash) that must be skipped, and a
// Tree benchmark with no size sub-benchmark (Tree_Get) that must also be skipped.
const sampleCSV = `goos: linux
goarch: amd64
pkg: github.com/pickeringtech/go-collections/collections/dicts
cpu: TestCPU
,bench.txt,
,sec/op,CI
Hash_Get/size_10-32,2.5e-08,∞
Hash_Get/size_1000-32,4e-08,∞
Comparison_Get/Hash-32,5e-08,∞
Tree_Get-32,2e-08,∞
geomean,3e-08,

,bench.txt,
,B/op,CI
Hash_Get/size_10-32,0,∞
Hash_Get/size_1000-32,0,∞
geomean,,

,bench.txt,
,allocs/op,CI
Hash_Get/size_10-32,0,∞
Hash_Get/size_1000-32,0,∞
geomean,,

pkg: github.com/pickeringtech/go-collections/collections/sets
,bench.txt,
,sec/op,CI
Hash_Contains/size_1000-32,2.6e-08,∞
geomean,2.6e-08,

,bench.txt,
,B/op,CI
Hash_Contains/size_1000-32,8,∞
geomean,,

,bench.txt,
,allocs/op,CI
Hash_Contains/size_1000-32,1,∞
geomean,,
`

func parseSample(t *testing.T) ([]Sample, Provenance, int) {
	t.Helper()
	samples, prov, skipped, err := Parse(strings.NewReader(sampleCSV), Provenance{Benchtime: "50ms", Count: "8"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return samples, prov, skipped
}

func TestParseConformingSamples(t *testing.T) {
	samples, prov, skipped := parseSample(t)

	if skipped != 2 { // Comparison_Get/Hash and Tree_Get
		t.Errorf("skipped = %d, want 2", skipped)
	}
	if len(samples) != 3 { // 2 Hash_Get + 1 Hash_Contains
		t.Fatalf("len(samples) = %d, want 3", len(samples))
	}
	if prov.GOOS != "linux" || prov.GOARCH != "amd64" || prov.CPU != "TestCPU" {
		t.Errorf("provenance config = %+v, want linux/amd64/TestCPU", prov)
	}
	if len(prov.Packages) != 2 {
		t.Errorf("packages = %v, want 2 entries", prov.Packages)
	}

	idx := indexSamples(samples)
	get1k, ok := idx[sampleKey{"dicts", "Hash", "Get", 1000}]
	if !ok {
		t.Fatal("missing dicts Hash Get size 1000")
	}
	if get1k.NsOp != 40 { // 4e-08 s -> 40 ns
		t.Errorf("Hash_Get/1000 ns/op = %v, want 40", get1k.NsOp)
	}
	contains, ok := idx[sampleKey{"sets", "Hash", "Contains", 1000}]
	if !ok {
		t.Fatal("missing sets Hash Contains size 1000")
	}
	if contains.BytesOp != 8 || contains.AllocsOp != 1 {
		t.Errorf("Contains B/op=%v allocs/op=%v, want 8 and 1", contains.BytesOp, contains.AllocsOp)
	}
}

func TestParseRejectsNonConforming(t *testing.T) {
	samples, _, _ := parseSample(t)
	for _, s := range samples {
		if s.Op == "" || s.Impl == "" || s.Size == 0 {
			t.Errorf("non-conforming sample leaked through: %+v", s)
		}
	}
}

func TestWithThousands(t *testing.T) {
	cases := map[string]string{
		"0":         "0",
		"42":        "42",
		"999":       "999",
		"1000":      "1,000",
		"12345":     "12,345",
		"1234567":   "1,234,567",
		"-12345":    "-12,345",
		"1234.5":    "1,234.5",
		"65767.123": "65,767.123",
	}
	for in, want := range cases {
		if got := withThousands(in); got != want {
			t.Errorf("withThousands(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFormatNsPrecision(t *testing.T) {
	cases := map[float64]string{
		0:        "0",
		5.246:    "5.25",
		40:       "40.0",
		906.6:    "907",
		65767.49: "65,767",
	}
	for in, want := range cases {
		if got := formatNs(in); got != want {
			t.Errorf("formatNs(%v) = %q, want %q", in, got, want)
		}
	}
}

func TestFormatCount(t *testing.T) {
	cases := map[float64]string{0: "0", 8: "8", 3208: "3,208", 7.5: "7.5"}
	for in, want := range cases {
		if got := formatCount(in); got != want {
			t.Errorf("formatCount(%v) = %q, want %q", in, got, want)
		}
	}
}

func TestInjectRegionIdempotent(t *testing.T) {
	doc := "intro\n" + MarkerStart + "\n\nold content\n" + MarkerEnd + "\noutro\n"
	region := "new content line\n"

	once, err := InjectRegion(doc, region)
	if err != nil {
		t.Fatalf("InjectRegion: %v", err)
	}
	if !strings.Contains(once, "new content line") || strings.Contains(once, "old content") {
		t.Errorf("region not replaced:\n%s", once)
	}
	if !strings.HasPrefix(once, "intro\n") || !strings.HasSuffix(once, "outro\n") {
		t.Errorf("surrounding text not preserved:\n%s", once)
	}

	twice, err := InjectRegion(once, region)
	if err != nil {
		t.Fatalf("InjectRegion (2nd): %v", err)
	}
	if once != twice {
		t.Errorf("InjectRegion not idempotent:\n--- once ---\n%s\n--- twice ---\n%s", once, twice)
	}
}

func TestInjectRegionMissingMarkers(t *testing.T) {
	if _, err := InjectRegion("no markers here", "x\n"); err == nil {
		t.Error("expected error when start marker missing")
	}
	if _, err := InjectRegion(MarkerStart+"\nbut no end", "x\n"); err == nil {
		t.Error("expected error when end marker missing")
	}
	if _, err := InjectRegion(MarkerEnd+" before "+MarkerStart, "x\n"); err == nil {
		t.Error("expected error when end precedes start")
	}
}

func TestResolveHeadlinesReportsMissing(t *testing.T) {
	samples, _, _ := parseSample(t)
	idx := indexSamples(samples)
	rows, missing := resolveHeadlines(idx)
	if len(rows) == 0 {
		t.Fatal("expected at least one resolved headline")
	}
	// The sample data lacks the ConcurrentHash/ConcurrentHashRW/Tree/Array
	// headlines, so they must be reported missing rather than rendered blank.
	if len(missing) == 0 {
		t.Error("expected some headlines to be reported missing")
	}
	for _, r := range rows {
		if r.Label == "" {
			t.Error("resolved headline has empty label")
		}
	}
}

func TestRenderChartDeterministicAndSafe(t *testing.T) {
	rows := []headlineRow{
		{"A & B <test>", 40, 0, 0},
		{"C", 0, 0, 0}, // zero value must not divide-by-zero or vanish
	}
	first := RenderChart(rows)
	second := RenderChart(rows)
	if first != second {
		t.Error("RenderChart not deterministic")
	}
	if !strings.Contains(first, "&amp; B &lt;test&gt;") {
		t.Errorf("labels not XML-escaped:\n%s", first)
	}
	if !strings.HasPrefix(first, "<svg") || !strings.HasSuffix(first, "</svg>\n") {
		t.Error("output is not a complete SVG document")
	}
}

func TestRenderReportStructure(t *testing.T) {
	samples, prov, _ := parseSample(t)
	prov.Commit = "deadbee"
	prov.Date = "2026-06-16T00:00:00Z"
	out := RenderReport(samples, prov, "docs/bench.svg")

	for _, want := range []string{
		"# Benchmark report",
		"## Provenance",
		"`deadbee`",
		"indicative, not authoritative",
		"![Benchmark chart](docs/bench.svg)",
		"## Full results",
		"### dicts",
		"### sets",
		"#### Get",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("report missing %q", want)
		}
	}
}

func TestRenderReadmeRegionContent(t *testing.T) {
	samples, prov, _ := parseSample(t)
	prov.Commit = "cafe123"
	prov.Date = "2026-06-16T09:30:00Z"
	region := RenderReadmeRegion(samples, prov, "docs/bench.svg", "BENCHMARKS.md")

	for _, want := range []string{
		"Indicative numbers",
		"![Benchmark chart](docs/bench.svg)",
		"Provenance: `cafe123` · 2026-06-16 ·",
		"Full report → [BENCHMARKS.md](BENCHMARKS.md)",
	} {
		if !strings.Contains(region, want) {
			t.Errorf("README region missing %q\n%s", want, region)
		}
	}
	if !strings.HasSuffix(region, "\n") {
		t.Error("region must end with a newline so the end marker lands on its own line")
	}
}
