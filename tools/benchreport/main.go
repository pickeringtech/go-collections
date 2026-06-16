// Command benchreport turns benchstat CSV output into the project's performance
// artifacts. It works in two stages so two environments — a controlled
// reference machine and the noisy CI runner — can each refresh independently
// while sharing one combined report (see issue #50):
//
//	capture: wrap one run's benchstat CSV with a provenance preamble and write
//	         it to docs/bench/<env>.csv (the committed per-environment dataset).
//	render:  read every docs/bench/*.csv and emit the combined BENCHMARKS.md,
//	         docs/bench.svg chart, and README preview region.
//
// Each stage is a pure function of its inputs, so re-running with unchanged data
// and flags produces byte-identical output. See the `bench-report` Make target
// for the full local pipeline; the main-only CI job runs the same target with
// the CI environment overrides.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "benchreport: expected a subcommand: capture | render")
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "capture":
		err = runCapture(os.Args[2:])
	case "render":
		err = runRender(os.Args[2:])
	default:
		err = fmt.Errorf("unknown subcommand %q (want: capture | render)", os.Args[1])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "benchreport:", err)
		os.Exit(1)
	}
}

// runCapture wraps one run's benchstat CSV with a provenance preamble and writes
// the committed per-environment dataset. It validates that the CSV holds at
// least one conforming benchmark before writing.
func runCapture(args []string) error {
	fs := flag.NewFlagSet("capture", flag.ContinueOnError)
	var (
		in     = fs.String("in", "-", "raw benchstat -format=csv input (\"-\" for stdin)")
		out    = fs.String("out", "", "output dataset path, e.g. docs/bench/reference.csv")
		env    = fs.String("env", "", "environment id, e.g. reference / ci")
		label  = fs.String("label", "", "display label for the environment")
		tier   = fs.String("tier", tierSecondary, "primary (headline/chart) or secondary (indicative)")
		mach   = fs.String("machine", "", "optional one-line machine description")
		commit = fs.String("commit", "", "commit SHA for the provenance stamp")
		date   = fs.String("date", "", "generation timestamp (UTC, ISO-8601)")
		goVer  = fs.String("goversion", runtime.Version(), "Go version that ran the benchmarks")
		bt     = fs.String("benchtime", "", "the -benchtime value used")
		count  = fs.String("count", "", "the -count value used")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *out == "" {
		return fmt.Errorf("capture: -out is required")
	}
	if *env == "" {
		return fmt.Errorf("capture: -env is required")
	}
	if *tier != tierPrimary && *tier != tierSecondary {
		return fmt.Errorf("capture: -tier must be %q or %q", tierPrimary, tierSecondary)
	}

	raw, err := readInput(*in)
	if err != nil {
		return err
	}

	meta := Meta{
		Env: *env, Label: *label, Tier: *tier, Machine: *mach,
		Commit: *commit, Date: *date, GoVersion: *goVer, Benchtime: *bt, Count: *count,
	}

	// Validate: the CSV must contain conforming benchmarks before we commit it.
	ds, skipped, err := LoadDataset(strings.NewReader(raw), meta)
	if err != nil {
		return err
	}
	if len(ds.Samples) == 0 {
		return fmt.Errorf("capture: no conforming benchmarks in input (expected names like Hash_Get/size_1000)")
	}
	if skipped > 0 {
		fmt.Fprintf(os.Stderr, "note: %d non-conforming benchmark name(s) skipped\n", skipped)
	}

	content := metaPreamble(meta) + raw
	if err := writeFile(*out, content); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "captured %d samples for env %q → %s\n", len(ds.Samples), *env, *out)
	return nil
}

// runRender reads every committed dataset and emits the combined report, chart,
// and README region.
func runRender(args []string) error {
	fs := flag.NewFlagSet("render", flag.ContinueOnError)
	var (
		dir        = fs.String("dir", "docs/bench", "directory of per-environment dataset CSVs")
		readmePath = fs.String("readme", "README.md", "README file to inject the preview region into")
		reportPath = fs.String("report", "BENCHMARKS.md", "output path for the full report")
		svgPath    = fs.String("svg", "docs/bench.svg", "output path for the chart (also the in-doc reference)")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}

	paths, err := filepath.Glob(filepath.Join(*dir, "*.csv"))
	if err != nil {
		return err
	}
	sort.Strings(paths)
	if len(paths) == 0 {
		return fmt.Errorf("render: no dataset CSVs found in %s (run `benchreport capture` first)", *dir)
	}

	var datasets []Dataset
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		ds, skipped, err := LoadDataset(f, Meta{})
		f.Close()
		if err != nil {
			return fmt.Errorf("loading %s: %w", p, err)
		}
		if len(ds.Samples) == 0 {
			fmt.Fprintf(os.Stderr, "warning: %s has no conforming benchmarks; skipping\n", p)
			continue
		}
		if skipped > 0 {
			fmt.Fprintf(os.Stderr, "note: %s — skipped %d non-conforming name(s)\n", p, skipped)
		}
		datasets = append(datasets, ds)
	}
	if len(datasets) == 0 {
		return fmt.Errorf("render: no usable datasets in %s", *dir)
	}

	primary := primaryDataset(datasets)
	if _, missing := resolveHeadlines(indexSamples(primary.Samples)); len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "note: %d headline(s) absent from %q data, skipped: %v\n",
			len(missing), datasetLabel(primary.Meta), missing)
	}
	rows, _ := resolveHeadlines(indexSamples(primary.Samples))

	if err := writeFile(*svgPath, RenderChart(rows, datasetLabel(primary.Meta))); err != nil {
		return err
	}
	if err := writeFile(*reportPath, RenderReport(datasets, *svgPath)); err != nil {
		return err
	}

	readme, err := os.ReadFile(*readmePath)
	if err != nil {
		return err
	}
	updated, err := InjectRegion(string(readme), RenderReadmeRegion(datasets, *svgPath, *reportPath))
	if err != nil {
		return err
	}
	if err := writeFile(*readmePath, updated); err != nil {
		return err
	}

	envs := make([]string, len(datasets))
	for i, d := range datasets {
		envs[i] = d.Meta.Env
	}
	fmt.Fprintf(os.Stderr, "rendered %s, %s, and the %s preview region from %d environment(s): %v\n",
		*reportPath, *svgPath, *readmePath, len(datasets), envs)
	return nil
}

func readInput(path string) (string, error) {
	if path == "-" {
		b, err := io.ReadAll(os.Stdin)
		return string(b), err
	}
	b, err := os.ReadFile(path)
	return string(b), err
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
