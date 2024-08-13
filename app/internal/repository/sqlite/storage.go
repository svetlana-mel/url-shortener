package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/svetlana-mel/url-shortener/internal/repository"
)

var _ repository.URLRepository = (*storage)(nil)

type storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*storage, error) {
	const op = "storage.sqlite/NewStorage"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	initStmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL,
		create_time TEXT NOT NULL,
		update_time TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = initStmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &storage{db}, nil
}
