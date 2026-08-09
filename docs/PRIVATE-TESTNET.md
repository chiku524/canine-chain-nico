# Private testnet (fork-owned)

Use this fork as a **separate chain** from Jackal public testnet/mainnet. No handoff from Jackal Labs is required.

| Field | Suggested value |
|-------|-----------------|
| Chain ID | `jackal-nico-1` |
| Bech32 prefix | `jkl` (same as Jackal — or change in `app/app.go` if you want isolation) |
| **Binary branch** | `feat/cosmos-modernization-phase4` (SDK 0.54 / `v630`) |
| Upgrade ladder | `v600` → `v610` → `v620` → `v630` (sequential; see [COSMOS-MODERNIZATION.md](./COSMOS-MODERNIZATION.md)) |
| Purpose | Validate modernization before upstream merge / Jackal devnet coordination |

---

## Why a private testnet?

- Jackal public testnet is operated by the core team; your fork can run **`v630` (0.54)** on infrastructure you control while coordinating with Jackal devnet ([handoff](./JACKAL-DEVNET-HANDOFF.md)).
- Proves export/import, wasmvm v3 contracts, storage proofs, and sequential upgrade handlers without waiting for governance on `jackal-1`.
- Becomes the evidence package when you open a PR to [JackalLabs/canine-chain](https://github.com/JackalLabs/canine-chain).

---

## Option A — Local (WSL, fastest)

Three-validator net on one machine:

```bash
# WSL or Linux — Go 1.25+, CGO, wasmvm v3 (see README)
git checkout feat/cosmos-modernization-phase4
make install
./scripts/multinode-local-testnet.sh
```

Chain ID in that script is `testing`. For a named private net, use Option B.

---

## Option B — Single-node private net (`jackal-nico-1`)

One command (WSL, after `make install`):

```bash
./scripts/init-nico-testnet.sh
canined start --home ~/.canine-nico
```

Or wipe + init + start together:

```bash
RESET=1 START=1 ./scripts/init-nico-testnet.sh
```

For a **full migration rehearsal**, schedule upgrades in order: `v600` → `v610` → `v620` → `v630` (handlers in `app/upgrades/`). For **0.54-only smoke**, start directly on the phase4 binary and test module behavior without historical state.

Manual steps (SDK 0.54 CLI — equivalent to the script):

```bash
export CHAIN_ID=jackal-nico-1
export MONIKER=nico-validator
export HOME_DIR=$HOME/.canine-nico

canined init $MONIKER --chain-id $CHAIN_ID --home $HOME_DIR
canined config set client chain-id $CHAIN_ID --home $HOME_DIR
canined config set client keyring-backend test --home $HOME_DIR
canined keys add validator --keyring-backend test --home $HOME_DIR

canined genesis add-genesis-account validator 100000000000000ujkl \
  --keyring-backend test --home $HOME_DIR

canined genesis gentx validator 1000000ujkl \
  --chain-id $CHAIN_ID --keyring-backend test --home $HOME_DIR

canined genesis collect-gentxs --home $HOME_DIR
canined genesis validate --home $HOME_DIR

canined start --home $HOME_DIR
```

---

## Option C — Cloud VPS testnet (multi-validator)

1. Provision 3+ Linux VPS (Ubuntu 22.04, 4 vCPU, 200 GB disk).
2. On each: run `scripts/bootstrap-wsl-dev.sh` equivalent (apt + **wasmvm v3** + Go **1.25** + build).
3. Build once: `make build-linux` → distribute `canined-linux-amd64` + checksum.
4. Validator 1: `init` + `gentx` + `collect-gentxs` → share `genesis.json`.
5. Validators 2–3: `init` + `gentx` → merge via `collect-gentxs` on validator 1.
6. Open RPC `26657`, peer `26656` between nodes.
7. Schedule **`v600`** upgrade at height H via gov proposal (chain starts on pre-v600 genesis if testing full migration path, or start directly on v600 binary for module testing only).

Document your peer IDs and genesis hash in `docs/inventory/nico-testnet-genesis.json`.

---

## Smoke tests on private net

```bash
CHAIN_ID=jackal-nico-1 NODE=http://127.0.0.1:26657 KEY=validator \
  ./scripts/smoke-v630-testnet.sh
```

Manual matrix: `docs/V630-TESTNET-UPGRADE.md` §5.

---

## Relationship to Jackal public networks

| Network | Who runs it | Your fork role |
|---------|-------------|----------------|
| `jackal-1` mainnet | Jackal validators | Read-only queries until upstream merges + gov |
| Jackal public testnet | Jackal team | Optional cross-check after private bake |
| **`jackal-nico-1`** | **You** | Primary modernization testbed |

After ≥2 weeks stable on `jackal-nico-1` (or Jackal devnet), use results in [JACKAL-DEVNET-HANDOFF.md](./JACKAL-DEVNET-HANDOFF.md) to support upstream PR and eventual Jackal mainnet governance.

---

## Local storage provider (nico-provider)

Official [Sequoia](https://github.com/JackalLabs/sequoia) still targets cosmos-sdk **0.45**, so this fork ships a lightweight Sequoia-compatible HTTP provider for private testing:

```bash
# Terminal A — chain (if not already running)
canined start --home ~/.canine-nico

# Terminal B — register + start provider on :3333
./scripts/start-nico-provider.sh
```

Endpoints: `GET /`, `GET /version`, `POST /v2/upload`, `GET /download/{merkle}`, `GET /list`.

Buy space and upload a file:

```bash
HOME=~/.canine-nico
ADDR=$(canined keys show validator -a --home $HOME --keyring-backend test)
# 1 GB for 30 days
canined tx storage buy-storage $ADDR 30 1000000000 ujkl \
  --from validator --home $HOME --keyring-backend test --chain-id jackal-nico-1 \
  --fees 5000ujkl --gas auto -y

# Post + upload to local provider
HEIGHT=$(canined status --home $HOME | jq -r .sync_info.latest_block_height)
canined tx storage post ./README.md $((HEIGHT+100000)) \
  --dest http://127.0.0.1:3333 --max_proofs 1 \
  --from validator --home $HOME --keyring-backend test --chain-id jackal-nico-1 \
  --fees 5000ujkl --gas auto -y
```

---

## Nico Lab (click-to-test web UI)

Local web app that exercises chain + provider actions with buttons:

```bash
# Terminal A — chain
canined start --home ~/.canine-nico

# Terminal B — provider
./scripts/start-nico-provider.sh

# Terminal C — lab UI
./scripts/start-nico-lab.sh
# open http://127.0.0.1:3456
```

