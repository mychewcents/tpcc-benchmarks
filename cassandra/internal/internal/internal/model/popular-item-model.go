package model

import (
	"github.com/gocql/gocql"
	"time"
)

type PopularItemRequest struct {
	WId            int
	DId            int
	NoOfLastOrders int
}

type PopularItemResponse struct {
	WId int
	DId int

	NoOfLastOrders int

	OrderItemInfoList []*OrderItemInfo

	PopularItemStatList []*PopularItemStat
}

type OrderItemInfo struct {
	OId     gocql.UUID
	OEntryD time.Time

	CName *Name

	PopularItemInfoList []*PopularItemInfo
}

type PopularItemInfo struct {
	IName      string
	OlQuantity int
}

type PopularItemStat struct {
	IName           string
	OrderPercentage float32
}
