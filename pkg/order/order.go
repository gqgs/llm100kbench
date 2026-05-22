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
	Updates  []*Update `json:"updates"`
	Context  []string  `json:"context"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	Alias       string   `json:"alias,omitempty"`
	Provider    string   `json:"provider,omitempty"`
	Model       string   `json:"model,omitempty"`
	Status      string   `json:"status,omitempty"`
	GeneratedAt string   `json:"generated_at,omitempty"`
	Notes       []string `json:"notes,omitempty"`
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
