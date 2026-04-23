#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
OUT_DIR="$SCRIPT_DIR/linux-debug"
EXE_NAME="big-black-box"

echo "🐧 Building Linux DEBUG..."
mkdir -p "$OUT_DIR"
cd "$PROJECT_ROOT" || exit

# No stripping (-s -w) so we get good stack traces
go build -o "$OUT_DIR/$EXE_NAME" ./packages

if [ $? -eq 0 ]; then
    echo "✅ Success: $OUT_DIR/$EXE_NAME"
    chmod +x "$OUT_DIR/$EXE_NAME"
fi
