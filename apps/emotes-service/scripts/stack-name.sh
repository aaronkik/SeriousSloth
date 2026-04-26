#!/bin/bash
set -euo pipefail

if [[ "${CI:-}" == "true" && "${GITHUB_REF:-}" =~ ^refs/pull/([0-9]+)/ ]]; then
  echo "pr-${BASH_REMATCH[1]}"
  exit 0
fi

whoami | cut -c1-5
