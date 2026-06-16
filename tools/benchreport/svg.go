package main

import (
	"fmt"
	"strconv"
	"strings"
)

// SVG chart geometry. SVG (not PNG) keeps the committed artifact textual, so
// git diffs stay reviewable and the file is hand-emitted with zero dependencies.
const (
	svgWidth   = 760
	svgPadX    = 16
	svgTop     = 56 // room for the title
	svgRowH    = 34
	svgBarH    = 20
	svgLabelW  = 230 // left gutter for the operation labels
	svgValueW  = 80  // right gutter for the value text
	svgBarPalH = 8
)

// barColors cycles through a small, colour-blind-friendly palette so adjacent
// bars are distinguishable. Deterministic: bar i uses barColors[i % len].
var barColors = []string{"#4c78a8", "#72b7b2", "#54a24b", "#eeca3b", "#e45756", "#b279a2"}

// RenderChart hand-emits an SVG horizontal bar chart of ns/op for the resolved
// headline operations. Output is deterministic for fixed input (no timestamps,
// no map iteration), so re-running with unchanged data is a no-op diff.
func RenderChart(rows []headlineRow) string {
	height := svgTop + len(rows)*svgRowH + 16
	barAreaW := svgWidth - 2*svgPadX - svgLabelW - svgValueW

	maxNs := 0.0
	for _, r := range rows {
		if r.NsOp > maxNs {
			maxNs = r.NsOp
		}
	}
	if maxNs <= 0 {
		maxNs = 1 // avoid divide-by-zero when every bar is zero
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" role="img" aria-label="Benchmark ns/op by operation">`,
		svgWidth, height, svgWidth, height))
	b.WriteByte('\n')
	b.WriteString(`<style>` +
		`text{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif}` +
		`.title{font-size:15px;font-weight:600;fill:#1f2328}` +
		`.label{font-size:12px;fill:#1f2328}` +
		`.value{font-size:12px;fill:#57606a}` +
		`</style>` + "\n")
	b.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#ffffff"/>`, svgWidth, height))
	b.WriteByte('\n')
	b.WriteString(fmt.Sprintf(`<text class="title" x="%d" y="32">Benchmark — lower is faster (ns/op, n=1000)</text>`, svgPadX))
	b.WriteByte('\n')

	barX := svgPadX + svgLabelW
	for i, r := range rows {
		rowY := svgTop + i*svgRowH
		barY := rowY + (svgRowH-svgBarH)/2
		w := r.NsOp / maxNs * float64(barAreaW)
		if w < 1 {
			w = 1 // keep a sub-pixel bar visible
		}
		labelY := barY + svgBarH/2 + 4
		color := barColors[i%len(barColors)]

		b.WriteString(fmt.Sprintf(
			`<text class="label" x="%d" y="%d" text-anchor="end">%s</text>`,
			barX-10, labelY, escapeXML(r.Label)))
		b.WriteByte('\n')
		b.WriteString(fmt.Sprintf(
			`<rect x="%d" y="%d" width="%s" height="%d" rx="3" fill="%s"/>`,
			barX, barY, trimFloat(w), svgBarH, color))
		b.WriteByte('\n')
		b.WriteString(fmt.Sprintf(
			`<text class="value" x="%s" y="%d">%s</text>`,
			trimFloat(float64(barX)+w+6), labelY, formatNs(r.NsOp)))
		b.WriteByte('\n')
	}

	b.WriteString("</svg>\n")
	return b.String()
}

// trimFloat formats a coordinate with at most two decimals and no trailing
// zeros, keeping the SVG compact and stable.
func trimFloat(v float64) string {
	s := strconv.FormatFloat(v, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-0" {
		return "0"
	}
	return s
}

func escapeXML(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	return r.Replace(s)
}
