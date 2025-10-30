package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	Db *sql.DB
}

func New(dsn string) (*Storage, error) {
	op := "storage.postgres.New"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{Db: db}, nil
}

func (s *Storage) Close() error {
	return s.Db.Close()
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	op := "storage.postgres.SaveURL"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO urls (url, alias) VALUES ($1, $2)`

	_, err := s.Db.ExecContext(ctx, query, urlToSave, alias)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	op := "storage.postgres.GetURL"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT url FROM urls WHERE alias = $1`

	var url string

	err := s.Db.QueryRowContext(ctx, query, alias).Scan(&url)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	op := "storage.postgres.DeleteURL"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM urls WHERE alias = $1`

	_, err := s.Db.ExecContext(ctx, query, alias)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) checkExists(ctx context.Context, url, alias string) (urlExists bool, aliasExists bool, err error) {
	query := `
        SELECT 
            EXISTS(SELECT 1 FROM urls WHERE url = $1),
            EXISTS(SELECT 1 FROM urls WHERE alias = $2)
    `
	err = s.Db.QueryRowContext(ctx, query, url, alias).Scan(&urlExists, &aliasExists)
	return urlExists, aliasExists, err
}
