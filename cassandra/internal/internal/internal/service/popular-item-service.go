package service

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
)

type PopularItemService interface {
	ProcessPopularItemService(request *model.PopularItemRequest) (*model.PopularItemResponse, error)
	io.Closer
}

type popularItemServiceImpl struct {
}

func NewPopularItemService(cluster *gocql.ClusterConfig) PopularItemService {
	return &popularItemServiceImpl{}
}

func (p *popularItemServiceImpl) ProcessPopularItemService(request *model.PopularItemRequest) (*model.PopularItemResponse, error) {
	panic("implement me")
}

func (p *popularItemServiceImpl) Close() error {
	panic("implement me")
}
