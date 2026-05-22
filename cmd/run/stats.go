package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gqgs/llminvestbench/pkg/holding"
	"github.com/gqgs/llminvestbench/pkg/manager"
	"github.com/gqgs/llminvestbench/pkg/modelconfig"
	"github.com/gqgs/llminvestbench/pkg/service"
)

type modelTotal struct {
	Alias string
	Total float64
}

func writeStats(ctx context.Context, svc service.Service, cfg *modelconfig.Config, date string) error {
	if err := os.MkdirAll("stats", 0o755); err != nil {
		return err
	}

	var builder strings.Builder
	builder.WriteString("| Model | Ticket | Sum | Quantity |\n")
	builder.WriteString("|-------|-------|-------|--------|\n")

	totals := []modelTotal{}
	for _, model := range cfg.EnabledModels() {
		mgr := manager.New(svc, model.Alias)
		holdings, err := mgr.GetHoldings(ctx)
		if err != nil {
			return err
		}
		sortHoldings(holdings)
		var total float64
		for _, h := range holdings {
			builder.WriteString(fmt.Sprintf("|`%s`|`%s`|%.0f|%d|\n", model.Alias, h.Ticket, h.Sum, h.Quantity))
			total += h.Sum
		}
		totals = append(totals, modelTotal{Alias: model.Alias, Total: total})
	}

	builder.WriteString("\n\n")
	builder.WriteString("| Model | Total Sum | Change |\n")
	builder.WriteString("|-------|-----------|--------|\n")
	sort.SliceStable(totals, func(i, j int) bool {
		return totals[i].Total > totals[j].Total
	})
	for _, total := range totals {
		builder.WriteString(fmt.Sprintf("|`%s`|%.0f|—|\n", total.Alias, total.Total))
	}

	return os.WriteFile(filepath.Join("stats", date+".md"), []byte(builder.String()), 0o644)
}

func sortHoldings(holdings holding.Holdings) {
	sort.SliceStable(holdings, func(i, j int) bool {
		if holdings[i].Ticket == "USD" {
			return true
		}
		if holdings[j].Ticket == "USD" {
			return false
		}
		return holdings[i].Ticket < holdings[j].Ticket
	})
}

func updateReadme(date string) error {
	statsPath := filepath.Join("stats", date+".md")
	content, err := os.ReadFile(statsPath)
	if err != nil {
		return err
	}

	readme, err := os.ReadFile("README.md")
	if err != nil {
		return err
	}

	header := "## Current Portfolio (" + date + ")"
	next := header + "\n\n" + string(content)
	text := string(readme)
	idx := strings.Index(text, "\n## Current Portfolio")
	if idx < 0 {
		text = strings.TrimRight(text, "\n") + "\n\n" + next + "\n"
	} else {
		text = strings.TrimRight(text[:idx], "\n") + "\n\n" + next
	}
	return os.WriteFile("README.md", []byte(text), 0o644)
}
