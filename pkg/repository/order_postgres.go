package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"wb_l0/pkg/model"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{
		db: db,
	}
}

func (r *OrderPostgres) Create(order model.Order) (string, error) {
	fmt.Printf("%+v\n", order)
	return "", nil
}
