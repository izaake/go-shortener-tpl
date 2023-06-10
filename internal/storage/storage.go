package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx"
)

type Storage interface {
	Exec(ctx context.Context, stmt *sql.Stmt, args ...interface{}) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Ping() error
	BeginTx() (*sql.Tx, error)
}

type SQLStorage struct {
	DB *sql.DB
}

func (s *SQLStorage) Exec(ctx context.Context, stmt *sql.Stmt, args ...interface{}) error {
	_, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLStorage) Query(query string, args ...interface{}) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *SQLStorage) QueryRow(query string, args ...interface{}) *sql.Row {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...)
}

func (s *SQLStorage) Ping() error {
	return s.DB.Ping()
}

func (s *SQLStorage) BeginTx() (*sql.Tx, error) {
	return s.DB.Begin()
}
