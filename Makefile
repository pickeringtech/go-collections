# Makefile — developer entry points. The benchmark-report pipeline (issue #50)
# lives here so the main-only CI job can stay a thin wrapper over `make
# bench-report`, keeping local and CI generation in lockstep (no logic drift).

# Bench knobs — mirror the CI env in .github/workflows/ci.yml. Override on the
# command line for a fast local preview, e.g.:
#   make bench-report BENCH_TIME=10ms BENCH_COUNT=2
BENCH_TIME  ?= 50ms
BENCH_COUNT ?= 8
# Scope the report to the standardized collections matrix (dicts/lists/sets).
BENCH_PKGS  ?= ./collections/...

# Which environment this run captures. The report surfaces two datasets, each
# refreshed independently and committed under docs/bench/:
#   - reference (primary): a fixed, controlled machine — the trustworthy
#     baseline that drives the README headline table + chart. Refreshed by a
#     maintainer running `make bench-report` locally.
#   - ci (secondary): the shared, noisy GitHub-hosted runner — indicative only.
#     Refreshed by the main-only CI job, which overrides these vars.
# Default to the reference environment for a local run. If you refresh the
# reference from a different machine, override BENCH_LABEL/BENCH_MACHINE.
BENCH_ENV     ?= reference
BENCH_TIER    ?= primary
BENCH_LABEL   ?= Reference — Framework Desktop
BENCH_MACHINE ?= Framework Desktop · AMD Ryzen AI MAX+ 395 · 128 GB unified memory · Arch Linux

# benchstat from PATH by default; CI installs the pinned version first.
BENCHSTAT ?= benchstat

# Tool versions for the local CI mirror (`make lint` / `make security` / `make
# ci`). These MUST track the pins in .github/workflows/ci.yml so a green local
# run predicts a green PR — golangci-lint in particular changes its findings
# between releases, so parity matters (issue #88). Bump here and there together.
# We run each tool via `go run <module>@<version>` rather than expecting it on
# PATH: that pins the exact version with zero setup and never touches this
# repo's (dependency-free) go.mod. The first run builds the tool (then caches).
# NB: keep these as bare values — a trailing inline `# comment` would fold its
# leading spaces into the variable (Make quirk) and leak into the `go run` arg.
GOLANGCI_VERSION    ?= v2.1.6
GOVULNCHECK_VERSION ?= v1.3.0
GOSEC_VERSION       ?= v2.27.1
COVERAGE_MIN        ?= 100

BENCHREPORT_DIR := tools/benchreport
BUILD_DIR       := build
BENCH_DATA_DIR  := docs/bench

# Long-term trend store (issue #51). Each push to `main` archives that run's raw
# (multi-sample) bench output under docs/bench/history/<timestamp>_<sha>.txt so
# drift across commits is queryable and significance is recoverable. Archiving is
# OPT-IN (BENCH_HISTORY non-empty): only the consistent CI environment should
# feed the store — a maintainer's ad-hoc local `make bench-report` must not
# pollute it with numbers from a different machine. The cap bounds repo growth by
# pruning the oldest snapshots. bench-render always *reads* the store (cheap when
# empty), so the trend section renders for everyone.
BENCH_HISTORY     ?=
BENCH_HISTORY_DIR := $(BENCH_DATA_DIR)/history
BENCH_HISTORY_CAP ?= 100
BENCH_ALERT       := $(BUILD_DIR)/bench-alert.md

# Nested Go modules (examples/, tools/benchreport/, …) are SEPARATE modules that
# the root `go test ./...` never descends into — so a `make test` that only ran
# the root module gave contributors false confidence while CI tested more (#79).
# Discover them dynamically (any nested go.mod; `-not -path ./go.mod` drops the
# root's own) so new modules are picked up automatically and this can't drift as
# modules are added. Prune .git and the build dir so the walk doesn't descend
# into large/irrelevant trees on every `make` invocation.
NESTED_MODULES := $(shell find . \( -name .git -o -path ./$(BUILD_DIR) \) -prune \
	-o -name go.mod -not -path ./go.mod -print | xargs -r -n1 dirname | sort)

# Provenance, computed once so the generator stays a pure function of its inputs.
GIT_SHA    := $(shell git rev-parse --short HEAD 2>/dev/null)
GEN_DATE   := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GO_VERSION := $(shell go env GOVERSION)

.PHONY: help test test-root test-nested bench bench-report bench-render \
	ci hygiene cover lint security cross-arch fuzz

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

test: test-root test-nested ## Run the full test suite — root module + every nested module

test-root: ## Run the root library test suite with race + shuffle
	go test -race -shuffle=on ./...

# Run each discovered nested module in its own dir. -shuffle=on (no -race) mirrors
# CI's examples job: those tests shell out via `go run`, so the race detector
# can't see inside the child process anyway (see .github/workflows/ci.yml). A
# failure in any module aborts the loop with a non-zero status.
test-nested: ## Run the tests of every nested module (examples/, tools/benchreport, …)
	@for dir in $(NESTED_MODULES); do \
		echo ">> testing nested module $$dir"; \
		(cd "$$dir" && go test -shuffle=on ./...) || exit 1; \
	done

# ----------------------------------------------------------------------------
# Local CI mirror (issue #88). `make ci` runs the SAME blocking gates the PR
# `CI Gate` aggregates (.github/workflows/ci.yml → ci-gate `needs:`), so a green
# `make ci` predicts a green PR. Each gate is also a standalone target so you can
# reproduce a single failing job. The matching CI jobs are:
#   hygiene    → build (compile · go.mod tidy · zero-dep · integrity)
#   cover      → test  (root suite, -race -shuffle, 100% coverage floor)
#   test-nested→ examples-e2e (+ every other nested module, see #79)
#   lint       → lint  (gofmt · go vet · golangci-lint @ pinned version)
#   security   → security (govulncheck · gosec)
#   cross-arch → cross-arch (386/arm64/s390x build+vet, 386 tests)
#   fuzz       → fuzz  (count-based smoke run of every Fuzz target)
# Report-only CI jobs (test-tip, benchmarks, api-compat) are deliberately NOT
# mirrored — they never gate a merge. Ordered cheap-/common-failure-first so a
# typical mistake (formatting, a failing test) aborts before the slow arches.
ci: hygiene lint cover test-nested security cross-arch fuzz ## Run every blocking CI gate locally — a green run predicts a green PR
	@echo ">> all blocking CI gates passed locally ✔"

hygiene: ## CI 'build' gate: compile, go.mod tidy + zero-dependency + module integrity
	go build ./...
	@echo ">> checking go.mod/go.sum are tidy"
	@go mod tidy; \
	changed=$$(git status --porcelain -- go.mod go.sum); \
	if [ -n "$$changed" ]; then \
		echo "::error::go.mod/go.sum are not tidy — run 'go mod tidy' and commit:"; \
		echo "$$changed"; git diff -- go.mod go.sum; exit 1; \
	fi; \
	echo "go.mod/go.sum are tidy ✔"
	@echo ">> enforcing zero runtime dependencies"
	@if grep -qE '^[[:space:]]*require' go.mod; then \
		echo "::error::unexpected dependency in go.mod — this library stays zero-dependency"; \
		grep -nE '^[[:space:]]*require' go.mod; exit 1; \
	fi; \
	echo "go.mod declares no dependencies ✔"
	go mod verify

cover: ## CI 'test' gate: root suite with -race -shuffle, then enforce the coverage floor
	@mkdir -p $(BUILD_DIR)
	go test -race -shuffle=on -coverprofile=$(BUILD_DIR)/coverage.out ./...
	@go tool cover -func=$(BUILD_DIR)/coverage.out | tail -1
	@pct=$$(go tool cover -func=$(BUILD_DIR)/coverage.out | tail -1 | awk '{print $$3}' | tr -d '%'); \
	echo "total coverage: $${pct}% (floor: $(COVERAGE_MIN)%)"; \
	if awk -v pct="$$pct" -v min="$(COVERAGE_MIN)" 'BEGIN{exit !(pct >= min)}'; then \
		echo "coverage meets the $(COVERAGE_MIN)% floor ✔"; \
	else \
		echo "::error::coverage $${pct}% dropped below the $(COVERAGE_MIN)% floor"; exit 1; \
	fi

lint: ## CI 'lint' gate: gofmt check + go vet + golangci-lint (pinned to CI's version)
	@echo ">> checking gofmt"
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "::error::these files are not gofmt-clean:"; echo "$$unformatted"; exit 1; \
	fi; \
	echo "all files are gofmt-clean ✔"
	go vet ./...
	@echo ">> golangci-lint run ($(GOLANGCI_VERSION))"
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION) run

security: ## CI 'security' gate: govulncheck + gosec (pinned to CI's versions)
	@echo ">> govulncheck ($(GOVULNCHECK_VERSION))"
	go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...
	@echo ">> gosec ($(GOSEC_VERSION))"
	go run github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION) ./...

# Cross-compile every package for the non-amd64 arches CI gates on. 386 binaries
# execute on the amd64 host (pure-Go, static), so we also run its tests; arm64
# and s390x (big-endian) can't run here, so they are build+vet only — matching
# the CI matrix in .github/workflows/ci.yml.
cross-arch: ## CI 'cross-arch' gate: 386/arm64/s390x build+vet (+ run 386 tests)
	@for arch in 386 arm64 s390x; do \
		echo ">> GOARCH=$$arch build + vet"; \
		GOARCH=$$arch go build ./... || exit 1; \
		GOARCH=$$arch go vet ./...   || exit 1; \
	done
	@echo ">> GOARCH=386 go test -shuffle=on (executes on the amd64 host)"
	GOARCH=386 go test -shuffle=on ./...

# Count-based smoke run of every Fuzz target (mirrors CI's -fuzztime=2000x): a
# fixed iteration count is sub-second per target and deterministic, so it just
# confirms each target compiles, runs and survives its seed corpus. CI fans the
# targets out across cores; running them serially here keeps the Makefile simple.
fuzz: ## CI 'fuzz' gate: run every Fuzz target for a fixed iteration count
	@for pkg in $$(go list ./...); do \
		dir=$$(go list -f '{{.Dir}}' "$$pkg"); \
		for fn in $$(grep -rhoE '^func (Fuzz[A-Za-z0-9_]+)' "$$dir"/*_test.go 2>/dev/null | awk '{print $$2}' | sort -u); do \
			echo ">> fuzzing $$pkg $$fn"; \
			go test -run='^$$' -fuzz="^$$fn$$" -fuzztime=2000x "$$pkg" || exit 1; \
		done; \
	done
	@echo "all fuzz targets passed the count-based smoke run ✔"

bench: ## Run the collections benchmarks once (no report), printing results
	go test -run='^$$' -bench=. -benchmem -benchtime=$(BENCH_TIME) -count=$(BENCH_COUNT) $(BENCH_PKGS)

# Scope bash + pipefail to bench-report so the `go test … | tee` pipeline fails
# the target when the benchmarks fail, instead of being masked by tee's exit
# status (the rest of the Makefile keeps POSIX-sh defaults).
bench-report: SHELL := bash
bench-report: .SHELLFLAGS := -o pipefail -c
bench-report: ## Benchmark this environment, capture its dataset, and regenerate the combined report
	@mkdir -p $(BUILD_DIR) $(BENCH_DATA_DIR)
	@echo ">> benchmarking $(BENCH_PKGS) for env '$(BENCH_ENV)' (-benchtime=$(BENCH_TIME) -count=$(BENCH_COUNT))"
	go test -run='^$$' -bench=. -benchmem -benchtime=$(BENCH_TIME) -count=$(BENCH_COUNT) $(BENCH_PKGS) \
		| tee $(BUILD_DIR)/bench.txt
	@echo ">> summarising with benchstat"
	$(BENCHSTAT) -format=csv $(BUILD_DIR)/bench.txt > $(BUILD_DIR)/bench.csv
	go build -C $(BENCHREPORT_DIR) -o $(CURDIR)/$(BUILD_DIR)/benchreport .
	@echo ">> capturing dataset → $(BENCH_DATA_DIR)/$(BENCH_ENV).csv"
	$(CURDIR)/$(BUILD_DIR)/benchreport capture \
		-in $(BUILD_DIR)/bench.csv \
		-out $(BENCH_DATA_DIR)/$(BENCH_ENV).csv \
		-env "$(BENCH_ENV)" \
		-label "$(BENCH_LABEL)" \
		-tier "$(BENCH_TIER)" \
		-machine "$(BENCH_MACHINE)" \
		-commit "$(GIT_SHA)" \
		-date "$(GEN_DATE)" \
		-goversion "$(GO_VERSION)" \
		-benchtime "$(BENCH_TIME)" \
		-count "$(BENCH_COUNT)"
	@if [ -n "$(BENCH_HISTORY)" ]; then \
		echo ">> archiving trend snapshot → $(BENCH_HISTORY_DIR) (cap $(BENCH_HISTORY_CAP))"; \
		$(CURDIR)/$(BUILD_DIR)/benchreport history \
			-in $(BUILD_DIR)/bench.txt \
			-dir $(BENCH_HISTORY_DIR) \
			-commit "$(GIT_SHA)" \
			-date "$(GEN_DATE)" \
			-cap $(BENCH_HISTORY_CAP); \
	fi
	@$(MAKE) --no-print-directory bench-render

bench-render: ## Re-render BENCHMARKS.md, docs/bench.svg, and the README preview from committed datasets
	go build -C $(BENCHREPORT_DIR) -o $(CURDIR)/$(BUILD_DIR)/benchreport .
	@echo ">> rendering combined report from $(BENCH_DATA_DIR)/*.csv (+ trend from $(BENCH_HISTORY_DIR))"
	@mkdir -p $(BUILD_DIR)
	$(CURDIR)/$(BUILD_DIR)/benchreport render \
		-dir $(BENCH_DATA_DIR) \
		-history $(BENCH_HISTORY_DIR) \
		-alert $(BENCH_ALERT)
