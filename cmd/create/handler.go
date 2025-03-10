package main

import (
	"context"
	_ "embed"
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
		return fmt.Errorf("failed to open storage: %w", err)
	}
	defer storage.Close()

	manager := manager.New(storage)
	if err = manager.CreateHoldings(ctx); err != nil {
		return fmt.Errorf("failed creating holdings: %w", err)
	}

	fmt.Println(prompt)

	holdings, err := manager.GetHoldings(ctx)
	if err != nil {
		return fmt.Errorf("failed getting holdings: %w", err)
	}

	fmt.Println("Your current portfolio:")
	fmt.Println("```")

	fmt.Println(potfolio.New(holdings, []string{"Initial holdings"}))
	fmt.Println("```")

	return nil
}
