# Jackal public network endpoints

Use these when running inventory capture, smoke tests, or state export against live networks.  
Override with env vars if a provider is down.

| Network | Chain ID | RPC | REST / API |
|---------|----------|-----|------------|
| Mainnet | `jackal-1` | `https://rpc.jackalprotocol.com:443` | `https://api.jackalprotocol.com` |
| Testnet | verify with `canined status` (often `jackal-testnet-1`) | `https://testnet-rpc.jackalprotocol.com:443` | `https://testnet-api.jackalprotocol.com` |

Community mirrors (mainnet):

| Provider | RPC | REST |
|----------|-----|------|
| Polkachu | `https://jackal-rpc.polkachu.com` | `https://jackal-api.polkachu.com` |
| Brochain | `https://jackal-rpc.brocha.in:443` | `https://jackal-rest.brocha.in` |

## Quick checks

```bash
# Mainnet height
curl -s "https://api.jackalprotocol.com/cosmos/base/tendermint/v1beta1/blocks/latest" | jq -r '.block.header.height'

# Applied upgrade plan
curl -s "https://api.jackalprotocol.com/cosmos/upgrade/v1beta1/applied_plan" | jq .

# Wasm code count
curl -s "https://api.jackalprotocol.com/cosmwasm/wasm/v1/code?pagination.limit=1" | jq '.pagination.total // .code_infos | length'
```

## Scripts in this repo

| Script | Purpose |
|--------|---------|
| `scripts/capture-chain-inventory.sh` | Snapshot wasm codes + IBC channels to `docs/inventory/` |
| `scripts/smoke-v600-testnet.sh` | Post-upgrade query smoke (testnet or mainnet) |
| `scripts/verify-v600-candidate.sh` | Pre-release build + unit smoke on `master` |

See [V600-TESTNET-UPGRADE.md](./V600-TESTNET-UPGRADE.md) and [PHASE0-INVENTORY.md](./PHASE0-INVENTORY.md).
