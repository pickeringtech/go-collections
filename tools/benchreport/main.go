// Command benchreport turns benchstat CSV output into the project's performance
// artifacts: the full BENCHMARKS.md report, a docs/bench.svg chart, and the
// marker-delimited preview region inside the README. It is a pure function of
// its inputs (CSV + provenance flags) so re-running with unchanged data and
// flags produces byte-identical output — see issue #50.
//
// Typical invocation (see the `bench-report` Make target):
//
//	go test ... -bench=. -benchmem ./collections/... > bench.txt
//	benchstat -format=csv bench.txt > bench.csv
//	go run . -csv bench.csv -commit "$(git rev-parse --short HEAD)" -date "$(date -u +%FT%TZ)" ...
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	if err := run(os.Args[1:], os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "benchreport:", err)
		os.Exit(1)
	}
}

func run(args []string, logw *os.File) error {
	fs := flag.NewFlagSet("benchreport", flag.ContinueOnError)
	var (
		csvPath    = fs.String("csv", "", "path to benchstat -format=csv output (\"-\" for stdin)")
		readmePath = fs.String("readme", "README.md", "README file to inject the preview region into")
		reportPath = fs.String("report", "BENCHMARKS.md", "output path for the full report")
		svgPath    = fs.String("svg", "docs/bench.svg", "output path for the chart (also the in-doc reference)")
		commit     = fs.String("commit", "", "commit SHA for the provenance stamp")
		date       = fs.String("date", "", "generation timestamp (UTC, ISO-8601) for the provenance stamp")
		goVersion  = fs.String("goversion", runtime.Version(), "Go version that ran the benchmarks")
		benchtime  = fs.String("benchtime", "", "the -benchtime value used (for the flags line)")
		count      = fs.String("count", "", "the -count value used (for the flags line)")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *csvPath == "" {
		return fmt.Errorf("-csv is required")
	}

	in := os.Stdin
	if *csvPath != "-" {
		f, err := os.Open(*csvPath)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	prov := Provenance{
		Commit:    *commit,
		Date:      *date,
		GoVersion: *goVersion,
		Benchtime: *benchtime,
		Count:     *count,
	}
	samples, prov, skipped, err := Parse(in, prov)
	if err != nil {
		return err
	}
	if len(samples) == 0 {
		return fmt.Errorf("no conforming benchmarks found in CSV (expected names like Hash_Get/size_1000)")
	}
	if skipped > 0 {
		fmt.Fprintf(logw, "note: skipped %d non-conforming benchmark name(s) (not <Impl>_<Op>/size_<N>)\n", skipped)
	}

	idx := indexSamples(samples)
	if _, missing := resolveHeadlines(idx); len(missing) > 0 {
		fmt.Fprintf(logw, "note: %d headline(s) absent from data, skipped: %v\n", len(missing), missing)
	}

	// Chart.
	rows, _ := resolveHeadlines(idx)
	if err := writeFile(*svgPath, RenderChart(rows)); err != nil {
		return err
	}

	// Full report.
	if err := writeFile(*reportPath, RenderReport(samples, prov, *svgPath)); err != nil {
		return err
	}

	// README preview region.
	readme, err := os.ReadFile(*readmePath)
	if err != nil {
		return err
	}
	region := RenderReadmeRegion(samples, prov, *svgPath, *reportPath)
	updated, err := InjectRegion(string(readme), region)
	if err != nil {
		return err
	}
	if err := writeFile(*readmePath, updated); err != nil {
		return err
	}

	fmt.Fprintf(logw, "wrote %s, %s, and the %s preview region (%d samples)\n",
		*reportPath, *svgPath, *readmePath, len(samples))
	return nil
}

// writeFile creates any missing parent directories, then writes the file.
func writeFile(path, content string) error {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
