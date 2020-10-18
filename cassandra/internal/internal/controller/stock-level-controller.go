package controller

import (
	"bufio"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
)

type StockLevelController interface {
	handler.TransactionHandler
}

type stockLevelControllerImpl struct {
	s service.StockLevelService
	r *bufio.Reader
}

func NewStockLevelController(cassandraSession *common.CassandraSession, reader *bufio.Reader) StockLevelController {
	return &stockLevelControllerImpl{
		s: service.NewStockLevelService(cassandraSession),
		r: reader,
	}
}

func (s *stockLevelControllerImpl) HandleTransaction(i []string) {

}

func (s *stockLevelControllerImpl) Close() error {
	panic("implement me")
}
