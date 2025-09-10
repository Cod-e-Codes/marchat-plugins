# Echo Plugin

A simple test plugin for the `marchat` terminal chat application that demonstrates basic plugin functionality.

## Description

The Echo Plugin is designed to test and demonstrate the marchat plugin system. It provides basic echo functionality and serves as a reference implementation for plugin development.

## Features

- **Message Echoing**: Repeats messages that start with `echo:`
- **Command Support**: Provides `:echo` command for direct echoing
- **Admin Commands**: Includes admin-only echo functionality
- **Plugin Testing**: Perfect for verifying plugin system functionality

## Commands

### User Commands
- `:echo <message>` — Echoes the message back to the chat

### Admin Commands  
- `:echo-admin <message>` — Echoes the message (admin-only)

### Message Processing
- `echo: <message>` — Automatically echoes messages starting with "echo:"

## Installation

### From Plugin Store
```bash
:store
# Navigate to echo plugin and press Enter to install
```

### Direct Installation
```bash
:plugin install echo
```

## Build and package (single-binary per OS/arch)

Builds produce a single archive targeted to your OS/architecture. The host validates the binary at runtime against `GOOS/GOARCH` and will error on mismatches.

### PowerShell (Windows)
```powershell
# From this directory
./build.ps1 -Version 2.0.0 -Os windows -Arch amd64
# Output: dist/echo-plugin-v2.0.0-windows-amd64.zip and .sha256
```

### Bash (Linux/macOS)
```bash
# From this directory
VERSION=2.0.0 bash build.sh # uses GOOS/GOARCH from env
# or override
VERSION=2.0.0 GOOS=linux GOARCH=amd64 bash build.sh
# Output: dist/echo-plugin-v2.0.0-linux-amd64.zip and .sha256
```

### Archive contents
- `echo` or `echo.exe` — the plugin binary
- `plugin.json` — optional manifest with metadata and commands
- `README.md` — documentation

### Naming convention
- `echo-plugin-v<version>-<goos>-<goarch>.zip`
  - Examples: `echo-plugin-v2.0.0-windows-amd64.zip`, `echo-plugin-v2.0.0-darwin-amd64.zip`

## Usage Examples

### Basic Echo Command
```bash
:echo Hello, world!
# Output: EchoBot: Hello, world!
```

### Message Echoing
```bash
echo: This is a test message
# Output: EchoBot: This is a test message
```

### Admin Echo Command
```bash
:echo-admin Admin test message
# Output: EchoBot: Admin test message
```

## Plugin Management

### List Installed Plugins
```bash
:plugin list
```

### Enable/Disable Plugin
```bash
:plugin enable echo
:plugin disable echo
```

### Uninstall Plugin (Admin Only)
```bash
:plugin uninstall echo
```

## Technical Details

### Plugin Structure
```
echo/
├── echo              # Binary executable
├── plugin.json       # Plugin manifest
└── README.md         # This documentation
```

### Plugin Manifest
```json
{
  "name": "echo",
  "version": "1.0.0",
  "description": "Simple echo plugin for testing the marchat plugin system",
  "author": "Cod-e-Codes",
  "license": "MIT",
  "commands": [
    {
      "name": "echo",
      "description": "Echo a message",
      "usage": ":echo <message>",
      "admin_only": false
    }
  ]
}
```

## Development

This plugin serves as a reference implementation for marchat plugin development. It demonstrates:

- **Plugin Interface**: Implements the required `Plugin` interface
- **Message Handling**: Processes incoming chat messages
- **Command Registration**: Registers commands with the host
- **Configuration**: Handles plugin configuration
- **Error Handling**: Graceful error handling and logging

### Source Code
The plugin source code is available in the main marchat repository at `plugin/examples/echo/echo.go` and serves as a template for developing new plugins.

## Version History

- **v1.0.0** - Initial release with basic echo functionality

## License

MIT License - see [LICENSE](../../LICENSE) for details.

## Contributing

This plugin is part of the marchat project. For plugin development guidelines, see the main repository's [PLUGIN_ECOSYSTEM.md](https://github.com/Cod-e-Codes/marchat/blob/main/PLUGIN_ECOSYSTEM.md).

## Support

For issues or questions about this plugin:
- Create an issue in the [marchat repository](https://github.com/Cod-e-Codes/marchat/issues)
- Check the [plugin documentation](https://github.com/Cod-e-Codes/marchat/blob/main/PLUGIN_ECOSYSTEM.md)
- Review the [plugin examples](https://github.com/Cod-e-Codes/marchat/tree/main/plugin/examples)

---

**Note**: This plugin is primarily intended for testing and demonstration purposes. For production use, consider developing more feature-rich plugins. 