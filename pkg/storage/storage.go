package storage

import (
	"database/sql"
	"io"
)

type Storage interface {
	io.Closer
	Begin() (*sql.Tx, error)
	DB() *sql.DB
}
