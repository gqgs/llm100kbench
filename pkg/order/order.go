package order

import (
	"encoding/json"
	"os"
)

type Update struct {
	Ticket   string  `json:"ticket"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Action   string  `json:"action"`
}

type Order struct {
	Updates []*Update `json:"updates"`
	Context []string  `json:"context"`
}

func Parse(orderFile string) (*Order, error) {
	order := new(Order)

	file, err := os.Open(orderFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(order); err != nil {
		return nil, err
	}

	return order, nil
}
