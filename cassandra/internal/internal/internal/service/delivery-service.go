package service

import (
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
)

type DeliveryService interface {
	ProcessDeliveryTransaction(request *model.DeliveryRequest) error
	io.Closer
}

type deliveryServiceImpl struct {
}

func NewDeliveryService(cassandraSession *common.CassandraSession) DeliveryService {
	return &deliveryServiceImpl{}
}

func (d *deliveryServiceImpl) ProcessDeliveryTransaction(request *model.DeliveryRequest) error {
	panic("implement me")
}

func (d *deliveryServiceImpl) Close() error {
	panic("implement me")
}
