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

type Data struct {
	Rows Stocks `json:"rows"`
}
type NasdaqResponse struct {
	Data Data `json:"data"`
}

func DecodeNasdaqResponse(reader io.Reader) (Stocks, error) {
	var response NasdaqResponse
	return response.Data.Rows, json.NewDecoder(reader).Decode(&response)
}
