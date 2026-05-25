package main

import (
	"fmt"
	"testing"

	"github.com/gqgs/llminvestbench/pkg/holding"
	"github.com/gqgs/llminvestbench/pkg/stocks"
	"github.com/stretchr/testify/require"
)

func TestSelectUniverseKeepsHeldTickers(t *testing.T) {
	rows := stocks.Stocks{
		{Symbol: "BIG", Lastsale: "$10.00", MarketCap: "1000"},
		{Symbol: "HELD", Lastsale: "$20.00", MarketCap: "1"},
		{Symbol: "MID", Lastsale: "$30.00", MarketCap: "500"},
	}
	holdings := holding.Holdings{{Ticket: "HELD", Quantity: 1}}

	selected := selectUniverse(rows, holdings, 2)

	require.Len(t, selected, 2)
	require.Equal(t, "HELD", selected[0].Symbol)
	require.Equal(t, "BIG", selected[1].Symbol)
}

func TestBuildPromptOmitsUnneededMarketFields(t *testing.T) {
	rows := stocks.Stocks{{
		Symbol:    "AAPL",
		Lastsale:  "$200.00",
		MarketCap: "3000000",
		Sector:    "Technology",
		Industry:  "Computer Manufacturing",
		Volume:    "12345",
		PctChange: "1.2%",
	}}

	prompt := buildPrompt(nil, nil, rows)

	require.Contains(t, prompt, `"sector":"Technology"`)
	require.NotContains(t, prompt, "industry")
	require.NotContains(t, prompt, "volume")
	require.NotContains(t, prompt, "pct_change")
}

func TestLimitPromptUniverseKeepsPromptUnderBudget(t *testing.T) {
	rows := make(stocks.Stocks, 0, 150)
	rows = append(rows, stocks.Stock{Symbol: "HELD", Lastsale: "$10.00", MarketCap: "1", Sector: "Held position"})
	for i := 0; i < 149; i++ {
		rows = append(rows, stocks.Stock{
			Symbol:    fmt.Sprintf("T%03d", i),
			Lastsale:  "$10.00",
			MarketCap: "1000000",
			Sector:    "A deliberately verbose market sector used to exercise prompt limits",
		})
	}

	limited := limitPromptUniverse(holding.Holdings{{Ticket: "HELD", Quantity: 1}}, nil, rows)

	require.Less(t, len(limited), len(rows))
	require.Equal(t, "HELD", limited[0].Symbol)
	require.LessOrEqual(t, len(buildPrompt(holding.Holdings{{Ticket: "HELD", Quantity: 1}}, nil, limited)), maxPromptBytes)
}
