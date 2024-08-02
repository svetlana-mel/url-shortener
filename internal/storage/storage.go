package storage

// интерфейсы к бд лежат рядом с контроллерами

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrAliasExists = errors.New("alias already exists")
)
