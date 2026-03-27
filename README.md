# go-ethereum-butler

A hybrid CLI+TUI application for managing multi-chain EVM transactions.

## Features

- **CLI mode** — scriptable blockchain queries with `--json` output (AI agent friendly)
- **TUI mode** — interactive keyboard-driven interface
- Multi-chain support (currently: Chiliz Chain)
- Address info: balance, nonce, transaction history, token holdings
- Transaction and block lookups
- Send native currency (CHZ) and ERC-20 token transactions
- Secure private key handling (environment variables)
- Config-driven: add chains/tokens/contacts via JSON files

## Quick Start

### 1. Setup Configuration

```bash
cp .env.example .env
```

Edit `.env` and add your private keys (without `0x` prefix):

```ini
BUTLER_WALLET_MAIN=your_private_key_here
BUTLER_WALLET_TEST=your_test_private_key_here
```

### 2. Configure Chains

Create `chains.json`:
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

Optionally configure `tokens.json` and `contacts.json` (see examples below).

### 3. Build and Run

```bash
go build -o butler ./cmd/butler

# TUI mode (interactive)
./butler

# CLI mode (scriptable)
./butler address 0xC3B2A6D869868916b1f5D46f9b7C62eD2f1D2c49
./butler tx 0x9f978b07484bb439e790079afe192e0a562a93a26c9d893ea8001ddda88e9981
./butler block latest
./butler chain-info
```

## CLI Usage

```
butler                              Interactive TUI mode
butler address <addr>               Address info (balance, nonce, tx history, tokens)
butler tx <hash>                    Transaction details with receipt
butler block [number|latest]        Block information
butler chain-info                   Chain status (latest block, gas price)

Global flags:
  --chain <name>     Blockchain network name (default: first in chains.json)
  --json             Output in JSON format
  --config <path>    Path to config directory
  --help             Help for any command
```

### Examples

```bash
# Human-readable output
$ butler address 0xC3B2A6D869868916b1f5D46f9b7C62eD2f1D2c49

  Address:  0xC3B2A6D869868916b1f5D46f9b7C62eD2f1D2c49
  Chain:    Chiliz Chain (88888)
  Balance:  5045.818772 CHZ
  Nonce:    310
  Type:     EOA

  Recent Transactions (last 10):
    Hash            Method       Value              Time
    -----------------------------------------------------------------
    0x9f97...9981   Transfer     +3631.087415       4d ago
    0x12a8...b3e5   Transfer     +316.752300        7d ago
    ...

# JSON output (for scripts and AI agents)
$ butler address 0xC3B2...D2c49 --json
{
  "address": "0xC3B2A6D869868916b1f5D46f9b7C62eD2f1D2c49",
  "chain": "Chiliz Chain",
  "chain_id": 88888,
  "native_balance": "5045.818772",
  "native_symbol": "CHZ",
  "nonce": 310,
  "is_contract": false,
  "recent_txs": [...]
}

# Pipe-friendly
$ butler chain-info --json | jq .latest_block
32617854
```

## TUI Usage

Run `butler` without arguments for interactive mode.

### Navigation

- **Up/Down arrows** or **j/k** — Navigate menu items
- **Enter** — Select/Confirm
- **Esc** — Go back to main menu
- **Ctrl+C** — Quit

### Main Menu

1. **Send Transaction** — Select wallet > chain > token > recipient > amount > confirm
2. **Check Balance** — Select wallet > chain > token > view balance
3. **Exit** — Quit the application

## Architecture

```
go-ethereum-butler/
├── cmd/butler/
│   ├── main.go              # Entry point (Cobra Execute)
│   └── cmd/                 # CLI commands
│       ├── root.go          # Root command + global flags
│       ├── tui.go           # TUI launcher (no-args fallback)
│       ├── address.go       # butler address
│       ├── tx.go            # butler tx
│       ├── block.go         # butler block
│       └── chaininfo.go     # butler chain-info
├── internal/
│   ├── domain/              # Pure data models + output types
│   ├── infra/
│   │   ├── config/          # Config loading + path resolution
│   │   ├── ethereum/        # RPC client (go-ethereum)
│   │   │   ├── abi/         # Raw ABI JSON files
│   │   │   └── contracts/   # Generated Go bindings (abigen)
│   │   └── explorer/        # Block explorer API client (Chiliscan)
│   ├── output/              # Human/JSON output formatter
│   └── tui/                 # Bubble Tea interactive UI
│       ├── app.go           # Router
│       ├── style/           # Shared styles
│       └── pages/           # Page components
├── chains.json              # Chain configs (gitignored)
├── tokens.json              # Token configs (gitignored)
├── contacts.json            # Address book (gitignored)
└── .env                     # Private keys (gitignored)
```

### Design

- **Clean Architecture**: domain (models) > infra (RPC, config, explorer) > presentation (CLI, TUI)
- **Dual data sources**: RPC for real-time state, Explorer API for indexed data (tx history)
- **Hybrid CLI+TUI**: Cobra routes commands; no args = Bubbletea TUI
- **Config-driven**: chains/tokens/contacts as JSON, zero code changes to add new ones

See `CLAUDE.md` for detailed developer documentation.

## Configuration

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

The `explorer_api_url` enables transaction history and token discovery. Any Etherscan-compatible API works.

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

### contacts.json

```json
[
  {
    "name": "My Friend",
    "address": "0x..."
  }
]
```

### Config Directory

Butler looks for config files in this order:
1. `--config /path/to/dir` flag
2. `BUTLER_CONFIG_DIR` environment variable
3. `~/.butler/` (if `chains.json` exists there)
4. Current working directory

## Security

- Private keys are **never** stored in code or config files
- Keys are only loaded from environment variables at signing time
- `.env` is excluded from git via `.gitignore`
- CLI read-only commands (`address`, `tx`, `block`, `chain-info`) never access private keys

## Dependencies

- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [Bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — TUI styling
- [go-ethereum](https://github.com/ethereum/go-ethereum) — Ethereum client
- [godotenv](https://github.com/joho/godotenv) — Environment variable management

## Requirements

- Go 1.25.1 or higher
- Access to EVM-compatible RPC endpoints

## License

MIT
