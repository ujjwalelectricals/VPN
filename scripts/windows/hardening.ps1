param(
  [ValidateSet('check')]
  [string]$Mode = 'check'
)

$ErrorActionPreference = 'Stop'
Write-Host "FreeVPN local privacy checklist" -ForegroundColor Cyan

$checks = @()
$checks += [pscustomobject]@{Name='WireGuard installed'; Pass=(Test-Path "$env:ProgramFiles\WireGuard\wireguard.exe")}
$checks += [pscustomobject]@{Name='Windows Firewall service'; Pass=((Get-Service MpsSvc -ErrorAction SilentlyContinue).Status -eq 'Running')}
$checks += [pscustomobject]@{Name='IPv6 enabled (review required)'; Pass=$true; Note='Ensure your active tunnel routes IPv6 (::/0) before relying on the VPN for IPv6 privacy.'}
$checks | Format-Table -AutoSize

Write-Host "No firewall rules are changed by this script. Use the generated WireGuard configuration's full-tunnel AllowedIPs plus Windows Firewall controls deliberately." -ForegroundColor Yellow
