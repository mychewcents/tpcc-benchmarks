package service

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/dao"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/table"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/udt"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
	"time"
)

type NewOrderService interface {
	ProcessNewOrderTransaction(request *model.NewOrderRequest) (*model.NewOrderResponse, error)
	io.Closer
}

type newOrderServiceImpl struct {
	c dao.CustomerDao
	o dao.
}

func NewNewOrderService(cluster *gocql.ClusterConfig) NewOrderService {
	return &newOrderServiceImpl{}
}

func (n *newOrderServiceImpl) ProcessNewOrderTransaction(request *model.NewOrderRequest) (*model.NewOrderResponse, error) {

}

func (n *newOrderServiceImpl) Close() error {
	panic("implement me")
}
