#!/usr/bin/env sh

set -eu

echo "Generating gogo proto code"
cd proto

proto_count=0
for file in $(find . -name '*.proto' | sort); do
  if grep -q go_package "$file"; then
    echo "  buf generate ${file#./}"
    buf generate --template buf.gen.gogo.yaml "$file"
    proto_count=$((proto_count + 1))
  fi
done

echo "processed ${proto_count} proto files"
if [ "$proto_count" -eq 0 ]; then
  echo "no .proto files found under proto/"
  exit 1
fi

cd ..

if [ ! -d github.com/jackalLabs/canine-chain/x ]; then
  echo "protobuf generation produced no output (expected github.com/jackalLabs/canine-chain/x)"
  exit 1
fi

# Merge generated module tree into repo (proto go_package omits /v5 major version).
cp -r github.com/jackalLabs/canine-chain/x/. x/
rm -rf github.com
