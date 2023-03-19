package repository

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	items    map[string]*list.Element
	queue    *list.List
	capacity int
}

type CacheItem struct {
	Key       string
	Value     any
	CreatedAt time.Time
}

func NewCache(capacity int) *Cache {
	items := make(map[string]*list.Element)
	cache := Cache{items: items, queue: list.New(), capacity: capacity}
	return &cache
}

func (c *Cache) Set(key string, value any) {
	c.Lock()
	defer c.Unlock()

	if element, exists := c.items[key]; exists {
		c.queue.MoveToFront(element)
		element.Value.(*CacheItem).Value = value
		return
	}

	if c.queue.Len() == c.capacity {
		c.purge()
	}

	item := &CacheItem{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
	}

	element := c.queue.PushFront(item)
	c.items[key] = element
}

func (c *Cache) purge() {
	if element := c.queue.Back(); element != nil {
		item := c.queue.Remove(element).(*CacheItem)
		delete(c.items, item.Key)
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.RLock()
	defer c.RUnlock()

	element, exists := c.items[key]
	if !exists {
		return nil, false
	}
	c.queue.MoveToFront(element)
	return element.Value.(*CacheItem).Value, true
}

func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	element, exists := c.items[key]
	if !exists {
		return errors.New("item with such key was not found in cache")
	}

	delete(c.items, key)
	c.queue.Remove(element)

	return nil
}
