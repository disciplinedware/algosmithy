# Architecture Overview

Algosmithy Core is a deterministic, actor-based runtime that wires a data dispatcher, strategy managers, and accounting into a reproducible backtest/sandbox loop.


## 1. Actor Engine

* Two execution modes in `ActorsEngine`:  
  * **SingleThread** — synchronous dispatch without queues for deterministic backtests.  
  * **MultiThread** — per-actor buffered inboxes with goroutines for sandbox/live.
* Standard lifecycle signals (`Started{}`, `Stopped{}`) are delivered to every actor.


## 2. Backtest Pipeline

* `BackTester` builds the actor system in **SingleThread** mode, registers a single `BacktestDispatcher` plus one `StrategyManager` per strategy class from YAML config.
* Strategies declare subscriptions (currently candles); the dispatcher loads candles via `HistoryDataLoader`, sorts all feeds into a single timestamped timeline, and streams them to subscribers while keeping the latest candle per pair.
* Market orders are executed at the latest candle price; fees are applied according to metadata (`quote`, `base`, or `side-dependent` modes) and emitted as trades back to the requesting strategy actor.
* `StrategyManager` wraps the user strategy, decodes params, keeps a multi-currency account (FIFO cash ledger), exposes `GetAvailableSizeToBuy/Sell` and `RegisterOrder`, and collects a final report.


## 3. Strategy Interfaces & Registration

* Base interface: `IStrategy` with `Init`, `OnStart`, `OnStopped`; extensions:  
  * `ICandleStrategy`  
  * `IOrderBookStrategy`  
  * `ITradeStrategy`
* Strategies are discovered via `strategy_factories.RegisterStrategy("<class>", factory)` and instantiated by `StrategyManager` based on the YAML config.
* Strategies interact with the platform through `ISmartContext` (logging, metadata, actors engine, data loader, LLM components).


## 4. Event Playback

* Candle data is merged across pairs/timeframes, globally sorted, and replayed in order.
* Progress is logged during long replays; cancellation of the parent context stops playback cleanly.


## 5. Sandbox & Reproducibility

* Backtests run inside an isolated Docker container (mounted code + config + core), making LLM-generated strategies safe and repeatable.
* All monetary calculations use `decimal.Decimal`; accounting handles multi-currency balances, PnL, and commissions precisely.


This architecture keeps strategies isolated while providing deterministic playback, realistic execution (fees, metadata-aware sizing), and a clean hook for the agentic LLM loop.
