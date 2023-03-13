package model

import "time"

type Payment struct {
	transaction     string
	requestID       string
	currency        string
	provider        string
	amount          int
	paymentDatetime time.Time
	bank            string
	deliveryCost    int
	goodsTotal      int
	customFee       int
}
