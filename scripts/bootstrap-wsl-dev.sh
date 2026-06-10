#!/usr/bin/env bash
# One-time dev environment setup inside WSL Ubuntu.
# Run from repo root: ./scripts/bootstrap-wsl-dev.sh
#
# Open Ubuntu first:  wsl -d Ubuntu
# Then:  cd /mnt/c/Users/chiku/Desktop/Jackal/canine-chain-nico

set -euo pipefail

echo "== Installing build dependencies =="
sudo apt-get update
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y \
  build-essential make jq curl git wget

WASMVM_TAG="${WASMVM_TAG:-v1.5.9}"
LIB="/usr/lib/libwasmvm.x86_64.so"
if [[ ! -f "$LIB" ]]; then
  echo "== Installing wasmvm ${WASMVM_TAG} =="
  sudo wget -q "https://github.com/CosmWasm/wasmvm/raw/${WASMVM_TAG}/internal/api/libwasmvm.x86_64.so" -O "$LIB"
fi

if ! command -v go >/dev/null 2>&1; then
  echo "== Installing Go 1.23.8 =="
  GO_VER=1.23.8
  curl -fsSL "https://go.dev/dl/go${GO_VER}.linux-amd64.tar.gz" -o /tmp/go.tar.gz
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf /tmp/go.tar.gz
  grep -q '/usr/local/go/bin' "$HOME/.bashrc" || echo 'export PATH=/usr/local/go/bin:$PATH' >> "$HOME/.bashrc"
  export PATH=/usr/local/go/bin:$PATH
fi

export CGO_ENABLED=1
export PATH="/usr/local/go/bin:${PATH}"

echo "== Tool versions =="
go version
make --version | head -1
gcc --version | head -1

echo ""
echo "Setup complete. From this WSL shell:"
echo "  cd $(pwd)"
echo "  ./scripts/verify-v600-candidate.sh"
echo "  make inventory-mainnet"
echo "  CGO_ENABLED=1 WASMVM_TAG=v1.5.9 make build-linux"
