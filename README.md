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

## Notes

- Removed __Gemini__ for now because the available free chat UI can't search for updated prices nor does it support the upload of CSV or JSON :grimacing:.
- Removed __Claude__ for now because the available free chat UI can't search for updated prices and its context window is too small for uploaded files :grimacing:.
- Removed __ChatGPT__ for now becaue the available free chat UI can't no longer do complex data analysis :grimacing:.

## Project Structure

- `cmd`: Contains the main command implementations
  - `create`: Initialize new portfolios
  - `list`: Display current holdings and context
  - `update`: Process investment orders and update holdings
  - `stocks`: Fetch most recent stock prices

## Prompt

The most recent prompt with the clear guidelines can be see [here](./cmd/create/prompt.txt) and [here](./cmd/list/prompt.txt).

## Current Portfolio (2025-03-24)

| Model | Ticket | Sum | Quantity |
|-------|-------|-------|--------|
|`deepseek`|`AVGO`|10754|55|
|`deepseek`|`AMD`|9930|99|
|`deepseek`|`AAPL`|29883|125|
|`deepseek`|`AMZN`|29887|150|
|`deepseek`|`NVDA`|18846|159|
|`grok`|`AMZN`|25968|131|
|`grok`|`NVDA`|69108|583|
|`perplexity`|`USD`|10|10|
|`perplexity`|`COST`|9849|11|
|`perplexity`|`AMD`|1607|15|
|`perplexity`|`AXON`|9492|17|
|`perplexity`|`CRWD`|9963|27|
|`perplexity`|`AMGN`|9766|31|
|`perplexity`|`DUOL`|9797|32|
|`perplexity`|`AAPL`|9848|46|
|`perplexity`|`AMZN`|9942|51|
|`perplexity`|`CTAS`|9917|51|
|`perplexity`|`COIN`|9899|52|
|`perplexity`|`DDOG`|9904|96|



| Model | Total Sum | Change |
|-------|-----------|--------|
|`perplexity`|99994|—|
|`deepseek`|99300|-0.70%|
|`grok`|95076|-4.92%|

