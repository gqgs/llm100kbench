# LLM Investment Benchmark

A tool for benchmarking and tracking Large Language Model (LLM) investment decisions.

## Overview

This project provides a framework to create, manage, and track investment portfolios generated by LLM models. It allows you to:
- Create new portfolios
- List current holdings and recent context
- Update portfolios based on model decisions

The model executions and their current context can be seen [here](./orders).

## Why?

To optimize their portfolio, the primary objective defined for the LLMs, it is imperative to evaluate the risk-reward ratio, formulate cogent assumptions about future market conditions, and leverage tools and their understanding of human psychology and financial market dynamics.

This benchmark may be a good proxy to measure how well LLMs are able to coordinate the aforementioned efforts.

## Project Structure

- `cmd`: Contains the main command implementations
  - `create`: Initialize new portfolios
  - `list`: Display current holdings and context
  - `update`: Process investment orders and update holdings

## Prompt

The most recent prompt with the clear guidelines can be see [here](./cmd/create/prompt.txt) and [here](./cmd/list/prompt.txt).

## Current Portfolio (2025-03-01)

| Model | Ticket | Sum | Quantity |
|-------|-------|-------|--------|
|`chatgpt`|`USD`|69|69|
|`chatgpt`|`AAPL`|99931|418|
|`deepseek`|`AMD`|9930|99|
|`deepseek`|`AAPL`|29883|125|
|`deepseek`|`ALAB`|9971|149|
|`deepseek`|`AMZN`|29887|150|
|`deepseek`|`ARM`|19962|159|
|`grok`|`ADBE`|29660|66|
|`grok`|`AAPL`|49965|209|
|`grok`|`ACGL`|19929|219|


| Model | Total Sum | Change |
|-------|-----------|--------|
|`chatgpt`|100000|—|
|`deepseek`|99633|—|
|`grok`|99554|—|
