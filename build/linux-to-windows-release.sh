#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
OUT_DIR="$SCRIPT_DIR/windows-release"
EXE_NAME="big-black-box.exe"

echo "🪟 Building Windows RELEASE..."
rm -rf "$OUT_DIR" && mkdir -p "$OUT_DIR"
cd "$PROJECT_ROOT" || exit

# -H=windowsgui is the magic flag that hides the black console box
env GOOS=windows GOARCH=amd64 go build \
    -ldflags="-s -w -H=windowsgui" \
    -trimpath \
    -o "$OUT_DIR/$EXE_NAME" ./packages

echo "✅ Success: $OUT_DIR/$EXE_NAME"
