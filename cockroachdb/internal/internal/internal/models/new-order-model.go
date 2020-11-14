package models

// NewOrder defines the new order object
type NewOrder struct {
	WarehouseID       int
	DistrictID        int
	CustomerID        int
	IsOrderLocal      int
	UniqueItems       int
	TotalAmount       float64
	OrderTimestamp    string
	NewOrderLineItems map[int]*NewOrderOrderLineItem
}

// NewOrderOrderLineItem defines the order lines items for new orders
type NewOrderOrderLineItem struct {
	Name                string
	SupplierWarehouseID int
	Quantity            int
	IsRemote            int
	StartStock          int
	FinalStock          int
	Data                string
	Price               float64
	CurrYTD             float64
	CurrOrderCnt        int
	Amount              float64
}

// NewOrderOutput stores the final output of the new order transaction
type NewOrderOutput struct {
	Customer         *NewOrderCustomerInfo
	OrderID          int
	OrderTimestamp   string
	UniqueItems      int
	TotalOrderAmount float64
	DistrictTax      float64
	WarehouseTax     float64
	OrderLineItems   map[int]*NewOrderOrderLineItem
}

// NewOrderCustomerInfo to be used by the Output state
type NewOrderCustomerInfo struct {
	WarehouseID int
	DistrictID  int
	CustomerID  int
	LastName    string
	Credit      string
	Discount    float64
}
