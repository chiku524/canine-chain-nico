#!/usr/bin/env bash

set -eo pipefail

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find . -name '*.proto' -print0 | xargs -0 -n1 dirname | sort -u)
for dir in $proto_dirs; do
  for file in "$dir"/*.proto; do
    [ -f "$file" ] || continue
    if grep -q go_package "$file"; then
      buf generate --template buf.gen.gogo.yaml "$file"
    fi
  done
done

cd ..

if [ ! -d github.com/jackalLabs/canine-chain/x ]; then
  echo "protobuf generation produced no output (expected github.com/jackalLabs/canine-chain/x)"
  exit 1
fi

# Merge generated module tree into repo (proto go_package omits /v5 major version).
cp -r github.com/jackalLabs/canine-chain/x/. x/
rm -rf github.com
