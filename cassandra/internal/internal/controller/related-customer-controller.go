package controller

import (
	"bufio"
	"github.com/gocql/gocql"
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

func NewRelatedCustomerController(cluster *gocql.ClusterConfig, reader *bufio.Reader) RelatedCustomerController {
	return &relatedCustomerControllerImpl{
		s: service.NewRelatedCustomerService(cluster),
		r: reader,
	}
}

func (r *relatedCustomerControllerImpl) HandleTransaction(i []string) {
}

func (r *relatedCustomerControllerImpl) Close() error {
	panic("implement me")
}
