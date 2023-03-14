package service

import (
	"wb_l0/pkg/model"
	"wb_l0/pkg/repository"
)

type Order interface {
	Create(order model.Order) (string, error)
}

type Service struct {
	Order
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repos.Order),
	}
}
