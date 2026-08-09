#!/usr/bin/env bash
# Register and start a local nico-provider against jackal-nico-1.
#
# Prerequisites: canined running (init-nico-testnet.sh), Go on PATH.
#
# Usage (WSL):
#   ./scripts/start-nico-provider.sh           # init on-chain + start HTTP :3333
#   INIT_ONLY=1 ./scripts/start-nico-provider.sh
#   START_ONLY=1 ./scripts/start-nico-provider.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CHAIN_ID="${CHAIN_ID:-jackal-nico-1}"
HOME_DIR="${HOME_DIR:-$HOME/.canine-nico}"
KEYRING="${KEYRING:-test}"
CANINED="${CANINED:-canined}"
PROVIDER_KEY="${PROVIDER_KEY:-provider}"
FUNDER_KEY="${FUNDER_KEY:-validator}"
PROVIDER_IP="${PROVIDER_IP:-http://127.0.0.1:3333}"
TOTAL_SPACE="${TOTAL_SPACE:-1092616192000}" # ~1 TB
COLLATERAL_FUND="${COLLATERAL_FUND:-20000000000ujkl}" # 20k JKL (10k collateral + gas)
REST_API="${REST_API:-http://127.0.0.1:1317}"
LISTEN="${LISTEN:-:3333}"

if ! command -v "$CANINED" >/dev/null 2>&1; then
  if [[ -x "$HOME/go/bin/canined" ]]; then
    CANINED="$HOME/go/bin/canined"
  else
    echo "canined not found" >&2
    exit 1
  fi
fi

export PATH="/usr/local/go/bin:$HOME/go/bin:${PATH}"

tx() {
  "$CANINED" "$@" --home "$HOME_DIR" --chain-id "$CHAIN_ID" \
    --keyring-backend "$KEYRING" --gas auto --gas-adjustment 1.5 \
    --fees 5000ujkl -y
}

q() {
  "$CANINED" query "$@" --home "$HOME_DIR" --chain-id "$CHAIN_ID" --output json
}

if [[ "${START_ONLY:-0}" != "1" ]]; then
  echo "== Ensuring provider key =="
  if ! "$CANINED" keys show "$PROVIDER_KEY" --home "$HOME_DIR" --keyring-backend "$KEYRING" >/dev/null 2>&1; then
    "$CANINED" keys add "$PROVIDER_KEY" --home "$HOME_DIR" --keyring-backend "$KEYRING"
  fi
  PROVIDER_ADDR="$("$CANINED" keys show "$PROVIDER_KEY" -a --home "$HOME_DIR" --keyring-backend "$KEYRING")"
  echo "Provider: $PROVIDER_ADDR"

  echo "== Funding provider ($COLLATERAL_FUND) =="
  BAL="$(q bank balances "$PROVIDER_ADDR" 2>/dev/null | jq -r '.balances[]? | select(.denom=="ujkl") | .amount' || true)"
  BAL="${BAL:-0}"
  NEED=10000000000
  if [[ "$BAL" -lt "$NEED" ]]; then
    tx tx bank send "$FUNDER_KEY" "$PROVIDER_ADDR" "$COLLATERAL_FUND"
    sleep 3
  else
    echo "Provider already funded ($BAL ujkl)"
  fi

  echo "== Init provider on-chain (ip=$PROVIDER_IP) =="
  if q storage show-providers "$PROVIDER_ADDR" 2>/dev/null | jq -e '.providers.address' >/dev/null 2>&1; then
    echo "Provider already registered — updating IP"
    tx provider set-ip "$PROVIDER_IP" --from "$PROVIDER_KEY" || true
  else
    tx provider init "$PROVIDER_IP" "$TOTAL_SPACE" "" --from "$PROVIDER_KEY"
  fi
  sleep 2
  q storage show-providers "$PROVIDER_ADDR" | jq .
fi

if [[ "${INIT_ONLY:-0}" == "1" ]]; then
  echo "INIT_ONLY=1 — skipping HTTP start"
  exit 0
fi

PROVIDER_ADDR="$("$CANINED" keys show "$PROVIDER_KEY" -a --home "$HOME_DIR" --keyring-backend "$KEYRING")"

echo "== Starting nico-provider on $LISTEN =="
cd "$ROOT/tools/nico-provider"
exec go run . \
  -address "$PROVIDER_ADDR" \
  -listen "$LISTEN" \
  -rest "$REST_API" \
  -chain-id "$CHAIN_ID"
