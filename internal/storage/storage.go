package storage

import (
	"database/sql"

	_ "github.com/jackc/pgx"
)

type Storage interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Ping() error
}

type SQLStorage struct {
	DB *sql.DB
}

func (s *SQLStorage) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.Exec(query, args)
}

func (s *SQLStorage) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.DB.Query(query, args)
}

func (s *SQLStorage) Ping() error {
	return s.DB.Ping()
}
