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

func loadSample(t *testing.T, base Meta) (Dataset, int) {
	t.Helper()
	ds, skipped, err := LoadDataset(strings.NewReader(sampleCSV), base)
	if err != nil {
		t.Fatalf("LoadDataset: %v", err)
	}
	return ds, skipped
}

func TestLoadDatasetConformingSamples(t *testing.T) {
	ds, skipped := loadSample(t, Meta{Benchtime: "50ms", Count: "8"})

	if skipped != 2 { // Comparison_Get/Hash and Tree_Get
		t.Errorf("skipped = %d, want 2", skipped)
	}
	if len(ds.Samples) != 3 { // 2 Hash_Get + 1 Hash_Contains
		t.Fatalf("len(samples) = %d, want 3", len(ds.Samples))
	}
	m := ds.Meta
	if m.GOOS != "linux" || m.GOARCH != "amd64" || m.CPU != "TestCPU" {
		t.Errorf("config = %+v, want linux/amd64/TestCPU", m)
	}
	if len(m.Packages) != 2 {
		t.Errorf("packages = %v, want 2 entries", m.Packages)
	}

	idx := indexSamples(ds.Samples)
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

func TestMetaPreambleRoundTrip(t *testing.T) {
	base := Meta{
		Env: "reference", Label: "Reference — Test Box", Tier: tierPrimary,
		Machine: "Test Box · 128 GB", Commit: "deadbee", Date: "2026-06-16T00:00:00Z",
		GoVersion: "go1.25.5", Benchtime: "50ms", Count: "8",
	}
	// A captured file is the preamble followed by the raw benchstat CSV; loading
	// it back must recover every capture-supplied field plus the CSV config.
	captured := metaPreamble(base) + sampleCSV
	ds, _, err := LoadDataset(strings.NewReader(captured), Meta{})
	if err != nil {
		t.Fatalf("LoadDataset: %v", err)
	}
	got := ds.Meta
	if got.Env != base.Env || got.Label != base.Label || got.Tier != base.Tier ||
		got.Machine != base.Machine || got.Commit != base.Commit || got.Date != base.Date ||
		got.GoVersion != base.GoVersion || got.Benchtime != base.Benchtime || got.Count != base.Count {
		t.Errorf("round-trip meta mismatch:\n got %+v\nwant %+v", got, base)
	}
	if !got.IsPrimary() {
		t.Error("expected primary tier to round-trip")
	}
	if got.GOOS != "linux" { // config still parsed from the CSV body
		t.Errorf("GOOS = %q, want linux", got.GOOS)
	}
}

func TestOrderDatasetsPrimaryFirst(t *testing.T) {
	ci := Dataset{Meta: Meta{Env: "ci", Tier: tierSecondary}}
	ref := Dataset{Meta: Meta{Env: "reference", Tier: tierPrimary}}
	ordered := orderDatasets([]Dataset{ci, ref})
	if !ordered[0].Meta.IsPrimary() || ordered[0].Meta.Env != "reference" {
		t.Errorf("primary not first: %+v", ordered)
	}
	if primaryDataset([]Dataset{ci, ref}).Meta.Env != "reference" {
		t.Error("primaryDataset did not pick the primary tier")
	}
}

func TestWithThousands(t *testing.T) {
	cases := map[string]string{
		"0": "0", "42": "42", "999": "999", "1000": "1,000", "12345": "12,345",
		"1234567": "1,234,567", "-12345": "-12,345", "1234.5": "1,234.5", "65767.123": "65,767.123",
	}
	for in, want := range cases {
		if got := withThousands(in); got != want {
			t.Errorf("withThousands(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFormatNsPrecision(t *testing.T) {
	cases := map[float64]string{0: "0", 5.246: "5.25", 40: "40.0", 906.6: "907", 65767.49: "65,767"}
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
	ds, _ := loadSample(t, Meta{})
	rows, missing := resolveHeadlines(indexSamples(ds.Samples))
	if len(rows) == 0 {
		t.Fatal("expected at least one resolved headline")
	}
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
	first := RenderChart(rows, "Env & <caption>")
	second := RenderChart(rows, "Env & <caption>")
	if first != second {
		t.Error("RenderChart not deterministic")
	}
	if !strings.Contains(first, "&amp; B &lt;test&gt;") || !strings.Contains(first, "Env &amp; &lt;caption&gt;") {
		t.Errorf("labels/caption not XML-escaped:\n%s", first)
	}
	if !strings.HasPrefix(first, "<svg") || !strings.HasSuffix(first, "</svg>\n") {
		t.Error("output is not a complete SVG document")
	}
}

func TestRenderReportTwoEnvironments(t *testing.T) {
	ref, _ := loadSample(t, Meta{Env: "reference", Label: "Reference — Test Box", Tier: tierPrimary,
		Commit: "deadbee", Date: "2026-06-16T00:00:00Z", Machine: "Test Box · 128 GB"})
	ci, _ := loadSample(t, Meta{Env: "ci", Label: "CI — shared runner", Tier: tierSecondary,
		Commit: "cafe123", Date: "2026-06-16T01:00:00Z"})

	out := RenderReport([]Dataset{ci, ref}, "docs/bench.svg") // pass out of order on purpose
	for _, want := range []string{
		"# Benchmark report",
		"## Environments",
		"### Reference — Test Box (primary)",
		"### CI — shared runner (secondary)",
		"Test Box · 128 GB",
		"indicative, not authoritative",
		"## Headlines — Reference — Test Box",
		"![Benchmark chart](docs/bench.svg)",
		"## Full results — Reference — Test Box",
		"## Full results — CI — shared runner (indicative)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("report missing %q", want)
		}
	}
	// Reference (primary) section must precede the CI (secondary) section.
	if strings.Index(out, "Full results — Reference") > strings.Index(out, "Full results — CI") {
		t.Error("primary environment did not lead the full-results sections")
	}
}

func TestRenderReadmeRegionContent(t *testing.T) {
	ref, _ := loadSample(t, Meta{Env: "reference", Label: "Reference — Box", Tier: tierPrimary,
		Commit: "cafe123", Date: "2026-06-16T09:30:00Z", GoVersion: "go1.25.5"})
	ci, _ := loadSample(t, Meta{Env: "ci", Label: "CI", Tier: tierSecondary,
		Commit: "0b87bdf", Date: "2026-06-16T10:00:00Z", GoVersion: "go1.24"})

	region := RenderReadmeRegion([]Dataset{ci, ref}, "docs/bench.svg", "BENCHMARKS.md")
	for _, want := range []string{
		"controlled **Reference — Box** baseline",
		"![Benchmark chart](docs/bench.svg)",
		"Reference — Box: `cafe123` · 2026-06-16 ·",
		"CI: `0b87bdf` · 2026-06-16 ·",
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
