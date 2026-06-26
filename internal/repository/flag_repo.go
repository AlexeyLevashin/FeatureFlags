package repository

import (
	"github.com/jmoiron/sqlx"
)

type FlagRepo struct {
	db *sqlx.DB
}

func NewFlagRepository(db *sqlx.DB) *FlagRepo {
	return &FlagRepo{db: db}
}

//func Create
