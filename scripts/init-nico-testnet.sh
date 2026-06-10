#!/usr/bin/env bash
# Initialize a single-validator private net (jackal-nico-1) for v600 testing.
#
# Prerequisites: canined on PATH (WSL: make install after bootstrap-wsl-dev.sh)
#   or CANINED=/path/to/canined
#
# Usage:
#   ./scripts/init-nico-testnet.sh              # init only
#   START=1 ./scripts/init-nico-testnet.sh      # init + canined start
#   RESET=1 ./scripts/init-nico-testnet.sh      # wipe ~/.canine-nico first

set -euo pipefail

CHAIN_ID="${CHAIN_ID:-jackal-nico-1}"
MONIKER="${MONIKER:-nico-validator}"
HOME_DIR="${HOME_DIR:-$HOME/.canine-nico}"
KEYRING="${KEYRING:-test}"
DENOM="${DENOM:-ujkl}"
CANINED="${CANINED:-canined}"

if ! command -v "$CANINED" >/dev/null 2>&1; then
  if [[ -x "$(pwd)/build/canined-linux-amd64" ]]; then
    CANINED="$(pwd)/build/canined-linux-amd64"
  elif [[ -x "$HOME/go/bin/canined" ]]; then
    CANINED="$HOME/go/bin/canined"
  else
    echo "canined not found. Run: make install  (WSL) or bash build-linux.sh" >&2
    exit 1
  fi
fi

if [[ "${RESET:-0}" == "1" ]]; then
  echo "Removing $HOME_DIR"
  rm -rf "$HOME_DIR"
fi

if [[ -f "$HOME_DIR/config/genesis.json" ]]; then
  echo "Genesis already exists at $HOME_DIR — use RESET=1 to re-init" >&2
  exit 1
fi

echo "Using $CANINED ($("$CANINED" version 2>/dev/null | head -1 || true))"
echo "Chain ID: $CHAIN_ID"
echo "Home:     $HOME_DIR"

"$CANINED" init "$MONIKER" --chain-id "$CHAIN_ID" --home "$HOME_DIR"

if ! "$CANINED" keys show validator --home "$HOME_DIR" --keyring-backend "$KEYRING" >/dev/null 2>&1; then
  "$CANINED" keys add validator --home "$HOME_DIR" --keyring-backend "$KEYRING"
fi

VALIDATOR_ADDR="$("$CANINED" keys show validator -a --home "$HOME_DIR" --keyring-backend "$KEYRING")"

"$CANINED" add-genesis-account validator "100000000000000${DENOM}" \
  --keyring-backend "$KEYRING" --home "$HOME_DIR"

"$CANINED" gentx validator "1000000${DENOM}" \
  --chain-id "$CHAIN_ID" --keyring-backend "$KEYRING" --home "$HOME_DIR"

"$CANINED" collect-gentxs --home "$HOME_DIR"

GENESIS="$HOME_DIR/config/genesis.json"
TMP="$HOME_DIR/config/tmp_genesis.json"
jq \
  --arg denom "$DENOM" \
  '.app_state.staking.params.bond_denom = $denom
   | .app_state.crisis.constant_fee.denom = $denom
   | .app_state.jklmint.params.mintDenom = $denom
   | .app_state.gov.deposit_params.min_deposit[0].denom = $denom
   | .app_state.gov.voting_params.voting_period = "120s"
   | .app_state.gov.deposit_params.max_deposit_period = "120s"' \
  "$GENESIS" > "$TMP" && mv "$TMP" "$GENESIS"

"$CANINED" validate-genesis --home "$HOME_DIR"

echo ""
echo "Private testnet ready."
echo "  Validator: $VALIDATOR_ADDR"
echo "  Start:     $CANINED start --home $HOME_DIR"
echo "  Smoke:     CHAIN_ID=$CHAIN_ID NODE=http://localhost:26657 KEY=validator ./scripts/smoke-v600-testnet.sh"
echo "  v600 gov:  ./scripts/submit-v600-upgrade-proposal.sh"

if [[ "${START:-0}" == "1" ]]; then
  exec "$CANINED" start --home "$HOME_DIR"
fi
