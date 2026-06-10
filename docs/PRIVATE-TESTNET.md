# Private testnet (fork-owned)

Use this fork as a **separate chain** from Jackal public testnet/mainnet. No handoff from Jackal Labs is required.

| Field | Suggested value |
|-------|-----------------|
| Chain ID | `jackal-nico-1` |
| Bech32 prefix | `jkl` (same as Jackal — or change in `app/app.go` if you want isolation) |
| Upgrade name for SDK 0.47 | `v600` |
| Purpose | Validate modernization before upstream merge |

---

## Why a private testnet?

- Jackal public testnet is operated by the core team; your fork can ship **`v600` immediately** on infrastructure you control.
- Proves export/import, wasm, storage proofs, and upgrade handler without waiting for governance on `jackal-1`.
- Becomes the evidence package when you open a PR to [JackalLabs/canine-chain](https://github.com/JackalLabs/canine-chain).

---

## Option A — Local (WSL, fastest)

Three-validator net on one machine:

```bash
# WSL, after bootstrap-wsl-dev.sh and make install
make install
./scripts/multinode-local-testnet.sh
```

Chain ID in that script is `testing`. For a named private net, use Option B.

---

## Option B — Single-node private net (`jackal-nico-1`)

```bash
export CHAIN_ID=jackal-nico-1
export MONIKER=nico-validator
export HOME_DIR=$HOME/.canine-nico

canined init $MONIKER --chain-id $CHAIN_ID --home $HOME_DIR
canined keys add validator --keyring-backend test --home $HOME_DIR

canined add-genesis-account validator 100000000000000ujkl \
  --keyring-backend test --home $HOME_DIR

canined gentx validator 1000000ujkl \
  --chain-id $CHAIN_ID --keyring-backend test --home $HOME_DIR

canined collect-gentxs --home $HOME_DIR
canined validate-genesis --home $HOME_DIR

canined start --home $HOME_DIR
```

Test `v600` upgrade locally with a governance proposal (see `scripts/upgrade-test.sh` pattern).

---

## Option C — Cloud VPS testnet (multi-validator)

1. Provision 3+ Linux VPS (Ubuntu 22.04, 4 vCPU, 200 GB disk).
2. On each: run `scripts/bootstrap-wsl-dev.sh` equivalent (apt + wasmvm + build).
3. Build once: `make build-linux` → distribute `canined-linux-amd64` + checksum.
4. Validator 1: `init` + `gentx` + `collect-gentxs` → share `genesis.json`.
5. Validators 2–3: `init` + `gentx` → merge via `collect-gentxs` on validator 1.
6. Open RPC `26657`, peer `26656` between nodes.
7. Schedule **`v600`** upgrade at height H via gov proposal (chain starts on pre-v600 genesis if testing full migration path, or start directly on v600 binary for module testing only).

Document your peer IDs and genesis hash in `docs/inventory/nico-testnet-genesis.json`.

---

## Smoke tests on private net

```bash
CHAIN_ID=jackal-nico-1 NODE=http://localhost:26657 KEY=validator \
  ./scripts/smoke-v600-testnet.sh
```

Manual matrix: `docs/V600-TESTNET-UPGRADE.md` §5.

---

## Relationship to Jackal public networks

| Network | Who runs it | Your fork role |
|---------|-------------|----------------|
| `jackal-1` mainnet | Jackal validators | Read-only queries until upstream merges + gov |
| Jackal public testnet | Jackal team | Optional cross-check after private bake |
| **`jackal-nico-1`** | **You** | Primary modernization testbed |

After ≥2 weeks stable on `jackal-nico-1`, use results to support upstream PR and eventual Jackal mainnet governance.
