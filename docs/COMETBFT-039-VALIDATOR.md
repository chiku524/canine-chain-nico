# CometBFT 0.39 validator guide (Jackal v630)

Jackal `v630` (SDK **0.54.3**) ships **CometBFT 0.39.3**. Validators upgrading from 0.37/0.38 must review node configuration before joining devnet or mainnet post-migration.

**Applies to:** `feat/cosmos-modernization-phase4` / on-chain upgrade `v630`

---

## Summary of changes vs 0.37 / 0.38

| Area | 0.37 / 0.38 | 0.39 |
|------|-------------|------|
| P2P stack | Classic CometBFT peer exchange | **libp2p** integration (config-dependent) |
| Sync | Block sync | **Adaptive sync** options |
| Config schema | `config.toml` v0 | Updated defaults — **do not blind-copy** old `config.toml` |
| State sync | Supported | Re-verify snapshot providers |

Always generate a fresh `config.toml` from the **v630 binary** and merge only intentional customizations (seeds, peers, ports, pruning).

---

## Recommended validator workflow

### 1. Fresh config from v630 binary

```bash
canined init <moniker> --chain-id <CHAIN_ID> --home <HOME>
# Compare new config/ with backed-up pre-upgrade config/
diff -u <backup>/config/config.toml <HOME>/config/config.toml
```

### 2. Preserve these from your old config

- `moniker`
- `external_address` / `listen_addr` (if non-default)
- `persistent_peers` / `seeds` / `unconditional_peer_ids`
- `p2p.max_num_inbound_peers` / `max_num_outbound_peers` (tune for your network)
- `pruning` strategy (if set via app.toml — see below)
- Sentry topology (keep validators private; sentries public)

### 3. Review new 0.39 sections

Read inline comments in the new `config.toml` for:

- **P2P / libp2p** — peer discovery settings
- **Block sync / adaptive sync** — start with defaults on devnet; tune after soak
- **Mempool** — may differ from 0.34-era Jackal mainnet

### 4. `app.toml` checks

- `minimum-gas-prices` — unchanged Jackal policy unless governance changed it
- `api` / `grpc` / `grpc-web` — re-enable if you serve public endpoints
- State sync snapshots — verify `snapshot-interval` if you provide snapshots to peers

---

## Devnet checklist

- [ ] Single validator produces blocks for ≥ 100 heights
- [ ] Second peer connects via `persistent_peers`
- [ ] `canined status` shows `SyncInfo.catching_up: false` after sync
- [ ] No repeated P2P disconnect loops in `canined.log`
- [ ] State sync from peer works (if used)
- [ ] Sentry → validator architecture still holds (validators not directly exposed)

---

## Common issues

| Symptom | Likely cause | Action |
|---------|--------------|--------|
| Peer connection failures | Stale peer IDs / wrong port | Update `persistent_peers` from devnet coordinator |
| Stuck catching up | Adaptive sync / snapshot mismatch | Try block sync; verify genesis hash |
| Mempool overflow | Low `size` / rate limits | Compare defaults with Jackal network policy |
| High disk use | New DB backend paths | Confirm `cosmos-db` + pruning settings |

---

## Rollback

If CometBFT 0.39 fails on devnet:

1. Restore snapshot taken **before** `v630` upgrade.
2. Run **pre-v630** binary (CometBFT 0.38.x).
3. File issue with logs + `config.toml` (redact secrets).

---

## References

- [CometBFT 0.39 release notes](https://github.com/cometbft/cometbft/releases)
- [V630-TESTNET-UPGRADE.md](./V630-TESTNET-UPGRADE.md)
- [JACKAL-DEVNET-HANDOFF.md](./JACKAL-DEVNET-HANDOFF.md)
