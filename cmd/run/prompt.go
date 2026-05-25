package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gqgs/llminvestbench/pkg/holding"
	"github.com/gqgs/llminvestbench/pkg/stocks"
)

type promptStock struct {
	Ticket    string  `json:"ticket"`
	Price     string  `json:"price"`
	MarketCap float64 `json:"market_cap"`
	Sector    string  `json:"sector,omitempty"`
}

const maxPromptBytes = 7000

func buildPrompt(holdings holding.Holdings, contexts []string, universe stocks.Stocks) string {
	payload := map[string]any{
		"holdings": holdings,
		"context":  contexts,
		"market":   compactMarket(universe),
	}
	encoded, _ := json.Marshal(payload)

	return strings.TrimSpace(`You're managing a portfolio of publicly tradable assets.

Objective: maximize the total value of the portfolio.

You will receive current holdings, recent context, and a bounded market table. You may trade only USD or tickers present in the market table. Use integer share quantities only.

Return strictly valid JSON with this exact shape:
{
  "updates": [
    {"ticket": "USD", "quantity": 1000, "price": 1, "action": "SELL", "reason": "Free cash for the target allocation."},
    {"ticket": "AAPL", "quantity": 5, "price": 200.12, "action": "BUY", "reason": "Concise investment rationale for this specific trade."}
  ],
  "context": ["short note useful for the next run"]
}

Rules:
- Use action "BUY" or "SELL".
- Quantity must be a positive integer.
- Price must match the supplied market price, except USD which is always 1.
- Every update must include a concise reason explaining that specific operation. Do not include hidden reasoning or step-by-step analysis.
- Do not sell more shares than currently held.
- Every BUY must be funded by preceding SELL updates in this order, including selling USD cash to invest it.
- If no trade is justified, return {"updates":[],"context":["reason for holding"]}.

Input:
`) + "\n" + string(encoded)
}

func limitPromptUniverse(holdings holding.Holdings, contexts []string, universe stocks.Stocks) stocks.Stocks {
	for len(universe) > 0 && len(buildPrompt(holdings, contexts, universe)) > maxPromptBytes {
		universe = universe[:len(universe)-1]
	}
	return universe
}

func selectUniverse(rows stocks.Stocks, holdings holding.Holdings, limit int) stocks.Stocks {
	held := map[string]struct{}{}
	for _, h := range holdings {
		held[h.Ticket] = struct{}{}
	}

	filtered := make(stocks.Stocks, 0, len(rows))
	for _, row := range rows {
		if row.Symbol == "" {
			continue
		}
		if _, err := stocks.CleanPrice(row.Lastsale); err != nil {
			continue
		}
		filtered = append(filtered, row)
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		return stocks.CleanNumber(filtered[i].MarketCap) > stocks.CleanNumber(filtered[j].MarketCap)
	})

	selected := make(stocks.Stocks, 0, limit+len(holdings))
	seen := map[string]struct{}{}
	for _, row := range filtered {
		if _, ok := held[row.Symbol]; ok {
			selected = append(selected, row)
			seen[row.Symbol] = struct{}{}
		}
	}
	for _, row := range filtered {
		if len(selected) >= limit {
			break
		}
		if _, ok := seen[row.Symbol]; ok {
			continue
		}
		selected = append(selected, row)
		seen[row.Symbol] = struct{}{}
	}
	return selected
}

func compactMarket(rows stocks.Stocks) []promptStock {
	market := make([]promptStock, 0, len(rows))
	for _, row := range rows {
		market = append(market, promptStock{
			Ticket:    row.Symbol,
			Price:     row.Lastsale,
			MarketCap: stocks.CleanNumber(row.MarketCap),
			Sector:    row.Sector,
		})
	}
	return market
}

func marketTickers(rows stocks.Stocks) map[string]struct{} {
	tickers := map[string]struct{}{"USD": {}}
	for _, row := range rows {
		tickers[row.Symbol] = struct{}{}
	}
	return tickers
}

func formatMoney(value float64) string {
	return fmt.Sprintf("%.2f", value)
}
