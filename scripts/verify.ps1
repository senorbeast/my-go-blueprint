[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"
$workspace = Split-Path -Parent $PSScriptRoot
$cache = Join-Path $workspace ".cache/go-build"
$env:GOCACHE = $cache

$go = Get-Command go -ErrorAction Stop
$fallback = Join-Path $HOME ".local/opt/go/bin/go.exe"

& $go.Source list std 1>$null 2>$null
if ($LASTEXITCODE -ne 0 -and (Test-Path -LiteralPath $fallback)) {
    $goPath = $fallback
} else {
    $goPath = $go.Source
}

Write-Host "Using Go: $goPath"
& $goPath test ./...
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}
