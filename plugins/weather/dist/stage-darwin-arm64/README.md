# Weather Plugin

A marchat plugin that provides weather information and forecasts using the wttr.in service.

## Description

The Weather Plugin allows users to get current weather conditions and 3-day forecasts for any location worldwide. It uses the free wttr.in API to provide accurate weather data.

## Features

- **Current Weather**: Get real-time weather conditions for any location
- **3-Day Forecast**: View extended weather forecasts
- **Message Processing**: Automatically responds to "weather:" messages
- **Global Coverage**: Works with any location worldwide
- **Rich Information**: Temperature, humidity, wind speed, and conditions

## Commands

### User Commands
- `:weather [location]` — Get current weather for a location (defaults to Charlotte, NC)
- `:forecast [location]` — Get 3-day weather forecast for a location

### Message Processing
- `weather: <location>` — Automatically gets weather for the specified location

## Installation

### From Plugin Store
```bash
:store
# Navigate to weather plugin and press Enter to install
```

### Direct Installation
```bash
:plugin install weather
```

## Build and Package

Builds produce a single archive targeted to your OS/architecture.

### PowerShell (Windows)
```powershell
# From this directory
./build.ps1 -Version 1.0.0 -Os windows -Arch amd64
# Output: dist/weather-plugin-v1.0.0-windows-amd64.zip and .sha256
```

### Bash (Linux/macOS)
```bash
# From this directory
VERSION=1.0.0 bash build.sh # uses GOOS/GOARCH from env
# or override
VERSION=1.0.0 GOOS=linux GOARCH=amd64 bash build.sh
VERSION=1.0.0 GOOS=darwin GOARCH=amd64 bash build.sh
VERSION=1.0.0 GOOS=darwin GOARCH=arm64 bash build.sh
# Output: dist/weather-plugin-v1.0.0-{os}-{arch}.zip and .sha256
```

### Archive Contents
- `weather` or `weather.exe` — the plugin binary
- `plugin.json` — plugin manifest with metadata and commands
- `README.md` — this documentation

### Naming Convention
- `weather-plugin-v<version>-<goos>-<goarch>.zip`
  - Examples: `weather-plugin-v1.0.0-windows-amd64.zip`, `weather-plugin-v1.0.0-darwin-arm64.zip`

## Usage Examples

### Get Current Weather
```bash
:weather
# Output: Weather for Charlotte, United States (default location)

:weather London
# Output: Weather for London, United Kingdom

:weather Tokyo
# Output: Weather for Tokyo, Japan
```

### Get Weather Forecast
```bash
:forecast New York
# Output: 3-day forecast for New York, United States
```

### Message-Based Weather
```bash
weather: Paris
# Output: Weather for Paris, France
```

## Plugin Management

### List Installed Plugins
```bash
:plugin list
```

### Enable/Disable Plugin
```bash
:plugin enable weather
:plugin disable weather
```

### Uninstall Plugin (Admin Only)
```bash
:plugin uninstall weather
```

## Technical Details

### Plugin Structure
```
weather/
├── weather (.exe)   # Binary executable
├── plugin.json      # Plugin manifest
└── README.md        # This documentation
```

### Weather Information Provided
- Location (area name and country)
- Current conditions
- Temperature (°C)
- Feels like temperature
- Humidity percentage
- Wind speed (km/h)
- 3-day forecast (when requested)

### API Service
This plugin uses [wttr.in](https://wttr.in), a free weather service that provides weather data in JSON format.

## Version History

- **v1.0.0** - Initial release with weather and forecast commands

## License

MIT License - see [LICENSE](../../LICENSE) for details.

## Contributing

This plugin is part of the marchat-plugins project. For plugin development guidelines, see the main repository's [CONTRIBUTING.md](../../CONTRIBUTING.md).

## Support

For issues or questions about this plugin:
- Create an issue in the [marchat-plugins repository](https://github.com/Cod-e-Codes/marchat-plugins/issues)
- Check the [plugin documentation](https://github.com/Cod-e-Codes/marchat/blob/main/PLUGIN_ECOSYSTEM.md)

---

**Note**: Weather data is provided by wttr.in. The plugin requires an active internet connection to fetch weather information.

