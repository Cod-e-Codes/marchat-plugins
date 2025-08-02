# CONTRIBUTING.md

## Contributing Plugins to marchat

Thank you for your interest in contributing plugins to marchat! This guide will help you create and submit plugins to the community registry.

## Plugin Requirements

### Structure
Your plugin should have this structure:
```
myplugin/
├── plugin.json     # Plugin manifest (required)
├── myplugin        # Binary executable (required)
└── README.md       # Documentation (recommended)
```

### Plugin Manifest (`plugin.json`)
```json
{
  "name": "myplugin",
  "version": "1.0.0",
  "description": "Brief description of what your plugin does",
  "author": "Your Name",
  "homepage": "https://github.com/yourusername/yourrepo",
  "commands": [
    {
      "name": "mycommand",
      "description": "What this command does",
      "usage": ":mycommand <args>",
      "admin_only": false
    }
  ],
  "permissions": [],
  "settings": {},
  "min_version": "0.2.0-beta.1"
}
```

### Binary Requirements
- **Executable**: Must be a valid Go binary
- **Cross-platform**: Build for Linux, Windows, and macOS
- **Self-contained**: No external dependencies beyond the marchat SDK
- **Size**: Keep under 10MB for reasonable download times

## Development

### 1. Use the marchat Plugin SDK
```go
package main

import (
    "github.com/Cod-e-Codes/marchat/plugin/sdk"
    "time"
)

type MyPlugin struct {
    *sdk.BasePlugin
}

func (p *MyPlugin) OnMessage(msg sdk.Message) ([]sdk.Message, error) {
    // Handle incoming messages
    return nil, nil
}

func (p *MyPlugin) Commands() []sdk.PluginCommand {
    return []sdk.PluginCommand{
        {
            Name:        "mycommand",
            Description: "What this command does",
            Usage:       ":mycommand <args>",
            AdminOnly:   false,
        },
    }
}

func main() {
    plugin := &MyPlugin{
        BasePlugin: sdk.NewBasePlugin("MyPlugin", "1.0.0"),
    }
    sdk.RunPlugin(plugin)
}
```

### 2. Build Your Plugin
```bash
# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o myplugin-linux-amd64
GOOS=windows GOARCH=amd64 go build -o myplugin-windows-amd64.exe
GOOS=darwin GOARCH=amd64 go build -o myplugin-darwin-amd64
```

### 3. Create Distribution Package
```bash
# Create ZIP file with binary and manifest
zip myplugin-plugin.zip myplugin plugin.json
```

## Submitting Your Plugin

### 1. Fork the Repository
Fork [marchat-plugins](https://github.com/Cod-e-Codes/marchat-plugins) to your GitHub account.

### 2. Add Your Plugin
- Add your plugin ZIP file to the `plugins/` directory
- Update `registry.json` with your plugin information
- Include a README.md for your plugin

### 3. Update Registry
Add your plugin to `registry.json`:
```json
{
  "name": "myplugin",
  "version": "1.0.0",
  "description": "Brief description of what your plugin does",
  "author": "Your Name",
  "homepage": "https://github.com/yourusername/yourrepo",
  "download_url": "https://github.com/Cod-e-Codes/marchat-plugins/raw/main/plugins/myplugin/myplugin-plugin.zip",
  "checksum": "sha256:YOUR_CHECKSUM_HERE",
  "category": "utility",
  "tags": ["chat", "utility", "your-tags"],
  "min_version": "0.2.0-beta.1"
}
```

### 4. Calculate Checksum
```bash
# On Linux/macOS
sha256sum myplugin-plugin.zip

# On Windows
Get-FileHash myplugin-plugin.zip -Algorithm SHA256
```

### 5. Create Pull Request
- Create a new branch for your plugin
- Commit your changes
- Create a pull request with a clear description

## Plugin Guidelines

### Security
- **No malicious code**: Plugins must not harm users or systems
- **Safe file operations**: Only write to designated plugin directories
- **Input validation**: Validate all user input
- **Resource limits**: Don't consume excessive CPU/memory

### Quality
- **Documentation**: Include clear README and usage examples
- **Error handling**: Graceful error handling and user feedback
- **Testing**: Test your plugin thoroughly before submission
- **Performance**: Efficient resource usage

### Naming
- **Unique names**: Ensure your plugin name is unique
- **Descriptive**: Use clear, descriptive names
- **No conflicts**: Avoid names that conflict with existing plugins

## Categories

Use appropriate categories for your plugin:
- **utility**: General utility plugins
- **chat**: Chat enhancement plugins
- **admin**: Administrative tools
- **fun**: Entertainment plugins
- **productivity**: Productivity tools
- **integration**: Third-party service integrations

## Review Process

1. **Automated checks**: Basic validation of plugin structure
2. **Security review**: Manual review for security issues
3. **Functionality test**: Testing plugin functionality
4. **Documentation review**: Ensuring adequate documentation

## Support

- **Issues**: Use GitHub Issues for bug reports
- **Discussions**: Use GitHub Discussions for questions
- **Documentation**: Check the main marchat documentation

## License

By contributing plugins, you agree to license your plugin under the MIT License, the same license as marchat.

---

Thank you for contributing to the marchat plugin ecosystem!
