# Jackal Cosmos stack modernization

Living roadmap and checklist for bringing **canine-chain** in line with the supported Cosmos release families.

| Field | Value |
|-------|-------|
| **Last updated** | 2026-06-08 |
| **Active branch** | `feat/cosmos-modernization-phase4` |
| **Current phase** | Phase 4 — SDK 0.54 (`v630`) **code complete** → **devnet / testnet validation** |
| **North-star target** | [2026.1 release family](https://docs.cosmos.network/sdk/latest/release-family) (SDK **0.54.x**, not 0.55) |
| **Validation strategy** | Fork code path **complete**; next gate: Jackal devnet + storage providers ([handoff](./JACKAL-DEVNET-HANDOFF.md)) |

---

## Fork-only path to 0.54 (active)

Jackal public testnet/mainnet are **deferred** until the fork reaches SDK **0.54** with green CI. Each phase lands on a feature branch, merges when build + unit tests + sim pass in Docker/Linux CI.

| Phase | On-chain name | SDK | wasmd | wasmvm | ibc-go | Branch |
|-------|---------------|-----|-------|--------|--------|--------|
| 1 ✓ | `v600` | 0.47 | 0.45 | 1.5.x | v7 | `master` |
| 2 ✓ | `v610` | 0.50 | 0.53.3 | 2.1.x | v8 | `feat/cosmos-modernization-phase2` |
| 3 ✓ | `v620` | 0.53 | 0.60.1 | 2.2.x | v10 | `feat/cosmos-modernization-phase3` |
| 4 ✓ | `v630` | 0.54 | 0.70 | 3.x | v11 | `feat/cosmos-modernization-phase4` |

After Phase 4 (code): **Jackal devnet coordination** → private `jackal-nico-1` (optional parallel) → Jackal public testnet → mainnet governance.

See [JACKAL-DEVNET-HANDOFF.md](./JACKAL-DEVNET-HANDOFF.md) for the team brief to Jackal developers.


## How to keep this document updated

**Every PR that touches dependencies, `app/app.go`, upgrade handlers, or custom modules must update this file.**

1. Check off completed items (`[x]`).
2. Bump **Last updated** date.
3. Update **Active branch** / **Current phase** if changed.
4. Add notes under **Phase log** (short bullet: what merged, testnet height, governance prop ID).
5. If dependency pins change, update the version matrix tables.

Agents and reviewers: treat an outdated checklist as incomplete work.

---

## Version matrix

### North star (2026.1 — current supported Cosmos stack)

| Component | Target version | Notes |
|-----------|----------------|-------|
| Cosmos SDK | **0.54.x** | Latest patch in family (e.g. 0.54.3+) |
| CometBFT | **0.39.x** | libp2p, adaptive sync |
| ibc-go | **v11.x** | ICS27-GMP optional |
| wasmd | **v0.70.x** | See [wasmd releases](https://github.com/CosmWasm/wasmd/releases) |
| wasmvm | **v3.x** | Shared lib path changes per major |
| Database | **cosmos-db** | Replaces `cometbft-db` |
| Store | **store/v2** | Required at SDK 0.54 |
| Go | **1.25+** | Required by wasmd 0.70 |

> **Note:** Cosmos SDK **0.55** is not released. The next *release family* after 0.54 will be 0.55.x when announced by Cosmos Labs.

### Mainnet today (pre-migration)

| Component | Version |
|-----------|---------|
| Cosmos SDK | 0.45.17 (`JackalLabs/cosmos-sdk-new` replace) |
| CometBFT | 0.34.27 (`TheMarstonConnell/cometbft` replace) |
| ibc-go | v4 |
| wasmd | v0.32 |
| wasmvm | v1.5.x |
| Go | 1.23 |

### Phase 1 branch (`master` — `v600`)

| Component | Version |
|-----------|---------|
| Cosmos SDK | 0.47.17 |
| CometBFT | 0.37.15 |
| cometbft-db | 0.14.1 |
| ibc-go | v7.10.0 |
| wasmd | v0.45.0 |
| wasmvm | v1.5.9 |
| Go | 1.23.8 |

### Phase 4 branch (`feat/cosmos-modernization-phase4` — `v630`)

| Component | Version |
|-----------|---------|
| Cosmos SDK | 0.54.3 |
| CometBFT | 0.39.3 |
| cosmos-db | 1.1.3 |
| store/v2 | 2.0.0 |
| ibc-go | v11.0.0 |
| wasmd | v0.70.0 |
| wasmvm | v3.0.4 |
| Go | 1.25.9 |

### Upgrade path (on-chain names)

```
0.45 mainnet  →  v600 (0.47)  →  v610 (0.50)  →  v620 (0.53)  →  v630 (0.54)
```

---

## Phase log

Chronological notes; append new entries at the top.

| Date | Phase | Notes |
|------|-------|-------|
| 2026-06-08 | Phase 4 | Docs + [JACKAL-DEVNET-HANDOFF.md](./JACKAL-DEVNET-HANDOFF.md); checklists closed for Phases 2–4 code work; validation gate → devnet. |
| 2026-06-11 | Phase 4 | SDK 0.54.3 / wasmd 0.70.0 / wasmvm v3.0.4 / ibc-go v11 / CometBFT 0.39.3 / store/v2; removed x/crisis + x/circuit; `v630` handler; build + `make test-unit` green in Docker (Go 1.25). |
| 2026-06-11 | Phase 3 | SDK 0.53.5 / wasmd 0.60.1 / wasmvm v2.2.4 / ibc-go v10.5.0; removed x/capability + ibc-fee; `v620` handler; build + `make test-unit` green in Docker. |
| 2026-06-08 | Phase 2 | SDK 0.50.9 / wasmd 0.53.3 / wasmvm v2.1.4 / ibc-go v8; `v610` handler; build + `make test-unit` green in Docker; CI wasmvm v2. |
| 2026-06-11 | Phase 2 | Started fork-only path to 0.54; branch `feat/cosmos-modernization-phase2`; target wasmd 0.53.3 / SDK 0.50.9 / ibc-go v8. |
| 2026-06-11 | Phase 1 | Sim CI green (seeds 2,17,18,20); lint fixes; merged to `master`. |
| 2026-06-10 | Phase 0–1 | `make proto-gen` via `ghcr.io/cosmos/proto-builder:0.13.1`; regenerated `.pb.go` with `cosmos/gogoproto`; Phase 0 inventory + v600 testnet/mainnet playbooks; `scripts/smoke-v600-testnet.sh`; proto-gen CI workflow; `v600` upgrade unit test. |
| 2026-06-09 | Phase 1 | Pushed `feat/cosmos-modernization-phase1`; `go mod tidy`; sim tests migrated off `simapp` → `testutil/sims`; storage `mulStorageCharge` overflow guard; CI: CGO + wasmvm 1.5.9 on Linux; README install section updated. |
| 2026-06-08 | Phase 1 | Branch `feat/cosmos-modernization-phase1`: `go.mod` → 0.47 stack; `app/app.go` wasmd 0.45 rewrite; `v600` handler; free post-proof ante in `app/ante_fee.go`; build green; filetree keeper tests fixed (keyring codec). |
| — | Phase 0 | Inventory started; mainnet at SDK 0.45 / ibc-go v4. |

---

## Global prerequisites

Applies to all phases; check once and re-verify each phase.

- [x] Single source of truth: **this file** + pinned `go.mod` on release branches
- [x] Long-lived migration workflow: phase branches → testnet → mainnet governance (fork path complete; live testnet next)
- [x] **Linux CI with CGO** for wasmvm integration tests (`.github/workflows/test-unit.yml`)
- [ ] Public testnet mirroring mainnet modules + representative wasm contracts
- [x] Upgrade playbook: halt height, binary checksum, rollback, validator comms ([V600-TESTNET-UPGRADE.md](./V600-TESTNET-UPGRADE.md), [V600-MAINNET-GOVERNANCE.md](./V600-MAINNET-GOVERNANCE.md), [NETWORK-ENDPOINTS.md](./NETWORK-ENDPOINTS.md))
- [ ] Remove unnecessary forks:
  - [x] Free post-proof ante → `app/ante_fee.go` (no SDK fork)
  - [x] CometBFT fork → upstream CometBFT (Phase 1)
  - [ ] Audit `TheMarstonConnell/go-merkletree/v2` replace (storage proofs)
- [x] Proto pipeline: `make proto-gen` on Linux (`proto-builder:0.13.1`); generated code uses `cosmos/gogoproto`; CI in `.github/workflows/proto-gen.yml`
- [x] Go version ladder documented per phase (see version matrix + [PHASE0-INVENTORY.md](./PHASE0-INVENTORY.md))

---

## Phase 0 — Baseline & inventory

**Goal:** Know exactly what is being migrated.  
**On-chain upgrade:** none

### Code & dependencies
- [ ] Tag and record mainnet release binary (`v5.1.x`) `go.mod` + all `replace` directives
- [x] List Jackal-only patches (historical SDK fork, CometBFT fork, custom ante, wasm gas) — [PHASE0-INVENTORY.md](./PHASE0-INVENTORY.md)
- [x] Inventory custom modules: `storage`, `filetree`, `rns`, `oracle`, `jklmint`, `notifications`
- [x] Inventory IBC: transfer, ICA, fee middleware, wasm IBC handler
- [x] List mainnet wasm code IDs / pinned contracts (per wasmvm hop) — **`docs/inventory/captured-mainnet-20260610T2333Z.json`**
- [x] Document connected IBC chains (mainnet capture; relayer version still TBD)

### Operations
- [ ] Test state export at current mainnet height (`canined export --height <H>`; app round-trip: `TestWasmdExport`)
- [x] Validator / provider communication plan for upgrade windows — [V600-MAINNET-GOVERNANCE.md](./V600-MAINNET-GOVERNANCE.md)

**Exit criteria:** Inventory complete; export tested; team agrees on phased plan.

---

## Phase 1 — SDK 0.47 / CometBFT 0.37 / ibc-go v7 / wasmd 0.45

**On-chain upgrade name:** `v600`  
**Branch:** merged to `master` (origin: `feat/cosmos-modernization-phase1`)  
**Handler:** `app/upgrades/v600/`

### Dependency pins (target)

| Package | Version |
|---------|---------|
| `github.com/cosmos/cosmos-sdk` | v0.47.17 |
| `github.com/cometbft/cometbft` | v0.37.15 |
| `github.com/cometbft/cometbft-db` | v0.14.1 |
| `github.com/cosmos/ibc-go/v7` | v7.10.0 |
| `github.com/CosmWasm/wasmd` | v0.45.0 |
| `github.com/CosmWasm/wasmvm` | v1.5.9 |

### `go.mod` & imports
- [x] Bump core Cosmos / CometBFT / ibc-go / wasmd versions
- [x] Remove `JackalLabs/cosmos-sdk-new` replace
- [x] Remove `TheMarstonConnell/cometbft` replace
- [x] Mechanical imports: `tendermint` → `cometbft`, `ibc-go/v4` → `v7`
- [x] Run `go mod tidy` safely (pin versions; avoid pulling SDK 0.54 accidentally)

### `app/` layer
- [x] Rewrite `app/app.go` (wasmd v0.45 template + Jackal modules)
- [x] Add `x/consensus` + `ConsensusParamsKeeper` + `SetParamStore`
- [x] Register `v600` in `app/upgrades.go`
- [x] `v600` handler: legacy params subspace → `MigrateParams` → `consensus` store
- [x] Ante handler: wasmd 0.45 chain (`ExtensionOptions`, `TxFeeChecker`, `RedundantRelayDecorator`)
- [x] Jackal free storage fee waiver: `app/ante_fee.go` + `NewJackalDeductFeeDecorator`
- [x] `ExportAppStateAndValidators` 3-arg signature (`runtime.AppI`)
- [x] `go build ./cmd/canined` succeeds
- [x] Full unit tests on **Linux CI with CGO**
- [x] Simulation tests (`app/sim_test.go`, `//go:build simulation`) updated for 0.47 APIs (`testutil/sims`, `simulation/client/cli`)

### `cmd/canined`
- [x] Root command aligned with wasmd 0.47 (`GenesisCoreCommand`, `InterceptConfigsPreRunHandler`)
- [x] `appExport` signature + `DefaultBaseappOptions`
- [x] Remove duplicate `genaccounts.go` (SDK genesis command covers it)

### Custom modules (0.47 API)
- [x] Remove legacy `Route()` / `QuerierRoute()` / `LegacyQuerierHandler()` (all 6 modules)
- [x] Remove deprecated `RandomizedParams` / fix simulation `SafeSub` / encoding config
- [x] Keeper `storetypes.StoreKey` where required
- [x] Bulk `AccAddressFromHex` → `AccAddressFromHexUnsafe` where needed
- [x] `.pb.go` imports: `cosmos/gogoproto/proto` (not `gogo/protobuf`)
- [x] Filetree `MakePrivateKey`: keyring codec with `cryptocodec.RegisterInterfaces`
- [x] Full `make proto-gen` and commit regenerated protos (`ghcr.io/cosmos/proto-builder:0.13.1`)

### Wasmbinding & wasm
- [x] Custom wasm plugins (`owasm`) wired in `app.go`
- [x] Jackal gas register: `app/wasm_config.go` → `wasmtypes.GasRegister`
- [ ] Verify pinned mainnet contracts on wasmvm 1.5.9 (testnet)

### Testing
- [x] `x/filetree/keeper` tests pass
- [x] `x/rns/keeper`, `x/oracle/keeper`, `x/notifications/keeper` pass
- [x] `app` ante fee unit tests pass
- [x] `x/storage/keeper` — `mulStorageCharge` rejects int64 wrap; `TestOverflow_Finding3` documents semantics
- [x] Integration tests requiring CGO (storage/grpc suites) on Linux CI

### Docs & ops
- [x] This roadmap document created
- [x] Update root `README.md` install section (Go version, wasmvm lib path)
- [x] Testnet deploy candidate playbook — [V600-TESTNET-UPGRADE.md](./V600-TESTNET-UPGRADE.md), [NETWORK-ENDPOINTS.md](./NETWORK-ENDPOINTS.md)
- [ ] Testnet deploy candidate `v600` executed on Jackal testnet
- [ ] Testnet smoke: storage post-proof (zero fee), filetree, rns, oracle, wasm execute — `scripts/smoke-v600-testnet.sh`, `scripts/verify-v600-candidate.sh`
- [ ] IBC transfer smoke on testnet
- [ ] Mainnet governance: upgrade name `v600`, halt height, binary checksum — **template:** [V600-MAINNET-GOVERNANCE.md](./V600-MAINNET-GOVERNANCE.md)
- [ ] Post-upgrade monitoring (48–72h) — checklist in governance doc

**Exit criteria:** Testnet stable ≥2 weeks; mainnet `v600` executed without consensus halt.

---

## Phase 2 — SDK 0.50 / CometBFT 0.38 / ibc-go v8 / wasmd 0.53

**On-chain upgrade name:** `v610`  
**Branch:** `feat/cosmos-modernization-phase2` (pushed)  
**Handler:** `app/upgrades/v610/`

### Dependency pins (shipped)

| Package | Version |
|---------|---------|
| `github.com/cosmos/cosmos-sdk` | 0.50.9 |
| `github.com/cometbft/cometbft` | 0.38.11 |
| `github.com/cosmos/ibc-go/v8` | v8.4.0 |
| `github.com/CosmWasm/wasmd` | 0.53.3 |
| `github.com/CosmWasm/wasmvm/v2` | 2.1.4 |

### Major breaking changes
- [x] **Capability** moved to `ibc-go/modules/capability`; rewired IBC + wasm scoped keepers
- [x] **Gov v1 only** — removed legacy `v1beta1` proposal handlers and wasm legacy gov handlers
- [x] **cosmossdk.io** x-modules (evidence, feegrant, upgrade, circuit)
- [x] Rebased `app/app.go` on wasmd 0.53 template

### Checklist
- [x] Branch from `master` (post-`v600`)
- [x] `go.mod` bump + import migration
- [x] `app/upgrades/v610/` handler (`circuit` store added)
- [x] Custom modules: store imports, `context.Context` keeper APIs, handler removal
- [x] Wasmbinding updated for wasmvm v2 Messenger API
- [x] Re-verified `app/ante_fee.go` against new ante APIs
- [x] wasmvm **2.x** lib in Docker / CI / README
- [x] `make test-unit` green in Docker
- [ ] Contract regression on testnet (all pinned code IDs)
- [ ] IBC relayer upgrade + counterparty check
- [ ] Testnet `v610` ≥2 weeks → mainnet governance

**Exit criteria (code):** met on fork. **Live network:** wasm contracts on wasmvm 2.x + testnet bake pending.

---

## Phase 3 — SDK 0.53 / ibc-go v10 (2025.1 family)

**On-chain upgrade name:** `v620`  
**Branch:** `feat/cosmos-modernization-phase3` (pushed)  
**Handler:** `app/upgrades/v620/`

### Dependency pins (shipped)

| Package | Version |
|---------|---------|
| `github.com/cosmos/cosmos-sdk` | 0.53.5 |
| `github.com/cometbft/cometbft` | 0.38.20 |
| `github.com/cosmos/ibc-go/v10` | v10.5.0 |
| `github.com/CosmWasm/wasmd` | 0.60.1 |
| `github.com/CosmWasm/wasmvm/v2` | 2.2.4 |

### Checklist
- [x] SDK 0.50 → 0.53 migration per wasmd 0.60.1 template
- [x] ibc-go v8 → v10; **removed** `x/capability` + IBC fee module
- [x] `app/upgrades/v620/` (deletes `capability`, `feeibc` stores)
- [x] IBC callbacks middleware, RouterV2, wasm NodeConfig
- [x] `make test-unit` green in Docker
- [ ] Storage proof + filetree ACL regression on devnet
- [ ] Performance baseline (block time, proof tx throughput)
- [ ] Testnet `v620` → mainnet governance

**Exit criteria (code):** met on fork. **Live network:** testnet ≥2 weeks pending.

---

## Phase 4 — SDK 0.54 / CometBFT 0.39 / ibc-go v11 / wasmd 0.70 (2026.1 target)

**On-chain upgrade name:** `v630`  
**Branch:** `feat/cosmos-modernization-phase4` (pushed)  
**Handler:** `app/upgrades/v630/`

### Dependency pins (shipped)

| Package | Version |
|---------|---------|
| `github.com/cosmos/cosmos-sdk` | 0.54.3 |
| `github.com/cometbft/cometbft` | 0.39.3 |
| `github.com/cosmos/ibc-go/v11` | v11.0.0 |
| `github.com/CosmWasm/wasmd` | v0.70.0 |
| `github.com/CosmWasm/wasmvm/v3` | v3.0.4 |
| `github.com/cosmos/cosmos-db` | 1.1.3 |
| `github.com/cosmos/cosmos-sdk/store/v2` | 2.0.0 |
| Go | 1.25.9 |

### Major breaking changes (0.54)
- [x] Migrate to **store/v2** (app + all custom modules)
- [x] `cosmossdk.io/log/v2` imports
- [x] Remove **x/crisis** module + store (`v630` deletes store)
- [x] Remove **x/circuit** (no SDK 0.54–compatible release; `v630` deletes store)
- [x] Evidence/feegrant/upgrade moved to in-tree SDK `x/*` modules
- [ ] CometBFT 0.39 config (libp2p, adaptive sync) — **validator runbook on devnet**
- [ ] Evaluate **BlockSTM** for Jackal workload

### Checklist
- [x] Rebased `app/app.go` on wasmd 0.70
- [x] ibc-go v11 middleware wiring
- [x] wasmvm **v3** lib in Docker / CI / README
- [x] `app/upgrades/v630/`
- [x] Ante / fee / wasm gas preserved (`app/ante_fee.go`, `app/wasm_config.go`)
- [x] CI: Go 1.25, wasmvm v3, CGO build workflows
- [x] `make test-unit` green in Docker (`-p 1` for isolation)
- [ ] Contract full regression on devnet (all pinned code IDs)
- [ ] OpenAPI / grpc-gateway regen (optional before mainnet)
- [ ] External review recommended before mainnet
- [ ] Devnet / testnet ≥4 weeks → mainnet `v630`
- [ ] 30-day post-upgrade monitoring

**Exit criteria (code):** Jackal on **2026.1** stack in fork; unit tests green. **Live network:** devnet handoff → testnet bake → governance.

---

## Recurring PR checklist (every dependency bump)

Copy into PR descriptions or use as self-review.

### `go.mod`
- [ ] Exact version pins for cosmos-sdk, cometbft, ibc-go, wasmd, wasmvm
- [ ] Obsolete `replace` directives removed
- [ ] `go` / `toolchain` meets wasmd minimum
- [ ] **This document** updated

### Mechanical imports
- [ ] `tendermint/*` → `cometbft/*` (if any remain)
- [ ] `ibc-go/vN` → correct major
- [ ] `gogo/protobuf` → `cosmos/gogoproto` in generated code
- [ ] `cometbft-db` → `cosmos-db` (Phase 4+)
- [ ] `store` → `store/v2` (Phase 4)

### `app/app.go` (diff against target wasmd tag)
- [ ] Keeper constructors and store keys
- [ ] Module manager order (begin/end block, init genesis)
- [ ] IBC router + middleware
- [ ] Gov + upgrade keepers
- [ ] Ante / post-handler
- [ ] Wasm keeper + custom options

### Jackal custom modules

| Module | Focus areas |
|--------|-------------|
| `x/storage` | Proofs, deals, providers, free-fee msgs, merkletree |
| `x/filetree` | Keyring/crypto, viewers/editors, ACL queries |
| `x/rns` | Names, bids, records, marketplace |
| `x/oracle` | Feeds, params |
| `x/jklmint` | Inflation, mint params |
| `x/notifications` | RNS-linked notifications |

- [ ] `module.go` — no legacy routes; `RegisterServices` only
- [ ] Migrations registered per consensus version
- [ ] Keeper + CLI tests pass on Linux

### Upgrades
- [ ] New `app/upgrades/vXXX/upgrades.go`
- [ ] Registered in `app/upgrades.go`
- [ ] `StoreUpgrades` (add / delete / rename stores)
- [ ] Tested with exported genesis on testnet

### Ecosystem
- [ ] Release binaries (linux amd64/arm64)
- [ ] Provider / storage node compatibility
- [ ] Explorers / indexers notified
- [ ] Relayer version pinned
- [ ] Cosmovisor / upgrade docs

---

## Testing matrix

| Layer | Action |
|-------|--------|
| Unit | `make test-unit` (Linux, CGO, Go 1.25 on phase4, `-tags='ledger test_ledger_mock test'`, `-p 1`) |
| Build | `go build ./cmd/canined` |
| App export | Export / import at height H-1 |
| Wasm | Instantiate + execute pinned contracts |
| IBC | Transfer + wasm or ICA path |
| Storage | `MsgPostProof` with zero fee (ante waiver) |
| Filetree | Post key, add/remove viewers |
| Upgrade | Run handler on testnet from pre-upgrade export |

---

## Risk register

| Risk | Mitigation |
|------|------------|
| Wasm contracts break on wasmvm major hop | Per-phase testnet regression on all code IDs; keep rollback binary |
| Storage proofs after store migrations | Dedicated proof test suite; merkletree lib pin |
| IBC channels during long migration | Coordinate counterparties; document relayer versions per phase |
| Lost free post-proof fee waiver | `app/ante_fee_test.go` + testnet zero-fee proof txs |
| Validator ops on CometBFT 0.39 | Early validator testnet; config migration guide |
| Skipping phases | **Do not** jump 0.47 → 0.54 in one upgrade |

---

## Indicative timeline

| Phase | Engineering (fork) | Testnet / devnet bake | Mainnet upgrade |
|-------|-------------------|----------------------|-----------------|
| 0 — Inventory | ✓ done | — | — |
| 1 — 0.47 (`v600`) | ✓ done | pending Jackal devnet | `v600` |
| 2 — 0.50 (`v610`) | ✓ done | pending | `v610` |
| 3 — 0.53 (`v620`) | ✓ done | pending | `v620` |
| 4 — 0.54 (`v630`) | ✓ done | **next** (6–8 weeks) | `v630` |

**Total (conservative):** fork engineering complete; live validation ~3–6 months before mainnet `v630` (sequential hops or agreed schedule with Jackal).

---

## References

- [Cosmos release families](https://docs.cosmos.network/sdk/latest/release-family)
- [SDK 0.54 upgrade guide](https://docs.cosmos.network/sdk/latest/upgrade/upgrade)
- [ibc-go releases](https://github.com/cosmos/ibc-go/releases)
- [wasmd releases](https://github.com/CosmWasm/wasmd/releases)
- [Jackal protocol research (local)](../../jackal-protocol-research.md) — optional context outside repo

---

## Jackal-specific patches (track re-ports)

| Patch | Original location | Modern location | Status |
|-------|-------------------|-----------------|--------|
| Free post-proof / attest / report / request-attestation fees | `JackalLabs/cosmos-sdk-new` `x/auth/ante/fee.go` | `app/ante_fee.go` | Done (Phase 1) |
| Custom wasm gas (compile/instance costs) | `app/wasm_config.go` | `app/wasm_config.go` | Done (Phase 1) |
| Storage merkletree | `TheMarstonConnell/go-merkletree/v2` replace | `go.mod` replace | Keep — verify each phase |
| CometBFT fork | `TheMarstonConnell/cometbft` | Upstream cometbft | Removed Phase 1 |
