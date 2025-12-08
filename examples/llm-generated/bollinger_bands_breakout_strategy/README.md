# Bollinger Bands Breakout — LLM-Generated Example

This directory contains artifacts produced by the LLM for a Bollinger Bands breakout strategy.
All files were generated autonomously without manual corrections.

### Files

- **strategy.go** — Go implementation of the strategy following the `IStrategy` interface  
- **config.yaml** — a complete backtest experiment definition:
  - grid search over `pair`, `timeframe`, `stddev`, `window_size`
  - system parameters (initial balances, subscriptions)
  - executor definition (`trailing_smart`)
- **result.json** — raw output of all backtest runs for every parameter combination

### What this example demonstrates

- correct use of `decimal.Decimal`
- adherence to the unified strategy API (`Init`, `OnStart`, `OnCandle`, etc.)
- automatic construction of parameterized grid-search backtests
- multi-instrument and multi-timeframe experimentation
- correct YAML templating (`$pair`, `$window_size`, `$timeframe`)
- deterministic execution through the core backtesting engine

This example is provided for demonstration purposes only.
It shows how the Algosmithy agent generates runnable Go strategies and their associated experiment definitions.
