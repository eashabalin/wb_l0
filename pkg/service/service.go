package service

import (
	"wb_l0/pkg/model"
	"wb_l0/pkg/repository"
)

type Order interface {
	Create(order model.Order) (string, error)
	GetByUID(uid string) (*model.Order, bool)
	GetAllFromDB() []model.Order
}

type Service struct {
	Order
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repos.OrderDB, repos.OrderCache),
	}
}
