# v630 testnet / devnet upgrade playbook

Deploy and validate the **Cosmos SDK 0.54** migration (`v630`) on Jackal devnet, private `jackal-nico-1`, or public testnet **after** prior hops (`v600`–`v620`) unless starting fresh on the phase4 binary.

| Field | Value |
|-------|-------|
| Upgrade name | `v630` |
| Branch | `feat/cosmos-modernization-phase4` |
| Handler | `app/upgrades/v630/upgrades.go` |
| Binary requirements | Go **1.25.9+**, CGO, wasmvm **v3.x** |
| Stores deleted at upgrade | `crisis`, `circuit` |

---

## 1. Pre-flight (local)

```bash
git checkout feat/cosmos-modernization-phase4
./scripts/verify-v630-candidate.sh
```

Skip slow sim in verify script: `SKIP_SIM=1 ./scripts/verify-v630-candidate.sh`

---

## 2. Build release candidate

```bash
WASMVM_VERSION=$(go list -m -f '{{.Version}}' github.com/CosmWasm/wasmvm/v3)
sudo wget -q "https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm.x86_64.so" \
  -O /usr/lib/libwasmvm.x86_64.so

CGO_ENABLED=1 make build-linux
sha256sum build/canined-linux-amd64
ldd build/canined-linux-amd64 | grep wasmvm || true
```

Publish binary + SHA256 in the devnet upgrade announcement.

---

## 3. Schedule upgrade

### Option A — governance proposal

```bash
canined tx gov submit-proposal software-upgrade v630 \
  --title "Jackal v630: Cosmos SDK 0.54 migration" \
  --description "Migration to SDK 0.54 / CometBFT 0.39 / ibc-go v11 / wasmd 0.70 / store/v2. See docs/V630-TESTNET-UPGRADE.md." \
  --upgrade-height <HALT_HEIGHT> \
  --deposit 1000000000000ujkl \
  --from <authority> \
  --chain-id <CHAIN_ID> \
  --gas auto
```

Private net helper:

```bash
UPGRADE_NAME=v630 UPGRADE_HEIGHT=200 ./scripts/submit-upgrade-proposal.sh
```

### Option B — Cosmovisor (recommended)

See [COSMOVISOR-LADDER.md](./COSMOVISOR-LADDER.md). Place binary at:

```bash
mkdir -p $HOME/.canine/cosmovisor/upgrades/v630/bin
cp build/canined-linux-amd64 $HOME/.canine/cosmovisor/upgrades/v630/bin/canined
```

---

## 4. Validator prep

1. Install **wasmvm v3** shared library (see [README](../README.md)).
2. Review [COMETBFT-039-VALIDATOR.md](./COMETBFT-039-VALIDATOR.md) — CometBFT **0.39** config changes.
3. Install Cosmovisor or plan manual binary swap at halt height − 1.
4. Ensure **≥ 2× disk** free for upgrade replay.
5. Snapshot at **H−1000**.

---

## 5. Smoke test matrix

Automated (queries):

```bash
CHAIN_ID=jackal-nico-1 NODE=http://127.0.0.1:26657 KEY=validator \
  ./scripts/smoke-v630-testnet.sh
```

### Manual functional tests

| # | Test | Pass criteria |
|---|------|---------------|
| 1 | **Storage post-proof zero fee** | Tx with only `MsgPostProof` pays **0 ujkl** |
| 2 | **Storage post-proof + send** | Mixed tx charges normal fee on non-free msgs |
| 3 | **Filetree** | `postkey` → `post-file` → `add-viewers` |
| 4 | **RNS** | `register` / resolve query |
| 5 | **Oracle** | Feed query returns data |
| 6 | **Wasm** | Instantiate + execute pinned contracts on **wasmvm v3** (code IDs in inventory) |
| 7 | **IBC transfer** | Transfer on ibc-go **v11** channel |
| 8 | **State export** | `canined export` at H+1000 succeeds |
| 9 | **Provider proof** | Provider submits proof; no panic |
| 10 | **Store/v2** | No IAVL panic; queries return consistent results post-upgrade |

### Wasm regression (mainnet code IDs)

From `docs/inventory/captured-mainnet-20260610T2333Z.json` — test code IDs **1–5** on wasmvm v3 before devnet sign-off.

---

## 6. Bake period

- Run devnet/testnet on `v630` for **≥ 4 weeks** before mainnet governance discussion.
- Monitor: missed blocks, wasm errors, proof throughput, IBC timeouts, CometBFT 0.39 peer connectivity.

---

## 7. Rollback plan

1. Stop validators.
2. Restore pre-upgrade snapshot (height H−1000).
3. Restart with **pre-v630** binary.
4. Root-cause and reschedule.

---

## References

- [V610-V620-UPGRADE-NOTES.md](./V610-V620-UPGRADE-NOTES.md) — earlier hops
- [COSMOVISOR-LADDER.md](./COSMOVISOR-LADDER.md)
- [COMETBFT-039-VALIDATOR.md](./COMETBFT-039-VALIDATOR.md)
- [JACKAL-DEVNET-HANDOFF.md](./JACKAL-DEVNET-HANDOFF.md)
- Handler: `app/upgrades/v630/upgrades.go`
