# Git Hooks Plugin

A marchat plugin that provides git repository management and status updates directly in your chat.

## Description

The Git Hooks Plugin allows users to interact with git repositories without leaving marchat. It provides commands for viewing git status, commit history, branches, and diffs, making it easy to monitor repository changes during conversations.

## Features

- **Git Status**: View the current status of your git repository
- **Commit History**: Browse recent commits with author and time information
- **Branch Management**: View current branch and available branches
- **Diff Viewer**: See uncommitted changes at a glance
- **Repository Watching**: Monitor repositories for changes (admin only)

## Commands

### User Commands
- `:git-status [path]` — Show git status of current or specified directory
- `:git-log [n] [path]` — Show recent git commits (default: 5)
- `:git-branch [path]` — Show current branch and available branches
- `:git-diff [path]` — Show git diff of uncommitted changes

### Admin Commands
- `:git-watch <path>` — Watch a repository for changes (admin only)

## Installation

### From Plugin Store
```bash
:store
# Navigate to githooks plugin and press Enter to install
```

### Direct Installation
```bash
:plugin install githooks
```

## Build and Package

Builds produce a single archive targeted to your OS/architecture.

### PowerShell (Windows)
```powershell
# From this directory
./build.ps1 -Version 1.0.0 -Os windows -Arch amd64
# Output: dist/githooks-plugin-v1.0.0-windows-amd64.zip and .sha256
```

### Bash (Linux/macOS)
```bash
# From this directory
VERSION=1.0.0 bash build.sh # uses GOOS/GOARCH from env
# or override
VERSION=1.0.0 GOOS=linux GOARCH=amd64 bash build.sh
VERSION=1.0.0 GOOS=darwin GOARCH=amd64 bash build.sh
VERSION=1.0.0 GOOS=darwin GOARCH=arm64 bash build.sh
# Output: dist/githooks-plugin-v1.0.0-{os}-{arch}.zip and .sha256
```

### Archive Contents
- `githooks` or `githooks.exe` — the plugin binary
- `plugin.json` — plugin manifest with metadata and commands
- `README.md` — this documentation

### Naming Convention
- `githooks-plugin-v<version>-<goos>-<goarch>.zip`
  - Examples: `githooks-plugin-v1.0.0-windows-amd64.zip`, `githooks-plugin-v1.0.0-darwin-arm64.zip`

## Usage Examples

### Check Git Status
```bash
:git-status
# Output: Shows status of current directory

:git-status /path/to/repo
# Output: Shows status of specified repository
```

### View Recent Commits
```bash
:git-log
# Output: Shows last 5 commits

:git-log 10
# Output: Shows last 10 commits

:git-log 10 /path/to/repo
# Output: Shows last 10 commits for specified repo
```

### View Branches
```bash
:git-branch
# Output: Shows current branch and all branches

:git-branch /path/to/repo
# Output: Shows branches for specified repository
```

### View Uncommitted Changes
```bash
:git-diff
# Output: Shows diff statistics for uncommitted changes

:git-diff /path/to/repo
# Output: Shows diff for specified repository
```

### Watch Repository (Admin Only)
```bash
:git-watch /path/to/repo
# Output: Now watching repository: /path/to/repo
```

## Plugin Management

### List Installed Plugins
```bash
:plugin list
```

### Enable/Disable Plugin
```bash
:plugin enable githooks
:plugin disable githooks
```

### Uninstall Plugin (Admin Only)
```bash
:plugin uninstall githooks
```

## Technical Details

### Plugin Structure
```
githooks/
├── githooks (.exe)  # Binary executable
├── plugin.json      # Plugin manifest
└── README.md        # This documentation
```

### Requirements
- Git must be installed and available in PATH
- Plugin will verify git availability during initialization
- Repository paths must contain a valid `.git` directory

### Git Commands Used
- `git status --short --branch` - Repository status
- `git log --pretty=format` - Commit history
- `git branch --show-current` - Current branch
- `git branch -a` - All branches
- `git diff --stat` - Change statistics

## Version History

- **v1.0.0** - Initial release with git status, log, branch, diff, and watch commands

## License

MIT License - see [LICENSE](../../LICENSE) for details.

## Contributing

This plugin is part of the marchat-plugins project. For plugin development guidelines, see the main repository's [CONTRIBUTING.md](../../CONTRIBUTING.md).

## Support

For issues or questions about this plugin:
- Create an issue in the [marchat-plugins repository](https://github.com/Cod-e-Codes/marchat-plugins/issues)
- Check the [plugin documentation](https://github.com/Cod-e-Codes/marchat/blob/main/PLUGIN_ECOSYSTEM.md)

## Security Notes

- Admin-only commands (git-watch) require proper authorization
- Plugin validates repository paths and git directory existence
- All git commands run with user's current permissions
- No git write operations are performed (read-only)

---

**Note**: This plugin executes git commands on your system. Ensure you trust the repositories you're querying and have proper access permissions.

