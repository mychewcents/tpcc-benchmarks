package service

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
)

type DeliveryService interface {
	ProcessDeliveryTransaction(request *model.DeliveryRequest) error
	io.Closer
}

type deliveryServiceImpl struct {
}

func NewDeliveryService(cluster *gocql.ClusterConfig) DeliveryService {
	return &deliveryServiceImpl{}
}

func (d *deliveryServiceImpl) ProcessDeliveryTransaction(request *model.DeliveryRequest) error {
	panic("implement me")
}

func (d *deliveryServiceImpl) Close() error {
	panic("implement me")
}
