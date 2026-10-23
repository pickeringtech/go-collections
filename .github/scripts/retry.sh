#!/usr/bin/env bash
# retry.sh — run a command, retrying with exponential backoff while it fails.
#
# Why: installing a pinned tool (`go install tool@version`) downloads the tool's
# module graph, and under GOTOOLCHAIN=auto a Go toolchain too, from
# proxy.golang.org. A transient proxy or network error ("stream error: ...
# INTERNAL_ERROR; received from peer") then fails the job, and in the merge
# queue drops every PR in the group, though nothing is wrong with the change.
# Retrying gets past the blip. A real failure (a bad version, a build error)
# fails every attempt, so the step still fails, with the command's exit status.
#
# Usage:  .github/scripts/retry.sh go install "example.com/tool@${VERSION}"
# Tune:   RETRY_ATTEMPTS (default 4) and RETRY_DELAY, the seconds before the
#         first retry (default 5), which doubles after each further failure.
set -euo pipefail
attempts=${RETRY_ATTEMPTS:-4}
delay=${RETRY_DELAY:-5}
for ((i = 1; ; i++)); do
	"$@" && exit 0
	status=$?
	if ((i >= attempts)); then
		echo "retry.sh: giving up after $i attempts (exit $status): $*" >&2
		exit "$status"
	fi
	echo "retry.sh: attempt $i/$attempts failed (exit $status); retrying in ${delay}s: $*" >&2
	sleep "$delay"
	delay=$((delay * 2))
done
