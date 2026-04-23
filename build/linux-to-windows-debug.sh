#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
OUT_DIR="$SCRIPT_DIR/windows-debug"
EXE_NAME="big-black-box.exe"

echo "🪟 Building Windows DEBUG..."
mkdir -p "$OUT_DIR"
cd "$PROJECT_ROOT" || exit

env GOOS=windows GOARCH=amd64 go build -o "$OUT_DIR/$EXE_NAME" ./packages

echo "✅ Success: $OUT_DIR/$EXE_NAME"
