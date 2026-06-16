package main

import (
	"bufio"
	"fmt"
	"io"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Sample is one standardized benchmark cell: a single (package, implementation,
// operation, size) measurement with its three metrics. It is the unit the
// report and chart are built from.
type Sample struct {
	Pkg      string // friendly package name, e.g. "dicts"
	Impl     string // implementation, e.g. "Hash" / "ConcurrentHashRW"
	Op       string // operation, e.g. "Get" / "PutInPlace"
	Size     int    // element count the benchmark ran against
	NsOp     float64
	BytesOp  float64
	AllocsOp float64
}

// Meta is the per-environment provenance: who ran the benchmarks, where, how,
// and how they should be presented. The commit/date/flags/label/tier come from
// the capture step (flags); goos/goarch/cpu/packages are discovered from the
// benchstat CSV config lines. It is persisted as a `# benchreport-meta:`
// preamble at the top of each committed dataset CSV so the dataset is
// self-describing and the render step needs no side files.
type Meta struct {
	Env       string // short id, e.g. "reference" / "ci"
	Label     string // display label, e.g. "Reference — Framework Desktop …"
	Tier      string // "primary" (headline/chart) or "secondary" (indicative)
	Machine   string // optional one-line machine description
	Commit    string
	Date      string // ISO-8601 UTC
	GoVersion string
	Benchtime string
	Count     string
	GOOS      string
	GOARCH    string
	CPU       string
	Packages  []string // full import paths, in first-seen order
}

// IsPrimary reports whether this dataset drives the headline table and chart.
func (m Meta) IsPrimary() bool { return m.Tier == tierPrimary }

// Dataset couples one environment's provenance with its samples.
type Dataset struct {
	Meta    Meta
	Samples []Sample
}

const (
	tierPrimary   = "primary"
	tierSecondary = "secondary"
	metaPrefix    = "# benchreport-meta:"
)

// benchstat emits benchmark names without the "Benchmark" prefix and with a
// trailing "-<GOMAXPROCS>". The standardized suite is "<Impl>_<Op>/size_<N>";
// anything that doesn't match (cross-impl Comparison/Integration benchmarks,
// or benchmarks without a size sub-benchmark) is intentionally excluded from
// the structured tables.
var nameRe = regexp.MustCompile(`^([A-Za-z][A-Za-z0-9]*)_([A-Za-z][A-Za-z0-9]*)/size_(\d+)(?:-\d+)?$`)

// isUnit reports whether a CSV column header names one of the three metrics.
func isUnit(s string) bool {
	switch s {
	case "sec/op", "B/op", "allocs/op":
		return true
	default:
		return false
	}
}

type cell struct {
	ns, bytes, allocs float64
}

// LoadDataset reads a committed dataset CSV — an optional `# benchreport-meta:`
// preamble followed by benchstat -format=csv output — into a Dataset. The
// provided base Meta (from capture-step flags) supplies fields the CSV can't
// carry; preamble lines, when present, override it; and the benchstat config
// lines fill in goos/goarch/cpu/packages. Returns the count of skipped
// non-conforming benchmark names for logging.
func LoadDataset(r io.Reader, base Meta) (Dataset, int, error) {
	meta := base
	cells := map[string]map[string]*cell{} // pkgPath -> name -> metrics
	var pkgOrder []string

	currentPkg := ""
	currentUnit := ""

	addPkg := func(p string) {
		if _, ok := cells[p]; ok {
			return
		}
		cells[p] = map[string]*cell{}
		pkgOrder = append(pkgOrder, p)
	}

	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r")
		if line == "" {
			currentUnit = ""
			continue
		}

		if rest, ok := strings.CutPrefix(line, metaPrefix); ok {
			applyMetaLine(&meta, strings.TrimSpace(rest))
			continue
		}

		if !strings.Contains(line, ",") {
			// Config line "key: value" (goos/goarch/pkg/cpu), or a stray token.
			key, val, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}
			key = strings.TrimSpace(key)
			val = strings.TrimSpace(val)
			switch key {
			case "goos":
				meta.GOOS = val
			case "goarch":
				meta.GOARCH = val
			case "cpu":
				meta.CPU = val
			case "pkg":
				currentPkg = val
				addPkg(val)
			}
			continue
		}

		fields := strings.Split(line, ",")
		if fields[0] == "" {
			// Header row: ",<unit>,CI" sets the unit; ",<file>," is ignored.
			if len(fields) > 1 && isUnit(fields[1]) {
				currentUnit = fields[1]
			}
			continue
		}

		name := fields[0]
		if name == "geomean" || currentUnit == "" || currentPkg == "" {
			continue
		}
		if len(fields) < 2 || fields[1] == "" {
			continue
		}
		v, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			continue // not a numeric data row
		}

		c := cells[currentPkg][name]
		if c == nil {
			c = &cell{}
			cells[currentPkg][name] = c
		}
		switch currentUnit {
		case "sec/op":
			c.ns = v * 1e9 // benchstat reports seconds; the report shows ns/op
		case "B/op":
			c.bytes = v
		case "allocs/op":
			c.allocs = v
		}
	}
	if err := sc.Err(); err != nil {
		return Dataset{}, 0, fmt.Errorf("reading dataset csv: %w", err)
	}

	meta.Packages = pkgOrder

	var samples []Sample
	skipped := 0
	for _, pkgPath := range pkgOrder {
		names := make([]string, 0, len(cells[pkgPath]))
		for n := range cells[pkgPath] {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			m := nameRe.FindStringSubmatch(n)
			if m == nil {
				skipped++
				continue
			}
			size, _ := strconv.Atoi(m[3])
			c := cells[pkgPath][n]
			samples = append(samples, Sample{
				Pkg:      path.Base(pkgPath),
				Impl:     m[1],
				Op:       m[2],
				Size:     size,
				NsOp:     c.ns,
				BytesOp:  c.bytes,
				AllocsOp: c.allocs,
			})
		}
	}

	return Dataset{Meta: meta, Samples: samples}, skipped, nil
}

// applyMetaLine parses one "key=value" provenance field from the preamble.
func applyMetaLine(m *Meta, kv string) {
	key, val, ok := strings.Cut(kv, "=")
	if !ok {
		return
	}
	switch strings.TrimSpace(key) {
	case "env":
		m.Env = val
	case "label":
		m.Label = val
	case "tier":
		m.Tier = val
	case "machine":
		m.Machine = val
	case "commit":
		m.Commit = val
	case "date":
		m.Date = val
	case "goversion":
		m.GoVersion = val
	case "benchtime":
		m.Benchtime = val
	case "count":
		m.Count = val
	}
}

// metaPreamble serializes a Meta into the deterministic `# benchreport-meta:`
// preamble written at the top of a captured dataset CSV. Only the
// capture-supplied fields are written; goos/goarch/cpu/packages stay in the
// benchstat body below it.
func metaPreamble(m Meta) string {
	var b strings.Builder
	write := func(k, v string) {
		if v != "" {
			fmt.Fprintf(&b, "%s %s=%s\n", metaPrefix, k, v)
		}
	}
	write("env", m.Env)
	write("label", m.Label)
	write("tier", m.Tier)
	write("machine", m.Machine)
	write("commit", m.Commit)
	write("date", m.Date)
	write("goversion", m.GoVersion)
	write("benchtime", m.Benchtime)
	write("count", m.Count)
	return b.String()
}
