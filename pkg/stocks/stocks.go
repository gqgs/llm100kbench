package stocks

import (
	"encoding/json"
	"io"
)

type Stocks []Stock

type Stock struct {
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	Lastsale  string `json:"lastsale"`
	Volume    string `json:"volume"`
	MarketCap string `json:"marketCap"`
	Country   string `json:"country"`
	Ipoyear   string `json:"ipoyear"`
	Industry  string `json:"industry"`
	Sector    string `json:"sector"`
}

func DecodeTickers(reader io.Reader) (Stocks, error) {
	var stocks Stocks
	return stocks, json.NewDecoder(reader).Decode(&stocks)
}
