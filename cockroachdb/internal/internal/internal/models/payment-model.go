package models

import (
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"
)

// Payment stores the input parameters for the payment transaction
type Payment struct {
	WarehouseID int
	DistrictID  int
	CustomerID  int
	Amount      float64
}

// PaymentOutput stores the output of the Payment transaction
type PaymentOutput struct {
	Customer      *dbdatamodel.Customer
	WarehouseAddr *dbdatamodel.Address
	DistrictAddr  *dbdatamodel.Address
	PaidAmount    float64
}
