package models

// RelatedCustomer stores the input of the related customer transaction
type RelatedCustomer struct {
	WarehouseID int
	DistrictID  int
	CustomerID  int
}

// RelatedCustomerOutput stores the output of the Related Customer transaction
type RelatedCustomerOutput struct {
	Customers map[int]map[int][]int
}
