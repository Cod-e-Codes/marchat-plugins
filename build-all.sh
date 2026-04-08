#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-1.0.0}"
PLUGINS=("echo" "weather" "githooks")
PLATFORMS=("windows/amd64" "linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64")

echo "Building all plugins for all platforms (v$VERSION)"
echo ""

for plugin in "${PLUGINS[@]}"; do
  echo "=== Building $plugin ==="
  cd "plugins/$plugin"
  
  for platform in "${PLATFORMS[@]}"; do
    OS="${platform%/*}"
    ARCH="${platform#*/}"
    echo "  - $OS/$ARCH"
    VERSION=$VERSION GOOS=$OS GOARCH=$ARCH bash build.sh
  done
  
  cd ../..
  echo ""
done

echo "✅ All builds complete!"
echo ""
echo "Next steps:"
echo "1. Create GitHub release: gh release create v$VERSION"
echo "2. Upload all ZIP files from plugins/*/dist/"
echo "3. Update registry.json with checksums"
echo "4. Push changes"

