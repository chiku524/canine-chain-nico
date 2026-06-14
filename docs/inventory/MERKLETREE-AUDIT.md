# Storage merkletree fork audit

**Date:** 2026-06-08  
**Status:** Retained — devnet proof regression required  
**Related:** [PHASE0-INVENTORY.md](../PHASE0-INVENTORY.md), [COSMOS-MODERNIZATION.md](../COSMOS-MODERNIZATION.md)

---

## Replace directive

```go
github.com/wealdtech/go-merkletree/v2 => github.com/TheMarstonConnell/go-merkletree/v2 v2.0.0-20250829184252-ad65f46fbd22
```

Jackal storage proofs depend on this fork for protocol-specific tree semantics. It is **not** removable without rewriting proof verification in `x/storage`.

---

## Why it stays through v630

| Phase | store backend | Merkletree impact |
|-------|---------------|-------------------|
| v600–v620 | IAVL (cosmossdk.io/store) | Proof bytes unchanged at app level |
| v630 | **store/v2** | Keeper + proof paths must still produce/consume same merkle semantics |

Unit tests in `x/storage/keeper` cover proof math; **devnet** must confirm live `MsgPostProof` against real providers.

---

## Pre-devnet verification (completed in fork)

- [x] Replace still present in `go.mod` on `feat/cosmos-modernization-phase4`
- [x] `x/storage/keeper` unit tests pass (`make test-unit`)
- [x] `mulStorageCharge` overflow guard in place
- [ ] Devnet: provider submits proof against active deal
- [ ] Devnet: post-proof tx **zero fee** (ante waiver)

---

## Devnet regression script

1. Buy storage + post file on devnet.
2. Run provider proof submission.
3. Confirm proof accepted in `query storage` proofs-by-address.
4. Compare proof root/hash with pre-upgrade behavior if testing sequential upgrades.

---

## Risk

If the merkletree fork is incompatible with a future Go or SDK release, proofs break silently until devnet catches it. Pin version in `go.mod` until upstream Jackal merges an audited path.

---

## References

- Fork: https://github.com/TheMarstonConnell/go-merkletree/v2
- Storage keeper: `x/storage/keeper/proofs.go`
