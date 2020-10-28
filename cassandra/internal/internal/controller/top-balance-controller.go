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
	response, _ := t.s.ProcessTopBalanceTransaction()
	printTopBalanceResponse(response)
}

func printTopBalanceResponse(r *model.TopBalanceResponse) {
	//fmt.Println(r)
	fmt.Sprintf("For each customer in C ranked in descending order of C_BALANCE:%v\n", ci)
}

func (t *topBalanceControllerImpl) Close() error {
	panic("implement me")
}
