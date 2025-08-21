# marchat Plugins

Community plugin registry and releases for [marchat](https://github.com/Cod-e-Codes/marchat) - a terminal-native, offline-first group chat application.

## Platform model (single-binary per OS/arch)

Plugins are now distributed as a single binary targeted to a specific platform (OS/architecture). The host validates the plugin binary at runtime against `runtime.GOOS` and `runtime.GOARCH`, and will return a clear error if the binary does not match the current system. This keeps plugin sizes minimal and avoids complex cross-platform packaging, while maintaining subprocess isolation and JSON-based communication.

- Expected archive name: `PLUGIN-NAME-v<version>-<goos>-<goarch>.zip`
- Inside the archive: the single plugin binary (`.exe` on Windows), an optional `plugin.json` manifest, and a `README.md`.
- Examples:
  - `echo-plugin-v1.0.0-windows-amd64.zip`
  - `echo-plugin-v1.0.0-linux-amd64.zip`
  - `echo-plugin-v1.0.0-darwin-arm64.zip`

If a user installs a plugin binary for the wrong platform, the host will refuse to start the plugin and display a helpful message indicating the required `GOOS/GOARCH`.

## Available Plugins

- **echo**: Simple echo plugin for testing the plugin system

## Installing Plugins

Use the marchat plugin store:
```bash
:store
```

Or install directly:
```bash
:plugin install echo
```

## Publishing plugins

When publishing a new plugin release, build and upload one archive per target platform and provide the corresponding checksums. A typical process:

1. Build platform-specific archives (see `plugins/echo/README.md` for examples and scripts):
   - `echo-plugin-v2.0.0-windows-amd64.zip`
   - `echo-plugin-v2.0.0-linux-amd64.zip`
   - `echo-plugin-v2.0.0-darwin-amd64.zip`
2. Compute SHA-256 checksums and publish alongside the archives.
3. Update `registry.json` to point to the platform-specific archive you intend to distribute for your audience. The host will still validate at runtime and show errors on mismatch.

See `REGISTRY.md` for detailed steps and an example entry.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to submit your own plugins.

## License

MIT License - see [LICENSE](LICENSE) for details.
