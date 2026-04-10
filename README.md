# marchat Plugins

Community plugin registry and releases for [marchat](https://github.com/Cod-e-Codes/marchat) - a terminal-native, offline-first group chat application.

## Platform model (single-binary per OS/arch)

Plugins are now distributed as a single binary targeted to a specific platform (OS/architecture). The host validates the plugin binary at runtime against `runtime.GOOS` and `runtime.GOARCH`, and will return a clear error if the binary does not match the current system. This keeps plugin sizes minimal and avoids complex cross-platform packaging, while maintaining subprocess isolation and JSON-based communication.

- Expected archive name: `PLUGIN-NAME-v<version>-<goos>-<goarch>.zip`
- Inside the archive: the single plugin binary (`.exe` on Windows), an optional `plugin.json` manifest, and a `README.md`.
- Examples:
  - `echo-plugin-v2.0.1-windows-amd64.zip`
  - `echo-plugin-v2.0.1-linux-amd64.zip`
  - `echo-plugin-v2.0.1-darwin-amd64.zip`

If a user installs a plugin binary for the wrong platform, the host will refuse to start the plugin and display a helpful message indicating the required `GOOS/GOARCH`.

## Available Plugins

- **echo** (v2.0.1): Simple echo plugin for testing the plugin system
- **weather** (v1.0.0): Get weather information and forecasts using wttr.in
- **githooks** (v1.0.0): Git repository management with status, log, branch, and diff commands

## Installing Plugins

Use the marchat plugin store:
```bash
:store
```

Or install directly:
```bash
:install echo
```

## Publishing plugins

### Automated (GitHub Actions)

Use **Actions** then **Release plugin archives** (`release.yml`). Pick the plugin folder (`echo`, `weather`, or `githooks`), the **plugin version** string that appears inside each zip name, and the **GitHub release tag** to create or append to. The workflow checks out this repo plus [marchat](https://github.com/Cod-e-Codes/marchat) as a sibling (for the `go.mod` `replace` on `plugin/sdk`), runs each platform `build.sh`, then creates the tag if missing and uploads all `dist/*.zip` files with `gh release upload`. You still need to refresh `registry.json` checksums and URLs after publishing (see `REGISTRY.md` and `collect-checksums.sh` / `collect-checksums.ps1`).

### Manual

When publishing without CI, build and upload one archive per target platform and provide the corresponding checksums. A typical process:

1. Build platform-specific archives (see `plugins/echo/README.md` for examples and scripts):
   - `echo-plugin-v2.0.1-windows-amd64.zip`
   - `echo-plugin-v2.0.1-linux-amd64.zip`
   - `echo-plugin-v2.0.1-darwin-amd64.zip`
2. Compute SHA-256 checksums and publish alongside the archives.
3. Update `registry.json` to point to the platform-specific archive you intend to distribute for your audience. The host will still validate at runtime and show errors on mismatch.

See `REGISTRY.md` for detailed steps and an example entry.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to submit your own plugins.

## License

MIT License - see [LICENSE](LICENSE) for details.
