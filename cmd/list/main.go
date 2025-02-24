package main

import (
	"cmp"
	"context"
	"log"
	"os"
)

//go:generate go tool argsgen

type options struct {
	db        string `arg:"database name (sqlite),required"`
	model     string `arg:"name of the model,required"`
	roundsums bool   `arg:"if it should round sums"`
}

func main() {
	opts := options{
		db: cmp.Or(os.Getenv("LLM_BENCH_DATABASE"), "llm100kbench.db"),
	}
	opts.MustParse()

	if err := handler(context.Background(), opts); err != nil {
		log.Fatal(err)
	}
}
