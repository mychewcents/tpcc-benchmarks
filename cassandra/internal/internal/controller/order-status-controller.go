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
	n.s.ProcessOrderStatusTransaction(request)
	//printOrderStatusResponse(response)
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
	if r == nil {
		return
	}
	fmt.Println("*********************** Order Status Transaction Output ***********************")
	fmt.Printf("1. Customer's name: %v %v %v, Customer Balance: %.2f.\n", r.CName.FirstName, r.CName.MiddleName, r.CName.LastName, r.CBalance)
	fmt.Printf("2. For the customer's last order OId:%v OEntryD:%v OCarrierId:%v\n", r.OId, r.OEntryD, r.OCarrierId)
	fmt.Printf("3. For each item in the customer's last order:\n")
	for _, o := range r.OrderLineStatusList {
		fmt.Println()
		fmt.Printf("\ta. Item Number: %v\n", o.OlIId)
		fmt.Printf("\tb. Supplying warehouse number: %v\n", o.OlSupplyWId)
		fmt.Printf("\tc. Quantity ordered: %v\n", o.OlQuantity)
		fmt.Printf("\td. Total price for ordered item: %.2f\n", o.OlAmount)
		if o.OlDeliveryD.IsZero() {
			fmt.Printf("\te. Data and time of delivery: null\n")
		} else {
			fmt.Printf("\te. Data and time of delivery: %v\n", o.OlDeliveryD)
		}
	}
	fmt.Println()
}

func (n *orderStatusControllerImpl) Close() error {
	panic("implement me")
}
