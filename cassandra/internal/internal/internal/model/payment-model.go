package model

import "time"

type PaymentRequest struct {
	CWId    int
	CDId    int
	CId     int
	Payment float64
}

type PaymentResponse struct {
	CWId int
	CDId int
	CId  int

	CName      *Name
	CAddress   *Address
	CPhone     string
	CSince     time.Time
	CCredit    string
	CCreditLim float64
	CDiscount  float32
	CBalance   float64

	WAddress *Address
	DAddress *Address

	Payment float64
}
