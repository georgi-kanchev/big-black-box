#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Output folder is now build/release
OUT_DIR="$SCRIPT_DIR/web-release"
WASM_FILE="game.wasm"
PORT=9090

echo "🚀 Building RELEASE for Web..."

rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$OUT_DIR/"

cd "$PROJECT_ROOT" || exit
env GOOS=js GOARCH=wasm go build -ldflags="-s -w" -trimpath -o "$OUT_DIR/$WASM_FILE" ./packages

if [ $? -ne 0 ]; then
    echo "❌ Release build failed."
    exit 1
fi

cat <<EOF > "$OUT_DIR/index.html"
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Big Black Box</title>
    <style>body { margin: 0; background: #000; overflow: hidden; }</style>
    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("./$WASM_FILE"), go.importObject).then(s => go.run(s.instance));
    </script>
</head>
<body></body>
</html>
EOF

echo "✅ Release Ready! Output in build/release"
echo "🚀 Opening Preview..."

(sleep 1 && xdg-open "http://localhost:$PORT") &

cd "$OUT_DIR"
python3 -c "import http.server; \
           m = http.server.SimpleHTTPRequestHandler.extensions_map; \
           m['.wasm'] = 'application/wasm'; \
           http.server.test(HandlerClass=http.server.SimpleHTTPRequestHandler, port=$PORT)"
