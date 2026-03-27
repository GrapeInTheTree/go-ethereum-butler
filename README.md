# go-ethereum-butler

A hybrid CLI+TUI application for managing multi-chain EVM transactions. Built with Go, optimized for Chiliz Chain.

## Install

### Homebrew (macOS / Linux)

```bash
brew tap GrapeInTheTree/tap
brew install butler
```

### Go Install

```bash
go install github.com/GrapeInTheTree/go-ethereum-butler/cmd/butler@latest
```

### Download Binary

Pre-built binaries for macOS (Intel/Apple Silicon) and Linux (amd64/arm64) are available on the [Releases](https://github.com/GrapeInTheTree/go-ethereum-butler/releases) page.

### Build From Source

```bash
git clone https://github.com/GrapeInTheTree/go-ethereum-butler.git
cd go-ethereum-butler
go build -o butler ./cmd/butler
```

## Features

- **CLI mode** — scriptable blockchain queries with `--json` output for AI agents and scripts
- **TUI mode** — interactive keyboard-driven interface for balance checks and token transfers
- **Dual data sources** — RPC for real-time chain state + Chiliscan API for transaction history
- **Multi-chain ready** — add any EVM chain via JSON config (currently: Chiliz Chain)
- **Config-driven** — chains, tokens, and contacts as JSON files, zero code changes to extend
- **Secure** — private keys only in `.env`, loaded on-demand at signing time, never cached

## CLI Usage

```
butler                              Interactive TUI mode (no subcommand)
butler address <addr>               Address info: balance, nonce, tx history, token holdings
butler tx <hash>                    Transaction details with receipt
butler block [number|latest]        Block information
butler chain-info                   Chain status: latest block, gas price

Global flags:
  --chain <name>     Blockchain network (default: first in chains.json)
  --json             Machine-readable JSON output
  --config <path>    Config directory path
  -h, --help         Help for any command
```

### Examples

```bash
# Address overview with balance, nonce, and recent transactions
$ butler address 0xC3B2A6D869868916b1f5D46f9b7C62eD2f1D2c49

  Address:  0xC3B2A6D869868916b1f5D46f9b7C62eD2f1D2c49
  Chain:    Chiliz Chain (88888)
  Balance:  5045.818772 CHZ
  Nonce:    310
  Type:     EOA

  Token Holdings:
    PEPPER     12340.000000

  Recent Transactions (last 10):
    Hash            Method       Value              Time
    -----------------------------------------------------------------
    0x9f97...9981   Transfer     +3631.087415       4d ago
    0x12a8...b3e5   Transfer     +316.752300        7d ago
    ...

# JSON output for scripts and AI agents
$ butler address 0xC3B2A6D869868916b1f5D46f9b7C62eD2f1D2c49 --json
{
  "address": "0xC3B2A6D869868916b1f5D46f9b7C62eD2f1D2c49",
  "chain": "Chiliz Chain",
  "chain_id": 88888,
  "native_balance": "5045.818772",
  "native_symbol": "CHZ",
  "nonce": 310,
  "is_contract": false,
  "token_balances": [...],
  "recent_txs": [...]
}

# Transaction details
$ butler tx 0x9f978b07484bb439e790079afe192e0a562a93a26c9d893ea8001ddda88e9981

  Tx Hash:    0x9f978b...88e9981
  Status:     success
  Block:      32492204
  Time:       2026-03-23 00:23:15 UTC
  From:       0xa3DF8880d1D1BfC5Bea208AC3f1662420B2E2657
  To:         0xC3B2A6D869868916b1f5D46f9b7C62eD2f1D2c49
  Value:      3631.087415 CHZ
  Gas Price:  5001.00 Gwei
  Gas Used:   21000 / 21000
  Tx Fee:     0.105021 CHZ
  Nonce:      5343
  Logs:       0 events

# Block info
$ butler block latest

  Block:       #32617854
  Hash:        0x3f1317eb...75aeb5a
  Time:        2026-03-27 09:05:54 UTC
  Miner:       0xc10ae5Cd2C63E4065f81E241c93237E06E12d41b
  Gas Used:    70012 / 30000000
  Base Fee:    2500.00 Gwei
  Txs:         1 transactions

# Chain status
$ butler chain-info --json | jq .latest_block
32617854
```

## TUI Usage

Run `butler` with no arguments for interactive mode.

### Navigation

| Key | Action |
|-----|--------|
| `up` / `k` | Move cursor up |
| `down` / `j` | Move cursor down |
| `enter` | Select / Confirm |
| `esc` | Back to main menu |
| `ctrl+c` | Quit |
| `0-9`, `.` | Amount input (send flow) |
| `backspace` | Delete last character |

### Flows

- **Check Balance** — Select wallet > chain > token (native or ERC-20) > view balance
- **Send Transaction** — Select wallet > chain > token > recipient > enter amount > confirm > tx hash

## Configuration

Butler uses JSON config files. All are gitignored by default — create your own.

### chains.json

```json
[
  {
    "name": "Chiliz Chain",
    "rpc_url": "https://rpc.ankr.com/chiliz",
    "chain_id": 88888,
    "currency_symbol": "CHZ",
    "logo_url": "",
    "explorer_api_url": "https://api.routescan.io/v2/network/mainnet/evm/88888/etherscan/api"
  }
]
```

- `rpc_url` — any EVM-compatible JSON-RPC endpoint
- `explorer_api_url` — Etherscan-compatible API (enables tx history and token discovery). Optional; if omitted, CLI commands still show RPC-based data (balance, nonce, blocks)
- Add multiple chains to the array for multi-chain support; select with `--chain <name>`

### tokens.json

```json
[
  {
    "symbol": "PEPPER",
    "name": "Pepper Token",
    "address": "0x60F397acBCfB8f4e3234C659A3E10867e6fA6b67",
    "decimals": 18,
    "chain_id": 88888,
    "logo_url": ""
  }
]
```

- `decimals` must match the contract (18 for most, 6 for USDC/USDT, 8 for WBTC)
- `chain_id` must match a chain in `chains.json`
- Native tokens (CHZ, ETH) are automatically available — no entry needed
- Same token on different chains needs separate entries

### contacts.json

```json
[
  {
    "name": "My Friend",
    "address": "0x..."
  }
]
```

Used by the TUI send flow for recipient selection.

### .env

```ini
BUTLER_WALLET_MAIN=your_private_key_here_without_0x_prefix
BUTLER_WALLET_TEST=your_test_private_key_here
```

Copy from `.env.example`. Private keys are **never** logged, cached, or included in git. CLI read-only commands (`address`, `tx`, `block`, `chain-info`) do not access private keys.

### Config Directory Resolution

Butler searches for config files in this order:

1. `--config /path/to/dir` flag (explicit)
2. `BUTLER_CONFIG_DIR` environment variable
3. `~/.butler/` directory (if `chains.json` exists there)
4. Current working directory (default)

## Architecture

```
go-ethereum-butler/
├── cmd/butler/
│   ├── main.go                  # Entry point: Cobra Execute()
│   └── cmd/
│       ├── root.go              # Root command, global flags, PersistentPreRunE
│       ├── tui.go               # TUI launcher (no-args fallback)
│       ├── address.go           # butler address — parallel RPC + Explorer
│       ├── tx.go                # butler tx — tx + receipt lookup
│       ├── block.go             # butler block — block by number
│       └── chaininfo.go         # butler chain-info — chain status
│
├── internal/
│   ├── domain/
│   │   ├── models.go            # Chain, Token, Wallet, Contact structs
│   │   └── output.go            # AddressInfo, TxDetail, BlockInfo, ChainStatus
│   │
│   ├── infra/
│   │   ├── config/config.go     # JSON/env config loading + path resolution
│   │   ├── ethereum/
│   │   │   ├── client.go        # RPC: balance, nonce, code, blocks, tx, gas
│   │   │   ├── erc20.go         # ERC-20: balance, transfer, formatting
│   │   │   ├── abi/erc20.json   # Standard ERC-20 ABI
│   │   │   └── contracts/       # Auto-generated Go bindings (abigen)
│   │   └── explorer/
│   │       └── etherscan.go     # Chiliscan API: tx history, token discovery
│   │
│   ├── output/
│   │   └── formatter.go         # Human-readable + JSON dual formatter
│   │
│   └── tui/
│       ├── app.go               # Bubble Tea router
│       ├── style/style.go       # Lipgloss styles
│       └── pages/               # mainmenu, balance, send
│
├── .goreleaser.yml              # Cross-compile + Homebrew formula generation
├── .github/workflows/
│   └── release.yml              # Tag-triggered GoReleaser CI
│
├── chains.json                  # Chain configs (gitignored, user-created)
├── tokens.json                  # Token configs (gitignored, user-created)
├── contacts.json                # Address book (gitignored, user-created)
└── .env                         # Private keys (gitignored, user-created)
```

### Design Principles

- **Clean Architecture** — domain (pure models) > infra (RPC, config, explorer) > presentation (CLI, TUI)
- **Hybrid CLI+TUI** — Cobra routes subcommands; no args = Bubbletea TUI. Single binary.
- **Dual Data Sources** — RPC for real-time chain state (balance, blocks, tx by hash). Explorer API for indexed data (tx history by address, token discovery). Explorer is optional with graceful degradation.
- **Config-Driven** — chains, tokens, contacts as JSON files. Zero code changes to add new chains or tokens.
- **Concurrent RPC Calls** — `butler address` fires 5 parallel requests (3 RPC + 2 Explorer) for fast response.

### Data Source Strategy

| Data | Source | Why |
|------|--------|-----|
| Native balance | RPC `eth_getBalance` | Standard RPC method |
| ERC-20 balance | RPC `eth_call` (balanceOf) | Direct contract call |
| Nonce, code, gas price | RPC | Standard RPC methods |
| Tx by hash + receipt | RPC | Standard RPC methods |
| Block by number | RPC | Standard RPC method |
| **Tx history by address** | **Explorer API** | No RPC method exists for this |
| **All token holdings** | **Explorer API** | Token discovery requires indexer |

## Releasing

Releases are automated via GoReleaser + GitHub Actions.

```bash
# 1. Commit your changes
git add . && git commit -m "feat: ..."

# 2. Tag and push
git tag v0.3.0
git push && git push --tags

# 3. GitHub Actions automatically:
#    - Cross-compiles linux/darwin x amd64/arm64
#    - Creates GitHub Release with changelog and binaries
#    - Updates Homebrew formula in GrapeInTheTree/homebrew-tap
```

Users upgrade with `brew upgrade butler` or `go install ...@latest`.

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [spf13/cobra](https://github.com/spf13/cobra) | v1.10.2 | CLI framework |
| [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) | v1.3.10 | TUI framework |
| [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) | v1.1.0 | TUI styling |
| [ethereum/go-ethereum](https://github.com/ethereum/go-ethereum) | v1.16.7 | Ethereum RPC client, tx signing, ABI bindings |
| [joho/godotenv](https://github.com/joho/godotenv) | v1.5.1 | .env file loader |

## Requirements

- Go 1.25.1+ (build from source / go install)
- Access to EVM-compatible RPC endpoints
- Chiliscan API access for tx history (free, no API key)

## License

MIT
