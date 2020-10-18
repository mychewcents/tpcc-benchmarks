package service

import (
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
)

type StockLevelService interface {
	ProcessStockLevelTransaction(request *model.StockLevelRequest) error
	io.Closer
}

type stockLevelServiceImpl struct {
}

func NewStockLevelService(cassandraSession *common.CassandraSession) StockLevelService {
	return &stockLevelServiceImpl{}
}

func (s *stockLevelServiceImpl) ProcessStockLevelTransaction(request *model.StockLevelRequest) error {
	panic("implement me")
}

func (s *stockLevelServiceImpl) Close() error {
	panic("implement me")
}
