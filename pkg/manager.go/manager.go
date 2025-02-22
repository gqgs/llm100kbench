package manager

import (
	"context"

	"github.com/gqgs/llminvestbench/pkg/holding"
	"github.com/gqgs/llminvestbench/pkg/storage"
)

const (
	defaultTicket   = "USD"
	defaultSum      = 100_000
	defaultQuantity = 100_000
)

type Manager struct {
	storage storage.Storage
	model   string
}

func New(storage storage.Storage) *Manager {
	return &Manager{
		storage: storage,
	}
}

func (m *Manager) CreateHoldings(ctx context.Context) error {
	return m.storage.CreateHolding(ctx, defaultTicket, defaultSum, defaultQuantity)
}

func (m *Manager) GetHoldings(ctx context.Context) (holding.Holdings, error) {
	return m.storage.GetHoldings(ctx)
}

func (m *Manager) GetHolding(ctx context.Context, ticket string) (*holding.Holding, error) {
	return m.storage.GetHolding(ctx, ticket)
}

func (m *Manager) SaveHoldings(ctx context.Context, holdings holding.Holdings) error {
	return m.storage.SaveHoldings(ctx, holdings)
}
