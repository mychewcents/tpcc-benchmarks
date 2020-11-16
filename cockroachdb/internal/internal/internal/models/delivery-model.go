package models

// Delivery stores the input for the Delivery transaction
type Delivery struct {
	WarehouseID int
	CarrierID   int
}

// DeliveryOutput stores the output of the Delivery transaction
type DeliveryOutput struct {
}
