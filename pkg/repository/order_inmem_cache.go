package repository

import "wb_l0/pkg/model"

type OrderInMemCache struct {
	cache *Cache
}

func NewOrderInMemCache() *OrderInMemCache {
	cache := NewCache(3)
	return &OrderInMemCache{cache: cache}
}

func (c *OrderInMemCache) Set(order model.Order) {
	c.cache.Set(order.UID, order)
}

func (c *OrderInMemCache) Get(uid string) (*model.Order, bool) {
	value, exists := c.cache.Get(uid)
	if !exists {
		return nil, false
	}
	order := value.(model.Order)
	return &order, exists
}

func (c *OrderInMemCache) Delete(uid string) error {
	return c.cache.Delete(uid)
}
