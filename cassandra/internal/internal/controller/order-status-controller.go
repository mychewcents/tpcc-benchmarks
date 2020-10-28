package controller

import (
	"fmt"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
	"strconv"
)

type OrderStatusController interface {
	handler.TransactionHandler
}

type orderStatusControllerImpl struct {
	s service.OrderStatusService
}

func NewOrderStatusTransactionController(cassandraSession *common.CassandraSession) OrderStatusController {
	return &orderStatusControllerImpl{
		s: service.NewOrderStatusService(cassandraSession),
	}
}

func (n *orderStatusControllerImpl) HandleTransaction(cmd []string) {
	request := makeOrderStatusRequest(cmd)
	response, _ := n.s.ProcessOrderStatusTransaction(request)
	printOrderStatusResponse(response)
}

func makeOrderStatusRequest(cmd []string) *model.OrderStatusRequest {
	cWId, _ := strconv.Atoi(cmd[1])
	cDId, _ := strconv.Atoi(cmd[2])
	cId, _ := strconv.Atoi(cmd[3])

	return &model.OrderStatusRequest{
		CWId: cWId,
		CDId: cDId,
		CId:  cId,
	}
}

func printOrderStatusResponse(r *model.OrderStatusResponse) {
	//fmt.Println(r)
	fmt.Sprintf("1.Customer's name:%v\n", model.NameModelFromUDT(&ct.CName))
    fmt.Sprintf("2.For the customer's last order OId:%v\n OEntryD:%v\n OCarrierId:%v\n", ov.OId, ov.OEntryD, ov.OCarrierId)
    fmt.Sprintf("3.For each item in the customer's last order:%v\n", olS)
}

func (n *orderStatusControllerImpl) Close() error {
	panic("implement me")
}
