package repository

import "errors"

var (
	ErrURLNotFound        = errors.New("url not found")
	ErrAliasAlreadyExists = errors.New("alias already exists")
	ErrAliasNotFound      = errors.New("alias not found")
	ErrAliasNotUnique     = errors.New("alias not unique")
)

type URLRepository interface {
	SaveURL(urlString, alias string) error
	GetURL(alias string) (string, error)
	GetAlias(urlString string) (string, error)
	DeleteURL(alias string) error
}
