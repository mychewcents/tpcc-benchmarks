package model

import (
	"github.com/gocql/gocql"
	"time"
)

type OrderStatusRequest struct {
	CWId int
	CDId int
	CId  int
}

type OrderLineStatus struct {
	OlIId       int
	OlSupplyWId int
	OlQuantity  int
	OlAmount    float32
	OlDeliveryD time.Time
}

type OrderStatusResponse struct {
	CName    *Name
	CBalance float64

	OId        gocql.UUID
	OEntryD    time.Time
	OCarrierId int

	OrderLineStatusList []*OrderLineStatus
}
