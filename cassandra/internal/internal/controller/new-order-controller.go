package controller

import (
	"bufio"
	"fmt"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/common"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/handler"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/model"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/service"
	"strconv"
	"strings"
)

type NewOrderController interface {
	handler.TransactionHandler
}

type newOrderControllerImpl struct {
	s service.NewOrderService
	r *bufio.Reader
}

func NewNewOrderTransactionController(cassandraSession *common.CassandraSession, reader *bufio.Reader) NewOrderController {
	return &newOrderControllerImpl{
		s: service.NewNewOrderService(cassandraSession),
		r: reader,
	}
}

func (n *newOrderControllerImpl) HandleTransaction(cmd []string) {
	request := makeNewOrderRequest(cmd, n.r)
	resp, _ := n.s.ProcessNewOrderTransaction(request)
	printNewOrderResponse(resp)
}

func makeNewOrderRequest(cmd []string, r *bufio.Reader) *model.NewOrderRequest {
	cId, _ := strconv.Atoi(cmd[1])
	wId, _ := strconv.Atoi(cmd[2])
	dId, _ := strconv.Atoi(cmd[3])
	m, _ := strconv.Atoi(cmd[4])

	return &model.NewOrderRequest{
		WId:              wId,
		DId:              dId,
		CId:              cId,
		NewOrderLineList: makeNewOrderLineList(m, r),
	}
}

func makeNewOrderLineList(m int, r *bufio.Reader) []*model.NewOrderLine {
	newOrderLineList := make([]*model.NewOrderLine, m)

	for i := 0; i < m; i++ {
		text, _ := r.ReadString('\n')
		orderLineSplit := strings.Split(strings.Trim(text, "\n"), ",")

		olIId, _ := strconv.Atoi(orderLineSplit[0])
		olSupplyWId, _ := strconv.Atoi(orderLineSplit[1])
		olQuantity, _ := strconv.Atoi(orderLineSplit[2])

		newOrderLine := &model.NewOrderLine{
			OlIId:       olIId,
			OlSupplyWId: olSupplyWId,
			OlQuantity:  olQuantity,
		}

		newOrderLineList[i] = newOrderLine
	}

	return newOrderLineList
}

func printNewOrderResponse(r *model.NewOrderResponse) {
	fmt.Println("*********************** New Order Transaction Output ***********************")
	fmt.Printf("1. Customer identifier - (WId:%v DId:%v CId:%v), ", r.WId, r.DId, r.CId)
	fmt.Printf("Customer Lastname:%v, CCredit:%v, CDiscount:%v.\n", r.CLast, r.CCredit, r.CDiscount)
	fmt.Printf("2. Warehouse tax rate: %v, District tax rate: %v.\n", r.WTax, r.DTax)
	fmt.Printf("3. Order number: %v, Entry Date: %v.\n", r.OId, r.OEntryD)
	fmt.Printf("4. Number of items NoOfItems: %v, Total amount for order: %.2f.\n", len(r.NewOrderLineInfoList), r.TotalAmount)
	fmt.Printf("5. For each ordered item: \n")
	for _, info := range r.NewOrderLineInfoList {
		fmt.Println()
		fmt.Printf("\ta. Item Number: %v\n", info.IId)
		fmt.Printf("\tb. Item Name: %v\n", info.IName)
		fmt.Printf("\tc. Supplier Warehouse: %v\n", info.SupplierWId)
		fmt.Printf("\td. Quantity: %v\n", info.Quantity)
		fmt.Printf("\te. OlAmount: %.2f\n", info.OlAmount)
		fmt.Printf("\tf. SQuantity: %v\n", info.SQuantity)
	}
	fmt.Println()
}

func (n *newOrderControllerImpl) Close() error {
	panic("implement me")
}
