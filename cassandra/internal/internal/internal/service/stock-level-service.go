package service

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
)

type StockLevelService interface {
	processStockLevelTransaction(request *model.StockLevelRequest) error
	io.Closer
}

type stockLevelServiceImpl struct {
}

func NewStockLevelService(cluster *gocql.ClusterConfig) StockLevelService {
	return &stockLevelServiceImpl{}
}

func (s *stockLevelServiceImpl) processStockLevelTransaction(request *model.StockLevelRequest) error {
	panic("implement me")
}

func (s *stockLevelServiceImpl) Close() error {
	panic("implement me")
}
