package nats_streaming

import (
	"fmt"
	"github.com/nats-io/stan.go"
)

func (h *Handler) CreateOrder(m *stan.Msg) {
	fmt.Printf("received: %s\n", m)
}
