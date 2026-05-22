package main

import (
	"fmt"

	"github.com/gqgs/llminvestbench/pkg/stocks"
)

func handler(opts options) error {
	rows, err := stocks.FetchNasdaqStocks(opts.url)
	if err != nil {
		return err
	}

	if err = stocks.WriteCSV(opts.output, rows); err != nil {
		return fmt.Errorf("failed to create csv file: %w", err)
	}

	return nil
}
