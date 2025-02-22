package main

import (
	"context"
	"log"
)

//go:generate go run github.com/gqgs/argsgen

type options struct {
	db    string `arg:"database name (sqlite),required"`
	model string `arg:"name of the model,required"`
	order string `arg:"path to the order file,required"`
}

func main() {
	opts := options{
		db: "llminvestbench.db",
	}
	opts.MustParse()

	if err := handler(context.Background(), opts); err != nil {
		log.Fatal(err)
	}
}
