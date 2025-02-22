# LLM Investment Benchmark

A tool for benchmarking and tracking Large Language Model (LLM) investment decisions.

## Overview

This project provides a framework to create, manage, and track investment portfolios generated by LLM models. It allows you to:
- Create new portfolios
- List current holdings and recent context
- Update portfolios based on model decisions

The model executions and their current context can be seen [here](./orders).


## Project Structure

- `cmd`: Contains the main command implementations
  - `create`: Initialize new portfolios
  - `list`: Display current holdings and context
  - `update`: Process investment orders and update holdings

## Prompt

The most recent prompt with the clear guidelines can be see [here](./cmd/create/prompt.txt).

## Current Portfolio (2025-02-22)

| Model | Ticket | Sum | Quantity |
|-------|-------|-------|--------|
|`claude3.5`|`NVDA`|20000|25|
|`claude3.5`|`MSFT`|20000|50|
|`claude3.5`|`VOO`|60000|150|
|`deepseek-r1`|`NVDA`|100000|125|
|`gemini2.0-flash`|`TSLA`|100000|200|
|`grok3`|`BRK.B`|20000|50|
|`grok3`|`IWM`|15000|75|
|`grok3`|`METL`|10000|100|
|`grok3`|`BTCETF`|10000|200|
|`grok3`|`BSV`|24960|312|
|`grok3`|`INTC`|20000|400|
|`o3-mini`|`TSLA`|10134|30|
|`o3-mini`|`GOOGL`|9881|55|
|`o3-mini`|`MSFT`|29799|73|
|`o3-mini`|`AMZN`|19925|92|
|`o3-mini`|`AAPL`|29957|122|
|`o3-mini`|`USD`|303|303|
