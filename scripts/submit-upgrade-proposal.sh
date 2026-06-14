#!/usr/bin/env bash
# Submit and auto-vote a software-upgrade proposal on a local/private net.
#
# Usage:
#   UPGRADE_NAME=v630 UPGRADE_HEIGHT=200 ./scripts/submit-upgrade-proposal.sh
#   UPGRADE_NAME=v610 UPGRADE_HEIGHT=50 HOME_DIR=~/.canine-nico ./scripts/submit-upgrade-proposal.sh

set -euo pipefail

CHAIN_ID="${CHAIN_ID:-jackal-nico-1}"
HOME_DIR="${HOME_DIR:-$HOME/.canine-nico}"
KEYRING="${KEYRING:-test}"
DENOM="${DENOM:-ujkl}"
UPGRADE_NAME="${UPGRADE_NAME:-v630}"
UPGRADE_HEIGHT="${UPGRADE_HEIGHT:-30}"
DEPOSIT="${DEPOSIT:-20000000${DENOM}}"
CANINED="${CANINED:-canined}"
PROPOSAL_ID="${PROPOSAL_ID:-}"

if ! command -v "$CANINED" >/dev/null 2>&1; then
  CANINED="${CANINED:-$HOME/go/bin/canined}"
fi

"$CANINED" tx gov submit-proposal software-upgrade "$UPGRADE_NAME" \
  --upgrade-height "$UPGRADE_HEIGHT" \
  --upgrade-info "${UPGRADE_NAME} modernization" \
  --title "Upgrade to $UPGRADE_NAME" \
  --description "Private/devnet upgrade to $UPGRADE_NAME" \
  --deposit "$DEPOSIT" \
  --from validator \
  --keyring-backend "$KEYRING" \
  --chain-id "$CHAIN_ID" \
  --home "$HOME_DIR" \
  -y \
  --broadcast-mode sync

sleep 6

if [[ -z "$PROPOSAL_ID" ]]; then
  PROPOSAL_ID=$("$CANINED" query gov proposals --home "$HOME_DIR" --output json | jq -r '.proposals[-1].id')
fi

"$CANINED" tx gov vote "$PROPOSAL_ID" yes \
  --from validator \
  --keyring-backend "$KEYRING" \
  --chain-id "$CHAIN_ID" \
  --home "$HOME_DIR" \
  -y \
  --broadcast-mode sync

echo ""
echo "Submitted and voted yes on proposal $PROPOSAL_ID for $UPGRADE_NAME at height $UPGRADE_HEIGHT."
echo "Query: $CANINED q upgrade plan --home $HOME_DIR"
