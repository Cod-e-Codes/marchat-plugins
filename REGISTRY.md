# Registry: publishing platform-specific plugins

This repository hosts a simple JSON registry (`registry.json`) that points to downloadable plugin archives. With the single-binary model, each plugin release contains one archive targeted to a specific OS/architecture.

## Steps to publish

1. Build the platform-specific archive
   - Use the scripts under your plugin directory (example: `plugins/echo`):
     - Windows (PowerShell): `./build.ps1 -Version 1.0.0 -Os windows -Arch amd64`
     - Linux/macOS (Bash): `VERSION=1.0.0 GOOS=linux GOARCH=amd64 bash build.sh`
   - Output: `dist/<name>-plugin-v<version>-<goos>-<goarch>.zip` and a matching `.sha256` file.

2. Upload the archive to your release hosting
   - Example: GitHub Releases for `marchat-plugins`.
   - Copy the final download URL of the uploaded archive.

3. Capture the SHA-256 checksum
   - Use the contents of the generated `.sha256` file (the hex digest).

4. Update `registry.json`
   - Keep one entry per plugin. Point `download_url` and `checksum` to the specific platform you are publishing.
   - Update `version` and `last_updated` accordingly.

### Example `registry.json` entry
```json
{
  "version": "1.0.0",
  "last_updated": "2025-01-01T00:00:00Z",
  "plugins": [
    {
      "name": "echo",
      "version": "1.0.0",
      "description": "Simple echo plugin for testing the plugin system",
      "author": "Cod-e-Codes",
      "homepage": "https://github.com/Cod-e-Codes/marchat",
      "download_url": "https://github.com/Cod-e-Codes/marchat-plugins/releases/download/v1.0.0/echo-plugin-v1.0.0-windows-amd64.zip",
      "checksum": "sha256:<hex-digest>",
      "category": "utility",
      "tags": ["chat", "utility", "echo", "test"],
      "min_version": "0.2.0-beta.1"
    }
  ]
}
```

Notes:
- The host validates `GOOS/GOARCH` at runtime and will refuse to launch mismatched binaries with a clear error.
- If you need to support multiple platforms, publish separate releases and update `registry.json` to point at the one most relevant for your audience. The store and `:plugin list` remain unchanged.
