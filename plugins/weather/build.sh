#!/usr/bin/env bash
set -euo pipefail

VERSION="${VERSION:-1.0.0}"
OS="${GOOS:-${OS_OVERRIDE:-$(uname | tr '[:upper:]' '[:lower:]')}}"
ARCH="${GOARCH:-amd64}"
OUT_DIR="${OUT_DIR:-dist}"

case "$OS" in
  mingw*|msys*|cygwin*|windows) OS=windows ;;
  linux) OS=linux ;;
  darwin) OS=darwin ;;
  *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

if ! command -v go >/dev/null 2>&1; then
  echo "Go toolchain not found. Install Go and ensure 'go' is on PATH." >&2
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="$SCRIPT_DIR/$OUT_DIR"
STAGE_DIR="$DIST_DIR/stage-$OS-$ARCH"
mkdir -p "$STAGE_DIR"

EXT=""
if [[ "$OS" == "windows" ]]; then EXT=".exe"; fi
BIN_NAME="weather$EXT"
BIN_PATH="$STAGE_DIR/$BIN_NAME"

echo "Building weather plugin for $OS/$ARCH v$VERSION"
(
  cd "$SCRIPT_DIR"
  GOOS="$OS" GOARCH="$ARCH" go build -trimpath -ldflags "-s -w" -o "$BIN_PATH" .
)

if [[ ! -f "$BIN_PATH" ]]; then
  echo "Build failed; binary not found at $BIN_PATH" >&2
  exit 1
fi

# Create plugin.json
cat > "$STAGE_DIR/plugin.json" <<JSON
{
  "name": "weather",
  "version": "$VERSION",
  "description": "Get weather information using wttr.in",
  "author": "Cod-e-Codes",
  "license": "MIT",
  "commands": [
    { "name": "weather", "description": "Get current weather for a location", "usage": ":weather [location]", "admin_only": false },
    { "name": "forecast", "description": "Get weather forecast for a location", "usage": ":forecast [location]", "admin_only": false }
  ]
}
JSON

cp "$SCRIPT_DIR/README.md" "$STAGE_DIR/README.md"

ZIP_NAME="weather-plugin-v${VERSION}-${OS}-${ARCH}.zip"
ZIP_PATH="$DIST_DIR/$ZIP_NAME"
rm -f "$ZIP_PATH"

echo "Creating archive $ZIP_NAME"
(
  cd "$STAGE_DIR"
  # Prefer 7z or zip
  if command -v 7z >/dev/null 2>&1; then
    7z a -tzip -mx=9 "$ZIP_PATH" * >/dev/null
  else
    zip -9 -r "$ZIP_PATH" . >/dev/null
  fi
)

echo "Computing SHA-256 checksum"
if command -v shasum >/dev/null 2>&1; then
  shasum -a 256 "$ZIP_PATH" > "$ZIP_PATH.sha256"
elif command -v sha256sum >/dev/null 2>&1; then
  (cd "$DIST_DIR" && sha256sum "$ZIP_NAME" > "$ZIP_NAME.sha256")
else
  echo "No sha256 tool found (shasum/sha256sum). Skipping checksum." >&2
fi

echo "Done: $ZIP_PATH"

