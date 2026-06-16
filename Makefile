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

BENCHREPORT_DIR := tools/benchreport
BUILD_DIR       := build
BENCH_DATA_DIR  := docs/bench

# Nested Go modules (examples/, tools/benchreport/, …) are SEPARATE modules that
# the root `go test ./...` never descends into — so a `make test` that only ran
# the root module gave contributors false confidence while CI tested more (#79).
# Discover them dynamically (any go.mod below the root; -mindepth 2 skips the
# root's own go.mod) so new modules are picked up automatically and this can't
# drift as modules are added.
NESTED_MODULES := $(shell find . -mindepth 2 -name go.mod -exec dirname {} \; | sort)

# Provenance, computed once so the generator stays a pure function of its inputs.
GIT_SHA    := $(shell git rev-parse --short HEAD 2>/dev/null)
GEN_DATE   := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GO_VERSION := $(shell go env GOVERSION)

.PHONY: help test test-root test-nested bench bench-report bench-render

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
	@$(MAKE) --no-print-directory bench-render

bench-render: ## Re-render BENCHMARKS.md, docs/bench.svg, and the README preview from committed datasets
	go build -C $(BENCHREPORT_DIR) -o $(CURDIR)/$(BUILD_DIR)/benchreport .
	@echo ">> rendering combined report from $(BENCH_DATA_DIR)/*.csv"
	$(CURDIR)/$(BUILD_DIR)/benchreport render -dir $(BENCH_DATA_DIR)
