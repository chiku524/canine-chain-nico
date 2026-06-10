#!/usr/bin/env bash
# Pre-release checks for the v600 (SDK 0.47) candidate on master.
# Run on Linux with CGO + wasmvm 1.5.9 installed.
#
# Usage: ./scripts/verify-v600-candidate.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "== v600 candidate: go version =="
go version

echo "== v600 candidate: build canined =="
CGO_ENABLED=1 go build -o build/canined ./cmd/canined
./build/canined version 2>/dev/null || true

echo "== v600 candidate: v600 upgrade unit tests =="
go test ./app/upgrades/v600/...

echo "== v600 candidate: ante fee tests (requires CGO) =="
if [[ "${CGO_ENABLED:-1}" == "1" ]]; then
  go test -tags=cgo -count=1 ./app -run 'TestJackalDeductFee|TestWasmdExport|TestBlockedAddrs' || {
    echo "CGO app tests failed or skipped — ensure wasmvm 1.5.9 lib is installed" >&2
    exit 1
  }
else
  echo "Skipping CGO app tests (CGO_ENABLED=0)"
fi

echo "== v600 candidate: export/import round-trip =="
echo "Covered by TestWasmdExport when CGO is enabled."

echo ""
echo "Next steps:"
echo "  1. NETWORK=mainnet ./scripts/capture-chain-inventory.sh"
echo "  2. CGO_ENABLED=1 WASMVM_TAG=v1.5.9 make build-linux"
echo "  3. Deploy v600 on testnet — docs/V600-TESTNET-UPGRADE.md"
echo "  4. CHAIN_ID=<testnet> NODE=<rpc> KEY=<key> ./scripts/smoke-v600-testnet.sh"
