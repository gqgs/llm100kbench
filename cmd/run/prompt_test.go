package main

import (
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
