# Ujjwal FreeVPN

A lightweight, privacy-focused VPN control project for Windows built around the free and open-source WireGuard protocol.

## What this project provides

- Windows-first WireGuard tunnel control in Go.
- Country/location profiles you can populate with your own servers.
- Local-only dashboard on `127.0.0.1:7070`.
- Node health/status and a `doctor` diagnostic command.
- Full-tunnel profile template with IPv4 + IPv6 routing and DNS settings.
- Conservative privacy hardening checks without silently changing firewall rules.
- Free cross-platform release builds through GitHub Actions.
- No paid API, subscription, proprietary VPN backend, or telemetry requirement.

## Important: what “100% free” means

The **software in this repository is free to use and modify**. A VPN service still needs an exit server. This repository does not provide free worldwide exit servers because those require infrastructure somewhere in the requested country.

For a genuinely $0 setup, you can run your own WireGuard server on hardware you already own, or use an Always Free cloud VM when its limits and region availability fit your needs. Oracle currently documents Always Free compute resources, including Ampere A1, with tenancy/home-region and capacity limits. See the official Oracle Free Tier documentation before deploying.

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

Each location is a WireGuard peer configuration. To add a country, create a server there, export a client config, place it in `profiles/`, and enable the matching entry in `nodes.json`.

WireGuard is the encrypted tunnel layer; this project is the lightweight management/control layer around it.

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
go build -trimpath -ldflags='-s -w' -o vpnctl.exe ./cmd/vpnctl
.\vpnctl.exe doctor
.\vpnctl.exe list
.\vpnctl.exe connect india
```

## Local dashboard

Run:

```powershell
.\vpnctl.exe dashboard
```

Then open:

`http://127.0.0.1:7070`

The dashboard is intentionally local-only. Connect/disconnect controls require the process to have enough Windows privilege to manage WireGuard tunnel services.

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

## Self-hosted Ubuntu exit node

The repository includes `server/ubuntu-wireguard-setup.sh` for a server you own or are authorized to administer. It installs WireGuard, enables IPv4 forwarding, configures NAT, and creates a peer entry. You still need to open UDP 51820 in the server/cloud firewall and put the server public key and endpoint into the client profile.

## Privacy model

A VPN changes the network path and can hide your normal public IP from destination websites, but it does not guarantee anonymity. Accounts, cookies, browser fingerprinting, device identifiers, GPS, and the VPN exit server itself can still associate activity with you.

## Security principles

- Never ship real client private keys in this public repository.
- Keep the control dashboard local-only.
- Do not provide an unrestricted arbitrary shell executor.
- Do not silently modify Windows Firewall rules.
- Prefer explicit user actions for privileged tunnel operations.
- Keep country profiles disabled until their configuration is valid.

## Project status

The current repository is a **real client/control layer**, not a public commercial VPN network. It is ready for you to attach your own nodes. Country switching works by selecting between your configured WireGuard exit servers.

## Roadmap

- [x] WireGuard profile catalog
- [x] Windows tunnel service control
- [x] Local dashboard
- [x] Diagnostics
- [x] Full-tunnel configuration template
- [x] Windows/Linux CI
- [x] Cross-platform release artifacts
- [ ] DNS leak test endpoint integration
- [ ] Connection latency/packet-loss scoring
- [ ] System-tray GUI
- [ ] Automatic reconnect policy
- [ ] Optional multi-hop profile orchestration
- [ ] Signed Windows release binaries
