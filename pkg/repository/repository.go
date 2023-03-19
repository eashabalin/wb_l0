package repository

import (
	"github.com/jmoiron/sqlx"
	"wb_l0/pkg/model"
)

type Order interface {
	Create(order model.Order) (string, error)
}

type OrderCache interface {
	Set(order model.Order)
	Get(uid string) (*model.Order, bool)
	Delete(uid string) error
}

type Repository struct {
	Order
	OrderCache
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order:      NewOrderPostgres(db),
		OrderCache: NewOrderInMemCache(),
	}
}
