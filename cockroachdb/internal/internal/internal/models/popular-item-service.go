package models

// PopularItem stores the popular item input
type PopularItem struct {
	WarehouseID int
	DistrictID  int
	LastNOrders int
}

// PopularItemOutput stores the popular item output
type PopularItemOutput struct {
}
