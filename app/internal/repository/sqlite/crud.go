package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/svetlana-mel/url-shortener/internal/repository"
)

// gcc -c C:\Users\79267\go\pkg\mod\github.com\mattn\go-sqlite3@v1.14.22\sqlite3-binding.c

func (s *storage) SaveURL(urlString, alias string) error {
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
			return fmt.Errorf("%s: %w", op, repository.ErrAliasAlreadyExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *storage) GetURL(alias string) (string, error) {
	const op = "repository.sqlite.crud.GetURL"

	row := s.db.QueryRow(`
		SELECT url
		FROM url
		WHERE alias=?
	`, alias)

	var urlResult string

	err := row.Scan(&urlResult)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, repository.ErrURLNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return urlResult, nil
}

func (s *storage) GetAlias(urlString string) (string, error) {
	const op = "repository.sqlite.crud.GetAlias"

	row := s.db.QueryRow(`
		SELECT alias
		FROM url
		WHERE url=?
	`, urlString)

	var alias string

	err := row.Scan(&alias)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", repository.ErrAliasNotFound
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return alias, nil
}

func (s *storage) DeleteURL(alias string) error {
	const op = "repository.sqlite.crud.DeleteURL"

	stmt, err := s.db.Prepare(`
		DELETE FROM url
		WHERE alias=?
	`)
	if err != nil {
		return fmt.Errorf("%s prepare stmt error: %w", op, err)
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return repository.ErrURLNotFound
	}

	return nil
}
