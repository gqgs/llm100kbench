package main

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/gqgs/llminvestbench/pkg/manager"
	"github.com/gqgs/llminvestbench/pkg/potfolio"
	"github.com/gqgs/llminvestbench/pkg/service"
	"github.com/gqgs/llminvestbench/pkg/storage"
)

//go:embed prompt.txt
var prompt string

func handler(ctx context.Context, opts options) error {
	storage, err := storage.NewSqlite(opts.db)
	if err != nil {
		return err
	}
	defer storage.Close()

	service := service.New(storage)
	manager := manager.New(service, opts.model)
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
		fmt.Println("```")
	}

	fmt.Println(potfolio.New(holdings, contexts))

	if opts.prompt {
		fmt.Println("```")
	}

	return nil
}
