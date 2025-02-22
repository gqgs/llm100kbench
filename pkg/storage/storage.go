package storage

import (
	"context"
	"io"

	"github.com/gqgs/llminvestbench/pkg/holding"
)

type Storage interface {
	io.Closer
	CreateHolding(ctx context.Context, ticket string, sum float64, quantity float64) error
	GetHoldings(ctx context.Context) (holding.Holdings, error)
	GetHolding(ctx context.Context, ticket string) (*holding.Holding, error)
	SaveHoldings(ctx context.Context, holdings holding.Holdings) error
	SaveContext(ctx context.Context, context string) error
	GetRecentContext(ctx context.Context, limit int) ([]string, error)
}

var storage Storage
