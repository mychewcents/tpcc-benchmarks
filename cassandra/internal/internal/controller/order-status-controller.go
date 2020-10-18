package controller

import (
	"bufio"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
)

type OrderStatusController interface {
	handler.TransactionHandler
}

type orderStatusControllerImpl struct {
	s service.OrderStatusService
	r *bufio.Reader
}

func NewOrderStatusTransactionController(cassandraSession *common.CassandraSession, reader *bufio.Reader) OrderStatusController {
	return &orderStatusControllerImpl{
		s: service.NewOrderStatusService(cassandraSession),
		r: reader,
	}
}

func (n *orderStatusControllerImpl) HandleTransaction(i []string) {

}

func (n *orderStatusControllerImpl) Close() error {
	panic("implement me")
}
