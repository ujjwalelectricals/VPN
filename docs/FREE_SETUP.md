# $0 setup plan

## Option A — use a computer you already own

1. Install WireGuard on the server computer.
2. Run `server/ubuntu-wireguard-setup.sh` on an Ubuntu server you administer.
3. Create a client key pair.
4. Add the client's public key to the server peer configuration.
5. Put the client private key, tunnel address, server public key, and endpoint into a local file such as `profiles/india.conf`.
6. Keep that `.conf` file out of GitHub; `.gitignore` already excludes `profiles/*.conf`.
7. Change the matching node in `nodes.json` to `enabled: true`.
8. Build `vpnctl` and run `vpnctl connect india` as Administrator on Windows.

## Option B — Always Free cloud

Oracle's current Free Tier documents Always Free compute resources, including Ampere A1. The Always Free compute allocation is limited and provisioning is tied to the tenancy's home region. Capacity can also be temporarily unavailable in a region. This can be useful for one or two self-hosted experiments, but it does not provide unlimited countries.

Do not put a credit-card secret, cloud secret, SSH private key, or VPN private key in this repository.

## Country switching

Country switching is not a magic API feature. Each country needs a reachable VPN exit server in that country. The client switches between those server profiles.

## Privacy checklist

- Use full-tunnel `AllowedIPs = 0.0.0.0/0, ::/0` when your server is prepared for IPv4 and IPv6 forwarding.
- Configure DNS in the tunnel profile.
- Verify IPv6 is actually routed before relying on it for privacy.
- Use a kill switch deliberately; test recovery before making it permanent.
- Remember that a VPN does not make you anonymous.
