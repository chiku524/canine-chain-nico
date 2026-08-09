#!/usr/bin/env bash
# One-time dev environment setup inside WSL Ubuntu.
# Run from repo root: ./scripts/bootstrap-wsl-dev.sh
#
# Open Ubuntu first:  wsl -d Ubuntu
# Then:  cd /mnt/c/Users/chiku/Desktop/Jackal/canine-chain-nico

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "== Installing build dependencies =="
sudo apt-get update
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y \
  build-essential make jq curl git wget

export PATH="/usr/local/go/bin:${PATH}"

GO_VER="${GO_VER:-$(awk '/^go / {print $2}' go.mod)}"
GO_VER="${GO_VER:-1.25.9}"
NEED_GO=1
if command -v go >/dev/null 2>&1; then
  CURRENT_GO="$(go env GOVERSION 2>/dev/null | sed 's/^go//')"
  if [[ "$CURRENT_GO" == "$GO_VER" ]]; then
    NEED_GO=0
  fi
fi

if [[ "$NEED_GO" == "1" ]]; then
  echo "== Installing Go ${GO_VER} =="
  curl -fsSL "https://go.dev/dl/go${GO_VER}.linux-amd64.tar.gz" -o /tmp/go.tar.gz
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf /tmp/go.tar.gz
  grep -q '/usr/local/go/bin' "$HOME/.bashrc" || echo 'export PATH=/usr/local/go/bin:$PATH' >> "$HOME/.bashrc"
  export PATH=/usr/local/go/bin:$PATH
fi

WASMVM_VERSION="$(go list -m -f '{{.Version}}' github.com/CosmWasm/wasmvm/v3 2>/dev/null || true)"
if [[ -z "$WASMVM_VERSION" ]]; then
  WASMVM_VERSION="$(go list -m -f '{{.Version}}' github.com/CosmWasm/wasmvm 2>/dev/null || true)"
fi
WASMVM_VERSION="${WASMVM_VERSION:-v3.0.4}"

LIB="/usr/lib/libwasmvm.x86_64.so"
echo "== Installing wasmvm ${WASMVM_VERSION} =="
sudo wget -q "https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm.x86_64.so" \
  -O "$LIB"

export CGO_ENABLED=1
export PATH="/usr/local/go/bin:${PATH}"

echo "== Tool versions =="
go version
make --version | head -1
gcc --version | head -1
ls -la "$LIB"

echo ""
echo "Setup complete. From this WSL shell:"
echo "  cd $ROOT"
echo "  make install"
echo "  SKIP_SIM=1 ./scripts/verify-v630-candidate.sh"
echo "  RESET=1 START=1 ./scripts/init-nico-testnet.sh"
