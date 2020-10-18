package controller

import (
	"bufio"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
)

type TopBalanceController interface {
	handler.TransactionHandler
}

type topBalanceControllerImpl struct {
	s service.TopBalanceService
	r *bufio.Reader
}

func NewTopBalanceController(cassandraSession *common.CassandraSession, reader *bufio.Reader) TopBalanceController {
	return &topBalanceControllerImpl{
		s: service.NewTopBalanceService(cassandraSession),
		r: reader,
	}
}

func (t *topBalanceControllerImpl) HandleTransaction(i []string) {
}

func (t *topBalanceControllerImpl) Close() error {
	panic("implement me")
}
