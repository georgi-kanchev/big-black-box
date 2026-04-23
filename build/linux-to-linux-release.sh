#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
OUT_DIR="$SCRIPT_DIR/linux-release"
EXE_NAME="big-black-box"

echo "🐧 Building Linux RELEASE..."
rm -rf "$OUT_DIR" && mkdir -p "$OUT_DIR"
cd "$PROJECT_ROOT" || exit

go build -ldflags="-s -w" -trimpath -o "$OUT_DIR/$EXE_NAME" ./packages

if [ $? -eq 0 ]; then
    echo "✅ Success: $OUT_DIR/$EXE_NAME"
    chmod +x "$OUT_DIR/$EXE_NAME"
fi
