package dbdatamodel

// Customer represents the Database state model
type Customer struct {
	WarehouseID int
	DistrictID  int
	CustomerID  int
	LastName    string
	Credit      string
	Discount    float64
}
