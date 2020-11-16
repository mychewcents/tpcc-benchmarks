package service

import (
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/common"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/internal/dao"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/internal/datamodel/table"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/internal/datamodel/view"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/model"
	"io"
)

type OrderStatusService interface {
	ProcessOrderStatusTransaction(request *model.OrderStatusRequest) (*model.OrderStatusResponse, error)
	io.Closer
}

type orderStatusServiceImpl struct {
	o  dao.OrderDao
	c  dao.CustomerDao
	ol dao.OrderLineDao
}

func NewOrderStatusService(cassandraSession *common.CassandraSession) OrderStatusService {
	return &orderStatusServiceImpl{
		o:  dao.NewOrderDao(cassandraSession),
		c:  dao.NewCustomerDao(cassandraSession),
		ol: dao.NewOrderLineDao(cassandraSession),
	}
}

func (o *orderStatusServiceImpl) ProcessOrderStatusTransaction(request *model.OrderStatusRequest) (*model.OrderStatusResponse, error) {
	chO := make(chan *view.OrderByCustomerView)
	chC := make(chan *table.CustomerTab)
	chOl := make(chan []*table.OrderLineTab)

	go o.o.GetCustomerLatestOrder(request.CWId, request.CDId, request.CId, chO)
	go o.c.GetCustomerByKey(request.CWId, request.CDId, request.CId, chC)

	ov := <-chO
	if ov == nil {
		return nil, nil
	}
	go o.ol.GetOrderLineListByKey(request.CWId, request.CDId, ov.OId, chOl)
	ct := <-chC
	olt := <-chOl

	olS := make([]*model.OrderLineStatus, 0)
	for _, ol := range olt {
		for sWId, quantity := range ol.OlWToQuantity {
			ol := &model.OrderLineStatus{
				OlIId:       ol.OlIId,
				OlSupplyWId: sWId,
				OlQuantity:  quantity,
				OlAmount:    ol.OlAmount * float32(quantity) / float32(ol.OlQuantity),
				OlDeliveryD: ov.OlDeliveryD,
			}
			olS = append(olS, ol)
		}

	}

	return &model.OrderStatusResponse{
		CName:               model.NameModelFromUDT(&ct.CName),
		CBalance:            ct.CBalance,
		OId:                 ov.OId,
		OEntryD:             ov.OEntryD,
		OCarrierId:          ov.OCarrierId,
		OrderLineStatusList: olS,
	}, nil
}

func (o *orderStatusServiceImpl) Close() error {
	panic("implement me")
}
