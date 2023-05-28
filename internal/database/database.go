package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB(dbConnection *string) (*sql.DB, error) {
	db, err := sql.Open("pgx", *dbConnection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	query := `
create table if not exists public.urls(
   	 user_id uuid not null,
   	 short_url text not null,
   	 original_url text not null,
   	 UNIQUE (user_id, short_url)
);`
	db.Exec(query)

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(time.Minute * 10)

	return db, nil
}
