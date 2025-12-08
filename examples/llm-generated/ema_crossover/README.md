# EMA Crossover — LLM-Generated Example

This directory contains artifacts produced autonomously by the Algosmithy agentic loop while generating a classic EMA-based crossover strategy.

All files were generated without manual edits.

## Files
- **strategy.go** — Go implementation of the strategy following the `IStrategy` interface.  
  Includes:
  - internal EMA implementation,
  - MACD-style signal logic,
  - multi-currency order execution through the strategy manager.

- **config.yaml** — backtest definition with a parameter grid over:
  - `pair` (ETH, BTC),
  - `timeframe` (5m, 1h),
  - EMA windows (`short_window`, `long_window`, `signal_window`),
  - system parameters (initial balances, subscriptions, executor).

- **result.json** — raw output from all backtest runs, one per parameter combination.

## What this example demonstrates

- autonomous generation of a multi-indicator strategy (short EMA, long EMA, signal EMA);
- correct lifecycle structure (`Init`, `OnStart`, `OnCandle`, `OnStopped`);
- correct use of `decimal.Decimal` for financial precision;
- proper wiring of subscriptions and order management via `IStrategyManager`;
- automatic generation of a parameterized grid-search backtest;
- successful execution inside the Docker-based deterministic engine.

This example is provided for illustration only — it demonstrates the autonomous generation pipeline rather than trading performance.
