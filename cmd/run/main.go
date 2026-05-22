package main

import (
	"context"
	"flag"
	"log"
	"path/filepath"
	"time"
)

const defaultURL = "https://api.nasdaq.com/api/screener/stocks"

type options struct {
	db         string
	models     string
	pricesURL  string
	date       string
	maxSymbols int
}

func main() {
	opts := options{}
	flag.StringVar(&opts.db, "db", "llm100kbench.db", "database name (sqlite)")
	flag.StringVar(&opts.models, "models", "models.json", "model config file")
	flag.StringVar(&opts.pricesURL, "prices-url", defaultURL, "url to get prices from")
	flag.StringVar(&opts.date, "date", time.Now().Format(time.DateOnly), "run date")
	flag.IntVar(&opts.maxSymbols, "max-symbols", 120, "maximum market symbols to include in prompts")
	flag.Parse()

	if err := run(context.Background(), opts); err != nil {
		log.Fatal(err)
	}
}

func pricePath(date string) string {
	return filepath.Join("prices", date+".csv")
}

func orderPath(alias, date string) string {
	return filepath.Join("orders", alias, date+".json")
}

func logPath(alias, date string) string {
	return filepath.Join("logs", alias, date+".md")
}
