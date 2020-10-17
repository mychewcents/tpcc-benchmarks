package service

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
)

type TopBalanceService interface {
	ProcessTopBalanceTransaction(N int) (*model.TopBalanceResponse, error)
	io.Closer
}

type topBalanceServiceImpl struct {
}

func NewTopBalanceService(cluster *gocql.ClusterConfig) TopBalanceService {
	return &topBalanceServiceImpl{}
}

func (t *topBalanceServiceImpl) ProcessTopBalanceTransaction(N int) (*model.TopBalanceResponse, error) {
	panic("implement me")
}

func (t *topBalanceServiceImpl) Close() error {
	panic("implement me")
}
