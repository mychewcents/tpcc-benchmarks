package service

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
)

type PaymentService interface {
	ProcessPaymentTransaction(request *model.PaymentRequest) (*model.PaymentResponse, error)
	io.Closer
}

type paymentServiceImpl struct {
}

func NewPaymentService(cluster *gocql.ClusterConfig) PaymentService {
	return &paymentServiceImpl{}
}

func (p *paymentServiceImpl) ProcessPaymentTransaction(request *model.PaymentRequest) (*model.PaymentResponse, error) {
	panic("implement me")
}

func (p *paymentServiceImpl) Close() error {
	panic("implement me")
}
