#!/bin/bash
set -euo pipefail

# CI_STACK_NAME is set by CI workflows.
if [[ -n "${CI_STACK_NAME:-}" ]]; then
  echo "${CI_STACK_NAME}"
  exit 0
fi

whoami | cut -c1-5
