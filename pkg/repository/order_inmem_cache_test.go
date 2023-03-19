package repository

import (
	"fmt"
	"testing"
	"wb_l0/pkg/model"
)

func TestSet(t *testing.T) {
	orderCache := NewOrderInMemCache()
	orderCache.Set(model.Order{})
}

func TestGet(t *testing.T) {
	orderCache := NewOrderInMemCache()
	orderCache.Set(model.Order{UID: "fgdsfas"})
	_, exists := orderCache.Get("fgdsfas")
	if !exists {
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	orderCache := NewOrderInMemCache()
	orderCache.Set(model.Order{UID: "fgdsfas"})
	err := orderCache.Delete("fgdsfas")
	if err != nil {
		t.Fail()
		return
	}
	_, exists := orderCache.Get("fgdsfas")
	if exists {
		t.Fail()
	}
}

func Test(t *testing.T) {
	orderCache := NewOrderInMemCache()

	orderCache.Set(model.Order{UID: "1"})
	fmt.Println(orderCache.cache.queue.Len())
	orderCache.Set(model.Order{UID: "2"})
	fmt.Println(orderCache.cache.queue.Len())
	orderCache.Set(model.Order{UID: "3"})
	fmt.Println(orderCache.cache.queue.Len())
	orderCache.Set(model.Order{UID: "4"})
	fmt.Println(orderCache.cache.queue.Len())
	_, exists := orderCache.Get("1")
	if exists {
		t.Fail()
	}
}
