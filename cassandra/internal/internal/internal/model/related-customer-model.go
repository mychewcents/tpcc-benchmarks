package model

type RelatedCustomerRequest struct {
	CWId int
	CDId int
	CId  int
}

type RelatedCustomerResponse struct {
	CWId int
	CDId int
	CId  int

	RelatedCIds []int
}
