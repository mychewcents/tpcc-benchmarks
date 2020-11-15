package models

// OrderStatus stores the input for the OrderStatus transaction
type OrderStatus struct {
	WarehouseID int
	DistrictID  int
	CustomerID  int
}

// OrderStatusOutput stores the output of the order status transaction
type OrderStatusOutput struct {
}
