package controller

import (
	"fmt"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
	"strconv"
)

type RelatedCustomerController interface {
	handler.TransactionHandler
}

type relatedCustomerControllerImpl struct {
	s service.RelatedCustomerService
}

func NewRelatedCustomerController(cassandraSession *common.CassandraSession) RelatedCustomerController {
	return &relatedCustomerControllerImpl{
		s: service.NewRelatedCustomerService(cassandraSession),
	}
}

func (r *relatedCustomerControllerImpl) HandleTransaction(cmd []string) {
	request := makeRelatedCustomerRequest(cmd)
	r.s.ProcessRelatedCustomerTransaction(request)
	//printRelatedCustomerResponse(response)
}

func makeRelatedCustomerRequest(cmd []string) *model.RelatedCustomerRequest {
	cWId, _ := strconv.Atoi(cmd[1])
	cDId, _ := strconv.Atoi(cmd[2])
	cId, _ := strconv.Atoi(cmd[3])

	return &model.RelatedCustomerRequest{CustomerIdentifier: &model.CustomerIdentifier{
		CWId: cWId,
		CDId: cDId,
		CId:  cId,
	}}
}

func printRelatedCustomerResponse(r *model.RelatedCustomerResponse) {
	if r == nil {
		return
	}
	fmt.Println("*********************** Related Customer Transaction Output ***********************")
	fmt.Printf("1. Customer Identifier: %+v\n", r.CustomerIdentifier)
	fmt.Println("2. Related Customers:")
	for _, c := range r.RelatedCustomerIdentifiers {
		fmt.Printf("\tRelated Customer Identifier: %+v\n", c)
	}
	fmt.Println()
}

func (r *relatedCustomerControllerImpl) Close() error {
	panic("implement me")
}
