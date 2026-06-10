#!/usr/bin/env bash
# Submit and auto-vote a v600 software-upgrade proposal on a local/private net.
#
# Usage (node must be running):
#   ./scripts/submit-v600-upgrade-proposal.sh
#   UPGRADE_HEIGHT=50 HOME_DIR=~/.canine-nico ./scripts/submit-v600-upgrade-proposal.sh

set -euo pipefail

CHAIN_ID="${CHAIN_ID:-jackal-nico-1}"
HOME_DIR="${HOME_DIR:-$HOME/.canine-nico}"
KEYRING="${KEYRING:-test}"
DENOM="${DENOM:-ujkl}"
UPGRADE_NAME="${UPGRADE_NAME:-v600}"
UPGRADE_HEIGHT="${UPGRADE_HEIGHT:-30}"
DEPOSIT="${DEPOSIT:-20000000${DENOM}}"
CANINED="${CANINED:-canined}"

if ! command -v "$CANINED" >/dev/null 2>&1; then
  CANINED="${CANINED:-$HOME/go/bin/canined}"
fi

"$CANINED" tx gov submit-proposal software-upgrade "$UPGRADE_NAME" \
  --upgrade-height "$UPGRADE_HEIGHT" \
  --upgrade-info "v600 SDK 0.47 modernization" \
  --title "Upgrade to $UPGRADE_NAME" \
  --description "Private testnet v600 validation" \
  --deposit "$DEPOSIT" \
  --from validator \
  --keyring-backend "$KEYRING" \
  --chain-id "$CHAIN_ID" \
  --home "$HOME_DIR" \
  -y \
  --broadcast-mode sync

sleep 6

"$CANINED" tx gov vote 1 yes \
  --from validator \
  --keyring-backend "$KEYRING" \
  --chain-id "$CHAIN_ID" \
  --home "$HOME_DIR" \
  -y \
  --broadcast-mode sync

echo ""
echo "Submitted proposal 1 for $UPGRADE_NAME at height $UPGRADE_HEIGHT."
echo "Query: $CANINED q upgrade plan --home $HOME_DIR"
echo "At height $UPGRADE_HEIGHT, restart with the v600 binary (same build if already on Phase 1)."
