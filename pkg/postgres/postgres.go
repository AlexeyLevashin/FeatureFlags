package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к бд: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	return db, nil
}
