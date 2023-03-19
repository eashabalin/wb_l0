package service

import (
	"wb_l0/pkg/model"
	"wb_l0/pkg/repository"
)

type OrderService struct {
	postgres repository.OrderDB
	cache    repository.OrderCache
}

func NewOrderService(repo repository.OrderDB, cache repository.OrderCache) *OrderService {
	return &OrderService{postgres: repo, cache: cache}
}

func (s *OrderService) Create(order model.Order) (string, error) {
	uid, err := s.postgres.Create(order)
	if err != nil {
		return "", err
	}
	s.cache.Set(order)
	return uid, nil
}

func (s *OrderService) Get(uid string) (*model.Order, bool) {
	order, exists := s.cache.Get(uid)
	if exists {
		return order, true
	}
	order, exists = s.postgres.Get(uid)
	if exists {
		return order, true
	}
	return nil, false
}
