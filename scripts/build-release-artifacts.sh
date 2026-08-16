#!/bin/bash
set -euo pipefail

umask 022
export TZ=UTC

OUTPUT_DIR="${1:-./dist}"
VERSION="${VERSION:-0.0.0-dev}"
COMMIT="${COMMIT:-unknown}"
SOURCE_DATE_EPOCH="${SOURCE_DATE_EPOCH:-$(git log -1 --format=%ct 2>/dev/null || printf '315532800')}"
BUILD_DATE="${BUILD_DATE:-$(date -u -d "@$SOURCE_DATE_EPOCH" +%Y-%m-%d)}"

if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+([-][0-9A-Za-z.-]+)?$ ]]; then
    echo "Error: version must be in semver format (e.g. 1.0.0 or 1.0.0-dry-run), got: $VERSION" >&2
    exit 1
fi
if [[ ! "$COMMIT" =~ ^[0-9A-Za-z.-]+$ ]]; then
    echo "Error: commit contains unsupported characters" >&2
    exit 1
fi
if [[ ! "$BUILD_DATE" =~ ^[0-9]{4}-[0-9]{2}-[0-9]{2}$ ]]; then
    echo "Error: build date must use YYYY-MM-DD" >&2
    exit 1
fi
if [[ ! "$SOURCE_DATE_EPOCH" =~ ^[0-9]+$ ]] || (( SOURCE_DATE_EPOCH < 315532800 )); then
    echo "Error: SOURCE_DATE_EPOCH must be a Unix timestamp on or after 1980-01-01" >&2
    exit 1
fi

mkdir -p "$OUTPUT_DIR"
OUTPUT_DIR="$(cd "$OUTPUT_DIR" && pwd -P)"
STAGING_ROOT="$(mktemp -d "$OUTPUT_DIR/.taskcapsule-release.XXXXXX")"
trap 'rm -rf -- "$STAGING_ROOT"' EXIT
ARTIFACT_FILES=()

build_target() {
    local goos="$1" goarch="$2" ext="$3"
    local binary="taskcapsule${ext}"
    local artifact="taskcapsule_${VERSION}_${goos}_${goarch}"
    local staging="$STAGING_ROOT/$artifact"

    echo "Building $artifact..."
	mkdir -p "$staging"

    GOOS="$goos" GOARCH="$goarch" go build \
        -trimpath \
        -ldflags "-s -w \
            -X github.com/vtino17/taskcapsule/internal/version.Version=${VERSION} \
            -X github.com/vtino17/taskcapsule/internal/version.Commit=${COMMIT} \
            -X github.com/vtino17/taskcapsule/internal/version.BuildDate=${BUILD_DATE}" \
        -o "$staging/$binary" \
        ./cmd/taskcapsule

    cp LICENSE README.md "$staging/"
    touch -d "@$SOURCE_DATE_EPOCH" "$staging/$binary" "$staging/LICENSE" "$staging/README.md"

    if [ "$goos" = "windows" ]; then
        (cd "$staging" && zip -X -q "$OUTPUT_DIR/${artifact}.zip" ./*)
        ARTIFACT_FILES+=("${artifact}.zip")
    else
        tar --sort=name \
            --mtime="@$SOURCE_DATE_EPOCH" \
            --owner=0 --group=0 --numeric-owner \
            -cf - -C "$staging" . | gzip -n > "$OUTPUT_DIR/${artifact}.tar.gz"
        ARTIFACT_FILES+=("${artifact}.tar.gz")
    fi

    if [ "$goos" = "windows" ]; then
        echo "  -> $OUTPUT_DIR/${artifact}.zip"
    else
        echo "  -> $OUTPUT_DIR/${artifact}.tar.gz"
    fi
}

build_target linux   amd64   ""
build_target linux   arm64   ""
build_target darwin  amd64   ""
build_target darwin  arm64   ""
build_target windows amd64   ".exe"

cd "$OUTPUT_DIR"
printf '%s\n' "${ARTIFACT_FILES[@]}" | sort | xargs sha256sum > checksums.txt
echo "Checksums written to $OUTPUT_DIR/checksums.txt"
