package main

import (
	"log"
	"path/filepath"
	"time"
)

const (
	defaultURL = "https://api.nasdaq.com/api/screener/stocks"
)

//go:generate go tool argsgen

type options struct {
	url    string `arg:"url to get the tickets from,required"`
	output string `arg:"output file,required"`
}

func main() {
	opts := options{
		url:    defaultURL,
		output: filepath.Join("prices", time.Now().Format(time.DateOnly)+".csv"),
	}
	opts.MustParse()

	if err := handler(opts); err != nil {
		log.Fatal(err)
	}
}
