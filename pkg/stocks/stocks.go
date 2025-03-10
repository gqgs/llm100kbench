package stocks

import (
	"encoding/json"
	"io"
)

type Stocks []Stock

type Stock struct {
	Symbol   string `json:"symbol"`
	Lastsale string `json:"lastsale"`
}

func DecodeTickers(reader io.Reader) (Stocks, error) {
	var stocks Stocks
	return stocks, json.NewDecoder(reader).Decode(&stocks)
}
