# v600 testnet upgrade playbook

Deploy and validate the **Cosmos SDK 0.47** migration (`v600`) on a Jackal testnet **before** mainnet governance.

| Field | Value |
|-------|-------|
| Upgrade name | `v600` |
| Branch | `master` |
| Handler | `app/upgrades/v600/upgrades.go` |
| Binary requirements | Go 1.23.8, CGO, wasmvm **v1.5.9** |

---

## 1. Build release candidate

```bash
git checkout master
git pull

# Linux amd64 (validators)
CGO_ENABLED=1 WASMVM_TAG=v1.5.9 make build-linux

# Verify
./build/canined version
ldd ./build/canined | grep wasmvm
```

Publish the binary + SHA256 checksum in the testnet upgrade announcement.

---

## 2. Schedule upgrade (governance or coordinator)

### Option A — on-chain software upgrade proposal

```bash
canined tx gov submit-proposal software-upgrade v600 \
  --title "Jackal v600: Cosmos SDK 0.47 migration" \
  --description "Testnet migration to SDK 0.47 / CometBFT 0.37 / ibc-go v7 / wasmd 0.45. See docs/V600-TESTNET-UPGRADE.md." \
  --upgrade-height <HALT_HEIGHT> \
  --deposit 1000000000000ujkl \
  --from <authority> \
  --chain-id <TESTNET_CHAIN_ID> \
  --gas auto
```

### Option B — coordinated halt (testnet only)

Validators agree on halt height and run Cosmovisor with the upgrade name `v600` in `upgrade-info.json`:

```json
{
  "name": "v600",
  "height": <HALT_HEIGHT>,
  "info": "{\"binaries\":{\"linux/amd64\":\"<sha256>:<download-url>\"}}"
}
```

---

## 3. Validator prep

Each validator / sentry:

1. Install wasmvm **v1.5.9** shared library (see [README](../README.md)).
2. Install Cosmovisor (recommended) or plan manual binary swap at halt height − 1.
3. Place the `v600` binary where Cosmovisor expects `upgrades/v600/bin/canined`.
4. Ensure **≥ 2× disk** free for upgrade replay.
5. Snapshot the node at **H−1000** (rollback safety).

```bash
# Cosmovisor layout example
mkdir -p $HOME/.canine/cosmovisor/upgrades/v600/bin
cp ./build/canined $HOME/.canine/cosmovisor/upgrades/v600/bin/
```

---

## 4. Upgrade window

| Step | Action |
|------|--------|
| H−24h | Announce halt height, binary URL, checksum |
| H−1h | Validators confirm Cosmovisor / binary in place |
| H | Chain halts; nodes restart with `v600` binary |
| H+1 | Confirm blocks resume; `canined status` shows new version |
| H+10 | Run smoke tests (below) |

### Post-upgrade verification

```bash
canined query upgrade applied v600
canined query consensus params
canined query storage params
```

---

## 5. Smoke test matrix

Automated script (queries only):

```bash
CHAIN_ID=jackal-testnet-1 NODE=https://testnet-rpc.jackalprotocol.com:443 KEY=test \
  ./scripts/smoke-v600-testnet.sh
```

### Manual functional tests

| # | Test | Pass criteria |
|---|------|---------------|
| 1 | **Storage post-proof zero fee** | Tx with *only* `MsgPostProof` (or attest/report/request-attestation) pays **0 ujkl** fee |
| 2 | **Storage post-proof with fee** | Tx mixing free msg + `MsgSend` charges normal fee |
| 3 | **Filetree** | `postkey` → `post-file` → `add-viewers` / `remove-viewers` |
| 4 | **RNS** | `register` / `resolve` query |
| 5 | **Oracle** | Feed query returns data |
| 6 | **Wasm** | Instantiate + execute a pinned testnet contract (record code ID) |
| 7 | **IBC transfer** | Transfer to counterparty testnet channel and back |
| 8 | **State export** | `canined export` at H+1000 succeeds |
| 9 | **Provider proof path** | Provider submits proof; no panic in storage keeper |

### Zero-fee post-proof example

```bash
# Build a tx with ONLY MsgPostProof — fee should be 0
canined tx storage post-proof ... --fees 0ujkl --gas auto --gas-adjustment 1.3 -y
canined query tx <hash> | jq '.tx.auth_info.fee'
```

---

## 6. Bake period

- Run testnet on `v600` for **≥ 2 weeks**.
- Monitor: missed blocks, wasm execution errors, proof throughput, IBC timeouts.
- File issues on the migration branch; do **not** proceed to mainnet until exit criteria in [COSMOS-MODERNIZATION.md](./COSMOS-MODERNIZATION.md) are met.

---

## 7. Rollback plan

If consensus fails at upgrade:

1. Stop all validators.
2. Restore pre-upgrade snapshot (height H−1000).
3. Restart with **pre-v600** binary.
4. Reschedule governance with fixed binary after root-cause.

---

## References

- Upgrade handler: `app/upgrades/v600/upgrades.go`
- Ante fee waiver: `app/ante_fee.go`, `app/ante_fee_test.go`
- [PHASE0-INVENTORY.md](./PHASE0-INVENTORY.md)
- [V600-MAINNET-GOVERNANCE.md](./V600-MAINNET-GOVERNANCE.md)
