#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-1.0.0}"
PLUGINS=("weather" "githooks")

echo "Collecting checksums for version $VERSION"
echo ""

for plugin in "${PLUGINS[@]}"; do
    echo "=== $plugin ==="
    for checksum_file in plugins/$plugin/dist/*-v${VERSION}-*.zip.sha256; do
        if [[ -f "$checksum_file" ]]; then
            basename "$checksum_file" .sha256
            cat "$checksum_file" | awk '{print "  sha256:" $1}'
            echo ""
        fi
    done
done

echo "Copy these checksums into registry.json"

