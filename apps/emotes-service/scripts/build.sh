#!/usr/bin/env bash
set -euo pipefail

HANDLERS_DIR="src/adapters/primary"
DIST_DIR="dist"

rm -rf "$DIST_DIR"

for handler in "$HANDLERS_DIR"/*/; do
  [ -f "${handler}main.go" ] || continue
  name=$(basename "$handler")
  echo "Building $name..."
  mkdir -p "$DIST_DIR/$name"
  GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
    go build -tags lambda.norpc -ldflags="-s -w" \
    -o "$DIST_DIR/$name/bootstrap" "./$handler"
done
