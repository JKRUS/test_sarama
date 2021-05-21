package repository

import (
	"github.com/jmoiron/sqlx"
)

type Data interface {
}

type Repository struct {
	Data
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Data: NewMsgPostgres(db),
	}
}
