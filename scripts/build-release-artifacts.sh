#!/bin/bash
set -eu

OUTPUT_DIR="${1:-./dist}"
VERSION="${VERSION:-0.0.0-dev}"
COMMIT="${COMMIT:-unknown}"
BUILD_DATE="${BUILD_DATE:-$(date -u +%Y-%m-%d)}"

if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+ ]]; then
    echo "Error: version must be in semver format (e.g. 1.0.0), got: $VERSION" >&2
    exit 1
fi

mkdir -p "$OUTPUT_DIR"

build_target() {
    local goos="$1" goarch="$2" ext="$3"
    local binary="taskcapsule${ext}"
    local artifact="taskcapsule_${VERSION}_${goos}_${goarch}"

    echo "Building $artifact..."

    GOOS="$goos" GOARCH="$goarch" go build \
        -trimpath \
        -ldflags "-s -w \
            -X github.com/vtino17/taskcapsule/internal/version.Version=${VERSION} \
            -X github.com/vtino17/taskcapsule/internal/version.Commit=${COMMIT} \
            -X github.com/vtino17/taskcapsule/internal/version.BuildDate=${BUILD_DATE}" \
        -o "$OUTPUT_DIR/$artifact/$binary" \
        ./cmd/taskcapsule

    cp LICENSE README.md "$OUTPUT_DIR/$artifact/"

    if [ "$goos" = "windows" ]; then
        (cd "$OUTPUT_DIR/$artifact" && zip -q "../${artifact}.zip" ./*)
    else
        tar czf "$OUTPUT_DIR/${artifact}.tar.gz" -C "$OUTPUT_DIR/${artifact}" .
    fi

    rm -rf "$OUTPUT_DIR/$artifact"
    echo "  -> $OUTPUT_DIR/${artifact}.tar.gz"
}

build_target linux   amd64   ""
build_target linux   arm64   ""
build_target darwin  amd64   ""
build_target darwin  arm64   ""
build_target windows amd64   ".exe"

cd "$OUTPUT_DIR"
sha256sum *.tar.gz *.zip > checksums.txt
echo "Checksums written to $OUTPUT_DIR/checksums.txt"
