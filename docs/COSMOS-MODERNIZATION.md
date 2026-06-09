# Jackal Cosmos stack modernization

Living roadmap and checklist for bringing **canine-chain** in line with the supported Cosmos release families.

| Field | Value |
|-------|-------|
| **Last updated** | 2026-06-09 |
| **Active branch** | `feat/cosmos-modernization-phase1` |
| **Current phase** | Phase 1 — SDK 0.47 (`v600`) |
| **North-star target** | [2026.1 release family](https://docs.cosmos.network/sdk/latest/release-family) (SDK **0.54.x**, not 0.55) |

---

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

### Migration branch (`feat/cosmos-modernization-phase1`)

| Component | Version |
|-----------|---------|
| Cosmos SDK | 0.47.17 |
| CometBFT | 0.37.15 |
| cometbft-db | 0.14.1 |
| ibc-go | v7.10.0 |
| wasmd | v0.45.0 |
| wasmvm | v1.5.9 |
| Go | 1.23.8 |

### Upgrade path (on-chain names)

```
0.45 mainnet  →  v600 (0.47)  →  v610 (0.50)  →  v620 (0.53)  →  v630 (0.54)
```

---

## Phase log

Chronological notes; append new entries at the top.

| Date | Phase | Notes |
|------|-------|-------|
| 2026-06-09 | Phase 1 | Pushed `feat/cosmos-modernization-phase1`; `go mod tidy`; sim tests migrated off `simapp` → `testutil/sims`; storage `mulStorageCharge` overflow guard; CI: CGO + wasmvm 1.5.9 on Linux; README install section updated. |
| 2026-06-08 | Phase 1 | Branch `feat/cosmos-modernization-phase1`: `go.mod` → 0.47 stack; `app/app.go` wasmd 0.45 rewrite; `v600` handler; free post-proof ante in `app/ante_fee.go`; build green; filetree keeper tests fixed (keyring codec). |
| — | Phase 0 | Inventory started; mainnet at SDK 0.45 / ibc-go v4. |

---

## Global prerequisites

Applies to all phases; check once and re-verify each phase.

- [ ] Single source of truth: **this file** + pinned `go.mod` on release branches
- [ ] Long-lived migration workflow: phase branches → testnet → mainnet governance
- [x] **Linux CI with CGO** for wasmvm integration tests (`.github/workflows/test-unit.yml`)
- [ ] Public testnet mirroring mainnet modules + representative wasm contracts
- [ ] Upgrade playbook: halt height, binary checksum, rollback, validator comms
- [ ] Remove unnecessary forks:
  - [x] Free post-proof ante → `app/ante_fee.go` (no SDK fork)
  - [ ] CometBFT fork → upstream CometBFT
  - [ ] Audit `TheMarstonConnell/go-merkletree/v2` replace (storage proofs)
- [ ] Proto pipeline: `make proto-gen` on Linux; generated code uses `cosmos/gogoproto`
- [ ] Go version ladder documented per phase (1.23 → 1.24 → 1.25+)

---

## Phase 0 — Baseline & inventory

**Goal:** Know exactly what is being migrated.  
**On-chain upgrade:** none

### Code & dependencies
- [ ] Tag and record mainnet release binary (`v5.1.x`) `go.mod` + all `replace` directives
- [ ] List Jackal-only patches (historical SDK fork, CometBFT fork, custom ante, wasm gas)
- [ ] Inventory custom modules: `storage`, `filetree`, `rns`, `oracle`, `jklmint`, `notifications`
- [ ] Inventory IBC: transfer, ICA, fee middleware, wasm IBC handler
- [ ] List mainnet wasm code IDs / pinned contracts (per wasmvm hop)
- [ ] Document connected IBC chains and relayer versions

### Operations
- [ ] Test state export at current mainnet height
- [ ] Validator / provider communication plan for upgrade windows

**Exit criteria:** Inventory complete; export tested; team agrees on phased plan.

---

## Phase 1 — SDK 0.47 / CometBFT 0.37 / ibc-go v7 / wasmd 0.45

**On-chain upgrade name:** `v600`  
**Branch:** `feat/cosmos-modernization-phase1`  
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
- [ ] Full unit tests on **Linux CI with CGO**
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
- [ ] Full `make proto-gen` and commit regenerated protos (replace manual patches)

### Wasmbinding & wasm
- [x] Custom wasm plugins (`owasm`) wired in `app.go`
- [x] Jackal gas register: `app/wasm_config.go` → `wasmtypes.GasRegister`
- [ ] Verify pinned mainnet contracts on wasmvm 1.5.9 (testnet)

### Testing
- [x] `x/filetree/keeper` tests pass
- [x] `x/rns/keeper`, `x/oracle/keeper`, `x/notifications/keeper` pass
- [x] `app` ante fee unit tests pass
- [x] `x/storage/keeper` — `mulStorageCharge` rejects int64 wrap; `TestOverflow_Finding3` documents semantics
- [ ] Integration tests requiring CGO (storage/grpc suites) on Linux CI

### Docs & ops
- [x] This roadmap document created
- [x] Update root `README.md` install section (Go version, wasmvm lib path)
- [ ] Testnet deploy candidate `v600`
- [ ] Testnet smoke: storage post-proof (zero fee), filetree, rns, oracle, wasm execute
- [ ] IBC transfer smoke on testnet
- [ ] Mainnet governance: upgrade name `v600`, halt height, binary checksum
- [ ] Post-upgrade monitoring (48–72h)

**Exit criteria:** Testnet stable ≥2 weeks; mainnet `v600` executed without consensus halt.

---

## Phase 2 — SDK 0.50 / CometBFT 0.38 / ibc-go v8 / wasmd ~0.50

**On-chain upgrade name:** `v610` (proposed)  
**Branch:** `feat/cosmos-modernization-phase2` (create when Phase 1 mainnet is stable)

### Dependency pins (target — verify against wasmd release notes at kickoff)

| Package | Version |
|---------|---------|
| `github.com/cosmos/cosmos-sdk` | 0.50.x |
| `github.com/cometbft/cometbft` | 0.38.x |
| `github.com/cosmos/ibc-go/v8` | v8.x |
| `github.com/CosmWasm/wasmd` | ~0.50.x |
| `github.com/CosmWasm/wasmvm` | 2.x |

### Major breaking changes
- [ ] **Remove `x/capability`** — largest structural change; rewire IBC + wasm scoped keepers
- [ ] **Gov v1 only** — remove legacy `v1beta1` proposal handlers and wasm `NewLegacyWasmProposalHandler`
- [ ] Continue **params module** migration off legacy subspaces
- [ ] Rebase `app/app.go` on wasmd 0.50 template

### Checklist
- [ ] New branch from post-`v600` mainnet tag
- [ ] `go.mod` bump + import migration
- [ ] `app/upgrades/v610/` handler + store migrations
- [ ] Custom modules: keeper constructors, `RegisterServices`, consensus versions
- [ ] Wasmbinding + `owasm` plugins updated
- [ ] Re-verify `app/ante_fee.go` against new ante APIs
- [ ] wasmvm **2.x** lib in Docker / install docs
- [ ] Contract regression on testnet (all pinned code IDs)
- [ ] IBC relayer upgrade + counterparty check
- [ ] Testnet `v610` ≥2 weeks → mainnet governance

**Exit criteria:** No `x/capability` in code or required state; wasm contracts verified on wasmvm 2.x.

---

## Phase 3 — SDK 0.53 / ibc-go v10 (2025.1 family)

**On-chain upgrade name:** `v620` (proposed)

### Dependency pins (target)

| Package | Version |
|---------|---------|
| `github.com/cosmos/cosmos-sdk` | 0.53.x |
| `github.com/cometbft/cometbft` | 0.38.x |
| `github.com/cosmos/ibc-go/v10` | v10.x |
| `github.com/CosmWasm/wasmd` | ~0.61.x (confirm matrix) |
| `github.com/CosmWasm/wasmvm` | 2.x → 3.x (per wasmd tag) |

### Checklist
- [ ] Follow [SDK upgrade guide](https://docs.cosmos.network/sdk/latest/upgrade/upgrade) 0.50 → 0.53
- [ ] ibc-go v8 → v10 (intermediate v9 if required by guide)
- [ ] `app/upgrades/v620/`
- [ ] Storage proof + filetree ACL regression suite
- [ ] Performance baseline (block time, proof tx throughput)
- [ ] Testnet `v620` → mainnet governance

**Exit criteria:** Matches 2025.1 release family on testnet ≥2 weeks.

---

## Phase 4 — SDK 0.54 / CometBFT 0.39 / ibc-go v11 / wasmd 0.70 (2026.1 target)

**On-chain upgrade name:** `v630` (proposed)

### Dependency pins (target)

| Package | Version |
|---------|---------|
| `github.com/cosmos/cosmos-sdk` | 0.54.x |
| `github.com/cometbft/cometbft` | 0.39.x |
| `github.com/cosmos/ibc-go/v11` | v11.x |
| `github.com/CosmWasm/wasmd` | v0.70.x |
| `github.com/CosmWasm/wasmvm` | v3.x |
| `github.com/cosmos/cosmos-db` | 1.x |
| `github.com/cosmos/cosmos-sdk/store/v2` | 2.x |
| Go | 1.25+ |

### Major breaking changes (0.54)
- [ ] Migrate to **store/v2**
- [ ] `cosmossdk.io/log/v2` imports
- [ ] Remove **x/crisis** module + store
- [ ] Drop legacy gov / unused modules per [0.54 release notes](https://docs.cosmos.network/sdk/next/upgrade/release)
- [ ] CometBFT 0.39 config (libp2p, adaptive sync)
- [ ] Evaluate **BlockSTM** for Jackal workload

### Checklist
- [ ] Rebase `app/app.go` on wasmd 0.70
- [ ] ibc-go v11 middleware wiring
- [ ] wasmvm **v3** lib + contract full regression
- [ ] `app/upgrades/v630/`
- [ ] Final ante / fee / wasm gas audit
- [ ] OpenAPI / grpc-gateway regen
- [ ] External review recommended before mainnet
- [ ] Testnet ≥4 weeks → mainnet `v630`
- [ ] 30-day post-upgrade monitoring

**Exit criteria:** Jackal on **2026.1** stack; CI green; docs and operator tooling updated.

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
| Unit | `make test-unit` (Linux, CGO, `-tags='ledger test_ledger_mock'`) |
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

| Phase | Engineering + testnet bake | Mainnet upgrade |
|-------|---------------------------|-----------------|
| 0 — Inventory | 1–2 weeks | — |
| 1 — 0.47 (`v600`) | 4–8 weeks (in progress) | `v600` |
| 2 — 0.50 (`v610`) | 8–12 weeks | `v610` |
| 3 — 0.53 (`v620`) | 6–10 weeks | `v620` |
| 4 — 0.54 (`v630`) | 8–12 weeks | `v630` |

**Total (conservative):** ~12–18 months with testnet gates between mainnet upgrades.

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
