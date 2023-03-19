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
	orders := repo.GetAll()
	for _, order := range orders {
		cache.Set(order)
	}
	service := OrderService{postgres: repo, cache: cache}
	return &service
}

func (s *OrderService) Create(order model.Order) (string, error) {
	uid, err := s.postgres.Create(order)
	if err != nil {
		return "", err
	}
	s.cache.Set(order)
	return uid, nil
}

func (s *OrderService) GetByUID(uid string) (*model.Order, bool) {
	order, exists := s.cache.Get(uid)
	if exists {
		return order, true
	}
	order, exists = s.postgres.GetByUID(uid)
	if exists {
		return order, true
	}
	return nil, false
}

func (s *OrderService) GetAllFromDB() []model.Order {
	return s.postgres.GetAll()
}
