package main

import (
	"context"
	"log"
)

//go:generate go tool argsgen

type options struct {
	db    string `arg:"database name (sqlite),required"`
	model string `arg:"name of the model,required"`
}

func main() {
	opts := options{
		db: "llm100kbench.db",
	}
	opts.MustParse()

	if err := handler(context.Background(), opts); err != nil {
		log.Fatal(err)
	}
}
