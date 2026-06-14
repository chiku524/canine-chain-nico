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
4. [**Cosmos stack modernization**](docs/COSMOS-MODERNIZATION.md) — living roadmap (Phases 1–4 **code complete**)
5. [**Jackal devnet handoff**](docs/JACKAL-DEVNET-HANDOFF.md) — team brief for Jackal developer coordination
6. [Phase 0 inventory](docs/PHASE0-INVENTORY.md)
7. [v630 testnet upgrade playbook](docs/V630-TESTNET-UPGRADE.md) — SDK 0.54 devnet validation
8. [v610 / v620 upgrade notes](docs/V610-V620-UPGRADE-NOTES.md)
9. [CometBFT 0.39 validator guide](docs/COMETBFT-039-VALIDATOR.md)
10. [Cosmovisor upgrade ladder](docs/COSMOVISOR-LADDER.md)
11. [v600 testnet upgrade playbook](docs/V600-TESTNET-UPGRADE.md)
12. [v600 mainnet governance template](docs/V600-MAINNET-GOVERNANCE.md)
13. [Public network endpoints](docs/NETWORK-ENDPOINTS.md)
14. [Windows dev setup (WSL / Docker)](docs/WINDOWS-DEV.md)
15. [Private testnet (fork-owned)](docs/PRIVATE-TESTNET.md)


## Installing the Canine CLI
### Prerequisites
* **Go 1.25.9** on `feat/cosmos-modernization-phase4` (SDK 0.54); **Go 1.23.8** on `master` (Phase 1 / v600)
* GNU Make and a C toolchain (`build-essential` on Debian/Ubuntu)
* **CGO enabled** for `canined` and wasm tests (`CGO_ENABLED=1`)
* **wasmvm shared library** matching `go.mod`:

```sh
# Linux amd64 — version from go.mod (Phase 1: v1.5.9; Phase 2+: v2.x from releases/)
WASMVM_VERSION=$(go list -m -f '{{.Version}}' github.com/CosmWasm/wasmvm/v3 2>/dev/null \
  || go list -m -f '{{.Version}}' github.com/CosmWasm/wasmvm)
sudo wget -q "https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm.x86_64.so" \
  -O /usr/lib/libwasmvm.x86_64.so
# Phase 1 (wasmvm v1.x) only — if the command above fails:
# sudo wget -q "https://github.com/CosmWasm/wasmvm/raw/v1.5.9/internal/api/libwasmvm.x86_64.so" -O /usr/lib/libwasmvm.x86_64.so
```

On macOS use `libwasmvm.dylib` from the same release URL pattern (`releases/download/${WASMVM_VERSION}/`).

### Installing
> if you want to use pebble follow this: https://github.com/JackalLabs/canine-chain/pull/511

To install `canined` on your Linux machine:

```shell
git clone https://github.com/JackalLabs/canine-chain.git
cd canine-chain
make install
```

### Pre-built Binary
[Releases](https://github.com/JackalLabs/canine-chain/releases) — download the latest release for your network. Install the **wasmvm** shared library version that matches the release (see `go.mod` and [docs/COSMOS-MODERNIZATION.md](docs/COSMOS-MODERNIZATION.md)):

```sh
WASMVM_VERSION=$(go list -m -f '{{.Version}}' github.com/CosmWasm/wasmvm/v3 2>/dev/null \
  || go list -m -f '{{.Version}}' github.com/CosmWasm/wasmvm)
sudo wget -q "https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm.x86_64.so" \
  -O /usr/lib/libwasmvm.x86_64.so
```

You may also need to run `sudo chmod +x canined` inside the executables directory to allow it to run.

## Testing this chain

```shell
make test
```

### v630 candidate verification (Linux + CGO, phase4 branch)

```shell
git checkout feat/cosmos-modernization-phase4
./scripts/verify-v630-candidate.sh            # SKIP_SIM=1 for faster run
./scripts/capture-chain-inventory.sh          # mainnet wasm + IBC snapshot
RESET=1 ./scripts/init-nico-testnet.sh      # private jackal-nico-1
CHAIN_ID=jackal-nico-1 NODE=http://127.0.0.1:26657 KEY=validator \
  ./scripts/smoke-v630-testnet.sh
```

Phase 1 (`master`): `./scripts/verify-v600-candidate.sh`

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

