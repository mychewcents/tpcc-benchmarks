package models

import "github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"

// NewOrder defines the new order object
type NewOrder struct {
	WarehouseID       int
	DistrictID        int
	CustomerID        int
	IsOrderLocal      int
	UniqueItems       int
	TotalAmount       float64
	OrderTimestamp    string
	NewOrderLineItems map[int]*NewOrderOrderLineItem
}

// NewOrderOrderLineItem stores the input for the NewOrder Order Lines
type NewOrderOrderLineItem struct {
	Name                string
	SupplierWarehouseID int
	Quantity            int
	IsRemote            int
}

// NewOrderOutput stores the final output of the new order transaction
type NewOrderOutput struct {
	Customer         *dbdatamodel.Customer
	OrderID          int
	OrderTimestamp   string
	UniqueItems      int
	TotalOrderAmount float64
	DistrictTax      float64
	WarehouseTax     float64
	OrderLineItems   map[int]*dbdatamodel.OrderLineItem
}
