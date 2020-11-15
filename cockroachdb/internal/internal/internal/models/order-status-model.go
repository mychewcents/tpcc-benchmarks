package models

import "github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"

// OrderStatus stores the input for the OrderStatus transaction
type OrderStatus struct {
	WarehouseID int
	DistrictID  int
	CustomerID  int
}

// OrderStatusOutput stores the output of the order status transaction
type OrderStatusOutput struct {
	Order      *dbdatamodel.Order
	Customer   *dbdatamodel.Customer
	OrderLines map[int]*dbdatamodel.OrderLineItem
}
