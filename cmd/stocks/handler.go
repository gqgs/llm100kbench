package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gqgs/llminvestbench/pkg/stocks"
)

func handler(opts options) error {
	values := make(url.Values)
	values.Set("tableonly", "true")
	values.Set("limit", "25")
	values.Set("offset", "0")
	values.Set("exchange", "nasdaq")
	values.Set("download", "true")
	req, err := http.NewRequest(http.MethodGet, opts.url+"?"+values.Encode(), nil)
	if err != nil {
		return fmt.Errorf("failed creating request: %w", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:85.0) Gecko/20100101 Firefox/85.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get url: %w", err)
	}
	defer resp.Body.Close()

	stocks, err := stocks.DecodeNasdaqResponse(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to decode tickers: %w", err)
	}

	records := make([][]string, 0, len(stocks)+1)
	records = append(records, []string{"ticket", "price"})
	for _, stock := range stocks {
		records = append(records, []string{stock.Symbol, stock.Lastsale})
	}

	file, err := os.Create(opts.output)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err = csv.NewWriter(file).WriteAll(records); err != nil {
		return fmt.Errorf("failed to create csv file: %w", err)
	}

	return nil
}
