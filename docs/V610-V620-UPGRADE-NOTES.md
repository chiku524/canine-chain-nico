# v610 and v620 upgrade notes

Supplement to [V600-TESTNET-UPGRADE.md](./V600-TESTNET-UPGRADE.md) and [V630-TESTNET-UPGRADE.md](./V630-TESTNET-UPGRADE.md) for the **middle hops** in the modernization ladder.

**Full ladder:** `v600` (0.47) → `v610` (0.50) → `v620` (0.53) → `v630` (0.54)

---

## v610 — SDK 0.50 / ibc-go v8 / wasmd 0.53

| Field | Value |
|-------|-------|
| Branch | `feat/cosmos-modernization-phase2` |
| Handler | `app/upgrades/v610/upgrades.go` |
| Go | 1.23.8 |
| wasmvm | **v2.x** |
| Stores **added** | `circuit` |

### Validator actions

1. Install wasmvm **v2** lib (`releases/download/v2.x.x/libwasmvm.x86_64.so`).
2. Build from `feat/cosmos-modernization-phase2` or use merged artifact after upstream integration.
3. Cosmovisor path: `upgrades/v610/bin/canined`.

### What changes in the app

- Capability module from `ibc-go/modules/capability` (still present at v610).
- Gov v1-only; legacy wasm gov proposals removed.
- cosmossdk.io x-modules (evidence, feegrant, upgrade, circuit).

### Smoke focus

- IBC transfer + ICA still work on ibc-go v8.
- Wasm contracts on **wasmvm v2** (re-test all pinned code IDs).
- Circuit module present (removed again at v630).

```bash
UPGRADE_NAME=v610 UPGRADE_HEIGHT=<H> ./scripts/submit-upgrade-proposal.sh
```

---

## v620 — SDK 0.53 / ibc-go v10 / wasmd 0.60

| Field | Value |
|-------|-------|
| Branch | `feat/cosmos-modernization-phase3` |
| Handler | `app/upgrades/v620/upgrades.go` |
| Go | 1.23.8 |
| wasmvm | **v2.2.x** |
| Stores **deleted** | `capability`, `feeibc` |

### Validator actions

1. wasmvm **v2.2+** lib (match `go.mod`).
2. **No IBC fee middleware** after upgrade — relayers must not expect fee module.
3. Cosmovisor path: `upgrades/v620/bin/canined`.

### What changes in the app

- **x/capability removed** — ibc-go v10 manages port routing differently.
- **IBC fee module removed** — delete `feeibc` store at upgrade.
- IBC callbacks middleware, RouterV2, wasm `NodeConfig`.

### Smoke focus

- IBC wasm contracts (IBC entrypoints) — security patches in wasmd 0.60.1.
- Storage proofs + filetree ACL (regression suite).
- No `query ibc-fee` — module gone.

```bash
UPGRADE_NAME=v620 UPGRADE_HEIGHT=<H> ./scripts/submit-upgrade-proposal.sh
```

---

## Sequential rehearsal (devnet)

When Jackal provides a **state export** from mainnet or testnet:

| Step | Halt at | Binary branch | Upgrade name |
|------|---------|---------------|--------------|
| 1 | H₁ | `master` / v600 | `v600` |
| 2 | H₂ | phase2 | `v610` |
| 3 | H₃ | phase3 | `v620` |
| 4 | H₄ | phase4 | `v630` |

Allow **≥ 1 week** soak between hops on devnet; **≥ 2 weeks** recommended.

Document halt heights and block hashes in `docs/inventory/devnet-upgrade-log.md` (create during rehearsal).

---

## References

- [COSMOVISOR-LADDER.md](./COSMOVISOR-LADDER.md)
- [COSMOS-MODERNIZATION.md](./COSMOS-MODERNIZATION.md)
