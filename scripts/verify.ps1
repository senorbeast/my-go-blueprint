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

$fixture = Join-Path ([System.IO.Path]::GetTempPath()) ("go-blueprint-auth-" + [guid]::NewGuid().ToString("N"))
try {
    & $goPath run . create --name verify-auth --module example.com/verify-auth --database postgres --seed minimal --no-frontend --feature auth --output $fixture
    if ($LASTEXITCODE -ne 0) {
        exit $LASTEXITCODE
    }
    Push-Location (Join-Path $fixture "backend")
    try {
        & $goPath mod tidy
        if ($LASTEXITCODE -ne 0) {
            exit $LASTEXITCODE
        }
        & $goPath tool sqlc generate
        if ($LASTEXITCODE -ne 0) {
            exit $LASTEXITCODE
        }
        & $goPath test ./...
        if ($LASTEXITCODE -ne 0) {
            exit $LASTEXITCODE
        }
    } finally {
        Pop-Location
    }
} finally {
    if (Test-Path -LiteralPath $fixture) {
        Remove-Item -LiteralPath $fixture -Recurse -Force
    }
}
