package controller

import (
	"bufio"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
)

type PaymentController interface {
	handler.TransactionHandler
}

type paymentControllerImpl struct {
	s service.PaymentService
	r *bufio.Reader
}

func NewPaymentController(cassandraSession *common.CassandraSession, reader *bufio.Reader) PaymentController {
	return &paymentControllerImpl{
		s: service.NewPaymentService(cassandraSession),
		r: reader,
	}
}

func (p *paymentControllerImpl) HandleTransaction(i []string) {

}

func (p *paymentControllerImpl) Close() error {
	panic("implement me")
}
