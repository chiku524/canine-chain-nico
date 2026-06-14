# Cosmovisor upgrade ladder (v600 → v630)

Guide for validators running **sequential** Jackal modernization upgrades with [Cosmovisor](https://docs.cosmos.network/sdk/v0.54/build/tooling/cosmovisor).

---

## Install Cosmovisor

```bash
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.7.0
# or: ./scripts/cosmovisor.sh
```

Set environment (example):

```bash
export DAEMON_NAME=canined
export DAEMON_HOME=$HOME/.canine
export DAEMON_ALLOW_DOWNLOAD_BINARIES=false
export UNSAFE_SKIP_BACKUP=false
```

Run the node via Cosmovisor:

```bash
cosmovisor run start
```

---

## Directory layout

```
$DAEMON_HOME/cosmovisor/
├── genesis/
│   └── bin/
│       └── canined          # binary active before any upgrade
├── upgrades/
│   ├── v600/bin/canined
│   ├── v610/bin/canined
│   ├── v620/bin/canined
│   └── v630/bin/canined
└── current -> upgrades/v630  # or genesis/ before first upgrade
```

**Rule:** The on-chain upgrade **name** must match the folder under `upgrades/` exactly (`v600`, `v610`, `v620`, `v630`).

---

## Upgrade ladder

| Order | Name | Branch / stack | wasmvm | Notes |
|-------|------|----------------|--------|-------|
| 1 | `v600` | `master` — SDK 0.47 | v1.5.x | Adds `consensus`, `crisis` stores |
| 2 | `v610` | phase2 — SDK 0.50 | v2.x | Adds `circuit` |
| 3 | `v620` | phase3 — SDK 0.53 | v2.2.x | Deletes `capability`, `feeibc` |
| 4 | `v630` | phase4 — SDK 0.54 | **v3.x** | Deletes `crisis`, `circuit`; CometBFT 0.39 |

### wasmvm lib per hop

Install the matching shared library **before** each binary swap:

```bash
WASMVM_VERSION=$(go list -m -f '{{.Version}}' github.com/CosmWasm/wasmvm/v3)  # adjust module path per phase
sudo wget -q "https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm.x86_64.so" \
  -O /usr/lib/libwasmvm.x86_64.so
sudo ldconfig 2>/dev/null || true
```

| Phase | go list module |
|-------|----------------|
| v600 | `github.com/CosmWasm/wasmvm` (v1) |
| v610, v620 | `github.com/CosmWasm/wasmvm/v2` |
| v630 | `github.com/CosmWasm/wasmvm/v3` |

---

## Private net rehearsal

```bash
# 1. Start on genesis binary (phase4 binary for greenfield 0.54 testing)
RESET=1 ./scripts/init-nico-testnet.sh
cosmovisor run start --home ~/.canine-nico

# 2. Schedule upgrade (example — only meaningful if chain started on older binary)
UPGRADE_NAME=v630 UPGRADE_HEIGHT=50 HOME_DIR=~/.canine-nico \
  ./scripts/submit-upgrade-proposal.sh

# 3. Pre-place next binary
mkdir -p ~/.canine-nico/cosmovisor/upgrades/v630/bin
cp build/canined ~/.canine-nico/cosmovisor/upgrades/v630/bin/
```

---

## upgrade-info.json (coordinated testnet)

```json
{
  "name": "v630",
  "height": 1234567,
  "info": "{\"binaries\":{\"linux/amd64\":\"<sha256>:https://example.com/canined-v630-linux-amd64\"}}"
}
```

---

## References

- [V600-TESTNET-UPGRADE.md](./V600-TESTNET-UPGRADE.md)
- [V610-V620-UPGRADE-NOTES.md](./V610-V620-UPGRADE-NOTES.md)
- [V630-TESTNET-UPGRADE.md](./V630-TESTNET-UPGRADE.md)
