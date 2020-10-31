package service

import (
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/dao"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
)

type DatabaseStateService interface {
	GetDatabaseState() (*model.DatabaseStateResponse, error)
	io.Closer
}

type databaseStateServiceImpl struct {
	d dao.DatabaseStateDao
}

func NewDatabaseStateService(cassandraSession *common.CassandraSession) DatabaseStateService {
	return &databaseStateServiceImpl{d: dao.NewDatabaseStateDao(cassandraSession)}
}

func (d *databaseStateServiceImpl) GetDatabaseState() (*model.DatabaseStateResponse, error) {
	response := &model.DatabaseStateResponse{}

	response.SumWYTD = d.d.GetWarehouseState()
	response.SumDYTD = d.d.GetDistrictState()
	response.SumCBalance, response.SumCYTDPayment, response.SumCPaymentCnt, response.SumCDeliveryCnt = d.d.GetCustomerState()
	response.MaxOId, response.SumOOlCnt = d.d.GetOrderState()
	response.SumOlAmount, response.SumOlQuantity = d.d.GetOrderLineState()
	response.SumSQuantity, response.SumSYTD, response.SumSOrderCnt, response.SumSRemoteCnt = d.d.GetStockState()

	return response, nil
}

func (d *databaseStateServiceImpl) Close() error {
	panic("implement me")
}
