package dbdatamodel

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
