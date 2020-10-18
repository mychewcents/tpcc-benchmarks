package service

import (
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
)

type PaymentService interface {
	ProcessPaymentTransaction(request *model.PaymentRequest) (*model.PaymentResponse, error)
	io.Closer
}

type paymentServiceImpl struct {
}

func NewPaymentService(cassandraSession *common.CassandraSession) PaymentService {
	return &paymentServiceImpl{}
}

func (p *paymentServiceImpl) ProcessPaymentTransaction(request *model.PaymentRequest) (*model.PaymentResponse, error) {
	panic("implement me")
}

func (p *paymentServiceImpl) Close() error {
	panic("implement me")
}
