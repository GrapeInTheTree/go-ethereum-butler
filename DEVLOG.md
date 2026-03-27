# Development Log

## 2026-03-27 — Cobra CLI Integration (v0.2.0)

### Motivation

butler was TUI-only. AI agents and scripts had no way to query Chiliz Chain data programmatically. Standard EVM JSON-RPC has no method to retrieve transaction history by address — `eth_getTransactionsByAddress` simply doesn't exist. Block explorers (Chiliscan) solve this via indexing, but there was no CLI tool bridging that gap for Chiliz Chain.

Research confirmed CLI tools beat MCP for AI agent integration (10-32x cheaper in token usage, 100% reliability vs 72%). The decision was CLI-first, with Foundry's `cast` as the benchmark.

### Architecture Decisions

**Cobra + Bubbletea hybrid in a single binary.**
`butler` with no args launches the existing TUI. Subcommands (`butler address`, `butler tx`, etc.) run as CLI and exit. This required zero changes to the TUI package — Cobra's root command `RunE` simply calls the TUI launcher.

**Dual data source: RPC + Chiliscan API.**
RPC covers balance, nonce, blocks, individual tx lookups. Chiliscan (Routescan, Etherscan-compatible) covers transaction history and token discovery — data that fundamentally requires an indexer. The explorer client degrades gracefully: if a chain has no `explorer_api_url`, those sections are simply omitted.

**Package-level `appContext` instead of dependency injection.**
CLI tools run a single command and exit. There's no concurrency at the command level. A package-level struct holding resolved config, chain, and explorer client is simpler and more readable than context injection or interfaces. This can be revisited if butler ever becomes a long-running server.

**Output types in `domain/`, not in command packages.**
`AddressInfo`, `TxDetail`, `BlockInfo`, `ChainStatus` are stable JSON contracts. Placing them in `domain/` keeps them framework-agnostic and reusable across CLI, TUI, and future API layers.

### Key Technical Details

- **Chiliscan API**: `https://api.routescan.io/v2/network/mainnet/evm/88888/etherscan/api`. Free tier, no API key, 2 req/sec, 10k calls/day. Confirmed working for `txlist`, `addresstokenbalance`, `tokentx`, `balance`.
- **pow10 bug**: The original `pow10(n int) int64` would overflow for decimals > 18 (int64 max ~9.2e18). Changed to `*big.Int` via `Exp(10, n, nil)`. This was a latent bug that would have manifested with non-standard tokens.
- **Config resolution**: Added 4-level cascade (`--config` > env > `~/.butler/` > CWD) to support CLI usage from any directory while preserving backward compatibility for TUI users who run from the project root.
- **Parallel RPC calls**: `butler address` fires 5 concurrent requests (3 RPC + 2 explorer) via goroutines + sync.WaitGroup. Cuts latency from ~5x RTT to ~1x RTT.

### What's Not Included (Phase 2+)

- `butler send` CLI command (write operations deferred to Phase 2)
- Connection pooling for ethclient (unnecessary for single-shot CLI)
- Wallet management CLI (`butler wallet add/list`)
- ABI encode/decode utilities
- Chiliz-specific commands (validators, staking via system contracts)
- Unit tests for the new packages

### Files Changed

| Action | Count | Files |
|--------|-------|-------|
| New | 9 | `cmd/butler/cmd/{root,tui,address,tx,block,chaininfo}.go`, `internal/domain/output.go`, `internal/infra/explorer/etherscan.go`, `internal/output/formatter.go` |
| Modified | 6 | `cmd/butler/main.go`, `internal/infra/ethereum/{client,erc20}.go`, `internal/infra/config/config.go`, `internal/domain/models.go`, `chains.json` |
| Total | +1234 / -38 lines | |

---

## 2024-11-19 — Initial TUI (v0.1.0)

First working version. Bubble Tea TUI with balance checks and native/ERC-20 transfers on Chiliz Chain. Clean architecture with domain/infra/tui separation. Config-driven via JSON files.
