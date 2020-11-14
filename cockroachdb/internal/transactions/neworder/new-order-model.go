package neworder

// NewOrder defines the new order object
type NewOrder struct {
	WarehouseID       int
	DistrictID        int
	CustomerID        int
	IsOrderLocal      int
	UniqueItems       int
	TotalAmount       float64
	OrderTimestamp    string
	NewOrderLineItems map[int]*OrderLineItem
}

// OrderLineItem defines the order lines items for new orders
type OrderLineItem struct {
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

// Output stores the final output of the new order transaction
type Output struct {
	Customer         *CustomerInfo
	OrderID          int
	OrderTimestamp   string
	UniqueItems      int
	TotalOrderAmount float64
	DistrictTax      float64
	WarehouseTax     float64
	OrderLineItems   map[int]*OrderLineItem
}

// CustomerInfo to be used by the Output state
type CustomerInfo struct {
	WarehouseID int
	DistrictID  int
	CustomerID  int
	LastName    string
	Credit      string
	Discount    float64
}
