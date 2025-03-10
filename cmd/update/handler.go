package main

import (
	"context"
	"fmt"

	"github.com/gqgs/llminvestbench/pkg/manager"
	"github.com/gqgs/llminvestbench/pkg/order"
	"github.com/gqgs/llminvestbench/pkg/service"
	"github.com/gqgs/llminvestbench/pkg/storage"
)

func handler(ctx context.Context, opts options) error {
	storage, err := storage.NewSqlite(opts.db)
	if err != nil {
		return err
	}
	defer storage.Close()

	order, err := order.Parse(opts.order)
	if err != nil {
		return fmt.Errorf("failed parsing order: %w", err)
	}

	service := service.New(storage)
	manager := manager.New(service, opts.model)
	holdings, err := manager.GetHoldings(ctx)
	if err != nil {
		return fmt.Errorf("failed getting holdings: %w", err)
	}

	if err := holdings.ProcessOrder(order); err != nil {
		return fmt.Errorf("failed processing order: %w", err)
	}

	if err := manager.Save(ctx, holdings, order.Context[len(order.Context)-1]); err != nil {
		return fmt.Errorf("failed saving: %w", err)
	}

	return nil
}
