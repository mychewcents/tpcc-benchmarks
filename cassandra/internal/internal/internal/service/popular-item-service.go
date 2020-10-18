package service

import (
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
)

type PopularItemService interface {
	ProcessPopularItemService(request *model.PopularItemRequest) (*model.PopularItemResponse, error)
	io.Closer
}

type popularItemServiceImpl struct {
}

func NewPopularItemService(cassandraSession *common.CassandraSession) PopularItemService {
	return &popularItemServiceImpl{}
}

func (p *popularItemServiceImpl) ProcessPopularItemService(request *model.PopularItemRequest) (*model.PopularItemResponse, error) {
	panic("implement me")
}

func (p *popularItemServiceImpl) Close() error {
	panic("implement me")
}
