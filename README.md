![canine banner](assets/jackal_logo.png)
# Canine Chain
**Canine Chain** is a core component to the Jackal Protocol, a distributed cloud storage platform running on blockchain rails. This component is the replicated state machine that manages storage deals, payments and data-permissions. It is built using the Cosmos-SDK and CometBFT (formerly Tendermint).

[![golangci-lint](https://github.com/JackalLabs/canine-chain/actions/workflows/golangci.yml/badge.svg)](https://github.com/JackalLabs/canine-chain/actions/workflows/golangci.yml)
[![Test](https://github.com/JackalLabs/canine-chain/actions/workflows/test-unit.yml/badge.svg)](https://github.com/JackalLabs/canine-chain/actions/workflows/test-unit.yml)
[![Build](https://github.com/JackalLabs/canine-chain/actions/workflows/build.yml/badge.svg)](https://github.com/JackalLabs/canine-chain/actions/workflows/build.yml)

## Wiki Pages

1. [Modules](x/README.md)
2. [Tokens](TOKENS.md)
3. [Storage Providers](cmd/canined/README.md)
4. [**Cosmos stack modernization**](docs/COSMOS-MODERNIZATION.md) — living roadmap & checklist (SDK 0.47 → 0.54)
5. [Phase 0 inventory](docs/PHASE0-INVENTORY.md)
6. [v600 testnet upgrade playbook](docs/V600-TESTNET-UPGRADE.md)
7. [v600 mainnet governance template](docs/V600-MAINNET-GOVERNANCE.md)
8. [Public network endpoints](docs/NETWORK-ENDPOINTS.md)
9. [Windows dev setup (WSL / Docker)](docs/WINDOWS-DEV.md)
10. [Private testnet (fork-owned)](docs/PRIVATE-TESTNET.md)


## Installing the Canine CLI
### Prerequisites
* **Go 1.23.8** on `master` (Phase 1 / v600; see `go.mod`; north-star 0.54 stack: Go 1.25+)
* GNU Make and a C toolchain (`build-essential` on Debian/Ubuntu)
* **CGO enabled** for `canined` and wasm tests (`CGO_ENABLED=1`)
* **wasmvm shared library** matching `go.mod` (migration branch: **v1.5.9**):

```sh
# Linux amd64 — replace TAG with the wasmvm version from go.mod (e.g. v1.5.9)
WASMVM_TAG=v1.5.9
sudo wget -q "https://github.com/CosmWasm/wasmvm/raw/${WASMVM_TAG}/internal/api/libwasmvm.x86_64.so" \
  -O /usr/lib/libwasmvm.x86_64.so
```

On macOS use `libwasmvm.dylib` from the same tag path under `internal/api/`.

### Installing
> if you want to use pebble follow this: https://github.com/JackalLabs/canine-chain/pull/511

To install `canined` on your Linux machine:

```shell
git clone https://github.com/JackalLabs/canine-chain.git
cd canine-chain
make install
```

### Pre-built Binary
[Releases](https://github.com/JackalLabs/canine-chain/releases) — download the latest release for your network. Install the **wasmvm** shared library version that matches the release (see `go.mod` and [docs/COSMOS-MODERNIZATION.md](docs/COSMOS-MODERNIZATION.md)). Example for wasmvm **1.5.x** (migration branch):

```sh
# Replace TAG with the wasmvm version from go.mod (e.g. v1.5.9)
sudo wget "https://github.com/CosmWasm/wasmvm/raw/TAG/internal/api/libwasmvm.x86_64.so" -O /lib/libwasmvm.x86_64.so
```

You may also need to run `sudo chmod +x canined` inside the executables directory to allow it to run.

## Testing this chain

```shell
make test
```

### v600 candidate verification (Linux + CGO)

```shell
./scripts/verify-v600-candidate.sh
./scripts/capture-chain-inventory.sh          # mainnet wasm + IBC snapshot
NETWORK=testnet ./scripts/capture-chain-inventory.sh
```

## Version Map

When Syncing, you **MUST** use the flag `--unsafe-skip-upgrades 118040` after `canined start` or else you will crash at height 118040.

|block height|canined version|
|------------|---------------|
|45381       |1.1.2          |
|0           |1.1.0          |

## License

Canine by Jackal uses the [MIT License](/LICENSE.md).

## Bug Bounty

Refer to the bug bounty program proposed by Jackal Labs [Here](https://jackaldao.medium.com/announcement-jackal-bug-bounty-program-31d4e03ab7e2)

### [Developer Contact](/ABOUT.md)

