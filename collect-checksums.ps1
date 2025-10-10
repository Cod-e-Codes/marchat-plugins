param(
    [string]$Version = "1.0.0"
)

$ErrorActionPreference = "Stop"

$Plugins = @("weather", "githooks")

Write-Host "Collecting checksums for version $Version" -ForegroundColor Cyan
Write-Host ""

foreach ($plugin in $Plugins) {
    Write-Host "=== $plugin ===" -ForegroundColor Yellow
    $checksumFiles = Get-ChildItem -Path "plugins/$plugin/dist" -Filter "*-v$Version-*.zip.sha256" -ErrorAction SilentlyContinue
    
    foreach ($file in $checksumFiles) {
        $zipName = $file.BaseName
        $content = Get-Content $file.FullName
        $hash = $content.Split(' ')[0]
        
        Write-Host $zipName -ForegroundColor Gray
        Write-Host "  sha256:$hash" -ForegroundColor Green
        Write-Host ""
    }
}

Write-Host "Copy these checksums into registry.json" -ForegroundColor Cyan

