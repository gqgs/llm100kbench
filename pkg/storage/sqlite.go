package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var _ = (Storage)((*sqliteStorage)(nil))

type sqliteStorage struct {
	db *sql.DB
}

// NewSqlite creates a new sqlite storage
// It should be closed after being used
func NewSqlite(dbPath string) (*sqliteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_synchronous=off&_journal_mode=off&cache=shared")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS context (
			model TEXT,
			context TEXT,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (model)
		)
	`); err != nil {
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

	for _, index := range []string{
		"CREATE INDEX IF NOT EXISTS idx_ticket on holdings(ticket)",
		"CREATE INDEX IF NOT EXISTS idx_model on holdings(model)",
		"CREATE INDEX IF NOT EXISTS idx_model_quantity on holdings(model, quantity)",
		"CREATE INDEX IF NOT EXISTS idx_positions on holdings(model, ticket, quantity)",
	} {
		if _, err := db.Exec(index); err != nil {
			return nil, err
		}
	}

	storage := &sqliteStorage{
		db: db,
	}

	return storage, nil
}

func (s *sqliteStorage) DB() *sql.DB {
	return s.db
}

func (s *sqliteStorage) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *sqliteStorage) Close() error {
	return s.db.Close()
}
