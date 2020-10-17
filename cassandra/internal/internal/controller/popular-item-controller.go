package controller

import (
	"bufio"
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
)

type PopularItemController interface {
	handler.TransactionHandler
}

type popularItemControllerImpl struct {
	s service.PopularItemService
	r *bufio.Reader
}

func NewPopularItemController(cluster *gocql.ClusterConfig, reader *bufio.Reader) PopularItemController {
	return &popularItemControllerImpl{
		s: service.NewPopularItemService(cluster),
		r: reader,
	}
}

func (p *popularItemControllerImpl) HandleTransaction(i []string) {
}

func (p *popularItemControllerImpl) Close() error {
	panic("implement me")
}
