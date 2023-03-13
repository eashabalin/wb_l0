package model

import "github.com/gogo/protobuf/types"

type Order struct {
	orderUID          string
	trackNumber       string
	entry             string
	delivery          Delivery
	payment           Payment
	locale            string
	internalSignature string
	customerID        string
	deliveryService   string
	shardKey          string
	smID              int
	dateCreated       types.Timestamp
	oofShard          string
}
