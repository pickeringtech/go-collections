package main

import (
	"math"
	"strconv"
	"strings"
)

// headline names one marquee benchmark cell to surface in the README preview
// table and the chart. This is a deliberately small, explicit allowlist — not
// every benchmark — so the README stays a glance-able highlight, with the full
// matrix living in BENCHMARKS.md. Entries missing from the data are skipped.
type headline struct {
	Label string
	Pkg   string
	Impl  string
	Op    string
	Size  int
}

// headlines is the curated set of "we take performance seriously" numbers. They
// centre on lookup cost (the operation a reader most wants to compare across
// implementations) at a mid-range size, plus one list and one set marquee.
var headlines = []headline{
	{"Dict — Hash.Get", "dicts", "Hash", "Get", 1000},
	{"Dict — ConcurrentHash.Get", "dicts", "ConcurrentHash", "Get", 1000},
	{"Dict — ConcurrentHashRW.Get", "dicts", "ConcurrentHashRW", "Get", 1000},
	{"Dict — Tree.Get", "dicts", "Tree", "Get", 1000},
	{"List — Array.Get", "lists", "Array", "Get", 1000},
	{"Set — Hash.Contains", "sets", "Hash", "Contains", 1000},
}

// pkgOrder fixes a stable family ordering for the report regardless of the order
// packages appear in the CSV. Unknown packages sort after these, alphabetically.
var pkgOrder = []string{"dicts", "lists", "sets"}

func pkgRank(p string) int {
	for i, n := range pkgOrder {
		if n == p {
			return i
		}
	}
	return len(pkgOrder)
}

// formatNs renders a ns/op value with adaptive precision (more decimals for
// small, sub-100ns numbers) and thousands separators, so a 40ns lookup and a
// 65,000ns insert are both readable and aligned.
func formatNs(v float64) string {
	if v == 0 {
		return "0"
	}
	var decimals int
	switch {
	case v < 10:
		decimals = 2
	case v < 100:
		decimals = 1
	default:
		decimals = 0
	}
	return withThousands(strconv.FormatFloat(v, 'f', decimals, 64))
}

// formatCount renders an integer-valued metric (B/op, allocs/op). benchstat
// reports these as whole numbers; a fractional median is shown to one decimal.
func formatCount(v float64) string {
	if v == 0 {
		return "0"
	}
	if v == math.Trunc(v) {
		return withThousands(strconv.FormatFloat(v, 'f', 0, 64))
	}
	return withThousands(strconv.FormatFloat(v, 'f', 1, 64))
}

// withThousands inserts commas into the integer part of a decimal string,
// preserving any sign and fractional part.
func withThousands(s string) string {
	sign := ""
	if strings.HasPrefix(s, "-") {
		sign, s = "-", s[1:]
	}
	intPart, frac, hasFrac := strings.Cut(s, ".")
	n := len(intPart)
	if n <= 3 {
		if hasFrac {
			return sign + intPart + "." + frac
		}
		return sign + intPart
	}
	var b strings.Builder
	lead := n % 3
	if lead > 0 {
		b.WriteString(intPart[:lead])
		if n > lead {
			b.WriteByte(',')
		}
	}
	for i := lead; i < n; i += 3 {
		b.WriteString(intPart[i : i+3])
		if i+3 < n {
			b.WriteByte(',')
		}
	}
	out := sign + b.String()
	if hasFrac {
		out += "." + frac
	}
	return out
}
