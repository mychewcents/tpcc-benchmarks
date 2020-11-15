package dbdatamodel

// Order denotes the Order object from the database
type Order struct {
	ID                int
	WarehouseID       int
	DistrictID        int
	CustomerID        int
	Timestamp         string
	CarrierID         int
	DeliveryTimestamp string
}
