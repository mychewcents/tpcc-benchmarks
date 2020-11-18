package models

// DatabaseState stores the output of the database state
type DatabaseState struct {
	TotalYTDWarehouse     float64
	TotalYTDDistrict      float64
	SumOrderIDs           int
	CBalance              float64
	CYTDPayment           float64
	CPaymentCount         int
	CDeliveryCount        int
	MaxOrderID            int
	TotalOrderLineCount   int
	TotalOrderAmount      float64
	TotalQuantity         int
	TotalStock            int
	TotalYTDStock         float64
	TotalOrderCount       int
	TotalRemoteOrderCount int
}
