# Changelog

All notable changes to this project will be documented in this file.

## [0.2.0] - 2026-03-27

### Added
- **Cobra CLI framework** — `butler` is now a hybrid CLI+TUI tool
  - `butler address <addr>` — comprehensive address info (balance, nonce, tx history, token holdings)
  - `butler tx <hash>` — full transaction details with receipt
  - `butler block [number|latest]` — block information
  - `butler chain-info` — chain status (latest block, gas price)
- **`--json` flag** on all commands for machine-readable/AI-agent-friendly output
- **`--chain` flag** for multi-chain support (default: first chain in chains.json)
- **`--config` flag** for custom config directory
- **Chiliscan API client** (`internal/infra/explorer/etherscan.go`)
  - Etherscan-compatible API integration via Routescan
  - Transaction history retrieval (not possible via standard RPC)
  - Token balance discovery for all ERC-20 holdings
  - Built-in rate limiting (2 req/sec for free tier)
- **8 new RPC query functions** in `internal/infra/ethereum/client.go`:
  - `GetNonce`, `GetCode`, `GetChainID`, `GetGasPrice`
  - `GetLatestBlockNumber`, `GetTransaction`, `GetTransactionReceipt`, `GetBlock`
- **Output formatter** (`internal/output/formatter.go`) — dual mode human-readable + JSON
- **Domain output types** (`internal/domain/output.go`) — stable JSON-serializable structs
- **Config path resolution cascade**: `--config` flag > `BUTLER_CONFIG_DIR` env > `~/.butler/` > CWD
- `ExplorerAPIURL` field in Chain model for per-chain block explorer API support

### Fixed
- **`pow10()` integer overflow** in `erc20.go` — changed from `int64` to `*big.Int`, preventing overflow for tokens with >18 decimals

### Changed
- `cmd/butler/main.go` refactored from direct TUI launch to Cobra `Execute()` (3 lines)
- Config loading functions now use `configPath()` for directory-aware file resolution
- `.env` loading attempts config directory first, then falls back to CWD

## [0.1.0] - 2024-11-19

### Added
- Initial TUI application with Bubble Tea framework
- Native currency (CHZ) balance checks via RPC
- Native currency transfers with EIP-1559 gas pricing
- ERC-20 token balance checks and transfers via abigen bindings
- Multi-wallet support (Main Wallet, Test Wallet)
- Address book (contacts.json)
- Config-driven chain/token/contact management via JSON files
- Chiliz Chain configuration with PEPPER token
- Structured logging to `butler.log` via slog
