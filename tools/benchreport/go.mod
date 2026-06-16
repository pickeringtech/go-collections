// benchreport lives in its OWN module (like examples/) so it never touches the
// library's root module: it adds no dependency to go.mod, sits outside the
// root's 100% coverage floor / lint / gosec gates, and is invoked only via
// `make bench-report` and the main-only CI job. It is intentionally
// dependency-free — it consumes benchstat's CSV output and emits Markdown + SVG
// using only the standard library.
module github.com/pickeringtech/go-collections/tools/benchreport

go 1.24
