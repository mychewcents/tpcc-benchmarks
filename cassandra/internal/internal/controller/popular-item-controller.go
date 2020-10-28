package controller

import (
	"fmt"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
	"strconv"
)

type PopularItemController interface {
	handler.TransactionHandler
}

type popularItemControllerImpl struct {
	s service.PopularItemService
}

func NewPopularItemController(cassandraSession *common.CassandraSession) PopularItemController {
	return &popularItemControllerImpl{
		s: service.NewPopularItemService(cassandraSession),
	}
}

func (p *popularItemControllerImpl) HandleTransaction(cmd []string) {
	request := makePopularItemRequest(cmd)
	response, _ := p.s.ProcessPopularItemService(request)
	printPopularItemResponse(response)
}

func makePopularItemRequest(cmd []string) *model.PopularItemRequest {
	wId, _ := strconv.Atoi(cmd[1])
	dId, _ := strconv.Atoi(cmd[2])
	l, _ := strconv.Atoi(cmd[3])

	return &model.PopularItemRequest{
		WId:            wId,
		DId:            dId,
		NoOfLastOrders: l,
	}
}

func printPopularItemResponse(r *model.PopularItemResponse) {
	fmt.Println("*********************** Popular Item Transaction Output ***********************")
	fmt.Printf("1. District Identifier - WId:%v DId:%v\n", r.WId, r.DId)
	fmt.Printf("2. Number of last orders to be examined: %v\n", r.NoOfLastOrders)
	fmt.Printf("3. For each order examined:\n")
	for _, o := range r.OrderItemInfoList {
		fmt.Printf("\ta. Order Number: %v, Order Entry Date: %v\n", o.OId, o.OEntryD)
		fmt.Printf("\tb. Name of customer: %+v\n", o.CName)
		fmt.Printf("\tc. Popular Items:\n")
		for _, p := range o.PopularItemInfoList {
			fmt.Printf("\t\ti. Item Name:%v\n", p.IName)
			fmt.Printf("\t\tii. Quantity Ordered:%v\n", p.OlQuantity)
		}
	}
	fmt.Printf("4. The percentage of examined orders that contain each popular item:\n")
	fmt.Println()
	for _, p := range r.PopularItemStatList {
		fmt.Printf("\ti. Item Name:%v\n", p.IName)
		fmt.Printf("\tii. The percentage of orders that contain the popular item:%.2f\n", p.OrderPercentage)
	}
	fmt.Println()
}

func (p *popularItemControllerImpl) Close() error {
	panic("implement me")
}
