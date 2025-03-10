package manager

import (
	"context"

	"github.com/gqgs/llminvestbench/pkg/holding"
	"github.com/gqgs/llminvestbench/pkg/repository"
	"github.com/gqgs/llminvestbench/pkg/service"
)

const (
	defaultTicket   = "USD"
	defaultSum      = 100_000
	defaultQuantity = 100_000
	contextLimit    = 3
)

type Manager struct {
	service service.Service
	model   string
}

func New(service service.Service, model string) *Manager {
	return &Manager{
		service: service,
		model:   model,
	}
}

func (m *Manager) CreateHoldings(ctx context.Context) error {
	return m.service.Exec(func(r repository.Repository) error {
		return r.CreateHolding(ctx, m.model, defaultTicket, defaultSum, defaultQuantity)
	})
}

func (m *Manager) GetHoldings(ctx context.Context) (holdings holding.Holdings, err error) {
	m.service.Exec(func(r repository.Repository) error {
		holdings, err = r.GetHoldings(ctx, m.model)
		return err
	})
	return
}

func (m *Manager) GetHolding(ctx context.Context, ticket string) (holding *holding.Holding, err error) {
	m.service.Exec(func(r repository.Repository) error {
		holding, err = r.GetHolding(ctx, m.model, ticket)
		return err
	})
	return
}

func (m *Manager) Save(ctx context.Context, holdings holding.Holdings, orderCtx string) error {
	return m.service.ExecTx(func(r repository.Repository) error {
		if err := r.SaveHoldings(ctx, m.model, holdings); err != nil {
			return err
		}
		if err := r.SaveContext(ctx, m.model, orderCtx); err != nil {
			return err
		}
		return nil
	})
}

func (m *Manager) GetRecentContext(ctx context.Context) (context []string, err error) {
	m.service.Exec(func(r repository.Repository) error {
		context, err = r.GetRecentContext(ctx, m.model, contextLimit)
		return err
	})
	return
}
