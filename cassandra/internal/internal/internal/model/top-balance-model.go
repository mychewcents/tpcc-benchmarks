package model

type TopBalanceCustomerInfo struct {
	CName    Name
	CBalance float64
	WName    string
	DName    string
}

type TopBalanceResponse struct {
	CustomerInfoList [10]*TopBalanceCustomerInfo
}
