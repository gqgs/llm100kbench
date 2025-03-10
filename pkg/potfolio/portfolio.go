package potfolio

import (
	"encoding/json"

	"github.com/gqgs/llminvestbench/pkg/holding"
)

type portifolio struct {
	Holdings holding.Holdings `json:"holdings"`
	Context  []string         `json:"context"`
}

func New(holdings holding.Holdings, context []string) *portifolio {
	return &portifolio{
		Holdings: holdings,
		Context:  context,
	}
}

func (p portifolio) String() string {
	encoded, _ := json.MarshalIndent(p, "", " ")
	return string(encoded)
}
