package dbdatamodel

// Customer represents the Database state model
type Customer struct {
	WarehouseID, DistrictID, CustomerID int
	FirstName, MiddleName, LastName     string
	Addr                                Address
	Credit                              string
	Discount                            float64
	Balance                             float64
	Phone                               string
	Since                               string
	CreditLimit                         float64
}
