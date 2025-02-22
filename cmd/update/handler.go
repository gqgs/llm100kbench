package main

import (
	"context"
	"fmt"

	"github.com/gqgs/llminvestbench/pkg/manager.go"
	"github.com/gqgs/llminvestbench/pkg/order"
	"github.com/gqgs/llminvestbench/pkg/storage"
)

func handler(ctx context.Context, opts options) error {
	storage, err := storage.NewSqlite(opts.db, opts.model)
	if err != nil {
		return err
	}
	defer storage.Close()

	order, err := order.Parse(opts.order)
	if err != nil {
		return fmt.Errorf("failed parsing order: %w", err)
	}

	manager := manager.New(storage)
	holdings, err := manager.GetHoldings(ctx)
	if err != nil {
		return fmt.Errorf("failed getting holdings: %w", err)
	}

	if err := holdings.ProcessOrder(order); err != nil {
		return fmt.Errorf("failed processing order: %w", err)
	}

	// FIXME: ideally both would be executed in a transaction
	if err := manager.SaveHoldings(ctx, holdings); err != nil {
		return fmt.Errorf("failed saving holdings: %w", err)
	}

	if err := manager.SaveContext(ctx, order.Context[len(order.Context)-1]); err != nil {
		return fmt.Errorf("failed saving context: %w", err)
	}

	return nil
}
