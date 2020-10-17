package neworder

type NewOrder struct {
	// INPUTS
	CustomerID, DistrictID, WarehouseID, NumItems int
	ItemIDs, SupplierWarehouseIDs, ItemQuantities []int64

	// OUTPUTS
	lastName, creditStatus, orderTimestamp string
	custDiscount, totalAmount, warehouseTax, districtTax float64
	orderID int
	itemNames []string
	itemAmount []float64
	itemStock []int
}

func (f *NewOrder) ProcessTransaction() string {
	fmt.Printf("Hello Akarsh!")
	fmt.Printf("%d", f.CustomerID)

	return f.printOutputState()
}

func (f *NewOrder) printOutputState() string {
	return "constant for now"
}
