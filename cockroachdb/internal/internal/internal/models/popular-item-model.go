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
	Order         *dbdatamodel.Order
	MaxOLQuantity int
	Customer      *dbdatamodel.Customer
	Items         map[int]*dbdatamodel.Item
}

// PopularItemOccuranceAndPercentage stores the occurrance count and the percentage for each item
type PopularItemOccuranceAndPercentage struct {
	Name       string
	Occurances int
	Percentage float64
}

// PopularItemOutput stores the popular item output
type PopularItemOutput struct {
	WarehouseID    int
	DistrictID     int
	StartOrderID   int
	LastOrderID    int
	Orders         map[int]*PopularItemOrderDetails
	ItemOccurances map[int]*PopularItemOccuranceAndPercentage
}
