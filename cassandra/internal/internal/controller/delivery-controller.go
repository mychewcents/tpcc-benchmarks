package controller

import (
	"bufio"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
)

type DeliveryController interface {
	handler.TransactionHandler
}

type deliveryControllerImpl struct {
	s service.DeliveryService
	r *bufio.Reader
}

func NewDeliveryTransactionController(cassandraSession *common.CassandraSession, reader *bufio.Reader) DeliveryController {
	return &deliveryControllerImpl{
		s: service.NewDeliveryService(cassandraSession),
		r: reader,
	}
}

func (d *deliveryControllerImpl) HandleTransaction(i []string) {

}

func (d *deliveryControllerImpl) Close() error {
	panic("implement me")
}
