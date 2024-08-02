package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/svetlana-mel/url-shortener/internal/storage"
)

// gcc -c C:\Users\79267\go\pkg\mod\github.com\mattn\go-sqlite3@v1.14.22\sqlite3-binding.c

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
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

	return &Storage{db}, nil
}

func (s *Storage) SaveURL(urlString, alias string) error {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare(`
		INSERT INTO url
		(url, alias, create_time, update_time) 
		VALUES(?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	timestamp := time.Now()

	_, err = stmt.Exec(urlString, alias, timestamp, timestamp)

	if err != nil {
		sqliteError, ok := err.(sqlite3.Error)
		if ok && sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrAliasExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	row := s.db.QueryRow(`
		SELECT url
		FROM url
		WHERE alias=?
	`, alias)

	var urlResult string

	err := row.Scan(&urlResult)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return urlResult, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare(`
		DELETE FROM url
		WHERE alias=?
	`)
	if err != nil {
		return fmt.Errorf("%s prepare stmt error: %w", op, err)
	}

	_, err = stmt.Exec(alias)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
