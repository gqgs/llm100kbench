package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gqgs/llminvestbench/pkg/holding"
	_ "github.com/mattn/go-sqlite3"
)

var _ = (Storage)((*sqliteStorage)(nil))

type sqliteStorage struct {
	db    *sql.DB
	model string
}

// NewSqlite creates a new sqlite storage instace namescaped to the model
// It should be closed after being used
func NewSqlite(dbPath, model string) (*sqliteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_synchronous=off&_journal_mode=off&cache=shared")
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS holdings (
			model TEXT,
			ticket TEXT,
			sum REAL,
			quantity REAL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (model, ticket)
		)
	`); err != nil {
		return nil, err
	}

	if _, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_ticket on holdings(ticket)"); err != nil {
		return nil, err
	}

	if _, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_model on holdings(model)"); err != nil {
		return nil, err
	}

	storage = &sqliteStorage{
		db:    db,
		model: model,
	}

	return storage.(*sqliteStorage), nil
}

func (s *sqliteStorage) Close() error {
	return s.db.Close()
}

func (s *sqliteStorage) CreateHolding(ctx context.Context, ticket string, sum, quantity float64) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO holdings (model, ticket, sum, quantity) VALUES (?, ?, ?, ?)", s.model, ticket, sum, quantity)
	return err
}

func (s *sqliteStorage) GetHoldings(ctx context.Context) (holding.Holdings, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT ticket, sum, quantity, created_at, updated_at FROM holdings WHERE model = ?", s.model)
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

func (s *sqliteStorage) GetHolding(ctx context.Context, ticket string) (*holding.Holding, error) {
	row := s.db.QueryRowContext(ctx, "SELECT ticket, sum, quantity FROM holdings WHERE ticket = ? AND model = ?", ticket, s.model)
	var holding holding.Holding
	err := row.Scan(&holding.Ticket, &holding.Sum, &holding.Quantity)
	return &holding, err
}

func (s *sqliteStorage) SaveHoldings(ctx context.Context, holdings holding.Holdings) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, h := range holdings {
		if _, err := tx.ExecContext(ctx, "INSERT INTO holdings (model, ticket, sum, quantity, created_at, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) ON CONFLICT(model, ticket) DO UPDATE SET sum = ?, quantity = ? WHERE model = ? AND ticket = ?", s.model, h.Ticket, h.Sum, h.Quantity, h.Sum, h.Quantity, s.model, h.Ticket); err != nil {
			return fmt.Errorf("failed to save holding: %w", err)
		}
	}
	return tx.Commit()
}
