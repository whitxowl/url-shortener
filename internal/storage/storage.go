package storage

import (
	"errors"
	"fmt"
)

type ErrNoUrlWithAlias struct {
	Alias string
}

func (e *ErrNoUrlWithAlias) Error() string {
	return fmt.Sprintf("no URL found with alias: %s", e.Alias)
}

func (e *ErrNoUrlWithAlias) Is(target error) bool {
	var errNoUrlWithAlias *ErrNoUrlWithAlias
	ok := errors.As(target, &errNoUrlWithAlias)
	return ok
}

var (
	ErrURLNotFound      = errors.New("URL not found")
	ErrURLAlreadyExists = errors.New("URL already exists")
)
