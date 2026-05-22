package stocks

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Stocks []Stock

type Stock struct {
	Symbol    string `json:"symbol"`
	Lastsale  string `json:"lastsale"`
	Name      string `json:"name"`
	MarketCap string `json:"marketCap"`
	Volume    string `json:"volume"`
	PctChange string `json:"pctchange"`
	Sector    string `json:"sector"`
	Industry  string `json:"industry"`
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

func FetchNasdaqStocks(rawurl string) (Stocks, error) {
	values := make(url.Values)
	values.Set("tableonly", "true")
	values.Set("limit", "25")
	values.Set("offset", "0")
	values.Set("exchange", "nasdaq")
	values.Set("download", "true")

	req, err := http.NewRequest(http.MethodGet, rawurl+"?"+values.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:85.0) Gecko/20100101 Firefox/85.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status from nasdaq: %s", resp.Status)
	}

	stocks, err := DecodeNasdaqResponse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tickers: %w", err)
	}
	return stocks, nil
}

func WriteCSV(path string, stocks Stocks) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"ticket", "price", "market_cap", "sector", "industry", "volume", "pct_change"}); err != nil {
		return err
	}
	for _, stock := range stocks {
		if err := writer.Write([]string{stock.Symbol, stock.Lastsale, stock.MarketCap, stock.Sector, stock.Industry, stock.Volume, stock.PctChange}); err != nil {
			return err
		}
	}
	return writer.Error()
}

func CleanPrice(price string) (float64, error) {
	cleaned := strings.TrimSpace(price)
	cleaned = strings.TrimPrefix(cleaned, "$")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	if cleaned == "" {
		return 0, fmt.Errorf("empty price")
	}
	return strconv.ParseFloat(cleaned, 64)
}

func CleanNumber(value string) float64 {
	cleaned := strings.TrimSpace(value)
	cleaned = strings.TrimSuffix(cleaned, "%")
	cleaned = strings.TrimPrefix(cleaned, "$")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	parsed, _ := strconv.ParseFloat(cleaned, 64)
	return parsed
}
