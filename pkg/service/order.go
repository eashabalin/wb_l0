package service

import (
	"wb_l0/pkg/model"
	"wb_l0/pkg/repository"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(order model.Order) (string, error) {
	return s.repo.Create(order)
}
