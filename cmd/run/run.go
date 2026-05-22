package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gqgs/llminvestbench/pkg/holding"
	"github.com/gqgs/llminvestbench/pkg/llm"
	"github.com/gqgs/llminvestbench/pkg/manager"
	"github.com/gqgs/llminvestbench/pkg/modelconfig"
	"github.com/gqgs/llminvestbench/pkg/order"
	"github.com/gqgs/llminvestbench/pkg/service"
	"github.com/gqgs/llminvestbench/pkg/stocks"
	"github.com/gqgs/llminvestbench/pkg/storage"
)

func run(ctx context.Context, opts options) error {
	cfg, err := modelconfig.Load(opts.models)
	if err != nil {
		return err
	}
	if err := requireSecrets(cfg.EnabledModels()); err != nil {
		return err
	}

	store, err := storage.NewSqlite(opts.db)
	if err != nil {
		return err
	}
	defer store.Close()

	rows, err := stocks.FetchNasdaqStocks(opts.pricesURL)
	if err != nil {
		return err
	}
	if err := os.MkdirAll("prices", 0o755); err != nil {
		return err
	}
	if err := stocks.WriteCSV(pricePath(opts.date), rows); err != nil {
		return err
	}

	priceMap := buildPriceMap(rows)
	svc := service.New(store)
	for _, model := range cfg.EnabledModels() {
		if err := runModel(ctx, svc, model, rows, priceMap, opts); err != nil {
			return fmt.Errorf("%s: %w", model.Alias, err)
		}
	}

	if err := writeStats(ctx, svc, cfg, opts.date); err != nil {
		return err
	}
	return updateReadme(opts.date)
}

func requireSecrets(models []modelconfig.Model) error {
	missing := []string{}
	seen := map[string]struct{}{}
	for _, model := range models {
		if model.Env == "" {
			continue
		}
		if _, ok := seen[model.Env]; ok {
			continue
		}
		seen[model.Env] = struct{}{}
		if os.Getenv(model.Env) == "" {
			missing = append(missing, model.Env)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required model API environment variables: %v", missing)
	}
	return nil
}

func runModel(ctx context.Context, svc service.Service, model modelconfig.Model, rows stocks.Stocks, priceMap map[string]float64, opts options) error {
	mgr := manager.New(svc, model.Alias)
	holdings, err := mgr.GetHoldings(ctx)
	if err != nil {
		return err
	}
	if len(holdings) == 0 {
		if err := mgr.CreateHoldings(ctx); err != nil {
			return err
		}
		holdings, err = mgr.GetHoldings(ctx)
		if err != nil {
			return err
		}
	}

	contexts, err := mgr.GetRecentContext(ctx)
	if err != nil {
		return err
	}
	universe := selectUniverse(rows, holdings, opts.maxSymbols)
	prompt := buildPrompt(holdings, contexts, universe)

	client, err := llm.New(model, os.Getenv(model.Env))
	if err != nil {
		return err
	}

	parsed, notes := generateOrder(ctx, client, prompt, holdings, priceMap, universe)
	parsed.Metadata = &order.Metadata{
		Alias:       model.Alias,
		Provider:    model.Provider,
		Model:       model.Model,
		Status:      "ok",
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Notes:       notes,
	}
	if len(notes) > 0 && notes[len(notes)-1] == "model failed validation after retry; no trades executed" {
		parsed.Metadata.Status = "failed"
	}

	if err := os.MkdirAll(filepath.Join("orders", model.Alias), 0o755); err != nil {
		return err
	}
	if err := writeOrder(orderPath(model.Alias, opts.date), parsed); err != nil {
		return err
	}

	if err := holdings.ProcessOrder(parsed); err != nil {
		return err
	}
	revalueHoldings(holdings, priceMap)

	orderContext := fmt.Sprintf("%s run completed with %d updates.", model.Model, len(parsed.Updates))
	if len(parsed.Context) > 0 {
		orderContext = parsed.Context[len(parsed.Context)-1]
	}
	return mgr.Save(ctx, holdings, orderContext)
}

func generateOrder(ctx context.Context, client llm.Client, prompt string, holdings holding.Holdings, priceMap map[string]float64, universe stocks.Stocks) (*order.Order, []string) {
	notes := []string{}
	raw, err := client.Generate(ctx, prompt)
	if err == nil {
		parsed, parseErr := parseOrder(raw)
		if parseErr == nil {
			if validateErr := validateOrder(parsed, holdings, priceMap, universe); validateErr == nil {
				return parsed, notes
			} else {
				err = validateErr
			}
		} else {
			err = parseErr
		}
	}

	notes = append(notes, "first response rejected: "+err.Error())
	retryPrompt := prompt + "\n\nYour previous response was rejected: " + err.Error() + "\nReturn corrected JSON only."
	raw, err = client.Generate(ctx, retryPrompt)
	if err == nil {
		parsed, parseErr := parseOrder(raw)
		if parseErr == nil {
			if validateErr := validateOrder(parsed, holdings, priceMap, universe); validateErr == nil {
				return parsed, notes
			} else {
				err = validateErr
			}
		} else {
			err = parseErr
		}
	}

	notes = append(notes, "retry response rejected: "+err.Error())
	notes = append(notes, "model failed validation after retry; no trades executed")
	return &order.Order{
		Updates: []*order.Update{},
		Context: []string{"No trades executed because the model did not return a valid order."},
	}, notes
}

func writeOrder(path string, parsed *order.Order) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(parsed)
}

func buildPriceMap(rows stocks.Stocks) map[string]float64 {
	prices := map[string]float64{"USD": 1}
	for _, row := range rows {
		price, err := stocks.CleanPrice(row.Lastsale)
		if err == nil && price > 0 {
			prices[row.Symbol] = price
		}
	}
	return prices
}

func revalueHoldings(holdings holding.Holdings, prices map[string]float64) {
	for _, h := range holdings {
		if price, ok := prices[h.Ticket]; ok {
			h.Sum = price * float64(h.Quantity)
		}
	}
}

var errInvalidOrder = errors.New("invalid order")
