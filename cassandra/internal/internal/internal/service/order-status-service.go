package service

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
)

type OrderStatusService interface {
	ProcessOrderStatusTransaction(request *model.OrderStatusRequest) (*model.OrderStatusResponse, error)
	io.Closer
}

type orderStatusServiceImpl struct {
}

func NewOrderStatusService(cluster *gocql.ClusterConfig) OrderStatusService {
	return &orderStatusServiceImpl{}
}

func (o *orderStatusServiceImpl) ProcessOrderStatusTransaction(request *model.OrderStatusRequest) (*model.OrderStatusResponse, error) {
	panic("implement me")
}

func (o *orderStatusServiceImpl) Close() error {
	panic("implement me")
}
