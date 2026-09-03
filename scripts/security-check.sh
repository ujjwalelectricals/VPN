#!/usr/bin/env bash
set -euo pipefail

# Fail if a non-example WireGuard config contains a likely real PrivateKey.
matches=$(git grep -n -I -E '^\s*PrivateKey\s*=' -- ':!profiles/*.conf.example' ':!README.md' || true)
if [[ -n "$matches" ]]; then
  echo "Potential private key material found:" >&2
  echo "$matches" >&2
  exit 1
fi

echo "Secret scan passed: no non-example WireGuard PrivateKey lines found."
