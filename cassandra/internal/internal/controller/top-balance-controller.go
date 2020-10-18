package controller

import (
	"bufio"
	"github.com/gocql/gocql"
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

func NewTopBalanceController(cluster *gocql.ClusterConfig, reader *bufio.Reader) TopBalanceController {
	return &topBalanceControllerImpl{
		s: service.NewTopBalanceService(cluster),
		r: reader,
	}
}

func (t *topBalanceControllerImpl) HandleTransaction(i []string) {
}

func (t *topBalanceControllerImpl) Close() error {
	panic("implement me")
}
