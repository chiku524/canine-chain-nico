#!/usr/bin/env bash
# Initialize a single-validator private net (jackal-nico-1) for modernization testing.
# Default binary: feat/cosmos-modernization-phase4 (SDK 0.54 / v630).
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

"$CANINED" config set client chain-id "$CHAIN_ID" --home "$HOME_DIR"
"$CANINED" config set client output json --home "$HOME_DIR"
"$CANINED" config set client keyring-backend "$KEYRING" --home "$HOME_DIR"

if ! "$CANINED" keys show validator --home "$HOME_DIR" --keyring-backend "$KEYRING" >/dev/null 2>&1; then
  "$CANINED" keys add validator --home "$HOME_DIR" --keyring-backend "$KEYRING"
fi

VALIDATOR_ADDR="$("$CANINED" keys show validator -a --home "$HOME_DIR" --keyring-backend "$KEYRING")"

"$CANINED" genesis add-genesis-account validator "100000000000000${DENOM}" \
  --keyring-backend "$KEYRING" --home "$HOME_DIR"

"$CANINED" genesis gentx validator "1000000${DENOM}" \
  --chain-id "$CHAIN_ID" --keyring-backend "$KEYRING" --home "$HOME_DIR"

"$CANINED" genesis collect-gentxs --home "$HOME_DIR"

GENESIS="$HOME_DIR/config/genesis.json"
TMP="$HOME_DIR/config/tmp_genesis.json"
# SDK 0.54 uses gov.params (legacy deposit_params / voting_params are null).
jq \
  --arg denom "$DENOM" \
  '.app_state.staking.params.bond_denom = $denom
   | .app_state.jklmint.params.mint_denom = $denom
   | .app_state.gov.params.min_deposit[0].denom = $denom
   | .app_state.gov.params.expedited_min_deposit[0].denom = $denom
   | .app_state.gov.params.voting_period = "120s"
   | .app_state.gov.params.max_deposit_period = "120s"
   | .app_state.gov.params.expedited_voting_period = "60s"
   | if .app_state.crisis then .app_state.crisis.constant_fee.denom = $denom else . end' \
  "$GENESIS" > "$TMP" && mv "$TMP" "$GENESIS"

"$CANINED" genesis validate --home "$HOME_DIR"

# Local smoke convenience: enable REST API (1317).
APP_TOML="$HOME_DIR/config/app.toml"
if [[ -f "$APP_TOML" ]]; then
  sed -i '0,/^enable = false/{s/^enable = false/enable = true/}' "$APP_TOML" || true
fi

echo ""
echo "Private testnet ready."
echo "  Validator: $VALIDATOR_ADDR"
echo "  Start:     $CANINED start --home $HOME_DIR"
echo "  Smoke:     CHAIN_ID=$CHAIN_ID NODE=http://127.0.0.1:26657 KEY=validator ./scripts/smoke-v630-testnet.sh"
echo "  Verify:    SKIP_SIM=1 ./scripts/verify-v630-candidate.sh"

if [[ "${START:-0}" == "1" ]]; then
  exec "$CANINED" start --home "$HOME_DIR"
fi
