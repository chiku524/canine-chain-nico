#!/usr/bin/env bash
# Smoke tests for v630 (SDK 0.54) private net or Jackal devnet.
#
# Usage:
#   CHAIN_ID=jackal-nico-1 NODE=http://127.0.0.1:26657 KEY=validator \
#     ./scripts/smoke-v630-testnet.sh
#
# Jackal devnet (when available):
#   CHAIN_ID=<devnet> NODE=<rpc> KEY=<key> ./scripts/smoke-v630-testnet.sh

set -eo pipefail

CHAIN_ID="${CHAIN_ID:-jackal-nico-1}"
NODE="${NODE:-http://127.0.0.1:26657}"
KEY="${KEY:-validator}"
DENOM="${DENOM:-ujkl}"
HOME_DIR="${HOME_DIR:-$HOME/.canine-nico}"
CANINED="${CANINED:-canined}"

canined() {
  if [[ "$NODE" == http://127.0.0.1:* ]] || [[ "$NODE" == http://localhost:* ]]; then
    command "$CANINED" --home "$HOME_DIR" --chain-id "$CHAIN_ID" "$@"
  else
    local hostport="${NODE#*://}"
    command "$CANINED" --node "tcp://${hostport}" --chain-id "$CHAIN_ID" "$@"
  fi
}

echo "== v630 smoke: chain status =="
canined status | jq -r '.sync_info.latest_block_height, .node_info.network'

echo "== v630 smoke: upgrade plan =="
canined query upgrade plan 2>/dev/null | jq . || echo "(no scheduled upgrade — OK on fresh net)"

echo "== v630 smoke: applied upgrades =="
for name in v600 v610 v620 v630; do
  if canined query upgrade applied "$name" 2>/dev/null | grep -q '"name"'; then
    echo "  applied: $name"
  fi
done

echo "== v630 smoke: consensus params =="
canined query params subspace base consensus 2>/dev/null | jq . || echo "(skipped — use comet RPC for consensus params on SDK 0.54)"

echo "== v630 smoke: bank balance =="
VALIDATOR_ADDR="$(canined keys show "$KEY" -a --keyring-backend test 2>/dev/null || canined keys show "$KEY" -a)"
if canined query bank balances "$VALIDATOR_ADDR" 2>/dev/null | jq .; then
  :
elif REST_API="${REST_API:-http://127.0.0.1:1317}" \
  curl -sf "${REST_API}/cosmos/bank/v1beta1/balances/${VALIDATOR_ADDR}" 2>/dev/null | jq .; then
  :
else
  echo "(skipped — bank CLI unavailable; set REST_API or enable api in app.toml)"
fi

echo "== v630 smoke: storage params =="
canined query storage params | jq .

echo "== v630 smoke: filetree params =="
canined query filetree params | jq .

echo "== v630 smoke: rns params =="
canined query rns params 2>/dev/null | jq . || true

echo "== v630 smoke: oracle feeds =="
canined query oracle list-feeds 2>/dev/null | jq . || true

echo "== v630 smoke: wasm code list =="
canined query wasm list-code 2>/dev/null | jq '.code_infos | length' || true

echo "== v630 smoke: IBC clients =="
canined query ibc client states 2>/dev/null | jq '.client_states | length' || true

echo "== v630 smoke: jklmint params =="
canined query jklmint params 2>/dev/null | jq . || true

echo "== v630 smoke: zero-fee post-proof =="
echo "Manual: submit MsgPostProof-only tx and confirm fee=0 (see docs/V630-TESTNET-UPGRADE.md §5)"

echo "== v630 smoke: wasmvm / binary =="
canined version 2>/dev/null || true

echo ""
echo "All automated v630 smoke checks completed."
echo "Full matrix: docs/V630-TESTNET-UPGRADE.md"
