#!/usr/bin/env bash
# Pre-devnet checks for the v630 (SDK 0.54) candidate on feat/cosmos-modernization-phase4.
# Run on Linux with CGO + wasmvm v3 installed.
#
# Usage: ./scripts/verify-v630-candidate.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

GO_MIN=1.25
TEST_TAGS='ledger test_ledger_mock test'

echo "== v630 candidate: go version =="
go version

echo "== v630 candidate: wasmvm shared library =="
WASMVM_VERSION=$(go list -m -f '{{.Version}}' github.com/CosmWasm/wasmvm/v3)
echo "go.mod wasmvm: $WASMVM_VERSION"
if [[ ! -f /usr/lib/libwasmvm.x86_64.so ]] && [[ ! -f /lib/libwasmvm.x86_64.so ]]; then
  echo "WARNING: libwasmvm.x86_64.so not found in /usr/lib or /lib" >&2
  echo "  sudo wget -q https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm.x86_64.so -O /usr/lib/libwasmvm.x86_64.so" >&2
fi

echo "== v630 candidate: build canined =="
CGO_ENABLED=1 go build -o build/canined ./cmd/canined
./build/canined version 2>/dev/null || true

echo "== v630 candidate: modernization upgrade unit tests =="
go test ./app/upgrades/v600/... ./app/upgrades/v610/... ./app/upgrades/v620/... ./app/upgrades/v630/...

echo "== v630 candidate: ante fee + export tests (requires CGO) =="
if [[ "${CGO_ENABLED:-1}" == "1" ]]; then
  go test -tags="cgo ${TEST_TAGS}" -count=1 ./app -run 'TestJackalDeductFee|TestWasmdExport|TestBlockedAddrs' || {
    echo "CGO app tests failed — ensure wasmvm v3 lib is installed" >&2
    exit 1
  }
else
  echo "Skipping CGO app tests (CGO_ENABLED=0)"
fi

echo "== v630 candidate: unit tests =="
CGO_ENABLED=1 make test-unit

echo "== v630 candidate: simulation tests (optional, slow) =="
if [[ "${SKIP_SIM:-0}" != "1" ]]; then
  CGO_ENABLED=1 make test-sim-import-export
  CGO_ENABLED=1 make test-sim-full-app
else
  echo "Skipped (set SKIP_SIM=0 to run)"
fi

echo ""
echo "v630 candidate verification passed."
echo ""
echo "Next steps:"
echo "  1. NETWORK=mainnet ./scripts/capture-chain-inventory.sh"
echo "  2. CGO_ENABLED=1 make build-linux"
echo "  3. RESET=1 ./scripts/init-nico-testnet.sh && START=1 ./scripts/init-nico-testnet.sh"
echo "  4. CHAIN_ID=jackal-nico-1 NODE=http://127.0.0.1:26657 KEY=validator ./scripts/smoke-v630-testnet.sh"
echo "  5. Jackal devnet coordination — docs/JACKAL-DEVNET-HANDOFF.md"
