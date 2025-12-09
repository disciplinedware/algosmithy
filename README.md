# Algosmithy

**Algosmithy is a deterministic Go-based engine that drives a fully autonomous agentic loop** — where an LLM generates a strategy, executes it in isolated sandboxes, evaluates the results, and iterates until performance targets are met.

**What Algosmithy is for**  
Algosmithy is designed for quant teams that need to **accelerate and automate the strategy-development workflow** — replacing manual coding and fragmented backtesting tools with a unified, deterministic system that generates, runs, evaluates, and iterates strategies at high throughput.

If your bottlenecks are slow R&D, inconsistent results, manual experimentation, or scattered tooling, Algosmithy turns strategy research into a **continuous, reproducible, high-velocity experimentation pipeline**.



## Architecture & Strategy Execution

* **Custom actor engine** with two execution modes:
  * **Synchronous, deterministic** — for reproducible backtests.
  * **Asynchronous, multi-threaded** — for sandbox and live trading.
* **Standard actor signals** (`Started{}`, `Stopped{}`) with direct message dispatch in sync mode.
* **Time-ordered market-event playback** — currently implemented for candles; other feeds follow the same event pipeline.
* **Unified strategy interfaces:** `ICandleStrategy`, `IOrderBookStrategy`, registered via `RegisterStrategy(...)`.
* **TradingEngine** manages capital, positions, order execution; strategies stay isolated from infrastructure logic.
* **Order execution as sub-strategies:** market orders implemented; architecture prepared for trailing, limit, TWAP, etc.
* **Multi-currency positions & PnL:** FIFO accounting, separate PnL in base/quote currencies, and correct fee application.


## Backtesting & Optimization

* **Parallel backtest execution** with full parameter-space exploration.
* **Single YAML configuration** defining environment, tests, and strategy parameters (`system_params`, `params`).
* **Multi-strategy runs** — a single experiment can evaluate several strategies simultaneously.
* **Backtests run inside Docker**, mounting the generated strategy code, `config.yaml`, and the core framework.
* **Docker images are reused or auto-built** when missing.


## Data & Connectors

* **Pluggable exchange connectors** — current implementation: Bybit (others follow the same interface).
* **Historical data caching** — in-memory for parallel backtests, file-based for reuse across runs.
* **Metadata layer** describes trading pairs (source, quote currency, fees, etc.).


## LLM Integration

* **Structured prompts and examples** (Go code + YAML) ensure correct strategy generation.
* **Tool-based interface for the LLM** — `RunBacktest`, `ListAvailableFunctions`, `GetMemory`, `SetMemory`, `InternalThoughts`.
* **LLM interacts with the engine through structured tools** — requesting market data, inspecting available indicators, running backtests, and finalizing strategies.
* **Automatic injection of relevant files and examples** into model messages.
* **Generated code is inserted into the project automatically** — including strategy implementation, parameter definitions, and registration via `RegisterStrategy(...)`.


## Tooling & Financial Accuracy
* **Code generation** for automatic name-based strategy registration.
* **`RunBacktest` executes the strategy inside an isolated Docker sandbox** using the generated code and YAML configuration.
* **Strict financial precision** — `decimal.Decimal` everywhere, `float64` forbidden in tests, PnL validated down to the last cent.


## Documentation & Examples
- [Architecture overview](https://github.com/disciplinedware/algosmithy/blob/main/docs/architecture.md)
- [Agentic loop details](https://github.com/disciplinedware/algosmithy/blob/main/docs/agentic-loop.md)
- [LLM-generated strategy examples](https://github.com/disciplinedware/algosmithy/blob/main/examples/llm-generated/README.md)


## Contact
For collaboration or enterprise integration: **konstantin@disciplinedware.com**  
Maintainer: **Konstantin Trunin** — https://www.linkedin.com/in/ktrunin
