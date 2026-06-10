#!/bin/bash
# Build canined for Linux using Docker (works from Windows Git Bash / macOS / Linux).

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ARCH="${1:-amd64}"

mkdir -p "$SCRIPT_DIR/build"

echo "Building canined for Linux ($ARCH) via Docker..."
echo "Build context: $SCRIPT_DIR"

docker build \
  --build-arg GOARCH="$ARCH" \
  -f "$SCRIPT_DIR/Dockerfile.linux-build" \
  -t canined-linux-builder \
  "$SCRIPT_DIR"

CONTAINER_ID=$(docker create canined-linux-builder)
docker cp "$CONTAINER_ID:/workspace/canine-chain/build/canined" "$SCRIPT_DIR/build/canined-linux-$ARCH"

WASMVM_PRESENT=false
if docker cp "$CONTAINER_ID:/workspace/canine-chain/build/libwasmvm.so" "$SCRIPT_DIR/build/libwasmvm-linux-$ARCH.so" 2>/dev/null; then
  WASMVM_PRESENT=true
fi
docker rm "$CONTAINER_ID" >/dev/null

echo ""
echo "Build complete:"
echo "  $SCRIPT_DIR/build/canined-linux-$ARCH"
if [ "$WASMVM_PRESENT" = true ]; then
  echo "  $SCRIPT_DIR/build/libwasmvm-linux-$ARCH.so"
fi
