package model

type CustomerIdentifier struct {
	CWId int
	CDId int
	CId  int
}

type RelatedCustomerRequest struct {
	CustomerIdentifier *CustomerIdentifier
}

type RelatedCustomerResponse struct {
	CustomerIdentifier         *CustomerIdentifier
	RelatedCustomerIdentifiers []*CustomerIdentifier
}
