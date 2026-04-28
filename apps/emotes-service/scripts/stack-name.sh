#!/bin/bash
set -euo pipefail

# PR_NUMBER is set in CI environments
if [[ -n "${PR_NUMBER:-}" ]]; then
  echo "pr-${PR_NUMBER}"
  exit 0
fi

whoami | cut -c1-5
