# Jackal devnet testing — team handoff

**Purpose:** Brief for conversations with Jackal Labs / core developers about validating the Cosmos modernization fork on a **Jackal-operated devnet** (or coordinated testnet) before any mainnet governance.

**Fork repo:** [github.com/chiku524/canine-chain-nico](https://github.com/chiku524/canine-chain-nico)  
**Target branch:** `feat/cosmos-modernization-phase4` (SDK **0.54.3** — north-star stack)  
**Last updated:** 2026-06-08

---

## Executive summary (talking points)

We have completed a **four-phase fork-only migration** of `canine-chain` from the current mainnet baseline (**SDK 0.45 / ibc-go v4 / wasmd 0.32**) to the **2026.1 Cosmos release family** (**SDK 0.54 / CometBFT 0.39 / ibc-go v11 / wasmd 0.70 / wasmvm v3 / store/v2**).

Work so far is **code + CI + unit tests** on our fork. We deliberately deferred live chain testing until **0.54** was reached. We are now ready to collaborate on **real network validation** — ideally on infrastructure Jackal already operates (devnet or public testnet) with validator and storage-provider participation.

We are **not** asking Jackal to merge blindly. We want a structured bake: export/import, sequential upgrades, wasm contract regression, storage proofs, filetree/RNS/oracle flows, and IBC where applicable.

---

## What we built (fork)

| Phase | On-chain name | Branch | SDK | wasmd | wasmvm | ibc-go | Status |
|-------|---------------|--------|-----|-------|--------|--------|--------|
| 1 | `v600` | `master` | 0.47 | 0.45 | 1.5.x | v7 | Merged; sim CI green |
| 2 | `v610` | `feat/cosmos-modernization-phase2` | 0.50 | 0.53.3 | 2.1.x | v8 | Pushed |
| 3 | `v620` | `feat/cosmos-modernization-phase3` | 0.53 | 0.60.1 | 2.2.x | v10 | Pushed |
| 4 | `v630` | `feat/cosmos-modernization-phase4` | **0.54.3** | **0.70** | **3.x** | **v11** | **Pushed — test gate** |

**Full upgrade ladder (proposed on-chain):**

```
jackal-1 today  →  v600  →  v610  →  v620  →  v630
   (0.45)         (0.47)   (0.50)   (0.53)   (0.54)
```

### Jackal-specific behavior preserved

- **Free post-proof / attest / report fees** — `app/ante_fee.go` (no SDK fork)
- **Custom wasm gas register** — `app/wasm_config.go`
- **Custom modules unchanged in scope:** `storage`, `filetree`, `rns`, `oracle`, `jklmint`, `notifications`
- **Wasmbinding** — storage, filetree, notifications custom messages
- **Storage merkletree** — `TheMarstonConnell/go-merkletree/v2` replace retained (needs proof regression on devnet)

### Store deletions per upgrade (validator ops)

| Upgrade | Stores removed | Reason |
|---------|----------------|--------|
| `v620` | `capability`, `feeibc` | ibc-go v10 removes capability + IBC fee module |
| `v630` | `crisis`, `circuit` | SDK 0.54 / wasmd 0.70; circuit not compatible with SDK 0.54 API set |

---

## What we need from Jackal

### 1. Devnet / testnet access

- [ ] **Chain ID** and RPC/gRPC endpoints for a non-mainnet environment we can reset or fork
- [ ] **Genesis or state export** at a representative height (or permission to run our own `jackal-nico-1`-style net with Jackal validators observing)
- [ ] **2–3 validator seats** (or cosigner agreement) for upgrade rehearsal
- [ ] **Storage provider test nodes** (at least one miner/provider) for proof and deal flows

### 2. Operational alignment

- [ ] Agreed **halt heights** and **upgrade names** (`v600` … `v630`) — or Jackal-preferred naming if different
- [ ] **Binary distribution** — who builds/signs `canined` for devnet (us vs Jackal CI)
- [ ] **Cosmovisor** / rollback plan per hop
- [ ] **CometBFT 0.39** config review (libp2p, sync mode) for validator runbooks

### 3. Wasm & IBC

- [ ] Confirm **mainnet code IDs** to regression-test on **wasmvm v3** (inventory: `docs/inventory/captured-mainnet-20260610T2333Z.json`)
- [ ] Pinned contracts / production-critical wasm — priority list from Jackal
- [ ] **IBC**: which channels/counterparties matter for testnet; relayer version for ibc-go v11
- [ ] Whether **IBC v2 / callbacks** paths need inclusion in smoke tests

### 4. Governance & timeline

- [ ] Jackal’s preferred path: **sequential mainnet upgrades** vs **single jump** after devnet (we recommend **no skip** — one upgrade name per SDK hop)
- [ ] Minimum **soak duration** on devnet/testnet before mainnet prop (we suggest ≥2 weeks per hop, ≥4 weeks on `v630`)
- [ ] Communication template for validators and storage providers

---

## What we bring

| Deliverable | Location |
|-------------|----------|
| Modernized codebase (0.54) | `feat/cosmos-modernization-phase4` |
| Roadmap & checklists | [COSMOS-MODERNIZATION.md](./COSMOS-MODERNIZATION.md) |
| Mainnet wasm/IBC inventory | `docs/inventory/captured-mainnet-20260610T2333Z.json` |
| Private testnet playbook | [PRIVATE-TESTNET.md](./PRIVATE-TESTNET.md) |
| v600 upgrade playbooks (templates) | [V600-TESTNET-UPGRADE.md](./V600-TESTNET-UPGRADE.md), [V600-MAINNET-GOVERNANCE.md](./V600-MAINNET-GOVERNANCE.md) |
| Smoke scripts | `scripts/smoke-v600-testnet.sh`, `scripts/verify-v600-candidate.sh`, `scripts/init-nico-testnet.sh` |
| Upgrade handlers | `app/upgrades/v600`, `v610`, `v620`, `v630` |

### Build requirements (0.54 binary)

- **Go 1.25.9+**
- **CGO enabled**
- **wasmvm v3** shared library:

```sh
WASMVM_VERSION=$(go list -m -f '{{.Version}}' github.com/CosmWasm/wasmvm/v3)
sudo wget -q "https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm.x86_64.so" \
  -O /usr/lib/libwasmvm.x86_64.so
make install
```

Docker build: see `Dockerfile.linux-build` and [WINDOWS-DEV.md](./WINDOWS-DEV.md).

---

## Proposed devnet test plan (6–8 weeks)

### Week 1–2 — Environment & baseline

1. Deploy `v630` binary on devnet **or** start private `jackal-nico-1` with Jackal observers
2. Validator sync, peer connectivity, block production
3. Baseline metrics: block time, tx throughput, memory

### Week 3–4 — Core protocol

| Area | Tests |
|------|--------|
| Storage | Buy storage, post file, **post-proof (zero fee)**, provider init/claim, attest/report |
| Filetree | Post key, post file, add/remove viewers & editors |
| RNS | Register, list, bid, records |
| Oracle | Feed updates, dependent pricing paths |
| Mint | Inflation / block rewards sanity |

### Week 5 — Wasm

- Re-instantiate or migrate **all pinned mainnet code IDs** on wasmvm v3
- Execute + query + IBC-enabled contracts if any
- Custom wasmbinding messages (storage/filetree/notifications)

### Week 6 — Upgrades (if devnet forked from mainnet state)

- Rehearse **`v600` → `v610` → `v620` → `v630`** with halt heights
- State export/import at each hop
- Verify store deletions (`capability`, `feeibc`, `crisis`, `circuit`)

### Week 7–8 — IBC & soak

- Transfer (+ ICA if used)
- Relayer compatibility (ibc-go v11)
- 72h+ soak, incident log, go/no-go for testnet/mainnet discussion

---

## Open questions for Jackal devs

1. Is there an existing **devnet** chain ID, or should we use **`jackal-nico-1`** privately and share artifacts?
2. Can Jackal provide **state export** at a recent testnet/mainnet height for upgrade rehearsal?
3. Which **wasm contracts** are production-critical beyond the five code IDs in our inventory?
4. Are **storage providers** on a separate release cycle that must align with `canined` v630?
5. Is **`x/circuit` removal** at `v630` acceptable, or does Jackal require a circuit-breaker replacement on 0.54?
6. Preferred **PR / upstream merge** process after devnet sign-off?

---

## Suggested meeting agenda (60 min)

1. **10 min** — Fork overview & upgrade ladder (`v600`–`v630`)
2. **15 min** — Store migrations, wasmvm hops, CometBFT 0.39 ops impact
3. **15 min** — Devnet access, validators, storage providers
4. **10 min** — Wasm + IBC regression scope
5. **10 min** — Timeline, governance, upstream merge

---

## Contacts & links

| Item | Value |
|------|-------|
| Fork | `https://github.com/chiku524/canine-chain-nico` |
| Phase 4 branch | `feat/cosmos-modernization-phase4` |
| Upstream | `https://github.com/JackalLabs/canine-chain` |
| Roadmap | [COSMOS-MODERNIZATION.md](./COSMOS-MODERNIZATION.md) |

*Fill in your team contacts below before sending:*

| Role | Name | Email / Discord |
|------|------|-----------------|
| Engineering lead | | |
| Validator liaison | | |
| Storage / provider liaison | | |

---

## Copy-paste email / Discord intro

> Hi Jackal team — we've completed a fork of `canine-chain` bringing the stack from SDK 0.45 to **0.54** (wasmd 0.70, wasmvm v3, ibc-go v11, CometBFT 0.39) across four on-chain upgrade steps (`v600`–`v630`). Code and unit tests are green on our side; we're ready for **devnet/testnet validation** with your validators and at least one storage provider.
>
> We'd like to align on: devnet access (or coordinated private net), wasm contract regression list, IBC/relayer requirements, and a rehearsal schedule for sequential upgrades. Our handoff doc: `docs/JACKAL-DEVNET-HANDOFF.md` on branch `feat/cosmos-modernization-phase4`.
>
> Are you open to a 60-minute technical sync in the next 1–2 weeks?
