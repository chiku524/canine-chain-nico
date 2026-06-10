#!/usr/bin/env bash
# Capture wasm code IDs and IBC channels from a live Jackal network via REST.
#
# Usage:
#   ./scripts/capture-chain-inventory.sh                    # mainnet
#   NETWORK=testnet ./scripts/capture-chain-inventory.sh
#   REST_API=https://jackal-api.polkachu.com ./scripts/capture-chain-inventory.sh
#
# Requires: curl, jq

set -euo pipefail

NETWORK="${NETWORK:-mainnet}"
DATE_STAMP="$(date -u +%Y%m%dT%H%MZ)"
OUT_DIR="${OUT_DIR:-docs/inventory}"
OUT_FILE="${OUT_FILE:-${OUT_DIR}/captured-${NETWORK}-${DATE_STAMP}.json}"

case "$NETWORK" in
  mainnet)
    CHAIN_ID="${CHAIN_ID:-jackal-1}"
    REST_API="${REST_API:-https://api.jackalprotocol.com}"
    ;;
  testnet)
    CHAIN_ID="${CHAIN_ID:-jackal-testnet-1}"
    REST_API="${REST_API:-https://testnet-api.jackalprotocol.com}"
    ;;
  *)
    echo "NETWORK must be mainnet or testnet (got: $NETWORK)" >&2
    exit 1
    ;;
esac

mkdir -p "$OUT_DIR"

echo "Capturing inventory from $REST_API ($NETWORK / $CHAIN_ID)..."

HEIGHT="$(curl -sf --max-time 30 "${REST_API}/cosmos/base/tendermint/v1beta1/blocks/latest" \
  | jq -r '.block.header.height // empty')"
if [[ -z "$HEIGHT" ]]; then
  echo "Failed to reach REST API at $REST_API" >&2
  exit 1
fi

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

curl -sf --max-time 60 "${REST_API}/cosmwasm/wasm/v1/code?pagination.limit=500" -o "$TMP/wasm.json" || echo '{}' > "$TMP/wasm.json"
curl -sf --max-time 60 "${REST_API}/ibc/core/channel/v1/channels?pagination.limit=200" -o "$TMP/ibc_channels.json" || echo '{}' > "$TMP/ibc_channels.json"
curl -sf --max-time 60 "${REST_API}/ibc/core/client/v1/client_states?pagination.limit=100" -o "$TMP/ibc_clients.json" || echo '{}' > "$TMP/ibc_clients.json"
curl -sf --max-time 30 "${REST_API}/cosmos/upgrade/v1beta1/applied_plan" -o "$TMP/upgrade.json" || echo '{}' > "$TMP/upgrade.json"

# Use slurpfile instead of --argjson so large REST payloads work on Windows Git Bash.
jq -n \
  --arg network "$NETWORK" \
  --arg chain_id "$CHAIN_ID" \
  --arg rest_api "$REST_API" \
  --arg height "$HEIGHT" \
  --arg captured_at "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  --slurpfile wasm "$TMP/wasm.json" \
  --slurpfile ibc_channels "$TMP/ibc_channels.json" \
  --slurpfile ibc_clients "$TMP/ibc_clients.json" \
  --slurpfile upgrade "$TMP/upgrade.json" \
  '{
    network: $network,
    chain_id: $chain_id,
    rest_api: $rest_api,
    block_height: ($height | tonumber),
    captured_at: $captured_at,
    wasm_code_count: ($wasm[0].code_infos // [] | length),
    wasm_codes: [($wasm[0].code_infos // [])[] | {code_id: .code_id, creator: .creator, data_hash: .data_hash}],
    ibc_channel_count: ($ibc_channels[0].channels // [] | length),
    ibc_channels: [($ibc_channels[0].channels // [])[] | {state: .state, ordering: .ordering, counterparty: .counterparty, connection_hops: .connection_hops, port_id: .port_id, channel_id: .channel_id, version: .version}],
    ibc_client_count: ($ibc_clients[0].client_states // [] | length),
    ibc_clients: [($ibc_clients[0].client_states // [])[] | {client_id: .client_id, client_state: .client_state["@type"] // .client_state.type_url // "unknown"}],
    applied_upgrade: $upgrade[0]
  }' > "$OUT_FILE"

echo "Wrote $OUT_FILE"
echo "  height: $HEIGHT"
echo "  wasm codes: $(jq '.wasm_code_count' "$OUT_FILE")"
echo "  ibc channels: $(jq '.ibc_channel_count' "$OUT_FILE")"
echo ""
echo "Merge highlights into docs/PHASE0-INVENTORY.md before v600 testnet wasm/IBC regression."
