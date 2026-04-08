param(
    [string]$Version = "1.0.0"
)

$ErrorActionPreference = "Stop"

$Plugins = @("echo", "weather", "githooks")
$Platforms = @(
    @{Os = "windows"; Arch = "amd64"},
    @{Os = "linux"; Arch = "amd64"},
    @{Os = "linux"; Arch = "arm64"},
    @{Os = "darwin"; Arch = "amd64"},
    @{Os = "darwin"; Arch = "arm64"}
)

Write-Host "Building all plugins for all platforms (v$Version)" -ForegroundColor Cyan
Write-Host ""

foreach ($plugin in $Plugins) {
    Write-Host "=== Building $plugin ===" -ForegroundColor Yellow
    Push-Location "plugins/$plugin"
    
    try {
        foreach ($platform in $Platforms) {
            $os = $platform.Os
            $arch = $platform.Arch
            Write-Host "  - $os/$arch" -ForegroundColor Gray
            & ./build.ps1 -Version $Version -Os $os -Arch $arch
        }
    }
    finally {
        Pop-Location
    }
    
    Write-Host ""
}

Write-Host "✅ All builds complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Create GitHub release: gh release create v$Version"
Write-Host "2. Upload all ZIP files from plugins/*/dist/"
Write-Host "3. Update registry.json with checksums"
Write-Host "4. Push changes"

