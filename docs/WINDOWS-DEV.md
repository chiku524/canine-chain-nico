# Windows development setup

Cosmos chains with wasmvm need **CGO**, **gcc**, and **make**. Git Bash on Windows does not include these by default.

## Recommended: WSL2 Ubuntu

You already have Ubuntu WSL installed. Use it for all build and test commands.

### 1. Open WSL and go to the repo

```powershell
wsl -d Ubuntu
cd /mnt/c/Users/chiku/Desktop/Jackal/canine-chain-nico
```

### 2. Bootstrap (once)

```bash
chmod +x scripts/bootstrap-wsl-dev.sh
./scripts/bootstrap-wsl-dev.sh
source ~/.bashrc
```

### 3. Daily commands

```bash
./scripts/verify-v600-candidate.sh
make inventory-mainnet
CGO_ENABLED=1 WASMVM_TAG=v1.5.9 make build-linux
sha256sum build/canined-linux-amd64
```

---

## Alternative: Docker (no WSL gcc needed)

If Docker Desktop is running, build a **Linux** binary from Git Bash or PowerShell:

```bash
cd ~/Desktop/Jackal/canine-chain-nico
bash build-linux.sh
# Output: build/canined-linux-amd64
```

This uses `Dockerfile.linux-build` with CGO + wasmvm inside the container.

---

## Git Bash limitations

| Tool | Git Bash | WSL | Docker |
|------|----------|-----|--------|
| `make` | Not installed | Yes (after bootstrap) | N/A |
| `gcc` / CGO | Not installed | Yes | Inside container |
| `make test-unit` | Fails (no CGO) | Works | Use CI or WSL |
| `make build-linux` | Needs `make` | Works | Use `bash build-linux.sh` |

### Optional: Chocolatey (admin PowerShell)

If you prefer native Windows tools, run **PowerShell as Administrator**:

```powershell
choco install make mingw -y
```

Then restart Git Bash and set:

```bash
export CGO_ENABLED=1
export PATH="/c/ProgramData/mingw64/mingw64/bin:$PATH"
```

Wasmvm on native Windows uses `libwasmvm.dll` — WSL is still simpler for parity with validators.

---

## Inventory capture from WSL

Jackal public APIs may timeout from some networks. From WSL:

```bash
make inventory-mainnet
# or mirror:
REST_API=https://jackal-api.polkachu.com ./scripts/capture-chain-inventory.sh
```
