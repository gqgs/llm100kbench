package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gqgs/llminvestbench/pkg/holding"
)

var _ = (Repository)((*repository)(nil))

type repository struct {
	db DB
}

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Repository interface {
	CreateHolding(ctx context.Context, model, ticket string, sum float64, quantity float64) error
	GetHoldings(ctx context.Context, model string) (holding.Holdings, error)
	GetHolding(ctx context.Context, model, ticket string) (*holding.Holding, error)
	SaveHoldings(ctx context.Context, model string, holdings holding.Holdings) error
	SaveContext(ctx context.Context, model, context string) error
	GetRecentContext(ctx context.Context, model string, limit int) ([]string, error)
}

func New(db DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateHolding(ctx context.Context, model, ticket string, sum, quantity float64) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO holdings (model, ticket, sum, quantity) VALUES (?, ?, ?, ?)", model, ticket, sum, quantity)
	return err
}

func (r *repository) GetHoldings(ctx context.Context, model string) (holding.Holdings, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT ticket, sum, quantity, created_at, updated_at FROM holdings WHERE model = ? AND quantity > 0", model)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	holdings := holding.Holdings{}
	for rows.Next() {
		var h holding.Holding
		if err := rows.Scan(&h.Ticket, &h.Sum, &h.Quantity, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		holdings = append(holdings, &h)
	}
	return holdings, nil
}

func (r *repository) GetHolding(ctx context.Context, model, ticket string) (*holding.Holding, error) {
	row := r.db.QueryRowContext(ctx, "SELECT ticket, sum, quantity FROM holdings WHERE ticket = ? AND model = ? AND quantity > 0", ticket, model)
	var holding holding.Holding
	err := row.Scan(&holding.Ticket, &holding.Sum, &holding.Quantity)
	return &holding, err
}

func (r *repository) SaveHoldings(ctx context.Context, model string, holdings holding.Holdings) error {
	for _, h := range holdings {
		if _, err := r.db.ExecContext(ctx, "INSERT INTO holdings (model, ticket, sum, quantity, created_at, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) ON CONFLICT(model, ticket) DO UPDATE SET sum = ?, quantity = ? WHERE model = ? AND ticket = ?", model, h.Ticket, h.Sum, h.Quantity, h.Sum, h.Quantity, model, h.Ticket); err != nil {
			return fmt.Errorf("failed to save holding: %w", err)
		}
	}
	return nil
}

func (r *repository) SaveContext(ctx context.Context, model, context string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO context (model, context, created_at, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) ON CONFLICT(model) DO UPDATE SET context = ?, updated_at = CURRENT_TIMESTAMP WHERE model = ?", model, context, context, model)
	return err
}

func (r *repository) GetRecentContext(ctx context.Context, model string, limit int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT context FROM context WHERE model = ? ORDER BY created_at DESC LIMIT ?", model, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contexts := []string{}
	for rows.Next() {
		var context string
		if err := rows.Scan(&context); err != nil {
			return nil, err
		}
		contexts = append(contexts, context)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return contexts, nil
}
