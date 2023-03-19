package repository

import (
	"github.com/jmoiron/sqlx"
	"wb_l0/pkg/model"
)

type OrderDB interface {
	Create(order model.Order) (string, error)
	Get(uid string) (*model.Order, bool)
}

type OrderCache interface {
	Set(order model.Order)
	Get(uid string) (*model.Order, bool)
	Delete(uid string) error
}

type Repository struct {
	OrderDB
	OrderCache
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		OrderDB:    NewOrderPostgres(db),
		OrderCache: NewOrderInMemCache(),
	}
}
