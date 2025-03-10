package main

import (
	"encoding/json"
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

	encoded, err := json.MarshalIndent(stocks, "", " ")
	if err != nil {
		return fmt.Errorf("failled to encode: %w", err)
	}

	return os.WriteFile(opts.output, encoded, os.ModePerm)
}
