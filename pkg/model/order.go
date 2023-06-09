package model

type Order struct {
	UID               string `json:"order_uid"`
	TrackNumber       string `json:"track_number"`
	Entry             string `json:"entry"`
	Delivery          Delivery
	Payment           Payment
	Items             []Item
	Locale            string `json:"locale"`
	InternalSignature string `json:"internal_signature"`
	CustomerID        string `json:"customer_id"`
	DeliveryService   string `json:"delivery_service"`
	ShardKey          string `json:"shardkey"`
	SmID              int    `json:"sm_id"`
	DateCreated       string `json:"date_created"`
	OofShard          string `json:"oof_shard"`
}
