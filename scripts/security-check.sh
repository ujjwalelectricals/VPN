#!/usr/bin/env bash
set -euo pipefail

# Fail when a repository file appears to contain an actual WireGuard private key.
# WireGuard keys are base64-encoded and typically 44 characters long.
matches=$(git grep -n -I -E '^\s*PrivateKey\s*=\s*[A-Za-z0-9+/]{43}={0,1}\s*$' -- ':!profiles/*.conf.example' || true)
if [[ -n "$matches" ]]; then
  echo "Potential real WireGuard private key found:" >&2
  echo "$matches" >&2
  exit 1
fi

echo "Secret scan passed: no real WireGuard PrivateKey values found."
