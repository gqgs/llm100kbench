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
		return fmt.Errorf("failed to open storage: %w", err)
	}
	defer storage.Close()

	service := service.New(storage)
	manager := manager.New(service, opts.model)
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
