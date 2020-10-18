package controller

import (
	"fmt"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
)

type PaymentController interface {
	handler.TransactionHandler
}

type paymentControllerImpl struct {
	s service.PaymentService
}

func NewPaymentController(cassandraSession *common.CassandraSession) PaymentController {
	return &paymentControllerImpl{
		s: service.NewPaymentService(cassandraSession),
	}
}

func (p *paymentControllerImpl) HandleTransaction(cmd []string) {
	request := makePaymentRequest(cmd)
	response, _ := p.s.ProcessPaymentTransaction(request)
	printPaymentResponse(response)
}

func makePaymentRequest(cmd []string) *model.PaymentRequest {
	panic("implement me")
}

func printPaymentResponse(r *model.PaymentResponse) {
	fmt.Println(r)
}

func (p *paymentControllerImpl) Close() error {
	panic("implement me")
}
