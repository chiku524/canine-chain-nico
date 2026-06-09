# v600 mainnet governance proposal

Template for the **mainnet** upgrade to Cosmos SDK 0.47 after testnet bake (≥ 2 weeks stable).

> Do **not** submit until [V600-TESTNET-UPGRADE.md](./V600-TESTNET-UPGRADE.md) exit criteria are met.

---

## Proposal summary

| Field | Value |
|-------|-------|
| Type | `software-upgrade` |
| Upgrade name | **`v600`** (must match `app/upgrades/v600/upgrades.go`) |
| Chain ID | `jackal-1` |
| Pre-upgrade line | v5.1.x (`v510` handler) |
| Post-upgrade stack | SDK 0.47.17, CometBFT 0.37.15, ibc-go v7.10, wasmd v0.45, wasmvm v1.5.9 |

---

## Title (example)

```
Jackal v600: Cosmos SDK 0.47 / CometBFT 0.37 migration
```

---

## Description (template)

```markdown
## Summary

This proposal schedules the **v600** on-chain upgrade for Jackal mainnet (`jackal-1`),
migrating the chain from Cosmos SDK 0.45 to **0.47**, CometBFT 0.34 to **0.37**,
ibc-go v4 to **v7**, and wasmd v0.32 to **v0.45**.

## Upgrade details

- **Upgrade name:** v600
- **Halt height:** <HALT_HEIGHT> (approx. <UTC_DATETIME>)
- **Binary:** canined v6.0.0-rc1 (or tagged release)
- **SHA256 (linux/amd64):** <CHECKSUM>
- **Download:** <RELEASE_URL>

## Validator instructions

1. Install wasmvm **v1.5.9** (`libwasmvm.x86_64.so`).
2. Install the v600 binary via Cosmovisor or manual swap before halt height.
3. Ensure ≥ 2× disk space; snapshot at H−1000.
4. Full playbook: https://github.com/JackalLabs/canine-chain/blob/feat/cosmos-modernization-phase1/docs/V600-TESTNET-UPGRADE.md

## Notable changes

- New `x/consensus` module; params migrated from legacy subspace
- Jackal free-fee storage msgs preserved (`MsgPostProof`, `MsgAttest`, `MsgReport`, `MsgRequestAttestationForm`)
- Storage overflow guard on `FileSize * MaxProofs`
- CometBFT and SDK fork replaces removed

## Risks

- Wasm contracts must be verified on wasmvm 1.5.9 (completed on testnet)
- IBC relayers must support ibc-go v7
- Plan rollback via snapshot if consensus fails

## Deposit

Minimum deposit: <GOV_DEPOSIT> ujkl
```

---

## CLI submission

Replace placeholders before submitting:

```bash
HALT_HEIGHT=<height>
DEPOSIT=1000000000000ujkl   # adjust to current min deposit
FROM=<proposer-key>

canined tx gov submit-proposal software-upgrade v600 \
  --title "Jackal v600: Cosmos SDK 0.47 / CometBFT 0.37 migration" \
  --description "$(cat docs/V600-MAINNET-GOVERNANCE.md | sed -n '/## Description/,/## CLI/p')" \
  --upgrade-height "$HALT_HEIGHT" \
  --deposit "$DEPOSIT" \
  --from "$FROM" \
  --chain-id jackal-1 \
  --gas auto \
  --gas-adjustment 1.3 \
  -y
```

Or use a JSON proposal file for longer descriptions:

```bash
canined tx gov submit-proposal proposal.json --from "$FROM" --chain-id jackal-1 -y
```

---

## Halt height selection

1. Pick target UTC time (low-traffic window; communicate ≥ 2 weeks ahead).
2. Estimate blocks: `blocks_remaining = (target_time - now) / 6s`.
3. `halt_height = current_height + blocks_remaining`.
4. Add buffer for governance voting period end.

```bash
canined status | jq -r '.SyncInfo.latest_block_height'
```

---

## Binary checksum manifest

Publish alongside the release (Cosmovisor `upgrade-info.json`):

```json
{
  "name": "v600",
  "height": <HALT_HEIGHT>,
  "info": "{\"binaries\":{\"linux/amd64\":\"<sha256>:https://github.com/JackalLabs/canine-chain/releases/download/v6.0.0/canined-linux-amd64\",\"linux/arm64\":\"<sha256>:https://...\"}}"
}
```

Generate checksum:

```bash
sha256sum build/canined
```

---

## Communications timeline

| When | Audience | Content |
|------|----------|---------|
| T−14d | Validators + providers | Testnet results summary; proposed halt window |
| T−7d | Public | Governance proposal live; deposit rally |
| T−48h | Validators | Final binary URL + checksum; Cosmovisor instructions |
| T−24h | Providers / indexers | API / grpc breaking changes (if any) |
| T−0 | All | Halt height reached; monitor #validators channel |
| T+72h | All | Post-upgrade all-clear or rollback notice |

---

## Post-upgrade monitoring (48–72h)

- [ ] Block time / missed blocks normal
- [ ] `canined query upgrade applied v600` returns applied
- [ ] Storage proofs flowing; zero-fee txs confirmed on mainnet
- [ ] Wasm contracts executing
- [ ] IBC channels healthy
- [ ] Explorers / indexers caught up

---

## References

- [COSMOS-MODERNIZATION.md](./COSMOS-MODERNIZATION.md)
- [V600-TESTNET-UPGRADE.md](./V600-TESTNET-UPGRADE.md)
- [PHASE0-INVENTORY.md](./PHASE0-INVENTORY.md)
