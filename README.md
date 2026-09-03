# Ujjwal FreeVPN

A lightweight, privacy-focused VPN control project for Windows built around the free and open-source WireGuard protocol.

## What this project provides

- Windows-first WireGuard tunnel control in Go.
- Country/location profiles you can populate with your own servers.
- Local-only dashboard on `127.0.0.1:7070`.
- Node health/status and a `doctor` diagnostic command.
- Full-tunnel profile template with IPv4 + IPv6 routing and DNS settings.
- Conservative privacy hardening checks without silently changing firewall rules.
- No paid API, subscription, proprietary VPN backend, or telemetry requirement.

## Important: what “100% free” means

The **software in this repository is free to use and modify**. A VPN service still needs an exit server. This repository does not provide free worldwide exit servers because those require infrastructure somewhere in the requested country.

For a genuinely $0 setup, you can run your own WireGuard server on hardware you already own, or use an Always Free cloud VM when its limits and region availability fit your needs. Oracle currently lists Always Free compute, including Ampere A1 resources, but those resources are limited to the tenancy's home region and subject to capacity/usage limits. citeturn919278search0turn919278search2

That means this project can be **$0 software + $0 server** for a small self-hosted setup, but it cannot honestly promise unlimited countries for free.

## Architecture

```text
                 Windows PC
                     |
                 Ujjwal VPN
                     |
                WireGuard tunnel
                     |
           +---------+---------+
           |                   |
       India node          USA node
           |                   |
       Internet           Internet
```

Each location is just a WireGuard peer configuration. To add a country, create a server there, export a client config, place it in `profiles/`, and enable the matching entry in `nodes.json`.

WireGuard is the actual encrypted tunnel layer; this project is the lightweight management/control layer around it. citeturn564602search14

## Windows setup

1. Install WireGuard for Windows from the official WireGuard project.
2. Run `scripts/windows/setup.ps1` in an elevated PowerShell if WinGet is available.
3. Copy `profiles/india.conf.example` to `profiles/india.conf`.
4. Put your **real** WireGuard client keys, server public key, address, and endpoint into the local `.conf` file.
5. Never commit a real `PrivateKey` to GitHub.
6. Set the matching `nodes.json` entry to `"enabled": true`.

Then build and run:

```powershell
go test ./...
go build -o vpnctl.exe ./cmd/vpnctl
.\vpnctl.exe doctor
.\vpnctl.exe list
.\vpnctl.exe connect india
```

WireGuard's Windows application and tooling are available from the official project. citeturn564602search12

## Commands

```text
vpnctl list
vpnctl connect <node-id>
vpnctl disconnect <node-id>
vpnctl status <node-id>
vpnctl doctor
vpnctl dashboard
vpnctl version
```

The dashboard is intentionally bound to `127.0.0.1`; it is not a public remote-control service.

## Privacy model

A VPN changes the network path and can hide your normal public IP from destination websites, but it does not guarantee anonymity. Accounts, cookies, browser fingerprinting, device identifiers, GPS, and the VPN exit server itself can still associate activity with you.

## Security principles

- Never ship real client private keys in this public repository.
- Keep the control dashboard local-only.
- Do not provide an unrestricted arbitrary shell executor.
- Do not silently modify Windows Firewall rules.
- Prefer explicit user actions for privileged tunnel operations.
- Keep country profiles disabled until their configuration is valid.

## Roadmap

- [x] WireGuard profile catalog
- [x] Windows tunnel service control
- [x] Local dashboard
- [x] Diagnostics
- [x] Full-tunnel configuration template
- [ ] DNS leak test endpoint integration
- [ ] Connection latency/packet-loss scoring
- [ ] System-tray GUI
- [ ] Automatic reconnect policy
- [ ] Optional multi-hop profile orchestration
- [ ] Signed Windows release builds
