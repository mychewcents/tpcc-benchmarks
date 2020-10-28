package model

type RelatedCustomerRequest struct {
	CWId int
	CDId int
	CId  int
}

type CustomerIdentifier struct {
	CWId int
	CDId int
	CId  int
}

type RelatedCustomerResponse struct {
	CustomerIdentifier         *CustomerIdentifier
	RelatedCustomerIdentifiers []*CustomerIdentifier
}
