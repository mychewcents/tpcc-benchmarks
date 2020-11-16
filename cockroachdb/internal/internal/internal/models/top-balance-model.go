package models

// TopBalance stores the input for the top balance transaction
type TopBalance struct {
}

// TopBalanceOutputRow stores each row for the top balance output
type TopBalanceOutputRow struct {
	WarehouseID, DistrictID, CustomerID                          int
	Balance                                                      float64
	FirstName, MiddleName, LastName, DistrictName, WarehouseName string
}

// TopBalanceOutput stores the output for the top balance transaction
type TopBalanceOutput struct {
	Rows []*TopBalanceOutputRow
}
