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

// Provenance is the trust header for the report: exactly what was run, where,
// and how. The runner/CPU/package facts come from benchstat's CSV config lines;
// the commit/date/flags are supplied by the caller (Make target or CI job).
type Provenance struct {
	Commit    string
	Date      string // ISO-8601 UTC, supplied by the caller
	GoVersion string
	GOOS      string
	GOARCH    string
	CPU       string
	Benchtime string
	Count     string
	Packages  []string // full import paths, in first-seen order
}

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

// Parse reads benchstat -format=csv output and returns the conforming samples
// plus the provenance facts discoverable from the CSV. Non-conforming benchmark
// names are skipped; the count of skipped names is returned for logging.
//
// The CSV is a sequence of per-config blocks (one per package, since the config
// key includes pkg). Each block holds three unit tables (sec/op, B/op,
// allocs/op). Config lines ("goos: …", "pkg: …") carry no comma; table header
// rows start with a comma; data rows are "<name>,<value>,<CI>".
func Parse(r io.Reader, prov Provenance) ([]Sample, Provenance, int, error) {
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
				prov.GOOS = val
			case "goarch":
				prov.GOARCH = val
			case "cpu":
				prov.CPU = val
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
		return nil, prov, 0, fmt.Errorf("reading benchstat csv: %w", err)
	}

	prov.Packages = pkgOrder

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

	return samples, prov, skipped, nil
}
