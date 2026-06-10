# Phase 0 — Baseline inventory

Recorded **2026-06-11** for the Jackal (`canine-chain`) Cosmos modernization program.  
See [COSMOS-MODERNIZATION.md](./COSMOS-MODERNIZATION.md) for the living roadmap.

---

## Current codebase (`master` — Phase 1 / v600)

| Package | Version |
|---------|---------|
| Cosmos SDK | 0.47.17 |
| CometBFT | 0.37.15 |
| cometbft-db | 0.14.1 |
| ibc-go | v7.10.0 |
| wasmd | v0.45.0 |
| wasmvm | v1.5.9 |
| Go / toolchain | 1.23.1 / go1.23.8 |

---

## Mainnet baseline (pre-`v600`)

| Item | Value |
|------|-------|
| Chain ID | `jackal-1` (mainnet) |
| Bech32 prefix | `jkl` |
| Latest mainnet line | **v5.1.x** (`app/upgrades/v510` — protocol label v5.1.0) |
| On-chain upgrade before modernization | `v510` (last in `registerMainnetUpgradeHandlers` before `v600`) |

### Dependency pins (mainnet `master` before migration)

| Package | Version |
|---------|---------|
| Cosmos SDK | **0.45.17** via `replace` → `JackalLabs/cosmos-sdk-new` |
| CometBFT | **0.34.27** via `replace` → `TheMarstonConnell/cometbft` |
| ibc-go | **v4** |
| wasmd | **v0.32** |
| wasmvm | **v1.5.x** |
| Go | **1.23** |

### `replace` directives (mainnet — removed in Phase 1 except merkletree)

| Replace | Purpose |
|---------|---------|
| `JackalLabs/cosmos-sdk-new` | Free post-proof fee waiver in SDK fork |
| `TheMarstonConnell/cometbft` | Jackal-specific CometBFT patches |
| `github.com/wealdtech/go-merkletree/v2` → `TheMarstonConnell/go-merkletree/v2` | Storage proof merkletree (**retained** on migration branch) |

---

## Public network endpoints

See [NETWORK-ENDPOINTS.md](./NETWORK-ENDPOINTS.md). Capture live wasm + IBC data:

```bash
./scripts/capture-chain-inventory.sh
NETWORK=testnet ./scripts/capture-chain-inventory.sh
```

Output lands in `docs/inventory/captured-*.json` — merge code IDs and channel IDs into the tables below.

---

## Migration branch (historical: `feat/cosmos-modernization-phase1`, now merged to `master`)

| Package | Version |
|---------|---------|
| Cosmos SDK | 0.47.17 |
| CometBFT | 0.37.15 |
| cometbft-db | 0.14.1 |
| ibc-go | v7.10.0 |
| wasmd | v0.45.0 |
| wasmvm | v1.5.9 |
| Go / toolchain | 1.23.1 / go1.23.8 |

### Remaining `replace` directives

| Replace | Status |
|---------|--------|
| `99designs/keyring` → `cosmos/keyring` | Standard Cosmos pin |
| `gin-gonic/gin` | Security CVE pin |
| `syndtr/goleveldb` | Store stability pin |
| `go-merkletree/v2` → MarstonConnell fork | **Audit each phase** — storage proofs |

---

## Jackal-only patches

| Patch | Legacy location | Modern location | Phase 1 status |
|-------|-----------------|-----------------|----------------|
| Free post-proof / attest / report / request-attestation fees | SDK fork `x/auth/ante/fee.go` | `app/ante_fee.go` | Ported |
| Custom wasm gas register | `app/wasm_config.go` | `app/wasm_config.go` | Ported |
| Storage int64 overflow on `FileSize * MaxProofs` | Unbounded multiply | `x/storage/keeper/mulStorageCharge` | Fixed |
| CometBFT fork | `TheMarstonConnell/cometbft` | Upstream `cometbft` 0.37.15 | Removed |

---

## Custom modules (`x/`)

| Module | Store key | Purpose |
|--------|-----------|---------|
| `storage` | `storage` | Deals, providers, proofs, payments, attestations |
| `filetree` | `filetree` | File ACLs, viewers/editors, encrypted tree |
| `rns` | `rns` | Name service, bids, marketplace |
| `oracle` | `oracle` | Price / data feeds |
| `jklmint` | `mint` (custom inflation) | Jackal inflation schedule |
| `notifications` | `notifications` | RNS-linked notifications |

All six modules: **gRPC `Msg` + `Query` only** (legacy `Route` / `Querier` removed in Phase 1).

---

## Standard + IBC modules (migration branch)

| Category | Modules |
|----------|---------|
| Core SDK | auth, bank, staking, slashing, mint (jklmint replaces default mint), distr, gov, params, consensus, crisis, evidence, feegrant, authz, upgrade |
| IBC v7 | core IBC, transfer, ICA (host + controller), **IBC fee middleware** |
| Wasm | wasmd `x/wasm` + custom `wasmbinding` (storage, filetree, notifications) |

---

## Upgrade handler registry (mainnet)

Handlers registered in `app/upgrades.go` → `registerMainnetUpgradeHandlers`:

`bouncybulldog`, `v3`, `v4`, `v410`, `v420`, `v430`, `v440`, `v450`, `v460`, `v500`, `v510`, **`v600`** (Phase 1 target).

**`v600` stores added:** `consensus`, `crisis`  
**`v600` migrations:** `baseapp.MigrateParams` → `x/consensus` param store

---

## Wasm / CosmWasm

| Item | Notes |
|------|-------|
| wasmvm (migration branch) | **v1.5.9** — lib path: `internal/api/libwasmvm.x86_64.so` |
| wasmd | v0.45.0 — legacy gov wasm proposals still enabled (`EnableAllProposals`) |
| Custom bindings | `wasmbinding/` — storage, filetree, notifications message plugins |
| Mainnet code IDs | **Captured 2026-06-10** — see `docs/inventory/captured-mainnet-20260610T2333Z.json` (height **18380416**) |

| Code ID | Creator | Data hash |
|---------|---------|-----------|
| 1 | jkl19j955sucqvyyk4l2cdxe3xfgy54qluf49pddnd | `0ADF4677…17E01E` |
| 2 | jkl1qr9g68c4kgy00meyy4tk2yycpw5r6us3vccy7v | `CA0C4BF7…B55DD8` |
| 3 | jkl1qr9g68c4kgy00meyy4tk2yycpw5r6us3vccy7v | `56CCC5F3…10B9B9` |
| 4 | jkl1gues8xhxwcp7k76dpm86dh3z4ac08gxkx4t5jh | `E0A23691…4B70F8` |
| 5 | jkl1gues8xhxwcp7k76dpm86dh3z4ac08gxkx4t5jh | `A30BB174…F0F425D` |

> Action: export `canined query wasm list-code` from mainnet and attach code IDs + checksums before testnet wasm smoke tests.

---

## IBC

| Item | Notes |
|------|-------|
| ibc-go major | v4 (mainnet) → **v7** (migration) |
| Middleware | IBC fee module wired in `app/app.go` |
| Connected chains | **157 channels** at height 18380416 — **93 OPEN** (`transfer`: 60, `icahost`: 33). Full list in captured JSON. Counterparties include Archway (ICA), Osmosis, and others — filter `STATE_OPEN` in inventory file. |
| Relayer | **TBD** — pin Hermes / Go relayer version compatible with ibc-go v7 (recommend Hermes ≥ 1.8 for ibc-go v7). Testnet capture pending (Jackal testnet REST timed out from this environment). |

---

## Operations checklist (Phase 0 exit)

- [x] Inventory documented (this file)
- [ ] Tag mainnet release binary (`v5.1.x`) and archive `go.mod` + replaces
- [x] App export/import round-trip tested (`app.TestWasmdExport` with CGO)
- [ ] State export tested at **current mainnet height** (`canined export --height <H>`)
- [x] Mainnet wasm code ID list captured (`scripts/capture-chain-inventory.sh` → `docs/inventory/captured-mainnet-20260610T2333Z.json`)
- [x] IBC channel inventory captured (mainnet; relayer version still TBD)
- [x] Validator communication plan drafted ([V600-MAINNET-GOVERNANCE.md](./V600-MAINNET-GOVERNANCE.md))

---

## References

- [COSMOS-MODERNIZATION.md](./COSMOS-MODERNIZATION.md)
- [V600-TESTNET-UPGRADE.md](./V600-TESTNET-UPGRADE.md)
- [V600-MAINNET-GOVERNANCE.md](./V600-MAINNET-GOVERNANCE.md)
- [NETWORK-ENDPOINTS.md](./NETWORK-ENDPOINTS.md)
