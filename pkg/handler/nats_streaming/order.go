package nats_streaming

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"wb_l0/pkg/model"
)

func (h *Handler) CreateOrder(m *stan.Msg) {
	var order model.Order
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		fmt.Println(err)
	}
	_, err = h.services.Order.Create(order)
	if err != nil {
		fmt.Println(err)
	}
}
