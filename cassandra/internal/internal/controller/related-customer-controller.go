package controller

import (
	"bufio"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
)

type RelatedCustomerController interface {
	handler.TransactionHandler
}

type relatedCustomerControllerImpl struct {
	s service.RelatedCustomerService
	r *bufio.Reader
}

func NewRelatedCustomerController(cassandraSession *common.CassandraSession, reader *bufio.Reader) RelatedCustomerController {
	return &relatedCustomerControllerImpl{
		s: service.NewRelatedCustomerService(cassandraSession),
		r: reader,
	}
}

func (r *relatedCustomerControllerImpl) HandleTransaction(i []string) {
}

func (r *relatedCustomerControllerImpl) Close() error {
	panic("implement me")
}
