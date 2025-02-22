package potfolio

import "github.com/gqgs/llminvestbench/pkg/holding"

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
