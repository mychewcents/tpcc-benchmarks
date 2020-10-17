package model

import (
	"github.com/gocql/gocql"
	"time"
)

type NewOrderLine struct {
	OlIId       int
	OlSupplyWId int
	OlQuantity  int
}

type NewOrderRequest struct {
	WId              int
	DId              int
	CId              int
	NewOrderLineList []*NewOrderLine
}

type NewOrderLineInfo struct {
	IId         int
	IName       string
	SupplierWId int
	Quantity    int
	OlAmount    int
	SQuantity   int
}

type NewOrderResponse struct {
	WId                  int
	DId                  int
	CId                  int
	CCredit              string
	CDiscount            float32
	WTax                 float32
	DTax                 float32
	OId                  gocql.UUID
	OEntryD              time.Time
	NoOfItems            int
	TotalAmount          float64
	NewOrderLineInfoList []*NewOrderLineInfo
}
