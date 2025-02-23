package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gqgs/llminvestbench/pkg/manager"
	"github.com/gqgs/llminvestbench/pkg/potfolio"
	"github.com/gqgs/llminvestbench/pkg/storage"
)

func handler(ctx context.Context, opts options) error {
	storage, err := storage.NewSqlite(opts.db, opts.model)
	if err != nil {
		return err
	}
	defer storage.Close()

	manager := manager.New(storage)
	holdings, err := manager.GetHoldings(ctx)
	if err != nil {
		return fmt.Errorf("failed getting holdings: %w", err)
	}

	if opts.roundsums {
		holdings.RoundSums()
	}

	contexts, err := manager.GetRecentContext(ctx)
	if err != nil {
		return fmt.Errorf("failed getting recent context: %w", err)
	}

	return json.NewEncoder(os.Stdout).Encode(potfolio.New(holdings, contexts))
}
