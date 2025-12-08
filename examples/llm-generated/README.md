# LLM-Generated Examples

This directory contains **real artifacts produced by the Algosmithy autonomous agentic loop**.  
Each subfolder represents a complete generation cycle where the LLM:

1. requested market data & available indicators,
2. generated Go strategy code,
3. produced a full YAML backtest configuration (including param grids),
4. executed the strategy inside a sandboxed Docker environment,
5. analyzed the results and refined the code if necessary.

All files in these examples were generated **without manual editing**.

---

## Structure of Each Example

Every strategy folder contains:

- **strategy.go** — the generated Go implementation of the strategy  
- **config.yaml** — the backtest experiment (parameters, instruments, subscriptions, system settings)  
- **result.json** — raw output of all backtest runs across the generated parameter grid  
- (optional) **description.txt** — short notes returned by the LLM during refinement  

These examples are not meant to be profitable strategies.  
They are intentionally simple (EMA crossover, Bollinger breakout) in order to demonstrate:

- correct adherence to the Algosmithy strategy interface,
- proper handling of `decimal.Decimal` for financial precision,
- autonomous generation of multi-run backtest configurations,
- correct event-driven lifecycle (`Init`, `OnCandle`, etc.),
- deterministic execution inside the core engine.

---

## Included Examples

### **ema_crossover/**
A multi-indicator EMA/MACD-style crossover strategy.  
Demonstrates multi-EMA dependency graph, order-management integration, and grid-search generation.

### **bollinger_bands_breakout_strategy/**
A Bollinger breakout strategy with LLM-generated indicator implementation and multi-instrument experimentation.

---

## Purpose of This Directory

These examples serve as **technical validation of the autonomous flow**:

- strategy synthesis → compilation → execution → evaluation → refinement  
- fully deterministic, Go-based, Docker-sandboxed execution  
- correct wiring of market data, actor events, and order logic  
- structured backtest orchestration through YAML

They illustrate how Algosmithy can generate **complete, runnable quant research artifacts end-to-end**.

