#!/usr/bin/env sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
export GOCACHE="$root/.cache/go-build"

go test ./...
