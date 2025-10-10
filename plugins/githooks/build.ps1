param(
    [string]$Version = "1.0.0",
    [ValidateSet("windows", "linux", "darwin")]
    [string]$Os = $(if ($env:GOOS) { $env:GOOS } elseif ($IsWindows) { "windows" } else { "linux" }),
    [ValidateSet("amd64", "arm64", "386")]
    [string]$Arch = $(if ($env:GOARCH) { $env:GOARCH } else { "amd64" }),
    [string]$OutDir = "dist"
)

$ErrorActionPreference = "Stop"

Write-Host "Building githooks plugin for $Os/$Arch v$Version" -ForegroundColor Cyan

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    throw "Go toolchain not found. Please install Go and ensure 'go' is on PATH."
}

# Resolve paths
$PluginDir = $PSScriptRoot
$DistDir = Join-Path $PluginDir $OutDir
New-Item -ItemType Directory -Force -Path $DistDir | Out-Null

# Build output paths
$ext = $(if ($Os -eq "windows") { ".exe" } else { "" })
$binName = "githooks$ext"
$stageDir = Join-Path $DistDir "stage-$Os-$Arch"
if (Test-Path $stageDir) { Remove-Item -Recurse -Force $stageDir }
New-Item -ItemType Directory -Force -Path $stageDir | Out-Null

# Set target env
$env:GOOS = $Os
$env:GOARCH = $Arch

# Build
$binPath = Join-Path $stageDir $binName
Push-Location $PluginDir
try {
    Write-Host "Running: GOOS=$Os GOARCH=$Arch go build -trimpath -ldflags '-s -w' -o $binPath ./..." -ForegroundColor DarkGray
    go build -trimpath -ldflags '-s -w' -o $binPath .
}
finally {
    Pop-Location
}

if (-not (Test-Path $binPath)) {
    throw "Build failed; binary not found at $binPath"
}

# Create plugin.json
$pluginJson = @{
    name = "githooks"
    version = $Version
    description = "Git repository management and status updates"
    author = "Cod-e-Codes"
    license = "MIT"
    commands = @(
        @{ name = "git-status"; description = "Show git status of current directory"; usage = ":git-status [path]"; admin_only = $false },
        @{ name = "git-log"; description = "Show recent git commits"; usage = ":git-log [n] [path]"; admin_only = $false },
        @{ name = "git-branch"; description = "Show current git branch and available branches"; usage = ":git-branch [path]"; admin_only = $false },
        @{ name = "git-diff"; description = "Show git diff of uncommitted changes"; usage = ":git-diff [path]"; admin_only = $false },
        @{ name = "git-watch"; description = "Watch a repository for changes (admin only)"; usage = ":git-watch <path>"; admin_only = $true }
    )
} | ConvertTo-Json -Depth 5

Set-Content -Path (Join-Path $stageDir "plugin.json") -Value $pluginJson -Encoding UTF8

# Include README
Copy-Item (Join-Path $PluginDir "README.md") -Destination (Join-Path $stageDir "README.md") -Force

# Zip archive
$zipName = "githooks-plugin-v$Version-$Os-$Arch.zip"
$zipPath = Join-Path $DistDir $zipName
if (Test-Path $zipPath) { Remove-Item $zipPath -Force }

Write-Host "Creating archive $zipName" -ForegroundColor Cyan
Compress-Archive -Path (Join-Path $stageDir '*') -DestinationPath $zipPath -Force

# Compute checksum
Write-Host "Computing SHA-256 checksum" -ForegroundColor DarkGray
$hash = (Get-FileHash -Algorithm SHA256 -Path $zipPath).Hash.ToLower()
Set-Content -Path ($zipPath + ".sha256") -Value "$hash  $zipName" -Encoding ASCII

Write-Host "Done: $zipPath" -ForegroundColor Green
Write-Host "Checksum: $($zipPath).sha256" -ForegroundColor Green

