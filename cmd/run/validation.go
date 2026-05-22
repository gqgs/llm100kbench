package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gqgs/llminvestbench/pkg/holding"
	"github.com/gqgs/llminvestbench/pkg/order"
	"github.com/gqgs/llminvestbench/pkg/stocks"
)

func parseOrder(raw string) (*order.Order, error) {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "```") {
		raw = strings.TrimPrefix(raw, "```json")
		raw = strings.TrimPrefix(raw, "```")
		raw = strings.TrimSuffix(raw, "```")
		raw = strings.TrimSpace(raw)
	}
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start < 0 || end < start {
		return nil, fmt.Errorf("%w: response did not contain a JSON object", errInvalidOrder)
	}
	raw = raw[start : end+1]

	var parsed order.Order
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil, fmt.Errorf("%w: %v", errInvalidOrder, err)
	}
	if parsed.Updates == nil {
		parsed.Updates = []*order.Update{}
	}
	return &parsed, nil
}

func validateOrder(parsed *order.Order, holdings holding.Holdings, priceMap map[string]float64, universe stocks.Stocks) error {
	if parsed == nil {
		return fmt.Errorf("%w: order is nil", errInvalidOrder)
	}
	if parsed.Context == nil || len(parsed.Context) == 0 {
		return fmt.Errorf("%w: context must include at least one entry", errInvalidOrder)
	}

	tickers := marketTickers(universe)
	quantities := map[string]int{}
	for _, h := range holdings {
		quantities[h.Ticket] = h.Quantity
	}

	cash := float64(quantities["USD"])
	for i, update := range parsed.Updates {
		if update == nil {
			return fmt.Errorf("%w: update %d is nil", errInvalidOrder, i)
		}
		update.Ticket = strings.ToUpper(strings.TrimSpace(update.Ticket))
		update.Action = strings.ToUpper(strings.TrimSpace(update.Action))
		if update.Quantity <= 0 {
			return fmt.Errorf("%w: %s quantity must be positive", errInvalidOrder, update.Ticket)
		}
		if _, ok := tickers[update.Ticket]; !ok {
			return fmt.Errorf("%w: %s is not in the market table", errInvalidOrder, update.Ticket)
		}
		price, ok := priceMap[update.Ticket]
		if !ok {
			return fmt.Errorf("%w: missing price for %s", errInvalidOrder, update.Ticket)
		}
		if update.Price <= 0 || abs(update.Price-price) > 0.01 {
			return fmt.Errorf("%w: %s price %.4f does not match market price %.4f", errInvalidOrder, update.Ticket, update.Price, price)
		}

		value := update.Price * float64(update.Quantity)
		switch update.Action {
		case "SELL":
			if quantities[update.Ticket] < update.Quantity {
				return fmt.Errorf("%w: cannot sell %d %s with %d held", errInvalidOrder, update.Quantity, update.Ticket, quantities[update.Ticket])
			}
			quantities[update.Ticket] -= update.Quantity
			cash += value
		case "BUY":
			if cash+0.01 < value {
				return fmt.Errorf("%w: insufficient USD to buy %s; need %s have %s", errInvalidOrder, update.Ticket, formatMoney(value), formatMoney(cash))
			}
			quantities[update.Ticket] += update.Quantity
			cash -= value
		default:
			return fmt.Errorf("%w: unsupported action %q", errInvalidOrder, update.Action)
		}
	}
	return nil
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
