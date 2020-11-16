package models

// StockLevel stores the transaction details for the Stock Level
type StockLevel struct {
	WarehouseID int
	DistrictID  int
	Threshold   int
	LastNOrders int
}

// StockLevelOutput stores the output for the Stock Level transaction
type StockLevelOutput struct {
	TotalItems   int
	StartOrderID int
	LastOrderID  int
}
