package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/gqgs/llminvestbench/pkg/manager"
	"github.com/gqgs/llminvestbench/pkg/potfolio"
	"github.com/gqgs/llminvestbench/pkg/storage"
)

//go:embed prompt.txt
var prompt string

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

	if opts.prompt {
		fmt.Println(prompt)

		fmt.Println("Your current portfolio:")
	}

	encoded, err := json.MarshalIndent(potfolio.New(holdings, contexts), "", " ")
	if err != nil {
		return fmt.Errorf("error marshaling file: %w", err)
	}

	fmt.Println(string(encoded))

	return nil
}
