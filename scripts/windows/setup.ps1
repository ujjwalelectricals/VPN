$ErrorActionPreference = 'Stop'
Write-Host "Ujjwal FreeVPN setup" -ForegroundColor Cyan
Write-Host "This installs the free WireGuard application if WinGet can provide it, then checks the local client."

if (-not (Get-Command winget -ErrorAction SilentlyContinue)) {
  Write-Warning "WinGet is not available. Install WireGuard for Windows from the official WireGuard site, then rerun this script."
  exit 1
}

winget install --id WireGuard.WireGuard --exact --accept-package-agreements --accept-source-agreements

$paths = @(
  "$env:ProgramFiles\WireGuard\wireguard.exe",
  "${env:ProgramFiles(x86)}\WireGuard\wireguard.exe"
)
$found = $paths | Where-Object { Test-Path $_ } | Select-Object -First 1
if ($found) {
  Write-Host "WireGuard found at $found" -ForegroundColor Green
} else {
  Write-Warning "WireGuard was not detected in the standard install locations."
}

Write-Host "Next: create profiles from profiles/*.conf.example, put real values in a local .conf file, and keep private keys out of Git."
