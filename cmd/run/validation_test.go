package main

import (
	"context"
	"errors"
	"testing"

	"github.com/gqgs/llminvestbench/pkg/holding"
	"github.com/gqgs/llminvestbench/pkg/order"
	"github.com/gqgs/llminvestbench/pkg/stocks"
	"github.com/stretchr/testify/require"
)

type fakeClient struct {
	responses []string
	errs      []error
	calls     int
}

func (f *fakeClient) Generate(ctx context.Context, prompt string) (string, error) {
	idx := f.calls
	f.calls++
	if idx < len(f.errs) && f.errs[idx] != nil {
		return "", f.errs[idx]
	}
	return f.responses[idx], nil
}

func TestParseOrderExtractsJSONFromMarkdown(t *testing.T) {
	parsed, err := parseOrder("```json\n{\"updates\":[],\"context\":[\"hold\"]}\n```")
	require.NoError(t, err)
	require.Empty(t, parsed.Updates)
	require.Equal(t, []string{"hold"}, parsed.Context)
}

func TestValidateOrderRejectsOverspend(t *testing.T) {
	holdings := holding.Holdings{{Ticket: "USD", Sum: 100, Quantity: 100}}
	universe := stocks.Stocks{{Symbol: "AAPL", Lastsale: "$50.00"}}
	err := validateOrder(&order.Order{
		Updates: []*order.Update{{Ticket: "AAPL", Quantity: 3, Price: 50, Action: "BUY"}},
		Context: []string{"buy"},
	}, holdings, map[string]float64{"USD": 1, "AAPL": 50}, universe)

	require.Error(t, err)
	require.ErrorIs(t, err, errInvalidOrder)
}

func TestValidateOrderAcceptsSellThenBuy(t *testing.T) {
	holdings := holding.Holdings{
		{Ticket: "USD", Sum: 0, Quantity: 0},
		{Ticket: "MSFT", Sum: 100, Quantity: 1},
	}
	universe := stocks.Stocks{
		{Symbol: "MSFT", Lastsale: "$100.00"},
		{Symbol: "AAPL", Lastsale: "$50.00"},
	}
	err := validateOrder(&order.Order{
		Updates: []*order.Update{
			{Ticket: "MSFT", Quantity: 1, Price: 100, Action: "SELL"},
			{Ticket: "AAPL", Quantity: 2, Price: 50, Action: "BUY"},
		},
		Context: []string{"rebalance"},
	}, holdings, map[string]float64{"USD": 1, "MSFT": 100, "AAPL": 50}, universe)

	require.NoError(t, err)
}

func TestGenerateOrderFallsBackToNoTrade(t *testing.T) {
	client := &fakeClient{errs: []error{errors.New("boom"), errors.New("still boom")}}
	parsed, notes := generateOrder(context.Background(), client, "prompt", holding.Holdings{{Ticket: "USD", Quantity: 100}}, map[string]float64{"USD": 1}, nil)

	require.Empty(t, parsed.Updates)
	require.Len(t, notes, 3)
	require.Equal(t, 2, client.calls)
}
