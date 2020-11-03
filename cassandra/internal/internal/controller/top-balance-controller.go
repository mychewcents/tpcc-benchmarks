package controller

import (
	"fmt"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
)

type TopBalanceController interface {
	handler.TransactionHandler
}

type topBalanceControllerImpl struct {
	s service.TopBalanceService
}

func NewTopBalanceController(cassandraSession *common.CassandraSession) TopBalanceController {
	return &topBalanceControllerImpl{
		s: service.NewTopBalanceService(cassandraSession),
	}
}

func (t *topBalanceControllerImpl) HandleTransaction(cmd []string) {
	t.s.ProcessTopBalanceTransaction()
	//printTopBalanceResponse(response)
}

func printTopBalanceResponse(r *model.TopBalanceResponse) {
	fmt.Println("*********************** Top Balance Transaction Output ***********************")
	fmt.Println("1. Customers ranked in descending order of their Balance:")
	fmt.Println()
	for _, c := range r.CustomerInfoList {
		fmt.Printf("\ta. Name of Customer: %+v\n", c.CName)
		fmt.Printf("\tb. Balance of customer’s outstanding payment: %.2f\n", c.CBalance)
		fmt.Printf("\tc. Warehouse name of customer: %+v\n", c.WName)
		fmt.Printf("\td. District name of customer: %+v\n\n", c.DName)
	}
	fmt.Println()
}

func (t *topBalanceControllerImpl) Close() error {
	panic("implement me")
}
