#!/bin/bash
set -euo pipefail

OUTPUT_DIR="${1:-./dist}"
SOURCE_DATE_EPOCH="${SOURCE_DATE_EPOCH:-$(git log -1 --format=%ct)}"
export SOURCE_DATE_EPOCH

VERIFY_ROOT="$(mktemp -d "${TMPDIR:-/tmp}/taskcapsule-reproducible.XXXXXX")"
FIRST="$VERIFY_ROOT/first"
SECOND="$VERIFY_ROOT/second"
trap 'rm -rf -- "$VERIFY_ROOT"' EXIT

mkdir -p "$FIRST" "$SECOND" "$OUTPUT_DIR"
bash scripts/build-release-artifacts.sh "$FIRST"
bash scripts/build-release-artifacts.sh "$SECOND"

diff -r "$FIRST" "$SECOND"
cp "$FIRST"/* "$OUTPUT_DIR"/
echo "Release artifacts are reproducible for SOURCE_DATE_EPOCH=$SOURCE_DATE_EPOCH"
