#!/usr/bin/env bash
set -euo pipefail

# Minimal Ubuntu WireGuard exit-node bootstrap.
# Run on a server you own/control. This script does not create accounts or pay for anything.
# It expects WG_SERVER_ADDRESS, WG_CLIENT_PUBLIC_KEY and WG_CLIENT_ADDRESS to be set.

: "${WG_SERVER_ADDRESS:?Set WG_SERVER_ADDRESS, e.g. 10.8.0.1/24}"
: "${WG_CLIENT_PUBLIC_KEY:?Set WG_CLIENT_PUBLIC_KEY to your client's public key}"
: "${WG_CLIENT_ADDRESS:?Set WG_CLIENT_ADDRESS, e.g. 10.8.0.2/32}"
: "${WG_EXTERNAL_INTERFACE:=eth0}"
: "${WG_LISTEN_PORT:=51820}"

sudo apt-get update
sudo apt-get install -y wireguard iptables

sudo install -d -m 700 /etc/wireguard
if [[ ! -f /etc/wireguard/server.key ]]; then
  umask 077
  wg genkey | sudo tee /etc/wireguard/server.key >/dev/null
  sudo cat /etc/wireguard/server.key | wg pubkey | sudo tee /etc/wireguard/server.pub >/dev/null
fi
SERVER_PRIVATE_KEY=$(sudo cat /etc/wireguard/server.key)

cat <<EOF | sudo tee /etc/wireguard/wg0.conf >/dev/null
[Interface]
Address = ${WG_SERVER_ADDRESS}
ListenPort = ${WG_LISTEN_PORT}
PrivateKey = ${SERVER_PRIVATE_KEY}
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -A FORWARD -o wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o ${WG_EXTERNAL_INTERFACE} -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -D FORWARD -o wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o ${WG_EXTERNAL_INTERFACE} -j MASQUERADE

[Peer]
PublicKey = ${WG_CLIENT_PUBLIC_KEY}
AllowedIPs = ${WG_CLIENT_ADDRESS}
EOF

sudo sysctl -w net.ipv4.ip_forward=1 >/dev/null
printf 'net.ipv4.ip_forward=1\n' | sudo tee /etc/sysctl.d/99-freevpn.conf >/dev/null
sudo systemctl enable --now wg-quick@wg0
sudo wg show

echo "WireGuard server ready. Open UDP ${WG_LISTEN_PORT} in your cloud/router firewall." 
echo "Server public key: $(sudo cat /etc/wireguard/server.pub)"
