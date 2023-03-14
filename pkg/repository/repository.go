package repository

import (
	"github.com/jmoiron/sqlx"
	"wb_l0/pkg/model"
)

type Order interface {
	Create(order model.Order) (string, error)
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order: NewOrderPostgres(db),
	}
}
