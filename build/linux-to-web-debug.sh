#!/bin/bash

# --- Setup Paths ---
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Changed output folder name to 'debug' inside the 'build' directory
OUT_DIR="$SCRIPT_DIR/web-debug"
WASM_FILE="game.wasm"
PORT=8080

echo "🛠️  Building DEBUG for Web..."

# 1. Prepare debug folder inside /build
# We wipe the old debug folder to ensure a fresh build
rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$OUT_DIR/"

# 2. Compile (Running from root, outputting to build/debug)
cd "$PROJECT_ROOT" || exit
echo "Compiling..."
env GOOS=js GOARCH=wasm go build -o "$OUT_DIR/$WASM_FILE" ./packages

if [ $? -ne 0 ]; then
    echo "❌ Build failed. Check your Go code/imports."
    exit 1
fi

# 3. Create index.html in build/debug
cat <<EOF > "$OUT_DIR/index.html"
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Big Black Box [debug]</title>
    <style>
        body { margin: 0; background: #1a1a1a; color: #00ff00; font-family: monospace; }
        canvas { display: block; margin: 0 auto; }
    </style>
    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("./$WASM_FILE"), go.importObject).then((result) => {
            go.run(result.instance);
        }).catch(err => {
            console.error("Wasm Loading Error:", err);
        });
    </script>
</head>
<body></body>
</html>
EOF

echo "✅ Success! Files are in build/debug"
echo "🌐 Server starting at http://localhost:$PORT"

# 4. Launch Browser
(sleep 1 && xdg-open "http://localhost:$PORT") &

# 5. Serve from build/debug
cd "$OUT_DIR"
python3 -c "import http.server; \
           m = http.server.SimpleHTTPRequestHandler.extensions_map; \
           m['.wasm'] = 'application/wasm'; \
           http.server.test(HandlerClass=http.server.SimpleHTTPRequestHandler, port=$PORT)"
