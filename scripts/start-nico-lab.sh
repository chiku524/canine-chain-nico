#!/usr/bin/env bash
# Start Nico Lab — click-to-test web UI for jackal-nico-1.
#
# Prerequisites: canined running, nico-provider on :3333 (optional but needed for uploads).
#
# Usage (WSL):
#   ./scripts/start-nico-lab.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="/usr/local/go/bin:$HOME/go/bin:${PATH}"

HOME_DIR="${HOME_DIR:-$HOME/.canine-nico}"
LISTEN="${LISTEN:-:3456}"
CHAIN_ID="${CHAIN_ID:-jackal-nico-1}"
PROVIDER="${PROVIDER:-http://127.0.0.1:3333}"

cd "$ROOT/tools/nico-lab"
WSL_IP="$(hostname -I 2>/dev/null | awk '{print $1}')"
echo "Nico Lab → http://127.0.0.1${LISTEN}"
if [[ -n "$WSL_IP" ]]; then
  echo "  WSL/Windows browser fallback → http://${WSL_IP}${LISTEN}"
fi
echo "  home=$HOME_DIR chain=$CHAIN_ID provider=$PROVIDER"
exec go run . \
  -listen "$LISTEN" \
  -home "$HOME_DIR" \
  -chain-id "$CHAIN_ID" \
  -provider "$PROVIDER"
