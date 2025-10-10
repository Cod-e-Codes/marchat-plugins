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
BIN_NAME="githooks$EXT"
BIN_PATH="$STAGE_DIR/$BIN_NAME"

echo "Building githooks plugin for $OS/$ARCH v$VERSION"
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
  "name": "githooks",
  "version": "$VERSION",
  "description": "Git repository management and status updates",
  "author": "Cod-e-Codes",
  "license": "MIT",
  "commands": [
    { "name": "git-status", "description": "Show git status of current directory", "usage": ":git-status [path]", "admin_only": false },
    { "name": "git-log", "description": "Show recent git commits", "usage": ":git-log [n] [path]", "admin_only": false },
    { "name": "git-branch", "description": "Show current git branch and available branches", "usage": ":git-branch [path]", "admin_only": false },
    { "name": "git-diff", "description": "Show git diff of uncommitted changes", "usage": ":git-diff [path]", "admin_only": false },
    { "name": "git-watch", "description": "Watch a repository for changes (admin only)", "usage": ":git-watch <path>", "admin_only": true }
  ]
}
JSON

cp "$SCRIPT_DIR/README.md" "$STAGE_DIR/README.md"

ZIP_NAME="githooks-plugin-v${VERSION}-${OS}-${ARCH}.zip"
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

