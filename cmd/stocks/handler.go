package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"

	"github.com/gqgs/llminvestbench/pkg/stocks"
)

func handler(opts options) error {
	resp, err := http.Get(opts.url)
	if err != nil {
		return fmt.Errorf("failed to get url: %w", err)
	}
	defer resp.Body.Close()

	stocks, err := stocks.DecodeTickers(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to decode tickers: %w", err)
	}

	file, err := os.Create(opts.output)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	records := make([][]string, 0, len(stocks)+1)
	records = append(records, []string{"ticket", "price"})
	for _, stock := range stocks {
		records = append(records, []string{stock.Symbol, stock.Lastsale})
	}

	return csv.NewWriter(file).WriteAll(records)
}
