# Agentic Loop

Algosmithy runs an autonomous research loop where an LLM generates Go strategies, backtests them inside an isolated sandbox, and iterates until a candidate is finalized.


## 1. Data Discovery (Tools)

The model works through structured tools defined in `llm_strat_gen/tools.go`:

* `GetCandles`, `GetOrderFlow`, `GetFundingRates`, `GetLiquidations`, `GetVolumeProfile`, `GetMarketLevels`, `GetTrendlines`, `GetPatterns`, `GetFibonacciLevels`, `ListAvailableIndicators`.

## 2. Strategy Generation

* System messages are loaded from config and interpolated with embedded examples (`example_go_code.go`, `example_yaml_config.yaml`, `example_json_result.json`) so the model sees the exact expected formats.
* The LLM is expected to emit:
  * Go code implementing `ICandleStrategy` or `IOrderBookStrategy`;
  * `RegisterStrategy(...)` call in `init()` with the strategy class name;
  * parameter structs and YAML config that describe environment, test grid, and params.

## 3. Execution Loop

* `LlmStratGen.Start` runs up to 20 iterations. Each iteration:
  * calls the model with the current transcript and the tool schema;
  * executes every requested tool and appends results/errors back to the chat;
  * logs chain-of-thought content without blocking generation.
* Available tools:
  * `RunBacktest` — compiles the provided Go code + YAML and calls `LlmStratContext.RunBacktest`, which runs `apps/back-tester-v1` inside Docker with the mounted code/config/core;
  * `FinalizeAlgorithm` — stores name/description through `LlmStratContext.FinalizeAlgorithm`;
  * `SetShortTermMemory` / `GetShortTermMemory` — in-memory scratchpad for intermediate results;
  * `InternalThoughts` — persists reasoning traces;
  * all data-discovery tools listed above.


## 4. Feedback & Finalization

* Every tool result (including errors) is returned to the model, allowing it to self-correct code or configs in the next turn.
* The loop exits when the model calls `FinalizeAlgorithm`; otherwise it stops after 20 iterations with a warning.
* Memory and finalization hooks are pluggable — current implementation is in-memory but can be swapped for vectors or persistent storage.

