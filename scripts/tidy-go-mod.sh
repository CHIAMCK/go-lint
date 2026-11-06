#!/usr/bin/env bash
set -euo pipefail

repo_root=$(git rev-parse --show-toplevel) || {
  echo "Error: run this inside a Git repository." >&2
  exit 1
}
cd "$repo_root"

[[ -f go.mod ]] || { echo "Error: go.mod not found at $repo_root." >&2; exit 1; }

go mod tidy

if git diff --quiet -- go.mod go.sum; then
  echo "No tracked go.mod or go.sum changes."
  exit 0
fi

echo "Review these go.mod and go.sum changes:"
git --no-pager diff -- go.mod go.sum || [[ $? -eq 1 ]]
