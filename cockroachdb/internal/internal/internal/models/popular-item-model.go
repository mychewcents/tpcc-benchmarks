package models

import "github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"

// PopularItem stores the popular item input
type PopularItem struct {
	WarehouseID int
	DistrictID  int
	LastNOrders int
}

// PopularItemOrderDetails stores the per order details
type PopularItemOrderDetails struct {
	ID        int
	Timestamp string
	Customer  *dbdatamodel.Customer
	Items     []*dbdatamodel.Item
}

// PopularItemOutput stores the popular item output
type PopularItemOutput struct {
	WarehouseID     int
	DistrictID      int
	StartOrderID    int
	LastOrderID     int
	Orders          map[int]*PopularItemOrderDetails
	ItemPercentages map[int]float64
}
