// The examples live in their OWN module so they consume the library as a
// genuine outside package: only the exported API is reachable, and they compile
// as a separate unit rather than recompiling in lockstep with the library.
//
// The `replace` points the require at the local checkout (the parent directory)
// so PR CI builds and E2E-tests the examples against the CURRENT tree — proving
// the public API still compiles and runs for a downstream consumer. The
// API-compat guarantee for *released* versions is covered separately by the
// gorelease gate (see issue #29), not here.
module github.com/pickeringtech/go-collections/examples

go 1.24

require github.com/pickeringtech/go-collections v0.0.0

replace github.com/pickeringtech/go-collections => ../
