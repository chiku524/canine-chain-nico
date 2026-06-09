#!/usr/bin/env bash
# Smoke tests for v600 testnet candidate. Run against a funded testnet key.
#
# Usage:
#   CHAIN_ID=jackal-testnet-1 NODE=https://rpc.example.com KEY=mykey \
#     ./scripts/smoke-v600-testnet.sh
#
# Requires: canined, jq, configured key in canined keyring.

set -eo pipefail

CHAIN_ID="${CHAIN_ID:-jackal-testnet-1}"
NODE="${NODE:-http://localhost:26657}"
KEY="${KEY:-test}"
DENOM="${DENOM:-ujkl}"

canined() {
  command canined --node "tcp://${NODE#*://}" --chain-id "$CHAIN_ID" "$@"
}

echo "== v600 smoke: chain status =="
canined status | jq -r '.SyncInfo.latest_block_height, .NodeInfo.version'

echo "== v600 smoke: upgrade plan =="
canined query upgrade plan 2>/dev/null | jq . || echo "(no scheduled upgrade — OK before governance)"

echo "== v600 smoke: consensus params module =="
canined query consensus params | jq .

echo "== v600 smoke: bank balance =="
canined query bank balances "$(canined keys show "$KEY" -a)" | jq .

echo "== v600 smoke: storage params =="
canined query storage params | jq .

echo "== v600 smoke: filetree params =="
canined query filetree params | jq .

echo "== v600 smoke: oracle feeds =="
canined query oracle list-feeds 2>/dev/null | jq . || true

echo "== v600 smoke: wasm code list =="
canined query wasm list-code 2>/dev/null | jq '.code_infos | length' || true

echo "== v600 smoke: IBC clients =="
canined query ibc client states 2>/dev/null | jq '.client_states | length' || true

echo "== v600 smoke: zero-fee post-proof tx dry-run =="
# Replace with a valid MsgPostProof from your provider when running on live testnet.
echo "Manual: submit MsgPostProof-only tx and confirm fee=0 (see docs/V600-TESTNET-UPGRADE.md)"

echo "All automated smoke checks completed."
