#!/usr/bin/env bash
set -euo pipefail

repo_root=$(git rev-parse --show-toplevel) || {
  echo "Error: run this inside a Git repository." >&2
  exit 1
}
cd "$repo_root"

config=/mnt/c/Users/ChingKokChiam/.config/golangci-lint/config.yml

[[ -f go.mod ]] || { echo "Error: go.mod not found at $repo_root." >&2; exit 1; }
[[ -f "$config" ]] || { echo "Error: lint config not found: $config" >&2; exit 1; }
git rev-parse --verify origin/develop >/dev/null 2>&1 || {
  echo "Error: origin/develop is unavailable; fetch it before linting." >&2
  exit 1
}

go get github.com/quasilyte/go-ruleguard/dsl@v0.3.23

# Custom rules can live outside the repository and may not invalidate cached results.
golangci-lint cache clean

GIT_CONFIG_COUNT=1 \
GIT_CONFIG_KEY_0=core.autocrlf \
GIT_CONFIG_VALUE_0=true \
golangci-lint run \
  --config "$config" \
  --new-from-rev=origin/develop
